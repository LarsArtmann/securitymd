package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"
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

// TemplateData represents data for template rendering
type TemplateData struct {
	Organization     string
	ContactEmail     string
	LatestVersion    string
	SupportEndDate   string
	LastUpdated      string
	Versions         []Version
	AdditionalFields map[string]interface{}
}

// GeneratePolicy generates a security policy
func (st *SecurityTool) GeneratePolicy(ctx context.Context, config PolicyConfig) error {
	// Prepare template data
	templateData := st.prepareTemplateData(config)

	// Read template
	templateContent, err := st.readTemplate(config.Type)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Process template
	content, err := st.processTemplate(templateContent, templateData)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

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

// prepareTemplateData prepares template data with defaults
func (st *SecurityTool) prepareTemplateData(config PolicyConfig) TemplateData {
	// Start with required fields
	data := TemplateData{
		Organization:     config.Organization,
		ContactEmail:     config.ContactEmail,
		LatestVersion:    "1.0.0",
		SupportEndDate:   time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
		LastUpdated:      time.Now().Format("2006-01-02"),
		Versions:         st.buildVersions(),
		AdditionalFields: make(map[string]interface{}),
	}

	// Add user-provided variables
	for k, v := range config.Variables {
		data.AdditionalFields[k] = v
	}

	return data
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

// processTemplate processes template using Go's text/template
func (st *SecurityTool) processTemplate(templateContent string, data TemplateData) (string, error) {
	// Create template
	tmpl, err := template.New("security").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (st *SecurityTool) buildVersions() []Version {
	now := time.Now()

	return []Version{
		{
			Name:            "v2.x",
			SemanticVersion: "2.0.0",
			SupportedUntil:  now.AddDate(1, 0, 0),
			Status:          StatusSupported,
			IsLatest:        true,
		},
	}
}