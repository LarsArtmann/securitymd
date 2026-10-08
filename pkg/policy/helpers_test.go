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

// writePolicy places the valid fixture policy at path (creating parents).
func writePolicy(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(validPolicy), 0o600)
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}
