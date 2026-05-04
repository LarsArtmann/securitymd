# Public or Private? — Decision Analysis

> **Project:** template-SECURITY
> **Date:** 2026-05-04
> **Current Visibility:** Private
> **Recommendation:** **Make public, conditionally**

---

## Project Summary

A Go CLI tool (`template-security`) that validates and generates `SECURITY.md` files for GitHub repositories. It checks for required sections (vulnerability reporting, supported versions, security practices, contact info), validates content quality, and generates compliant templates.

- **Language:** Go 1.26
- **License:** MIT
- **Lines of Code:** ~1,200 (Go) + ~1,300 (shell scripts)
- **Dependencies:** cobra, viper, ginkgo/gomega, testify, fatih/color, go-composable-business-types
- **CI/CD:** GitHub Actions workflow for security validation
- **Tests:** Unit tests (testify) + BDD acceptance tests (ginkgo) — structure exists, coverage is incomplete

---

## Pro Publica (Arguments FOR Making Public)

### 1. Zero Security Risk — No Secrets, No Vulnerable Logic

This tool **validates markdown files**. It reads `SECURITY.md`, checks for sections, and generates templates. There is:
- No network server, no API endpoints, no database
- No secrets, credentials, or private keys anywhere in the codebase
- No proprietary algorithms or business logic
- No user data processing

**Verdict:** No security downside to exposing this code.

### 2. Solves a Real, Unsolved Problem

As noted in `PARTS.md`: **no dedicated SECURITY.md validator exists** in the open-source ecosystem. Markdown linters are generic; GitHub's guidance is documentation-only. This tool fills a genuine gap:
- Security-specific validation rules
- CI/CD integration ready
- Opinionated standards based on GitHub's official guidelines

### 3. Strong Reusability Potential

- **`ProjectDetector`** (polyglot project metadata detection from git/package.json/go.mod/Cargo.toml/pyproject.toml) — already analyzed as extraction-worthy in `PARTS.md`
- **`SecurityValidator`** — could become a GitHub Action or pre-commit hook
- **Template engine** — reusable for any policy-as-code workflow

### 4. Community Benefit

Making this public would:
- Help other projects improve their security disclosure practices
- Serve as a reference implementation for SECURITY.md best practices
- Enable contributions (bug fixes, new validation rules, additional templates)
- Potentially be listed on GitHub's security tools ecosystem page

### 5. Portfolio & Reputation Value

- Demonstrates Go CLI engineering (cobra, viper, structured errors)
- Shows testing discipline (BDD with ginkgo, table-driven tests)
- Signals care for developer security tooling
- MIT license already chosen — the intent to share is there

### 6. Low Maintenance Overhead Once Public

