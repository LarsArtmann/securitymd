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
	"github.com/onsi/gomega"
)

var updateGolden = flag.Bool("update", false, "overwrite golden files with current output")

// pathHashPattern matches the hash segment go-finding derives from the file
// path in line-0 finding IDs. The hash is stable per path, but golden
// fixtures run under t.TempDir, so the value itself must not be pinned.
var pathHashPattern = regexp.MustCompile(`securitymd:([a-z-]+):[0-9a-f]{16}`)

// flawedPolicy and its contract live in fixtures_test.go / testdata/policy.

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
				gomega.NewWithT(t).Expect(err).NotTo(gomega.HaveOccurred())

				return ctx
			},
		},
		{
			name: "flawed",
			setup: func(t *testing.T, dir string) context.Context {
				t.Helper()

				path := filepath.Join(dir, "SECURITY.md")
				gomega.NewWithT(t).Expect(os.WriteFile(path, []byte(flawedPolicy(t)), 0o600)).To(gomega.Succeed())

				return withWorkingDir(t.Context(), dir)
			},
		},
		{
			name: "suppressed",
			setup: func(t *testing.T, dir string) context.Context {
				t.Helper()

				// The flawed policy plus an in-file suppression for its
				// error-severity contact rule: the finding stays visible with
				// suppression metadata in JSON and SARIF, while
				// unresolved-template remains active (the exit gate).
				suppressed := strings.Replace(
					flawedPolicy(t),
					"Please email {{.ContactEmail}}",
					"<!-- securitymd:ignore(missing-contact) organization email pending -->\nPlease email {{.ContactEmail}}",
					1,
				)

				path := filepath.Join(dir, "SECURITY.md")
				gomega.NewWithT(t).Expect(os.WriteFile(path, []byte(suppressed), 0o600)).To(gomega.Succeed())

				return withWorkingDir(t.Context(), dir)
			},
		},
		{
			name: "all-suppressed",
			setup: func(t *testing.T, dir string) context.Context {
				t.Helper()

				// Every finding the flawed fixture produces is silenced by
				// one unexpired until-suppression: the escape hatch at full
				// extent. JSON pins the evidence (reason + expiry); SARIF
				// pins that a fully-suppressed report exports zero results.
				allSuppressed := strings.Replace(
					flawedPolicy(t),
					"Please email {{.ContactEmail}}",
					"<!-- securitymd:ignore(missing-contact,missing-response-time,unresolved-template) until 2999-01-01 deferred debt with an expiry -->\nPlease email {{.ContactEmail}}",
					1,
				)

				path := filepath.Join(dir, "SECURITY.md")
				gomega.NewWithT(t).Expect(os.WriteFile(path, []byte(allSuppressed), 0o600)).To(gomega.Succeed())

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
			gomega.NewWithT(t).Expect(err).NotTo(gomega.HaveOccurred())

			assertGolden(t, test.name+".json", normalizeOutput(t, dir, renderReportJSON(t, report)))
			assertGolden(t, test.name+".sarif", normalizeOutput(t, dir, renderReportSARIF(t, ctx, report)))
		})
	}
}

func renderReportJSON(t *testing.T, report *finding.Report) string {
	t.Helper()

	var buf bytes.Buffer
	gomega.NewWithT(t).Expect(report.WriteJSON(&buf)).To(gomega.Succeed())

	return buf.String()
}

func renderReportSARIF(t *testing.T, ctx context.Context, report *finding.Report) string {
	t.Helper()

	var buf bytes.Buffer
	gomega.NewWithT(t).Expect(report.WriteSARIFWithOpts(ctx, &buf)).To(gomega.Succeed())

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
	g := gomega.NewWithT(t)

	golden := filepath.Join("testdata", "golden", name)
	if *updateGolden {
		g.Expect(os.MkdirAll(filepath.Dir(golden), 0o750)).To(gomega.Succeed())
		g.Expect(os.WriteFile(golden, []byte(got), 0o600)).To(gomega.Succeed())

		return
	}

	want, err := os.ReadFile(golden)
	g.Expect(err).NotTo(gomega.HaveOccurred(),
		"golden file missing; regenerate with go test ./pkg/policy -run TestReport_golden -update")
	g.Expect(got).To(gomega.Equal(string(want)),
		"output contract drifted; inspect the diff, or if intended refresh with -update")
}
