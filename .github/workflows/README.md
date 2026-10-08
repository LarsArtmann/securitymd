# GitHub Security Policy Validation

This workflow builds and tests securitymd, then validates this repository's own `SECURITY.md` with the tool itself (dogfooding).

## What it does

- Builds the module and runs the full test suite (`go build ./...`, `go test ./...`)
- Runs `securitymd validate --file SECURITY.md` — fails the build on error-severity findings
- Go version follows `go.mod` (`go-version-file`)

## Run locally

```bash
go run ./cmd/securitymd validate --file SECURITY.md
```

Or via BuildFlow: `buildflow -s securitymd`
