# Status Report — Fleet-Policy Compliance Pass (library-policy violations, PUBLIC flip)

**Date:** 2026-10-09 21:13 CEST
**Session scope:** `securitymd` repo — visibility flip PRIVATE→PUBLIC, then fix both `library-policy` violations (testify ban + cobra missing companions), full BuildFlow/nix/docs gate recovery.
**Format note:** HTML is this skill's canonical output; user explicitly demanded `.md` at this path — honored as instructed override.

---

## Executive summary

Both library-policy violations are fixed and the scanner passes. The CLI now runs on cmdguard v4 + fang (fleet house pattern), testify is fully removed in favor of Gomega, and every quality gate is green: build, full test suite, `nix flake check` (new vendorHash), golangci-lint 0 findings, full `buildflow` gate PASSED (warnings only), `library-policy` PASSED. Three gate-fallout surprises (stale vendorHash short-circuit, daemon-committed binary, embedded structure-linter ignoring project config) were diagnosed and resolved in-session and recorded as AGENTS.md gotchas.

---

## a) FULLY DONE

| #  | Work                                                                                                                                                                                                                                                                                             | Verification                                                                                                         |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------- |
| 1  | **Repo flipped PUBLIC** (`gh repo edit --visibility public`)                                                                                                                                                                                                                                     | `gh repo view` → `PUBLIC`; `go list -m github.com/LarsArtmann/securitymd@v1.0.0` resolves via proxy.golang.org       |
| 2  | **cmdguard v4.1.0 + fang adoption** (satisfies cobra-companions policy): `cmdguard.NewCLI[appConfig]` root, typed flag structs (`validateFlags`/`setupFlags`/`statusFlags`) with `flag:`/`values:` enum tags, `WithNoArgs`, panic-recovery middleware, signal handling, `--no-color` global flag | Build + binary smoke test (`validate` exit 1, `status` exit 0, `--version` string preserved, piped output ANSI-free) |
| 3  | **Exit-code contract preserved** on the new execution stack: findings → 1, operational → 2; cmdguard prints errors exactly once, `main` maps codes only                                                                                                                                          | All 16 contract scenarios green (`cmd/securitymd/contract_test.go` subprocess re-exec)                               |
| 4  | **testify → Gomega migration**, 15 test files (12 `pkg/policy`, 1 `pkg/provider` + helper, 2 `cmd/securitymd`), testify gone from go.mod and code                                                                                                                                                | `rg 'stretchr/testify'` → 0 hits; full suite green                                                                   |
| 5  | **library-policy scanner**: `PASSED — 0 violations` (122 banned libraries enforced)                                                                                                                                                                                                              | Re-ran after all changes                                                                                             |
| 6  | **vendorHash rotation** for the new dependency tree: `sha256-0sPYOSiuReUMrhH/RDC2+l0U7RAxPYHtbkRPMv/0y9c=` (fakeHash → capture got-hash → paste, the canonical nix flow)                                                                                                                         | `nix build .#securitymd` green (12.7 MB binary)                                                                      |
| 7  | **`nix flake check` all passed** (docs-gate + treefmt + package) — after gofumpt import-order fix                                                                                                                                                                                                | CLI output: "all checks passed!"                                                                                     |
| 8  | **golangci-lint 0 findings** via `buildflow -s golangci-lint`; varnamelen `g` exemption added to `.golangci.yml` (gomega receiver; same exemption cmdguard curates)                                                                                                                              | buildflow step green                                                                                                 |
| 9  | **Full `buildflow` gate PASSED** (49/56 steps, warnings only: link-scan 404s + small token dupes; zero error findings) with `BUILDFLOW_NO_RESULT_CACHE=1`                                                                                                                                        | Gate verdict line                                                                                                    |
| 10 | **Daemon-committed binary removed from git** (`bin/securitymd` had been auto-committed — Rule 007 violation): `git rm --cached` + `bin/` gitignored                                                                                                                                              | `git ls-files bin/` empty                                                                                            |
| 11 | **Stale gitignored `vendor/` deleted** (`trash`) — it was a relic causing 73 gomod-check findings                                                                                                                                                                                                | findings gone on re-run                                                                                              |
| 12 | **flake.nix meta license** added (`licenses.mit`, matching LICENSE)                                                                                                                                                                                                                              | flake-meta-checker finding gone                                                                                      |
| 13 | **CI coverage gate hardened** with `-race` (buildflow race-detector check demanded a configured home)                                                                                                                                                                                            | `.github/workflows/security-validation.yml:31`                                                                       |
| 14 | **`.buildflow.yml` + `.go-structure-linter.yaml`** created: skip + documented rationale for the `pkg/`-is-public false positive (fleet precedent: go-error-family)                                                                                                                               | gate green; rationale recorded                                                                                       |
| 15 | **Docs refreshed**: AGENTS.md (PUBLIC flip, cmdguard CLI architecture, Gomega testing, 3 new gotchas, known-state refresh), FEATURES.md rows 17+20, README architecture line                                                                                                                     | readme_drift_test still green                                                                                        |

