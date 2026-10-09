package policy

import (
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// noResponsePolicy drops the response commitment from validPolicy: it yields
// exactly one finding (missing-response-time) to suppress.
const noResponsePolicy = `# Security Policy

## Supported Versions

| Version | Supported Until |
| ------- | --------------- |
| v2.x    | 2026-12-31     |

Only the latest release receives security fixes.

## Reporting a Vulnerability

Email us at security@example.com with any security issues.

## Security Practices

We follow security best practices in development and operations.

All changes are reviewed before merge and CI runs security scanning.
`

func TestValidate_honors_suppression_comment(t *testing.T) {
	t.Parallel()

	content := noResponsePolicy + "\n<!-- securitymd:ignore(missing-response-time) response time published in the wiki -->\n"

	findings := validateContent(t, content)

	suppressed := findingsSuppressedByRule(findings, "missing-response-time")
	require.Len(t, suppressed, 1, "the warning must be suppressed")

	assert.Equal(t, finding.SuppressionInSource, suppressed[0].Suppression.Kind)
	assert.Equal(t, finding.RuleName("missing-response-time"), suppressed[0].Suppression.Rule)
	assert.Equal(t, "response time published in the wiki", suppressed[0].Suppression.Reason)
	assert.Equal(t, finding.SeverityWarning, suppressed[0].Severity,
		"suppression marks metadata, it must not rewrite severity")
}

func TestValidate_suppression_requires_reason(t *testing.T) {
	t.Parallel()

	content := noResponsePolicy + "\n<!-- securitymd:ignore(missing-response-time) -->\n"

	findings := validateContent(t, content)

	assert.Empty(t, findingsSuppressedByRule(findings, "missing-response-time"),
		"a suppression without a reason is inert")
	assert.Contains(t, ruleIDs(findings), "missing-response-time",
		"the finding must survive an inert suppression")
}

func TestValidate_suppression_unknown_rule_is_inert(t *testing.T) {
	t.Parallel()

	content := validPolicy + "\n<!-- securitymd:ignore(not-a-real-rule) typo must fail safe -->\n"

	findings := validateContent(t, content)

	assert.NotContains(t, ruleIDs(findings), "not-a-real-rule")

	for _, f := range findings {
		assert.Nil(t, f.Suppression, "no real rule may be suppressed by an unknown rule name")
	}
}

func TestValidate_suppression_multiple_rules(t *testing.T) {
	t.Parallel()

	thin := "# Security Policy\n\nContact security@example.com about anything.\n"
	content := thin + "\n<!-- securitymd:ignore(too-short,no-version-info) bootstrap stage, filled next sprint -->\n"

	findings := validateContent(t, content)

	assert.Len(t, findingsSuppressedByRule(findings, "too-short"), 1)
	assert.Len(t, findingsSuppressedByRule(findings, "no-version-info"), 1)
	assert.Equal(t, "bootstrap stage, filled next sprint",
		findingsSuppressedByRule(findings, "too-short")[0].Suppression.Reason)
}

func TestValidate_missing_file_cannot_be_suppressed(t *testing.T) {
	t.Parallel()

	// missing-file has no file content to carry a comment: the rule must not
	// appear in the suppressible set, so a config can never silence it.
	assert.NotContains(t, KnownRuleIDs(), finding.RuleName("suppressible-nonexistent"))

	findings, err := Detect(withWorkingDir(t.Context(), t.TempDir()))
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Nil(t, findings[0].Suppression)
}

func findingsSuppressedByRule(findings []finding.Finding, rule string) []finding.Finding {
	var matched []finding.Finding

	for _, f := range findings {
		if string(f.Rule) == rule && f.Suppression != nil {
			matched = append(matched, f)
		}
	}

	return matched
}
