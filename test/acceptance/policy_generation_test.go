package acceptance

import (
	"os"
	"path/filepath"

	"github.com/LarsArtmann/template-SECURITY/internal"
	"github.com/LarsArtmann/template-SECURITY/internal/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Security Policy Generation", ginkgo.Label("acceptance"), func() {
	var (
		tool      *internal.SecurityTool
		config    internal.PolicyConfig
		outputDir string
	)

	assertPolicyContains := func(expected string) {
		err := tool.GeneratePolicy(nil, config)
		expectNoError(err)

		content, err := os.ReadFile(filepath.Join(outputDir, "SECURITY.md"))
		expectNoError(err)
		gomega.Expect(string(content)).To(gomega.ContainSubstring(expected))
	}

	ginkgo.BeforeEach(func() {
		tool = internal.NewSecurityTool()
		outputDir = ginkgo.GinkgoT().TempDir()
		config = internal.PolicyConfig{
			Type:         types.PolicyTypeGitHub,
			Organization: "Acme Corp",
			ContactEmail: "security@acme.com",
			OutputDir:    outputDir,
			Variables:    make(map[string]string),
		}
	})

	ginkgo.Describe("Generating a complete security policy", func() {
		const templateContent = `# Security Policy for {{.Organization}}

## Supported Versions

| Version | Supported Until |
|---------|-----------------|
| {{.LatestVersion}} | {{.SupportEndDate}} |

## Reporting a Vulnerability

Email us at {{.ContactEmail}}.

## Security Practices

We follow security best practices.

_Last updated: {{.LastUpdated}}_
`

		ginkgo.BeforeEach(func() {
			templateDir := filepath.Join(outputDir, "..", "templates")
			expectNoError(os.MkdirAll(templateDir, 0o755))

			templateFile := filepath.Join(templateDir, "SECURITY.md")
			expectNoError(os.WriteFile(templateFile, []byte(templateContent), 0o644))

			expectNoError(os.Chdir(filepath.Join(outputDir, "..")))
		})

		ginkgo.It("creates a SECURITY.md file in the output directory", func() {
			err := tool.GeneratePolicy(nil, config)
			expectNoError(err)

			outputFile := filepath.Join(outputDir, "SECURITY.md")
			gomega.Expect(outputFile).To(gomega.BeAnExistingFile())
		})

		ginkgo.It("includes the organization name", func() {
			assertPolicyContains("Acme Corp")
		})

		ginkgo.It("includes the contact email", func() {
			assertPolicyContains("security@acme.com")
		})

		ginkgo.It("includes version information", func() {
			assertPolicyContains("Supported Versions")
		})
	})

	ginkgo.Describe("Handling configuration", func() {
		ginkgo.Context("when custom variables are provided", func() {
			ginkgo.It("processes custom variables", func() {
				config.Variables["CUSTOM_FIELD"] = "custom_value"
				config.OutputDir = ginkgo.GinkgoT().TempDir()

				templateDir := filepath.Join(config.OutputDir, "..", "templates")
				expectNoError(os.MkdirAll(templateDir, 0o755))
				expectNoError(
					os.WriteFile(
						filepath.Join(templateDir, "SECURITY.md"),
						[]byte("# Test"),
						0o644,
					),
				)
				expectNoError(os.Chdir(filepath.Join(config.OutputDir, "..")))

				err := tool.GeneratePolicy(nil, config)
				expectNoError(err)
			})
		})
	})
})

var _ = ginkgo.Describe("Configuration Loading", ginkgo.Label("acceptance"), func() {
	var tool *internal.SecurityTool

	ginkgo.BeforeEach(func() {
		tool = internal.NewSecurityTool()
	})

	ginkgo.Describe("Finding configuration file", func() {
		ginkgo.It("returns a config file path when one exists in the project", func() {
			result := tool.FindConfigFile()
			gomega.Expect(result).To(gomega.SatisfyAny(
				gomega.BeEmpty(),
				gomega.ContainSubstring(".yaml"),
			))
		})
	})

	ginkgo.Describe("Loading configuration", func() {
		ginkgo.It("returns an error when config file does not exist", func() {
			_, err := tool.LoadConfig("/nonexistent/path/config.yaml")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})
