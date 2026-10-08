# Status: securitymd Rebuild + BuildFlow Provider Integration

**Date:** 2026-10-08 23:40
**Session scope:** Rename `template-SECURITY` → `securitymd`, rebuild as a go-finding/toolsdk-based linter, wire into BuildFlow. Skills used: naming-review, linter-building (+ ecosystem/mechanism refs), buildflow. Report covers only this session.

---

## Verdict in one line

The tool is rebuilt, fully tested, lint-clean, and **works end-to-end inside BuildFlow** (detect → repair → verify verified live); the one real blocker is that BuildFlow's **nix FOD build cannot resolve the not-yet-published module** — the local `replace` works for `go build`/`go test` but is ignored by the nix prepared-source fetch.

---

## a) FULLY DONE (verified this session)

1. **Research/verification pass**: read toolsdk (spec/registry/triggers/options/dryrun), linter-autoconfigure-sdk (full source), reference providers (licenseforge, dependabot-auto-configure, erraudit), BuildFlow sdk_registry/sdk_imports/AGENTS.md. All API claims verified against source, not docs.
2. **Naming decision (10 candidates)** — winner `securitymd`:
   | #  | Name                           | Verdict                                                                                 |
   | -- | ------------------------------ | --------------------------------------------------------------------------------------- |
   | 1  | **securitymd**                 | ✅ chosen: names exactly what it owns, short namespace for `-s securitymd`/suppressions |
   | 2  | policyforge                    | licenseforge sibling, but "policy" is overbroad (OPA? CI policy?)                       |
   | 3  | security-policy-auto-configure | behavior-exact family name, but mechanism-named and long                                |
   | 4  | security-policy-sync           | "sync" lies — no source of truth to sync from                                           |
   | 5  | security-policy-lint           | honest for validate half, hides generate half                                           |
   | 6  | disclosure-policy              | the true domain term (vulnerability disclosure), obscure to newcomers                   |
   | 7  | vulnpolicy                     | short but coined/unfamiliar                                                             |
   | 8  | secguard                       | vague guard-noun, names nothing                                                         |
   | 9  | securitymd-check               | suffix adds nothing over securitymd                                                     |
   | 10 | secmd                          | unpronounceable abbreviation                                                            |
