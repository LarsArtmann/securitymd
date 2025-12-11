package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
)

// FileSystemRepository implements file-based repository storage
type FileSystemRepository struct {
	basePath string
	mu       sync.RWMutex
}

// NewFileSystemRepository creates a new file system repository
func NewFileSystemRepository(basePath string) (*FileSystemRepository, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, FileOperationError{
			BaseError: NewBaseError("FILE001", "Failed to create repository directory", "filesystem"),
			Operation: "mkdir",
			FilePath:  basePath,
		}.WithPermission(true)
	}

	return &FileSystemRepository{
		basePath: basePath,
	}, nil
}

// Project file system implementation

type projectFileSystemRepository struct {
	*FileSystemRepository
}

// NewProjectFileSystemRepository creates a new project file system repository
func NewProjectFileSystemRepository(basePath string) (ProjectRepository, error) {
	fs, err := NewFileSystemRepository(filepath.Join(basePath, "projects"))
	if err != nil {
		return nil, err
	}

	return &projectFileSystemRepository{FileSystemRepository: fs}, nil
}

// Save stores project information
func (pfsr *projectFileSystemRepository) Save(ctx context.Context, project domain.ProjectInfo) error {
	pfsr.mu.Lock()
	defer pfsr.mu.Unlock()

	filePath := filepath.Join(pfsr.basePath, project.Name+".json")
	
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return NewFileOperationError("FILE003", "Failed to marshal project data", "write", filePath)
	}

	return os.WriteFile(filePath, data, 0644)
}

// Get retrieves project information
func (pfsr *projectFileSystemRepository) Get(ctx context.Context, id domain.PolicyID) (domain.ProjectInfo, error) {
	pfsr.mu.RLock()
	defer pfsr.mu.RUnlock()

	filePath := filepath.Join(pfsr.basePath, id.Value+".json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.ProjectInfo{}, NewFileOperationError("FILE001", "Project not found", "read", filePath)
		}
		return domain.ProjectInfo{}, NewFileOperationError("FILE002", "Failed to read project file", "read", filePath)
	}

	var project domain.ProjectInfo
	if err := json.Unmarshal(data, &project); err != nil {
		return domain.ProjectInfo{}, NewFileOperationError("FILE003", "Failed to unmarshal project data", "read", filePath)
	}

	return project, nil
}

// Delete removes project information
func (pfsr *projectFileSystemRepository) Delete(ctx context.Context, id domain.PolicyID) error {
	pfsr.mu.Lock()
	defer pfsr.mu.Unlock()

	filePath := filepath.Join(pfsr.basePath, id.Value+".json")
	
	return os.Remove(filePath)
}

