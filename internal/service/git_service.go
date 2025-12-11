package service

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
)

// gitFileSystemService implements GitService interface
type gitFileSystemService struct {
	errorHandler ErrorHandler
}

// NewGitService creates a new Git service
func NewGitService() GitService {
	errorHandler := NewErrorHandler(true)
	return &gitFileSystemService{
		errorHandler: errorHandler,
	}
}

// GetRepositoryName gets repository name from git config
func (g *gitFileSystemService) GetRepositoryName(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewProjectDetectionError(
			"DET001",
			"Failed to get git remote URL",
			"git",
			"GetRepositoryName",
		)
	}

	url := strings.TrimSpace(string(output))
	if url == "" {
		return "", nil
	}

	// Parse repository name from URL
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", nil
	}

	repoName := parts[len(parts)-1]
	return strings.TrimSuffix(repoName, ".git"), nil
}

// GetUserEmail gets user email from git config
func (g *gitFileSystemService) GetUserEmail(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewProjectDetectionError(
			"DET001",
			"Failed to get git user email",
			"git",
			"GetUserEmail",
		)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetRemoteURL gets git remote URL
func (g *gitFileSystemService) GetRemoteURL(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewProjectDetectionError(
			"DET001",
			"Failed to get git remote URL",
			"git",
			"GetRemoteURL",
		)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch gets current branch
func (g *gitFileSystemService) GetCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to main as default
		return "main", nil
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "main", nil
	}

	return branch, nil
}

// IsGitRepository checks if current directory is a git repository
func (g *gitFileSystemService) IsGitRepository(ctx context.Context) bool {
	_, err := os.Stat(".git")
	return err == nil
}