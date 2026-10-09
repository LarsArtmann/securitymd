# Status Report — Pareto Execution Sweep + Brutal Self-Review

**Date:** 2026-10-09 16:48 CEST
**Scope of this report:** the "NOW GET SHIT DONE" execution sweep of `docs/planning/2026-10-09_15-37_pareto-publish-and-hardening-plan.md` (this session, ~15:40–16:45), plus the honest self-review the prompt asked for. Everything below is evidence-checked against this session's own runs; claims I could not verify are marked as such.
**Format note:** skill default is HTML; explicit `.md` path demand in the prompt wins (same override as the 15-11 report; not propagated into the skill).

**Gate state at writing time (all re-run this session, no filters):** build ✅ · full test suite (4 packages, zero skips) ✅ · golangci-lint 0 issues ✅ · gofumpt clean ✅ · `nix flake check` (format + package build + docs-gate) ✅ · docs-gate ✅.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | **M2 SARIF+JSON golden tests** — 6 goldens (clean/mixed/missing × json/sarif) pin the exact output contract; path entropy + path-hash IDs normalized; `-update` refresh documented | `pkg/policy/golden_test.go`, `pkg/policy/testdata/golden/` |
| 2 | **M3 mutation/discrimination proof** — 3 sabotage targets (missing-contact patterns emptied → 7 red; too-short 20→0 → 1 red; unresolved-template off-by-one → 2 red incl. golden); revert green; run in a `git archive HEAD` copy so the concurrent session's WIP couldn't pollute the baseline; recorded in FEATURES | `FEATURES.md` mutation-proof note |
| 3 | **M4 OSS survey** — no dedicated SECURITY.md validator CLI exists (npm/PyPI/GitHub searched); Scorecard named as closest neighbor; ROADMAP claim upgraded from "unverified belief" to verified-with-caveat wording | `docs/status/2026-10-09_15-57_oss-landscape-survey.md` |
| 4 | **M6 README.md trigger** — docs-only repos activate; spec-shape coverage test | `pkg/provider/provider.go`, `provider_test.go` |
| 5 | **M8 fuzz + parser hardening** — 10.9M execs green; scheme allowlist (https/ssh/git@), clean-coordinates invariant, ≥3 segments, slash/CRLF tolerance; 7 regression table cases | `pkg/policy/project.go`, `project_fuzz_test.go` |
| 6 | **M10+M24 docs-gate** — un-annotated-archive gate + dangling `docs/` reference alarm as flake app AND flake check; negative-tested in both directions | `flake.nix` (`docsGate`), verified fail+pass runs |
| 7 | **M14 flake package** — `packages.<system>.securitymd`, Go 1.27 pinned (nixpkgs default is 1.26.8 — would fail go.mod), git fixture hermeticity fixed for the FOD; binary e2e smoke (setup → validate clean) | `flake.nix`, `vendorHash` pinned |
| 8 | **M15 interactive setup** — TTY-gated org/repo prompt; proven inert on piped stdin (CI-safe property pinned by test) | `cmd/securitymd/prompt.go`, `prompt_test.go` |
| 9 | **M16 `--force`** — in-place regeneration with timestamped `.bak` + restore hint; refuse-by-default untouched (existing test still green) | `pkg/policy/generate.go`, `generate_options_test.go` |
| 10 | **M17 `--location`** — write target (root/.github/docs) + candidate-order override for `validate`; unknown values are static-sentinel errors | `pkg/policy/generate.go`, `policy.go`, `cmd/securitymd/{setup,validate}.go` |
| 11 | **M23 scope verdicts** — golangci-plugin: won't (architecture mismatch); baseline/ratchet + localization: defer with revisit triggers | `ROADMAP.md` |
| 12 | **M1 runbook staging** — every mechanical step pre-verified (mod tidy clean, local install path, BuildFlow replace/flake-input line numbers); Lars's part reduced to 4 commands | `docs/planning/2026-10-09_16-34_publish-runbook.md` |
| 13 | **M22 upstream feedback** — verified at source: exhaustruct note has **no defect to file** (curated patterns already anchored type-form, `linter_settings.go:244`); `SaveJSONBytes` already byte-faithful in SDK v0.8.0; the two real items written to crush-config (lessons.md entry + nix-private-go-repos gotcha row) | `ROADMAP.md`, `/home/lars/projects/crush-config` |
| 14 | **M5/M7/M9 (concurrent session's work)** — suppressions, severity overrides, version-cell caching landed and verified green in the full suite; I consolidated on their M9 design and dropped my duplicate cache | `pkg/policy/suppress.go`, `severity.go`, `version_cache_test.go` |
| 15 | **Docs sync** — TODO_LIST trimmed to open items, FEATURES gained rows 9–17, CHANGELOG [Unreleased] batch entry, AGENTS commands/testing/architecture refreshed, README gained the new CLI surface | all five files |

## b) PARTIALLY DONE

1. **M17 symmetry gap (found in self-review):** `validate --location` reorders detection, but `status` was never given the flag — `status` always uses the default order. One-flag fix, real inconsistency.
2. **M24:** the gate exists and runs in `nix flake check`, but the fleet cron (run it per-repo on a schedule) and the `docs/reviews/` convention decision remain open (Lars).
3. **M1:** runbook staged and pre-verified; the publish itself is unexecuted by design (Lars).
4. **M22:** resolved-as-moot for two claims, written for two — but nothing was actually *filed* upstream (correctly, per verification; noting so nobody expects issues to exist).
5. **FEATURES row 15 (CLI tests):** prompt helpers are now unit-tested; cobra wiring itself still only covered via acceptance/provider tests.
6. **Suppression coverage in goldens:** the goldens pin unsuppressed output only; a suppressed-finding fixture (🔇 + exit-neutral through SARIF/JSON) is not yet pinned.
7. **README rule tables vs `sectionRules` table vs FEATURES row 1:** three hand-maintained lists of the same 11 rules — drift-guarded only by accident (goldens pin behavior, not docs). See (e).

## c) NOT STARTED (all with known gates)

1. **M11 fleet sweep + announcement** — Lars timing gate. Honest miss: the *draft* announcement and the sweep script (F11.1/F11.2) were preparable without the gate; I skipped them entirely.
2. **M12 BuildFlow post-sweep nix verify** — waits on BuildFlow's `execution/` split landing.
3. **M13 BuildFlow hygiene** (gotcha re-anchor, workspace gates, tool_options guard test) — BuildFlow repo.
4. **M18 own SECURITY.md regen + metadata.yaml tags** — post-rename (Lars).
5. **M19 website / M20 GitHub Action** — post-publish by design.
6. **M21 structure-aware validation** — deliberately deferred (M5-first ordering satisfied; big parity-tested port).
7. **metadata.yaml tags, contact policy, fleet announcement timing** — Lars decisions in TODO_LIST.

## d) TOTALLY FUCKED UP

Nothing shipped is broken — but three things came close, and one documented lie exists:

1. **Exit-code lie (pre-existing, found now):** README + plan say "exit 2 = operational failure". `cmd/securitymd/main.go:38` exits **1 for everything**; exit 2 is unreachable. Findings-vs-crash are distinguishable only by stderr text. Needs a real fix (`errPolicyFindings` → 1, other errors → 2) + a test. This survived two docs-health passes because nobody executed the documented contract.
2. **docs-gate first version had a false-failure bug:** `set -o pipefail` + grep's exit-1-on-no-match made the gate fail for any living doc without backticked `docs/` refs. It passed locally by luck (every doc happened to have refs) and only failed in the flake check against the older store tree. I negative-tested *after* fixing — the order should have been negative-test-first.
3. **M9 duplicate implementation:** I built a public-`LatestTag` cache while the concurrent session built a private `lookupLatestTag` — two caches for one concept, briefly both in the tree. Caught because their `version_cache_test.go` appeared mid-build; without that, this ships as a split brain.
4. **Mutation-proof baseline pollution:** my first isolated copy carried the other session's failing WIP tests, which would have made the proof record meaningless (can't show red→green over a red baseline). Caught before recording; redone from `git archive HEAD`.

