package policy

import (
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "overwrite golden files with current output")

// pathHashPattern matches the hash segment go-finding derives from the file
// path in line-0 finding IDs. The hash is stable per path, but golden
// fixtures run under t.TempDir, so the value itself must not be pinned.
var pathHashPattern = regexp.MustCompile(`securitymd:([a-z-]+):[0-9a-f]{16}`)

// flawedPolicy deliberately violates two missing-section rules with different
// severities (missing-contact error, missing-response-time warning) and
// carries an unresolved template variable, exercising both go-finding ID
// formats: path-hash (line 0) and line-precise.
const flawedPolicy = `# Security Policy

## Supported Versions

| Version | Supported Until |
| ------- | --------------- |
| v2.x    | 2026-12-31     |

Only the latest release receives security fixes.

## Reporting a Vulnerability

Please email {{.ContactEmail}} with any security issues.

## Security Practices

We follow security best practices in development and operations.

All changes are reviewed before merge and CI runs security scanning.
`

func TestReport_golden_output_contract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T, dir string) context.Context
	}{
		{
			name: "clean",
			setup: func(t *testing.T, dir string) context.Context {
				t.Helper()

				ctx := withWorkingDir(t.Context(), dir)
				_, err := Generate(ctx, GenerateOptions{Organization: "AcmeCorp", Repository: "widget"})
				require.NoError(t, err)

				return ctx
			},
		},
		{
			name: "flawed",
			setup: func(t *testing.T, dir string) context.Context {
				t.Helper()

				path := filepath.Join(dir, "SECURITY.md")
				require.NoError(t, os.WriteFile(path, []byte(flawedPolicy), 0o600))

				return withWorkingDir(t.Context(), dir)
			},
		},
		{
			name: "missing-file",
			setup: func(t *testing.T, dir string) context.Context {
				t.Helper()

				return withWorkingDir(t.Context(), dir)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			ctx := test.setup(t, dir)

			report, err := Report(ctx)
			require.NoError(t, err)

			assertGolden(t, test.name+".json", normalizeOutput(t, dir, renderReportJSON(t, report)))
			assertGolden(t, test.name+".sarif", normalizeOutput(t, dir, renderReportSARIF(t, ctx, report)))
		})
	}
}

func renderReportJSON(t *testing.T, report *finding.Report) string {
	t.Helper()

	var buf bytes.Buffer
	require.NoError(t, report.WriteJSON(&buf))

	return buf.String()
}

func renderReportSARIF(t *testing.T, ctx context.Context, report *finding.Report) string {
	t.Helper()

	var buf bytes.Buffer
	require.NoError(t, report.WriteSARIFWithOpts(ctx, &buf))

	return buf.String()
}

// normalizeOutput strips run-volatile entropy: the fixture's absolute
// directory (position files, the missing-file message) and the ID hash
// derived from it. Everything else — rule IDs, severities, messages, SARIF
// shape, ID format — is pinned byte-for-byte.
func normalizeOutput(t *testing.T, dir, output string) string {
	t.Helper()

	normalized := strings.ReplaceAll(output, dir, "<REPO>")

	return pathHashPattern.ReplaceAllString(normalized, "securitymd:$1:<PATHHASH>")
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()

	golden := filepath.Join("testdata", "golden", name)
	if *updateGolden {
		require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o750))
		require.NoError(t, os.WriteFile(golden, []byte(got), 0o600))

		return
	}

	want, err := os.ReadFile(golden)
	require.NoError(t, err, "golden file missing; regenerate with go test ./pkg/policy -run TestReport_golden -update")
	assert.Equal(t, string(want), got, "output contract drifted; inspect the diff, or if intended refresh with -update")
}
