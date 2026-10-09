package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain re-executes the compiled test binary as the real CLI: the child
// branch runs main() directly, so its os.Exit code is exactly what a shell
// (and therefore CI) would observe. This pins README.md's exit-code contract
// end-to-end instead of trusting that the wiring matches the docs.
func TestMain(m *testing.M) {
	if os.Getenv("SECURITYMD_CONTRACT_CHILD") == "1" {
		main()

		return
	}

	os.Exit(m.Run())
}

// repoWithPolicy writes content as the repo's SECURITY.md and returns the dir.
func repoWithPolicy(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	if content != "" {
		require.NoError(t, os.WriteFile(filepath.Join(dir, "SECURITY.md"), []byte(content), 0o600))
	}

	return dir
}

// fixturePolicy reads a canonical fixture from the shared testdata/policy
// directory so the contract tests validate the same bytes the unit and
// golden tests do.
func fixturePolicy(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "..", "pkg", "policy", "testdata", "policy", name))
	require.NoError(t, err, "policy fixture %s missing", name)

	return string(content)
}

// fixturePolicyOf adapts a named shared fixture into a policy resolver for
// the scenario table, keeping rows self-describing without sentinel strings.
func fixturePolicyOf(name string) func(t *testing.T) string {
	return func(t *testing.T) string {
		t.Helper()

		return fixturePolicy(t, name)
	}
}

// policyWithoutResponseTime strips the response commitment from the compliant
// fixture: exactly one warning remains, ready for severity escalation.
func policyWithoutResponseTime(t *testing.T) string {
	t.Helper()

	return strings.ReplaceAll(
		fixturePolicy(t, "compliant.md"),
		"We commit to an initial response within 48 hours.\n\n", "")
}

// allErrorsSuppressed suppresses every error-severity rule the flawed fixture
// violates, leaving only the response-time warning: the full escape hatch.
func allErrorsSuppressed(t *testing.T) string {
	t.Helper()

	suppressed := strings.Replace(fixturePolicy(t, "flawed.md"),
		"Please email {{.ContactEmail}}",
		"<!-- securitymd:ignore(missing-contact) organization email pending -->\nPlease email {{.ContactEmail}}", 1)

	return suppressed +
		"\n<!-- securitymd:ignore(unresolved-template) template rendered by CI on release -->\n"
}

// allFindingsSuppressedUntilFuture silences every flawed-fixture finding with
// one unexpired until-suppression: the escape hatch with a calendar bound.
func allFindingsSuppressedUntilFuture(t *testing.T) string {
	t.Helper()

	return strings.Replace(fixturePolicy(t, "flawed.md"),
		"Please email {{.ContactEmail}}",
		"<!-- securitymd:ignore(missing-contact,missing-response-time,unresolved-template) "+
			"until 2999-01-01 deferred debt with an expiry -->\nPlease email {{.ContactEmail}}", 1)
}

// contactSuppressedUntilPast pins expiry semantics end-to-end: an until-clause
// in the past is as if the comment were never written.
func contactSuppressedUntilPast(t *testing.T) string {
	t.Helper()

	return strings.Replace(fixturePolicy(t, "flawed.md"),
		"Please email {{.ContactEmail}}",
		"<!-- securitymd:ignore(missing-contact) until 2020-01-01 long-past deferral -->\nPlease email {{.ContactEmail}}", 1)
}

func runCLI(t *testing.T, dir string, args ...string) (int, string) {
	t.Helper()

	command := exec.CommandContext(t.Context(), os.Args[0], args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "SECURITYMD_CONTRACT_CHILD=1")

	output, err := command.CombinedOutput()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else {
		require.NoError(t, err)
	}

	return exitCode, string(output)
}

