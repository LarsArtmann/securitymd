# TODO List — securitymd

_Open items only. Harvested 2026-10-09 from `docs/status/2026-10-09_16-48_pareto-execution-and-self-review-status.md` section (f) — its do-now items 1–10 were executed same-day (exit-code contract, `status --location`, suppressed goldens, fixture single source, README drift guard, severity-combination tests, CI hardening, fleet prep); long-tail items live in ROADMAP "Promoted from the execution report"._

## Blocked on Lars (needs human action)

- [ ] **P0: Publish securitymd** — rename local dir + GitHub repo `template-SECURITY` → `securitymd`, push, tag `v1.0.0`; then drop BuildFlow replaces + flake input, re-vendor, `nix run .#update-vendor-hash`; regenerate this repo's own SECURITY.md (current one predates the rename and would render the old repo name). Every mechanical step is staged and pre-verified in `docs/planning/2026-10-09_16-34_publish-runbook.md`. Input: `docs/archive/pre-rebuild/PUBLIC_OR_PRIVATE.md` (verdict: conditionally make public — cleanup done by the rebuild).
- [ ] **Fleet contact policy decision** — keep GitHub-advisory-only default (current) or bake a real address (e.g. `security@lars.software`) via BuildFlow `tool_options`. Default stays advisory-only until Lars rules (no fabricated addresses).
- [ ] **Fleet gate announcement + sweep** — draft and script are READY (`docs/planning/2026-10-09_fleet-gate-announcement-draft.md`, `scripts/fleet-securitymd-sweep.sh`); remaining: fill the two TODO placeholders (date + channel), run the `buildflow -s securitymd --fix` sweep, then flip the CI default (suppression comments and `severity-overrides` are the per-repo escape hatches).
- [ ] **Update `.config/metadata.yaml` tags** — still says `template`/`archived` (pre-rebuild relic; also flagged by PUBLIC_OR_PRIVATE's should-fix list). Decide: fresh tags reflecting the live securitymd tool, or confirm archiving. Consumer of the file is external tooling — Lars's call. (Examined 2026-10-09: `git-town.toml` is fine — `main = "master"`, GitHub API connector.)
- [ ] **docs/reviews/ convention decision** (plan M24) — whether review reports get a standing directory convention; the repo-side `docs-gate` tooling already exists (`nix run .#docs-gate`).
- [ ] **Exit-code follow-ups needing Lars's verdict** (from report section g): suppression expiry grammar (`until YYYY-MM-DD`) — yes/no; whether BuildFlow's provider should map securitymd's exit 2 distinctly (toolsdk error taxonomy fit, report f.39).

## Post-publish verification (first push triggers these)

- [ ] First real CI runner run of the hardened workflow (coverage floor, fuzz smoke, govulncheck, docs-gate) — all steps verified locally 2026-10-09, but never executed on a runner; expect action-SHA or runner-image surprises.
- [ ] Post-publish: pkg.go.dev render check + repo description/topics; GoReleaser or flake-based release workflow (ROADMAP).

## Process & ecosystem

- [ ] docs-health fleet cron: run `nix run .#docs-gate` per fleet repo on a schedule (the gate script exists in this repo's flake; wiring it fleet-wide needs a cron home — Lars).

## BuildFlow-side

- [ ] Once the concurrent `execution/` file-split sweep in BuildFlow lands (2026-10-09, still running at pipeline_1300+): verify `nix build .` green and re-run `buildflow -s securitymd --fix` e2e with the nix-built binary. The securitymd FOD resolution itself is already proven (deps phase passes; vendorHash invariant, gotcha #230).
- [ ] BuildFlow gotcha-anchor warnings (GOTCHAS.md:124/207 → `execution/pipeline.go:432/530`) need re-anchoring to the post-split files — belongs to the split sweep's follow-up.
- [ ] Guard test: assert `securitymd` tool_options validation error message text (report #34)
- [ ] Run `docs --check` after any provider-count change and fix table drift (report #32; clean at 0 fail as of 2026-10-09)
