# Roadmap — securitymd

_Long-term direction and raw ideas. Actionable items live in TODO_LIST.md; graduated items move there when scoped._

## Direction

- Stay a single-purpose, single-binary tool in the go-finding family; BuildFlow is the primary consumer, the CLI the manual escape hatch
- Publish as open source (see TODO_LIST publish checklist) — fills a real gap: **unverified belief** that no dedicated SECURITY.md validator exists in OSS; survey the landscape before publishing and drop the claim if contradicted

## Raw ideas

- Markdown-structure-aware validation phase (headings/table parsing instead of substring heuristics)
- GitHub Action wrapper (`securitymd-action`) for non-BuildFlow consumers
- golangci-lint plugin distribution of the detector (low priority — BuildFlow is the consumer)
- Baseline/ratchet mode for incremental fleet adoption
- Localized finding descriptions (English-only today)
- Website/docs launch once published (website-launch pattern)
- Feedback upstream: exhaustruct ignore-patterns must be anchored struct patterns (path-based ones are silent no-ops) — note to golangci-lint-auto-configure
- Proposal to linter-autoconfigure-sdk: first-class byte-faithful `SaveBytes` (or doc rename for `SaveJSONBytes`)
- Record the "nix FOD ignores local replaces" lesson in crush-config `references/lessons.md` (cross-project lesson)
- Update the nix-private-go-repos skill: nixpkgs `gotools` bundles an older Go than go.mod's floor, so treefmt's goimports fails in-sandbox with a toolchain-download attempt (found 2026-10-09 in securitymd's flake check)

## Open questions

- Fleet-wide contact default (advisory-only vs real email) — Lars
- Does `securitymd` belong in crush-config's project-discovery checklist so future sessions check SECURITY.md compliance?
