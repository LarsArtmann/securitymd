package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

// policyFixture loads a canonical SECURITY.md scenario from testdata/policy.
// Unit, golden, suppression, and acceptance tests all read these bytes, so a
// contract change lands everywhere at once instead of drifting across four
// hand-rolled string constants.
func policyFixture(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("testdata", "policy", name))
	gomega.NewWithT(t).Expect(err).NotTo(gomega.HaveOccurred(),
		"policy fixture %s missing from testdata/policy", name)

	return string(content)
}

// compliantPolicy is the zero-findings baseline every other fixture mutates.
func compliantPolicy(t *testing.T) string {
	t.Helper()

	return policyFixture(t, "compliant.md")
}

// flawedPolicy violates two missing-section rules with different severities
// (missing-contact error, missing-response-time warning) and carries an
// unresolved template variable, exercising both go-finding ID formats:
// path-hash (line 0) and line-precise.
func flawedPolicy(t *testing.T) string {
	t.Helper()

	return policyFixture(t, "flawed.md")
}

// policyWithoutResponseTime drops the response commitment from the compliant
// fixture: exactly one finding (missing-response-time) remains to suppress.
func policyWithoutResponseTime(t *testing.T) string {
	t.Helper()

	return strings.ReplaceAll(
		compliantPolicy(t),
		"We commit to an initial response within 48 hours.\n\n", "")
}

// TestPolicyFixtures_contract pins what each fixture means: the compliant
// fixture must satisfy the validator, and the flawed fixture must violate
// exactly the rules the golden and mutation tests rely on.
func TestPolicyFixtures_contract(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	g.Expect(ruleIDs(validateContent(t, compliantPolicy(t)))).
		To(gomega.BeEmpty(), "compliant.md is the zero-findings baseline")

	g.Expect(ruleIDs(validateContent(t, flawedPolicy(t)))).
		To(gomega.Equal([]string{"missing-contact", "missing-response-time", "unresolved-template"}),
			"flawed.md must violate exactly these rules, in emission order")
}

func TestPolicyWithoutResponseTime_derives_exactly_one_finding(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	g.Expect(ruleIDs(validateContent(t, policyWithoutResponseTime(t)))).
		To(gomega.Equal([]string{"missing-response-time"}),
			"the derivation must leave exactly the suppressible warning")
}
