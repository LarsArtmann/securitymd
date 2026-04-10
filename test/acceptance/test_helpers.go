package acceptance

import (
	"os"

	"github.com/onsi/gomega"
)

func withTempFile(content string, fn func(path string)) {
	tmpFile, err := os.CreateTemp("", "SECURITY.md")
	if err != nil {
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}

	if err := tmpFile.Close(); err != nil {
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}

	fn(tmpFile.Name())
}

func expectNoError(err error) {
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
