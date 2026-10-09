package main

import (
	"fmt"
	"os"

	autoconfigure "github.com/larsartmann/linter-autoconfigure-sdk"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	setupOrganization string
	setupRepository   string
	setupEmail        string
	setupDirectory    string
	setupLocation     string
	setupDryRun       bool
	setupForce        bool
)

func newSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Generate SECURITY.md",
		Long: `Generate a SECURITY.md from the embedded canonical template.

Organization and repository default to the git remote's GitHub coordinates;
when no remote exists and stdin is a terminal, securitymd prompts for them.
The contact email is optional — without it the policy points to GitHub
private vulnerability reporting. An existing policy is never overwritten
unless --force is given, which regenerates in place after writing a
timestamped backup.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := promptWhenIdentityMissing(cmd); err != nil {
				return err
			}

			result, err := policy.Generate(cmd.Context(), policy.GenerateOptions{
				Directory:    setupDirectory,
				Organization: setupOrganization,
				Repository:   setupRepository,
				ContactEmail: setupEmail,
				Location:     setupLocation,
				DryRun:       setupDryRun,
				Force:        setupForce,
			})
			if err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}

			return printSetupResult(result)
		},
	}

	cmd.Flags().StringVar(&setupOrganization, "organization", "", "GitHub organization/owner (default: git remote, else interactive prompt)")
	cmd.Flags().StringVar(&setupRepository, "repository", "", "Repository name (default: git remote, else interactive prompt)")
	cmd.Flags().StringVar(&setupEmail, "email", "", "Security contact email (optional)")
	cmd.Flags().StringVar(&setupDirectory, "directory", "", "Target repository directory (default: current)")
	cmd.Flags().StringVar(&setupLocation, "location", policy.LocationRoot, "Canonical policy location (root, .github, docs)")
	cmd.Flags().BoolVar(&setupDryRun, "dry-run", false, "Render without writing")
	cmd.Flags().BoolVar(&setupForce, "force", false, "Regenerate an existing policy in place (backs it up first)")

	return cmd
}

// promptWhenIdentityMissing asks for the GitHub coordinates on an interactive
// terminal when neither flags nor a git remote can provide them.
func promptWhenIdentityMissing(cmd *cobra.Command) error {
	if setupDryRun || setupOrganization != "" || setupRepository != "" || !stdinIsInteractive() {
		return nil
	}

	dir := setupDirectory
	if dir == "" {
		dir = autoconfigure.WorkingDir(cmd.Context())
	}

	if policy.DetectRepoIdentity(cmd.Context(), dir).IsComplete() {
		return nil
	}

	org, repo, err := promptIdentity(os.Stdin, color.Output)
	if err != nil {
		return fmt.Errorf("interactive setup failed: %w", err)
	}

	setupOrganization, setupRepository = org, repo

	return nil
}

func printSetupResult(result policy.GenerateResult) error {
	if result.Wrote {
		color.Green("✅ %s", result.Description)

		if result.BackupPath != "" {
			color.Cyan("   Restore with: cp %s %s", result.BackupPath, result.Path)
		}

		color.Cyan("   Run 'securitymd validate' to confirm compliance")

		return nil
	}

	if setupDryRun {
		color.Cyan("🔍 %s", result.Description)

		return nil
	}

	color.Yellow("⚠️  %s", result.Description)

	return nil
}