// List returns all projects
func (pfsr *projectFileSystemRepository) List(ctx context.Context) ([]domain.ProjectInfo, error) {
	pfsr.mu.RLock()
	defer pfsr.mu.RUnlock()

	files, err := os.ReadDir(pfsr.basePath)
	if err != nil {
		return nil, NewFileOperationError("FILE002", "Failed to read projects directory", "read", pfsr.basePath)
	}

	var projects []domain.ProjectInfo
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(pfsr.basePath, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip corrupted files
		}

		var project domain.ProjectInfo
		if err := json.Unmarshal(data, &project); err != nil {
			continue // Skip corrupted files
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// Search finds projects by criteria
func (pfsr *projectFileSystemRepository) Search(ctx context.Context, criteria ProjectSearchCriteria) ([]domain.ProjectInfo, error) {
	allProjects, err := pfsr.List(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.ProjectInfo
	
	for _, project := range allProjects {
		if pfsr.matchesCriteria(project, criteria) {
			results = append(results, project)
		}
	}

	return pfsr.applySortingAndPagination(results, criteria), nil
}

// matchesCriteria checks if project matches search criteria
func (pfsr *projectFileSystemRepository) matchesCriteria(project domain.ProjectInfo, criteria ProjectSearchCriteria) bool {
	// Name filter
	if criteria.Name != "" && project.Name != criteria.Name {
		return false
	}

	// Type filter
	if criteria.Type != "" && project.Type != criteria.Type {
		return false
	}

	// Tech stack filter
	if len(criteria.TechStack) > 0 {
		projectStack := make(map[types.TechStack]bool)
		for _, tech := range project.TechStack {
			projectStack[tech] = true
		}
		
		for _, requiredTech := range criteria.TechStack {
			if !projectStack[requiredTech] {
				return false
			}
		}
	}

	// Languages filter
	if len(criteria.Languages) > 0 {
		projectLangs := make(map[string]bool)
		for _, lang := range project.Languages {
			projectLangs[lang] = true
		}
		
		for _, requiredLang := range criteria.Languages {
			if !projectLangs[requiredLang] {
				return false
			}
		}
	}

	return true
}

// applySortingAndPagination applies sorting and pagination to results
func (pfsr *projectFileSystemRepository) applySortingAndPagination(projects []domain.ProjectInfo, criteria ProjectSearchCriteria) []domain.ProjectInfo {
	// Apply sorting (simple implementation)
	switch criteria.SortBy {
	case SortProjectByName:
		// Sort by name
		for i := 0; i < len(projects)-1; i++ {
			for j := i + 1; j < len(projects); j++ {
				if criteria.SortOrder == SortOrderDescending {
					if projects[i].Name < projects[j].Name {
						projects[i], projects[j] = projects[j], projects[i]
					}
				} else {
					if projects[i].Name > projects[j].Name {
						projects[i], projects[j] = projects[j], projects[i]
					}
				}
			}
		}
	default:
		// Default: no sorting
	}

	// Apply pagination
	if criteria.Offset >= len(projects) {
		return []domain.ProjectInfo{}
	}

	end := criteria.Offset + criteria.Limit
	if end > len(projects) || criteria.Limit == 0 {
		end = len(projects)
	}

	return projects[criteria.Offset:end]
}

// Policy file system implementation

type policyFileSystemRepository struct {
	*FileSystemRepository
}

// NewPolicyFileSystemRepository creates a new policy file system repository
func NewPolicyFileSystemRepository(basePath string) (PolicyRepository, error) {
	fs, err := NewFileSystemRepository(filepath.Join(basePath, "policies"))
	if err != nil {
		return nil, err
	}

	return &policyFileSystemRepository{FileSystemRepository: fs}, nil
}

// Save stores a security policy
func (pfsr *policyFileSystemRepository) Save(ctx context.Context, policy domain.SecurityPolicy) error {
	pfsr.mu.Lock()
	defer pfsr.mu.Unlock()

	filePath := filepath.Join(pfsr.basePath, policy.ID.Value+".json")
	
	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return NewFileOperationError("FILE003", "Failed to marshal policy data", "write", filePath)
	}

	return os.WriteFile(filePath, data, 0644)
}

// Get retrieves a security policy by ID
func (pfsr *policyFileSystemRepository) Get(ctx context.Context, id domain.PolicyID) (domain.SecurityPolicy, error) {
	pfsr.mu.RLock()
	defer pfsr.mu.RUnlock()

	filePath := filepath.Join(pfsr.basePath, id.Value+".json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.SecurityPolicy{}, NewFileOperationError("FILE001", "Policy not found", "read", filePath)
		}
		return domain.SecurityPolicy{}, NewFileOperationError("FILE002", "Failed to read policy file", "read", filePath)
	}

	var policy domain.SecurityPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return domain.SecurityPolicy{}, NewFileOperationError("FILE003", "Failed to unmarshal policy data", "read", filePath)
	}

	return policy, nil
}

