package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/LarsArtmann/template-SECURITY/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate security policies",
		Long: `Validate security policies for completeness and compliance.
Checks SECURITY.md files for required sections and content quality.`,
		RunE: runValidate,
	}

	cmd.Flags().String("file", "", "Validate specific policy file")

	return cmd
}

func runValidate(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")

	validator := internal.NewSecurityValidator()

	if file != "" {
		return validateSpecificFile(file, validator)
	} else {
		return validateAllPolicies(validator)
	}
}

func validateSpecificFile(filename string, validator *internal.SecurityValidator) error {
	color.Cyan("🔍 Validating security policy: %s", filename)

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		color.Red("❌ File not found: %s", filename)

		return fmt.Errorf("file not found: %s", filename)
	}

	result, err := validator.ValidateSECURITYMd(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Print result
	validator.PrintResults([]*internal.SecurityValidationResult{result})

	if !result.Valid {
		return errors.New("policy validation failed")
	}

	return nil
}

func validateAllPolicies(validator *internal.SecurityValidator) error {
	color.Cyan("✅ Validating all security policies...")

	policyFiles := []string{
		"SECURITY.md",
		"security-policy.md",
	}

	var results []*internal.SecurityValidationResult

	overallValid := true

	for _, filename := range policyFiles {
		// Check if file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			color.Yellow("⚠️  File not found: %s (skipping)", filename)

			continue
		}

		result, err := validator.ValidateSECURITYMd(filename)
		if err != nil {
			color.Red("❌ Failed to validate %s: %v", filename, err)

			overallValid = false

			continue
		}

		results = append(results, result)

		if !result.Valid {
			overallValid = false
		}
	}

	if len(results) == 0 {
		color.Yellow("⚠️  No security policy files found to validate")

		return nil
	}

	// Print all results
	validator.PrintResults(results)

	// Overall summary
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
