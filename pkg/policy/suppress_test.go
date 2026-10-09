package policy

import (
	"testing"
	"time"

	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/gomega"
	"github.com/samber/lo"
)

func TestValidate_honors_suppression_comment(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(
		t,
	) + "\n<!-- securitymd:ignore(missing-response-time) commitment published in the support wiki -->\n"

	findings := validateContent(t, content)

	suppressed := findingsSuppressedByRule(findings, "missing-response-time")
	g.Expect(suppressed).To(gomega.HaveLen(1), "the warning must be suppressed")

	g.Expect(suppressed[0].Suppression.Kind).To(gomega.Equal(finding.SuppressionInSource))
	g.Expect(suppressed[0].Suppression.Rule).To(gomega.Equal(finding.RuleName("missing-response-time")))
	g.Expect(suppressed[0].Suppression.Reason).To(gomega.Equal("commitment published in the support wiki"))
	g.Expect(suppressed[0].Severity).To(gomega.Equal(finding.SeverityWarning),
		"suppression marks metadata, it must not rewrite severity")
}

func TestValidate_suppression_requires_reason(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(t) + "\n<!-- securitymd:ignore(missing-response-time) -->\n"

	findings := validateContent(t, content)

	g.Expect(findingsSuppressedByRule(findings, "missing-response-time")).
		To(gomega.BeEmpty(), "a suppression without a reason is inert")
	g.Expect(ruleIDs(findings)).To(gomega.ContainElement("missing-response-time"),
		"the finding must survive an inert suppression")
}

// The rules are substring-based, so suppression-comment TEXT itself counts as
// content: a reason mentioning "response time" satisfies missing-response-time
// before suppression is even consulted. Pinned here so a future structure-aware
// parser (ROADMAP) knows the incumbent behavior it must stay compatible with.
func TestValidate_suppression_text_counts_as_content(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(
		t,
	) + "\n<!-- securitymd:ignore(missing-response-time) we respond within 48 hours per wiki -->\n"

	findings := validateContent(t, content)

	g.Expect(ruleIDs(findings)).NotTo(gomega.ContainElement("missing-response-time"),
		"the reason text satisfies the substring rule, leaving nothing to suppress")
}

func TestValidate_suppression_unknown_rule_is_inert(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := compliantPolicy(t) + "\n<!-- securitymd:ignore(not-a-real-rule) typo must fail safe -->\n"

	findings := validateContent(t, content)

	g.Expect(ruleIDs(findings)).NotTo(gomega.ContainElement("not-a-real-rule"))

	for _, f := range findings {
		g.Expect(f.Suppression).To(gomega.BeNil(), "no real rule may be suppressed by an unknown rule name")
	}
}

func TestValidate_suppression_multiple_rules(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

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

	g.Expect(findingsSuppressedByRule(findings, "missing-versions")).To(gomega.HaveLen(1))
	g.Expect(findingsSuppressedByRule(findings, "missing-response-time")).To(gomega.HaveLen(1))
	g.Expect(findingsSuppressedByRule(findings, "missing-versions")[0].Suppression.Reason).
		To(gomega.Equal("tracked in the docs issue"))
}

func TestValidate_missing_file_cannot_be_suppressed(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	// missing-file has no file content to carry a comment: the rule must not
	// appear in the suppressible set, so a config can never silence it.
	g.Expect(KnownRuleIDs()).NotTo(gomega.ContainElement(finding.RuleName("suppressible-nonexistent")))

	findings, err := Detect(withWorkingDir(t.Context(), t.TempDir()))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.HaveLen(1))
	g.Expect(findings[0].Suppression).To(gomega.BeNil())
}

func TestValidate_suppression_until_keeps_evidence(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(t) +
		"\n<!-- securitymd:ignore(missing-response-time) until 2999-01-01 debt tracked in the roadmap -->\n"

	findings := validateContent(t, content)

	suppressed := findingsSuppressedByRule(findings, "missing-response-time")
	g.Expect(suppressed).To(gomega.HaveLen(1), "an unexpired until-suppression still silences the rule")

	g.Expect(suppressed[0].Suppression.Reason).To(gomega.Equal("debt tracked in the roadmap"))
	g.Expect(suppressed[0].Suppression.ExpiresAt).NotTo(gomega.BeNil())
	g.Expect(*suppressed[0].Suppression.ExpiresAt).
		To(gomega.Equal(time.Date(2999, 1, 2, 0, 0, 0, 0, time.UTC)),
			"until grants the whole given UTC day: expiry is the next midnight")
}

func TestValidate_suppression_expired_carries_no_metadata(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(t) +
		"\n<!-- securitymd:ignore(missing-response-time) until 2020-01-01 long-past deferral -->\n"

	findings := validateContent(t, content)

	g.Expect(ruleIDs(findings)).To(gomega.ContainElement("missing-response-time"),
		"an expired suppression must not silence anything")

	for _, f := range findings {
		g.Expect(f.Suppression).To(gomega.BeNil(),
			"an expired directive attaches nothing: every consumer must see plain debt")
	}
}

func TestValidate_suppression_until_requires_reason(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(t) + "\n<!-- securitymd:ignore(missing-response-time) until 2999-01-01 -->\n"

	findings := validateContent(t, content)

	g.Expect(findingsSuppressedByRule(findings, "missing-response-time")).
		To(gomega.BeEmpty(), "an until-clause does not replace the required reason")
}

func TestValidate_suppression_malformed_until_is_inert(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(t) +
		"\n<!-- securitymd:ignore(missing-response-time) until 17-05-2030 wrong date order -->\n"

	findings := validateContent(t, content)

	g.Expect(findingsSuppressedByRule(findings, "missing-response-time")).
		To(gomega.BeEmpty(), "an unparsable date must fail safe: no suppression, not an indefinite one")
	g.Expect(ruleIDs(findings)).To(gomega.ContainElement("missing-response-time"))
}

// TestSuppressionDirective_fullDayGrant pins the boundary directly: the
// suppression holds through the last nanosecond of the until-day and is gone
// the moment the next day starts.
func TestSuppressionDirective_fullDayGrant(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	directive, ok := parseSuppressionReason("until 2030-05-17 tracked debt")
	g.Expect(ok).To(gomega.BeTrue())
	g.Expect(directive.expiresAt).NotTo(gomega.BeNil())

	lastMomentOfDay := time.Date(2030, 5, 17, 23, 59, 59, 999999999, time.UTC)
	g.Expect(directive.activeAt(lastMomentOfDay)).
		To(gomega.BeTrue(), "active through the end of the until-day")

	firstMomentAfter := time.Date(2030, 5, 18, 0, 0, 0, 0, time.UTC)
	g.Expect(directive.activeAt(firstMomentAfter)).
		To(gomega.BeFalse(), "expired the instant the next day begins")
}

func findingsSuppressedByRule(findings []finding.Finding, rule string) []finding.Finding {
	return lo.Filter(findings, func(f finding.Finding, _ int) bool {
		return string(f.Rule) == rule && f.Suppression != nil
	})
}
