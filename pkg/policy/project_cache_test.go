package policy

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatestTag_returns_tag_and_caches_per_directory(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, runGitIn(t, dir, "init"))
	require.NoError(t, runGitIn(t, dir, "commit", "--allow-empty", "-m", "init"))
	require.NoError(t, runGitIn(t, dir, "tag", "v9.9.9"))

	ctx := t.Context()

	assert.Equal(t, "v9.9.9", LatestTag(ctx, dir), "first call resolves the tag")

	canceled, cancel := context.WithCancel(ctx)
	cancel()

	assert.Equal(t, "v9.9.9", LatestTag(canceled, dir),
		"cached hit must survive a dead context: a fresh git exec would fail on it")

	other := t.TempDir()
	assert.Equal(t, "", LatestTag(canceled, other),
		"a different directory must not inherit the cache entry")
}

func TestLatestTag_empty_when_no_tags(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, runGitIn(t, dir, "init"))
	require.NoError(t, runGitIn(t, dir, "commit", "--allow-empty", "-m", "init"))

	assert.Equal(t, "", LatestTag(t.Context(), dir))
}

func runGitIn(t *testing.T, dir string, args ...string) error {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_AUTHOR_NAME=test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)

	return cmd.Run()
}
