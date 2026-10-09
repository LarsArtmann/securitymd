package policy

import (
	"fmt"
	"slices"
	"strings"

	"github.com/larsartmann/go-finding"
)

// contentRuleIDs mirrors the content-quality rules built in validate.go.
var contentRuleIDs = []finding.RuleName{ //nolint:gochecknoglobals // declarative rule-ID table, read-only after init
	ruleTooShort,
	ruleUnresolvedTemplate,
	ruleNoContent,
	ruleNoVersionInfo,
}

// KnownRuleIDs lists every rule ID securitymd can emit, including
// missing-file. Severity configuration and suppressions key on these IDs.
func KnownRuleIDs() []finding.RuleName {
	ids := make([]finding.RuleName, 0, len(sectionRules)+len(contentRuleIDs)+1)
	for _, rule := range sectionRules {
		ids = append(ids, rule.rule)
	}

	ids = append(ids, contentRuleIDs...)
	ids = append(ids, RuleMissingFile)

	return ids
}

// SeverityOverrides remaps finding severities per rule — the escape hatch for
// fleets that must adopt incrementally (e.g. downgrade missing-file from
// error to warning while policies roll out). Keys must be known rule IDs
// (KnownRuleIDs); values must be valid finding severities.
type SeverityOverrides map[finding.RuleName]finding.Severity

// ParseSeverityOverrides parses a "rule=severity[,rule=severity...]" spec
// into overrides. An empty spec means no overrides. Unknown rules, invalid
// severities, and malformed parts fail with a validation error that names
// the accepted values.
func ParseSeverityOverrides(spec string) (SeverityOverrides, error) {
	overrides := make(SeverityOverrides)

	for part := range strings.SplitSeq(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		override, err := parseSeverityOverride(part)
		if err != nil {
			return nil, err
		}

		overrides[override.rule] = override.severity
	}

	return overrides, nil
}

// severityOverride is one parsed "rule=severity" pair.
type severityOverride struct {
	rule     finding.RuleName
	severity finding.Severity
}

func parseSeverityOverride(part string) (severityOverride, error) {
	var parsed severityOverride

	ruleName, severityName, found := strings.Cut(part, "=")
	if !found {
		return parsed, finding.NewValidationError(
			fmt.Sprintf("invalid severity override %q: want rule=severity (e.g. missing-file=warning)", part), nil)
	}

	rule := finding.RuleName(strings.TrimSpace(ruleName))
	if !slices.Contains(KnownRuleIDs(), rule) {
		return parsed, finding.NewValidationError(
			fmt.Sprintf("unknown rule %q in severity override (known rules: %s)",
				rule, strings.Join(ruleIDNames(), ", ")), nil)
	}

	severity := finding.Severity(strings.TrimSpace(severityName))
	if !severity.IsValid() {
		return parsed, finding.NewValidationError(
			fmt.Sprintf("invalid severity %q for rule %q (valid severities: info, warning, error, critical)",
				severityName, rule), nil)
	}

	parsed.rule = rule
	parsed.severity = severity

	return parsed, nil
}

// ApplySeverityOverrides returns the findings with overridden severities.
// Unknown-rule keys are rejected at parse time, so a no-op here means the
// spec matched nothing — an empty override map returns the input unchanged.
func ApplySeverityOverrides(findings []finding.Finding, overrides SeverityOverrides) []finding.Finding {
	if len(overrides) == 0 {
		return findings
	}

	adjusted := slices.Clone(findings)
	for i, candidate := range adjusted {
		if severity, ok := overrides[candidate.Rule]; ok {
			adjusted[i].Severity = severity
		}
	}

	return adjusted
}

func ruleIDNames() []string {
	ids := KnownRuleIDs()

	names := make([]string, 0, len(ids))
	for _, id := range ids {
		names = append(names, string(id))
	}

	return names
}
