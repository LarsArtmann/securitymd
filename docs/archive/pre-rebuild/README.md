# Pre-rebuild archive (template-SECURITY era)

These documents describe the **pre-2026-10 architecture** (`internal/` packages,
`cmd/template-security/`, `templates/`, `.template-security.yaml`, viper config,
policy-type system) — all of which was deleted in the securitymd rebuild
(`1c55390`, 2026-10-08).

~~Five root reports, moved here 2026-10-08.~~ Extended 2026-10-09 by the
docs-health pass: every pre-rebuild report in `docs/status/`, `docs/planning/`,
and `docs/modularization/` was annotated (inline ~~strikethrough~~ verdicts, no
open items remain) and archived here. `docs/status/` now holds only post-rebuild
reports. Superseded by
[`../../status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md`](../../status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md).

## Manifest

| File | Classification | Reason |
| ---- | -------------- | ------ |
| IMPROVEMENT_PLAN.md | ARCHIVE | Improvement plan against the deleted codebase; items either shipped in the rebuild or are moot |
| PARTS.md | ARCHIVE | Decomposition inventory of the old `internal/` layout; extraction never happened |
| PROJECT_SPLIT_EXECUTIVE_REPORT.md | ARCHIVE | Considered splitting a project that no longer exists in that shape; never executed |
| BDD_TESTS_REVIEW.md | ARCHIVE | Review of the old acceptance test files; its recommendation was adopted and survives in `test/acceptance/` |
| PUBLIC_OR_PRIVATE.md | ARCHIVE (decision doc) | Verdict "conditionally make public"; the condition (cleanup) was met by the rebuild — still the input for the pending publish decision (TODO_LIST), read before publishing |
| 2025-12-11_20-27_SHELL-TO-GO-MIGRATION-COMPLETE-STATUS.md | ARCHIVE | Celebrates the v2 template-CLI-SDK era, purged hours later by the simplified MVP |
| 2025-12-11_21-52_TEMPLATE-SECURITY-80-20-IMPLEMENTATION-STATUS.md | ARCHIVE | 80/20 phase tracking for files purged the same day |
| 2025-12-11_23-04_PROJECT-REFACTORING-STATUS.md | ARCHIVE | Refactoring status of the v2-era build errors, resolved by the simplification |
| 2025-12-11_23-05_ARCHITECTURE-REEVALUATION.md | ARCHIVE (decision doc) | Verdict "PURGE & SIMPLIFY" — executed same day via the simplified MVP |
| 2025-12-11_23-41_SIMPLIFIED-MVP-STATUS.md | ARCHIVE | MVP scope decision that survived into the rebuild (SECURITY.md only) |
| 2025-12-11_20-50_SECURITY-FOCUSED-80-20-PLAN.md | ARCHIVE | December 80/20 plan; Phase 1–2 ran, spirit kept by the rebuild |
| 2025-12-11_20-50_SECURITY-EXECUTION-GRAPH.md | ARCHIVE | Mermaid execution graph of the same plan |
| 2026-05-04_21-47_COMPREHENSIVE-PROJECT-STATUS.md | ARCHIVE | First go-finding-era comprehensive status; all items closed by later sessions or the rebuild |
| 2026-05-05_18-10_GO-FINDING-MIGRATION-COMPLETE.md | ARCHIVE | go-finding migration session 1; follow-ups closed same day |
| 2026-05-05_18-51_DEEP-GO-FINDING-INTEGRATION.md | ARCHIVE | go-finding deep-integration session 2; pipeline ideas superseded by toolsdk |
| 2026-05-05_23-27_CODE-DEDUPLICATION-AND-LINT-CLEANUP.md | ARCHIVE | Session 3 dedup/lint pass; split-brain question resolved by deletion |
| 2026-05-06_09-03_ARCHITECTURE-REVIEW-AND-IMPROVEMENTS.md | ARCHIVE | Architecture review whose core target (public `pkg/`, dead types gone) the rebuild adopted |
| 2026-05-06_09-03_FULL-CODE-REVIEW-AND-EXECUTION-PLAN.md | ARCHIVE | File-by-file review of deleted files; Pareto items all resolved |
| 2026-07-26_09-02_DEDUPLICATION-ZERO-CLONES-VERIFIED.md | ARCHIVE | art-dupl zero-clone verification of `internal/project_detector.go`, since deleted |
| go-composable-business-types-usage.md | ARCHIVE | Integration plan for a library that was renamed, adopted, removed within two days; types deleted |
| DEPENDENCY_GRAPH.md | ARCHIVE | Dependency analysis of the deleted `internal/` tree |
| EXECUTION_PLAN.md | ARCHIVE | 14-task modularization plan, never executed; superseded by the rebuild |
| PROPOSAL.md | ARCHIVE | The modularization proposal itself; solved differently by the rebuild |

Annotation convention: `~~item~~ done at <hash>` / `NOT-DO — <reason>` /
`moot — <reason>` / `routed — <where it lives now>`. Reference hashes:
`1c55390` (old tree deleted), `72085c2`/`4a8987a`/`519916d` (2026-10-08/09 docs
pass), `9d3094a` (flake.nix created 2026-06-17), `8864cb7` (CONTRIBUTING.md).
