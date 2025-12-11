package service

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
)

// projectDetectionService implements ProjectDetectionService interface
type projectDetectionService struct {
	gitService   GitService
	errorHandler ErrorHandler
}

// NewProjectDetectionService creates a new project detection service
func NewProjectDetectionService() ProjectDetectionService {
	errorHandler := NewErrorHandler(true)
	gitService := NewGitService()

	return &projectDetectionService{
		gitService:   gitService,
		errorHandler: errorHandler,
	}
}

// DetectProjectName detects project name from various sources
func (p *projectDetectionService) DetectProjectName(ctx context.Context) (string, error) {
	// Try git repository name
	if p.gitService.IsGitRepository(ctx) {
		repoName, err := p.gitService.GetRepositoryName(ctx)
		if err == nil && repoName != "" {
			return repoName, nil
		}
	}

	// Try go.mod
	if goModName := p.detectFromGoMod(); goModName != "" {
		return goModName, nil
	}

	// Try package.json
	if packageName := p.detectFromPackageJSON(); packageName != "" {
		return packageName, nil
	}

	// Try directory name
	dir, err := os.Getwd()
	if err == nil {
		dirName := filepath.Base(dir)
		if dirName != "." && dirName != "/" {
			return dirName, nil
		}
	}

	return "MyProject", nil
}

// DetectOrganization detects organization name
func (p *projectDetectionService) DetectOrganization(ctx context.Context) (string, error) {
	// Try git remote
	if p.gitService.IsGitRepository(ctx) {
		remoteURL, err := p.gitService.GetRemoteURL(ctx)
		if err == nil && remoteURL != "" {
			if org := p.extractOrganizationFromURL(remoteURL); org != "" {
				return org, nil
			}
		}
	}

	// Try go.mod
	if org := p.detectFromGoMod(); org != "" {
		return org, nil
	}

	return "MyProject", nil
}

// DetectDomain detects organization domain
func (p *projectDetectionService) DetectDomain(ctx context.Context) (string, error) {
	// Try git remote
	if p.gitService.IsGitRepository(ctx) {
		remoteURL, err := p.gitService.GetRemoteURL(ctx)
		if err == nil && remoteURL != "" {
			if domain := p.extractDomainFromURL(remoteURL); domain != "" {
				return domain, nil
			}
		}
	}

	return "example.com", nil
}

// DetectTechnologyStack detects technology stack
func (p *projectDetectionService) DetectTechnologyStack(ctx context.Context) ([]types.TechStack, error) {
	var stack []types.TechStack

	// Detect from files
	if p.fileExists("go.mod") {
		stack = append(stack, types.TechStackGo)
	}
	if p.fileExists("package.json") {
		stack = append(stack, types.TechStackNode)
	}
	if p.fileExists("pyproject.toml") || p.fileExists("requirements.txt") {
		stack = append(stack, types.TechStackPython)
	}
	if p.fileExists("Cargo.toml") {
		stack = append(stack, types.TechStackRust)
	}
	if p.fileExists("Dockerfile") {
		stack = append(stack, types.TechStackDocker)
	}
	if p.fileExists("docker-compose.yml") || p.fileExists("docker-compose.yaml") {
		stack = append(stack, types.TechStackDocker)
	}

	return stack, nil
}

// DetectProjectType detects project type
func (p *projectDetectionService) DetectProjectType(ctx context.Context) (types.ProjectType, error) {
	// Detect based on files
	if p.fileExists("package.json") {
		content, err := os.ReadFile("package.json")
		if err == nil {
			jsonContent := string(content)
			if strings.Contains(jsonContent, "\"react\"") || strings.Contains(jsonContent, "\"vue\"") || strings.Contains(jsonContent, "\"angular\"") {
				return types.TypeWeb, nil
			}
			if strings.Contains(jsonContent, "\"express\"") || strings.Contains(jsonContent, "\"fastify\"") {
				return types.TypeAPI, nil
			}
		}
	}

	if p.fileExists("go.mod") {
		return types.TypeAPI, nil
	}

	if p.fileExists("Cargo.toml") {
		content, _ := os.ReadFile("Cargo.toml")
		if strings.Contains(string(content), "[bin]") {
			return types.TypeCLI, nil
		}
		return types.TypeLibrary, nil
	}

	return types.TypeWeb, nil
}

// AnalyzeProject analyzes project structure and metadata
func (p *projectDetectionService) AnalyzeProject(ctx context.Context) (domain.ProjectInfo, error) {
	projectName, _ := p.DetectProjectName(ctx)
	projectType, _ := p.DetectProjectType(ctx)
	techStack, _ := p.DetectTechnologyStack(ctx)

	languages := p.detectLanguages()
	frameworks := p.detectFrameworks()
	dependencies := p.detectDependencies()

	return domain.ProjectInfo{
		Name:         projectName,
		Type:         projectType,
		TechStack:    techStack,
		Languages:    languages,
		Frameworks:   frameworks,
		Dependencies: dependencies,
		Size:         p.calculateProjectSize(),
		Complexity:   p.calculateComplexity(),
	}, nil
}

