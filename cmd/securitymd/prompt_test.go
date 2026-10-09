package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptIdentity_reads_both_coordinates(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer

	org, repo, err := promptIdentity(strings.NewReader("AcmeCorp\nwidget\n"), &out)
	require.NoError(t, err)
	assert.Equal(t, "AcmeCorp", org)
	assert.Equal(t, "widget", repo)
	assert.Contains(t, out.String(), "Organization")
	assert.Contains(t, out.String(), "Repository")
}

func TestPromptIdentity_trims_whitespace(t *testing.T) {
	t.Parallel()

	org, repo, err := promptIdentity(strings.NewReader("  AcmeCorp  \n\twidget\n"), &bytes.Buffer{})
	require.NoError(t, err)
	assert.Equal(t, "AcmeCorp", org)
	assert.Equal(t, "widget", repo)
}

func TestPromptIdentity_accepts_final_line_without_newline(t *testing.T) {
	t.Parallel()

	org, repo, err := promptIdentity(strings.NewReader("AcmeCorp\nwidget"), &bytes.Buffer{})
	require.NoError(t, err)
	assert.Equal(t, "AcmeCorp", org)
	assert.Equal(t, "widget", repo)
}

func TestPromptIdentity_rejects_empty_answers_with_flag_guidance(t *testing.T) {
	t.Parallel()

	_, _, err := promptIdentity(strings.NewReader("\nwidget\n"), &bytes.Buffer{})
	require.ErrorIs(t, err, errEmptyPromptAnswer)
	assert.Contains(t, err.Error(), "--organization",
		"the error must point non-interactive users at the flags")

	_, _, err = promptIdentity(strings.NewReader("AcmeCorp\n   \n"), &bytes.Buffer{})
	require.ErrorIs(t, err, errEmptyPromptAnswer)
	assert.Contains(t, err.Error(), "--repository")
}

func TestPromptIdentity_rejects_immediate_eof(t *testing.T) {
	t.Parallel()

	_, _, err := promptIdentity(strings.NewReader(""), &bytes.Buffer{})
	require.Error(t, err)
}

// In tests and CI, stdin is never a char device: the prompt helper must stay
// silent and leave the flow to the honest skip (never block a pipeline).
func TestPromptWhenIdentityMissing_is_inert_without_tty(t *testing.T) {
	t.Parallel()

	cmd := newSetupCmd()
	cmd.SetContext(t.Context())
	require.NoError(t, cmd.Flags().Set("dry-run", "false"))

	require.NoError(t, promptWhenIdentityMissing(cmd))
	assert.Empty(t, setupOrganization, "no prompt may run without an interactive terminal")
	assert.Empty(t, setupRepository)
}
