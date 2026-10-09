package policy

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/require"
)

// withWorkingDir pins Detect/Generate to dir, mirroring how BuildFlow scopes
// providers via finding.WithWorkingDir.
func withWorkingDir(ctx context.Context, dir string) context.Context {
	return finding.WithWorkingDir(ctx, dir)
}

// writePolicy places the canonical compliant fixture at path (creating
// parents).
func writePolicyFixture(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, []byte(compliantPolicy(t)), 0o600))
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}
