// Package acceptance provides BDD acceptance tests for securitymd.
package acceptance

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "securitymd Acceptance Suite")
}
