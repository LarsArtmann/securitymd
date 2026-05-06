// Package internal provides tests for the internal package.
package internal

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/LarsArtmann/template-SECURITY/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityTool_GeneratePolicy(t *testing.T) {
	tool := NewSecurityTool()

	// Create a temporary template for testing
	tmpDir := t.TempDir()
	templateContent := `# Security Policy for {{.Organization}}

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| {{.LatestVersion}} | {{.SupportEndDate}} |

## Reporting a Vulnerability

Email us at {{.ContactEmail}}.

## Security Practices

We follow security best practices.

*Last updated: {{.LastUpdated}}*
`

	templatePath := tmpDir + "/SECURITY.md"
	err := os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name        string
		config      PolicyConfig
		expectError bool
	}{
		{
			name: "basic github policy",
			config: PolicyConfig{
				Type:         types.PolicyTypeGitHub,
				Organization: "TestOrg",
				ContactEmail: "security@test.org",
				OutputDir:    t.TempDir(),
				Variables:    map[string]string{"CUSTOM_VAR": "custom_value"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create template in the expected location
			templateDir := tt.config.OutputDir + "/../templates"
			err := os.MkdirAll(templateDir, 0o750)
			require.NoError(t, err)

			templateFile := templateDir + "/SECURITY.md"
			err = os.WriteFile(templateFile, []byte(templateContent), 0o600)
			require.NoError(t, err)

			// Change working directory temporarily
			originalWD, err := os.Getwd()
			require.NoError(t, err)

			t.Chdir(tt.config.OutputDir + "/..")
			t.Cleanup(func() {
				err := os.Chdir(originalWD)
				require.NoError(t, err)
			})

			err = tool.GeneratePolicy(context.TODO(), tt.config)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				// Check if file was created
				outputFile := tt.config.OutputDir + "/SECURITY.md"
				_, err = os.Stat(outputFile)
				require.NoError(t, err)

				// Check content
				content, err := os.ReadFile(outputFile)
				require.NoError(t, err)

				contentStr := string(content)
				assert.Contains(t, contentStr, tt.config.Organization)
				assert.Contains(t, contentStr, tt.config.ContactEmail)
				assert.Contains(t, contentStr, "## Supported Versions")
			}
		})
	}
}

func TestSecurityTool_PrepareTemplateData(t *testing.T) {
	t.Parallel()

	tool := NewSecurityTool()

	config := PolicyConfig{
		Type:         types.PolicyTypeGitHub,
		Organization: "TestOrg",
		ContactEmail: "security@test.org",
		Variables: map[string]string{
			"CUSTOM_FIELD": "custom_value",
		},
	}

	data := tool.prepareTemplateData(config)

	assert.Equal(t, "TestOrg", data.Organization)
	assert.Equal(t, "security@test.org", data.ContactEmail)
	assert.Equal(t, "1.0.0", data.LatestVersion)
	assert.Equal(t, "custom_value", data.AdditionalFields["CUSTOM_FIELD"])
	assert.Len(t, data.Versions, 1)
	assert.True(t, data.Versions[0].IsLatest)
	assert.Equal(t, types.StatusSupported, data.Versions[0].Status)
}

func TestSecurityTool_BuildVersions(t *testing.T) {
	t.Parallel()

	tool := NewSecurityTool()

	versions := tool.buildVersions()

	assert.Len(t, versions, 1)

	latest := versions[0]
	assert.Equal(t, "v2.x", latest.Name)
	assert.Equal(t, "2.0.0", latest.SemanticVersion)
	assert.Equal(t, types.StatusSupported, latest.Status)
	assert.True(t, latest.IsLatest)
	assert.False(t, latest.IsPrevious)

	// Check supported until date is approximately 1 year from now
	expected := time.Now().AddDate(1, 0, 0)
	assert.WithinDuration(t, expected, latest.SupportedUntil, time.Minute)
}

func TestSecurityTool_ProcessTemplate(t *testing.T) {
	t.Parallel()

	tool := NewSecurityTool()

	templateContent := `# Security Policy for {{.Organization}}

Contact: {{.ContactEmail}}
Latest Version: {{.LatestVersion}}
Last Updated: {{.LastUpdated}}
`

	data := TemplateData{
		Organization:  "TestOrg",
		ContactEmail:  "security@test.org",
		LatestVersion: "2.0.0",
		LastUpdated:   "2025-12-11",
	}

	result, err := tool.processTemplate(templateContent, data)
	require.NoError(t, err)

	expectedLines := []string{
		"# Security Policy for TestOrg",
		"Contact: security@test.org",
		"Latest Version: 2.0.0",
		"Last Updated: 2025-12-11",
	}

	for _, line := range expectedLines {
		assert.Contains(t, result, line)
	}
}
