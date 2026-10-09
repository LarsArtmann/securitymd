package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_force_regenerates_in_place_with_backup(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	require.NoError(t, os.WriteFile(existing, []byte("# Stale hand-written policy\n"), 0o600))

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		Force:        true,
	})
	require.NoError(t, err)

	assert.True(t, result.Wrote)
	assert.NotEmpty(t, result.BackupPath)
	assert.Contains(t, result.Description, "backed up")

	backup, err := os.ReadFile(result.BackupPath)
	require.NoError(t, err)
	assert.Equal(t, "# Stale hand-written policy\n", string(backup), "the backup must carry the previous policy")

	assert.Contains(t, mustRead(t, existing), "AcmeCorp/widget", "the policy must be regenerated in place")
}

func TestGenerate_force_dry_run_writes_nothing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	require.NoError(t, os.WriteFile(existing, []byte("# Hand-written\n"), 0o600))

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		Force:        true,
		DryRun:       true,
	})
	require.NoError(t, err)

	assert.False(t, result.Wrote)
	assert.Empty(t, result.BackupPath)
	assert.Contains(t, result.Description, "dry-run")
	assert.Equal(t, "# Hand-written\n", mustRead(t, existing))

	backups, _ := filepath.Glob(filepath.Join(dir, "*.bak"))
	assert.Empty(t, backups, "dry-run must not write a backup")
}

func TestGenerate_force_refuses_second_policy_location(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	require.NoError(t, os.WriteFile(existing, []byte("# Root policy\n"), 0o600))

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		Location:     LocationDocs,
		Force:        true,
	})
	require.NoError(t, err)

	assert.False(t, result.Wrote)
	assert.Contains(t, result.Description, "two policies")
	assert.Equal(t, "# Root policy\n", mustRead(t, existing))
	assert.NoFileExists(t, filepath.Join(dir, "docs", "SECURITY.md"))
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

			dir := t.TempDir()

			result, err := Generate(t.Context(), GenerateOptions{
				Directory:    dir,
				Organization: "AcmeCorp",
				Repository:   "widget",
				Location:     test.location,
			})
			require.NoError(t, err)
			require.True(t, result.Wrote)
			assert.FileExists(t, filepath.Join(dir, test.relPath))
		})
	}
}

func TestGenerate_unknown_location_is_an_error(t *testing.T) {
	t.Parallel()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    t.TempDir(),
		Organization: "AcmeCorp",
		Repository:   "widget",
		Location:     "nowhere",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown policy location")
	assert.Empty(t, result.Path)
}

func TestOrderedCandidates_prefers_canonical_location(t *testing.T) {
	t.Parallel()

	defaultOrder, err := OrderedCandidates("")
	require.NoError(t, err)
	assert.Equal(t, CandidateLocations, defaultOrder)

	githubOrder, err := OrderedCandidates(LocationGitHub)
	require.NoError(t, err)
	assert.Equal(t, []string{".github/SECURITY.md", "SECURITY.md", "docs/SECURITY.md"}, githubOrder)

	docsOrder, err := OrderedCandidates(LocationDocs)
	require.NoError(t, err)
	assert.Equal(t, []string{"docs/SECURITY.md", "SECURITY.md", ".github/SECURITY.md"}, docsOrder)

	_, err = OrderedCandidates("nowhere")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown policy location")
}
