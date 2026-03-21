package main

import (
	"fmt"
	"maps"

	"github.com/LarsArtmann/template-SECURITY/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	orgName      string
	contactEmail string
	policyType   string
	outputDir    string
	quickMode    bool
)

func newSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Generate SECURITY.md template",
		Long: `Generate a SECURITY.md template for your project.
Supports basic GitHub and enterprise security policies.`,
		RunE: runSetup,
	}

	// Command line flags
	cmd.Flags().StringVarP(&policyType, "type", "t", "github", "Policy type (github, enterprise)")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	cmd.Flags().StringVarP(&contactEmail, "email", "e", "", "Security contact email")
	cmd.Flags().StringVar(&outputDir, "output", ".", "Output directory")
	cmd.Flags().BoolVar(&quickMode, "quick", false, "Quick setup with default values")

	return cmd
}

func runSetup(cmd *cobra.Command, args []string) error {
	// Check for config file first
	securityTool := internal.NewSecurityTool()
	configFile := securityTool.FindConfigFile()

	// Load config if found
	var (
		userConfig *internal.Config
		err        error
	)

	if configFile != "" {
		color.Cyan("📄 Loading configuration from: %s", configFile)

		userConfig, err = securityTool.LoadConfig(configFile)
		if err != nil {
			color.Yellow("⚠️  Warning: Failed to load config file: %v", err)
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
	if policyType != "github" && policyType != "enterprise" {
		return fmt.Errorf("invalid policy type: %s (supported: github, enterprise)", policyType)
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

	color.Cyan("🔒 Generating SECURITY.md for: %s", orgName)

	// Prepare variables from config
	variables := make(map[string]string)
	if userConfig != nil {
		maps.Copy(variables, userConfig.Variables)
	}

	config := internal.PolicyConfig{
		Type:         internal.PolicyType(policyType),
		Organization: orgName,
		ContactEmail: contactEmail,
		OutputDir:    outputDirToUse,
		Variables:    variables,
	}

	return securityTool.GeneratePolicy(cmd.Context(), config)
}
