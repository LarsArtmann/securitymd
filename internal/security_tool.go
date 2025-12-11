package internal

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/template-CLI/pkg/sdk/vfs"
	"github.com/LarsArtmann/template-CLI/pkg/sdk/fileops"
	"github.com/spf13/afero"
)

// SecurityTool represents the main security policy tool
type SecurityTool struct {
	templateManager *vfs.TemplateManager
	fileSystem      vfs.FileSystem
	fileOps        *fileops.Service
	processor       *vfs.TemplateProcessor
	security        *vfs.TemplateSecurity
}

// NewSecurityTool creates a new security tool instance
func NewSecurityTool() *SecurityTool {
	vfsImpl := vfs.NewOSVFS()
	fileOps := fileops.NewService(fileops.ServiceOptions{
		BaseDirectory: ".",
		ValidatePaths: true,
	})
	
	return &SecurityTool{
		templateManager: vfs.NewTemplateManager(vfsImpl),
		fileSystem:      vfsImpl,
		fileOps:        fileOps,
		processor:       vfs.NewTemplateProcessor(),
		security:        vfs.NewTemplateSecurity(),
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
	// Start with standard variables
	variables := map[string]string{
		"{{ORGANIZATION}}":  config.Organization,
		"{{CONTACT_EMAIL}}": config.ContactEmail,
		"{{DATE}}":          "2025-12-11", // TODO: Use current date
		"{{YEAR}}":          "2025",       // TODO: Use current year
	}

	// Add custom variables
	for k, v := range config.Variables {
		variables[k] = v
	}

	// Replace variables
	result := template
	for placeholder, value := range variables {
		result = replaceAll(result, placeholder, value)
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
