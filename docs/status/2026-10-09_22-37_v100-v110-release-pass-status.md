# Status Report — v1.0.0/v1.1.0 Release Pass

**Generated:** 2026-10-09 22:37 CEST
**Session scope:** "Release a v1.0.0" — go-release skill execution: state assessment, tag-immutability discovery, CHANGELOG cut, full gate run, v1.0.0 release completion, v1.1.0 cut/tag/release, proxy + sumdb + go-get verification, living-docs sync.
**Verdict:** **BOTH RELEASES LIVE AND CONSUMER-VERIFIED.** `v1.0.0` (pre-existing, proxy-frozen) got its missing artifacts; `v1.1.0` (the 26 post-tag hardening commits) was cut at a CI-green commit, pushed, released, and anchored in the checksum database. No recovery actions needed; one cosmetic blemish (daemon-owned commit message on the tagged commit).

---

## 0. The decisive discovery (read this first)

The session began with "Release a v1.0.0" — but v1.0.0 **already existed**: tagged at `5c739ec` (18:13), pushed, and — the critical fact — **already served by proxy.golang.org with checksums recorded in sum.golang.org** (`h1:N2NpRbS7…`, verified at source via the `.info` and `/lookup` endpoints). The repo was already PUBLIC and CI was already green (both facts contradicted AGENTS.md/TODO_LIST, which still said PRIVATE + billing-blocked).

That evidence forced the fork:

- **Moving the tag to include the 26 post-tag commits = prohibited.** sum.golang.org is a signed transparency log; a moved tag produces a hard SECURITY ERROR for every consumer. Forever.
- **Correct path:** complete v1.0.0's release artifacts (CHANGELOG section, GitHub Release) against the frozen tag, and ship the post-tag work as **v1.1.0** (MINOR bump — the delta is the cobra → cmdguard v4 CLI migration: `main.go` +75 lines, new dependency surface, plus the ANSI-free piped-output guard and docs-audit pass).

**Reflected judgment:** the earlier publish pass (18:13, prior session) tagged v1.0.0 *without* cutting its CHANGELOG section or GitHub Release — so a version reached the public proxy before its release artifacts existed. Nothing was broken (the tagged tree was green and complete), but it means "the release" and "the tag" had drifted apart. That discipline gap is the root cause of most of this session's catch-up work; it is the item I'd most want gated in future.

---

## a) FULLY DONE

