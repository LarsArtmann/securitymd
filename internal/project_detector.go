package internal

import (
	"os"
	"path/filepath"
	"strings"
)

// ProjectDetector detects project information from various sources
type ProjectDetector struct{}

// NewProjectDetector creates a new project detector
func NewProjectDetector() *ProjectDetector {
	return &ProjectDetector{}
}

// DetectProjectName attempts to detect the project name from various sources
func (pd *ProjectDetector) DetectProjectName() string {
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

// DetectDomain attempts to detect the organization domain
func (pd *ProjectDetector) DetectDomain() string {
	// Check git remote for domain
	if domain := pd.detectFromGitRemote(); domain != "" {
		return domain
	}

	// Fallback to common patterns
	if strings.Contains(pd.DetectProjectName(), "-") {
		parts := strings.Split(pd.DetectProjectName(), "-")
		if len(parts) > 1 {
			return parts[0] + ".com"
		}
	}

	return "example.com"
}

// detectFromPackageJSON tries to detect project name from package.json
func (pd *ProjectDetector) detectFromPackageJSON() string {
	content, err := os.ReadFile("package.json")
	if err != nil {
		return ""
	}

	// Simple string parsing for now
	// In a real implementation, use JSON parsing
	contentStr := string(content)

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

// detectFromGoMod tries to detect project name from go.mod
func (pd *ProjectDetector) detectFromGoMod() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			module := strings.TrimPrefix(line, "module ")
			module = strings.TrimSpace(module)

			// Extract last part of module path
			parts := strings.Split(module, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}

	return ""
}

// detectFromCargoToml tries to detect project name from Cargo.toml
func (pd *ProjectDetector) detectFromCargoToml() string {
	content, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return ""
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name = ") {
			name := strings.TrimPrefix(line, "name = ")
			name = strings.TrimSpace(name)

			// Remove quotes
			return strings.Trim(name, "\"")
		}
	}

	return ""
}

// detectFromPyProjectToml tries to detect project name from pyproject.toml
func (pd *ProjectDetector) detectFromPyProjectToml() string {
	content, err := os.ReadFile("pyproject.toml")
	if err != nil {
		return ""
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name = ") {
			name := strings.TrimPrefix(line, "name = ")
			name = strings.TrimSpace(name)

			// Remove quotes
			return strings.Trim(name, "\"")
		}
	}

	return ""
}

// detectFromDirectoryName uses the current directory name
func (pd *ProjectDetector) detectFromDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Base(dir)
}

// detectFromGitRemote tries to detect domain from git remote
func (pd *ProjectDetector) detectFromGitRemote() string {
	// This is a simplified implementation
	// In a real implementation, parse git remote output
	return ""
}
