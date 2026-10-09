# TODO List — securitymd

_Open items only. Refreshed 2026-10-09 late evening (docs-health AUDIT pass + release pass): local checkout renamed to `~/projects/securitymd`, BuildFlow `nix build` + E2E option channel verified green (recorded in AGENTS.md), post-rebuild status/planning reports annotated + archived under per-dir `archived/`, **v1.0.0 + v1.1.0 released** (GitHub Releases cut; proxy + sumdb + clean-dir `go get` verified; CI green post-billing-fix; visibility flipped). Done this pass (no longer listed): suppression expiry grammar, exit-2 verdict, setup CI-crash fix, provider gate contract, drift-guard red run, naming audit (kept both flags, documented), critical-gate parity check in BuildFlow (verified at source, no defect), GitHub Actions billing fix (runs 37985112024/37985651759/37986553713 green), visibility flip (repo is PUBLIC), first real CI runner runs of the hardened workflow._

## Blocked on Lars (needs human action)

- [ ] **BuildFlow consumer switch** — the visibility flip is done, so this is now actionable: drop BOTH local replaces (BuildFlow `go.mod:425` + `tools/go.mod:259`), `go get github.com/LarsArtmann/securitymd@v1.1.0`, then `nix flake update securitymd` + `nix run .#update-vendor-hash` + `nix build .`. NO vendoring (BuildFlow gotcha #229 retired `go work vendor`) and the flake input STAYS (tracks latest per the 2026-10-05 input policy). Runbook (historical, superseded in section 2): `docs/planning/archived/2026-10-09_16-34_publish-runbook.md`.
- [ ] **Fleet contact policy decision** — keep GitHub-advisory-only default (current) or bake a real address (e.g. `security@lars.software`) via BuildFlow `tool_options`. Default stays advisory-only until Lars rules (no fabricated addresses).
- [ ] **Fleet gate announcement + sweep** — draft and script are READY (`docs/planning/2026-10-09_fleet-gate-announcement-draft.md`, `scripts/fleet-securitymd-sweep.sh`); remaining: fill the two TODO placeholders (date + channel), add the `.bak`-gitignore guidance line + a pre-announcement failing-repo count to the announcement, run the `buildflow -s securitymd --fix` sweep, then flip the CI default (suppression comments, `until` expiry, and `severity-overrides` are the per-repo escape hatches).
- [ ] **Update `.config/metadata.yaml` tags** — still says `template`/`archived` (pre-rebuild relic). Decide: fresh tags reflecting the live securitymd tool, or confirm archiving. Consumer of the file is external tooling — Lars's call.
- [ ] **docs/reviews/ convention decision** (plan M24) — whether review reports get a standing directory convention; the repo-side `docs-gate` tooling already exists (`nix run .#docs-gate`).

## Post-release verification

- [ ] pkg.go.dev render check for v1.0.0 + v1.1.0 — proxy + sumdb + clean-dir `go get` all verified 2026-10-09 evening; the pkg.go.dev docs pages still 404 due to on-demand rendering lag (up to ~1h after indexing). Re-check `https://pkg.go.dev/github.com/LarsArtmann/securitymd@v1.1.0` later.
- [ ] GoReleaser or flake-based release workflow (ROADMAP) — releases are currently cut by hand (`gh release create` against pushed annotated tags).

## Process & ecosystem

- [ ] **Prove the extended docs-gate red** — mutate a throwaway copy (drop an un-annotated file into `docs/status/archived/`), watch the gate fail, record it (third-time lesson: negative-test-first; flagged in the 20:19 status report).
- [ ] docs-health fleet cron: run `nix run .#docs-gate` per fleet repo on a schedule (the gate script exists in this repo's flake; wiring it fleet-wide needs a cron home — Lars).

## BuildFlow-side

- [ ] **Rebuild the buildflow binary with the fixed go-auto-upgrade FindGoMod** — `FindGoMod` started its walk at `filepath.Dir(path)`, so a module-root WorkDir skipped its own go.mod and stdlib2lo's dependency gate always closed (`lo-dependency-missing` warnings even with samber/lo required; verified dir→false / file→true at pinned rev 64ca820). Fixed upstream in go-auto-upgrade `92b4dfd` (+ regression test). To activate: `nix flake update go-auto-upgrade` in BuildFlow, then `nix build . && nix run .#reinstall` (concurrent-session caution). This repo already applied the lo rewrites manually, so the expected post-rebuild delta is just the 5 stale `lo-dependency-missing` warnings disappearing.
- [ ] BuildFlow gotcha-anchor warnings (GOTCHAS.md:124/207 → `execution/pipeline.go:432/530`) need re-anchoring to the post-split files — belongs to the split sweep's follow-up.
- [ ] **Harden BuildFlow's findings gate against suppression** (found + verified at source 2026-10-09): `execution/workflow_result_2.go` `filterFindingsAtOrAbove` filters on severity only and ignores `Suppression` — securitymd's provider now strips suppressed findings at the boundary (pinned by `TestProvider_suppressed_findings_never_reach_the_gate`), but BuildFlow-side defense in depth would protect every tool.
- [ ] Guard test: assert `securitymd` tool_options validation error message text (report #34).
- [ ] Run `docs --check` after any provider-count change and fix table drift (report #32; clean as of 2026-10-09).
- [ ] **docs-gate manifest-coverage leg**: an archived file must have a row in its archive home's manifest README (annotation presence is gated today; classification completeness is not).

## Decided (recorded so nobody re-opens them)

- **Exit-2 propagation: CLI-only.** `buildflow -s securitymd` treats detect as advisory (exit 0 with findings) by design, and BuildFlow's own error taxonomy (errorfamily/ExitFindingsRemain=69) is that repo's contract. securitymd's 0/1/2 stays a shell/CI contract; revisit only if a real consumer needs the distinction surfaced through toolsdk.
- **Flag naming: keep `--severity` and `--set-severity`.** `--severity` is load-bearing (BuildFlow's provider passes `--severity=warning`); the two flags do different jobs and README documents the distinction.
- **Suppression grammar is final for v1: `securitymd:ignore(rule[,rule2]) [until YYYY-MM-DD] reason`.** Malformed dates inert; expiry grants the whole UTC day; expired directives attach nothing anywhere.
- **BuildFlow's findings gate handles escalated severities correctly — verified at source 2026-10-09, nothing to port.** `filterFindingsAtOrAbove` (`execution/workflow_result_2.go:187`) filters with `!Severity.LessThan(threshold)` (threshold semantics), so a `--set-severity …=critical` finding trips the gate. The gap was suppression metadata only — covered by the provider boundary strip above.