3. **Module renamed** to `github.com/LarsArtmann/securitymd`; deps: go-finding, go-finding/toolsdk, linter-autoconfigure-sdk (WorkingDir/FirstExisting), go-atomic-write, cobra, color, testify, ginkgo/gomega. **Dropped viper** (kills ~15 indirect deps) and go-finding/pipeline.
4. **`pkg/policy` core**: candidate discovery (`SECURITY.md`, `.github/`, `docs/`), 10 stable kebab rule IDs (new: `missing-file`), multi-signal section patterns (advisory link/security.txt count as contact; broader response-time phrasing), line-precise `unresolved-template`, **embedded template** (`go:embed` — fixes the old runtime dependency on a `templates/` dir relative to CWD), never-overwrite generation, dry-run, atomic idempotent write, git-remote identity parsing (https/ssh/gitlab-nested), latest-tag best effort. Missing-file finding carries the rendered policy as `AfterCode` + `FixStrategyDirect` when identity is derivable (go-finding's builder validation forced this honesty).
5. **`pkg/provider`**: toolsdk Spec (name `securitymd`, trigger on manifests + policy locations, `contact-email` tool option, Detect + Repair, dry-run aware, nil HealthCheck → NoOp via BuildFlow conversion).
6. **CLI `cmd/securitymd`**: `validate` (text/json/sarif, `--file`, `--severity`, exit 0/1 by error findings, SilenceUsage), `setup` (org/repo/email/directory/dry-run), `status`.
7. **Old code removed** (trash, not rm): `internal/` (incl. dead types package), `cmd/template-security/`, `templates/`, `scripts/` (~1,300 legacy shell lines), `.template-security.yaml`, `.golangci.yml.v1`, `.go-arch-lint.yml`, old acceptance tests.
8. **Tests all green** (`go test ./...`): unit (validate/generate/detect/project), provider contract tests (spec shape + option validation, full detect→repair→verify loop in a real git repo, dry-run, never-overwrite, contact-email option flow), 7 Ginkgo BDD specs. Includes the dogfood invariant: _the generated template must pass its own validator_.
9. **Lint: 0 issues** under the repo's 90+-linter config (was 22 on first run; fixed in code, not by suppression, except 3 documented nolint/#nosec with reasons). gofumpt clean, `go vet` clean. Config updated: dead `^internal/` exclusion rules removed, exhaustruct_v5 ignore-patterns rewritten to anchored struct-type patterns (the old path-based ones never matched anything — pre-existing latent misconfig).
10. **BuildFlow wiring (all committed by the auto-git daemon)**: require+replace in `go.mod` (root, indirect) and `tools/go.mod` (direct), blank import in `tools/providers/sdk_imports.go` with inventory comment, `go work vendor` (securitymd vendored), wiring guard test `TestSecuritymdProviderRegistered` (passes), AGENTS.md counts updated (118→119 steps, 9→10 toolsdk tools).
11. **Live end-to-end verification in BuildFlow**: `list providers` shows `securitymd detect+repair content`; `buildflow -s securitymd --fix` in a fresh git repo generated SECURITY.md ("1 fixed"); re-run detect = clean. Also verified on this repo's own hand-written policy: passes.
12. **README rewritten** for securitymd (rules table, exit codes, BuildFlow usage, tool_options example).
13. Dogfooding beyond tests: manual CLI loop in /tmp repos (dry-run → generate → validate → overwrite-refusal → missing-file finding).

## b) PARTIALLY DONE

1. **BuildFlow vendorHash**: `nix run .#update-vendor-hash` **failed** — the FOD fetches `github.com/LarsArtmann/securitymd@v0.0.0` from the network (local `replace` is not honored in the prepared-source build) and the repo does not exist on GitHub yet. Root cause understood, fix not applied (blocked on publishing, see questions). Note: BuildFlow's established sibling pattern is a **flake input + preparedSrc** (per GOTCHAS for licenseforge), which I did not replicate — I should have loaded the `nix-private-go-repos`/collector-extraction skills before choosing the replace-directive route.
2. **Docs**: README done. **Project AGENTS.md NOT rewritten** (still describes template-SECURITY/just/internal-layout — actively wrong now). CHANGELOG entry not added. `.github/workflows/security-validation.yml` still broken (Go 1.21 + `make`, both wrong) — known issue carried over, not fixed.
3. **flake.nix stale**: description "Security policy template for Go", devshell name `template-security-dev`. Cosmetic but lying.
4. **This status report** — being written now (counts as delivered with this file).

## c) NOT STARTED

1. Repo/GitHub rename `template-SECURITY` → `securitymd` (directory still `template-SECURITY`; BuildFlow replace points at the local path; `go install …@latest` impossible until renamed+pushed+tagged; then drop the replaces and re-vendor).
2. First release tag of securitymd (v0.0.1/v1.0.0) — needed for BuildFlow's nix build to resolve it without flake-input surgery.
3. BuildFlow GOTCHAS.md entry for the securitymd wiring (convention for toolsdk siblings).
4. BuildFlow `docs --check` doc-consistency run (its counts-vs-registry verifier may flag the 119 number or docs tables I didn't touch).
5. BuildFlow full workspace test gate (`nix run .#test` / erraudit) — I ran only `tools/providers` subset + root build.
6. Fleet rollout: running `buildflow -s securitymd --fix` across covered repos (every manifest-carrying repo without SECURITY.md will now gate-fail until generated — intended, but unannounced).
7. Trash-or-keep decision on stale root reports (IMPROVEMENT_PLAN.md, PARTS.md, PROJECT_SPLIT_EXECUTIVE_REPORT.md, BDD_TESTS_REVIEW.md, PUBLIC_OR_PRIVATE.md) and stale `docs/` planning docs from the old architecture.
8. FEATURES.md/TODO_LIST.md reconciliation for the rename (docs-health pass).
9. dependabot/SECURITY.md-of-this-repo regeneration with the new tool identity (repo's own SECURITY.md predates the rebuild and passes, but was hand-written under the old name).
10. golangci config's dead `exhaustruct` (v1-name) exclusion entry for `^cmd/` (leftover of the old config; harmless, cosmetic).

## d) TOTALLY FUCKED UP (honest accounting — all caught and fixed in-session)

1. **Wrote inverted section logic first** (`if !containsAnyFold → continue` emitted findings for PRESENT sections) — caught by reading my own diff one minute later, before tests existed.
2. **Left a garbage artifact** (`statOK`/`execLookStat`) in project.go from an abandoned approach — removed after re-reading.
3. **Wrong import path** for linter-autoconfigure-sdk (`…/autoconfigure` subpackage vs root package) — one wasted build cycle.
4. **Two `finding.ToolName` vs `string` type mismatches** (ToolInfo.Name) — compiler caught.
5. **`FixStrategyDirect` without fix content** — go-finding's builder validation rejected the missing-file finding at runtime; I discovered it via manual CLI testing rather than a test first (the provider tests now pin this).
6. **Test fixture bugs**: "compliant" fixtures were under the 20-line minimum (too-short fired); the reporting-removal case expected `missing-contact` that cannot fire when only the heading is removed. Both were _my_ fixture errors, not validator bugs — expectations fixed.
7. **`go test | tail` masked failures** — the first full-suite run showed only the last package's failure while pkg/policy was also failing; switched to filtering `ok|FAIL` lines. Lesson: never tail test output.
8. **multiedit staleness**: gofumpt reformatted provider.go between my read and my edit (edit refused — good); one validate.go batch applied 7/8 with a silent partial failure I only found via the compiler.
9. **Stray doc comment**: the DetectNamed removal left its explanatory comment attached to `Report`'s godoc (godoclint caught).
10. **Never delivered the 10-name list to the user mid-session** — it lived in my head/this report; the first half of the original ask should have been surfaced earlier.
11. **Did not anticipate the nix FOD/replace interaction** — the single biggest miss of the session: I wired the sibling via go.mod replace (works locally) without checking BuildFlow's flake-input pattern for siblings first, so `nix build` for BuildFlow is broken until the module is published or flake-input-wired.

## e) WHAT WE SHOULD IMPROVE (design/process reflections)

1. **Sibling-integration checklist**: before adding a Go sibling dep to BuildFlow, check flake inputs + preparedSrc (ADR-041) FIRST; go.mod replace is only ever the local-dev half of the wiring.
2. **Fixture honesty**: "valid" fixtures in validator tests should be generated-or-validated against the same threshold the rule enforces; my hand-written ones silently under-shot minLines.
3. **CLI exit-code semantics** are now 0/1/2 but the ecosystem's ternary ExitCodeByConfidence (go-linter-sdk) was considered and skipped — fine for a config-style tool, revisit if rules gain confidence variance.
4. **Trigger breadth**: the manifest trigger list is dependabot-inspired; docs-only repos with none of those manifests won't activate the tool. Possibly add `README.md` to the trigger set.
5. **Validator heuristics are still substring-based** (inherited from the old tool). A markdown-structure-aware pass (headings, tables) would raise precision; multi-signal patterns already reduced the worst FPs.
6. **`missing-file` severity=error fleet-wide** gates every non-compliant repo hard. That is the point — but it deserves an announcement + one `--fix` sweep day before it lands in CI defaults.
7. **Concurrent-session awareness**: another agent worked in BuildFlow during this session (dependency sweep, licenseforge integration) — its changes interleaved with mine in the same repo. My BuildFlow edits were additive and committed cleanly, but coordination there is luck, not process.
8. **`reports/` dir appeared mid-session** in this repo (coverage.out, jscpd-report.json) — not mine, left untouched, origin unknown (likely another session or hook).

## f) NEXT — prioritized backlog (~45 items)

**P0 — unblock BuildFlow nix build (pick one path)**

1. Publish: rename GitHub repo to securitymd, push, tag v1.0.0, drop both `replace` lines, `go work vendor`, `nix run .#update-vendor-hash`.
2. OR flake-input path: add securitymd as BuildFlow flake input + preparedSrc wiring (load `nix-private-go-repos` skill first).
3. Run `nix build .` for BuildFlow to green.
4. Run BuildFlow `docs --check` and fix any count/table drift (119 steps, 10 toolsdk).
5. Run BuildFlow full workspace test gate (`nix run .#test`) + erraudit (must exit 0).

**P1 — finish the rename & docs of this repo**
6. Rename local directory + GitHub repo to `securitymd`; update BuildFlow replace paths (or drop after publish).
7. Rewrite project AGENTS.md for the new architecture (commands, layout, conventions, known-good state).
8. CHANGELOG entry: the rebuild, breaking changes (config file gone, policy types gone, module path, binary name).
9. Fix `.github/workflows/security-validation.yml` (Go 1.27, no make, new paths) or replace with a buildflow-based workflow.
10. Update flake.nix (description, devshell name); optionally add a package output for the binary.
11. Trash stale root reports (IMPROVEMENT_PLAN.md, PARTS.md, PROJECT_SPLIT_EXECUTIVE_REPORT.md, BDD_TESTS_REVIEW.md, PUBLIC_OR_PRIVATE.md) pending your call.
12. docs-health pass: FEATURES.md, TODO_LIST.md, prune `docs/planning/` + `docs/modularization/` for the dead architecture.
13. Regenerate this repo's own SECURITY.md with the new tool (currently passes, but predates the rebuild).
14. Add `reports/` to .gitignore or investigate its origin.

**P2 — harden securitymd itself**
15. Decide default fleet contact policy (GitHub-advisory-only vs baked email via tool_options) and document it.
16. Consider `README.md` in the trigger manifest list so docs-only repos activate.
17. Add a SARIF golden-file test for the CLI output.
18. Add mutation/discrimination proof run for the rule table (temporarily break one rule, watch tests fail).
19. Markdown-structure-aware validation phase (headings/table parsing instead of substring heuristics).
20. `securitymd setup` interactive mode when no git remote (prompt for org/repo).
21. Support `--force`/`--regenerate` for setup (explicit overwrite with backup) — currently impossible by design.
22. Add `docs/SECURITY.md`-preferring repos config knob (canonical location choice).
23. Rule severity configuration (per-repo downgrade of missing-response-time etc.).
24. GitHub Action wrapper (`securitymd-action`) for non-BuildFlow consumers.
25. golangci-lint v2 plugin distribution of the detector (linter-building distribution doctrine) — low priority; BuildFlow is the primary consumer.
26. Suppression support: honor a `securitymd:ignore(rule) reason` comment in SECURITY.md.
27. Baseline/ratchet mode for incremental adoption.
28. Cache version cell (avoid re-running git describe per detect).
29. Property tests for parseGitRemote (fuzz remote URL shapes).
30. Localize descriptions of findings (they are English-only).

**P3 — BuildFlow-side polish**
31. GOTCHAS.md entry: securitymd wiring, replace-vs-flake-input lesson, publish dependency.
32. Update BuildFlow README/docs tables listing tools (docs --check will point them out).
33. Add securitymd to BuildFlow's overview dashboard if tools are enumerated there.
34. Guard test: assert `securitymd` tool_options validation error message text.
35. Consider `DependsOn: ["license-sync"]`? (both write repo-meta files; ordering irrelevant today, note for future collision rules).
36. Announce the fleet-wide gate change (missing SECURITY.md = error) before it hits default pipeline configs.

**P4 — ecosystem**
37. Feed the exhaustruct ignore-patterns fix (anchored struct patterns) back as a note to golangci-lint-auto-configure — the path-based pattern in this repo was a silent no-op.
38. Propose to linter-autoconfigure-sdk: rename `SaveJSONBytes` doc to make byte-faithful non-JSON use first-class (or add `SaveBytes`).
39. Record the "nix FOD ignores local replaces" lesson in crush-config lessons.md (it cost this session a failed FOD cycle).
40. go-release skill pass for v1.0.0 (tag on green CI, proxy propagation, pkg.go.dev check).
41. website-launch for securitymd once published (optional).
42. Consider adding securitymd detection to the crush-config skill's project-discovery checklist (so future sessions check SECURITY.md compliance).

## g) QUESTIONS (cannot be answered from the codebase)

1. **Publish path for the nix blocker:** rename + push + tag `securitymd` on GitHub now (I can prepare everything but cannot push), or wire it as a BuildFlow flake input + preparedSrc first and publish later? (Flake-input is more moving parts; publish unblocks fastest but is your call since it renames a public repo.)
2. **Default fleet contact policy:** keep GitHub-advisory-only links in generated policies, or bake a real address (e.g. `security@lars.software`) fleet-wide via tool_options/`.buildflow.yml`?
3. **Stale root reports** (IMPROVEMENT_PLAN.md, PARTS.md, PROJECT_SPLIT_EXECUTIVE_REPORT.md, BDD_TESTS_REVIEW.md, PUBLIC_OR_PRIVATE.md): trash as pre-rebuild history, or keep as archived status snapshots?

---

**State at report time:** template-SECURITY repo committed clean by the auto-git daemon; BuildFlow wiring committed (5a83c8603, 7657e7aca, a9e755b14); `go build`/`go test`/`golangci-lint` green in both repos for everything I touched; BuildFlow `nix build` broken pending publish-or-flake-input decision.