## b) PARTIALLY DONE

| # | Item                                              | State                                                                                                                                                                                                                                                                                                                                                 |
| - | ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **Contract-test coverage of the new CLI surface** | The 16 existing scenarios pass, but `--version`, bare `securitymd` (help), unknown-command, and piped-`--help` ANSI-freeness have NO scenarios. One scenario's assertion was deliberately loosened to fang's framing (pins domain message `unknown policy location "bogus"` instead of full original line) — legitimate but a real contract reshaping |
| 2 | **README accuracy post-swap**                     | Verified: no "invalid --location" wording in README (no drift), rule tables guarded by drift test. NOT verified/documented: the new global `--no-color` flag and fang's restyled help/error output shape                                                                                                                                              |
| 3 | **BuildFlow-side integration after this session** | BuildFlow's flake input tracks latest — this session changed go.mod/vendorHash on the consumed repo, so BuildFlow's own nix FOD likely needs its `nix-hash-fix` + a re-verified detect→repair E2E. Not started (their repo, their gate)                                                                                                               |
| 4 | **Known-state docs**                              | AGENTS.md "Known state" refreshed; TODO_LIST.md NOT harvested from this report yet (docs-health HARVEST pending), CHANGELOG.md untouched — the cmdguard swap has no changelog entry                                                                                                                                                                   |
| 5 | **Binary size tradeoff**                          | Measured (12.7 MB, charm stack aboard), called "acceptable" — but never surfaced as a decision for Lars or documented in FEATURES                                                                                                                                                                                                                     |

## c) NOT STARTED

