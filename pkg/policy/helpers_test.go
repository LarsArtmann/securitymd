package policy

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/gomega"
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
	gomega.NewWithT(t).Expect(os.MkdirAll(filepath.Dir(path), 0o750)).To(gomega.Succeed())
	gomega.NewWithT(t).Expect(os.WriteFile(path, []byte(compliantPolicy(t)), 0o600)).To(gomega.Succeed())
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	gomega.NewWithT(t).Expect(err).NotTo(gomega.HaveOccurred())
}
