package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/LarsArtmann/template-SECURITY/internal/types"
	finding "github.com/larsartmann/go-finding"
	"github.com/spf13/viper"
)

// Error types.
var (
	ErrConfigNotFound    = finding.NewValidationError("configuration file not found", nil)
	ErrInvalidConfig     = finding.NewValidationError("invalid configuration", nil)
	ErrTemplateNotFound  = finding.NewValidationError("template not found", nil)
	ErrTemplateParse     = finding.NewValidationError("failed to parse template", nil)
	ErrFileWrite         = finding.NewIOError("failed to write file", nil)
	ErrInvalidPolicyType = finding.NewValidationError("invalid policy type", nil)
	ErrMissingField      = finding.NewValidationError("required field missing", nil)
)

// SecurityTool represents a simple security policy tool.
type SecurityTool struct{}

// NewSecurityTool creates a new security tool instance.
func NewSecurityTool() *SecurityTool {
	return &SecurityTool{}
}

// LoadConfig loads configuration from file.
func (st *SecurityTool) LoadConfig(configPath string) (*Config, error) {
	configViper := viper.New()

	configViper.SetConfigFile(configPath)

	configViper.AutomaticEnv()
	configViper.SetEnvPrefix("TEMPLATE_SECURITY")

	err := configViper.ReadInConfig()
	if err != nil {
		return nil, finding.NewIOError(
			"failed to read config file "+configPath,
			err,
		)
	}

	var config Config

	err = configViper.Unmarshal(&config)
	if err != nil {
		return nil, finding.NewParseError("failed to unmarshal config", err)
	}

	if config.Type == "" {
		config.Type = types.PolicyTypeGitHub
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
	configNames := []string{
		".template-security.yaml",
		".template-security.yml",
		"template-security.yaml",
		"template-security.yml",
	}

	dir, _ := os.Getwd()

	for {
		for _, name := range configNames {
			configPath := filepath.Join(dir, name)
			if _, err := os.Stat(configPath); err == nil {
				return configPath
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	return ""
}

// Config represents a configuration file structure.
type Config struct {
	Organization   string            `mapstructure:"organization"   yaml:"organization"`
	ContactEmail   string            `mapstructure:"contactEmail"   yaml:"contactEmail"`
	Type           types.PolicyType  `mapstructure:"type"           yaml:"type"`
	OutputDir      string            `mapstructure:"outputDir"      yaml:"outputDir"`
	TemplateDir    string            `mapstructure:"templateDir"    yaml:"templateDir"`
	DefaultVersion string            `mapstructure:"defaultVersion" yaml:"defaultVersion"`
	SupportYears   int               `mapstructure:"supportYears"   yaml:"supportYears"`
	Variables      map[string]string `mapstructure:"variables"      yaml:"variables"`
}

// PolicyConfig represents security policy configuration.
type PolicyConfig struct {
	Type         types.PolicyType  `json:"type"`
	Organization string            `json:"organization"`
	ContactEmail string            `json:"contactEmail"`
	OutputDir    string            `json:"outputDir"`
	Variables    map[string]string `json:"variables"`
}

// TemplateData represents data for template rendering.
type TemplateData struct {
	Organization     string
	ContactEmail     string
	LatestVersion    string
	SupportEndDate   string
	LastUpdated      string
	Versions         []types.Version
	AdditionalFields map[string]any
}

// GeneratePolicy generates a security policy.
func (st *SecurityTool) GeneratePolicy(_ context.Context, config PolicyConfig) error {
	templateData := st.prepareTemplateData(config)

	templateContent, err := st.readTemplate(config.Type)
	if err != nil {
		return finding.NewIOError(
			fmt.Sprintf("failed to read template for type %s", config.Type),
			err,
		)
	}

	content, err := st.processTemplate(templateContent, templateData)
	if err != nil {
		return finding.NewParseError("failed to process template", err)
	}

	outputFile := "SECURITY.md"
	if config.OutputDir != "." {
		outputFile = strings.TrimSuffix(config.OutputDir, "/") + "/SECURITY.md"
	}

	err = os.WriteFile(outputFile, []byte(content), 0o600)
	if err != nil {
		return finding.NewIOError(
			"failed to write file "+outputFile,
			err,
		)
	}

	fmt.Printf("✅ Generated %s for %s\n", outputFile, config.Organization)

	return nil
}

// prepareTemplateData prepares template data with defaults.
func (st *SecurityTool) prepareTemplateData(config PolicyConfig) TemplateData {
	supportYears := 1

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

	for k, v := range config.Variables {
		data.AdditionalFields[k] = v
	}

	return data
}

func (st *SecurityTool) readTemplate(_ types.PolicyType) (string, error) {
	templatePath := "templates/SECURITY.md"

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	return string(content), nil
}

func (st *SecurityTool) processTemplate(templateContent string, data TemplateData) (string, error) {
	tmpl, err := template.New("security").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf(
			"failed to parse template (content length %d): %w",
			len(templateContent),
			err,
		)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf(
			"failed to execute template for org=%q email=%q: %w",
			data.Organization,
			data.ContactEmail,
			err,
		)
	}

	return buf.String(), nil
}

func (st *SecurityTool) buildVersions() []types.Version {
	now := time.Now()

	return []types.Version{
		{
			Name:            "v2.x",
			SemanticVersion: "2.0.0",
			SupportedUntil:  now.AddDate(1, 0, 0),
			Status:          types.StatusSupported,
			IsLatest:        true,
			IsPrevious:      false,
		},
	}
}
