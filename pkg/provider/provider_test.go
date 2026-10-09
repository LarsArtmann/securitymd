package provider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
	"github.com/onsi/gomega"
)

// The spec must be registered (package var side effect) and valid: this is
// the exact contract BuildFlow's ToolsFromSDK conversion enforces.
func TestProvider_registered_and_shaped(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	g.Expect(Provider).NotTo(gomega.BeNil())
	g.Expect(Provider.Name).To(gomega.Equal("securitymd"))
	g.Expect(Provider.Description).NotTo(gomega.BeEmpty())
	g.Expect(Provider.Detect).NotTo(gomega.BeNil())
	g.Expect(Provider.Repair).NotTo(gomega.BeNil(),
		"securitymd must repair a missing policy by generating it")
	g.Expect(Provider.Trigger.Files).NotTo(gomega.BeNil(),
		"the trigger must gate on project manifests, not run everywhere")
	g.Expect(Provider.Trigger.Files).To(gomega.ContainElement("SECURITY.md"))

	g.Expect(Provider.ValidateOptions(toolsdk.OptionValues{})).To(gomega.Succeed())
	g.Expect(Provider.ValidateOptions(toolsdk.OptionValues{"unknown-knob": true})).
		To(gomega.MatchError(toolsdk.ErrUnknownOption))
}

// TestProvider_suppressed_findings_never_reach_the_gate pins the provider's
// exit-neutral contract across the BuildFlow boundary: BuildFlow's findings
// gate counts the provider result, so an in-file suppression must strip the
// finding HERE — the gate cannot be trusted to honor suppression metadata
// (verified 2026-10-09: BuildFlow's filterFindingsAtOrAbove filters on
// severity only).
func TestProvider_suppressed_findings_never_reach_the_gate(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	g.Expect(os.WriteFile(filepath.Join(dir, "SECURITY.md"),
		[]byte(flawedWithSuppressedContact(t)), 0o600)).To(gomega.Succeed())

	findings, err := Provider.Detect.Detect(finding.WithWorkingDir(t.Context(), dir))
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(rulesOf(findings)).NotTo(gomega.ContainElement("missing-contact"),
		"a suppressed finding must not reach consumers' findings gates")
	g.Expect(rulesOf(findings)).
		To(gomega.ConsistOf("missing-response-time", "unresolved-template"),
			"active findings must pass through untouched")
}

// flawedWithSuppressedContact derives the shared flawed fixture with an
// in-file suppression on its error-severity rule.
func flawedWithSuppressedContact(t *testing.T) string {
	t.Helper()

	raw, err := os.ReadFile(filepath.Join("..", "policy", "testdata", "policy", "flawed.md"))
	gomega.NewWithT(t).Expect(err).NotTo(gomega.HaveOccurred(), "shared flawed fixture missing")

	return strings.Replace(string(raw),
		"Please email {{.ContactEmail}}",
		"<!-- securitymd:ignore(missing-contact) organization email pending -->\nPlease email {{.ContactEmail}}",
		1)
}

func rulesOf(findings []finding.Finding) []string {
	rules := make([]string, 0, len(findings))
	for _, f := range findings {
		rules = append(rules, string(f.Rule))
	}

	return rules
}

// A docs-only repo (README.md, no dependency manifests) must still activate:
// published projects without a policy are exactly the ones that need one.
func TestProvider_activates_on_docs_only_repos(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	g.Expect(Provider.Trigger.Files).To(gomega.ContainElement("README.md"))
	g.Expect(Provider.Inputs).To(gomega.ContainElement("README.md"))
}

