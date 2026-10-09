package policy

import (
	"slices"
	"strings"
	"time"

	"github.com/larsartmann/go-finding"
)

// suppressionMarker begins an in-file suppression comment anywhere in a
// SECURITY.md:
//
//	securitymd:ignore(rule[,rule2]) reason
//	securitymd:ignore(rule[,rule2]) until YYYY-MM-DD reason
//
// The comment suppresses the listed rules for the entire policy file. The
// reason is REQUIRED: a marker without a reason is inert, so a stripped
// rationale never silently silences a rule. Unknown rule names are inert
// (a typo must fail safe, not suppress silently-by-accident). missing-file
// cannot be suppressed: with no file there is no content to carry the
// comment — a repository must have a policy.
//
// An `until` clause scopes the suppression to a calendar day: the rules stay
// suppressed through the end of the given UTC day and surface again the next
// day (deferred debt comes back on its own). A malformed date is inert — a
// suppression must be unambiguous to hide anything.
const suppressionMarker = "securitymd:ignore("

// untilPrefix introduces the optional expiry clause of a suppression reason.
const untilPrefix = "until "

// suppressionDirective is one parsed in-file suppression: why the rules are
// silenced, and until when (nil = indefinite).
type suppressionDirective struct {
	reason    string
	expiresAt *time.Time
}

// activeAt reports whether the directive still silences rules at now. An
// expired directive is treated as if it were never written: the finding
// carries no suppression metadata anywhere (JSON, SARIF, gate), so every
// consumer agrees the debt is due again.
func (d suppressionDirective) activeAt(now time.Time) bool {
	return d.expiresAt == nil || now.Before(*d.expiresAt)
}

// applySuppressions marks findings whose rule is suppressed by an in-file
// suppression comment that is still active at now. Suppressed findings keep
// their severity but carry the suppression metadata, so consumers (SARIF
// export, active-finding gates) exclude them while the evidence stays
// visible. Expired directives attach nothing: the finding surfaces as if
// unsuppressed.
func applySuppressions(lines []string, findings []finding.Finding, now time.Time) []finding.Finding {
	directives := parseSuppressions(lines)
	if len(directives) == 0 {
		return findings
	}

	marked := slices.Clone(findings)
	for i, candidate := range marked {
		directive, ok := directives[string(candidate.Rule)]
		if !ok || !directive.activeAt(now) {
			continue
		}

		marked[i].Suppression = &finding.Suppression{
			Kind:      finding.SuppressionInSource,
			Rule:      candidate.Rule,
			Reason:    directive.reason,
			ExpiresAt: directive.expiresAt,
		}
	}

	return marked
}

// parseSuppressions extracts rule → directive from well-formed suppression
// comments; later comments for the same rule win.
func parseSuppressions(lines []string) map[string]suppressionDirective {
	directives := make(map[string]suppressionDirective)

	for _, line := range lines {
		_, afterMarker, found := strings.Cut(line, suppressionMarker)
		if !found {
			continue
		}

		ruleList, directive, ok := parseSuppressionComment(afterMarker)
		if !ok {
			continue
		}

		for rule := range strings.SplitSeq(ruleList, ",") {
			rule = strings.TrimSpace(rule)
			if rule != "" {
				directives[rule] = directive
			}
		}
	}

	return directives
}

// parseSuppressionComment parses the comma-separated rule list before the
// closing parenthesis, then the required reason with its optional `until`
// clause. The final result marks the comment malformed (missing paren, empty
// rule list, empty reason, or an unparsable expiry date); a malformed comment
// is inert.
func parseSuppressionComment(rest string) (string, suppressionDirective, bool) {
	ruleList, reason, found := strings.Cut(rest, ")")
	if !found {
		return "", suppressionDirective{}, false
	}

	reason = trimCommentTerminator(reason)
	if strings.TrimSpace(ruleList) == "" || reason == "" {
		return "", suppressionDirective{}, false
	}

	directive, ok := parseSuppressionReason(reason)
	if !ok {
		return "", suppressionDirective{}, false
	}

	return ruleList, directive, true
}

// parseSuppressionReason splits an optional `until YYYY-MM-DD` prefix off the
// required reason. The expiry grants the whole given UTC day: ExpiresAt is
// the NEXT midnight, so "until 2026-12-31" stays active through that entire
// day and expires the moment 2027-01-01 begins. A malformed date fails the
// whole comment (inert), never degrades to an indefinite suppression.
func parseSuppressionReason(reason string) (suppressionDirective, bool) {
	if !strings.HasPrefix(reason, untilPrefix) {
		return suppressionDirective{reason: reason}, true
	}

	dateAndReason := strings.TrimPrefix(reason, untilPrefix)

	dateToken, remainingReason, _ := strings.Cut(dateAndReason, " ")

	expiryDay, err := time.Parse(time.DateOnly, dateToken)
	if err != nil {
		return suppressionDirective{}, false
	}

	remainingReason = strings.TrimSpace(remainingReason)
	if remainingReason == "" {
		return suppressionDirective{}, false
	}

	expiresAt := expiryDay.AddDate(0, 0, 1)

	return suppressionDirective{reason: remainingReason, expiresAt: &expiresAt}, true
}

// trimCommentTerminator strips the HTML comment's closing "-->" so that
// `<!-- securitymd:ignore(rule) -->` is a reason-less (inert) comment, not a
// suppression whose reason is "-->".
func trimCommentTerminator(reason string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(reason), "-->"))
}
