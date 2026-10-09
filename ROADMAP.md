# Roadmap — securitymd

_Long-term direction and raw ideas. Actionable items live in TODO_LIST.md; graduated items move there when scoped._

## Direction

- Stay a single-purpose, single-binary tool in the go-finding family; BuildFlow is the primary consumer, the CLI the manual escape hatch
- Publish as open source (see TODO_LIST publish checklist) — fills a real gap: **verified 2026-10-09** that no dedicated SECURITY.md validator/linter CLI exists in OSS (survey: `docs/status/2026-10-09_15-57_oss-landscape-survey.md`); closest neighbor is OpenSSF Scorecard's heuristic Security-Policy scoring, so claims say "first dedicated validator CLI", never "first to validate policy content"

## Raw ideas

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
