package acceptance

import (
	"os"

	"github.com/onsi/gomega"
)

func withTempFile(content string, fn func(path string)) {
	tmpFile, err := os.CreateTemp("", "SECURITY.md")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(tmpFile.Close()).To(gomega.Succeed())

	fn(tmpFile.Name())
}
