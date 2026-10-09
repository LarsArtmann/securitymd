# TODO List — securitymd

_Open items only. Source: harvest of `docs/status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md` (2026-10-09), minus items completed since._

## Blocked on Lars (needs human action)

- [ ] **P0: Publish securitymd** — rename local dir + GitHub repo `template-SECURITY` → `securitymd`, push, tag `v1.0.0`; then drop BuildFlow replaces + flake input, re-vendor, `nix run .#update-vendor-hash`; regenerate this repo's own SECURITY.md (current one predates the rename and would render the old repo name). Input: `docs/archive/pre-rebuild/PUBLIC_OR_PRIVATE.md` (verdict: conditionally make public — cleanup done by the rebuild).
- [ ] **Fleet contact policy decision** — keep GitHub-advisory-only default (current) or bake a real address (e.g. `security@lars.software`) via BuildFlow `tool_options`. Default stays advisory-only until Lars rules (no fabricated addresses).
- [ ] **Fleet gate announcement** — `missing-file` is error severity; every manifest-carrying repo without SECURITY.md fails BuildFlow's findings gate. Announce + run one `buildflow -s securitymd --fix` sweep before it lands in CI defaults.
- [ ] **Update `.config/metadata.yaml` tags** — still says `template`/`archived` (pre-rebuild relic; also flagged by PUBLIC_OR_PRIVATE's should-fix list). Decide: fresh tags reflecting the live securitymd tool, or confirm archiving. Consumer of the file is external tooling — Lars's call. (Examined 2026-10-09: `git-town.toml` is fine — `main = "master"`, GitHub API connector.)

## This repo

- [ ] SARIF golden-file test for CLI output (`cmd/securitymd/validate.go`; report #17)
- [ ] Mutation/discrimination proof run for the rule table — break one rule, watch tests fail (report #18)
- [ ] Decide + wire `README.md` into trigger manifests so docs-only repos activate (report #16)
- [ ] `setup` interactive mode when no git remote — prompt for org/repo (report #20)
- [ ] `--force`/`--regenerate` for setup with explicit backup (report #21)
- [ ] Canonical policy-location knob (`docs/SECURITY.md`-preferring repos; report #22)
- [ ] Per-rule severity configuration (report #23)
- [ ] Suppression comments (`securitymd:ignore(rule) reason`; report #26)
- [ ] Cache the version cell (avoid re-running `git describe` per detect; report #28)
- [ ] Fuzz `parseGitRemote` over remote-URL shapes (report #29)
- [ ] flake.nix `packages.<system>.securitymd` binary output (plan 2026-10-09 M14)
- [ ] OSS-landscape survey before publishing — verify the "no dedicated SECURITY.md validator" positioning claim, update README/ROADMAP if contradicted (plan M4)
- [ ] One-command docs-health gate for this repo (per-file `~~` + check-rows; plan M10)

## Process & ecosystem (plan 2026-10-09 M22–M24)

- [ ] Upstream feedback batch: exhaustruct anchored-patterns note (golangci-lint-auto-configure), `SaveBytes` proposal (linter-autoconfigure-sdk), crush-config lessons.md "nix FOD ignores local replaces", nix-private-go-repos gotools/goimports gotcha (M22)
- [ ] Scope verdicts recorded in ROADMAP: baseline/ratchet, localization, golangci-plugin distribution — decide defer/adopt (M23)
- [ ] docs-health fleet cron/gate (un-annotated archives + TODO/CHANGELOG overlap drift alarm) + `docs/reviews/` convention decision (M24; needs Lars)

## BuildFlow-side

- [ ] Once the concurrent `execution/` file-split sweep in BuildFlow lands (2026-10-09, still running at pipeline_1300+): verify `nix build .` green and re-run `buildflow -s securitymd --fix` e2e with the nix-built binary. The securitymd FOD resolution itself is already proven (deps phase passes; vendorHash invariant, gotcha #230).
- [ ] BuildFlow gotcha-anchor warnings (GOTCHAS.md:124/207 → `execution/pipeline.go:432/530`) need re-anchoring to the post-split files — belongs to the split sweep's follow-up.
- [ ] Guard test: assert `securitymd` tool_options validation error message text (report #34)
- [ ] Run `docs --check` after any provider-count change and fix table drift (report #32; clean at 0 fail as of 2026-10-09)
