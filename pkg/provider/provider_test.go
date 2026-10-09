package provider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The spec must be registered (package var side effect) and valid: this is
// the exact contract BuildFlow's ToolsFromSDK conversion enforces.
func TestProvider_registered_and_shaped(t *testing.T) {
	t.Parallel()

	require.NotNil(t, Provider)
	assert.Equal(t, "securitymd", Provider.Name)
	assert.NotEmpty(t, Provider.Description)
	assert.NotNil(t, Provider.Detect)
	assert.NotNil(t, Provider.Repair, "securitymd must repair a missing policy by generating it")
	assert.NotNil(t, Provider.Trigger.Files,
		"the trigger must gate on project manifests, not run everywhere")
	assert.Contains(t, Provider.Trigger.Files, "SECURITY.md")

	require.NoError(t, Provider.ValidateOptions(toolsdk.OptionValues{}))
	require.ErrorIs(t,
		Provider.ValidateOptions(toolsdk.OptionValues{"unknown-knob": true}),
		toolsdk.ErrUnknownOption)
}

// TestProvider_suppressed_findings_never_reach_the_gate pins the provider's
// exit-neutral contract across the BuildFlow boundary: BuildFlow's findings
// gate counts the provider result, so an in-file suppression must strip the
// finding HERE — the gate cannot be trusted to honor suppression metadata
// (verified 2026-10-09: BuildFlow's filterFindingsAtOrAbove filters on
// severity only).
func TestProvider_suppressed_findings_never_reach_the_gate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SECURITY.md"),
		[]byte(flawedWithSuppressedContact(t)), 0o600))

	findings, err := Provider.Detect.Detect(finding.WithWorkingDir(t.Context(), dir))
	require.NoError(t, err)

	assert.NotContains(t, rulesOf(findings), "missing-contact",
		"a suppressed finding must not reach consumers' findings gates")
	assert.ElementsMatch(t, []string{"missing-response-time", "unresolved-template"},
		rulesOf(findings),
		"active findings must pass through untouched")
}

// flawedWithSuppressedContact derives the shared flawed fixture with an
// in-file suppression on its error-severity rule.
func flawedWithSuppressedContact(t *testing.T) string {
	t.Helper()

	raw, err := os.ReadFile(filepath.Join("..", "policy", "testdata", "policy", "flawed.md"))
	require.NoError(t, err, "shared flawed fixture missing")

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

	assert.Contains(t, Provider.Trigger.Files, "README.md")
	assert.Contains(t, Provider.Inputs, "README.md")
}

func TestProvider_detect_then_repair_then_verify(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := finding.WithWorkingDir(t.Context(), dir)

	// Detect: missing file is an error finding.
	findings, err := Provider.Detect.Detect(ctx)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, finding.RuleName("missing-file"), findings[0].Rule)

	// Repair: generates the policy (identity from options-less context needs
	// a git remote, which temp dirs lack — so pass the org/repo via the
	// generate path indirectly is not possible; repair must skip honestly).
	result, err := Provider.Repair.Repair(ctx)
	require.NoError(t, err)
	assert.Contains(t, result.Description, "organization/repository",
		"without a derivable identity, repair must explain rather than fabricate")

	// Verify the anti-lie loop end-to-end with a real git remote.
	gitRepo := t.TempDir()
	require.NoError(t, runGit(gitRepo, "init"))
	require.NoError(t, runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git"))

	gitCtx := finding.WithWorkingDir(t.Context(), gitRepo)

	findings, err = Provider.Detect.Detect(gitCtx)
	require.NoError(t, err)
	require.Len(t, findings, 1)

	result, err = Provider.Repair.Repair(gitCtx)
	require.NoError(t, err)
	assert.Contains(t, result.Description, "created")

	// Verify: BuildFlow re-runs Detect and measures the delta itself.
	findings, err = Provider.Detect.Detect(gitCtx)
	require.NoError(t, err)
	assert.Empty(t, findings, "after repair, detection must be clean (the verify step)")
}

func TestProvider_repair_respects_dry_run(t *testing.T) {
	t.Parallel()

	gitRepo := t.TempDir()
	require.NoError(t, runGit(gitRepo, "init"))
	require.NoError(t, runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git"))

	ctx := finding.WithWorkingDir(toolsdk.WithDryRun(t.Context(), true), gitRepo)

	result, err := Provider.Repair.Repair(ctx)
	require.NoError(t, err)
	assert.Contains(t, result.Description, "dry-run")
	assert.NoFileExists(t, filepath.Join(gitRepo, "SECURITY.md"))
}

func TestProvider_repair_never_touches_existing_policy(t *testing.T) {
	t.Parallel()

	gitRepo := t.TempDir()
	require.NoError(t, runGit(gitRepo, "init"))
	require.NoError(t, runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git"))

	existing := filepath.Join(gitRepo, "SECURITY.md")
	require.NoError(t, writeFile(existing, "# Hand-written policy\n"))

	ctx := finding.WithWorkingDir(t.Context(), gitRepo)

	result, err := Provider.Repair.Repair(ctx)
	require.NoError(t, err)
	assert.Contains(t, result.Description, "never overwrites")
	assert.Equal(t, "# Hand-written policy\n", mustRead(existing))
}

func TestProvider_contact_email_option_flows_into_repair(t *testing.T) {
	t.Parallel()

	gitRepo := t.TempDir()
	require.NoError(t, runGit(gitRepo, "init"))
	require.NoError(t, runGit(gitRepo, "remote", "add", "origin",
		"https://github.com/AcmeCorp/widget.git"))

	ctx := finding.WithWorkingDir(
		toolsdk.WithOptions(t.Context(), toolsdk.OptionValues{"contact-email": "security@acme.com"}),
		gitRepo,
	)

	_, err := Provider.Repair.Repair(ctx)
	require.NoError(t, err)
	assert.Contains(t, mustRead(filepath.Join(gitRepo, "SECURITY.md")), "security@acme.com")
}

func TestProvider_severity_overrides_downgrade_findings(t *testing.T) {
	t.Parallel()

	ctx := finding.WithWorkingDir(
		toolsdk.WithOptions(t.Context(), toolsdk.OptionValues{"severity-overrides": "missing-file=warning"}),
		t.TempDir(),
	)

	findings, err := Provider.Detect.Detect(ctx)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, finding.SeverityWarning, findings[0].Severity,
		"tool_options severity-overrides must reach the detector output")
}

func TestProvider_invalid_severity_override_fails_with_clear_error(t *testing.T) {
	t.Parallel()

	ctx := finding.WithWorkingDir(
		toolsdk.WithOptions(t.Context(), toolsdk.OptionValues{"severity-overrides": "missing-file=fatal"}),
		t.TempDir(),
	)

	_, err := Provider.Detect.Detect(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "severity-overrides")
	assert.Contains(t, err.Error(), "fatal", "the error must name the offending value")
}
