package policy

import (
	"slices"
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
// (a typo must fail safe, not suppress silently-by-accident). missing-file
// cannot be suppressed: with no file there is no content to carry the
// comment — a repository must have a policy.
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

	marked := slices.Clone(findings)
	for i, candidate := range marked {
		if reason, ok := reasons[string(candidate.Rule)]; ok {
			marked[i].Suppression = &finding.Suppression{
				Kind:      finding.SuppressionInSource,
				Rule:      candidate.Rule,
				Reason:    reason,
				ExpiresAt: nil,
			}
		}
	}

	return marked
}

// parseSuppressions extracts rule → reason from well-formed suppression
// comments; later comments for the same rule win.
func parseSuppressions(lines []string) map[string]string {
	reasons := make(map[string]string)

	for _, line := range lines {
		_, afterMarker, found := strings.Cut(line, suppressionMarker)
		if !found {
			continue
		}

		ruleList, reason, ok := parseSuppressionComment(afterMarker)
		if !ok {
			continue
		}

		for rule := range strings.SplitSeq(ruleList, ",") {
			rule = strings.TrimSpace(rule)
			if rule != "" {
				reasons[rule] = reason
			}
		}
	}

	return reasons
}

// parseSuppressionComment parses the comma-separated rule list before the
// closing parenthesis, then the required reason. The third result marks the
// comment malformed (missing paren, empty rule list, or empty reason); a
// malformed comment is inert.
func parseSuppressionComment(rest string) (string, string, bool) {
	ruleList, reason, found := strings.Cut(rest, ")")
	if !found {
		return "", "", false
	}

	reason = trimCommentTerminator(reason)
	if strings.TrimSpace(ruleList) == "" || reason == "" {
		return "", "", false
	}

	return ruleList, reason, true
}

// trimCommentTerminator strips the HTML comment's closing "-->" so that
// `<!-- securitymd:ignore(rule) -->` is a reason-less (inert) comment, not a
// suppression whose reason is "-->".
func trimCommentTerminator(reason string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(reason), "-->"))
}
