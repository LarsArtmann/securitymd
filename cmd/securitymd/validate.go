package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
var errPolicyFindings = errors.New("SECURITY.md validation failed: error-severity findings remain (see above)")

var (
	outputFormat     string
	minimumSeverity  string
	severitySpecs    []string
	validateLocation string
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
			overrides, err := parseSeverityFlag(severitySpecs)
			if err != nil {
				return err
			}

			file, _ := cmd.Flags().GetString("file")
			if file != "" {
				return validateFile(cmd, file, overrides)
			}

			return validateRepo(cmd, overrides)
		},
	}

	cmd.Flags().String("file", "", "Validate a specific policy file instead of discovering it")
	cmd.Flags().StringVar(&outputFormat, "format", outputFormatText, "Output format (text, json, sarif)")
	cmd.Flags().
		StringVar(&minimumSeverity, "severity", "info", "Minimum severity to report (info, warning, error, critical)")
	cmd.Flags().StringSliceVar(&severitySpecs, "set-severity", nil,
		"Override a rule's severity, rule=level (repeatable, e.g. --set-severity missing-file=warning)")
	cmd.Flags().StringVar(&validateLocation, "location", policy.LocationRoot,
		"Preferred canonical policy location (root, .github, docs)")

	return cmd
}

func parseSeverityFlag(specs []string) (policy.SeverityOverrides, error) {
	overrides, err := policy.ParseSeverityOverrides(strings.Join(specs, ","))
	if err != nil {
		return nil, fmt.Errorf("--set-severity: %w", err)
	}

	return overrides, nil
}

func validateFile(cmd *cobra.Command, filename string, overrides policy.SeverityOverrides) error {
	findings, err := policy.Validate(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return outputFindings(cmd, buildReport(findings, overrides))
}

func validateRepo(cmd *cobra.Command, overrides policy.SeverityOverrides) error {
	candidates, err := policy.OrderedCandidates(validateLocation)
	if err != nil {
		return fmt.Errorf("invalid --location: %w", err)
	}

	findings, err := policy.DetectIn(cmd.Context(), candidates)
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	return outputFindings(cmd, buildReport(findings, overrides))
}

// buildReport applies severity overrides and aggregates the findings; the
// summary therefore always reflects the severities the user actually sees.
func buildReport(findings []finding.Finding, overrides policy.SeverityOverrides) *finding.Report {
	report := finding.NewReportFromFindings(
		finding.ToolInfo{Name: string(policy.ToolName), Version: version},
		policy.ApplySeverityOverrides(findings, overrides))

	return report
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

	now := time.Now()

	for _, f := range findings {
		location := string(f.Position.File)
		if f.Position.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, f.Position.Line)
		}

		suppressed := f.IsSuppressedAt(now)

		switch {
		case suppressed:
			color.Cyan("🔇 [%s] %s (%s)\n    🔇 suppressed: %s", f.Rule, f.Message, location, f.Suppression.Reason)
		case f.Severity == finding.SeverityError:
			color.Red("❌ [%s] %s (%s)", f.Rule, f.Message, location)
		default:
			color.Yellow("⚠️  [%s] %s (%s)", f.Rule, f.Message, location)
		}

		if !suppressed && f.Suggestion != "" {
			fmt.Printf("    💡 %s\n", f.Suggestion)
		}
	}
}

// hasErrors gates the exit code on ACTIVE findings only: suppressed findings
// stay visible in the output but never fail CI — that is the escape hatch.
// The gate is threshold-based (error or critical): an escalated finding must
// never silently pass.
func hasErrors(findings []finding.Finding) bool {
	active := finding.Filter(findings, func(f finding.Finding) bool {
		return !f.IsSuppressedAt(time.Now())
	})

	return len(finding.Filter(active, finding.BySeverityAtLeast(finding.SeverityError))) > 0
}

func parseSeverity(severity string) finding.Severity {
	parsed := finding.Severity(severity)
	if parsed.IsValid() {
		return parsed
	}

	return finding.SeverityInfo
}
