package policy

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readmeRuleRow matches a rule-table row's leading `rule-id` cell.
var readmeRuleRow = regexp.MustCompile(`^\|\s*` + "`" + `([a-z-]+)` + "`")

// codeRuleSeverities assembles the code's single source of truth: every rule
// ID securitymd can emit mapped to the severity it emits with.
func codeRuleSeverities() map[finding.RuleName]finding.Severity {
	severities := make(map[finding.RuleName]finding.Severity, len(sectionRules)+len(contentRuleIDs)+1)

	for _, rule := range sectionRules {
		severities[rule.rule] = rule.severity
	}

	for rule, severity := range contentRuleSeverities {
		severities[rule] = severity
	}

	severities[RuleMissingFile] = missingFileSeverity

	return severities
}

// readmeDocumentedSeverities parses README.md's two rule tables into
// rule → documented severity ("error" or "warning" from the table heading).
func readmeDocumentedSeverities(t *testing.T) map[finding.RuleName]finding.Severity {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	require.NoError(t, err, "README.md must stay readable from the drift guard")

	documented := make(map[finding.RuleName]finding.Severity)
	section := ""

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "### ") {
			section = line

			continue
		}

		match := readmeRuleRow.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		switch {
		case strings.Contains(section, "(error)"):
			documented[finding.RuleName(match[1])] = finding.SeverityError
		case strings.Contains(section, "(warning)"):
			documented[finding.RuleName(match[1])] = finding.SeverityWarning
		}
	}

	return documented
}

// TestREADME_rule_tables_match_code pins README's rule tables to the code's
// rule table: a rule added, removed, or re-severitied in code without a
// README update fails here, and vice versa. This is the guard for the
// exit-2-class failure mode where docs describe a tool that no longer exists.
func TestREADME_rule_tables_match_code(t *testing.T) {
	t.Parallel()

	code := codeRuleSeverities()
	documented := readmeDocumentedSeverities(t)

	require.NotEmpty(t, documented, "the README rule tables must be found")

	assert.Equal(t, code, documented,
		"README rule tables and code severity tables have drifted; "+
			"update README.md and the rule tables together")
}

func TestKnownRuleIDs_matches_severity_table_keys(t *testing.T) {
	t.Parallel()

	assert.ElementsMatch(t, KnownRuleIDs(), mapKeys(codeRuleSeverities()),
		"the severity table must cover exactly the emittable rules")
}

func mapKeys(severities map[finding.RuleName]finding.Severity) []finding.RuleName {
	keys := make([]finding.RuleName, 0, len(severities))
	for rule := range severities {
		keys = append(keys, rule)
	}

	return keys
}
