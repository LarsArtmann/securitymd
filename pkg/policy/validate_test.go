package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/gomega"
)

// compliantPolicy lives in testdata/policy/compliant.md (see fixtures_test.go):
// one canonical baseline shared by unit, golden, suppression, and acceptance
// tests.

func validateContent(t *testing.T, content string) []finding.Finding {
	t.Helper()
	g := gomega.NewWithT(t)

	path := filepath.Join(t.TempDir(), "SECURITY.md")
	g.Expect(os.WriteFile(path, []byte(content), 0o600)).To(gomega.Succeed())

	findings, err := Validate(path)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	return findings
}

func ruleIDs(findings []finding.Finding) []string {
	ids := make([]string, 0, len(findings))
	for _, f := range findings {
		ids = append(ids, string(f.Rule))
	}

	return ids
}

func TestValidate_accepts_compliant_policy(t *testing.T) {
	t.Parallel()

	findings := validateContent(t, compliantPolicy(t))
	gomega.NewWithT(t).Expect(findings).
		To(gomega.BeEmpty(), "a compliant policy must yield zero findings")
}

func TestValidate_reports_each_missing_section(t *testing.T) {
	t.Parallel()

	compliant := compliantPolicy(t)

	tests := []struct {
		name string
		// remove strips one section from the valid policy
		remove    string
		expectIDs []string
	}{
		{"header", "# Security Policy\n", []string{"missing-header"}},
		{"reporting", "## Reporting a Vulnerability", []string{"missing-reporting"}},
		{"versions", "## Supported Versions", []string{"missing-versions"}},
		{"practices", "## Security Practices", []string{"missing-practices"}},
		{"contact", "security@example.com", []string{"missing-contact"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			g := gomega.NewWithT(t)

			broken := strings.ReplaceAll(compliant, test.remove, "")
			findings := validateContent(t, broken)

			for _, expected := range test.expectIDs {
				g.Expect(ruleIDs(findings)).To(gomega.ContainElement(expected))
			}
		})
	}
}

func TestValidate_accepts_advisory_link_as_contact(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	content := strings.ReplaceAll(compliantPolicy(t),
		"Email us at security@example.com with any security issues.",
		"Report via https://github.com/AcmeCorp/widget/security/advisories/new privately.")

	findings := validateContent(t, content)
	g.Expect(ruleIDs(findings)).NotTo(gomega.ContainElement("missing-contact"))
}

func TestValidate_flags_unresolved_template_variables_with_line(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	findings := validateContent(t, "# Security Policy\n\nContact {{.ContactEmail}} for issues.\n")

	var unresolved []finding.Finding

	for _, candidate := range findings {
		if candidate.Rule == "unresolved-template" {
			unresolved = append(unresolved, candidate)
		}
	}

	g.Expect(unresolved).To(gomega.HaveLen(1))
	g.Expect(unresolved[0].Position.Line).To(gomega.Equal(3),
		"unresolved-template must point at the exact offending line")
	g.Expect(unresolved[0].Severity).To(gomega.Equal(finding.SeverityError))
}

func TestValidate_flags_thin_and_placeholder_content(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		content   string
		expectIDs []string
	}{
		{
			"too short",
			"# Security Policy\n",
			[]string{"too-short", "no-content", "no-version-info"},
		},
		{
			"placeholder only",
			"# Security Policy\n\n## Reporting a Vulnerability\n\nContact security@example.com about anything.\n\n## Supported Versions\n\nSee versions.\n\n## Security Practices\n\nFill this in later.\n",
			[]string{"no-content"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			g := gomega.NewWithT(t)

			findings := validateContent(t, test.content)
			g.Expect(ruleIDs(findings)).To(gomega.ContainElements(test.expectIDs))
		})
	}
}

func TestValidate_generated_template_passes(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		ContactEmail: "security@acme.com",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Wrote).To(gomega.BeTrue())

	findings := validateContent(t, mustRead(t, result.Path))
	g.Expect(findings).To(gomega.BeEmpty(),
		"the embedded template must satisfy its own validator (dogfood invariant)")
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	g := gomega.NewWithT(t)

	content, err := os.ReadFile(path)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	return string(content)
}
