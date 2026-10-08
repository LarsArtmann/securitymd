package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_creates_policy_and_reports(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		ContactEmail: "security@acme.com",
	})
	require.NoError(t, err)

	assert.True(t, result.Wrote)
	assert.Equal(t, filepath.Join(dir, "SECURITY.md"), result.Path)
	assert.FileExists(t, result.Path)

	content := mustRead(t, result.Path)
	assert.Contains(t, content, "security@acme.com")
	assert.Contains(t, content, "AcmeCorp/widget")
	assert.NotContains(t, content, "{{", "no template variables may leak")
}

func TestGenerate_advisory_only_when_no_email(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	require.NoError(t, err)
	require.True(t, result.Wrote)

	content := mustRead(t, result.Path)
	assert.Contains(t, content, "security/advisories/new")
}

func TestGenerate_dry_run_writes_nothing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
		DryRun:       true,
	})
	require.NoError(t, err)

	assert.False(t, result.Wrote)
	assert.NoFileExists(t, filepath.Join(dir, "SECURITY.md"))
	assert.Contains(t, result.Description, "dry-run")
}

func TestGenerate_never_overwrites_existing_policy(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "SECURITY.md")
	require.NoError(t, os.WriteFile(existing, []byte("# Our own policy\n"), 0o600))

	result, err := Generate(t.Context(), GenerateOptions{
		Directory:    dir,
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	require.NoError(t, err)

	assert.False(t, result.Wrote)
	assert.Contains(t, result.Description, "never overwrites")
	assert.Equal(t, "# Our own policy\n", mustRead(t, existing))
}

func TestGenerate_skips_without_identity(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	result, err := Generate(t.Context(), GenerateOptions{Directory: dir})
	require.NoError(t, err)

	assert.False(t, result.Wrote)
	assert.Contains(t, result.Description, "organization/repository")
	assert.NoFileExists(t, filepath.Join(dir, "SECURITY.md"))
}

func TestGenerate_uses_working_dir_from_context(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := withWorkingDir(t.Context(), dir)

	result, err := Generate(ctx, GenerateOptions{
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	require.NoError(t, err)
	require.True(t, result.Wrote)
	assert.Equal(t, filepath.Join(dir, "SECURITY.md"), result.Path)
}
