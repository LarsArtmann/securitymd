package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/LarsArtmann/template-SECURITY/internal"
	finding "github.com/larsartmann/go-finding"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	outputFormat    string
	minimumSeverity string
)

func newValidateCmd() *cobra.Command {
	cmd := newCommand(
		"validate",
		"Validate security policies",
		"Validate security policies for completeness and compliance.\nChecks SECURITY.md files for required sections and content quality.",
		runValidate,
	)

	cmd.Flags().String("file", "", "Validate specific policy file")
	cmd.Flags().StringVar(&outputFormat, "format", "text", "Output format (text, json, sarif)")
	cmd.Flags().StringVar(&minimumSeverity, "severity", "info", "Minimum severity to report (info, warning, error, critical)")

	return cmd
}

func runValidate(cmd *cobra.Command, _ []string) error {
	file, _ := cmd.Flags().GetString("file")

	validator := internal.NewSecurityValidator()

	if file != "" {
		return validateSpecificFile(file, validator)
	}

	return validateAllPolicies(validator)
}

func parseSeverity(s string) finding.Severity {
	sev := finding.Severity(s)
	if sev.IsValid() {
		return sev
	}

	return finding.SeverityInfo
}

func validateSpecificFile(filename string, validator *internal.SecurityValidator) error {
	if outputFormat != "json" && outputFormat != "sarif" {
		color.Cyan("🔍 Validating security policy: %s", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if outputFormat == "json" || outputFormat == "sarif" {
			fmt.Fprintf(os.Stderr, "file not found: %s\n", filename)

			return fmt.Errorf("file not found: %s", filename)
		}

		color.Red("❌ File not found: %s", filename)

		return fmt.Errorf("file not found: %s", filename)
	}

	report, err := validator.ValidateSECURITYMd(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return outputReport(report)
}

func validateAllPolicies(validator *internal.SecurityValidator) error {
	if outputFormat != "json" && outputFormat != "sarif" {
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
			if outputFormat != "json" && outputFormat != "sarif" {
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
		if outputFormat != "json" && outputFormat != "sarif" {
			color.Yellow("⚠️  No security policy files found to validate")
		}

		return nil
	}

	if outputFormat != "text" {
		merged := mergeReports(reports)

		return outputReport(merged)
	}

	validator.PrintResults(reports)

	printSummary(reports, overallValid)

	if !overallValid {
		return errors.New("one or more policies failed validation")
	}

	return nil
}

func mergeReports(reports []*finding.Report) *finding.Report {
	if len(reports) == 1 {
		return reports[0]
	}

	merged := finding.NewReport(finding.ToolInfo{Name: "template-security"})

	for _, r := range reports {
		merged.AddFindings(r.Findings)
	}

	merged.ComputeSummary()

	return merged
}

func outputReport(report *finding.Report) error {
	minSev := parseSeverity(minimumSeverity)

	switch outputFormat {
	case "json":
		return report.WriteJSON(os.Stdout)
	case "sarif":
		return report.WriteSARIFFiltered(os.Stdout, minSev)
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

	for _, f := range filtered {
		icon := "⚠️"
		if f.Severity == finding.SeverityError || f.Severity == finding.SeverityCritical {
			icon = "❌"
		}

		fmt.Printf("  %s %s\n", icon, f.Message)
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


