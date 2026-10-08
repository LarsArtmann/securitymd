package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
	finding "github.com/larsartmann/go-finding"
	"github.com/spf13/cobra"
)

const (
	outputFormatJSON  = "json"
	outputFormatSarif = "sarif"
	outputFormatText  = "text"
)

// errPolicyFindings is returned when error-severity findings remain: exit
// code 1 means "findings", not "crash".
var errPolicyFindings = errors.New("one or more policies failed validation")

var (
	outputFormat    string
	minimumSeverity string
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate security policy",
		Long: `Validate the repository's SECURITY.md for completeness and compliance.

Checks the candidate locations (SECURITY.md, .github/SECURITY.md,
docs/SECURITY.md); a missing file is itself an error finding. Exits 1 when
error-severity findings remain.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			file, _ := cmd.Flags().GetString("file")
			if file != "" {
				return validateFile(cmd, file)
			}

			return validateRepo(cmd)
		},
	}

	cmd.Flags().String("file", "", "Validate a specific policy file instead of discovering it")
	cmd.Flags().StringVar(&outputFormat, "format", outputFormatText, "Output format (text, json, sarif)")
	cmd.Flags().
		StringVar(&minimumSeverity, "severity", "info", "Minimum severity to report (info, warning, error, critical)")

	return cmd
}

func validateFile(cmd *cobra.Command, filename string) error {
	findings, err := policy.Validate(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	report := finding.NewReport(finding.ToolInfo{Name: string(policy.ToolName), Version: version})
	report.AddFindings(findings)
	report.ComputeSummary()

	return outputFindings(cmd, report)
}

func validateRepo(cmd *cobra.Command) error {
	report, err := policy.Report(cmd.Context())
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	return outputFindings(cmd, report)
}

func outputFindings(cmd *cobra.Command, report *finding.Report) error {
	findings := report.FindingsSnapshot()

	switch outputFormat {
	case outputFormatJSON:
		if err := report.WriteJSON(os.Stdout); err != nil {
			return fmt.Errorf("write JSON report: %w", err)
		}
	case outputFormatSarif:
		if err := report.WriteSARIFWithOpts(
			cmd.Context(),
			os.Stdout,
			finding.WithMinSeverity(parseSeverity(minimumSeverity)),
		); err != nil {
			return fmt.Errorf("write SARIF report: %w", err)
		}
	default:
		printFindings(findings)
	}

	if hasErrors(findings) {
		return errPolicyFindings
	}

	return nil
}

func printFindings(findings []finding.Finding) {
	if len(findings) == 0 {
		color.Green("✅ SECURITY.md - All checks passed")

		return
	}

	for _, f := range findings {
		location := string(f.Position.File)
		if f.Position.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, f.Position.Line)
		}

		if f.Severity == finding.SeverityError {
			color.Red("❌ [%s] %s (%s)", f.Rule, f.Message, location)
		} else {
			color.Yellow("⚠️  [%s] %s (%s)", f.Rule, f.Message, location)
		}

		if f.Suggestion != "" {
			fmt.Printf("    💡 %s\n", f.Suggestion)
		}
	}
}

func hasErrors(findings []finding.Finding) bool {
	return len(finding.Filter(findings, finding.BySeverity(finding.SeverityError))) > 0
}

func parseSeverity(severity string) finding.Severity {
	parsed := finding.Severity(severity)
	if parsed.IsValid() {
		return parsed
	}

	return finding.SeverityInfo
}
