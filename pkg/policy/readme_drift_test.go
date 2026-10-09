package policy

import (
	"maps"
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

	maps.Copy(severities, contentRuleSeverities)

	severities[RuleMissingFile] = missingFileSeverity

	return severities
}

// readmeDocumentedSeverities parses a README's two rule tables into
// rule → documented severity ("error" or "warning" from the table heading).
func readmeDocumentedSeverities(t *testing.T, path string) map[finding.RuleName]finding.Severity {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err, "README must stay readable from the drift guard")

	documented := make(map[finding.RuleName]finding.Severity)
	section := ""

	for line := range strings.SplitSeq(string(content), "\n") {
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

// ruleTableDrift reports the first difference between the code's severity
// table and the documented one; "" means they match. Both directions are
// checked: an undocumented code rule and a phantom documented rule both
// count as drift.
func ruleTableDrift(code, documented map[finding.RuleName]finding.Severity) string {
	for rule, severity := range code {
		if documented[rule] != severity {
			return fmt.Sprintf("rule %s: code emits %s, README documents %s",
				rule, severity, documented[rule])
		}
	}

	for rule := range documented {
		if _, emitted := code[rule]; !emitted {
			return "README documents a rule the code never emits: " + string(rule)
		}
	}

	return ""
}

// TestREADME_rule_tables_match_code pins README's rule tables to the code's
// rule table: a rule added, removed, or re-severitied in code without a
// README update fails here, and vice versa. This is the guard for the
// exit-2-class failure mode where docs describe a tool that no longer exists.
func TestREADME_rule_tables_match_code(t *testing.T) {
	t.Parallel()

	documented := readmeDocumentedSeverities(t, filepath.Join("..", "..", "README.md"))

	require.NotEmpty(t, documented, "the README rule tables must be found")

	assert.Empty(t, ruleTableDrift(codeRuleSeverities(), documented),
		"README rule tables and code severity tables have drifted; "+
			"update README.md and the rule tables together")
}

// TestREADME_drift_guard_red_run is the guard's own red run: the guard only
// earns trust if a mutated README (and a mutated severity table) demonstrably
// fails it. Without this, a broken parser or a vacuous comparator could sit
// green forever while docs drift.
func TestREADME_drift_guard_red_run(t *testing.T) {
	t.Parallel()

	code := codeRuleSeverities()

	// End-to-end: a README whose rule ID has a typo must produce drift.
	mutated := strings.Replace(
		readmeContent(t),
		"`missing-contact`",
		"`missing-contact-typo`",
		1,
	)
	require.NotEqual(t, readmeContent(t), mutated, "the mutation must change the README")

	mutatedPath := filepath.Join(t.TempDir(), "README.md")
	require.NoError(t, os.WriteFile(mutatedPath, []byte(mutated), 0o600))

	drift := ruleTableDrift(code, readmeDocumentedSeverities(t, mutatedPath))
	assert.Contains(t, drift, "missing-contact",
		"a renamed rule id must be reported as drift")

	// Comparator: a re-severitied rule must be reported even when every ID
	// still matches.
	flipped := make(map[finding.RuleName]finding.Severity, len(documentedFor(code)))
	maps.Copy(flipped, code)
	flipped["missing-response-time"] = finding.SeverityError

	drift = ruleTableDrift(code, flipped)
	assert.Contains(t, drift, "missing-response-time",
		"a severity flip must be reported as drift")
}

// readmeContent loads the real README for mutation in the red run.
func readmeContent(t *testing.T) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	require.NoError(t, err, "README.md must stay readable from the drift guard")

	return string(content)
}

func documentedFor(code map[finding.RuleName]finding.Severity) map[finding.RuleName]finding.Severity {
	return code
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
