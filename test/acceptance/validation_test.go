package acceptance

import (
	"os"
	"strings"

	"github.com/LarsArtmann/template-SECURITY/internal"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Security Policy Validation", ginkgo.Label("acceptance"), func() {
	var validator *internal.SecurityValidator

	ginkgo.BeforeEach(func() {
		validator = internal.NewSecurityValidator()
	})

	ginkgo.Describe("Validating SECURITY.md completeness", func() {
		ginkgo.It("validates a complete security policy", func() {
			content := `# Security Policy

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| v2.x | 2026-12-31 |

## Reporting a Vulnerability

Email us at security@example.com with any security issues.

## Security Practices

We follow security best practices.

`
			tmpFile, err := os.CreateTemp("", "SECURITY.md")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			tmpFile.Close()

			result, err := validator.ValidateSECURITYMd(tmpFile.Name())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(result.Valid).To(gomega.BeTrue())
		})

		ginkgo.It("fails when reporting section is missing", func() {
			content := `# Security Policy

## Supported Versions

v1.0

## Security Practices

We follow best practices.

`
			tmpFile, err := os.CreateTemp("", "SECURITY.md")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			tmpFile.Close()

			result, err := validator.ValidateSECURITYMd(tmpFile.Name())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(result.Valid).To(gomega.BeFalse())
			gomega.Expect(strings.Join(result.Errors, "")).
				To(gomega.ContainSubstring("Reporting a Vulnerability"))
		})

		ginkgo.It("fails when contact email is missing", func() {
			content := `# Security Policy

## Supported Versions

v1.0

## Reporting a Vulnerability

Contact us through our website.

## Security Practices

We follow best practices.

`
			tmpFile, err := os.CreateTemp("", "SECURITY.md")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			tmpFile.Close()

			result, err := validator.ValidateSECURITYMd(tmpFile.Name())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(result.Valid).To(gomega.BeFalse())
			gomega.Expect(strings.Join(result.Errors, "")).
				To(gomega.ContainSubstring("contact email"))
		})
	})

	ginkgo.Describe("Detecting template variables", func() {
		ginkgo.It("fails when template variables are unresolved", func() {
			content := `# Security Policy

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| {{LATEST_VERSION}} | {{SUPPORT_END_DATE}} |

## Reporting a Vulnerability

Email us at {{CONTACT_EMAIL}}.

## Security Practices

We follow security best practices.
`
			tmpFile, err := os.CreateTemp("", "SECURITY.md")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			tmpFile.Close()

			result, err := validator.ValidateSECURITYMd(tmpFile.Name())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(result.Valid).To(gomega.BeFalse())

			hasTemplateError := false
			for _, errMsg := range result.Errors {
				if strings.Contains(errMsg, "template variable") {
					hasTemplateError = true
					break
				}
			}
			gomega.Expect(hasTemplateError).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("Checking content quality", func() {
		ginkgo.It("warns when content is too short", func() {
			content := `# Security Policy

Missing content.
`
			tmpFile, err := os.CreateTemp("", "SECURITY.md")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			tmpFile.Close()

			result, err := validator.ValidateSECURITYMd(tmpFile.Name())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			hasShortWarning := false
			for _, warning := range result.Warnings {
				if strings.Contains(warning, "too short") {
					hasShortWarning = true
					break
				}
			}
			gomega.Expect(hasShortWarning).To(gomega.BeTrue())
		})

		ginkgo.It("fails when policy lacks substantive content", func() {
			content := `# Security Policy

## Supported Versions

## Reporting a Vulnerability

## Security Practices

`
			tmpFile, err := os.CreateTemp("", "SECURITY.md")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			tmpFile.Close()

			result, err := validator.ValidateSECURITYMd(tmpFile.Name())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(result.Valid).To(gomega.BeFalse())

			hasContentError := false
			for _, errMsg := range result.Errors {
				if strings.Contains(errMsg, "substantive") {
					hasContentError = true
					break
				}
			}
			gomega.Expect(hasContentError).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("Handling non-existent files", func() {
		ginkgo.It("returns an error for non-existent files", func() {
			_, err := validator.ValidateSECURITYMd("/nonexistent/path/SECURITY.md")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})
