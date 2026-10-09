package policy

import (
	"time"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnownRuleIDs_covers_every_rule(t *testing.T) {
	t.Parallel()

	ids := KnownRuleIDs()

	for _, expected := range []string{
		"missing-header",
		"missing-reporting",
		"missing-versions",
		"missing-practices",
		"missing-contact",
		"missing-response-time",
		"too-short",
		"unresolved-template",
		"no-content",
		"no-version-info",
		"missing-file",
	} {
		assert.Contains(t, ids, finding.RuleName(expected))
	}

	assert.Len(t, ids, 11, "a new rule must be added to the known-ID table in the same change")
}

func TestParseSeverityOverrides(t *testing.T) {
	t.Parallel()

	t.Run("parses a downgrade spec", func(t *testing.T) {
		t.Parallel()

		overrides, err := ParseSeverityOverrides("missing-file=warning, too-short=info")
		require.NoError(t, err)

		assert.Equal(t, SeverityOverrides{
			"missing-file": finding.SeverityWarning,
			"too-short":    finding.SeverityInfo,
		}, overrides)
	})

	t.Run("empty spec means no overrides", func(t *testing.T) {
		t.Parallel()

		overrides, err := ParseSeverityOverrides("")
		require.NoError(t, err)
		assert.Empty(t, overrides)
	})

	t.Run("rejects unknown rule with the known list", func(t *testing.T) {
		t.Parallel()

		_, err := ParseSeverityOverrides("missing-fiel=warning")

		require.Error(t, err)
		assert.Contains(t, err.Error(), `unknown rule "missing-fiel"`)
		assert.Contains(t, err.Error(), "missing-file", "the error must name the accepted values")
	})

	t.Run("rejects invalid severity", func(t *testing.T) {
		t.Parallel()

		_, err := ParseSeverityOverrides("missing-file=fatal")

		require.Error(t, err)
		assert.Contains(t, err.Error(), `invalid severity "fatal"`)
		assert.Contains(t, err.Error(), "missing-file")
	})

	t.Run("rejects malformed part", func(t *testing.T) {
		t.Parallel()

		_, err := ParseSeverityOverrides("missing-file-warning")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "want rule=severity")
	})
}

func TestApplySeverityOverrides_downgrades_missing_file(t *testing.T) {
	t.Parallel()

	findings, err := Detect(withWorkingDir(t.Context(), t.TempDir()))
	require.NoError(t, err)
	require.Len(t, findings, 1)
	require.Equal(t, finding.SeverityError, findings[0].Severity)

	overrides, err := ParseSeverityOverrides("missing-file=warning")
	require.NoError(t, err)

	adjusted := ApplySeverityOverrides(findings, overrides)
	require.Len(t, adjusted, 1)
	assert.Equal(t, finding.SeverityWarning, adjusted[0].Severity,
		"the adoption unblocker: missing-file downgraded to warning")
	assert.Nil(t, adjusted[0].Suppression,
		"severity overrides must not masquerade as suppressions")
}

func TestApplySeverityOverrides_empty_is_noop(t *testing.T) {
	t.Parallel()

	findings := validateContent(t, policyWithoutResponseTime(t))
	require.NotEmpty(t, findings)

	assert.Equal(t, findings, ApplySeverityOverrides(findings, SeverityOverrides{}),
		"an empty override map must return the findings unchanged")
}

// The CLI composes the two escape hatches: Validate applies in-file
// suppressions first, then the --set-severity overrides rewrite severities on
// the result. This pins the combination on the same rule — the edge the
// report flagged as unproven.
func TestSeverityOverrides_combine_with_suppressions(t *testing.T) {
	t.Parallel()

	content := policyWithoutResponseTime(t) +
		"\n<!-- securitymd:ignore(missing-response-time) tracked in the support wiki -->\n"

	findings := validateContent(t, content)
	require.Len(t, findings, 1)
	require.NotNil(t, findings[0].Suppression, "precondition: the warning is suppressed")

	overrides, err := ParseSeverityOverrides("missing-response-time=critical")
	require.NoError(t, err)

	adjusted := ApplySeverityOverrides(findings, overrides)
	require.Len(t, adjusted, 1)

	assert.Equal(t, finding.SeverityCritical, adjusted[0].Severity,
		"the override rewrites severity even on a suppressed finding")
	assert.NotNil(t, adjusted[0].Suppression,
		"the override must not strip the suppression evidence")
	assert.True(t, adjusted[0].IsSuppressedAt(time.Now()),
		"severity escalation must not resurrect a suppressed finding")
}

// The gate contract the CLI implements: suppressed findings never trip exit 1,
// even when their severity was escalated by an override — and the same
// escalation without a suppression DOES trip it, so the suppression is what
// keeps CI green, not the override being harmless.
func TestSeverityOverrides_suppressed_findings_stay_exit_neutral(t *testing.T) {
	t.Parallel()

	overrides, err := ParseSeverityOverrides("missing-response-time=critical")
	require.NoError(t, err)

	activeErrors := func(findings []finding.Finding) []finding.Finding {
		active := finding.Filter(findings, func(f finding.Finding) bool {
			return !f.IsSuppressedAt(time.Now())
		})

		return finding.Filter(active, finding.BySeverityAtLeast(finding.SeverityError))
	}

	suppressed := validateContent(t, policyWithoutResponseTime(t) +
		"\n<!-- securitymd:ignore(missing-response-time) tracked in the support wiki -->\n")
	assert.Empty(t, activeErrors(ApplySeverityOverrides(suppressed, overrides)),
		"an escalated-but-suppressed finding must not activate: the escape hatch holds")

	unsuppressed := validateContent(t, policyWithoutResponseTime(t))
	assert.NotEmpty(t, activeErrors(ApplySeverityOverrides(unsuppressed, overrides)),
		"the same escalation without an in-file suppression activates")
}
