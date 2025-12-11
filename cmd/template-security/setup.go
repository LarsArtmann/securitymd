package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/LarsArtmann/template-SECURITY/v2/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	orgName      string
	contactEmail string
	interactive  bool
	policyType   string
	outputDir    string
	quickMode    bool
)

func newSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Interactive security policy wizard",
		Long: `Setup security policies interactively or via command line options.
Supports GitHub, enterprise, bug bounty, incident response, and privacy policies.`,
		RunE: runSetup,
	}

	// Command line flags
	cmd.Flags().StringVarP(&policyType, "type", "t", "", "Policy type (github, enterprise, bug-bounty, incident-response, privacy-policy)")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	cmd.Flags().StringVarP(&contactEmail, "email", "e", "", "Security contact email")
	cmd.Flags().StringVar(&outputDir, "output", ".", "Output directory")
	cmd.Flags().BoolVar(&interactive, "interactive", true, "Force interactive mode")
	cmd.Flags().BoolVar(&quickMode, "quick", false, "Quick setup with minimal prompts")

	return cmd
}

func runSetup(cmd *cobra.Command, args []string) error {
	color.Cyan("🔒 Security Policy Setup v%s", version)
	color.White("Let's secure your project with proper policies!")
	fmt.Println()

	// Initialize project detector
	detector := internal.NewProjectDetector()

	// If not interactive and missing required info, prompt for it
	if !cmd.Flags().Changed("interactive") && interactive {
		// Force interactive mode off if all flags are provided
		interactive = false
	}

	if interactive || !cmd.Flags().Changed("organization") {
		promptOrg(detector)
	}

	if interactive || !cmd.Flags().Changed("email") {
		promptEmail(detector)
	}

	if interactive || !cmd.Flags().Changed("type") {
		promptPolicyType()
	}

	// Validate required fields
	if policyType == "" {
		return fmt.Errorf("policy type is required")
	}
	if orgName == "" {
		return fmt.Errorf("organization name is required")
	}
	if contactEmail == "" {
		return fmt.Errorf("contact email is required")
	}

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

func promptOrg(detector *internal.ProjectDetector) {
	defaultName := detector.DetectProjectName()
	prompt := &survey.Input{
		Message: "Organization name",
		Default: defaultName,
	}
	survey.AskOne(prompt, &orgName, survey.WithValidator(survey.Required))
}

func promptEmail(detector *internal.ProjectDetector) {
	domain := detector.DetectDomain()
	defaultEmail := "security@" + domain
	prompt := &survey.Input{
		Message: "Security contact email",
		Default: defaultEmail,
	}
	survey.AskOne(prompt, &contactEmail, survey.WithValidator(survey.Required))
}

func promptPolicyType() {
	prompt := &survey.Select{
		Message: "Select policy type to generate:",
		Options: []string{
			"github - GitHub SECURITY.md",
			"enterprise - Enterprise security policy",
			"bug-bounty - Bug bounty program",
			"incident-response - Incident response plan",
			"privacy-policy - GDPR/CCPA privacy policy",
		},
	}
	var selection string
	survey.AskOne(prompt, &selection)
	
	// Extract policy type from selection (get text before " - ")
	parts := strings.Split(selection, " - ")
	if len(parts) > 0 {
		policyType = parts[0]
	}
}

// Remove the old functions as they're no longer needed
