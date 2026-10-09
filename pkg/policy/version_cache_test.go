package policy

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runGitIsolated runs git with fleet config neutralized so tests never see
// the developer's hooks or aliases.
func runGitIsolated(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
}

// TestLookupLatestTag_caches_per_directory proves the memoization
// behaviorally: after the cache is warmed, deleting the tag does not change
// the answer for that directory, while a fresh directory re-reads git. No
// call-counting seam is needed — and none would be race-free under
// t.Parallel anyway.
func TestLookupLatestTag_caches_per_directory(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	runGitIsolated(t, repo, "init")
	runGitIsolated(t, repo, "commit", "--allow-empty", "-m", "init")
	runGitIsolated(t, repo, "tag", "v9.9.9")

	ctx := t.Context()

	require.Equal(t, "v9.9.9", lookupLatestTag(ctx, repo), "first lookup reads git")

	runGitIsolated(t, repo, "tag", "-d", "v9.9.9")

	assert.Equal(t, "v9.9.9", lookupLatestTag(ctx, repo),
		"second lookup must hit the cache, not re-run git describe")

	untagged := t.TempDir()
	runGitIsolated(t, untagged, "init")
	runGitIsolated(t, untagged, "commit", "--allow-empty", "-m", "init")

	assert.Empty(t, lookupLatestTag(ctx, untagged),
		"a different directory must not inherit another directory's cache entry")
}

func TestVersionCell_falls_back_to_placeholder(t *testing.T) {
	t.Parallel()

	untagged := t.TempDir()
	runGitIsolated(t, untagged, "init")
	runGitIsolated(t, untagged, "commit", "--allow-empty", "-m", "init")

	assert.Equal(t, "Latest release", versionCell(t.Context(), untagged))
}

func TestVersionCell_caches_empty_result_for_tagless_repo(t *testing.T) {
	t.Parallel()

	untagged := t.TempDir()
	runGitIsolated(t, untagged, "init")
	runGitIsolated(t, untagged, "commit", "--allow-empty", "-m", "init")

	ctx := t.Context()
	require.Equal(t, "Latest release", versionCell(ctx, untagged))

	runGitIsolated(t, untagged, "tag", "v1.0.0")

	assert.Equal(t, "Latest release", versionCell(ctx, untagged),
		"the tagless cache entry must hold: repeated renders stay stable within a process")
}

func TestGenerate_uses_latest_tag_in_version_cell(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	runGitIsolated(t, repo, "init")
	runGitIsolated(t, repo, "commit", "--allow-empty", "-m", "init")
	runGitIsolated(t, repo, "tag", "v3.2.1")

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    repo,
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	require.NoError(t, err)
	require.True(t, result.Wrote)

	content, err := os.ReadFile(result.Path)
	require.NoError(t, err)
	assert.Contains(t, string(content), "v3.2.1",
		"the version cell must carry the repo's latest tag")
}
