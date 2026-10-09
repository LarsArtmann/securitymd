package main

import (
	"fmt"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
	finding "github.com/larsartmann/go-finding"
	"github.com/spf13/cobra"
)

var statusLocation string

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show security policy status",
		Long:  `Show where the repository's SECURITY.md lives and how compliant it is.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatus(cmd)
		},
	}

	cmd.Flags().StringVar(&statusLocation, "location", policy.LocationRoot,
		"Preferred canonical policy location (root, .github, docs)")

	return cmd
}

func runStatus(cmd *cobra.Command) error {
	candidates, err := policy.OrderedCandidates(statusLocation)
	if err != nil {
		return fmt.Errorf("invalid --location: %w", err)
	}

	report, err := policy.ReportIn(cmd.Context(), candidates)
	if err != nil {
		return fmt.Errorf("status failed: %w", err)
	}

	color.White("📋 Security Policy Status")
	color.White("========================")

	findings := report.FindingsSnapshot()

	if len(findings) == 0 {
		color.Green("✅ SECURITY.md found and compliant")

		return nil
	}

	errors := finding.Filter(findings, finding.BySeverityAtLeast(finding.SeverityError))
	warnings := finding.Filter(findings, finding.BySeverity(finding.SeverityWarning))

	color.Red("❌ %d error(s), %d warning(s)", len(errors), len(warnings))

	for _, f := range findings {
		if f.Severity.GreaterThanOrEqual(finding.SeverityError) {
			color.Red("  • [%s] %s", f.Rule, f.Message)
		} else {
			color.Yellow("  • [%s] %s", f.Rule, f.Message)
		}
	}

	fmt.Println()
	color.Cyan("Next steps:")
	color.Cyan("  1. securitymd setup    — generate a policy if missing")
	color.Cyan("  2. securitymd validate — full compliance check")

	return nil
}
