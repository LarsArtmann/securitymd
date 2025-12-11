package internal

import (
	"context"
	"fmt"
	"maps"

	"github.com/LarsArtmann/template-CLI/pkg/sdk/fileops"
	"github.com/LarsArtmann/template-CLI/pkg/sdk/vfs"
	"github.com/spf13/afero"
)

// SecurityTool represents the main security policy tool
type SecurityTool struct {
	templateManager   *vfs.TemplateManager
	fileSystem        vfs.FileSystem
	fileOps           *fileops.Service
	processor         *vfs.TemplateProcessor
	security          *vfs.TemplateSecurity
	detector          *ProjectDetector
	githubIntegration *GitHubIntegration
	variableDetector  *TemplateVariableDetector
}

// NewSecurityTool creates a new security tool instance
func NewSecurityTool() *SecurityTool {
	vfsImpl := vfs.NewOSVFS()
	fileOps := fileops.NewService(fileops.ServiceOptions{
		BaseDirectory: ".",
		ValidatePaths: true,
	})

	return &SecurityTool{
		templateManager:   vfs.NewTemplateManager(vfsImpl),
		fileSystem:        vfsImpl,
		fileOps:           fileOps,
		processor:         vfs.NewTemplateProcessor(),
		security:          vfs.NewTemplateSecurity(),
		detector:          NewProjectDetector(),
		githubIntegration: NewGitHubIntegration(),
		variableDetector:  NewTemplateVariableDetector(),
	}
}

// PolicyType represents different security policy types
type PolicyType string

const (
	PolicyTypeGitHub           PolicyType = "github"
	PolicyTypeEnterprise       PolicyType = "enterprise"
	PolicyTypeBugBounty        PolicyType = "bug-bounty"
	PolicyTypeIncidentResponse PolicyType = "incident-response"
	PolicyTypePrivacyPolicy    PolicyType = "privacy-policy"
)

// PolicyConfig holds configuration for policy generation
type PolicyConfig struct {
	Type         PolicyType
	Organization string
	ContactEmail string
	OutputDir    string
	Variables    map[string]string
}

// GeneratePolicy generates a security policy based on configuration
func (st *SecurityTool) GeneratePolicy(ctx context.Context, config PolicyConfig) error {
	// Load appropriate template
	templatePath := st.getTemplatePath(config.Type)
	templateFile, err := st.fileSystem.Open(templatePath)
	if err != nil {
		return fmt.Errorf("failed to open template %s: %w", templatePath, err)
	}
	defer templateFile.Close()

	// Read template content
	templateContent, err := afero.ReadAll(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template content: %w", err)
	}

	// Validate template security
	if err := st.security.ValidateTemplate(templatePath, templateContent); err != nil {
		return fmt.Errorf("template security validation failed: %w", err)
	}

	// Process template with variables
	processedContent := st.processTemplate(string(templateContent), config)

	// Determine output filename
	outputFile := st.getOutputFilename(config.Type, config.OutputDir)

	// Write output file using safe operations
	result := st.fileOps.WriteFile(outputFile, processedContent)
	if result.IsError() {
		return fmt.Errorf("failed to write policy file: %w", result.Error())
	}

	return nil
}

// getTemplatePath returns the path to the template file
func (st *SecurityTool) getTemplatePath(policyType PolicyType) string {
	switch policyType {
	case PolicyTypeGitHub:
		return "templates/github-security.md"
	case PolicyTypeEnterprise:
		return "templates/enterprise-policy.md"
	case PolicyTypeBugBounty:
		return "templates/bug-bounty-program.md"
	case PolicyTypeIncidentResponse:
		return "templates/incident-response.md"
	case PolicyTypePrivacyPolicy:
		return "templates/privacy-policy.md"
	default:
		return "templates/base-security.md"
	}
}

// getOutputFilename returns the output filename for the policy type
func (st *SecurityTool) getOutputFilename(policyType PolicyType, outputDir string) string {
	switch policyType {
	case PolicyTypeGitHub:
		return outputDir + "/SECURITY.md"
	case PolicyTypeEnterprise:
		return outputDir + "/security-policy.md"
	case PolicyTypeBugBounty:
		return outputDir + "/bug-bounty-policy.md"
	case PolicyTypeIncidentResponse:
		return outputDir + "/incident-response.md"
	case PolicyTypePrivacyPolicy:
		return outputDir + "/privacy-policy.md"
	default:
		return outputDir + "/security-policy.md"
	}
}

