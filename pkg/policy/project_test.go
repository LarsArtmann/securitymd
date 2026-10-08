package policy

import (
	"path/filepath"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitRemote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		remote string
		want   RepoIdentity
	}{
		{"https://github.com/AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"https://github.com/AcmeCorp/widget", RepoIdentity{"AcmeCorp", "widget"}},
		{"git@github.com:AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"ssh://git@github.com/AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"https://gitlab.com/group/subproject/repo.git", RepoIdentity{"subproject", "repo"}},
		{"not-a-remote", RepoIdentity{}},
		{"", RepoIdentity{}},
	}

	for _, test := range tests {
		t.Run(test.remote, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, parseGitRemote(test.remote))
		})
	}
}

func TestRepoIdentity_IsComplete(t *testing.T) {
	t.Parallel()

	assert.True(t, RepoIdentity{"a", "b"}.IsComplete())
	assert.False(t, RepoIdentity{"a", ""}.IsComplete())
	assert.False(t, RepoIdentity{}.IsComplete())
}

func TestDetect_reports_missing_file_as_direct_fix(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := withWorkingDir(t.Context(), dir)

	findings, err := Detect(ctx)
	require.NoError(t, err)
	require.Len(t, findings, 1)

	f := findings[0]
	assert.Equal(t, finding.RuleName("missing-file"), f.Rule)
	assert.Equal(t, finding.SeverityError, f.Severity)
	assert.Equal(t, finding.FixStrategySuggest, f.FixStrategy,
		"no git remote in a temp dir: the fix cannot carry content, so it stays suggest-only")
	assert.NotEmpty(t, f.Suggestion)
}

func TestDetect_validates_existing_policy(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := withWorkingDir(t.Context(), dir)

	_, err := Generate(ctx, GenerateOptions{
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	require.NoError(t, err)

	findings, err := Detect(ctx)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestDetect_finds_policy_in_github_dir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, writePolicy(filepath.Join(dir, ".github", "SECURITY.md")))

	findings, err := Detect(withWorkingDir(t.Context(), dir))
	require.NoError(t, err)
	assert.Empty(t, findings, "a compliant .github/SECURITY.md must be discovered and pass")
}
