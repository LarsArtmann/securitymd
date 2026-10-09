# Roadmap — securitymd

_Long-term direction and raw ideas. Actionable items live in TODO_LIST.md; graduated items move there when scoped._

## Direction

- Stay a single-purpose, single-binary tool in the go-finding family; BuildFlow is the primary consumer, the CLI the manual escape hatch
- Publish as open source — repo renamed + pushed + tagged `v1.0.0` 2026-10-09; repo is still PRIVATE, so the remaining button is the GitHub visibility flip (Lars), after which the module proxy/pkg.go.dev resolve. Fills a real gap: **verified 2026-10-09** that no dedicated SECURITY.md validator/linter CLI exists in OSS (survey: `docs/status/archived/2026-10-09_15-57_oss-landscape-survey.md`); closest neighbor is OpenSSF Scorecard's heuristic Security-Policy scoring, so claims say "first dedicated validator CLI", never "first to validate policy content"

## Raw ideas

### Promoted from the 2026-10-09 execution report (section f, long tail)

- `validate --repair`: detect + generate one-shot for CLI users
- ~~Suppression expiry: `securitymd:ignore(rule) until YYYY-MM-DD reason` via go-finding's time model — gated on Lars's yes/no (grammar change + golden refresh)~~ DONE 2026-10-09: shipped with full-day grant semantics, goldens, and contract scenarios
- Provider `location` toolsdk option so BuildFlow repos pick `.github`/`docs` canonically
- ~~`--no-color` / `NO_COLOR` confirmation and docs~~ DONE 2026-10-09: verified at source (fatih/color honors NO_COLOR/TERM=dumb/non-TTY) and pinned by an ANSI-free contract scenario; shell completions + man pages (cobra built-ins) still open
- Windows CI leg (`.bak` timestamp format is already colon-free; verify the rest)
- Benchmarks: `Detect` on a large repo — perf baseline before structure-aware validation
- `securitymd list-rules`: machine-readable rule metadata for CI config UIs
- `securitymd report --since`: diff two SARIF files, ratchet-lite without state
- Org-level SECURITY.md inheritance (central policy, per-repo overrides)
- SARIF `automationDetails`/run metadata embedding repo identity — go-finding upstream discussion first
- Version stamp: flake package should stamp the git rev, not the static `0.1.0-dev`
- Pre-commit hook recipe in README (validate on commit)
- ~~Naming audit: `--severity` (SARIF floor) vs `--set-severity` (override) confusion — rename or alias~~ RESOLVED 2026-10-09: keep both — BuildFlow's provider passes `--severity=warning` (the flag is load-bearing) and the two jobs differ; README now documents the distinction
- Verify before claiming: "config-free posture" (no config file) as an advertised feature
- Sync `.golangci.yml` build-tags/go version with go.mod on every toolchain bump (currently tolerated skew)
- Post-publish: nixpkgs PR packaging securitymd (the flake gives the derivation); Dependabot/renovate for fleet repos

### Promoted from the hardening-session next-lists (2026-10-09 docs-health harvest)

- Fuzz the severity-override parser (`parseSeverityOverrides`) — same corpus discipline as `FuzzParseGitRemote`
- Property fuzz over `GenerateOptions`: a generated policy never contains `{{`
- CLI-level golden for `--format json` stdout+exit combos beyond the single pinned scenario
- Per-package coverage floors (or a `cmd/` ratchet) — repo-total 87.6% hides `cmd/` decay
- golangci-lint into CI (needs a verified action SHA; decide the `cmd/` lint-exclusion question first)
- Sweep script: CRLF-tolerant repos file + `--json` summary mode
- One-off audit: gosec (inside golangci) vs govulncheck overlap gap
- Re-survey the OSS landscape before any "still first" marketing claim older than ~6 months (baseline survey 2026-10-09: `docs/status/archived/2026-10-09_15-57_oss-landscape-survey.md`)
- Contract test: replace fixture sentinel strings with a named resolver column
- Concurrent-session claim protocol (heartbeat/lock file) — process idea from three reports

### Older ideas

- Markdown-structure-aware validation phase (headings/table parsing instead of substring heuristics) — gated on suppressions-first (shipped 2026-10-09); parity-test any port against the substring detector
- GitHub Action wrapper (`securitymd-action`) for non-BuildFlow consumers — adopt after publish (needs an installed module)
- golangci-lint plugin distribution of the detector — **verdict 2026-10-09: won't implement.** golangci's plugin ABI targets Go-AST analysis; a markdown policy linter shoehorned into it is architecture mismatch, and BuildFlow + CLI + (future) Action already cover distribution
- Baseline/ratchet mode for incremental fleet adoption — **verdict 2026-10-09: defer.** Severity overrides (`--set-severity missing-file=warning`) already unblock incremental adoption; a baseline file adds stateful machinery with no fleet demand yet. Revisit at ~50 adopting repos
- Localized finding descriptions (English-only today) — **verdict 2026-10-09: defer.** Findings are machine-consumed (SARIF/CI logs); translation maintenance has zero current user signal
- Website/docs launch once published (website-launch pattern)
- Feedback upstream: exhaustruct ignore-patterns — **verdict 2026-10-09: no defect to file.** Verified at source (golangci-lint-auto-configure `linter_settings.go:244`): the curated patterns are already anchored type-form (`net/http.Client`, ...), and the v4→v5 key migration carries type-form values, so nothing path-based is emitted
- Proposal to linter-autoconfigure-sdk: byte-faithful `SaveBytes` — **verdict 2026-10-09: resolved upstream.** SDK v0.8.0's `SaveJSONBytes` is already byte-faithful ("no newline is appended or trimmed") with `MarshalJSONIndented` as the bytes helper; only the naming nit remains, not worth an issue
- ✅ Recorded the "nix FOD ignores local replaces" lesson in crush-config `references/lessons.md` (2026-10-09)
- ✅ Updated the nix-private-go-repos skill: nixpkgs `gotools` bundles an older Go than go.mod's floor, so treefmt's goimports fails in-sandbox with a toolchain-download attempt (verified 2026-10-09 in securitymd's flake.nix)

## Open questions

- Fleet-wide contact default (advisory-only vs real email) — Lars
- Does `securitymd` belong in crush-config's project-discovery checklist so future sessions check SECURITY.md compliance?
