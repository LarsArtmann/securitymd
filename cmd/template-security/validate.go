package main

import (
	"encoding/json"
	"fmt"

	"github.com/LarsArtmann/template-SECURITY/v2/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate security policies",
		Long: `Validate security policies for completeness and compliance.
Checks SECURITY.md files, policy documents, and compliance requirements.`,
		RunE: runValidate,
	}

	cmd.Flags().String("file", "", "Validate specific policy file")
	cmd.Flags().Bool("json", false, "Output validation results in JSON format")

	return cmd
}

func runValidate(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if file != "" {
		return validateSpecificFile(file, jsonOutput)
	} else {
		return validateAllPolicies(jsonOutput)
	}
}

func validateSpecificFile(filename string, jsonOutput bool) error {
	color.Cyan("🔍 Validating security policy: %s", filename)

	validator := internal.NewPolicyValidator()
	result := validator.ValidateFile(filename)

	if jsonOutput {
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	}

	displayValidationResult(filename, result)

	if !result.Valid {
		return fmt.Errorf("policy validation failed with score: %.1f", result.Score)
	}

	return nil
}

func validateAllPolicies(jsonOutput bool) error {
	color.Cyan("✅ Validating all security policies...")

	validator := internal.NewPolicyValidator()
	policyFiles := []string{
		"SECURITY.md",
		"security-policy.md",
		"incident-response.md",
		"privacy-policy.md",
		"bug-bounty-policy.md",
	}

	var allResults map[string]internal.ValidationResult
	allResults = make(map[string]internal.ValidationResult)

	overallValid := true

	for _, filename := range policyFiles {
		result := validator.ValidateFile(filename)
		allResults[filename] = result

		if !result.Valid {
			overallValid = false
		}

		if !jsonOutput {
			displayValidationResult(filename, result)
		}
	}

	if jsonOutput {
		jsonData, _ := json.MarshalIndent(allResults, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	}

	color.White("\n📊 Validation Summary:")
	if overallValid {
		color.Green("✅ All policies passed validation")
	} else {
		color.Red("❌ Some policies failed validation")
	}

	if !overallValid {
		return fmt.Errorf("one or more policies failed validation")
	}

	return nil
}

func displayValidationResult(filename string, result internal.ValidationResult) {
	fmt.Printf("\n📋 %s\n", filename)

	if result.Valid {
		color.Green("✅ VALID")
	} else {
		color.Red("❌ INVALID")
	}

	color.Cyan("📊 Score: %.1f/100", result.Score)

	if len(result.Issues) > 0 {
		color.Yellow("⚠️  Issues found:")
		for _, issue := range result.Issues {
			switch issue.Type {
			case "error":
				color.Red("   ❌ %s", issue.Description)
				color.White("      💡 %s", issue.Suggestion)
			case "warning":
				color.Yellow("   ⚠️  %s", issue.Description)
				color.White("      💡 %s", issue.Suggestion)
			case "info":
				color.Blue("   ℹ️  %s", issue.Description)
				color.White("      💡 %s", issue.Suggestion)
			}
		}
	}

	if len(result.Recommendations) > 0 {
		color.Cyan("💡 Recommendations:")
		for _, rec := range result.Recommendations {
			fmt.Printf("   • %s\n", rec)
		}
	}
}
