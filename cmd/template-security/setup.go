package main

import (
	"fmt"

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
	// Set defaults if not provided
	if policyType == "" {
		policyType = "github"
	}
	
	// Validate policy type
	if policyType != "github" && policyType != "enterprise" {
		return fmt.Errorf("invalid policy type: %s (supported: github, enterprise)", policyType)
	}

	// If no organization provided, try to detect it
	if orgName == "" {
		detector := internal.NewProjectDetector()
		orgName = detector.DetectProjectName()
		if orgName == "" {
			orgName = "YourOrganization"
		}
	}

	// If no email provided, try to detect domain
	if contactEmail == "" {
		detector := internal.NewProjectDetector()
		domain := detector.DetectDomain()
		if domain == "" {
			domain = "yourcompany.com"
		}
		contactEmail = "security@" + domain
	}

	color.Cyan("🔒 Generating SECURITY.md for: %s", orgName)
	
	// Create security tool and generate policy
	securityTool := internal.NewSecurityTool()
	config := internal.PolicyConfig{
		Type:         internal.PolicyType(policyType),
		Organization: orgName,
		ContactEmail: contactEmail,
		OutputDir:    outputDir,
		Variables:    make(map[string]string),
	}

	return securityTool.GeneratePolicy(cmd.Context(), config)
}