func TestProvider_detect_then_repair_then_verify(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	ctx := finding.WithWorkingDir(t.Context(), dir)

	// Detect: missing file is an error finding.
	findings, err := Provider.Detect.Detect(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.HaveLen(1))
	g.Expect(findings[0].Rule).To(gomega.Equal(finding.RuleName("missing-file")))

	// Repair: generates the policy (identity from options-less context needs
	// a git remote, which temp dirs lack — so pass the org/repo via the
	// generate path indirectly is not possible; repair must skip honestly).
	result, err := Provider.Repair.Repair(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Description).To(gomega.ContainSubstring("organization/repository"),
		"without a derivable identity, repair must explain rather than fabricate")

	// Verify the anti-lie loop end-to-end with a real git remote.
	gitRepo := t.TempDir()
	g.Expect(runGit(gitRepo, "init")).To(gomega.Succeed())
	g.Expect(runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git")).To(gomega.Succeed())

	gitCtx := finding.WithWorkingDir(t.Context(), gitRepo)

	findings, err = Provider.Detect.Detect(gitCtx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.HaveLen(1))

	result, err = Provider.Repair.Repair(gitCtx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Description).To(gomega.ContainSubstring("created"))

	// Verify: BuildFlow re-runs Detect and measures the delta itself.
	findings, err = Provider.Detect.Detect(gitCtx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.BeEmpty(), "after repair, detection must be clean (the verify step)")
}

func TestProvider_repair_respects_dry_run(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	gitRepo := t.TempDir()
	g.Expect(runGit(gitRepo, "init")).To(gomega.Succeed())
	g.Expect(runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git")).To(gomega.Succeed())

	ctx := finding.WithWorkingDir(toolsdk.WithDryRun(t.Context(), true), gitRepo)

	result, err := Provider.Repair.Repair(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Description).To(gomega.ContainSubstring("dry-run"))
	g.Expect(filepath.Join(gitRepo, "SECURITY.md")).NotTo(gomega.BeAnExistingFile())
}

func TestProvider_repair_never_touches_existing_policy(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	gitRepo := t.TempDir()
	g.Expect(runGit(gitRepo, "init")).To(gomega.Succeed())
	g.Expect(runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git")).To(gomega.Succeed())

	existing := filepath.Join(gitRepo, "SECURITY.md")
	g.Expect(writeFile(existing, "# Hand-written policy\n")).To(gomega.Succeed())

	ctx := finding.WithWorkingDir(t.Context(), gitRepo)

	result, err := Provider.Repair.Repair(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Description).To(gomega.ContainSubstring("never overwrites"))
	g.Expect(mustRead(existing)).To(gomega.Equal("# Hand-written policy\n"))
}

func TestProvider_contact_email_option_flows_into_repair(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	gitRepo := t.TempDir()
	g.Expect(runGit(gitRepo, "init")).To(gomega.Succeed())
	g.Expect(runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git")).To(gomega.Succeed())

	ctx := finding.WithWorkingDir(
		toolsdk.WithOptions(t.Context(), toolsdk.OptionValues{"contact-email": "security@acme.com"}),
		gitRepo,
	)

	_, err := Provider.Repair.Repair(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(mustRead(filepath.Join(gitRepo, "SECURITY.md"))).To(gomega.ContainSubstring("security@acme.com"))
}

func TestProvider_severity_overrides_downgrade_findings(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	ctx := finding.WithWorkingDir(
		toolsdk.WithOptions(t.Context(), toolsdk.OptionValues{"severity-overrides": "missing-file=warning"}),
		t.TempDir(),
	)

	findings, err := Provider.Detect.Detect(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.HaveLen(1))
	g.Expect(findings[0].Severity).To(gomega.Equal(finding.SeverityWarning),
		"tool_options severity-overrides must reach the detector output")
}

func TestProvider_invalid_severity_override_fails_with_clear_error(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	ctx := finding.WithWorkingDir(
		toolsdk.WithOptions(t.Context(), toolsdk.OptionValues{"severity-overrides": "missing-file=fatal"}),
		t.TempDir(),
	)

	_, err := Provider.Detect.Detect(ctx)

	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("severity-overrides"))
	g.Expect(err.Error()).To(gomega.ContainSubstring("fatal"), "the error must name the offending value")
}