- gitleaks run on the now-PUBLIC repo (`buildflow -s gitleaks` — see self-review: my pre-flip scan was a thin regex pass, not the fleet secret scanner)
- TODO_LIST/ROADMAP harvest of this report's section (f)
- Contract scenarios for `--version` / bare command / unknown command / piped help
- Acceptance-suite (Ginkgo) coverage of the CLI layer (it exercises `pkg/` only — zero coverage of the cmdguard rewrite; the contract subprocess test is the sole guard)
- cmdguard upgrade watch (v4.1.0 pinned) + decision whether its auditlog/go-health companion plugins matter here
- Diagnosing buildflow's "9 tools unavailable (health check failed)" warning (interrogate et al.)
- aarch64 flake checks (`nix flake check` omits them: "use --all-systems")
- CI proof run (still blocked on Lars's GitHub Actions billing — pre-existing)
- Fleet sweep execution + announcement draft update (Lars-gated, staged pre-session)

## d) TOTALLY FUCKED UP

Nothing is currently on fire — but two things WERE, and both are honest self-inflicted wounds:

1. **A compiled binary entered git history.** I ran `go build -o bin/securitymd` inside the repo without checking ignore coverage; the auto-commit daemon committed it. A fleet Rule 007 violation created by me, caught by buildflow's structure linter, not by me. Fixed in-session (`git rm --cached` + gitignore), but the blob lives in history until a rewrite Lars almost certainly doesn't want for a 12 MB blob. Lesson: build artifacts go to `/tmp` or an ignored path, always.
2. **The PUBLIC flip used a hand-rolled regex secret scan instead of gitleaks.** One `rg` pass for common token patterns before an irreversible-ish visibility change is not the fleet's secret scanner. No findings — but the rigor was beneath the moment. The repo is public NOW; the compensating action (`buildflow -s gitleaks`) is queued as next-step #1.

## e) WHAT WE SHOULD IMPROVE (self-review, brutally honest)

**What did I forget?**

- The `bin/` artifact hygiene check before building (above).
- gitleaks before the flip (above).
- To check what NEW flags/commands cmdguard injects onto the CLI surface (`--no-color` appeared; nothing else did — verified after the fact, should have been immediate).
- That "acceptance tests pass (cached)" means they never touch `cmd/` at all — my "all green" claim had a coverage hole I only noticed while writing this report.
- CHANGELOG.md — the one living doc I never updated despite a user-visible CLI framework swap.

**What is stupid that we do anyway?**

- Hand-rolling `exitCodeFor` when the fleet owns `go-error-family` (built for exactly this: sysexits mapping, ExitCoder semantics). Works today; fleet-standard migration candidate.
- The `.go-structure-linter.yaml` is currently a ghost: the embedded buildflow step ignores it and the standalone CLI isn't wired into CI. It's a rationale document pretending to be a config.
- Golden-file tests cover JSON/SARIF but text output — the thing humans and CI logs actually read — is only substring-pinned.

**What could I have done better?**

- **Read the whole rule before fixing it.** I jumped from the scanner output to cmdguard adoption; reading `library-policy/policies/companions.go` first would have shown the exclusion is presence-based and let me decide depth deliberately instead of discovering the execution-layer semantics (cmdguard's `rule_execute_bypass` lint AGAINST `fang.Execute` on cmdguard roots) mid-rewrite.
- **Failure triage discipline on nix.** I burned 3–4 tool calls on derivation forensics before the insight (stale hash ⇒ FOD short-circuit ⇒ no got-hash). `buildflow doctor` and the failure-triage reference were the prescribed first moves.
- **Scope the smoke test to its lifetime.** The `bin/securitymd` smoke test was fine; its location wasn't.
- **Verify parity claims at the surface.** I said "command surface parity" before ever listing `--help` output. (Now verified: completion/help/setup/status/validate + `--no-color`/`-v`.)

**What could we still improve?**

- Contract philosophy under fang: pin domain messages (current) AND add a custom fang error handler if we ever need byte-stable output for downstream parsers — decide explicitly, not by drift.
- Migrate exit mapping to `go-error-family`; delete the hand-rolled `errors.Is` ladder.
- Wire the standalone `go-structure-linter .` into CI so `.go-structure-linter.yaml` stops being a ghost.
- Add go-snaps snapshot tests for `validate` text output (fleet testing standard already lists go-snaps).

**Ghost systems?** One, borderline: `.go-structure-linter.yaml` (above). Value = recorded rationale; it should either gain a runner (CI step) or be explicitly demoted to a comment in `.buildflow.yml`.

**Split brains created or found?**

- Structure-lint policy now lives in TWO files (`.buildflow.yml` skip + `.go-structure-linter.yaml` rationale) — fleet-known pattern, still a drift surface.
- Exit-code knowledge lives in README table + contract test + `exitCodeFor` — pinned, acceptable.
- Flag documentation: struct tags (source of truth) vs README flag lists — UNGUARDED. The rule-table drift guard should get a flags-table sibling.

**Scope creep?** Bordered on it: a dependency-policy fix became a full CLI framework migration. Justified by the policy's intent and the contract safety net — but the honest sequencing alternative (fang-only for the immediate finding, cmdguard as a tracked follow-up) existed and I chose depth. This session had the budget for it; note the pattern.

**Did we remove something useful?** `promptWhenIdentityMissing`'s package-var mutation became pure `resolveIdentity(flags)` — strictly better. The stale `vendor/` dir and the committed binary were pure liabilities. Nothing useful lost. One behavior change to own: cmdguard's `--no-color` flag and fang's error framing are new user-visible surface; `WithNoArgs` now rejects stray positional args that cobra silently ignored (strictness, unpinned by contract).

**Tests: how are we doing?** Suite green and honest at the unit/contract layer; gaps: CLI-in-acceptance (zero), text-output snapshots, fuzz on `ParseSeverityOverrides`, TZ-boundary test for `until` expiry (pinned UTC, single-TZ only), and no coverage measurement taken this session (80% floor asserted by CI config, never re-measured post-migration).

**Did I lie?** No — but two claims were softer than they sounded until verified: "command surface parity" (verified only now, via `--help`) and "all green" (true per-gate, while the acceptance-cache subtlety above hid a real gap). Both closed in this report.

---

## f) NEXT 50 (impact-sorted brainstorm; most items are ROADMAP fuel, not commitments)

| #  | Task                                                                                                             | Impact                                               |
| -- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| 1  | Run `buildflow -s gitleaks` on the now-public repo                                                               | HIGH — open exposure window since flip               |
| 2  | Re-run full gate on FRESH buildflow binary (concurrent session shipped f9af40d; advisory seen)                   | HIGH — verdicts may shift                            |
| 3  | Sync BuildFlow side: their flake input pulls this session's revs → their `nix-hash-fix` + provider E2E re-verify | HIGH — fleet integration is the repo's raison d'être |
| 4  | Contract scenarios: `--version`, bare command, unknown command, piped `--help` ANSI-free                         | HIGH — pins the new execution surface                |
| 5  | HARVEST this list into TODO_LIST/ROADMAP (docs-health)                                                           | HIGH — anti-entombment                               |
| 6  | CHANGELOG entry for cmdguard/fang swap + PUBLIC flip                                                             | MED-HIGH                                             |
| 7  | README: document `--no-color` + fang help/error output shape                                                     | MED-HIGH                                             |
| 8  | Version bump decision: tag v1.1.0 for the CLI swap, or bundle with fleet rollout                                 | MED-HIGH — Lars-gated                                |
| 9  | Add Ginkgo acceptance spec exercising CLI via exec (close zero-acceptance-coverage gap)                          | MED-HIGH                                             |
| 10 | Flags-table drift guard (README flags vs struct tags, sibling of rule-table guard)                               | MED                                                  |
| 11 | CI: measure + publish real coverage % (80% floor never re-measured post-migration)                               | MED                                                  |
| 12 | go-snaps snapshot for `validate` text output                                                                     | MED                                                  |
| 13 | Wire standalone `go-structure-linter .` into CI (de-ghost the yaml)                                              | MED                                                  |
| 14 | Link-scan 404: README links `LarsArtmann/BuildFlow` (private repo) → replace with lars.software or drop          | MED                                                  |
| 15 | Migrate exit mapping to `go-error-family`, delete hand-rolled ladder                                             | MED                                                  |
| 16 | Binary size note in FEATURES + measure startup latency delta (charm init cost)                                   | MED                                                  |
| 17 | Fuzz `ParseSeverityOverrides` (second parser, same class of risk)                                                | MED                                                  |
| 18 | Fuzztime sync: CI runs 30s, AGENTS.md says 45s                                                                   | LOW-MED (doc truth)                                  |
| 19 | Contract scenario: malformed `--set-severity` value → exit 2                                                     | LOW-MED                                              |
| 20 | `resolveIdentity` FirstExisting short-circuit unit test (path uncovered)                                         | LOW-MED                                              |
| 21 | Diagnose buildflow "9 tools unavailable" health warnings (interrogate et al.)                                    | LOW-MED                                              |
| 22 | `nix flake check --all-systems` for aarch64                                                                      | LOW-MED                                              |
| 23 | CI: SARIF upload action (a SARIF-emitting linter should feed code scanning)                                      | MED                                                  |
| 24 | CI green badge in README once billing fixed + first run passes                                                   | LOW-MED                                              |
| 25 | pkg.go.dev appearance check post-flip (license detected, docs render)                                            | LOW-MED                                              |
| 26 | Dependabot config present/current for cmdguard tracking (buildflow step applicability)                           | LOW-MED                                              |
| 27 | Verify `securitymd completion bash` still works post-swap                                                        | LOW-MED                                              |
| 28 | Dogfood: regenerate repo SECURITY.md with new binary (version-cell drift check)                                  | LOW                                                  |
| 29 | cmdguard docgen → man pages in docs/                                                                             | LOW                                                  |
| 30 | Deduplicate 12–20-token test warnings via shared helpers                                                         | LOW                                                  |
| 31 | Pre-commit hook: confirm `buildflow precommit install` state                                                     | LOW                                                  |
| 32 | TZ-matrix (or UTC-pinned assertion doc) for `until` expiry boundary                                              | LOW                                                  |
| 33 | Announcement draft: update for cmdguard/fang adoption before fleet rollout                                       | MED (pre-rollout)                                    |
| 34 | `fleet-securitymd-sweep.sh` dry-run against 2–3 repos                                                            | MED (pre-rollout)                                    |
| 35 | cmdguard upgrade watch (v4.1.0 → note in ROADMAP)                                                                | LOW                                                  |
| 36 | Windows console behavior note (no Windows CI exists; color/fang assumptions)                                     | LOW                                                  |
| 37 | Golden test: JSON error envelope for mid-write failures (edge)                                                   | LOW                                                  |
| 38 | `.buildflow.yml` config lint (`buildflow config lint`) for unknown-key warnings                                  | LOW                                                  |
| 39 | ROADMAP cross-link: structure-aware suppression parser (existing item) stays tracked                             | LOW                                                  |
| 40 | Consider `WithConfigFile` for org-wide defaults (YAGNI-flag: only on real demand)                                | LOW                                                  |
| 41 | Consider `securitymd doctor` via cmdguard health checks (YAGNI-flag)                                             | LOW                                                  |
| 42 | Release automation: GoReleaser evaluation (go-release skill) when v1.1.0 lands                                   | LOW                                                  |
| 43 | nix devShell audit: gofumpt/golangci-lint availability as documented                                             | LOW                                                  |
| 44 | GitHub social preview + topics polish post-flip                                                                  | LOW                                                  |
| 45 | Archive this report series + annotate when superseded (docs-health ANNOTATE)                                     | LOW                                                  |
| 46 | TODO_LIST BuildFlow section: re-verify `filterFindingsAtOrAbove` still unfixed upstream                          | LOW                                                  |
| 47 | Error-output stability decision: custom fang error handler vs domain-message pinning (document choice)           | LOW                                                  |
| 48 | `WithNoArgs` strictness: confirm no documented usage relies on stray args                                        | LOW                                                  |
| 49 | Contract test count grows → consider table split for runtime                                                     | LOW                                                  |
| 50 | Post-billing CI proof: `gh workflow run "Security Policy Validation"` (Lars-gated)                               | MED when unblocked                                   |

---

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **Binary size vs policy strictness:** the cmdguard companion rule pulled the full charm stack into a fleet security scanner (12.7 MB, ~+40%). Is that an accepted fleet-wide cost, or should the companion policy gain a "lightweight CLI" carve-out (fang-only) that securitymd-sized tools could use?
2. **Release cadence for the CLI swap:** tag `v1.1.0` now (cmdguard/fang is a user-visible runtime change behind a compatible contract), or hold the tag until the fleet rollout proves the new binary across repos?
3. **CI after the PUBLIC flip:** with billing still blocking GitHub Actions, do you want the first proof run prioritized (billing fix → `workflow_dispatch`), or is CI dormant until the fleet-gate announcement lands?

---

_Point-in-time snapshot 2026-10-09 21:13 CEST. Auto-commit daemon owns commits (harness contract: no manual commit without explicit user request). Section (f) is the HARVEST input for TODO_LIST/ROADMAP._