// Private helper methods

func (p *projectDetectionService) detectFromGoMod() string {
	if !p.fileExists("go.mod") {
		return ""
	}

	content, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}

	// Extract module name
	moduleRegex := regexp.MustCompile(`module\s+([^\s]+)`)
	matches := moduleRegex.FindSubmatch(content)
	if len(matches) > 1 {
		moduleName := string(matches[1])
		parts := strings.Split(moduleName, "/")
		return parts[len(parts)-1]
	}

	return ""
}

func (p *projectDetectionService) detectFromPackageJSON() string {
	if !p.fileExists("package.json") {
		return ""
	}

	content, err := os.ReadFile("package.json")
	if err != nil {
		return ""
	}

	// Extract name field
	nameRegex := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`)
	matches := nameRegex.FindSubmatch(content)
	if len(matches) > 1 {
		return string(matches[1])
	}

	return ""
}

func (p *projectDetectionService) extractOrganizationFromURL(url string) string {
	// Handle GitHub URLs
	if strings.Contains(url, "github.com") {
		parts := strings.Split(url, "/")
		if len(parts) >= 3 {
			return parts[2] // github.com/owner/repo
		}
	}

	return ""
}

func (p *projectDetectionService) extractDomainFromURL(url string) string {
	// For now, generate domain from organization
	org := p.extractOrganizationFromURL(url)
	if org != "" {
		return strings.ToLower(org) + ".com"
	}

	return "example.com"
}

func (p *projectDetectionService) detectLanguages() []string {
	var languages []string

	// Simple file-based detection
	if p.fileExists("go.mod") {
		languages = append(languages, "Go")
	}
	if p.fileExists("package.json") {
		languages = append(languages, "JavaScript")
	}
	if p.fileExists("pyproject.toml") || p.fileExists("requirements.txt") {
		languages = append(languages, "Python")
	}
	if p.fileExists("Cargo.toml") {
		languages = append(languages, "Rust")
	}

	return languages
}

func (p *projectDetectionService) detectFrameworks() []string {
	var frameworks []string

	if p.fileExists("package.json") {
		content, _ := os.ReadFile("package.json")
		jsonContent := string(content)

		if strings.Contains(jsonContent, "\"react\"") {
			frameworks = append(frameworks, "React")
		}
		if strings.Contains(jsonContent, "\"vue\"") {
			frameworks = append(frameworks, "Vue.js")
		}
		if strings.Contains(jsonContent, "\"angular\"") {
			frameworks = append(frameworks, "Angular")
		}
		if strings.Contains(jsonContent, "\"express\"") {
			frameworks = append(frameworks, "Express.js")
		}
	}

	return frameworks
}

func (p *projectDetectionService) detectDependencies() []string {
	var dependencies []string

	// Go dependencies
	if p.fileExists("go.mod") {
		content, _ := os.ReadFile("go.mod")
		modContent := string(content)

		// Extract direct dependencies
		lines := strings.Split(modContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "\t//") && strings.Contains(line, " ") {
				dep := strings.Fields(line)[0]
				if dep != "" {
					dependencies = append(dependencies, dep)
				}
			}
		}
	}

	// Node.js dependencies
	if p.fileExists("package.json") {
		content, _ := os.ReadFile("package.json")
		jsonContent := string(content)

		// Simple regex to extract dependencies
		depRegex := regexp.MustCompile(`"([^"]+)"\s*:\s*"[^"]+"`)
		matches := depRegex.FindAllStringSubmatch(jsonContent, -1)
		for _, match := range matches {
			if len(match) > 1 && match[1] != "name" && match[1] != "version" {
				dependencies = append(dependencies, match[1])
			}
		}
	}

	return dependencies
}

func (p *projectDetectionService) calculateProjectSize() domain.ProjectSize {
	// Count files (simplified)
	fileCount := p.countFiles(".")

	switch {
	case fileCount <= 10:
		return domain.SizeMicro
	case fileCount <= 50:
		return domain.SizeSmall
	case fileCount <= 200:
		return domain.SizeMedium
	case fileCount <= 1000:
		return domain.SizeLarge
	default:
		return domain.SizeXLarge
	}
}

func (p *projectDetectionService) calculateComplexity() domain.ProjectComplexity {
	// Simple complexity calculation based on file count and dependencies
	fileCount := p.countFiles(".")
	depCount := len(p.detectDependencies())

	score := fileCount + (depCount * 2)

	switch {
	case score <= 20:
		return domain.ComplexitySimple
	case score <= 100:
		return domain.ComplexityModerate
	case score <= 500:
		return domain.ComplexityComplex
	default:
		return domain.ComplexityVeryComplex
	}
}

func (p *projectDetectionService) countFiles(dir string) int {
	count := 0

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") && path != "." {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip vendor, node_modules, etc.
		if strings.Contains(path, "vendor/") || strings.Contains(path, "node_modules/") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			count++
		}

		return nil
	})

	return count
}

func (p *projectDetectionService) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
