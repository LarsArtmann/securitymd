package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/LarsArtmann/template-SECURITY/internal"
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

func runValidate(cmd *cobra.Command, args []string) error {
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

	result, err := validator.ValidateSECURITYMd(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return printValidationResult(result, validator)
}

func validateAllPolicies(validator *internal.SecurityValidator) error {
	if outputFormat != "json" {
		color.Cyan("✅ Validating all security policies...")
	}

	policyFiles := []string{
		"SECURITY.md",
		"security-policy.md",
	}

	var results []*internal.SecurityValidationResult

	overallValid := true

	for _, filename := range policyFiles {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}

		result, err := validator.ValidateSECURITYMd(filename)
		if err != nil {
			if outputFormat != "json" {
				color.Red("❌ Failed to validate %s: %v", filename, err)
			}

			overallValid = false

			continue
		}

		results = append(results, result)

		if !result.Valid {
			overallValid = false
		}
	}

	if len(results) == 0 {
		if outputFormat != "json" {
			color.Yellow("⚠️  No security policy files found to validate")
		}

		return nil
	}

	if outputFormat == "json" {
		return printJSONResults(results)
	}

	validator.PrintResults(results)

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

func printValidationResult(
	result *internal.SecurityValidationResult,
	validator *internal.SecurityValidator,
) error {
	if outputFormat == "json" {
		return printJSONResults([]*internal.SecurityValidationResult{result})
	}

	validator.PrintResults([]*internal.SecurityValidationResult{result})

	if !result.Valid {
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

func printJSONResults(results []*internal.SecurityValidationResult) error {
	validCount, totalCount, allValid := countValidAndCheck(results)

	return printJSON(map[string]any{
		"valid": allValid,
		"files": results,
		"summary": map[string]int{
			"total": totalCount,
			"valid": validCount,
		},
	})
}

func printJSON(data map[string]any) error {
	jsonBytes, jErr := json.MarshalIndent(data, "", "  ")
	if jErr != nil {
		return fmt.Errorf("failed to marshal JSON: %w", jErr)
	}

	fmt.Println(string(jsonBytes))

	return nil
}

func countValidAndCheck(
	results []*internal.SecurityValidationResult,
) (validCount, totalCount int, allValid bool) {
	for _, r := range results {
		totalCount++

		if r.Valid {
			validCount++
		}
	}

	allValid = validCount == totalCount

	return validCount, totalCount, allValid
}

func allResultsValid(results []*internal.SecurityValidationResult) bool {
	_, _, allValid := countValidAndCheck(results)

	return allValid
}

func countValid(results []*internal.SecurityValidationResult) int {
	count, _, _ := countValidAndCheck(results)

	return count
}
