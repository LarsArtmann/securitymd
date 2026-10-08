package main

import (
	"fmt"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	setupOrganization string
	setupRepository   string
	setupEmail        string
	setupDirectory    string
	setupDryRun       bool
)

func newSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Generate SECURITY.md",
		Long: `Generate a SECURITY.md from the embedded canonical template.

Organization and repository default to the git remote's GitHub coordinates;
the contact email is optional — without it the policy points to GitHub
private vulnerability reporting. An existing policy is never overwritten.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := policy.Generate(cmd.Context(), policy.GenerateOptions{
				Directory:    setupDirectory,
				Organization: setupOrganization,
				Repository:   setupRepository,
				ContactEmail: setupEmail,
				DryRun:       setupDryRun,
			})
			if err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}

			if result.Wrote {
				color.Green("✅ %s", result.Description)
				color.Cyan("   Run 'securitymd validate' to confirm compliance")

				return nil
			}

			if setupDryRun {
				color.Cyan("🔍 %s", result.Description)

				return nil
			}

			color.Yellow("⚠️  %s", result.Description)

			return nil
		},
	}

	cmd.Flags().StringVar(&setupOrganization, "organization", "", "GitHub organization/owner (default: git remote)")
	cmd.Flags().StringVar(&setupRepository, "repository", "", "Repository name (default: git remote)")
	cmd.Flags().StringVar(&setupEmail, "email", "", "Security contact email (optional)")
	cmd.Flags().StringVar(&setupDirectory, "directory", "", "Target repository directory (default: current)")
	cmd.Flags().BoolVar(&setupDryRun, "dry-run", false, "Render without writing")

	return cmd
}