// TestCLI_exit_code_contract pins README.md:88 — `0` clean · `1`
// error-severity findings · `2` operational failure — across the surface CI
// scripts actually drive.
func TestCLI_exit_code_contract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		policy     func(t *testing.T) string
		args       []string
		wantExit   int
		wantOut    []string
		notWantOut []string
	}{
		{
			name:     "clean policy validates green",
			policy:   fixturePolicyOf("compliant.md"),
			args:     []string{"validate"},
			wantExit: 0,
		},
		{
			name:     "missing policy is a finding, not a crash",
			policy:   nil,
			args:     []string{"validate"},
			wantExit: 1,
			wantOut:  []string{"missing-file"},
		},
		{
			name:     "flawed policy is a finding",
			policy:   fixturePolicyOf("flawed.md"),
			args:     []string{"validate"},
			wantExit: 1,
			wantOut:  []string{"missing-contact", "unresolved-template"},
		},
		{
			name:     "JSON findings still exit 1",
			policy:   fixturePolicyOf("flawed.md"),
			args:     []string{"validate", "--format", "json"},
			wantExit: 1,
			wantOut:  []string{"\"rule\": \"missing-contact\""},
		},
		{
			name:     "suppressed error findings keep exit neutral",
			policy:   allErrorsSuppressed,
			args:     []string{"validate"},
			wantExit: 0,
			wantOut:  []string{"suppressed: organization email pending"},
		},
		{
			name:     "unexpired until-suppression keeps exit neutral",
			policy:   allFindingsSuppressedUntilFuture,
			args:     []string{"validate"},
			wantExit: 0,
			wantOut:  []string{"suppressed: deferred debt with an expiry"},
		},
		{
			name:       "expired until-suppression surfaces the finding",
			policy:     contactSuppressedUntilPast,
			args:       []string{"validate"},
			wantExit:   1,
			wantOut:    []string{"missing-contact"},
			notWantOut: []string{"suppressed"},
		},
		{
			name:     "critical escalation trips the threshold gate",
			policy:   policyWithoutResponseTime,
			args:     []string{"validate", "--set-severity", "missing-response-time=critical"},
			wantExit: 1,
		},
		{
			name:     "downgraded missing-file passes CI",
			policy:   nil,
			args:     []string{"validate", "--set-severity", "missing-file=warning"},
			wantExit: 0,
		},
		{
			name:     "setup refuses to overwrite an existing policy",
			policy:   fixturePolicyOf("compliant.md"),
			args:     []string{"setup"},
			wantExit: 0,
			wantOut:  []string{"already exists"},
		},
		{
			name:     "setup without identity or TTY skips honestly",
			policy:   nil,
			args:     []string{"setup"},
			wantExit: 0,
			wantOut:  []string{"could not derive"},
		},
		{
			name:     "explicit missing --file is operational failure",
			policy:   nil,
			args:     []string{"validate", "--file", "does-not-exist.md"},
			wantExit: 2,
		},
		{
			name:     "unknown validate --location is operational failure",
			policy:   nil,
			args:     []string{"validate", "--location", "bogus"},
			wantExit: 2,
			wantOut:  []string{"invalid --location"},
		},
		{
			name:     "unknown status --location is operational failure",
			policy:   nil,
			args:     []string{"status", "--location", "bogus"},
			wantExit: 2,
		},
		{
			name:     "status is informational on an empty repo",
			policy:   nil,
			args:     []string{"status"},
			wantExit: 0,
		},
		{
			name:     "status honors --location docs and stays informational",
			policy:   nil,
			args:     []string{"status", "--location", "docs"},
			wantExit: 0,
			wantOut:  []string{"missing-file"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			policyContent := ""
			if test.policy != nil {
				policyContent = test.policy(t)
			}

			dir := repoWithPolicy(t, policyContent)

			exitCode, output := runCLI(t, dir, test.args...)

			assert.Equal(t, test.wantExit, exitCode,
				"exit code must match the documented contract")
			for _, want := range test.wantOut {
				assert.Contains(t, output, want)
			}
			for _, notWant := range test.notWantOut {
				assert.NotContains(t, output, notWant)
			}
		})
	}
}
