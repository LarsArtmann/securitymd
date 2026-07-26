package internal

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectDetector detects project information from various sources.
type ProjectDetector struct{}

// lastPartFromEnd extracts the nth part from the end of a split string.
// For example, lastPartFromEnd("org/repo", "/", 1) returns "repo",
// and lastPartFromEnd("org/repo", "/", 2) returns "org".
func lastPartFromEnd(s, sep string, n int) string {
	parts := strings.Split(s, sep)

	idx := len(parts) - n
	if idx >= 0 {
		return parts[idx]
	}

	return ""
}

// loadFileResult holds a file's contents and tracks whether the file was
// readable. Callers should use Value() to obtain the string, which yields
// the empty string when the file was unreadable.
type loadFileResult struct {
	content string
	ok      bool
}

// Value returns the file's contents, or the empty string when unreadable.
func (r loadFileResult) Value() string {
	if !r.ok {
		return ""
	}

	return r.content
}

// loadFile reads the file at path and packages the result.
func loadFile(path string) loadFileResult {
	content, err := os.ReadFile(path)
	if err != nil {
		return loadFileResult{}
	}

	return loadFileResult{content: string(content), ok: true}
}

// firstPathPart returns the portion of value before the first "/" separator,
// or the whole value when no separator is present.
func firstPathPart(value string) string {
	if part, _, ok := strings.Cut(value, "/"); ok {
		return part
	}

	return value
}

// NewProjectDetector creates a new project detector.
func NewProjectDetector() *ProjectDetector {
	return &ProjectDetector{}
}

// DetectProjectName attempts to detect the project name from various sources.
func (pd *ProjectDetector) DetectProjectName() string {
	// Try to detect from git remote first (more reliable)
	if name := pd.detectFromGitRemote(); name != "" {
		name = strings.TrimSuffix(name, ".git")

		return lastPartFromEnd(name, "/", 1)
	}

	// Try to detect from package.json
	if name := pd.detectFromPackageJSON(); name != "" {
		return name
	}

	// Try to detect from go.mod
	if name := pd.detectFromGoMod(); name != "" {
		return name
	}

	// Try to detect from Cargo.toml
	if name := pd.detectFromCargoToml(); name != "" {
		return name
	}

	// Try to detect from pyproject.toml
	if name := pd.detectFromPyProjectToml(); name != "" {
		return name
	}

	// Fallback to directory name
	if name := pd.detectFromDirectoryName(); name != "" {
		return name
	}

	return "MyProject"
}

// DetectOrganization attempts to detect the organization name.
func (pd *ProjectDetector) DetectOrganization() string {
	// Try to detect from git remote first
	if org := pd.detectOrgFromGitRemote(); org != "" {
		return org
	}

	// Try to detect from package.json
	if org := pd.detectOrgFromPackageJSON(); org != "" {
		return org
	}

	// Fallback to project name
	return pd.DetectProjectName()
}

// DetectDomain attempts to detect the organization domain.
func (pd *ProjectDetector) DetectDomain() string {
	// Check git remote for domain
	if domain := pd.detectDomainFromGitRemote(); domain != "" {
		return domain
	}

	// Try to construct domain from organization name
	org := pd.DetectOrganization()
	if org != "" {
		// Common domain patterns
		if strings.Contains(org, "-") {
			parts := strings.Split(org, "-")
			if len(parts) > 1 {
				return strings.ToLower(parts[0]) + "." + strings.ToLower(parts[1]) + ".com"
			}
		}

		// Default pattern
		return strings.ToLower(org) + ".com"
	}

	return "example.com"
}

