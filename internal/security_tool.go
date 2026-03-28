package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

// Error types.
var (
	ErrConfigNotFound    = errors.New("configuration file not found")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrTemplateNotFound  = errors.New("template not found")
	ErrTemplateParse     = errors.New("failed to parse template")
	ErrFileWrite         = errors.New("failed to write file")
	ErrInvalidPolicyType = errors.New("invalid policy type")
	ErrMissingField      = errors.New("required field missing")
)

// SecurityError represents a structured error with code and context.
type SecurityError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Cause   error  `json:"cause,omitempty"`
}

func (e *SecurityError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SecurityError) Unwrap() error {
	return e.Cause
}

// NewSecurityError creates a new SecurityError.
func NewSecurityError(code, message, field string, cause error) *SecurityError {
	return &SecurityError{
		Code:    code,
		Message: message,
		Field:   field,
		Cause:   cause,
	}
}

// PolicyType represents different types of security policies.
type PolicyType string

const (
	PolicyTypeGitHub     PolicyType = "github"
	PolicyTypeEnterprise PolicyType = "enterprise"
)

// VersionStatus represents the status of a version.
type VersionStatus string

const (
	StatusSupported  VersionStatus = "supported"
	StatusDeprecated VersionStatus = "deprecated"
	StatusEOL        VersionStatus = "end-of-life"
)

// ContactType represents different types of contact methods.
type ContactType string

const (
	ContactTypeEmail ContactType = "email"
	ContactTypeWeb   ContactType = "web"
	ContactTypeAPI   ContactType = "api"
)

// Version represents a software version with support information.
type Version struct {
	Name            string        `json:"name"`
	SemanticVersion string        `json:"semantic_version"`
	SupportedUntil  time.Time     `json:"supported_until"`
	Status          VersionStatus `json:"status"`
	IsLatest        bool          `json:"is_latest"`
	IsPrevious      bool          `json:"is_previous"`
}

// Contact represents security contact information.
type Contact struct {
	Type         ContactType `json:"type"`
	Value        string      `json:"value"`
	ResponseTime string      `json:"response_time"`
	Description  string      `json:"description"`
}

// SecurityTool represents a simple security policy tool.
type SecurityTool struct{}

// NewSecurityTool creates a new security tool instance.
func NewSecurityTool() *SecurityTool {
	return &SecurityTool{}
}

// LoadConfig loads configuration from file.
func (st *SecurityTool) LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	v.SetConfigFile(configPath)

	// Enable environment variable support
	v.AutomaticEnv()
	v.SetEnvPrefix("TEMPLATE_SECURITY")

	// Read config file
	err := v.ReadInConfig()
	if err != nil {
		return nil, NewSecurityError(
			"CONFIG_READ_FAILED",
			"failed to read config file",
			configPath,
			err,
		)
	}

	// Unmarshal into config struct
	var config Config

	err = v.Unmarshal(&config)
	if err != nil {
		return nil, NewSecurityError(
			"CONFIG_UNMARSHAL_FAILED",
			"failed to unmarshal config",
			"",
			err,
		)
	}

	// Set defaults
	if config.Type == "" {
		config.Type = PolicyTypeGitHub
	}

	if config.SupportYears == 0 {
		config.SupportYears = 1
	}

	if config.DefaultVersion == "" {
		config.DefaultVersion = "1.0.0"
	}

	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if config.TemplateDir == "" {
		config.TemplateDir = "templates"
	}

	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}

	return &config, nil
}

// FindConfigFile finds configuration file in current directory or parent directories.
func (st *SecurityTool) FindConfigFile() string {
	// Check for config files in order of preference
	configNames := []string{
		".template-security.yaml",
		".template-security.yml",
		"template-security.yaml",
		"template-security.yml",
	}

	// Start in current directory and go up
	dir, _ := os.Getwd()

	for {
		for _, name := range configNames {
			configPath := filepath.Join(dir, name)
			if _, err := os.Stat(configPath); err == nil {
				return configPath
			}
		}

		// Go up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root
		}

		dir = parent
	}

	return ""
}

// Config represents a configuration file structure.
type Config struct {
	Organization   string            `yaml:"organization"`
	ContactEmail   string            `yaml:"contact_email"`
	Type           PolicyType        `yaml:"type"            mapstructure:"type"`
	OutputDir      string            `yaml:"output_dir"`
	TemplateDir    string            `yaml:"template_dir"`
	DefaultVersion string            `yaml:"default_version"`
	SupportYears   int               `yaml:"support_years"`
	Variables      map[string]string `yaml:"variables"`
}

// PolicyConfig represents security policy configuration.
type PolicyConfig struct {
	Type         PolicyType        `json:"type"`
	Organization string            `json:"organization"`
	ContactEmail string            `json:"email"`
	OutputDir    string            `json:"output_dir"`
	Variables    map[string]string `json:"variables"`
}

// TemplateData represents data for template rendering.
type TemplateData struct {
	Organization     string
	ContactEmail     string
	LatestVersion    string
	SupportEndDate   string
	LastUpdated      string
	Versions         []Version
	AdditionalFields map[string]any
}

// GeneratePolicy generates a security policy.
func (st *SecurityTool) GeneratePolicy(ctx context.Context, config PolicyConfig) error {
	// Prepare template data
	templateData := st.prepareTemplateData(config)

	// Read template
	templateContent, err := st.readTemplate(config.Type)
	if err != nil {
		return NewSecurityError(
			"TEMPLATE_READ_FAILED",
			"failed to read template",
			string(config.Type),
			err,
		)
	}

	// Process template
	content, err := st.processTemplate(templateContent, templateData)
	if err != nil {
		return NewSecurityError("TEMPLATE_PROCESS_FAILED", "failed to process template", "", err)
	}

	// Save to SECURITY.md
	outputFile := "SECURITY.md"
	if config.OutputDir != "." {
		outputFile = strings.TrimSuffix(config.OutputDir, "/") + "/SECURITY.md"
	}

	if err := os.WriteFile(outputFile, []byte(content), 0o644); err != nil {
		return NewSecurityError("FILE_WRITE_FAILED", "failed to write file", outputFile, err)
	}

	fmt.Printf("✅ Generated %s for %s\n", outputFile, config.Organization)

	return nil
}

// prepareTemplateData prepares template data with defaults.
func (st *SecurityTool) prepareTemplateData(config PolicyConfig) TemplateData {
	// Start with required fields
	supportYears := 1 // default

	if years, exists := config.Variables["SUPPORT_YEARS"]; exists {
		if parsed, err := time.ParseDuration(years + "y"); err == nil {
			supportYears = int(parsed.Hours() / (24 * 365))
		}
	}

	data := TemplateData{
		Organization:     config.Organization,
		ContactEmail:     config.ContactEmail,
		LatestVersion:    "1.0.0",
		SupportEndDate:   time.Now().AddDate(supportYears, 0, 0).Format("2006-01-02"),
		LastUpdated:      time.Now().Format("2006-01-02"),
		Versions:         st.buildVersions(),
		AdditionalFields: make(map[string]any),
	}

	// Add user-provided variables
	for k, v := range config.Variables {
		data.AdditionalFields[k] = v
	}

	return data
}

// readTemplate reads the appropriate template file.
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

// processTemplate processes template using Go's text/template.
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
