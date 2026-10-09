package policy

import (
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_honors_suppression_comment(t *testing.T) {
	t.Parallel()

	content := policyWithoutResponseTime(t) + "\n<!-- securitymd:ignore(missing-response-time) commitment published in the support wiki -->\n"

	findings := validateContent(t, content)

	suppressed := findingsSuppressedByRule(findings, "missing-response-time")
	require.Len(t, suppressed, 1, "the warning must be suppressed")

	assert.Equal(t, finding.SuppressionInSource, suppressed[0].Suppression.Kind)
	assert.Equal(t, finding.RuleName("missing-response-time"), suppressed[0].Suppression.Rule)
	assert.Equal(t, "commitment published in the support wiki", suppressed[0].Suppression.Reason)
	assert.Equal(t, finding.SeverityWarning, suppressed[0].Severity,
		"suppression marks metadata, it must not rewrite severity")
}

func TestValidate_suppression_requires_reason(t *testing.T) {
	t.Parallel()

	content := policyWithoutResponseTime(t) + "\n<!-- securitymd:ignore(missing-response-time) -->\n"

	findings := validateContent(t, content)

	assert.Empty(t, findingsSuppressedByRule(findings, "missing-response-time"),
		"a suppression without a reason is inert")
	assert.Contains(t, ruleIDs(findings), "missing-response-time",
		"the finding must survive an inert suppression")
}

// The rules are substring-based, so suppression-comment TEXT itself counts as
// content: a reason mentioning "response time" satisfies missing-response-time
// before suppression is even consulted. Pinned here so a future structure-aware
// parser (ROADMAP) knows the incumbent behavior it must stay compatible with.
func TestValidate_suppression_text_counts_as_content(t *testing.T) {
	t.Parallel()

	content := policyWithoutResponseTime(t) + "\n<!-- securitymd:ignore(missing-response-time) we respond within 48 hours per wiki -->\n"

	findings := validateContent(t, content)

	assert.NotContains(t, ruleIDs(findings), "missing-response-time",
		"the reason text satisfies the substring rule, leaving nothing to suppress")
}

func TestValidate_suppression_unknown_rule_is_inert(t *testing.T) {
	t.Parallel()

	content := compliantPolicy(t) + "\n<!-- securitymd:ignore(not-a-real-rule) typo must fail safe -->\n"

	findings := validateContent(t, content)

	assert.NotContains(t, ruleIDs(findings), "not-a-real-rule")

	for _, f := range findings {
		assert.Nil(t, f.Suppression, "no real rule may be suppressed by an unknown rule name")
	}
}

func TestValidate_suppression_multiple_rules(t *testing.T) {
	t.Parallel()

	// Missing the versions section and any response commitment, and short
	// enough for too-short. The pair under suppression (missing-versions +
	// missing-response-time) is deliberate: rule names carrying "version"
	// (no-version-info, missing-versions as a STRING in the comment) satisfy
	// hasVersionInformation themselves, so no-version-info can never be
	// suppressed in-file — its name defeats its own detection.
	thin := `# Security Policy

## Reporting a Vulnerability

Email security@example.com; we triage quickly.

Include reproduction steps and potential impact.

## Security Practices

All changes are reviewed and dependencies are scanned.

Access is audited quarterly across all systems.
`
	content := thin + "\n<!-- securitymd:ignore(missing-versions,missing-response-time) tracked in the docs issue -->\n"

	findings := validateContent(t, content)

	assert.Len(t, findingsSuppressedByRule(findings, "missing-versions"), 1)
	assert.Len(t, findingsSuppressedByRule(findings, "missing-response-time"), 1)
	assert.Equal(t, "tracked in the docs issue",
		findingsSuppressedByRule(findings, "missing-versions")[0].Suppression.Reason)
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
