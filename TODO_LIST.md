# TODO List — securitymd

_Open items only. Source: harvest of `docs/status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md` (2026-10-09), minus items completed since. The 2026-10-09 Pareto plan's executable set (M2–M10, M14–M18-mechanical, M22, M23, M24-gate) is complete — see `docs/planning/2026-10-09_15-37_pareto-publish-and-hardening-plan.md` and `docs/planning/2026-10-09_16-34_publish-runbook.md`._

## Blocked on Lars (needs human action)

- [ ] **P0: Publish securitymd** — rename local dir + GitHub repo `template-SECURITY` → `securitymd`, push, tag `v1.0.0`; then drop BuildFlow replaces + flake input, re-vendor, `nix run .#update-vendor-hash`; regenerate this repo's own SECURITY.md (current one predates the rename and would render the old repo name). Every mechanical step is staged and pre-verified in `docs/planning/2026-10-09_16-34_publish-runbook.md`. Input: `docs/archive/pre-rebuild/PUBLIC_OR_PRIVATE.md` (verdict: conditionally make public — cleanup done by the rebuild).
- [ ] **Fleet contact policy decision** — keep GitHub-advisory-only default (current) or bake a real address (e.g. `security@lars.software`) via BuildFlow `tool_options`. Default stays advisory-only until Lars rules (no fabricated addresses).
- [ ] **Fleet gate announcement** — `missing-file` is error severity; every manifest-carrying repo without SECURITY.md fails BuildFlow's findings gate. Announce + run one `buildflow -s securitymd --fix` sweep before it lands in CI defaults (suppression comments and `severity-overrides` are the per-repo escape hatches).
- [ ] **Update `.config/metadata.yaml` tags** — still says `template`/`archived` (pre-rebuild relic; also flagged by PUBLIC_OR_PRIVATE's should-fix list). Decide: fresh tags reflecting the live securitymd tool, or confirm archiving. Consumer of the file is external tooling — Lars's call. (Examined 2026-10-09: `git-town.toml` is fine — `main = "master"`, GitHub API connector.)
- [ ] **docs/reviews/ convention decision** (plan M24) — whether review reports get a standing directory convention; the repo-side `docs-gate` tooling already exists (`nix run .#docs-gate`).

## Process & ecosystem

- [ ] docs-health fleet cron: run `nix run .#docs-gate` per fleet repo on a schedule (the gate script exists in this repo's flake; wiring it fleet-wide needs a cron home — Lars).

## BuildFlow-side

- [ ] Once the concurrent `execution/` file-split sweep in BuildFlow lands (2026-10-09, still running at pipeline_1300+): verify `nix build .` green and re-run `buildflow -s securitymd --fix` e2e with the nix-built binary. The securitymd FOD resolution itself is already proven (deps phase passes; vendorHash invariant, gotcha #230).
- [ ] BuildFlow gotcha-anchor warnings (GOTCHAS.md:124/207 → `execution/pipeline.go:432/530`) need re-anchoring to the post-split files — belongs to the split sweep's follow-up.
- [ ] Guard test: assert `securitymd` tool_options validation error message text (report #34)
- [ ] Run `docs --check` after any provider-count change and fix table drift (report #32; clean at 0 fail as of 2026-10-09)
