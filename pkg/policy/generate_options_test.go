package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega"
)

func TestGenerate_force_regenerates_in_place_with_backup(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	g.Expect(os.WriteFile(existing, []byte("# Stale hand-written policy\n"), 0o600)).To(gomega.Succeed())

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		Force:        true,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(result.Wrote).To(gomega.BeTrue())
	g.Expect(result.BackupPath).NotTo(gomega.BeEmpty())
	g.Expect(result.Description).To(gomega.ContainSubstring("backed up"))

	backup, err := os.ReadFile(result.BackupPath)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(backup).To(gomega.Equal([]byte("# Stale hand-written policy\n")),
		"the backup must carry the previous policy")

	g.Expect(mustRead(t, existing)).To(gomega.ContainSubstring("AcmeCorp/widget"),
		"the policy must be regenerated in place")
}

func TestGenerate_force_dry_run_writes_nothing(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	g.Expect(os.WriteFile(existing, []byte("# Hand-written\n"), 0o600)).To(gomega.Succeed())

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		Force:        true,
		DryRun:       true,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(result.Wrote).To(gomega.BeFalse())
	g.Expect(result.BackupPath).To(gomega.BeEmpty())
	g.Expect(result.Description).To(gomega.ContainSubstring("dry-run"))
	g.Expect(mustRead(t, existing)).To(gomega.Equal("# Hand-written\n"))

	backups, _ := filepath.Glob(filepath.Join(dir, "*.bak"))
	g.Expect(backups).To(gomega.BeEmpty(), "dry-run must not write a backup")
}

func TestGenerate_force_refuses_second_policy_location(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	g.Expect(os.WriteFile(existing, []byte("# Root policy\n"), 0o600)).To(gomega.Succeed())

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		Location:     LocationDocs,
		Force:        true,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(result.Wrote).To(gomega.BeFalse())
	g.Expect(result.Description).To(gomega.ContainSubstring("two policies"))
	g.Expect(mustRead(t, existing)).To(gomega.Equal("# Root policy\n"))
	g.Expect(filepath.Join(dir, "docs", "SECURITY.md")).NotTo(gomega.BeAnExistingFile())
}

func TestGenerate_location_writes_canonical_target(t *testing.T) {
	t.Parallel()

	tests := []struct {
		location string
		relPath  string
	}{
		{LocationRoot, "SECURITY.md"},
		{LocationGitHub, filepath.Join(".github", "SECURITY.md")},
		{LocationDocs, filepath.Join("docs", "SECURITY.md")},
	}

	for _, test := range tests {
		t.Run(test.location, func(t *testing.T) {
			t.Parallel()

			g := gomega.NewWithT(t)

			dir := t.TempDir()

			result, err := Generate(t.Context(), GenerateOptions{
				Directory:    dir,
				Organization: "AcmeCorp",
				Repository:   "widget",
				Location:     test.location,
			})
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(result.Wrote).To(gomega.BeTrue())
			g.Expect(filepath.Join(dir, test.relPath)).To(gomega.BeAnExistingFile())
		})
	}
}

func TestGenerate_unknown_location_is_an_error(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    t.TempDir(),
		Organization: "AcmeCorp",
		Repository:   "widget",
		Location:     "nowhere",
	})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("unknown policy location"))
	g.Expect(result.Path).To(gomega.BeEmpty())
}

func TestOrderedCandidates_prefers_canonical_location(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	defaultOrder, err := OrderedCandidates("")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(defaultOrder).To(gomega.Equal(CandidateLocations))

	githubOrder, err := OrderedCandidates(LocationGitHub)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(githubOrder).To(gomega.Equal([]string{".github/SECURITY.md", "SECURITY.md", "docs/SECURITY.md"}))

	docsOrder, err := OrderedCandidates(LocationDocs)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(docsOrder).To(gomega.Equal([]string{"docs/SECURITY.md", "SECURITY.md", ".github/SECURITY.md"}))

	_, err = OrderedCandidates("nowhere")
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("unknown policy location"))
}
