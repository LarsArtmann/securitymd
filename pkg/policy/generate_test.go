package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega"
)

func TestGenerate_creates_policy_and_reports(t *testing.T) {
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
	g.Expect(result.Path).To(gomega.Equal(filepath.Join(dir, "SECURITY.md")))
	g.Expect(result.Path).To(gomega.BeAnExistingFile())

	content := mustRead(t, result.Path)
	g.Expect(content).To(gomega.ContainSubstring("security@acme.com"))
	g.Expect(content).To(gomega.ContainSubstring("AcmeCorp/widget"))
	g.Expect(content).NotTo(gomega.ContainSubstring("{{"), "no template variables may leak")
}

func TestGenerate_advisory_only_when_no_email(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Wrote).To(gomega.BeTrue())

	content := mustRead(t, result.Path)
	g.Expect(content).To(gomega.ContainSubstring("security/advisories/new"))
}

func TestGenerate_dry_run_writes_nothing(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		DryRun:       true,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(result.Wrote).To(gomega.BeFalse())
	g.Expect(filepath.Join(dir, "SECURITY.md")).NotTo(gomega.BeAnExistingFile())
	g.Expect(result.Description).To(gomega.ContainSubstring("dry-run"))
}

func TestGenerate_never_overwrites_existing_policy(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	g.Expect(os.WriteFile(existing, []byte("# Our own policy\n"), 0o600)).To(gomega.Succeed())

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(result.Wrote).To(gomega.BeFalse())
	g.Expect(result.Description).To(gomega.ContainSubstring("never overwrites"))
	g.Expect(mustRead(t, existing)).To(gomega.Equal("# Our own policy\n"))
}

func TestGenerate_skips_without_identity(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{Directory: dir})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(result.Wrote).To(gomega.BeFalse())
	g.Expect(result.Description).To(gomega.ContainSubstring("organization/repository"))
	g.Expect(filepath.Join(dir, "SECURITY.md")).NotTo(gomega.BeAnExistingFile())
}

func TestGenerate_uses_working_dir_from_context(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	ctx := withWorkingDir(t.Context(), dir)

	result, err := Generate(ctx, GenerateOptions{
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Wrote).To(gomega.BeTrue())
	g.Expect(result.Path).To(gomega.Equal(filepath.Join(dir, "SECURITY.md")))
}
