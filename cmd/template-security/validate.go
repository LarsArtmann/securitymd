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
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate security policies",
		Long: `Validate security policies for completeness and compliance.
Checks SECURITY.md files for required sections and content quality.`,
		RunE: runValidate,
	}

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
			return outputJSONError(filename, err)
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

func printJSONResults(results []*internal.SecurityValidationResult) error {
	data, err := json.MarshalIndent(map[string]interface{}{
		"valid": allResultsValid(results),
		"files": results,
		"summary": map[string]int{
			"total": len(results),
			"valid": countValid(results),
		},
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))

	if !allResultsValid(results) {
		return errors.New("validation failed")
	}

	return nil
}

func allResultsValid(results []*internal.SecurityValidationResult) bool {
	for _, r := range results {
		if !r.Valid {
			return false
		}
	}
	return true
}

func countValid(results []*internal.SecurityValidationResult) int {
	count := 0
	for _, r := range results {
		if r.Valid {
			count++
		}
	}
	return count
}

func outputJSONError(filename string, err error) error {
	data, jErr := json.MarshalIndent(map[string]interface{}{
		"error":   "file not found",
		"file":    filename,
		"details": err.Error(),
	}, "", "  ")
	if jErr != nil {
		return jErr
	}
	fmt.Println(string(data))
	return err
}
