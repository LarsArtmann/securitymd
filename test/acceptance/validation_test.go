// Package acceptance provides BDD acceptance tests for securitymd.
package acceptance

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func withPolicyFile(content string, callback func(path string)) {
	dir := gomegaTempDir()
	path := filepath.Join(dir, "SECURITY.md")

	gomega.Expect(os.WriteFile(path, []byte(content), 0o600)).To(gomega.Succeed())

	callback(path)
}

func gomegaTempDir() string {
	dir, err := os.MkdirTemp("", "securitymd-acceptance-*")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return dir
}

func ruleNames(findings []finding.Finding) []string {
	names := make([]string, 0, len(findings))
	for _, f := range findings {
		names = append(names, string(f.Rule))
	}

	return names
}

var _ = ginkgo.Describe("SECURITY.md validation", ginkgo.Label("acceptance"), func() {
	const compliantPolicy = `# Security Policy

## Supported Versions

| Version | Supported Until |
| ------- | --------------- |
| v2.x    | 2026-12-31     |

## Reporting a Vulnerability

Email security@example.com for any security issue; we respond within 48 hours.

## Security Practices

All changes are reviewed, dependencies scanned, releases signed.
`

	ginkgo.It("passes a compliant policy", func() {
		withPolicyFile(compliantPolicy, func(path string) {
			findings, err := policy.Validate(path)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(findings).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("missing essentials", func() {
		ginkgo.It("reports a missing vulnerability-reporting section", func() {
			withPolicyFile("# Security Policy\n\nNothing to see here, just filler text for the file.\n", func(path string) {
				findings, err := policy.Validate(path)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(ruleNames(findings)).To(gomega.ContainElement("missing-reporting"))
			})
		})

		ginkgo.It("reports a policy with no contact channel", func() {
			broken := strings.ReplaceAll(compliantPolicy, "Email security@example.com for any security issue; we respond within 48 hours.",
				"Please find a way to contact us about problems.")

			withPolicyFile(broken, func(path string) {
				findings, err := policy.Validate(path)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(ruleNames(findings)).To(gomega.ContainElement("missing-contact"))
			})
		})

		ginkgo.It("flags leftover template variables", func() {
			broken := strings.ReplaceAll(compliantPolicy,
				"Email security@example.com for any security issue; we respond within 48 hours.",
				"Email {{.ContactEmail}} for any security issue; we respond within 48 hours.")

			withPolicyFile(broken, func(path string) {
				findings, err := policy.Validate(path)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(ruleNames(findings)).To(gomega.ContainElement("unresolved-template"))
			})
		})
	})
})

var _ = ginkgo.Describe("SECURITY.md generation", ginkgo.Label("acceptance"), func() {
	ginkgo.It("generates a policy that passes its own validator", func() {
		dir := gomegaTempDir()

		result, err := policy.Generate(ginkgo.GinkgoT().Context(), policy.GenerateOptions{
			Directory:    dir,
			Organization: "AcmeCorp",
			Repository:   "widget",
			ContactEmail: "security@acme.com",
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(result.Wrote).To(gomega.BeTrue())

		findings, err := policy.Validate(result.Path)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(findings).To(gomega.BeEmpty())
	})

	ginkgo.It("refuses to overwrite a hand-written policy", func() {
		dir := gomegaTempDir()
		existing := filepath.Join(dir, "SECURITY.md")
		gomega.Expect(os.WriteFile(existing, []byte("# Hand-written\n"), 0o600)).To(gomega.Succeed())

		result, err := policy.Generate(ginkgo.GinkgoT().Context(), policy.GenerateOptions{
			Directory:    dir,
			Organization: "AcmeCorp",
			Repository:   "widget",
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(result.Wrote).To(gomega.BeFalse())
		gomega.Expect(result.Description).To(gomega.ContainSubstring("never overwrites"))
	})

	ginkgo.It("detects a missing policy as an error finding with a fix", func() {
		dir := gomegaTempDir()

		findings, err := policy.Detect(finding.WithWorkingDir(ginkgo.GinkgoT().Context(), dir))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(findings).To(gomega.HaveLen(1))
		gomega.Expect(string(findings[0].Rule)).To(gomega.Equal("missing-file"))
		gomega.Expect(findings[0].Severity).To(gomega.Equal(finding.SeverityError))
		gomega.Expect(findings[0].Suggestion).NotTo(gomega.BeEmpty())
	})
})
