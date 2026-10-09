# Publish Runbook — securitymd v1.0.0

**Date:** 2026-10-09 16:34 CEST
**Scope:** Everything mechanical about publishing `github.com/LarsArtmann/securitymd`, staged and pre-verified so the human steps are four commands and the robot steps are copy-paste. Input: plan 2026-10-09 M1; consumer wiring verified live in BuildFlow (go.mod:425 replace, flake.nix:122 input).
**Rule:** every command below was either executed while writing this runbook (marked ✅) or is gated on the rename (marked ⏳).

## 0. Pre-verified readiness (✅ done 2026-10-09)

| Check                               | Command                                                                                | Result                                                                  |
| ----------------------------------- | -------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| Module path is the new name         | `head -1 go.mod`                                                                       | `module github.com/LarsArtmann/securitymd`                              |
| `go mod tidy` clean                 | `GOWORK=off go mod tidy && git status --short go.mod go.sum`                           | no diff                                                                 |
| Local install path compiles + links | `GOWORK=off go build -o /tmp/securitymd ./cmd/securitymd && /tmp/securitymd --version` | `securitymd version dev`                                                |
| Nix package builds sandboxed        | `nix build .#securitymd` (plan M14)                                                    | green, e2e smoke passed                                                 |
| Gates green                         | `GOWORK=off go test ./... && GOWORK=off golangci-lint run && nix flake check`          | green (suppression-suite caveat resolved once the concurrent WIP lands) |
| Output contract pinned              | `go test ./pkg/policy -run TestReport_golden` (plan M2)                                | green                                                                   |

## 1. Lars: the four human steps

```bash
# 1. Rename the local directory (fleet workspace may need re-registration)
mv /home/lars/projects/template-SECURITY /home/lars/projects/securitymd

# 2. Rename the GitHub repo (Settings → General → Repository name):
#    template-SECURITY → securitymd
#    GitHub redirects the old name automatically; module consumers are unaffected.

# 3. Push master + tag (go-release skill governs the full release pass;
#    the tag is what makes @latest resolvable):
git push origin master
git tag -a v1.0.0 -m "securitymd v1.0.0 — first public release" && git push origin v1.0.0

# 4. Verify proxy + pkg.go.dev propagation (may take a few minutes):
GOPROXY=https://proxy.golang.org GOWORK=off go list -m github.com/LarsArtmann/securitymd@v1.0.0
```

Then: `go install github.com/LarsArtmann/securitymd/cmd/securitymd@latest` from a scratch dir outside the fleet workspace.

## 2. BuildFlow consumer: drop the scaffolding (⏳ after rename)

```bash
cd /home/lars/projects/BuildFlow

# a. Drop the local-path replace (go.mod line ~425):
#    replace github.com/LarsArtmann/securitymd => /home/lars/projects/template-SECURITY
GOWORK=off go mod edit -dropreplace=github.com/LarsArtmann/securitymd

# b. Point the require at the published version:
GOWORK=off go get github.com/LarsArtmann/securitymd@v1.0.0

# c. Re-vendor the workspace (BuildFlow AGENTS: go work vendor, NOT go mod vendor):
go work vendor

# d. Drop the flake input (flake.nix lines ~122-132, the git+ssh
#    template-SECURITY input with the pending-rename comment), then:
nix flake update   # or scoped update of consumers of that input

# e. Refresh the vendored hash:
nix run .#update-vendor-hash

# f. Verify:
nix build . && GOWORK=off go test ./...
```

## 3. This repo: dogfood regeneration (⏳ after rename)

```bash
cd /home/lars/projects/securitymd   # new path

# Current SECURITY.md predates the rename and renders the old repo name.
trash SECURITY.md && securitymd setup   # fresh template, advisory-only contact
securitymd validate                     # must exit 0

# metadata.yaml still says template/archived — Lars decides fresh tags (TODO_LIST P0).
```

## 4. Post-publish checklist

- [ ] `pkg.go.dev/github.com/LarsArtmann/securitymd` renders (public README)
- [ ] ROADMAP's "publish as open source" direction marked done
- [ ] FEATURES row 13 (Published module) → FULLY_FUNCTIONAL
- [ ] TODO_LIST P0 publish item checked off; M11 fleet sweep + M19 website unblock
- [ ] GitHub repo description + topics (`security`, `sarif`, `security-policy`, `linter`)
- [ ] Old repo redirect works: `git ls-remote git@github.com:LarsArtmann/template-SECURITY.git`