- No infrastructure to maintain (it's a CLI binary)
- No SaaS costs, no uptime requirements
- Issues can be community-triaged
- `go install github.com/LarsArtmann/template-SECURITY/cmd/template-security@latest` just works

---

## Contra Publica (Arguments AGAINST Making Public)

### 1. Code Quality Is Not Production-Ready

| Issue | Detail |
|-------|--------|
| `go.mod` has `replace` directive | Points to local `../go-composable-business-types` — won't compile for external users without fixing |
| `go mod tidy` required | Dependency graph is stale; LSP reports 20+ errors |
| BDD test suite is skeleton-only | `BDD_TESTS_REVIEW.md` scores it 0/50 — acceptance tests exist but are incomplete |
| Shell scripts are legacy | 1,300 lines of bash that duplicate Go functionality; no clear migration path documented |
| Workflow uses outdated Go 1.21 | `go.mod` says 1.26 but CI uses `setup-go@v4` with Go 1.21 |
| Config uses placeholder values | `.template-security.yaml` has `MyCompany`, `security@mycompany.com`, `hackerone.com` URLs |
| SECURITY.md itself has `security@github.com` | Not a real contact for this project |

### 2. Unfinished Architecture Refactoring

- `IMPROVEMENT_PLAN.md` lists 21 items, most uncompleted
- `PARTS.md` proposes extraction of `ProjectDetector` into `projectmeta` — not done
- `PROJECT_SPLIT_EXECUTIVE_REPORT.md` proposes 4-project split — not done
- The codebase is in a mid-refactoring state with stale documentation

### 3. Documentation Is Internally Inconsistent

- `CHANGELOG.md` says v0.1.0 from 2026-01-01 but no tags or releases exist
- `BDD_TESTS_REVIEW.md` is a self-critique document that paints a poor picture
- Multiple status reports in `docs/status/` describe various incomplete migration states
- `metadata.yaml` tags the project as `archived` with importance `25/100`

### 4. No Release Artifacts

- No GitHub releases, no tagged versions
- No goreleaser configuration
- No binary distribution (only `go install` which won't work due to `replace` directive)

### 5. Potential Maintenance Burden

- Opening to community means issue triage, PR reviews, support expectations
- The `go-composable-business-types` dependency is a sibling project — coupling risk
- Shell scripts would need to be either maintained or removed

---

## Conditional Recommendation

### Make public IF the following conditions are met:

#### Must-Fix Before Public (estimated 2-3 hours)

- [ ] **Fix `go.mod` `replace` directive** — Either publish `go-composable-business-types` first, or remove the dependency and inline the minimal types used (`IDID`, `PolicyID`)
- [ ] **Run `go mod tidy`** — Ensure clean dependency graph
- [ ] **Fix CI workflow Go version** — Align `setup-go` with `go.mod` version (1.26)
- [ ] **Replace placeholder config values** — Update `.template-security.yaml` and `SECURITY.md` with real values or clearly mark as examples
- [ ] **Remove or archive stale internal docs** — `BDD_TESTS_REVIEW.md`, `IMPROVEMENT_PLAN.md`, `PROJECT_SPLIT_EXECUTIVE_REPORT.md`, `docs/status/*` are internal working documents that create a negative impression

#### Should-Fix Before Public (estimated 1-2 hours)

- [ ] **Add a `.goreleaser.yml`** or at minimum document the `go install` path works
- [ ] **Update README.md** — Remove emojis (per project convention), add contribution guidelines, clarify the project status
- [ ] **Remove `metadata.yaml` `archived` tag** — Or confirm the project IS archived and decide accordingly
- [ ] **Tag v0.1.0** — Create the first release to signal intent

#### Nice-to-Have Post-Public

- [ ] Extract `ProjectDetector` into `projectmeta` (per `PARTS.md` plan)
- [ ] Replace `viper` with `koanf` (per architecture guidelines)
- [ ] Complete BDD test suite
- [ ] Publish as GitHub Action
- [ ] Clean up legacy shell scripts

---

## Decision Matrix

| Factor                    | Weight | Public | Private | Notes                          |
|---------------------------|--------|--------|---------|--------------------------------|
| Security risk             | High   | 10     | 10      | No difference — zero risk      |
| Community value           | High   | 9      | 2       | Fills a real gap               |
| Code quality perception   | High   | 4      | 7       | Needs cleanup first            |
| Maintenance burden        | Medium | 4      | 8       | Private = zero expectations    |
| Portfolio/reputation      | Medium | 8      | 3       | Shows engineering standards    |
| Reusability               | Medium | 9      | 3       | Others can build on it         |
| Dependency readiness      | High   | 3      | 8       | `replace` directive blocks use |
| **Weighted Total**        |        | **47/70** | **41/70** | **Public wins, marginally** |

---

## Final Verdict

**Conditionally make public.** The project solves a genuine unsolved problem in the OSS ecosystem and carries zero security risk. However, publishing in its current state would create a poor first impression due to stale dependencies, placeholder values, and unfinished refactoring.

**Recommended path:**

1. Spend **2-3 hours** on the must-fix items above
2. Publish
3. Announce in relevant communities (Go security tooling, GitHub security guides)
4. Iterate based on community feedback

If those 2-3 hours cannot be allocated in the near term, keeping private is acceptable — but the project should be tagged as `public-when-ready` to ensure the intent is not lost.
