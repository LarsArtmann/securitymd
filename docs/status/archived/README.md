# Post-rebuild status archive

Closed history from the 2026-10-09 rebuild→publish cycle. Every file here was
annotated by the 2026-10-09 docs-health AUDIT pass (inline `~~verdict~~` per
numbered item/table row: `done at <hash>` / `routed — <where>` / `decided` /
`NOT-DO` / `superseded`), then archived via `git mv` — history intact. Do not
open without a concrete need; forward-looking work lives in `TODO_LIST.md` and
`ROADMAP.md`, product history in `CHANGELOG.md`.

## Manifest (bulk archive, 2026-10-09 docs-health pass)

| File                                                              | Classification         | Deciding reason                                                                                                                                                                   |
| ----------------------------------------------------------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md  | ARCHIVE                | Rebuild session report; all items done (rebuild `1c55390`), superseded (flake-input FOD), or routed — 91 verdicts                                                                   |
| 2026-10-09_02-02_docs-pass-nix-flake-input-and-concurrent-refactor-triage.md | ARCHIVE    | Docs-pass + nix-unblock report; items done at `519916d`/`5c739ec` or routed — 92 verdicts                                                                                          |
| 2026-10-09_15-11_docs-health-annotation-archive-sweep.md          | ARCHIVE                | The previous docs-health pass (537 verdicts, pre-rebuild archive); its open items shipped (survey `f387d65`, gate `e195158`) or routed — 47 verdicts + 3 tilde-protected bare lines |
| 2026-10-09_15-57_oss-landscape-survey.md                          | ARCHIVE (reference)    | OSS-landscape survey backing ROADMAP's positioning claim; consequences shipped, re-survey cadence routed to ROADMAP — 3 verdicts                                                    |
| 2026-10-09_16-48_pareto-execution-and-self-review-status.md       | ARCHIVE                | Pareto execution sweep; all f-items done by later sessions or routed — 94 verdicts                                                                                                 |
| 2026-10-09_17-08_securitymd-buildflow-integration-session-status.md | ARCHIVE              | BuildFlow toolsdk integration; guard test verified passing in BuildFlow, e2e/nix verified 2026-10-09 — 85 verdicts                                                                  |
| 2026-10-09_17-24_post-hardening-execution-status.md               | ARCHIVE                | Hardening pass (exit contract, drift guard, CI); its f-list executed by the publish pass or routed — 80 verdicts                                                                    |

Still live: [`../2026-10-09_18-15_done-done-publish-pass-status.md`](../2026-10-09_18-15_done-done-publish-pass-status.md)
— the newest report; its blocked items (billing, visibility flip, BuildFlow
consumer switch) are the active TODO_LIST P0s, so it stays unarchived.

## Exceptions (documented, not silently skipped)

Three lines in the 15-11 report and one row (F10.1) in the archived plan
(`../../planning/archived/2026-10-09_15-37_pareto-publish-and-hardening-plan.md`)
contain literal `` `~~` `` inside backticked text; the annotator's
already-struck protection refuses those lines by design. They are
self-describing (each narrates its own resolution) and intentionally bare.

Convention note: this repo now uses per-directory `archived/` subdirectories
(skill default) for post-rebuild docs; `docs/archive/pre-rebuild/` remains the
flat pre-rebuild archive with its own manifest.
