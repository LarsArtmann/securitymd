package policy

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	atomicwrite "github.com/larsartmann/go-atomic-write"
	autoconfigure "github.com/larsartmann/linter-autoconfigure-sdk"
)

//go:embed template.md
var templateMarkdown string

// GenerateOptions controls one policy generation run.
type GenerateOptions struct {
	// Directory is the repository root. Defaults to the working directory
	// from ctx via autoconfigure.WorkingDir.
	Directory string

	// Organization and Repository name the GitHub coordinates used for the
	// advisory link. Auto-detected from the git remote when empty.
	Organization string
	Repository   string

	// ContactEmail is the optional security contact; when empty the policy
	// points to GitHub private vulnerability reporting only.
	ContactEmail string

	// DryRun renders the policy but skips the write.
	DryRun bool

	// Location is the canonical write target: "root" (default), ".github",
	// or "docs". An empty value means "root".
	Location string

	// Force regenerates an existing policy in place after backing it up
	// (SECURITY.md.<timestamp>.bak next to it). Refuse-by-default stays:
	// without Force an existing policy is never touched.
	Force bool
}

// Policy locations accepted by GenerateOptions.Location and the CLI's
// --location flags.
const (
	LocationRoot   = "root"
	LocationGitHub = ".github"
	LocationDocs   = "docs"
)

// GenerateResult describes what a Generate run did (or, under DryRun, would do).
type GenerateResult struct {
	// Path is the file that was written or would be written.
	Path string

	// Wrote reports that the file was actually written this run.
	Wrote bool

	// BackupPath is the timestamped backup of the previous policy, set only
	// when an existing policy was Force-regenerated.
	BackupPath string

	// Description is the human-readable outcome for CLI and BuildFlow output.
	Description string
}

// Generate renders the embedded canonical template and writes SECURITY.md to
// the repository root — atomically, and ONLY when no policy file exists yet:
// securitymd never overwrites a human-authored policy. The Force option is
// the explicit escape hatch: it regenerates an existing policy in place and
// backs the previous content up first. Skips (no error) when the identity
// cannot be derived, telling the caller what is missing.
func Generate(ctx context.Context, opts GenerateOptions) (GenerateResult, error) {
	dir := opts.Directory
	if dir == "" {
		dir = autoconfigure.WorkingDir(ctx)
	}

	target, err := policyTarget(dir, opts.Location)
	if err != nil {
		return GenerateResult{}, err
	}

	existing, found := autoconfigure.FirstExisting(dir, CandidateLocations...)
	if found {
		if !opts.Force {
			return GenerateResult{
				Path:        existing,
				Description: existing + " already exists — securitymd never overwrites an existing policy",
			}, nil
		}

		if existing != target {
			return GenerateResult{
				Path: existing,
				Description: fmt.Sprintf(
					"%s already exists — regenerating at %s would create two policies; move the file first or drop --location",
					existing, target),
			}, nil
		}
	}

	identity := RepoIdentity{Organization: opts.Organization, Repository: opts.Repository}
	if !identity.IsComplete() {
		identity = DetectRepoIdentity(ctx, dir)
	}

	if !identity.IsComplete() {
		return GenerateResult{
			Path: target,
			Description: "could not derive organization/repository from the git remote — " +
				"pass --organization and --repository (or add a git remote)",
		}, nil
	}

	content, err := renderForIdentity(ctx, dir, identity, opts.ContactEmail)
	if err != nil {
		return GenerateResult{}, fmt.Errorf("render SECURITY.md template for %s/%s: %w",
			identity.Organization, identity.Repository, err)
	}

	if opts.DryRun {
		description := "dry-run: would create " + target
		if found && opts.Force {
			description = "dry-run: would back up and regenerate " + target
		}

		return GenerateResult{Path: target, Description: description}, nil
	}

	backupPath, err := backupExisting(found, existing)
	if err != nil {
		return GenerateResult{}, err
	}

	if _, err := atomicwrite.WriteIfChanged(target, []byte(content)); err != nil {
		return GenerateResult{}, fmt.Errorf("write %s: %w", target, err)
	}

	result := GenerateResult{
		Path:  target,
		Wrote: true,
		Description: fmt.Sprintf("created %s for %s/%s",
			target, identity.Organization, identity.Repository),
	}

	if backupPath != "" {
		result.BackupPath = backupPath
		result.Description = fmt.Sprintf("regenerated %s for %s/%s (previous policy backed up to %s)",
			target, identity.Organization, identity.Repository, backupPath)
	}

	return result, nil
}

// policyTarget maps a Location option to its write path under dir.
func policyTarget(dir, location string) (string, error) {
	switch location {
	case "", LocationRoot:
		return filepath.Join(dir, CandidateLocations[0]), nil
	case LocationGitHub:
		return filepath.Join(dir, ".github", CandidateLocations[0]), nil
	case LocationDocs:
		return filepath.Join(dir, "docs", CandidateLocations[0]), nil
	default:
		return "", fmt.Errorf("unknown policy location %q (want %s, %s, or %s)",
			location, LocationRoot, LocationGitHub, LocationDocs)
	}
}

// backupExisting copies the existing policy to a timestamped .bak file next
// to it. found=false (nothing to back up) is a no-op returning "".
func backupExisting(found bool, existing string) (string, error) {
	if !found {
		return "", nil
	}

	previous, err := os.ReadFile(existing)
	if err != nil {
		return "", fmt.Errorf("read %s for backup: %w", existing, err)
	}

	backupPath := existing + "." + time.Now().Format("20060102T150405") + ".bak"
	if _, err := atomicwrite.WriteIfChanged(backupPath, previous); err != nil {
		return "", fmt.Errorf("write backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

// templateData is the render input for template.md.
type templateData struct {
	Organization string
	Repository   string
	ContactEmail string
	VersionCell  string
	LastUpdated  string
}

// renderForIdentity renders the embedded template for a known repo identity.
// Shared by Generate (write path) and the missing-file finding preview.
func renderForIdentity(ctx context.Context, dir string, identity RepoIdentity, contactEmail string) (string, error) {
	return renderTemplate(templateData{
		Organization: identity.Organization,
		Repository:   identity.Repository,
		ContactEmail: contactEmail,
		VersionCell:  versionCell(ctx, dir),
		LastUpdated:  time.Now().Format("2006-01-02"),
	})
}

// versionCell fills the Supported Versions table: the latest git tag when
// available, an honest placeholder otherwise. The tag lookup is cached per
// directory (see lookupLatestTag) so a run's repeated renders cost one
// `git describe`.
func versionCell(ctx context.Context, dir string) string {
	if tag := lookupLatestTag(ctx, dir); tag != "" {
		return tag
	}

	return "Latest release"
}

func renderTemplate(data templateData) (string, error) {
	tmpl, err := template.New("security").Parse(templateMarkdown)
	if err != nil {
		return "", fmt.Errorf("parse embedded template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template for %s/%s: %w",
			data.Organization, data.Repository, err)
	}

	return buf.String(), nil
}
