package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/LarsArtmann/template-SECURITY/internal"
	finding "github.com/larsartmann/go-finding"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var outputFormat string

func newValidateCmd() *cobra.Command {
	cmd := newCommand(
		"validate",
		"Validate security policies",
		"Validate security policies for completeness and compliance.\nChecks SECURITY.md files for required sections and content quality.",
		runValidate,
	)

	cmd.Flags().String("file", "", "Validate specific policy file")
	cmd.Flags().StringVar(&outputFormat, "format", "text", "Output format (text, json)")

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

func validateSpecificFile(filename string, validator *internal.SecurityValidator) error {
	if outputFormat != "json" {
		color.Cyan("🔍 Validating security policy: %s", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if outputFormat == "json" {
			return printJSONError(filename, err)
		}

		color.Red("❌ File not found: %s", filename)

		return fmt.Errorf("file not found: %s", filename)
	}

	report, err := validator.ValidateSECURITYMd(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return printValidationResult(report, validator)
}

func validateAllPolicies(validator *internal.SecurityValidator) error {
	if outputFormat != "json" {
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
			if outputFormat != "json" {
				color.Red("❌ Failed to validate %s: %v", filename, err)
			}

			overallValid = false

			continue
		}

		reports = append(reports, report)

		if !reportIsValid(reports) {
			overallValid = false
		}
	}

	if len(reports) == 0 {
		if outputFormat != "json" {
			color.Yellow("⚠️  No security policy files found to validate")
		}

		return nil
	}

	if outputFormat == "json" {
		return printJSONReports(reports)
	}

	validator.PrintResults(reports)

	color.White("\n📊 Validation Summary:")

	if overallValid {
		color.Green("✅ All policies passed validation")
	} else {
		color.Red("❌ Some policies failed validation")
	}

	if !overallValid {
		return errors.New("one or more policies failed validation")
	}

	return nil
}

func reportIsValid(reports []*finding.Report) bool {
	for _, r := range reports {
		for _, f := range r.Findings {
			if f.Severity == finding.SeverityError {
				return false
			}
		}
	}

	return true
}

func printValidationResult(
	report *finding.Report,
	validator *internal.SecurityValidator,
) error {
	if outputFormat == "json" {
		return printJSONReports([]*finding.Report{report})
	}

	validator.PrintResults([]*finding.Report{report})

	if !reportIsValid([]*finding.Report{report}) {
		return errors.New("policy validation failed")
	}

	return nil
}

func printJSONError(filename string, err error) error {
	return printJSON(map[string]any{
		"error":   "file not found",
		"file":    filename,
		"details": err.Error(),
	})
}

func printJSONReports(reports []*finding.Report) error {
	validCount := 0

	for _, r := range reports {
		if isValid(r) {
			validCount++
		}
	}

	return printJSON(map[string]any{
		"valid": validCount == len(reports),
		"files": reports,
		"summary": map[string]int{
			"total": len(reports),
			"valid": validCount,
		},
	})
}

func isValid(report *finding.Report) bool {
	for _, f := range report.Findings {
		if f.Severity == finding.SeverityError {
			return false
		}
	}

	return true
}

func printJSON(data map[string]any) error {
	jsonBytes, jErr := json.MarshalIndent(data, "", "  ")
	if jErr != nil {
		return fmt.Errorf("failed to marshal JSON: %w", jErr)
	}

	fmt.Println(string(jsonBytes))

	return nil
}