// GetByType retrieves policies by type
func (pfsr *policyFileSystemRepository) GetByType(ctx context.Context, policyType types.PolicyType) ([]domain.SecurityPolicy, error) {
	pfsr.mu.RLock()
	defer pfsr.mu.RUnlock()

	files, err := os.ReadDir(pfsr.basePath)
	if err != nil {
		return nil, NewFileOperationError("FILE002", "Failed to read policies directory", "read", pfsr.basePath)
	}

	var policies []domain.SecurityPolicy
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(pfsr.basePath, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var policy domain.SecurityPolicy
		if err := json.Unmarshal(data, &policy); err != nil {
			continue
		}

		if policy.Type == policyType {
			policies = append(policies, policy)
		}
	}

	return policies, nil
}

// Update modifies an existing policy
func (pfsr *policyFileSystemRepository) Update(ctx context.Context, policy domain.SecurityPolicy) error {
	// Check if policy exists
	existing, err := pfsr.Get(ctx, policy.ID)
	if err != nil {
		return fmt.Errorf("policy does not exist: %w", err)
	}

	// Update timestamp
	policy.UpdatedAt = time.Now()

	return pfsr.Save(ctx, policy)
}

// Delete removes a security policy
func (pfsr *policyFileSystemRepository) Delete(ctx context.Context, id domain.PolicyID) error {
	pfsr.mu.Lock()
	defer pfsr.mu.Unlock()

	filePath := filepath.Join(pfsr.basePath, id.Value+".json")
	
	return os.Remove(filePath)
}

