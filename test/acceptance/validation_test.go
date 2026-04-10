package acceptance

import (
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

	expectValidationFailsWithError := func(content, expectedErrorSubstring string) {
		withTempFile(content, func(path string) {
			result, err := validator.ValidateSECURITYMd(path)
			expectNoError(err)
			gomega.Expect(result.Valid).To(gomega.BeFalse())

			hasError := false

			for _, errMsg := range result.Errors {
				if strings.Contains(errMsg, expectedErrorSubstring) {
					hasError = true

					break
				}
			}

			gomega.Expect(hasError).To(gomega.BeTrue())
		})
	}

	createTempFileAndValidate := func(content string) *internal.SecurityValidationResult {
		var result *internal.SecurityValidationResult

		withTempFile(content, func(path string) {
			var err error

			result, err = validator.ValidateSECURITYMd(path)
			expectNoError(err)
		})

		return result
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
			result := createTempFileAndValidate(content)
			gomega.Expect(result.Valid).To(gomega.BeTrue())
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

We follow security best practices.
`,
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
			result := createTempFileAndValidate(content)

			hasShortWarning := false

			for _, warning := range result.Warnings {
				if strings.Contains(warning, "too short") {
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
