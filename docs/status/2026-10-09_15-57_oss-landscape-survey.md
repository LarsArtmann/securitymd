# OSS-Landscape Survey — Dedicated SECURITY.md Validators

**Date:** 2026-10-09 15:57 CEST
**Purpose:** Verify the positioning claim ("no dedicated SECURITY.md validator exists in OSS") before it ships on a README. Input to the publish checklist (plan 2026-10-09 M4).
**Method:** Web survey of npm, PyPI, GitHub search, and the repo-hygiene ecosystem (agentic research pass, 2026-10-09).

## Verdict

**No dedicated SECURITY.md validator/linter CLI exists** — the gap claim survives. One caveat: "first tool to validate SECURITY.md content" would be refutable, because OpenSSF Scorecard has scored SECURITY.md content heuristically since ~2020. Safe wording: **"the first dedicated CLI tool for linting and validating SECURITY.md files"** — defensible as of this survey date.

## Findings

### Dedicated SECURITY.md linters/validators: none found

- **npm:** nothing matching "security.md lint/validate" (only unrelated tools: `lockfile-lint`, `npm-package-json-lint`).
- **GitHub:** `SECURITY.md linter` → 0 results; `securitymd` → only test/demo fixture repos (e.g. `ThomasEvers64/securitymd-demo`, `CKSecOrg/test-repo-with-securitymd-*`), no shipped tool.
- **PyPI:** no `securitymd`/`secmd` equivalent.

### Closest neighbors (repo-hygiene tools that touch SECURITY.md)

| Tool                                                                | What it does re: SECURITY.md                                                                                                                                                                                                         | Validates content?                                                                                        |
| ------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------- |
| [OpenSSF Scorecard](https://github.com/ossf/scorecard)              | `Security-Policy` check + probes (`securityPolicyPresent`, `securityPolicyContainsLinks`, `securityPolicyContainsText`, `securityPolicyContainsVulnerabilityDisclosure`); awards points for links, disclosure text, timeline wording | Partial content heuristics (regex scoring), not section completeness; scoring framework, not a linter CLI |
| [Repolinter](https://github.com/todogroup/repolinter)               | Default ruleset: `security-file-exists`                                                                                                                                                                                              | Presence only                                                                                             |
| [hashload/boss](https://github.com/hashload/boss)                   | Checks SECURITY.md presence in root/`.github/`/`docs/`, offers `boss cra init` to generate a template                                                                                                                                | Presence + template generator (niche, Delphi ecosystem)                                                   |
| [MegaLinter/Super-Linter](https://github.com/oxsecurity/megalinter) | Generic markdown linting                                                                                                                                                                                                             | No SECURITY.md-specific rules                                                                             |
| [OSSF Security Insights](https://github.com/ossf/security-insights) | Schema-validates `security-insights.yml`, which references the policy                                                                                                                                                                | Different file entirely                                                                                   |

### GitHub's own tooling

Generation guidance and the "Add security policy" settings UI only — no validator is shipped; the repo sidebar detects presence.

## Consequences

- README/ROADMAP may claim "first dedicated SECURITY.md validator CLI" (with survey date), not "first to validate security policy content" (Scorecard contradicts the latter).
- Differentiation line for the website/README: securitymd checks **section completeness with line-precise, machine-readable findings** (go-finding/SARIF) and **generates a compliant policy** — Scorecard scores heuristically inside a much larger checksuite, Repolinter checks presence only.
- Re-survey before any "still first" claim in marketing copy older than ~6 months.
