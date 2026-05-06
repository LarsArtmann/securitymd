// Package acceptance provides test helper functions.
package acceptance

import (
	"os"

	"github.com/onsi/gomega"
)

func withTempFile(content string, callback func(path string)) {
	tmpFile, err := os.CreateTemp("", "SECURITY-*.md")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	_, err = tmpFile.WriteString(content)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = tmpFile.Close()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}
	}()

	callback(tmpFile.Name())
}

func expectNoError(err error) {
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
