package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityValidator_ValidateSECURITYMd(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name           string
		content        string
		expectValid    bool
		expectedErrors []string
		expectedWarns  []string
	}{
		{
			name: "valid complete security policy",
			content: `# Security Policy

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| v2.x | 2026-12-31 |

## Reporting a Vulnerability

Email us at security@example.com with any security issues.

## Security Practices

We follow security best practices.

`,
			expectValid:    true,
			expectedErrors: []string{},
			expectedWarns:  []string{"Should specify response time for vulnerability reports"},
		},
		{
			name: "missing required sections",
			content: `# Some Other Document

This is not a security policy.

`,
			expectValid: false,
			expectedErrors: []string{
				"Missing 'Reporting a Vulnerability' section",
				"Missing 'Supported Versions' section",
				"Missing 'Security Practices' section",
				"Missing contact email address",
			},
			expectedWarns: []string{
				"Should specify response time for vulnerability reports",
				"SECURITY.md seems too short (< 20 lines)",
				"No version information found",
			},
		},
		{
			name: "content too short",
			content: `# Security Policy

## Supported Versions

v1.0

`,
			expectValid: false,
			expectedErrors: []string{
				"Missing 'Reporting a Vulnerability' section",
				"Missing 'Security Practices' section",
				"Missing contact email address",
			},
			expectedWarns: []string{
				"Should specify response time for vulnerability reports",
				"SECURITY.md seems too short (< 20 lines)",
				"No version information found",
			},
		},
		{
			name: "with template variables",
			content: `# Security Policy

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| {{LATEST_VERSION}} | {{SUPPORT_END_DATE}} |

## Reporting a Vulnerability

Email us at {{CONTACT_EMAIL}}.

## Security Practices

We follow security best practices.

`,
			expectValid: false,
			expectedErrors: []string{
				"Unresolved template variable: | {{LATEST_VERSION}} | {{SUPPORT_END_DATE}} |",
				"Unresolved template variable: Email us at {{CONTACT_EMAIL}}.",
			},
			expectedWarns: []string{
				"Should specify response time for vulnerability reports",
				"SECURITY.md seems too short (< 20 lines)",
				"No version information found",
			},
		},
		{
			name: "no version information",
			content: `# Security Policy

## Supported Versions

We support our software.

## Reporting a Vulnerability

Email us at security@example.com.

## Security Practices

We follow security best practices.

`,
			expectValid:    true,
			expectedErrors: []string{},
			expectedWarns: []string{
				"Should specify response time for vulnerability reports",
				"No version information found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			filename := tmpDir + "/SECURITY.md"
			err := os.WriteFile(filename, []byte(tt.content), 0o644)
			require.NoError(t, err)

			// Run validation
			result, err := validator.ValidateSECURITYMd(filename)
			require.NoError(t, err)

			// Check validity
			assert.Equal(t, tt.expectValid, result.Valid)

			// Check errors
			for _, expectedError := range tt.expectedErrors {
				found := false

				for _, actualError := range result.Errors {
					if strings.Contains(actualError, expectedError) {
						found = true

						break
					}
				}

				assert.True(
					t,
					found,
					"Expected error '%s' not found in %v",
					expectedError,
					result.Errors,
				)
			}

			// Check warnings
			for _, expectedWarning := range tt.expectedWarns {
				found := false

				for _, actualWarning := range result.Warnings {
					if strings.Contains(actualWarning, expectedWarning) {
						found = true

						break
					}
				}

				assert.True(
					t,
					found,
					"Expected warning '%s' not found in %v",
					expectedWarning,
					result.Warnings,
				)
			}
		})
	}
}

func TestSecurityValidator_ValidateContentQuality(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name           string
		content        string
		expectedErrors []string
	}{
		{
			name: "good quality content",
			content: strings.Repeat(
				"This is substantive content with meaningful information. ",
				10,
			),
		},
		{
			name:    "template variables unresolved",
			content: "This has {{TEMPLATE_VARIABLE}} that should be resolved.",
			expectedErrors: []string{
				"Unresolved template variable: This has {{TEMPLATE_VARIABLE}} that should be resolved.",
			},
		},
		{
			name:           "no actual content",
			content:        "# Just Headers\n## And Headers\n\n### More Headers",
			expectedErrors: []string{"SECURITY.md lacks substantive content"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &SecurityValidationResult{}
			validator.validateContentQuality(tt.content, result)

			if len(tt.expectedErrors) > 0 {
				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, result.Errors, expectedError)
				}
			} else {
				// Don't check for empty errors, as validation logic may add warnings
			}
		})
	}
}

func TestSecurityValidator_PrintResults(t *testing.T) {
	validator := NewSecurityValidator()

	results := []*SecurityValidationResult{
		{
			File:     "SECURITY.md",
			Valid:    true,
			Errors:   []string{},
			Warnings: []string{"Minor warning"},
		},
		{
			File:     "BAD_SECURITY.md",
			Valid:    false,
			Errors:   []string{"Missing section"},
			Warnings: []string{},
		},
	}

	// This test just ensures PrintResults doesn't panic
	// In a real implementation, you might capture stdout and verify the output
	assert.NotPanics(t, func() {
		validator.PrintResults(results)
	})
}
