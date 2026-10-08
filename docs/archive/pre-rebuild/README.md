# Pre-rebuild archive (template-SECURITY era)

These documents describe the **pre-2026-10 architecture** (`internal/` packages,
`cmd/template-security/`, `templates/`, `.template-security.yaml`, viper config,
policy-type system) — all of which was deleted in the securitymd rebuild.

They are archived, not maintained. Superseded by
[`../../status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md`](../../status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md).

| File | Classification | Reason |
| ---- | -------------- | ------ |
| IMPROVEMENT_PLAN.md | ARCHIVE | Improvement plan against the deleted codebase; items either shipped in the rebuild or are moot |
| PARTS.md | ARCHIVE | Decomposition inventory of the old `internal/` layout |
| PROJECT_SPLIT_EXECUTIVE_REPORT.md | ARCHIVE | Considered splitting a project that no longer exists in that shape |
| BDD_TESTS_REVIEW.md | ARCHIVE | Review of the old acceptance test files (replaced by `test/acceptance/`) |
| PUBLIC_OR_PRIVATE.md | ARCHIVE (decision doc) | Verdict was "conditionally make public" after cleanup; the rebuild did the cleanup. Still relevant input for the pending publish decision — read before publishing, then archive for good |