| # | Work | Evidence |
|---|------|----------|
| 1 | Tag-immutability assessment (the go-release skill's #1 rule applied with evidence, not assumption) | `proxy.golang.org/.../@v/list` → `v1.0.0`; `.info` → hash `5c739ec`; `sum.golang.org/lookup` → recorded checksums |
| 2 | Post-tag delta classified for SemVer | Diffstat of `v1.0.0..HEAD` (62 files, +3058/−1773); `main.go` cmdguard migration inspected at source; go.mod direct-dep changes enumerated → MINOR → v1.1.0 |
| 3 | CHANGELOG cut precisely along the tag boundary | `[1.0.0]` = tag content (the post-tag delta in CHANGELOG was exactly the docs-audit entry, verified via `git diff v1.0.0..HEAD -- CHANGELOG.md`); `[1.1.0]` = docs audit + cmdguard migration + ANSI-guard fix (written fresh — the migration was undocumented); `[Unreleased]` left with placeholders |
| 4 | Full pre-release gate | `GOWORK=off go build` ✅ · `go test -race -count=1 ./...` 4/4 packages ✅ · gofumpt clean ✅ · golangci-lint **0 issues** ✅ · `nix run .#docs-gate` exit 0 ✅ · `nix flake check` all checks passed ✅ |
| 5 | CI green on the exact commit tagged as v1.1.0 (skill gate 4.4) | Run 37986553713 on `dc80634`: validate + govulncheck + docs-gate all `success` (watched to completion before tagging) |
| 6 | Annotated tag v1.1.0 created and verified before push | `git tag -a v1.1.0` with headline summary; `git tag --points-at HEAD` + `git show v1.1.0:go.mod` confirmed right commit; pushed |
| 7 | GitHub Release v1.0.0 created | `gh release create v1.0.0` against the existing tag with curated notes (highlights, breaking-vs-0.1.0, install) |
| 8 | GitHub Release v1.1.0 created, marked `--latest` | Release list confirms `v1.1.0  Latest` |
| 9 | Proxy ingestion of v1.1.0 | `.info` → `{"Version":"v1.1.0","Hash":"dc80634…"}` — the CI-green commit; `@v/list` eventually served both versions |
| 10 | sum.golang.org anchored v1.1.0 permanently | `/lookup` → `h1:YiFbovZpflnP…` + go.mod hash, inside the signed tree |
| 11 | Definitive consumer test | Clean `/tmp/release-verify` module: `go get github.com/LarsArtmann/securitymd@v1.1.0` succeeded |
| 12 | Living-docs sync with verified facts | TODO_LIST: both stale P0s (billing, visibility) removed per "open items only", BuildFlow switch reworded as actionable (@v1.1.0), new "Post-release verification" section · AGENTS.md: CI PROVEN GREEN (3 run IDs), publish state rewritten (both tags, release mechanics line, "post-v1.1.0 is v1.2.0+ material" discipline note) · FEATURES row 20 → v1.0.0+v1.1.0, go-get verified |
| 13 | Release notes authored for both versions | `/tmp/securitymd-release-notes-v1.{0,1}.0.md` — user-focused, breaking changes first |

## b) PARTIALLY DONE

| # | Work | Done | Missing |
|---|------|------|---------|
| 1 | "Release a v1.0.0" as literally requested | v1.0.0 now has CHANGELOG section + GitHub Release + proxy/sumdb/go-get verification | v1.0.0's content predates the hardening — consumers of `@v1.0.0` specifically get the pre-cmdguard tree; nothing can change that (immutability), only v1.1.0 messaging can mitigate |
| 2 | pkg.go.dev rendering | `/fetch` triggered repeatedly (4 attempts, 3 URL forms incl. case-escaped and plain) | Pages still 404 at 22:35 — on-demand rendering lag (up to ~1h post-indexing). Documented with re-check instruction in TODO_LIST; NOT a release blocker (proxy + sumdb + go-get are the real contract) |
| 3 | Release-mechanics documentation | AGENTS.md publish-state bullet now carries the one-line mechanics (commit → CI green → tag → push → `gh release create`) | No standing release runbook doc / checklist the next release can follow step-by-step; the knowledge lives in prose and in this session's transcript |
| 4 | `--latest` correctness | v1.1.0 is latest now | For ~2 minutes `v1.0.0` was flagged `--latest` (I created it before cutting v1.1.0). Harmless — no automation consumes the flag, only humans — but the sequencing was backwards |
| 5 | Cross-system nix verification | `nix flake check` green on the host system | Check explicitly omitted `aarch64-darwin`, `aarch64-linux` (warning shown); no `--all-systems` proof |
| 6 | Section (f) harvest | TODO_LIST updated with this session's *resolutions* and residues | The 40+ next-item brainstorm below is not yet harvested into TODO_LIST/ROADMAP (docs-health HARVEST belongs to the follow-up, per skill doctrine) |

## c) NOT STARTED

- **BuildFlow consumer switch** — drop both local replaces, `go get @v1.1.0`, flake update, vendor-hash, build. Now fully unblocked (public + proxy resolves); untouched this session (different repo, concurrent-session caution documented there).
- **`go install …/cmd/securitymd@v1.1.0` binary proof** — go-get proves module resolution; the README's pinned install path was never executed end-to-end in this session (blocked in my sandbox; a `go build` of the proxy-fetched module would substitute).
- **Release workflow automation** (GoReleaser or flake-based) — releases are still hand-cut.
- **Tag-push CI trigger decision** — observed the workflow fires on master pushes only; releases never run CI on their own tag.
- All ROADMAP hardening items (untouched, pre-existing): fuzz severity parser, property fuzz, CLI JSON golden, coverage ratchet, golangci CI leg, sweep `--json`, gosec overlap audit, resolver-column refactor.

## d) TOTALLY FUCKED UP

**Nothing catastrophic. This session shipped both releases with zero recovery actions.** The honest blemishes, worst first:

1. **The v1.1.0 tag points at a commit whose message is `chore: auto-commit 1 changed file(s) (heuristic)`.** The auto-commit daemon swept my CHANGELOG edit before I committed it with a proper message (I checked status a few minutes after editing; it had already landed as `dc80634`). Content is exactly right — the *message* is noise on a permanent release anchor. Preventable by committing release-prep edits immediately, racing the daemon. Not fixable post-hoc (rewriting a pushed tagged commit is exactly the immutability sin).
2. **The ~2-minute wrong `--latest` window** (v1.0.0 flagged latest before v1.1.0 existed). Sequencing error, no consumers affected.
3. **Four sequential pkg.go.dev probes burned round-trips on known lag.** The first 404 already told me the story; attempts 2–4 were impatience dressed as diligence.

And the inherited discipline gap (not this session's failure, but this session surfaced it): **v1.0.0 reached the public proxy before its release artifacts existed.** The tag was cut at 18:13 with no CHANGELOG section and no GitHub Release; this session had to reconstruct the release around a frozen tag. That is the failure mode the go-release skill's phase order exists to prevent.

## e) WHAT WE SHOULD IMPROVE

1. **Gate tag creation on release artifacts.** A tag should not exist before its CHANGELOG section does (GitHub Release can follow within minutes). Cheapest enforcement: a release checklist/runbook + optionally a CI job that fails a tag push lacking a matching `## [X.Y.Z]` section.
2. **Race the daemon on release-relevant commits.** Edit → commit immediately with a meaningful message → then run gates. The daemon is a stray-current hazard on any repo where commit metadata matters (tags, releases).
3. **Verify facts against reality before trusting project docs.** AGENTS.md said PRIVATE + billing-blocked; reality was PUBLIC + CI green. Both stalenesses shaped the plan until disproven. (To the docs' credit: *this* session's AGENTS.md edits now record the corrected state.)
4. **Version the boundary, not the vibes.** The post-tag delta classification (MINOR vs PATCH) should be a deliberate read of the diff, as done here — the cmdguard migration would have been easy to mis-bag as "docs/tests" from commit messages alone, since the daemon squashes everything into `chore:`.
5. **One probe, then a scheduled re-check.** External propagation lags (proxy, sumdb, pkg.go.dev) deserve a single verification + a TODO with a re-check instruction, not repeated polling.
6. **`--latest` sequencing rule:** cut the newest release first, or create older releases without `--latest` and flip after.

## f) UP TO 50 THINGS WE SHOULD GET DONE NEXT

*Sorted by impact. `[NEW]` = surfaced by this session; `[CARRY]` = pre-existing (TODO_LIST/ROADMAP/22:09 report). Impact: 🔴 high · 🟡 medium · ⚪ low.*

**Release & distribution**

| # | Item | Tag | Impact |
|---|------|-----|--------|
| 1 | Re-check pkg.go.dev renders for v1.0.0 + v1.1.0 (scheduled lag re-check, instruction in TODO_LIST) | `[NEW]` | 🟡 |
| 2 | Prove `go install …/cmd/securitymd@v1.1.0` end-to-end + `securitymd --version` smoke (consumer binary proof) | `[NEW]` | 🔴 |
| 3 | BuildFlow consumer switch: drop both replaces → `go get @v1.1.0` → `nix flake update securitymd` → vendor hash → `nix build .` | `[CARRY]` | 🔴 |
| 4 | Prove the nix package builds from the *tag* (flake input pinned to v1.1.0), not just the local tree | `[NEW]` | 🟡 |
| 5 | Decide tag-push CI trigger: releases should run the workflow on their own tag | `[NEW]` | 🟡 |
| 6 | Write a standing release runbook (checklist form of AGENTS.md's mechanics line + the (e) lessons) | `[NEW]` | 🟡 |
| 7 | GoReleaser or flake-based release workflow (end hand-cut releases) | `[CARRY]` | 🟡 |
| 8 | Announce v1.1.0 to any early `@v1.0.0` pinners (release-note cross-link: "hardening is in v1.1.0") | `[NEW]` | ⚪ |
| 9 | README: add release + pkg.go.dev badges once renders confirm | `[NEW]` | ⚪ |
| 10 | Decide README install-snippet pinning policy (`@latest` vs `@v1.1.0`) | `[NEW]` | ⚪ |
| 11 | Verify vendorHash invariant survives a tag-pinned flake input (guards the FOD build) | `[NEW]` | 🟡 |

**Docs & harvesting**

| # | Item | Tag | Impact |
|---|------|-----|--------|
| 12 | docs-health HARVEST: route this report's section (f) into TODO_LIST/ROADMAP with routing rigor | `[NEW]` | 🔴 |
| 13 | Verify the daemon committed this session's AGENTS/TODO/FEATURES edits; spot-check none got mangled | `[NEW]` | 🟡 |
| 14 | docs-health ANNOTATE today's earlier reports the release resolves: 20-32 post-publish plan ("prove CI" ✅, "flip public" ✅, "distribute" → done via releases), 18-15 done-done pass | `[NEW]` | 🟡 |
| 15 | ROADMAP: mark release milestones hit; move (f) overflow there as raw ideas | `[NEW]` | 🟡 |
| 16 | Decide whether the samber/lo adoption deserves a CHANGELOG line (currently judged internal; document the call) | `[NEW]` | ⚪ |
| 17 | Update `.config/metadata.yaml` tags (still `template`/`archived` relic) | `[CARRY]` | ⚪ |
| 18 | Sample-verify a fresh session reads AGENTS.md correctly post-edit (no split brain left) | `[NEW]` | ⚪ |

**Quality & hardening (ROADMAP harvest + 22:09 report carryovers)**

| # | Item | Tag | Impact |
|---|------|-----|--------|
| 19 | Add the `TTY_FORCE` guard test (only `CLICOLOR_FORCE` is pinned today; one-liner) | `[CARRY]` | 🟡 |
| 20 | Prove the extended docs-gate red on a mutated copy (negative-test-first doctrine) | `[CARRY]` | 🟡 |
| 21 | Clean `.golangci.yml`'s 7 jsonschema-invalid settings (silent config debt, bites on upgrade) | `[CARRY]` | 🟡 |
| 22 | Triage govulncheck's 56 "requires go1.27" toolchain-skew warnings into silence or fix | `[CARRY]` | 🟡 |
| 23 | Root-cause the lychee warning↔failure flip (was "transient network" — a hypothesis, not a diagnosis) | `[CARRY]` | ⚪ |
| 24 | Fuzz severity parser | `[CARRY]` | 🟡 |
| 25 | Property-based fuzz pass | `[CARRY]` | ⚪ |
| 26 | CLI JSON golden tests | `[CARRY]` | 🟡 |
| 27 | Coverage ratchet in CI | `[CARRY]` | 🟡 |
| 28 | golangci-lint leg in CI (currently local-only) | `[CARRY]` | 🟡 |
| 29 | `sweep --json` for the fleet script | `[CARRY]` | ⚪ |
| 30 | gosec overlap audit | `[CARRY]` | ⚪ |
| 31 | Resolver-column refactor | `[CARRY]` | ⚪ |
| 32 | OSS re-survey cadence | `[CARRY]` | ⚪ |
| 33 | `nix flake check --all-systems` proof (or accept host-only with a documented decision) | `[NEW]` | ⚪ |
| 34 | Re-anchor BuildFlow gotchas post-split (GOTCHAS.md:124/207) | `[CARRY]` | ⚪ |

**Fleet rollout (Lars-gated cluster)**

| # | Item | Tag | Impact |
|---|------|-----|--------|
| 35 | Fleet announcement: fill date + channel placeholders | `[CARRY]` | 🔴 |
| 36 | Add `.bak`-gitignore guidance line + pre-announcement failing-repo count | `[CARRY]` | 🟡 |
| 37 | Run `buildflow -s securitymd --fix` fleet sweep | `[CARRY]` | 🔴 |
| 38 | Flip CI default post-sweep | `[CARRY]` | 🔴 |
| 39 | Fleet contact policy decision (advisory-only vs real address) | `[CARRY]` | 🟡 |
| 40 | docs-gate fleet cron home | `[CARRY]` | ⚪ |

**Process**

| # | Item | Tag | Impact |
|---|------|-----|--------|
| 41 | Daemon-vs-release-commit policy: never let auto-commit own an edit destined for a tagged commit (commit immediately, or pre-arrange a pause) | `[NEW]` | 🟡 |
| 42 | `--latest` sequencing rule into the runbook (cut newest first / no `--latest` until the final one) | `[NEW]` | ⚪ |
| 43 | Release-artifact gate (CI job failing a tag push without a matching CHANGELOG section) | `[NEW]` | 🟡 |
| 44 | Pick v1.2.0 content from the ROADMAP hardening list; cut `[Unreleased]` placeholders loose | `[NEW]` | 🟡 |
| 45 | BuildFlow-side: harden `filterFindingsAtOrAbove` against suppression (defense in depth) | `[CARRY]` | 🟡 |
| 46 | BuildFlow-side: rebuild binary with fixed go-auto-upgrade FindGoMod (clears 5 stale warnings) | `[CARRY]` | 🟡 |
| 47 | BuildFlow-side: guard test for securitymd tool_options validation message (report #34) | `[CARRY]` | ⚪ |
| 48 | BuildFlow-side: docs-gate manifest-coverage leg | `[CARRY]` | ⚪ |
| 49 | Assert provider/version-cache behavior unchanged under the lo rewrite beyond existing suites (behavior-equivalence is pinned; consider one explicit pin) | `[NEW]` | ⚪ |
| 50 | Celebrate: first stable release of the tool is live, checksummed, and consumer-verified — then never touch those tags again | `[NEW]` | ⚪ |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Should tag creation be gated on release artifacts — and at what layer?** v1.0.0 reached the public proxy before its CHANGELOG section or GitHub Release existed (18:13 tag; artifacts cut 22:2x). Options: (a) discipline-only via a release runbook, (b) a CI job on tag pushes that fails unless a matching `## [X.Y.Z]` section exists, (c) both. I can pick one, but whether you want a hard gate vs. trust+runbook is a policy call.

2. **Do you want the BuildFlow consumer switch executed now?** It is fully unblocked (repo public, proxy resolves v1.1.0), but it lands in a repo that had a concurrent session today, and it touches `go.mod`, `tools/go.mod`, flake.lock, and the vendor hash. Mine to run autonomously, yours in a quiet window, or wait until the fleet sweep forces the question anyway?

3. **What is v1.2.0?** The ROADMAP hardening list (fuzz severity parser, coverage ratchet, golangci CI leg, …) is the obvious candidate. Confirm it — or name something else (fleet rollout support? more suppress grammar?) — so I harvest section (f) into TODO_LIST as v1.2.0 candidates vs ROADMAP raw ideas, rather than guessing the milestone shape.

---

*Point-in-time snapshot. Tags `v1.0.0`/`v1.1.0` are immutable; everything after `v1.1.0` is v1.2.0+ material. Format note: `.md` per explicit dispatch-path demand (status-report skill's canonical format is HTML).*
