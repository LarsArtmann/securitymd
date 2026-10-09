# TODO List — securitymd

_Open items only. Refreshed 2026-10-09 evening after the publish pass: repo renamed + pushed + tagged `v1.0.0`, workflow re-enabled, dogfood SECURITY.md regenerated. Done this pass (no longer listed): suppression expiry grammar, exit-2 verdict, setup CI-crash fix, provider gate contract, drift-guard red run, naming audit (kept both flags, documented)._

## Blocked on Lars (needs human action)

- [ ] **P0: Fix GitHub Actions billing** — the first real runner run (2026-10-09, run 37957194805) was refused by GitHub: "recent account payments have failed or your spending limit needs to be increased" (all three jobs, zero steps executed). Nothing in the workflow is wrong. After fixing billing: `gh workflow run "Security Policy Validation" -R LarsArtmann/securitymd` (workflow_dispatch is wired) and watch all four hardened steps.
- [ ] **P0: Visibility flip** — the repo is still PRIVATE; the tag `v1.0.0` is cut, so making it public immediately makes `go install github.com/LarsArtmann/securitymd/cmd/securitymd@v1.0.0` and pkg.go.dev resolve. Verdict input: `docs/archive/pre-rebuild/PUBLIC_OR_PRIVATE.md` (conditional public; cleanup done by the rebuild). After flipping: verify pkg.go.dev renders + drop the "Published module" PARTIALLY note in FEATURES.
- [ ] **BuildFlow consumer switch** — drop the local replace (BuildFlow go.mod:425) + flake input, `go get github.com/LarsArtmann/securitymd@v1.0.0` (needs GOPRIVATE while private), `go work vendor`, `nix run .#update-vendor-hash`, `nix build .` — runbook section 2 in `docs/planning/2026-10-09_16-34_publish-runbook.md`. The local directory rename (runbook step 1) stays cosmetic/optional until this lands.
- [ ] **Fleet contact policy decision** — keep GitHub-advisory-only default (current) or bake a real address (e.g. `security@lars.software`) via BuildFlow `tool_options`. Default stays advisory-only until Lars rules (no fabricated addresses).
- [ ] **Fleet gate announcement + sweep** — draft and script are READY (`docs/planning/2026-10-09_fleet-gate-announcement-draft.md`, `scripts/fleet-securitymd-sweep.sh`); remaining: fill the two TODO placeholders (date + channel), run the `buildflow -s securitymd --fix` sweep, then flip the CI default (suppression comments, `until` expiry, and `severity-overrides` are the per-repo escape hatches).
- [ ] **Update `.config/metadata.yaml` tags** — still says `template`/`archived` (pre-rebuild relic). Decide: fresh tags reflecting the live securitymd tool, or confirm archiving. Consumer of the file is external tooling — Lars's call.
- [ ] **docs/reviews/ convention decision** (plan M24) — whether review reports get a standing directory convention; the repo-side `docs-gate` tooling already exists (`nix run .#docs-gate`).

## Post-publish verification (as soon as billing is fixed)

- [ ] First real CI runner run of the hardened workflow (coverage floor, fuzz smoke, govulncheck, docs-gate) — first attempt 2026-10-09 was billing-refused before any step executed; expect action-SHA or runner-image surprises on the real run.
- [ ] After the visibility flip: pkg.go.dev render check; GoReleaser or flake-based release workflow (ROADMAP).

## Process & ecosystem

- [ ] docs-health fleet cron: run `nix run .#docs-gate` per fleet repo on a schedule (the gate script exists in this repo's flake; wiring it fleet-wide needs a cron home — Lars).

## BuildFlow-side

- [ ] Once the concurrent `execution/` file-split sweep in BuildFlow lands (2026-10-09, still running): verify `nix build .` green and re-run `buildflow -s securitymd --fix` e2e with the nix-built binary.
- [ ] BuildFlow gotcha-anchor warnings (GOTCHAS.md:124/207 → `execution/pipeline.go:432/530`) need re-anchoring to the post-split files — belongs to the split sweep's follow-up.
- [ ] **Harden BuildFlow's findings gate against suppression** (found + verified at source 2026-10-09): `execution/workflow_result_2.go` `filterFindingsAtOrAbove` filters on severity only and ignores `Suppression` — securitymd's provider now strips suppressed findings at the boundary (pinned by `TestProvider_suppressed_findings_never_reach_the_gate`), but BuildFlow-side defense in depth would protect every tool.
- [ ] Guard test: assert `securitymd` tool_options validation error message text (report #34).
- [ ] Run `docs --check` after any provider-count change and fix table drift (report #32; clean as of 2026-10-09).

## Decided (recorded so nobody re-opens them)

- **Exit-2 propagation: CLI-only.** `buildflow -s securitymd` treats detect as advisory (exit 0 with findings) by design, and BuildFlow's own error taxonomy (errorfamily/ExitFindingsRemain=69) is that repo's contract. securitymd's 0/1/2 stays a shell/CI contract; revisit only if a real consumer needs the distinction surfaced through toolsdk.
- **Flag naming: keep `--severity` and `--set-severity`.** `--severity` is load-bearing (BuildFlow's provider passes `--severity=warning`); the two flags do different jobs and README documents the distinction.
- **Suppression grammar is final for v1: `securitymd:ignore(rule[,rule2]) [until YYYY-MM-DD] reason`.** Malformed dates inert; expiry grants the whole UTC day; expired directives attach nothing anywhere.
