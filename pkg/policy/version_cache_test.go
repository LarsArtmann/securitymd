package policy

import (
	"os"
	"os/exec"
	"testing"

	"github.com/onsi/gomega"
)

// runGitIsolated runs git with fleet config neutralized so tests never see
// the developer's hooks or aliases.
func runGitIsolated(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), "git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_AUTHOR_NAME=test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)

	out, err := cmd.CombinedOutput()
	gomega.NewWithT(t).Expect(err).NotTo(gomega.HaveOccurred(), "git %v: %s", args, out)
}

// TestLookupLatestTag_caches_per_directory proves the memoization
// behaviorally: after the cache is warmed, deleting the tag does not change
// the answer for that directory, while a fresh directory re-reads git. No
// call-counting seam is needed — and none would be race-free under
// t.Parallel anyway.
func TestLookupLatestTag_caches_per_directory(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	repo := t.TempDir()
	runGitIsolated(t, repo, "init")
	runGitIsolated(t, repo, "commit", "--allow-empty", "-m", "init")
	runGitIsolated(t, repo, "tag", "v9.9.9")

	ctx := t.Context()

	g.Expect(lookupLatestTag(ctx, repo)).To(gomega.Equal("v9.9.9"), "first lookup reads git")

	runGitIsolated(t, repo, "tag", "-d", "v9.9.9")

	g.Expect(lookupLatestTag(ctx, repo)).To(gomega.Equal("v9.9.9"),
		"second lookup must hit the cache, not re-run git describe")

	untagged := t.TempDir()
	runGitIsolated(t, untagged, "init")
	runGitIsolated(t, untagged, "commit", "--allow-empty", "-m", "init")

	g.Expect(lookupLatestTag(ctx, untagged)).To(gomega.BeEmpty(),
		"a different directory must not inherit another directory's cache entry")
}

func TestVersionCell_falls_back_to_placeholder(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	untagged := t.TempDir()
	runGitIsolated(t, untagged, "init")
	runGitIsolated(t, untagged, "commit", "--allow-empty", "-m", "init")

	g.Expect(versionCell(t.Context(), untagged)).To(gomega.Equal("Latest release"))
}

func TestVersionCell_caches_empty_result_for_tagless_repo(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	untagged := t.TempDir()
	runGitIsolated(t, untagged, "init")
	runGitIsolated(t, untagged, "commit", "--allow-empty", "-m", "init")

	ctx := t.Context()
	g.Expect(versionCell(ctx, untagged)).To(gomega.Equal("Latest release"))

	runGitIsolated(t, untagged, "tag", "v1.0.0")

	g.Expect(versionCell(ctx, untagged)).To(gomega.Equal("Latest release"),
		"the tagless cache entry must hold: repeated renders stay stable within a process")
}

func TestGenerate_uses_latest_tag_in_version_cell(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	repo := t.TempDir()
	runGitIsolated(t, repo, "init")
	runGitIsolated(t, repo, "commit", "--allow-empty", "-m", "init")
	runGitIsolated(t, repo, "tag", "v3.2.1")

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    repo,
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(result.Wrote).To(gomega.BeTrue())

	content, err := os.ReadFile(result.Path)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(content).To(gomega.ContainSubstring("v3.2.1"),
		"the version cell must carry the repo's latest tag")
}
