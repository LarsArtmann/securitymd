package internal

import (
	"os/exec"
	"strings"
)

// GitHubIntegration provides GitHub-specific integrations
type GitHubIntegration struct {
	detector *ProjectDetector
}

// NewGitHubIntegration creates a new GitHub integration
func NewGitHubIntegration() *GitHubIntegration {
	return &GitHubIntegration{
		detector: NewProjectDetector(),
	}
}

// GitHubInfo represents GitHub repository information
type GitHubInfo struct {
	Owner         string
	Repository    string
	FullName      string
	URL           string
	IsGitHub      bool
	Host          string
	DefaultBranch string
	IsPrivate     bool
}

// GetGitHubInfo detects GitHub repository information
func (gi *GitHubIntegration) GetGitHubInfo() *GitHubInfo {
	info := &GitHubInfo{
		IsGitHub: false,
		Host:     "github.com",
	}

	// Get git remote URL
	remoteURL := gi.detector.detectFromGitRemote()
	if remoteURL == "" {
		return info
	}

	// Parse GitHub URL
	info = gi.parseGitHubURL(remoteURL)

	// Get additional info from git commands
	if info.IsGitHub {
		info.DefaultBranch = gi.getDefaultBranch()
		info.IsPrivate = gi.checkIfPrivate()
	}

	return info
}

// parseGitHubURL parses various GitHub URL formats
func (gi *GitHubIntegration) parseGitHubURL(url string) *GitHubInfo {
	info := &GitHubInfo{
		IsGitHub: false,
	}

	cleanURL := strings.TrimSpace(url)
	cleanURL = strings.TrimSuffix(cleanURL, ".git")

	// HTTPS format: https://github.com/owner/repo
	if strings.HasPrefix(cleanURL, "https://github.com/") {
		parts := strings.Split(strings.TrimPrefix(cleanURL, "https://github.com/"), "/")
		if len(parts) >= 2 {
			info.Owner = parts[0]
			info.Repository = parts[1]
			info.FullName = parts[0] + "/" + parts[1]
			info.URL = "https://github.com/" + info.FullName
			info.IsGitHub = true
			info.Host = "github.com"
		}
		return info
	}

	// SSH format: git@github.com:owner/repo
	if strings.HasPrefix(cleanURL, "git@github.com:") {
		parts := strings.Split(strings.TrimPrefix(cleanURL, "git@github.com:"), "/")
		if len(parts) >= 2 {
			info.Owner = parts[0]
			info.Repository = parts[1]
			info.FullName = parts[0] + "/" + parts[1]
			info.URL = "https://github.com/" + info.FullName
			info.IsGitHub = true
			info.Host = "github.com"
		}
		return info
	}

	// Git protocol: git://github.com/owner/repo
	if strings.HasPrefix(cleanURL, "git://github.com/") {
		parts := strings.Split(strings.TrimPrefix(cleanURL, "git://github.com/"), "/")
		if len(parts) >= 2 {
			info.Owner = parts[0]
			info.Repository = parts[1]
			info.FullName = parts[0] + "/" + parts[1]
			info.URL = "https://github.com/" + info.FullName
			info.IsGitHub = true
			info.Host = "github.com"
		}
		return info
	}

	// Generic GitHub Enterprise: https://github.enterprise.com/owner/repo
	if strings.Contains(cleanURL, "github") {
		if strings.HasPrefix(cleanURL, "https://") {
			hostAndPath := strings.TrimPrefix(cleanURL, "https://")
			parts := strings.Split(hostAndPath, "/")
			if len(parts) >= 3 {
				info.Host = parts[0]
				info.Owner = parts[1]
				info.Repository = parts[2]
				info.FullName = info.Owner + "/" + info.Repository
				info.URL = cleanURL
				info.IsGitHub = true
			}
		}
	}

	return info
}

// getDefaultBranch gets the default branch from git
func (gi *GitHubIntegration) getDefaultBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to common branch names
		return gi.detectCurrentBranch()
	}

	ref := strings.TrimSpace(string(output))
	if strings.HasPrefix(ref, "refs/remotes/origin/") {
		return strings.TrimPrefix(ref, "refs/remotes/origin/")
	}

	return gi.detectCurrentBranch()
}

// detectCurrentBranch gets the current branch
func (gi *GitHubIntegration) detectCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "main" // Default fallback
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "main"
	}

	return branch
}

// checkIfPrivate attempts to determine if repository is private
func (gi *GitHubIntegration) checkIfPrivate() bool {
	// This is a simplified check
	// In a real implementation, you might use GitHub API
	// For now, assume false (public)
	return false
}

// GetGitHubURLs returns various GitHub URLs for the repository
func (gi *GitHubIntegration) GetGitHubURLs(info *GitHubInfo) map[string]string {
	if !info.IsGitHub {
		return map[string]string{}
	}

	baseURL := "https://" + info.Host + "/" + info.FullName

	return map[string]string{
		"repository": baseURL,
		"issues":     baseURL + "/issues",
		"security":   baseURL + "/security",
		"advisories": baseURL + "/security/advisories",
		"pulse":      baseURL + "/pulse",
		"network":    baseURL + "/network",
		"settings":   baseURL + "/settings",
		"actions":    baseURL + "/actions",
		"releases":   baseURL + "/releases",
		"wiki":       baseURL + "/wiki",
		"api":        "https://api." + info.Host + "/repos/" + info.FullName,
	}
}

// GenerateGitHubTemplateVariables generates template variables from GitHub info
func (gi *GitHubIntegration) GenerateGitHubTemplateVariables(info *GitHubInfo) map[string]string {
	if !info.IsGitHub {
		return map[string]string{}
	}

	urls := gi.GetGitHubURLs(info)

	variables := map[string]string{
		// Basic repo info
		"{{GITHUB_OWNER}}":      info.Owner,
		"{{GITHUB_REPOSITORY}}": info.Repository,
		"{{GITHUB_FULL_NAME}}":  info.FullName,
		"{{GITHUB_URL}}":        info.URL,
		"{{GITHUB_HOST}}":       info.Host,
		"{{GITHUB_BRANCH}}":     info.DefaultBranch,

		// URLs
		"{{GITHUB_ISSUES_URL}}":     urls["issues"],
		"{{GITHUB_SECURITY_URL}}":   urls["security"],
		"{{GITHUB_ADVISORIES_URL}}": urls["advisories"],
		"{{GITHUB_SETTINGS_URL}}":   urls["settings"],
		"{{GITHUB_ACTIONS_URL}}":    urls["actions"],
		"{{GITHUB_RELEASES_URL}}":   urls["releases"],

		// Template-specific URLs
		"{{SECURITY_ADVISORIES_URL}}": urls["advisories"],
		"{{BOUNTY_PROGRAM_URL}}":      urls["security"] + "/policy",
		"{{SECURITY_BLOG_URL}}":       "https://" + info.Host + "/" + info.Owner + "/blog/security",

		// Organization info
		"{{ORGANIZATION}}": info.Owner,
		"{{DOMAIN}}":       info.Host,
	}

	// Generate domain from org for email
	if info.Owner != "" {
		variables["{{ORGANIZATION_EMAIL}}"] = "security@" + strings.ToLower(info.Owner) + ".com"
		variables["{{CONTACT_EMAIL}}"] = "security@" + strings.ToLower(info.Owner) + ".com"
	}

	return variables
}

// IsGitHubRepository checks if current directory is a GitHub repository
func (gi *GitHubIntegration) IsGitHubRepository() bool {
	info := gi.GetGitHubInfo()
	return info.IsGitHub
}
