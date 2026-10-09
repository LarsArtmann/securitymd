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
		policy     string
		args       []string
		wantExit   int
		wantOut    []string
		notWantOut []string
	}{
		{
			name:     "clean policy validates green",
			policy:   "fixture:compliant.md",
			args:     []string{"validate"},
			wantExit: 0,
		},
		{
			name:     "missing policy is a finding, not a crash",
			policy:   "",
			args:     []string{"validate"},
			wantExit: 1,
			wantOut:  []string{"missing-file"},
		},
		{
			name:     "flawed policy is a finding",
			policy:   "fixture:flawed.md",
			args:     []string{"validate"},
			wantExit: 1,
			wantOut:  []string{"missing-contact", "unresolved-template"},
		},
		{
			name:     "JSON findings still exit 1",
			policy:   "fixture:flawed.md",
			args:     []string{"validate", "--format", "json"},
			wantExit: 1,
			wantOut:  []string{"\"rule\": \"missing-contact\""},
		},
		{
			name:     "suppressed error findings keep exit neutral",
			policy:   "fixture-suppressed:all",
			args:     []string{"validate"},
			wantExit: 0,
			wantOut:  []string{"suppressed: organization email pending"},
		},
		{
			name:     "critical escalation trips the threshold gate",
			policy:   "fixture-derived:no-response",
			args:     []string{"validate", "--set-severity", "missing-response-time=critical"},
			wantExit: 1,
		},
		{
			name:     "downgraded missing-file passes CI",
			policy:   "",
			args:     []string{"validate", "--set-severity", "missing-file=warning"},
			wantExit: 0,
		},
		{
			name:     "explicit missing --file is operational failure",
			policy:   "",
			args:     []string{"validate", "--file", "does-not-exist.md"},
			wantExit: 2,
		},
		{
			name:     "unknown validate --location is operational failure",
			policy:   "",
			args:     []string{"validate", "--location", "bogus"},
			wantExit: 2,
			wantOut:  []string{"invalid --location"},
		},
		{
			name:     "unknown status --location is operational failure",
			policy:   "",
			args:     []string{"status", "--location", "bogus"},
			wantExit: 2,
		},
		{
			name:     "status is informational on an empty repo",
			policy:   "",
			args:     []string{"status"},
			wantExit: 0,
		},
		{
			name:     "status honors --location docs and stays informational",
			policy:   "",
			args:     []string{"status", "--location", "docs"},
			wantExit: 0,
			wantOut:  []string{"missing-file"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			policyContent := ""
			switch {
			case test.policy == "fixture:compliant.md":
				policyContent = fixturePolicy(t, "compliant.md")
			case test.policy == "fixture:flawed.md":
				policyContent = fixturePolicy(t, "flawed.md")
			case test.policy == "fixture-suppressed:all":
				policyContent = allErrorsSuppressed(t)
			case test.policy == "fixture-derived:no-response":
				policyContent = policyWithoutResponseTime(t)
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
