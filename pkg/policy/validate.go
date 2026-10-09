package policy

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/larsartmann/go-finding"
)

// sectionRule declares one required SECURITY.md element. Patterns match
// case-insensitively as substrings; any single pattern satisfies the rule.
// Rule IDs are stable forever: suppressions and configs key on them.
type sectionRule struct {
	rule     finding.RuleName
	message  string
	severity finding.Severity
	patterns []string
}

var sectionRules = []sectionRule{ //nolint:gochecknoglobals // declarative rule table keyed by stable IDs; read-only after init
	{
		rule:     "missing-header",
		message:  "Missing '# Security Policy' header",
		severity: finding.SeverityError,
		patterns: []string{"security policy"},
	},
	{
		rule:     "missing-reporting",
		message:  "Missing 'Reporting a Vulnerability' section",
		severity: finding.SeverityError,
		patterns: []string{
			"reporting a vulnerability",
			"report a vulnerability",
			"reporting vulnerabilities",
			"vulnerability report",
			"disclos",
		},
	},
	{
		rule:     "missing-versions",
		message:  "Missing 'Supported Versions' section",
		severity: finding.SeverityError,
		patterns: []string{"supported version"},
	},
	{
		rule:     "missing-practices",
		message:  "Missing 'Security Practices' section",
		severity: finding.SeverityError,
		patterns: []string{"security practice"},
	},
	{
		rule:     "missing-contact",
		message:  "No contact channel: add a security email, GitHub advisory link, or security.txt reference",
		severity: finding.SeverityError,
		patterns: []string{"@", "security/advisories", "security.txt"},
	},
	{
		rule:     "missing-response-time",
		message:  "Should state a response-time commitment for vulnerability reports",
		severity: finding.SeverityWarning,
		patterns: []string{
			"response time",
			"respond within",
			"response within",
			"initial response",
			"within 24 hours",
			"within 48 hours",
			"within 72 hours",
		},
	},
}

const (
	minLines      = 20
	minLineLength = 20
)

// Content-quality rule IDs. Like the section-rule IDs above, these are stable
// forever: suppressions and severity configs key on them.
const (
	ruleTooShort           = finding.RuleName("too-short")
	ruleUnresolvedTemplate = finding.RuleName("unresolved-template")
	ruleNoContent          = finding.RuleName("no-content")
	ruleNoVersionInfo      = finding.RuleName("no-version-info")
)

// contentRuleSeverities is the declarative severity table for the
// content-quality rules built in validateContentQuality — the same single
// source the code and the README drift guard read, so a severity change here
// forces the documented table to follow.
var contentRuleSeverities = map[finding.RuleName]finding.Severity{ //nolint:gochecknoglobals // declarative severity table keyed by stable IDs; read-only after init
	ruleTooShort:           finding.SeverityWarning,
	ruleUnresolvedTemplate: finding.SeverityError,
	ruleNoContent:          finding.SeverityError,
	ruleNoVersionInfo:      finding.SeverityWarning,
}

// Validate checks a SECURITY.md file's content against the required-section
// and content-quality rules, returning one finding per violated rule. The
// file must exist; missing files are Detect's concern (rule "missing-file").
func Validate(filePath string) ([]finding.Finding, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")

	findings := make([]finding.Finding, 0, len(sectionRules))

	findings = append(findings, validateSections(filePath, lines)...)
	findings = append(findings, validateContentQuality(filePath, lines)...)

	return applySuppressions(lines, findings, time.Now()), nil
}

func validateSections(filePath string, lines []string) []finding.Finding {
	findings := make([]finding.Finding, 0, len(sectionRules))

	for _, rule := range sectionRules {
		if containsAnyFold(lines, rule.patterns) {
			continue
		}

		f, err := buildFinding(rule.rule, rule.message, rule.severity, filePath, 0)
		if err != nil {
			continue
		}

		findings = append(findings, f)
	}

	return findings
}

func validateContentQuality(filePath string, lines []string) []finding.Finding {
	var findings []finding.Finding

	if len(lines) < minLines {
		f, err := buildFinding(ruleTooShort,
			fmt.Sprintf("SECURITY.md seems too short (< %d lines)", minLines),
			contentRuleSeverities[ruleTooShort], filePath, 0)
		if err == nil {
			findings = append(findings, f)
		}
	}

	for lineIndex, line := range lines {
		if !strings.Contains(line, "{{") || !strings.Contains(line, "}}") {
			continue
		}

		f, err := buildFinding(ruleUnresolvedTemplate,
			"Unresolved template variable: "+strings.TrimSpace(line),
			contentRuleSeverities[ruleUnresolvedTemplate], filePath, lineIndex+1)
		if err == nil {
			findings = append(findings, f)
		}
	}

	if !hasSubstantiveContent(lines, minLineLength) {
		f, err := buildFinding(ruleNoContent,
			"SECURITY.md lacks substantive content",
			contentRuleSeverities[ruleNoContent], filePath, 0)
		if err == nil {
			findings = append(findings, f)
		}
	}

	if !hasVersionInformation(lines) {
		f, err := buildFinding(ruleNoVersionInfo,
			"No version information found",
			contentRuleSeverities[ruleNoVersionInfo], filePath, 0)
		if err == nil {
			findings = append(findings, f)
		}
	}

	return findings
}

// buildFinding creates a finding with the tool's consistent defaults:
// security category, suggest strategy. Line 0 becomes a file-level position —
// never a fabricated line number.
func buildFinding(
	rule finding.RuleName,
	message string,
	severity finding.Severity,
	filePath string,
	line int,
) (finding.Finding, error) {
	pos := finding.FilePos(finding.FilePath(filePath))
	if line > 0 {
		pos = finding.Pos(finding.FilePath(filePath), line, 1)
	}

	built, err := finding.NewBuilder(rule, ToolName, message, severity, pos).
		WithCategory(finding.CategorySecurity).
		WithTags(finding.TagSecurity).
		WithFixStrategy(finding.FixStrategySuggest).
		Build()
	if err != nil {
		return built, fmt.Errorf("build finding {rule:%q file:%q line:%d}: %w", rule, filePath, line, err)
	}

	return built, nil
}

// containsAnyFold reports whether any line contains any pattern,
// case-insensitively.
func containsAnyFold(lines []string, patterns []string) bool {
	for _, line := range lines {
		folded := strings.ToLower(line)
		for _, pattern := range patterns {
			if strings.Contains(folded, pattern) {
				return true
			}
		}
	}

	return false
}

// hasSubstantiveContent reports whether any line carries real content:
// long enough, not a heading, not placeholder/example text.
func hasSubstantiveContent(lines []string, minLength int) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) <= minLength ||
			strings.HasPrefix(trimmed, "#") ||
			strings.Contains(strings.ToLower(trimmed), "example") {
			continue
		}

		return true
	}

	return false
}

func hasVersionInformation(lines []string) bool {
	for _, line := range lines {
		folded := strings.ToLower(line)
		if strings.Contains(folded, "version") &&
			(strings.Contains(folded, "v") || strings.Contains(folded, ".")) {
			return true
		}
	}

	return false
}
