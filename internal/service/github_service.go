package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
)

// gitHubService implements GitHubService interface
type gitHubService struct {
	gitService   GitService
	errorHandler ErrorHandler
}

// NewGitHubService creates a new GitHub service
func NewGitHubService() GitHubService {
	errorHandler := NewErrorHandler(true)
	gitService := NewGitService()
	
	return &gitHubService{
		gitService:   gitService,
		errorHandler: errorHandler,
	}
}

// GetGitHubInfo detects GitHub repository information
func (g *gitHubService) GetGitHubInfo(ctx context.Context) (domain.GitHubInfo, error) {
	if !g.IsGitHubRepository(ctx) {
		return domain.GitHubInfo{
			IsGitHub: false,
			Host:     "github.com",
		}, nil
	}

	remoteURL, err := g.gitService.GetRemoteURL(ctx)
	if err != nil {
		return domain.GitHubInfo{}, errors.NewGitHubError(
			"GH001",
			"Failed to get git remote URL",
			"",
			"GetRemoteURL",
		)
	}

	info := g.parseGitHubURL(remoteURL)
	
	// Get additional info
	info.DefaultBranch = g.getDefaultBranch(ctx)
	info.Status = domain.StatusPublic // Default assumption
	info.Protocol = g.detectProtocol(remoteURL)
	info.URLs = g.GetGitHubURLs(info)

	return info, nil
}

// GetGitHubURLs generates GitHub URLs for repository
func (g *gitHubService) GetGitHubURLs(info domain.GitHubInfo) map[string]string {
	if !info.IsGitHub {
		return map[string]string{}
	}

	baseURL := "https://" + info.Host + "/" + info.FullName
	
	return map[string]string{
		"repository":  baseURL,
		"issues":      baseURL + "/issues",
		"security":    baseURL + "/security",
		"advisories":  baseURL + "/security/advisories",
		"pulse":       baseURL + "/pulse",
		"network":     baseURL + "/network",
		"settings":    baseURL + "/settings",
		"actions":     baseURL + "/actions",
		"releases":    baseURL + "/releases",
		"wiki":        baseURL + "/wiki",
		"api":         "https://api." + info.Host + "/repos/" + info.FullName,
	}
}

// IsGitHubRepository checks if repository is on GitHub
func (g *gitHubService) IsGitHubRepository(ctx context.Context) bool {
	remoteURL, err := g.gitService.GetRemoteURL(ctx)
	if err != nil {
		return false
	}

	return g.isGitHubURL(remoteURL)
}

// GetOwnerInfo retrieves GitHub owner information
func (g *gitHubService) GetOwnerInfo(ctx context.Context, owner string) (domain.Organization, error) {
	// For now, return minimal organization info
	return domain.Organization{
		Name:  owner,
		Type:  domain.OrgTypeOpenSource,
		Size:  domain.SizeSmall,
		GitHub: domain.GitHubInfo{
			Owner:    owner,
			Host:     "github.com",
			IsGitHub: true,
		},
	}, nil
}

// GetSecurityAdvisories retrieves security advisories
func (g *gitHubService) GetSecurityAdvisories(ctx context.Context, owner, repo string) ([]SecurityAdvisory, error) {
	// For now, return empty slice
	// In a real implementation, this would call GitHub API
	return []SecurityAdvisory{}, nil
}

// Private helper methods

func (g *gitHubService) parseGitHubURL(url string) domain.GitHubInfo {
	info := domain.GitHubInfo{
		IsGitHub: false,
		Host:     "github.com",
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

func (g *gitHubService) isGitHubURL(url string) bool {
	return strings.Contains(url, "github.com") ||
		   strings.HasPrefix(url, "git@github.com:") ||
		   strings.HasPrefix(url, "git://github.com/")
}

func (g *gitHubService) detectProtocol(url string) types.GitProtocol {
	if strings.HasPrefix(url, "https://") {
		return types.ProtocolHTTPS
	}
	if strings.HasPrefix(url, "git@") {
		return types.ProtocolSSH
	}
	if strings.HasPrefix(url, "git://") {
		return types.ProtocolGit
	}
	return types.ProtocolHTTPS
}

func (g *gitHubService) getDefaultBranch(ctx context.Context) string {
	// Try to get default branch from git
	if g.gitService.IsGitRepository(ctx) {
		branch, err := g.gitService.GetCurrentBranch(ctx)
		if err == nil && branch != "" {
			return branch
		}
	}

	// Default fallbacks
	return "main"
}