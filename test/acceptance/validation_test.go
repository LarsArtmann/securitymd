// Package acceptance provides tests for security policy validation.
package acceptance

import (
	"strings"

	"github.com/LarsArtmann/template-SECURITY/internal"
	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Security Policy Validation", ginkgo.Label("acceptance"), func() {
	var validator *internal.SecurityValidator

	ginkgo.BeforeEach(func() {
		validator = internal.NewSecurityValidator()
	})

	expectValidationFailsWithError := func(content, expectedErrorSubstring string) {
		withTempFile(content, func(path string) {
			report, err := validator.ValidateSECURITYMd(path)
			expectNoError(err)
			gomega.Expect(internal.ReportIsValid(report)).To(gomega.BeFalse())

			hasError := false

			for _, f := range report.FindingsSnapshot() {
				if f.Severity == finding.SeverityError &&
					strings.Contains(f.Message, expectedErrorSubstring) {
					hasError = true

					break
				}
			}

			gomega.Expect(hasError).To(gomega.BeTrue())
		})
	}

	createTempFileAndValidate := func(content string) *finding.Report {
		var report *finding.Report

		withTempFile(content, func(path string) {
			var err error

			report, err = validator.ValidateSECURITYMd(path)
			expectNoError(err)
		})

		return report
	}

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
			report := createTempFileAndValidate(content)
			gomega.Expect(internal.ReportIsValid(report)).To(gomega.BeTrue())
		})

		ginkgo.It("fails validation for incomplete policies", func() {
			testCases := []struct {
				name        string
				content     string
				errContains string
			}{
				{
					name: "missing reporting section",
					content: `# Security Policy

## Supported Versions

v1.0

## Security Practices

We follow best practices.

`,
					errContains: "Reporting a Vulnerability",
				},
				{
					name: "missing contact email",
					content: `# Security Policy

## Supported Versions

v1.0

## Reporting a Vulnerability

Contact us through our website.

## Security Practices

We follow best practices.

`,
					errContains: "contact email",
				},
				{
					name: "unresolved template variables",
					content: `# Security Policy

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| {{LATEST_VERSION}} | {{SUPPORT_END_DATE}} |

## Reporting a Vulnerability

Email us at {{CONTACT_EMAIL}}.

## Security Practices

We follow security best practices.`,
					errContains: "template variable",
				},
				{
					name: "lacks substantive content",
					content: `# Security Policy

## Supported Versions

## Reporting a Vulnerability

## Security Practices

`,
					errContains: "substantive",
				},
			}

			for _, tc := range testCases {
				ginkgo.By(tc.name, func() {
					expectValidationFailsWithError(tc.content, tc.errContains)
				})
			}
		})
	})

	ginkgo.Describe("Checking content quality", func() {
		ginkgo.It("warns when content is too short", func() {
			content := `# Security Policy

Missing content.
`
			report := createTempFileAndValidate(content)

			hasShortWarning := false

			for _, f := range report.FindingsSnapshot() {
				if f.Severity == finding.SeverityWarning &&
					strings.Contains(f.Message, "too short") {
					hasShortWarning = true

					break
				}
			}

			gomega.Expect(hasShortWarning).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("Handling non-existent files", func() {
		ginkgo.It("returns an error for non-existent files", func() {
			_, err := validator.ValidateSECURITYMd("/nonexistent/path/SECURITY.md")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})
