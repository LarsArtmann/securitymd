package acceptance

import (
	"os"

	"github.com/onsi/gomega"
)

func withTempFile(content string, fn func(path string)) {
	tmpFile, err := os.CreateTemp("", "SECURITY.md")
	expectNoError(err)
	defer os.Remove(tmpFile.Name())

	expectNoError2(tmpFile.WriteString(content))
	expectNoError(tmpFile.Close())

	fn(tmpFile.Name())
}

func expectNoError(err error) {
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func expectNoError2(n int, err error) {
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
