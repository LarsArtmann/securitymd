# Domain Language

Ubiquitous language for **securitymd**. Same terms mean the same thing in code, docs, and conversations.

## Glossary

| Term | Definition | Context |
| ---- | ---------- | ------- |
| Policy | A repository's `SECURITY.md` file — the artifact this tool owns | "Validate the policy", "generate a policy" |
| Candidate location | One of `SECURITY.md`, `.github/SECURITY.md`, `docs/SECURITY.md` — first existing wins | `pkg/policy.CandidateLocations` |
| Rule | A single validation check with a stable kebab ID (`missing-file`, `missing-header`, `too-short`, …) | Findings, suppressions, and configs key on rule IDs |
| Finding | A `finding.Finding` emitted by detect/validate — the only output currency of the tool | go-finding data model |
| Missing-file finding | The error-severity finding for an absent policy; carries a rendered preview as `AfterCode` when identity is derivable | Drives BuildFlow `--fix` |
| Repair | Creating a missing policy via `Generate` — create-only, never overwrites an existing (human-written) policy | toolsdk Repairer |
| Generate | Rendering the embedded template with identity + optional contact email | `pkg/policy.Generate` |
| Identity | `RepoIdentity{Organization, Repository}` parsed from `git remote origin` | `pkg/policy.DetectRepoIdentity` |
| Contact | Where vulnerability reports go: GitHub advisory link always, email only when provided | No fabricated `security@domain` addresses — deliberate |
| Version cell | The "Supported Versions" table row: latest git tag, or "Latest release" | Best-effort via `git describe` |
| Dry run | Render and report without writing | CLI `--dry-run`, toolsdk dry-run context |
| Provider | The toolsdk self-registration (`securitymd`) that BuildFlow consumes via blank import | `pkg/provider.Provider` |

## Bounded contexts

| Context | Description |
| ------- | ----------- |
| `pkg/policy` | Domain core: validation rules, generation, detection, git identity |
| `pkg/provider` | Integration: toolsdk Spec exposing the core to BuildFlow |
| `cmd/securitymd` | CLI surface over the same core (validate, setup, status) |
