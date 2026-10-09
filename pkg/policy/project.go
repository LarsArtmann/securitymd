package policy

import (
	"context"
	"os/exec"
	"strings"
	"sync"
)

// RepoIdentity names the GitHub organization/owner and repository a policy
// file belongs to. Both parts are required to render truthful advisory links.
type RepoIdentity struct {
	Organization string
	Repository   string
}

// IsComplete reports whether both identity parts are known.
func (r RepoIdentity) IsComplete() bool {
	return r.Organization != "" && r.Repository != ""
}

// DetectRepoIdentity derives the organization/repository from the directory's
// git origin remote. Supports HTTPS and SSH GitHub URLs; returns an
// incomplete identity when the remote is missing or unparseable.
func DetectRepoIdentity(ctx context.Context, dir string) RepoIdentity {
	// #nosec G204 -- fixed argv; dir is the caller's repository path
	cmd := exec.CommandContext(ctx, "git", "-C", dir, "config", "--get", "remote.origin.url")

	output, err := cmd.Output()
	if err != nil {
		return RepoIdentity{}
	}

	return parseGitRemote(strings.TrimSpace(string(output)))
}

// minRemoteSegments is the shortest hosted-remote shape: host, organization,
// repository (e.g. github.com/ORG/REPO after scheme stripping).
const minRemoteSegments = 3

func parseGitRemote(remote string) RepoIdentity {
	remote = strings.TrimSpace(remote)
	remote = strings.TrimSuffix(remote, "/")
	remote = strings.TrimSuffix(remote, ".git")
	remote = strings.TrimSuffix(remote, "/")

	// Only the documented remote forms carry GitHub coordinates. Anything
	// else (git://, file://, local paths) degrades to an incomplete identity
	// rather than a plausible-looking wrong one.
	if !strings.HasPrefix(remote, "https://") &&
		!strings.HasPrefix(remote, "ssh://") &&
		!strings.HasPrefix(remote, "git@") {
		return RepoIdentity{}
	}

	remote = strings.TrimPrefix(remote, "https://")
	remote = strings.TrimPrefix(remote, "ssh://")
	remote = strings.TrimPrefix(remote, "git@")

	// git@github.com:ORG/REPO -> github.com/ORG/REPO
	remote = strings.Replace(remote, ":", "/", 1)

	parts := strings.Split(remote, "/")
	if len(parts) < minRemoteSegments {
		return RepoIdentity{}
	}

	org, repo := parts[len(parts)-2], parts[len(parts)-1]
	if org == "" || repo == "" {
		return RepoIdentity{}
	}

	// A complete identity must be clean coordinates: stray separators or
	// whitespace would render a garbage advisory link into the policy.
	if strings.ContainsAny(org, ":/ \t\r\n") || strings.ContainsAny(repo, ":/ \t\r\n") {
		return RepoIdentity{}
	}

	return RepoIdentity{Organization: org, Repository: repo}
}

// LatestTag returns the most recent git tag of the directory (e.g. "v1.2.3"),
// or "" when no tag exists. Best-effort: failures degrade to the empty string.
func LatestTag(ctx context.Context, dir string) string {
	// #nosec G204 -- fixed argv; dir is the caller's repository path
	cmd := exec.CommandContext(ctx, "git", "-C", dir, "describe", "--tags", "--abbrev=0")

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// latestTagCache memoizes the latest tag per directory for the process
// lifetime: Detect (missing-file preview) and Repair (Generate) render the
// same policy within one run, and `git describe` is not free. Tags do not
// change mid-run; a fresh process re-reads them.
var latestTagCache sync.Map //nolint:gochecknoglobals // per-process memo table keyed by repository directory

// lookupLatestTag is LatestTag with a per-directory process cache. The empty
// result (no tags) is cached too — re-probing a tagless repo on every render
// is exactly the repeated cost this eliminates.
func lookupLatestTag(ctx context.Context, dir string) string {
	if cached, ok := latestTagCache.Load(dir); ok {
		tag, _ := cached.(string)

		return tag
	}

	tag := LatestTag(ctx, dir)
	latestTagCache.Store(dir, tag)

	return tag
}
