package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validPolicy = `# Security Policy

## Supported Versions

| Version | Supported Until |
| ------- | --------------- |
| v2.x    | 2026-12-31     |

## Reporting a Vulnerability

Email us at security@example.com with any security issues.

## Security Practices

We follow security best practices in development and operations.
`

func validateContent(t *testing.T, content string) []finding.Finding {
	t.Helper()

	path := filepath.Join(t.TempDir(), "SECURITY.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	findings, err := Validate(path)
	require.NoError(t, err)

	return findings
}

func ruleIDs(findings []finding.Finding) []string {
	ids := make([]string, 0, len(findings))
	for _, f := range findings {
		ids = append(ids, string(f.Rule))
	}

	return ids
}

func TestValidate_accepts_compliant_policy(t *testing.T) {
	t.Parallel()

	findings := validateContent(t, validPolicy)
	assert.Empty(t, findings, "a compliant policy must yield zero findings")
}

func TestValidate_reports_each_missing_section(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// remove strips one section from the valid policy
		remove    string
		expectIDs []string
	}{
		{"header", "# Security Policy\n", []string{"missing-header"}},
		{"reporting", "## Reporting a Vulnerability", []string{"missing-reporting", "missing-contact"}},
		{"versions", "## Supported Versions", []string{"missing-versions"}},
		{"practices", "## Security Practices", []string{"missing-practices"}},
		{"contact", "security@example.com", []string{"missing-contact"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			broken := strings.ReplaceAll(validPolicy, test.remove, "")
			findings := validateContent(t, broken)

			for _, expected := range test.expectIDs {
				assert.Contains(t, ruleIDs(findings), expected)
			}
		})
	}
}

func TestValidate_accepts_advisory_link_as_contact(t *testing.T) {
	t.Parallel()

	content := strings.ReplaceAll(validPolicy,
		"Email us at security@example.com with any security issues.",
		"Report via https://github.com/AcmeCorp/widget/security/advisories/new privately.")

	findings := validateContent(t, content)
	assert.NotContains(t, ruleIDs(findings), "missing-contact")
}

func TestValidate_flags_unresolved_template_variables_with_line(t *testing.T) {
	t.Parallel()

	findings := validateContent(t, "# Security Policy\n\nContact {{.ContactEmail}} for issues.\n")

	var unresolved []finding.Finding
	for _, f := range findings {
		if f.Rule == "unresolved-template" {
			unresolved = append(unresolved, f)
		}
	}

	require.Len(t, unresolved, 1)
	assert.Equal(t, 3, unresolved[0].Position.Line,
		"unresolved-template must point at the exact offending line")
	assert.Equal(t, finding.SeverityError, unresolved[0].Severity)
}

func TestValidate_flags_thin_and_placeholder_content(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		content   string
		expectIDs []string
	}{
		{
			"too short",
			"# Security Policy\n",
			[]string{"too-short", "no-content", "no-version-info"},
		},
		{
			"placeholder only",
			"# Security Policy\n\n## Reporting a Vulnerability\n\nContact security@example.com about anything.\n\n## Supported Versions\n\nSee versions.\n\n## Security Practices\n\nFill this in later.\n",
			[]string{"no-content"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			findings := validateContent(t, test.content)
			assert.Subset(t, ruleIDs(findings), test.expectIDs)
		})
	}
}

func TestValidate_generated_template_passes(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		ContactEmail: "security@acme.com",
	})
	require.NoError(t, err)
	require.True(t, result.Wrote)

	findings := validateContent(t, mustRead(t, result.Path))
	assert.Empty(t, findings,
		"the embedded template must satisfy its own validator (dogfood invariant)")
}

func mustRead(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content)
}
