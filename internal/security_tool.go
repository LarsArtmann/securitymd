package internal

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// PolicyType represents different types of security policies
type PolicyType string

const (
	PolicyTypeGitHub     PolicyType = "github"
	PolicyTypeEnterprise PolicyType = "enterprise"
)

// VersionStatus represents the status of a version
type VersionStatus string

const (
	StatusSupported  VersionStatus = "supported"
	StatusDeprecated VersionStatus = "deprecated"
	StatusEOL        VersionStatus = "end-of-life"
)

// ContactType represents different types of contact methods
type ContactType string

const (
	ContactTypeEmail ContactType = "email"
	ContactTypeWeb   ContactType = "web"
	ContactTypeAPI   ContactType = "api"
)

// Version represents a software version with support information
type Version struct {
	Name            string        `json:"name"`
	SemanticVersion string        `json:"semantic_version"`
	SupportedUntil  time.Time     `json:"supported_until"`
	Status          VersionStatus `json:"status"`
	IsLatest        bool          `json:"is_latest"`
	IsPrevious      bool          `json:"is_previous"`
}

// Contact represents security contact information
type Contact struct {
	Type         ContactType `json:"type"`
	Value        string      `json:"value"`
	ResponseTime string      `json:"response_time"`
	Description  string      `json:"description"`
}

// SecurityTool represents a simple security policy tool
type SecurityTool struct{}

// NewSecurityTool creates a new security tool instance
func NewSecurityTool() *SecurityTool {
	return &SecurityTool{}
}

// PolicyConfig represents security policy configuration
type PolicyConfig struct {
	Type         PolicyType       `json:"type"`
	Organization string            `json:"organization"`
	ContactEmail string            `json:"email"`
	OutputDir    string            `json:"output_dir"`
	Variables    map[string]string `json:"variables"`
}

// GeneratePolicy generates a security policy
func (st *SecurityTool) GeneratePolicy(ctx context.Context, config PolicyConfig) error {
	// Prepare template variables
	variables := st.prepareVariables(config)

	// Read template
	templateContent, err := st.readTemplate(config.Type)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Process template
	content := st.processTemplate(templateContent, variables)

	// Save to SECURITY.md
	outputFile := "SECURITY.md"
	if config.OutputDir != "." {
		outputFile = fmt.Sprintf("%s/SECURITY.md", strings.TrimSuffix(config.OutputDir, "/"))
	}

	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Generated %s for %s\n", outputFile, config.Organization)
	return nil
}

// prepareVariables prepares template variables with defaults
func (st *SecurityTool) prepareVariables(config PolicyConfig) map[string]string {
	variables := make(map[string]string)

	// Add user-provided variables
	for k, v := range config.Variables {
		variables[k] = v
	}

	// Set required variables with defaults
	variables["ORGANIZATION"] = config.Organization
	variables["CONTACT_EMAIL"] = config.ContactEmail
	variables["LATEST_VERSION"] = "1.0.0"
	variables["SUPPORT_END_DATE"] = time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	variables["LAST_UPDATED"] = time.Now().Format("2006-01-02")

	return variables
}

// readTemplate reads the appropriate template file
func (st *SecurityTool) readTemplate(policyType PolicyType) (string, error) {
	templatePath := "templates/SECURITY.md"
	
	// For now, we use the same template for both types
	// In the future, we could have different templates
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// processTemplate replaces template variables with values
func (st *SecurityTool) processTemplate(template string, variables map[string]string) string {
	content := template
	
	for placeholder, value := range variables {
		templateVar := fmt.Sprintf("{{%s}}", placeholder)
		content = strings.ReplaceAll(content, templateVar, value)
	}

	return content
}