// processTemplate substitutes variables in the template
func (st *SecurityTool) processTemplate(template string, config PolicyConfig) string {
	// Start with GitHub-specific variables first (they should have priority)
	variables := map[string]string{
		// GitHub variables will be added here
	}

	// Add GitHub-specific variables if it's a GitHub repository
	githubInfo := st.githubIntegration.GetGitHubInfo()
	if githubInfo.IsGitHub {
		githubVariables := st.githubIntegration.GenerateGitHubTemplateVariables(githubInfo)

		// Add GitHub variables
		for key, value := range githubVariables {
			variables[key] = value
		}
	}

	// Add standard variables, but don't override GitHub-specific ones
	standardVariables := map[string]string{
		// Organization info (only if not already set by GitHub)
		"{{ORGANIZATION}}":        config.Organization,
		"{{CONTACT_EMAIL}}":       config.ContactEmail,
		"{{SECURITY_TEAM_EMAIL}}": config.ContactEmail,

		// Version info
		"{{LATEST_VERSION}}":            "v2.x",
		"{{PREVIOUS_VERSION}}":          "v1.x",
		"{{SUPPORT_END_DATE}}":          "2026-12-31",
		"{{PREVIOUS_SUPPORT_END_DATE}}": "2025-12-31",

		// Policy info
		"{{CURRENT_POLICY_VERSION}}":  "2.0",
		"{{PREVIOUS_POLICY_VERSION}}": "1.0",
		"{{LAST_UPDATED}}":            "2025-12-11",
		"{{PREVIOUS_UPDATED}}":        "2025-06-11",

		// Security resources
		"{{PGP_KEY_URL}}":             "https://{{ORGANIZATION}}.com/security/pgp",
		"{{BOUNTY_PROGRAM_URL}}":      "https://{{ORGANIZATION}}.com/security/bounty",
		"{{SECURITY_ADVISORIES_URL}}": "https://github.com/{{ORGANIZATION}}/security/advisories",
		"{{SECURITY_BLOG_URL}}":       "https://{{ORGANIZATION}}.com/blog/security",
		"{{SECURITY_DOCS_URL}}":       "https://{{ORGANIZATION}}.com/docs/security",
		"{{INCIDENT_RESPONSE_URL}}":   "https://{{ORGANIZATION}}.com/security/incident-response",

		// Legal resources
		"{{TERMS_URL}}":          "https://{{ORGANIZATION}}.com/terms",
		"{{PRIVACY_POLICY_URL}}": "https://{{ORGANIZATION}}.com/privacy",
		"{{LICENSE_URL}}":        "https://creativecommons.org/licenses/by-sa/4.0/",

		// Bounty rewards
		"{{CRITICAL_REWARD}}": "1000",
		"{{HIGH_REWARD}}":     "500",
		"{{MEDIUM_REWARD}}":   "200",
		"{{LOW_REWARD}}":      "50",

		// Default content
		"{{RESEARCHERS_LIST}}": "- Thank you to all security researchers who have helped us secure our products",
	}

	// Merge standard variables, but don't override GitHub-specific ones
	for key, value := range standardVariables {
		if _, exists := variables[key]; !exists {
			variables[key] = value
		}
	}

	// Add custom variables (highest priority)
	maps.Copy(variables, config.Variables)

	// Replace variables (multi-pass for nested variables)
	result := template
	maxPasses := 3
	for pass := 0; pass < maxPasses; pass++ {
		changed := false
		for placeholder, value := range variables {
			oldResult := result
			result = replaceAll(result, placeholder, value)
			if oldResult != result {
				changed = true
			}
		}
		if !changed {
			break // No more changes, exit early
		}
	}

	return result
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	// Simple implementation for now, could be optimized
	result := s
	for {
		before := result
		result = replaceString(result, old, new)
		if before == result {
			break
		}
	}
	return result
}

// replaceString replaces the first occurrence of old with new in s
func replaceString(s, old, new string) string {
	// Simple string replacement implementation
	if len(old) == 0 {
		return s
	}

	i := 0
	for i+len(old) <= len(s) {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
		i++
	}
	return s
}
