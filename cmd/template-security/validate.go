package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/LarsArtmann/template-SECURITY/internal"
	"github.com/fatih/color"
	finding "github.com/larsartmann/go-finding"
	"github.com/spf13/cobra"
)

const (
	// outputFormatJSON represents JSON output format.
	outputFormatJSON = "json"
	// outputFormatSarif represents SARIF output format.
	outputFormatSarif = "sarif"
	// outputFormatText represents text output format.
	outputFormatText = "text"
)

var (
	outputFormat    string
	minimumSeverity string
)

// errPolicyValidation is returned when one or more policies failed validation.
var errPolicyValidation = errors.New("one or more policies failed validation")

// errFileNotFound indicates that the policy file was not found.
var errFileNotFound = errors.New("policy file not found")

func newValidateCmd() *cobra.Command {
	cmd := newCommand(
		"validate",
		"Validate security policies",
		"Validate security policies for completeness and compliance.\nChecks SECURITY.md files for required sections and content quality.",
		nil,
	)

	cmd.Flags().String("file", "", "Validate specific policy file")
	cmd.Flags().StringVar(&outputFormat, "format", "text", "Output format (text, json, sarif)")
	cmd.Flags().
		StringVar(&minimumSeverity, "severity", "info", "Minimum severity to report (info, warning, error, critical)")

	cmd.RunE = func(_ *cobra.Command, _ []string) error {
		file, _ := cmd.Flags().GetString("file")

		validator := internal.NewSecurityValidator()

		if file != "" {
			return validateSpecificFile(file, validator)
		}

		return validateAllPolicies(validator)
	}

	return cmd
}

func parseSeverity(s string) finding.Severity {
	sev := finding.Severity(s)
	if sev.IsValid() {
		return sev
	}

	return finding.SeverityInfo
}

func validateSpecificFile(filename string, validator *internal.SecurityValidator) error {
	if outputFormat != outputFormatJSON && outputFormat != outputFormatSarif {
		color.Cyan("🔍 Validating security policy: %s", filename)
	}

	stat, err := os.Stat(filename)
	if err != nil {
		if outputFormat == outputFormatJSON || outputFormat == outputFormatSarif {
			if _, writeErr := fmt.Fprintf(
				os.Stderr,
				"file not found: %s\n",
				filename,
			); writeErr != nil {
				color.Red("Failed to write error: %v", writeErr)
			}

			return fmt.Errorf("%w: %s", errFileNotFound, filename)
		}

		color.Red("❌ File not found: %s", filename)

		return fmt.Errorf("%w: %s", errFileNotFound, filename)
	}

	_ = stat

	report, err := validator.ValidateSECURITYMd(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return outputReport(report)
}

func validateAllPolicies(validator *internal.SecurityValidator) error {
	if outputFormat != outputFormatJSON && outputFormat != outputFormatSarif {
		color.Cyan("✅ Validating all security policies...")
	}

	policyFiles := []string{
		"SECURITY.md",
		"security-policy.md",
	}

	var reports []*finding.Report

	overallValid := true

	for _, filename := range policyFiles {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}

		report, err := validator.ValidateSECURITYMd(filename)
		if err != nil {
			if outputFormat != outputFormatJSON && outputFormat != outputFormatSarif {
				color.Red("❌ Failed to validate %s: %v", filename, err)
			}

			overallValid = false

			continue
		}

		reports = append(reports, report)

		if !internal.ReportIsValid(report) {
			overallValid = false
		}
	}

	if len(reports) == 0 {
		if outputFormat != outputFormatJSON && outputFormat != outputFormatSarif {
			color.Yellow("⚠️  No security policy files found to validate")
		}

		return nil
	}

	if outputFormat != outputFormatText {
		merged := mergeReports(reports)

		return outputReport(merged)
	}

	validator.PrintResults(reports)

	printSummary(reports, overallValid)

	if !overallValid {
		return errPolicyValidation
	}

	return nil
}

func mergeReports(reports []*finding.Report) *finding.Report {
	if len(reports) == 1 {
		return reports[0]
	}

	merged := finding.NewReport(finding.ToolInfo{Name: "template-security", Version: "dev"})

	for _, r := range reports {
		merged.AddFindings(r.Findings)
	}

	merged.ComputeSummary()

	return merged
}

func outputReport(report *finding.Report) error {
	minSev := parseSeverity(minimumSeverity)

	switch outputFormat {
	case outputFormatJSON:
		err := report.WriteJSON(os.Stdout)
		if err != nil {
			return fmt.Errorf("failed to write JSON report: %w", err)
		}

		return nil
	case outputFormatSarif:
		err := report.WriteSARIFFiltered(os.Stdout, minSev)
		if err != nil {
			return fmt.Errorf("failed to write SARIF report: %w", err)
		}

		return nil
	default:
		return printTextReport(report)
	}
}

func printTextReport(report *finding.Report) error {
	minSev := parseSeverity(minimumSeverity)
	filtered := finding.Filter(report.Findings,
		finding.BySeverityAtLeast(minSev),
		finding.NotSuppressed,
	)

	for _, findingItem := range filtered {
		icon := "⚠️"
		if findingItem.Severity == finding.SeverityError ||
			findingItem.Severity == finding.SeverityCritical {
			icon = "❌"
		}

		_, err := fmt.Fprintf(os.Stdout, "  %s %s\n", icon, findingItem.Message)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	return nil
}

func printSummary(_ []*finding.Report, overallValid bool) {
	color.White("\n📊 Validation Summary:")

	if overallValid {
		color.Green("✅ All policies passed validation")
	} else {
		color.Red("❌ Some policies failed validation")
	}
}