## e) WHAT WE SHOULD IMPROVE (brutal answers to "what did you forget / what could be better")

1. **Concurrency protocol.** A second Crush session worked the same repo the whole time (M5/M7/M9), caused 3 stale-read refusals, one duplicate implementation, and one invalidated proof run. Cheap fix: before starting, `git status` + check for concurrent writers; claim files by touching a `.session-claim` note or simply re-checking `git diff HEAD` before each task. I adapted mid-flight instead of preventing.
2. **Negative-test-first for gates.** The docs-gate should have been proven to FAIL before being proven to pass. Same discipline as M3, applied to scripts — I applied it late.
3. **Fixture split brain.** `validPolicy` (validate_test), `flawedPolicy` (golden_test), `compliantPolicy` (acceptance) are three hand-rolled policies describing the same contract. Move to `testdata/*.md` fixtures with one canonical + deliberate mutations.
4. **Write-then-caught vs know-then-write.** This session's compile errors: `mustRead` arity (twice), a glued `}` from a sloppy edit, a dead import, an unused param, mnd's `3`, wsl whitespace, err113, cyclop 16/15, nix `inputs` arg, nix backtick escaping, cobra nil-context panic, wrong smoke-test subcommand. Each caught fast, but a dozen linter/tax rounds is the linter doing my reading. Slower drafting against the known lint config (mnd ignored-numbers, err113, wsl cuddle rules, exhaustruct patterns) would save a cycle per task.
5. **The plan's fine tasks were treated as all-or-nothing per Lars gate.** M11's draft/script and parts of M24 were preparable without their human decisions — "blocked" and "preparable" were conflated. Future plans: split each gated task into "preparable now" vs "needs the decision".
6. **Verify documented behavior, not just shipped behavior.** The exit-2 lie survived because every verification pass ran tests, never the README's contract table. A tiny docs-contract test (table-driven: scenario → expected exit code, asserted via CLI subprocess) would pin README/FEATURES claims mechanically.
7. **Go LSP is broken in this environment** (`documentSymbol` unsupported → `lsp_replace_symbol` unusable); I fell back to text edits, which caused the glued-brace bug. Worth fixing the toolchain, else always re-view after symbol-level edits.
8. **Suppression×severity interaction untested at the edges:** suppressed-then-overridden severity ordering (`--set-severity` + `securitymd:ignore` on the same rule) has no explicit test. The code paths exist independently; the combination is unproven.

