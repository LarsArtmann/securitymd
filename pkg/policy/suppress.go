package policy

import (
	"strings"

	"github.com/larsartmann/go-finding"
)

// suppressionMarker begins an in-file suppression comment anywhere in a
// SECURITY.md:
//
//	securitymd:ignore(rule[,rule2]) reason
//
// The comment suppresses the listed rules for the entire policy file. The
// reason is REQUIRED: a marker without a reason is inert, so a stripped
// rationale never silently silences a rule. Unknown rule names are inert
// (a typo must fail safe, not suppress nothing silently-by-accident).
// missing-file cannot be suppressed: with no file there is no content to
// carry the comment — a repository must have a policy.
const suppressionMarker = "securitymd:ignore("

// applySuppressions marks findings whose rule is suppressed by an in-file
// suppression comment. Suppressed findings keep their severity but carry the
// suppression metadata, so consumers (SARIF export, active-finding gates)
// exclude them while the evidence stays visible.
func applySuppressions(lines []string, findings []finding.Finding) []finding.Finding {
	reasons := parseSuppressions(lines)
	if len(reasons) == 0 {
		return findings
	}

	marked := make([]finding.Finding, len(findings))
	for i, f := range findings {
		if reason, ok := reasons[string(f.Rule)]; ok {
			f.Suppression = &finding.Suppression{
				Kind:   finding.SuppressionInSource,
				Rule:   f.Rule,
				Reason: reason,
			}
		}

		marked[i] = f
	}

	return marked
}

// parseSuppressions extracts rule → reason from well-formed suppression
// comments; later comments for the same rule win.
func parseSuppressions(lines []string) map[string]string {
	reasons := make(map[string]string)

	for _, line := range lines {
		markerIndex := strings.Index(line, suppressionMarker)
		if markerIndex < 0 {
			continue
		}

		rules, reason, ok := parseSuppressionComment(line[markerIndex+len(suppressionMarker):])
		if !ok {
			continue
		}

		for _, rule := range rules {
			reasons[rule] = reason
		}
	}

	return reasons
}

// parseSuppressionComment parses the comma-separated rule list up to the
// closing parenthesis, then the required reason. ok=false marks the comment
// malformed (missing paren, empty rule list, or empty reason) and inert.
func parseSuppressionComment(rest string) (rules []string, reason string, ok bool) {
	closing := strings.Index(rest, ")")
	if closing < 0 {
		return nil, "", false
	}

	ruleList := strings.TrimSpace(rest[:closing])
	reason = strings.TrimSpace(rest[closing+1:])
	if ruleList == "" || reason == "" {
		return nil, "", false
	}

	rules = make([]string, 0, 4)
	for _, rule := range strings.Split(ruleList, ",") {
		rule = strings.TrimSpace(rule)
		if rule != "" {
			rules = append(rules, rule)
		}
	}

	if len(rules) == 0 {
		return nil, "", false
	}

	return rules, reason, true
}
