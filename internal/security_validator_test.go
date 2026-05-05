package internal

import (
	"os"
	"strings"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertContainsAll(t *testing.T, expected, actual []string, itemType string) {
	for _, expectedItem := range expected {
		found := false

		for _, actualItem := range actual {
			if strings.Contains(actualItem, expectedItem) {
				found = true

				break
			}
		}

		assert.True(
			t,
			found,
			"Expected %s '%s' not found in %v",
			itemType,
			expectedItem,
			actual,
		)
	}
}

func collectMessages(report *finding.Report, severity finding.Severity) []string {
	var messages []string

	for _, f := range report.Findings {
		if f.Severity == severity {
			messages = append(messages, f.Message)
		}
	}

	return messages
}

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
			tmpDir := t.TempDir()
			filename := tmpDir + "/SECURITY.md"
			err := os.WriteFile(filename, []byte(tt.content), 0o644)
			require.NoError(t, err)

			report, err := validator.ValidateSECURITYMd(filename)
			require.NoError(t, err)

			valid := ReportIsValid(report)
			assert.Equal(t, tt.expectValid, valid)

			assertContainsAll(
				t,
				tt.expectedErrors,
				collectMessages(report, finding.SeverityError),
				"error",
			)
			assertContainsAll(
				t,
				tt.expectedWarns,
				collectMessages(report, finding.SeverityWarning),
				"warning",
			)
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
			report := finding.NewReport(finding.ToolInfo{Name: toolName})
			lines := strings.Split(tt.content, "\n")
			validator.validateContentQuality("test.md", tt.content, lines, report)

			if len(tt.expectedErrors) > 0 {
				errorMessages := collectMessages(report, finding.SeverityError)
				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, errorMessages, expectedError)
				}
			}
		})
	}
}

func TestSecurityValidator_PrintResults(t *testing.T) {
	validator := NewSecurityValidator()

	report1 := finding.NewReport(finding.ToolInfo{Name: toolName})
	report1.AddFinding(
		finding.NewFinding(
			"test",
			toolName,
			"Minor warning",
			finding.SeverityWarning,
			finding.Pos("SECURITY.md", 1, 0),
			1.0,
		),
	)

	report2 := finding.NewReport(finding.ToolInfo{Name: toolName})
	report2.AddFinding(
		finding.NewFinding(
			"test",
			toolName,
			"Missing section",
			finding.SeverityError,
			finding.Pos("BAD_SECURITY.md", 1, 0),
			1.0,
		),
	)

	reports := []*finding.Report{report1, report2}

	assert.NotPanics(t, func() {
		validator.PrintResults(reports)
	})
}