## f) NEXT 50 (sorted by impact; starred = this repo, do-now class)

1. ★ Fix the exit-code contract: findings → 1, operational failure → 2 (`main.go`), + subprocess test pinning the README table
2. ★ `status --location` flag (M17 symmetry)
3. ★ Suppressed-finding golden fixture (JSON + SARIF, 🔇/exit-neutral contract)
4. ★ Policy fixtures → `testdata/*.md` (kill the 3-way fixture split brain)
5. ★ Rule-table single source: test asserting README's rule table == `sectionRules` + `missing-file` (doc-drift guard)
6. ★ Suppression × severity-override combination test
7. ★ Wire docs-gate + fuzz smoke (30s) into CI workflow
8. ★ govulncheck in CI (gosec already runs inside golangci-lint; verify vulncheck separately)
9. ★ Coverage gate vs the 80% target (`go test -cover` in CI, floor pinned)
10. ★ Prepare F11.1/F11.2 now: fleet announcement draft + `buildflow -s securitymd --fix` sweep script (execution stays Lars-gated)
11. **P0 publish** (Lars: 4 commands, runbook ready)
12. Regenerate own SECURITY.md + metadata.yaml tags (M18, post-publish)
13. Provider `location` toolsdk option (BuildFlow repos pick .github/docs canonically)
14. `validate --repair`: detect+generate one-shot for CLI users
15. Suppression expiry: optional `until YYYY-MM-DD` in the grammar (go-finding `IsSuppressedAt` already supports time-based suppression)
16. `--no-color` / NO_COLOR confirmation + docs
17. Shell completions + man pages (cobra built-ins, cheap)
18. Windows CI leg (`.bak` timestamp format is colon-free; verify the rest)
19. Benchmarks: Detect on a large repo (perf baseline before M21)
20. M12: BuildFlow `nix build .` + nix-binary e2e after their sweep lands
21. M13: BuildFlow gotcha re-anchor + tool_options error-message guard test
22. Post-publish: pkg.go.dev render check + repo description/topics
23. Post-publish: GoReleaser or flake-based release workflow
24. Post-publish: M20 GitHub Action (`securitymd-action`)
25. Post-publish: M19 website (demo video centerpiece)
26. M21: structure-aware validation (goldmark port, parity-gated) — after 3–6 stabilize behavior
27. Pre-commit hook recipe in README (validate on commit)
28. Dependabot/renovate for fleet repos (post-publish)
29. nixpkgs PR packaging securitymd (post-publish, flake gives the derivation)
30. Cross-check `.golangci.yml` build-tags/go version vs go.mod on toolchain bumps (lint config pins 1.26.7, go.mod 1.27 — currently tolerated, should be synced)
31. `docs/DOMAIN_LANGUAGE.md`: add suppression/severity/location terms
32. Decide `.bak` gitignore policy for repos under sweep (fleet guidance line in the announcement)
33. churn guard: `go work vendor` drift check in BuildFlow (their side)
34. Read-model idea: `securitymd report --since` diffing two SARIF files (ratchet-lite without state)
35. Template: org-level SECURITY.md inheritance (central policy, per-repo overrides) — ROADMAP fuel
36. `securitymd list-rules` (machine-readable rule metadata for CI config UIs)
37. Audit `--severity` vs `--set-severity` flag-name confusion (rename or alias)
38. Golden test for `--format json` CLI exit/stdout combination (contract beyond Report level)
39. Add `exit 2` semantics to BuildFlow provider error mapping (verify toolsdk error taxonomy fit)
40. Explore gitleaks-style config-free posture claim in README (verify before claiming "no config file" is a feature)
41. Update crush-config project-discovery checklist question (ROADMAP open question) into a decision doc
42. Fix the 2 dangling skill links (`collector-extraction`, `linter-building` → `~/.agents/skills/`) in the crush-config fan-out
43. Repair Go LSP (documentSymbol) in this environment or document edit-tool fallback as standing practice
44. Add `securitymd` to crush-config's lessons index (nix-FOD lesson now exists; make it discoverable)
45. Fuzz `parseSeverityFlag` / severity-override parser (same corpus discipline as M8)
46. Property test: generated policy never contains `{{` (already asserted per-test; promote to fuzz over GenerateOptions)
47. Version injection: flake package should stamp git rev, not the static `0.1.0-dev`
48. SARIF `automationDetails`/run metadata: consider embedding repo identity (go-finding upstream discussion first)
49. Fleet rollout metric: count repos failing the gate before announcement (informs severity-default decision)
50. Close the loop on this report's section (f) → docs-health HARVEST into TODO_LIST/ROADMAP

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **Publish now or after BuildFlow's `execution/` split settles?** The runbook is ready and the tool is FOD-proven either way; but if the sweep lands soon, one publish + one BuildFlow re-vendor cycle is cheaper than two. Which order do you want?
2. **Exit-code contract:** should operational failure really be exit 2 (README's promise), or is exit 1-for-everything acceptable and the README should change? (This decides whether item 1 above is a code fix or a docs fix — I can argue both: 2 distinguishes crash from findings for CI scripting; 1 is what every consumer already observes.)
3. **Suppression expiry:** do you want `securitymd:ignore(rule) until YYYY-MM-DD reason` (auto-re-exposing stale suppressions via go-finding's time model), or is that machinery you'd rather not carry? It changes the grammar we just shipped and its golden fixtures.

---

*Point-in-time snapshot. Section (f) items are routing candidates for docs-health HARVEST — items 1–10 are TODO_LIST class, the rest are ROADMAP fuel until they earn promotion.*
