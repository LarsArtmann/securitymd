package main

import (
	"errors"
	"fmt"
	"maps"

	"github.com/LarsArtmann/template-SECURITY/internal"
	"github.com/LarsArtmann/template-SECURITY/internal/types"
	"github.com/spf13/cobra"
)

// errInvalidPolicyType indicates an invalid policy type was provided.
var errInvalidPolicyType = errors.New("invalid policy type")

// ValidatePolicyType validates the policy type and returns an error if invalid.
func ValidatePolicyType(policyType string) error {
	if policyType != "github" && policyType != "enterprise" {
		return fmt.Errorf(
			"%w: %s (supported: github, enterprise)",
			errInvalidPolicyType,
			policyType,
		)
	}

	return nil
}

var (
	// orgName is a command-line flag for the organization name.
	orgName string
	// contactEmail is a command-line flag for the contact email.
	contactEmail string
	// policyType is a command-line flag for the policy type.
	policyType string
	// outputDir is a command-line flag for the output directory.
	outputDir string
	// quickMode is a command-line flag for quick setup mode.
	quickMode bool
)

func newSetupCmd() *cobra.Command {
	cmd := newCommand(
		"setup",
		"Generate SECURITY.md template",
		"Generate a SECURITY.md template for your project.\nSupports basic GitHub and enterprise security policies.",
		nil,
	)

	// Command line flags
	cmd.Flags().StringVarP(&policyType, "type", "t", "github", "Policy type (github, enterprise)")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	cmd.Flags().StringVarP(&contactEmail, "email", "e", "", "Security contact email")
	cmd.Flags().StringVar(&outputDir, "output", ".", "Output directory")
	cmd.Flags().BoolVar(&quickMode, "quick", false, "Quick setup with default values")

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		// Check for config file first
		securityTool := internal.NewSecurityTool()
		configFile := securityTool.FindConfigFile()

		// Load config if found
		var (
			userConfig *internal.Config
			err        error
		)

		if configFile != "" {
			printInfof("📄 Loading configuration from: %s", configFile)

			userConfig, err = securityTool.LoadConfig(configFile)
			if err != nil {
				printWarningf("⚠️  Warning: Failed to load config file: %v", err)
			}
		}

		// Set defaults if not provided
		if policyType == "" {
			policyType = "github"
			if userConfig != nil {
				policyType = string(userConfig.Type)
			}
		}

		// Validate policy type
		if err := ValidatePolicyType(policyType); err != nil {
			return err
		}

		// If no organization provided, try to detect it or use config
		if orgName == "" {
			if userConfig != nil && userConfig.Organization != "" {
				orgName = userConfig.Organization
			} else {
				detector := internal.NewProjectDetector()

				orgName = detector.DetectProjectName()
				if orgName == "" {
					orgName = "YourOrganization"
				}
			}
		}

		// If no email provided, try to detect domain or use config
		if contactEmail == "" {
			if userConfig != nil && userConfig.ContactEmail != "" {
				contactEmail = userConfig.ContactEmail
			} else {
				detector := internal.NewProjectDetector()

				domain := detector.DetectDomain()
				if domain == "" {
					domain = "yourcompany.com"
				}

				contactEmail = "security@" + domain
			}
		}

		outputDirToUse := outputDir
		if outputDirToUse == "." && userConfig != nil && userConfig.OutputDir != "" {
			outputDirToUse = userConfig.OutputDir
		}

		printInfof("🔒 Generating SECURITY.md for: %s", orgName)

		// Prepare variables from config
		variables := make(map[string]string)
		if userConfig != nil {
			maps.Copy(variables, userConfig.Variables)
		}

		config := internal.PolicyConfig{
			Type:         types.PolicyType(policyType),
			Organization: orgName,
			ContactEmail: contactEmail,
			OutputDir:    outputDirToUse,
			Variables:    variables,
		}

		return securityTool.GeneratePolicy(cmd.Context(), config)
	}

	return cmd
}
