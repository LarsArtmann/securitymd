package policy

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
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
}

// GenerateResult describes what a Generate run did (or, under DryRun, would do).
type GenerateResult struct {
	// Path is the file that was written or would be written.
	Path string

	// Wrote reports that the file was actually written this run.
	Wrote bool

	// Description is the human-readable outcome for CLI and BuildFlow output.
	Description string
}

// Generate renders the embedded canonical template and writes SECURITY.md to
// the repository root — atomically, and ONLY when no policy file exists yet:
// securitymd never overwrites a human-authored policy. Skips (no error) when
// the identity cannot be derived, telling the caller what is missing.
func Generate(ctx context.Context, opts GenerateOptions) (GenerateResult, error) {
	dir := opts.Directory
	if dir == "" {
		dir = autoconfigure.WorkingDir(ctx)
	}

	if existing, found := autoconfigure.FirstExisting(dir, CandidateLocations...); found {
		return GenerateResult{
			Path:        existing,
			Description: existing + " already exists — securitymd never overwrites an existing policy",
		}, nil
	}

	identity := RepoIdentity{Organization: opts.Organization, Repository: opts.Repository}
	if !identity.IsComplete() {
		identity = DetectRepoIdentity(ctx, dir)
	}

	if !identity.IsComplete() {
		return GenerateResult{
			Path: filepath.Join(dir, CandidateLocations[0]),
			Description: "could not derive organization/repository from the git remote — " +
				"pass --organization and --repository (or add a git remote)",
		}, nil
	}

	content, err := renderForIdentity(ctx, dir, identity, opts.ContactEmail)
	if err != nil {
		return GenerateResult{}, fmt.Errorf("render SECURITY.md template for %s/%s: %w",
			identity.Organization, identity.Repository, err)
	}

	target := filepath.Join(dir, CandidateLocations[0])

	if opts.DryRun {
		return GenerateResult{
			Path:        target,
			Description: "dry-run: would create " + target,
		}, nil
	}

	if _, err := atomicwrite.WriteIfChanged(target, []byte(content)); err != nil {
		return GenerateResult{}, fmt.Errorf("write %s: %w", target, err)
	}

	return GenerateResult{
		Path:        target,
		Wrote:       true,
		Description: fmt.Sprintf("created %s for %s/%s", target, identity.Organization, identity.Repository),
	}, nil
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