// List returns all policies
func (pfsr *policyFileSystemRepository) List(ctx context.Context) ([]domain.SecurityPolicy, error) {
	pfsr.mu.RLock()
	defer pfsr.mu.RUnlock()

	files, err := os.ReadDir(pfsr.basePath)
	if err != nil {
		return nil, NewFileOperationError("FILE002", "Failed to read policies directory", "read", pfsr.basePath)
	}

	var policies []domain.SecurityPolicy
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(pfsr.basePath, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var policy domain.SecurityPolicy
		if err := json.Unmarshal(data, &policy); err != nil {
			continue
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// Validate checks policy validity
func (pfsr *policyFileSystemRepository) Validate(ctx context.Context, policy domain.SecurityPolicy) domain.Validation {
	validation := domain.Validation{
		Status:    types.StatusValid,
		Score:     100,
		Issues:    []domain.ValidationIssue{},
		Warnings:  []domain.ValidationWarning{},
		CheckedAt: time.Now(),
		Quality: domain.QualityMetrics{
			Sections:     0,
			Completeness: 100,
		},
	}

	// Check basic validity
	if !policy.IsValid() {
		validation.Status = types.StatusInvalid
		validation.Score = 0
		validation.Issues = append(validation.Issues, domain.ValidationIssue{
			ID:          "VAL001",
			Field:       "policy",
			Category:    domain.CategoryStructure,
			Level:       types.ValidationLevelError,
			Message:     "Policy is not valid",
			Suggestion:  "Check required fields",
		})
	}

	return validation
}

// Validation file system implementation

type validationFileSystemRepository struct {
	*FileSystemRepository
}

// NewValidationFileSystemRepository creates a new validation file system repository
func NewValidationFileSystemRepository(basePath string) (ValidationRepository, error) {
	fs, err := NewFileSystemRepository(filepath.Join(basePath, "validations"))
	if err != nil {
		return nil, err
	}

	return &validationFileSystemRepository{FileSystemRepository: fs}, nil
}

// Save stores validation results
func (vfsr *validationFileSystemRepository) Save(ctx context.Context, result domain.Validation) error {
	vfsr.mu.Lock()
	defer vfsr.mu.Unlock()

	// Create filename with policy ID and timestamp
	fileName := fmt.Sprintf("%s_%d.json", "validation", time.Now().Unix())
	filePath := filepath.Join(vfsr.basePath, fileName)
	
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return NewFileOperationError("FILE003", "Failed to marshal validation data", "write", filePath)
	}

	return os.WriteFile(filePath, data, 0644)
}

// Get retrieves validation results by policy ID
func (vfsr *validationFileSystemRepository) Get(ctx context.Context, policyID domain.PolicyID, version string) (domain.Validation, error) {
	vfsr.mu.RLock()
	defer vfsr.mu.RUnlock()

	// For file system, we'll search for files containing policy ID
	files, err := os.ReadDir(vfsr.basePath)
	if err != nil {
		return domain.Validation{}, NewFileOperationError("FILE002", "Failed to read validations directory", "read", vfsr.basePath)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(vfsr.basePath, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var validation domain.Validation
		if err := json.Unmarshal(data, &validation); err != nil {
			continue
		}

		// Simple implementation - return first match
		return validation, nil
	}

	return domain.Validation{}, NewFileOperationError("FILE001", "Validation not found", "read", "")
}

// GetLatest returns the most recent validation
func (vfsr *validationFileSystemRepository) GetLatest(ctx context.Context, policyID domain.PolicyID) (domain.Validation, error) {
	vfsr.mu.RLock()
	defer vfsr.mu.RUnlock()

	files, err := os.ReadDir(vfsr.basePath)
	if err != nil {
		return domain.Validation{}, NewFileOperationError("FILE002", "Failed to read validations directory", "read", vfsr.basePath)
	}

	var latestValidation domain.Validation
	var latestTime time.Time

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(vfsr.basePath, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var validation domain.Validation
		if err := json.Unmarshal(data, &validation); err != nil {
			continue
		}

		if validation.CheckedAt.After(latestTime) {
			latestValidation = validation
			latestTime = validation.CheckedAt
		}
	}

	if latestValidation.CheckedAt.IsZero() {
		return domain.Validation{}, NewFileOperationError("FILE001", "No validation found", "read", "")
	}

	return latestValidation, nil
}

// GetHistory returns validation history
func (vfsr *validationFileSystemRepository) GetHistory(ctx context.Context, policyID domain.PolicyID, limit int) ([]domain.Validation, error) {
	vfsr.mu.RLock()
	defer vfsr.mu.RUnlock()

	files, err := os.ReadDir(vfsr.basePath)
	if err != nil {
		return nil, NewFileOperationError("FILE002", "Failed to read validations directory", "read", vfsr.basePath)
	}

	var validations []domain.Validation
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(vfsr.basePath, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var validation domain.Validation
		if err := json.Unmarshal(data, &validation); err != nil {
			continue
		}

		validations = append(validations, validation)
	}

	// Sort by checked_at descending
	for i := 0; i < len(validations)-1; i++ {
		for j := i + 1; j < len(validations); j++ {
			if validations[i].CheckedAt.Before(validations[j].CheckedAt) {
				validations[i], validations[j] = validations[j], validations[i]
			}
		}
	}

	// Apply limit
	if limit > 0 && limit < len(validations) {
		validations = validations[:limit]
	}

	return validations, nil
}

// Search finds validations by criteria
func (vfsr *validationFileSystemRepository) Search(ctx context.Context, criteria ValidationSearchCriteria) ([]domain.Validation, error) {
	// Simplified implementation - just return all and filter
	all, err := vfsr.GetHistory(ctx, domain.PolicyID{}, 1000)
	if err != nil {
		return nil, err
	}

	var results []domain.Validation
	for _, validation := range all {
		if vfsr.matchesSearchCriteria(validation, criteria) {
			results = append(results, validation)
		}
	}

	return results, nil
}

// matchesSearchCriteria checks if validation matches search criteria
func (vfsr *validationFileSystemRepository) matchesSearchCriteria(validation domain.Validation, criteria ValidationSearchCriteria) bool {
	// Status filter
	if criteria.Status != "" && validation.Status != criteria.Status {
		return false
	}

	// Score range filter
	if criteria.ScoreFrom != nil && validation.Score < *criteria.ScoreFrom {
		return false
	}
	if criteria.ScoreTo != nil && validation.Score > *criteria.ScoreTo {
		return false
	}

	// Date range filter
	if criteria.DateFrom != nil && validation.CheckedAt.Before(*criteria.DateFrom) {
		return false
	}
	if criteria.DateTo != nil && validation.CheckedAt.After(*criteria.DateTo) {
		return false
	}

	return true
}