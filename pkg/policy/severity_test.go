package policy

import (
	"testing"
	"time"

	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/gomega"
)

func TestKnownRuleIDs_covers_every_rule(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

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
		g.Expect(ids).To(gomega.ContainElement(finding.RuleName(expected)))
	}

	g.Expect(ids).To(gomega.HaveLen(11),
		"a new rule must be added to the known-ID table in the same change")
}

func TestParseSeverityOverrides(t *testing.T) {
	t.Parallel()

	t.Run("parses a downgrade spec", func(t *testing.T) {
		t.Parallel()
		g := gomega.NewWithT(t)

		overrides, err := ParseSeverityOverrides("missing-file=warning, too-short=info")
		g.Expect(err).NotTo(gomega.HaveOccurred())

		g.Expect(overrides).To(gomega.Equal(SeverityOverrides{
			"missing-file": finding.SeverityWarning,
			"too-short":    finding.SeverityInfo,
		}))
	})

	t.Run("empty spec means no overrides", func(t *testing.T) {
		t.Parallel()
		g := gomega.NewWithT(t)

		overrides, err := ParseSeverityOverrides("")
		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(overrides).To(gomega.BeEmpty())
	})

	t.Run("rejects unknown rule with the known list", func(t *testing.T) {
		t.Parallel()
		g := gomega.NewWithT(t)

		_, err := ParseSeverityOverrides("missing-fiel=warning")

		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(err.Error()).To(gomega.ContainSubstring(`unknown rule "missing-fiel"`))
		g.Expect(err.Error()).To(gomega.ContainSubstring("missing-file"),
			"the error must name the accepted values")
	})

	t.Run("rejects invalid severity", func(t *testing.T) {
		t.Parallel()
		g := gomega.NewWithT(t)

		_, err := ParseSeverityOverrides("missing-file=fatal")

		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(err.Error()).To(gomega.ContainSubstring(`invalid severity "fatal"`))
		g.Expect(err.Error()).To(gomega.ContainSubstring("missing-file"))
	})

	t.Run("rejects malformed part", func(t *testing.T) {
		t.Parallel()
		g := gomega.NewWithT(t)

		_, err := ParseSeverityOverrides("missing-file-warning")

		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(err.Error()).To(gomega.ContainSubstring("want rule=severity"))
	})
}

func TestApplySeverityOverrides_downgrades_missing_file(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	findings, err := Detect(withWorkingDir(t.Context(), t.TempDir()))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.HaveLen(1))
	g.Expect(findings[0].Severity).To(gomega.Equal(finding.SeverityError))

	overrides, err := ParseSeverityOverrides("missing-file=warning")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	adjusted := ApplySeverityOverrides(findings, overrides)
	g.Expect(adjusted).To(gomega.HaveLen(1))
	g.Expect(adjusted[0].Severity).To(gomega.Equal(finding.SeverityWarning),
		"the adoption unblocker: missing-file downgraded to warning")
	g.Expect(adjusted[0].Suppression).To(gomega.BeNil(),
		"severity overrides must not masquerade as suppressions")
}

func TestApplySeverityOverrides_empty_is_noop(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	findings := validateContent(t, policyWithoutResponseTime(t))
	g.Expect(findings).NotTo(gomega.BeEmpty())

	g.Expect(ApplySeverityOverrides(findings, SeverityOverrides{})).
		To(gomega.Equal(findings), "an empty override map must return the findings unchanged")
}

// The CLI composes the two escape hatches: Validate applies in-file
// suppressions first, then the --set-severity overrides rewrite severities on
// the result. This pins the combination on the same rule — the edge the
// report flagged as unproven.
func TestSeverityOverrides_combine_with_suppressions(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := policyWithoutResponseTime(t) +
		"\n<!-- securitymd:ignore(missing-response-time) tracked in the support wiki -->\n"

	findings := validateContent(t, content)
	g.Expect(findings).To(gomega.HaveLen(1))
	g.Expect(findings[0].Suppression).NotTo(gomega.BeNil(), "precondition: the warning is suppressed")

	overrides, err := ParseSeverityOverrides("missing-response-time=critical")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	adjusted := ApplySeverityOverrides(findings, overrides)
	g.Expect(adjusted).To(gomega.HaveLen(1))

	g.Expect(adjusted[0].Severity).To(gomega.Equal(finding.SeverityCritical),
		"the override rewrites severity even on a suppressed finding")
	g.Expect(adjusted[0].Suppression).NotTo(gomega.BeNil(),
		"the override must not strip the suppression evidence")
	g.Expect(adjusted[0].IsSuppressedAt(time.Now())).
		To(gomega.BeTrue(), "severity escalation must not resurrect a suppressed finding")
}

// The gate contract the CLI implements: suppressed findings never trip exit 1,
// even when their severity was escalated by an override — and the same
// escalation without a suppression DOES trip it, so the suppression is what
// keeps CI green, not the override being harmless.
func TestSeverityOverrides_suppressed_findings_stay_exit_neutral(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	overrides, err := ParseSeverityOverrides("missing-response-time=critical")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	activeErrors := func(findings []finding.Finding) []finding.Finding {
		active := finding.Filter(findings, func(f finding.Finding) bool {
			return !f.IsSuppressedAt(time.Now())
		})

		return finding.Filter(active, finding.BySeverityAtLeast(finding.SeverityError))
	}

	suppressed := validateContent(t, policyWithoutResponseTime(t)+
		"\n<!-- securitymd:ignore(missing-response-time) tracked in the support wiki -->\n")
	g.Expect(activeErrors(ApplySeverityOverrides(suppressed, overrides))).
		To(gomega.BeEmpty(), "an escalated-but-suppressed finding must not activate: the escape hatch holds")

	unsuppressed := validateContent(t, policyWithoutResponseTime(t))
	g.Expect(activeErrors(ApplySeverityOverrides(unsuppressed, overrides))).
		NotTo(gomega.BeEmpty(), "the same escalation without an in-file suppression activates")
}
