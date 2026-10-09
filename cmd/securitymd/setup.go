package main

import (
	"context"
	"fmt"
	"os"

	cmdguard "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	autoconfigure "github.com/larsartmann/linter-autoconfigure-sdk"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
)

// setupFlags carries everything `securitymd setup` accepts.
type setupFlags struct {
	Organization string `flag:"organization"                 help:"GitHub organization/owner (default: git remote, else interactive prompt)"`
	Repository   string `flag:"repository"                   help:"Repository name (default: git remote, else interactive prompt)"`
	Email        string `flag:"email"                        help:"Security contact email (optional)"`
	Directory    string `flag:"directory"                    help:"Target repository directory (default: current)"`
	Location     string `flag:"location" default:"root"      help:"Canonical policy location (root, .github, docs)"`
	DryRun       bool   `flag:"dry-run" default:"false"      help:"Render without writing"`
	Force        bool   `flag:"force" default:"false"        help:"Regenerate an existing policy in place (backs it up first)"`
}

func registerSetupCmd(cli *cmdguard.CLI[appConfig]) error {
	cmd, err := cmdguard.NewCommand(
		"setup",
		setupFlags{},
		runSetup,
		cmdguard.WithShort("Generate SECURITY.md"),
		cmdguard.WithLong(`Generate a SECURITY.md from the embedded canonical template.

Organization and repository default to the git remote's GitHub coordinates;
when no remote exists and stdin is a terminal, securitymd prompts for them.
The contact email is optional — without it the policy points to GitHub
private vulnerability reporting. An existing policy is never overwritten
unless --force is given, which regenerates in place after writing a
timestamped backup.`),
		cmdguard.WithNoArgs(),
	)
	if err != nil {
		return fmt.Errorf("building setup command: %w", err)
	}

	return cmdguard.AddCommand(cli, cmd)
}

func runSetup(ctx context.Context, _ *appConfig, flags setupFlags) error {
	organization, repository, err := resolveIdentity(ctx, flags)
	if err != nil {
		return err
	}

	result, err := policy.Generate(ctx, policy.GenerateOptions{
		Directory:    flags.Directory,
		Organization: organization,
		Repository:   repository,
		ContactEmail: flags.Email,
		Location:     flags.Location,
		DryRun:       flags.DryRun,
		Force:        flags.Force,
	})
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	return printSetupResult(result, flags.DryRun)
}

// resolveIdentity asks for the GitHub coordinates on an interactive terminal
// when neither flags nor a git remote can provide them.
func resolveIdentity(ctx context.Context, flags setupFlags) (string, string, error) {
	if flags.DryRun || flags.Organization != "" || flags.Repository != "" || !stdinIsInteractive() {
		return flags.Organization, flags.Repository, nil
	}

	dir := flags.Directory
	if dir == "" {
		dir = autoconfigure.WorkingDir(ctx)
	}

	if _, found := autoconfigure.FirstExisting(dir, policy.CandidateLocations...); found {
		return flags.Organization, flags.Repository, nil
	}

	if policy.DetectRepoIdentity(ctx, dir).IsComplete() {
		return flags.Organization, flags.Repository, nil
	}

	organization, repository, err := promptIdentity(os.Stdin, color.Output)
	if err != nil {
		return "", "", fmt.Errorf("interactive setup failed: %w", err)
	}

	return organization, repository, nil
}

func printSetupResult(result policy.GenerateResult, dryRun bool) error {
	if result.Wrote {
		color.Green("✅ %s", result.Description)

		if result.BackupPath != "" {
			color.Cyan("   Restore with: cp %s %s", result.BackupPath, result.Path)
		}

		color.Cyan("   Run 'securitymd validate' to confirm compliance")

		return nil
	}

	if dryRun {
		color.Cyan("🔍 %s", result.Description)

		return nil
	}

	color.Yellow("⚠️  %s", result.Description)

	return nil
}