// detectFromGitRemote tries to detect git remote URL.
func (pd *ProjectDetector) detectFromGitRemote() string {
	cmd := exec.CommandContext(context.Background(), "git", "config", "--get", "remote.origin.url")

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// detectOrgFromGitRemote tries to extract organization from git remote.
func (pd *ProjectDetector) detectOrgFromGitRemote() string {
	gitURL := strings.TrimSpace(pd.detectFromGitRemote())
	if gitURL == "" {
		return ""
	}

	// Remove .git suffix
	gitURL = strings.TrimSuffix(gitURL, ".git")

	// Parse HTTPS URLs
	if after, ok := strings.CutPrefix(gitURL, "https://"); ok {
		gitURL = after

		parts := strings.Split(gitURL, "/")
		if len(parts) >= 2 {
			return parts[1] // github.com/ORG/repo -> parts[1] is ORG
		}
	}

	// Parse SSH URLs
	if strings.Contains(gitURL, "@") {
		parts := strings.Split(gitURL, ":")
		if len(parts) >= 2 {
			pathParts := strings.Split(parts[1], "/")
			if len(pathParts) >= 2 {
				return pathParts[0]
			}
		}
	}

	// Parse generic path
	org := lastPartFromEnd(gitURL, "/", 2)
	if org != "" {
		return org
	}

	return ""
}

// detectDomainFromGitRemote tries to extract domain from git remote.
func (pd *ProjectDetector) detectDomainFromGitRemote() string {
	// Extract domain from HTTPS URL
	if after, ok := strings.CutPrefix(strings.TrimSpace(pd.detectFromGitRemote()), "https://"); ok {
		return firstPathPart(after)
	}

	// Extract domain from SSH URL
	if _, after, ok := strings.Cut(strings.TrimSpace(pd.detectFromGitRemote()), "@"); ok {
		return firstPathPart(after)
	}

	return ""
}

// detectFromPackageJSON tries to detect project name from package.json.
func (pd *ProjectDetector) detectFromPackageJSON() string {
	// Simple string parsing for now
	// In a real implementation, use JSON parsing
	contentStr := loadFile("package.json").Value()

	// Look for "name": "project-name"
	start := strings.Index(contentStr, "\"name\"")
	if start == -1 {
		return ""
	}

	// Find the colon after "name"
	colon := strings.Index(contentStr[start:], ":")
	if colon == -1 {
		return ""
	}

	// Find the first quote after colon
	firstQuote := strings.Index(contentStr[start+colon:], "\"")
	if firstQuote == -1 {
		return ""
	}

	// Find the closing quote
	nameStart := start + colon + firstQuote + 1

	nameEnd := strings.Index(contentStr[nameStart:], "\"")
	if nameEnd == -1 {
		return ""
	}

	name := strings.TrimSpace(contentStr[nameStart : nameStart+nameEnd])

	return strings.Trim(name, "\"")
}

// detectOrgFromPackageJSON tries to detect organization from package.json.
func (pd *ProjectDetector) detectOrgFromPackageJSON() string {
	contentStr := loadFile("package.json").Value()

	// Try to detect from author field
	if strings.Contains(contentStr, "\"author\"") {
		start := strings.Index(contentStr, "\"author\"")
		if start != -1 {
			// Look for organization in author field
			if strings.Contains(contentStr[start:start+500], "\"organization\"") {
				orgStart := strings.Index(contentStr[start:], "\"organization\"")
				if orgStart != -1 {
					// Extract organization value
					remainder := contentStr[start+orgStart:]

					colon := strings.Index(remainder, ":")
					if colon != -1 {
						quote := strings.Index(remainder[colon:], "\"")
						if quote != -1 {
							orgStart := colon + quote + 1

							orgEnd := strings.Index(remainder[orgStart:], "\"")
							if orgEnd != -1 {
								org := strings.TrimSpace(remainder[orgStart : orgStart+orgEnd])

								return strings.Trim(org, "\"")
							}
						}
					}
				}
			}
		}
	}

	// Try to detect from scope in package name (@org/package)
	name := pd.detectFromPackageJSON()
	if after, ok := strings.CutPrefix(name, "@"); ok {
		if parts := strings.Split(after, "/"); len(parts) > 0 {
			return parts[0]
		}
	}

	return ""
}

// detectFromGoMod tries to detect project name from go.mod.
func (pd *ProjectDetector) detectFromGoMod() string {
	contentStr := loadFile("go.mod").Value()
	lines := strings.SplitSeq(contentStr, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "module "); ok {
			module := after
			module = strings.TrimSpace(module)

			// Extract last part of module path
			return lastPartFromEnd(module, "/", 1)
		}
	}

	return ""
}

// detectNameFromTomlFile extracts the project name from a TOML file.
func (pd *ProjectDetector) detectNameFromTomlFile(filename string) string {
	contentStr := loadFile(filename).Value()
	lines := strings.SplitSeq(contentStr, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "name = "); ok {
			name := after
			name = strings.TrimSpace(name)

			// Remove quotes
			return strings.Trim(name, "\"")
		}
	}

	return ""
}

// detectFromCargoToml tries to detect project name from Cargo.toml.
func (pd *ProjectDetector) detectFromCargoToml() string {
	return pd.detectNameFromTomlFile("Cargo.toml")
}

// detectFromPyProjectToml tries to detect project name from pyproject.toml.
func (pd *ProjectDetector) detectFromPyProjectToml() string {
	return pd.detectNameFromTomlFile("pyproject.toml")
}

// detectFromDirectoryName uses the current directory name.
func (pd *ProjectDetector) detectFromDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Base(dir)
}
