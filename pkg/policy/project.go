package policy

import (
	"os/exec"
	"strings"
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
func DetectRepoIdentity(dir string) RepoIdentity {
	// #nosec G204 -- fixed argv; dir is the caller's repository path
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")

	output, err := cmd.Output()
	if err != nil {
		return RepoIdentity{}
	}

	return parseGitRemote(strings.TrimSpace(string(output)))
}

func parseGitRemote(remote string) RepoIdentity {
	remote = strings.TrimSuffix(strings.TrimSpace(remote), ".git")

	remote = strings.TrimPrefix(remote, "https://")
	remote = strings.TrimPrefix(remote, "ssh://")
	remote = strings.TrimPrefix(remote, "git@")

	// git@github.com:ORG/REPO -> github.com/ORG/REPO
	remote = strings.Replace(remote, ":", "/", 1)

	parts := strings.Split(remote, "/")
	if len(parts) < 2 {
		return RepoIdentity{}
	}

	org, repo := parts[len(parts)-2], parts[len(parts)-1]
	if org == "" || repo == "" {
		return RepoIdentity{}
	}

	return RepoIdentity{Organization: org, Repository: repo}
}

// LatestTag returns the most recent git tag of the directory (e.g. "v1.2.3"),
// or "" when no tag exists. Best-effort: failures degrade to the empty string.
func LatestTag(dir string) string {
	// #nosec G204 -- fixed argv; dir is the caller's repository path
	cmd := exec.Command("git", "-C", dir, "describe", "--tags", "--abbrev=0")

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
