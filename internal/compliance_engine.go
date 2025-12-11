package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/LarsArtmann/template-CLI/pkg/sdk/fileops"
	"github.com/LarsArtmann/template-CLI/pkg/sdk/vfs"
	"github.com/samber/mo"
)

// ComplianceEngine handles compliance checking for various frameworks
type ComplianceEngine struct {
	fileSystem vfs.FileSystem
	fileOps    *fileops.Service
	validator  *PolicyValidator
}

// NewComplianceEngine creates new compliance engine
func NewComplianceEngine() *ComplianceEngine {
	vfsImpl := vfs.NewOSVFS()
	fileOps := fileops.NewService(fileops.ServiceOptions{
		BaseDirectory: ".",
		ValidatePaths: true,
	})

	return &ComplianceEngine{
		fileSystem: vfsImpl,
		fileOps:    fileOps,
		validator:  NewPolicyValidator(),
	}
}

// ComplianceFramework represents a compliance framework
type ComplianceFramework string

const (
	FrameworkGDPR     ComplianceFramework = "gdpr"
	FrameworkSOC2     ComplianceFramework = "soc2"
	FrameworkISO27001 ComplianceFramework = "iso27001"
	FrameworkNIST     ComplianceFramework = "nist"
)

// ComplianceResult represents compliance check results
type ComplianceResult struct {
	Framework       ComplianceFramework `json:"framework"`
	Score           float64             `json:"score"`
	TotalChecks     int                 `json:"total_checks"`
	PassedChecks    int                 `json:"passed_checks"`
	FailedChecks    []ComplianceIssue   `json:"failed_checks"`
	Recommendations []string            `json:"recommendations"`
	LastChecked     time.Time           `json:"last_checked"`
}

// ComplianceIssue represents a specific compliance issue
type ComplianceIssue struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "high", "medium", "low"
	Suggestion  string `json:"suggestion"`
}

// CheckCompliance runs compliance check for specified framework
func (ce *ComplianceEngine) CheckCompliance(ctx context.Context, framework ComplianceFramework) mo.Result[ComplianceResult] {
	switch framework {
	case FrameworkGDPR:
		return ce.checkGDPRCompliance(ctx)
	case FrameworkSOC2:
		return ce.checkSOC2Compliance(ctx)
	case FrameworkISO27001:
		return ce.checkISO27001Compliance(ctx)
	case FrameworkNIST:
		return ce.checkNISTCompliance(ctx)
	default:
		return mo.Err[ComplianceResult](fmt.Errorf("unsupported framework: %s", framework))
	}
}

// checkGDPRCompliance checks GDPR compliance
func (ce *ComplianceEngine) checkGDPRCompliance(ctx context.Context) mo.Result[ComplianceResult] {
	result := ComplianceResult{
		Framework:   FrameworkGDPR,
		LastChecked: time.Now(),
	}

	// Check for privacy policy
	privacyPolicyExists := ce.checkFileExists("privacy-policy.md")
	if privacyPolicyExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Data Protection",
			Description: "Privacy policy not found",
			Severity:    "high",
			Suggestion:  "Generate privacy-policy.md with GDPR compliance",
		})
		result.Recommendations = append(result.Recommendations, "Create GDPR-compliant privacy policy")
	}

	// Check for data protection officer
	dpoConfigured := ce.checkDPOConfigured()
	if dpoConfigured {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Governance",
			Description: "Data Protection Officer not configured",
			Severity:    "medium",
			Suggestion:  "Configure DPO contact information",
		})
		result.Recommendations = append(result.Recommendations, "Designate Data Protection Officer")
	}

	// Check for data breach procedures
	breachProceduresExist := ce.checkDataBreachProcedures()
	if breachProceduresExist {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Incident Response",
			Description: "Data breach procedures not documented",
			Severity:    "high",
			Suggestion:  "Create data breach response procedures",
		})
		result.Recommendations = append(result.Recommendations, "Document data breach response procedures")
	}

	// Check for data retention policy
	retentionPolicyExists := ce.checkDataRetentionPolicy()
	if retentionPolicyExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Data Management",
			Description: "Data retention policy not found",
			Severity:    "medium",
			Suggestion:  "Create data retention and deletion policy",
		})
		result.Recommendations = append(result.Recommendations, "Create data retention policy")
	}

	// Check for consent management
	consentConfigured := ce.checkConsentManagement()
	if consentConfigured {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "User Rights",
			Description: "Consent management not implemented",
			Severity:    "high",
			Suggestion:  "Implement user consent management system",
		})
		result.Recommendations = append(result.Recommendations, "Implement consent management")
	}

	result.TotalChecks = 5
	result.Score = float64(result.PassedChecks) / float64(result.TotalChecks) * 100

	return mo.Ok(result)
}

// checkSOC2Compliance checks SOC 2 Type II compliance
func (ce *ComplianceEngine) checkSOC2Compliance(ctx context.Context) mo.Result[ComplianceResult] {
	result := ComplianceResult{
		Framework:   FrameworkSOC2,
		LastChecked: time.Now(),
	}

	// Check for security policy
	securityPolicyExists := ce.checkFileExists("security-policy.md")
	if securityPolicyExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Security",
			Description: "Security policy not found",
			Severity:    "high",
			Suggestion:  "Generate comprehensive security policy",
		})
	}

	// Check for incident response plan
	incidentResponseExists := ce.checkFileExists("incident-response.md")
	if incidentResponseExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Incident Response",
			Description: "Incident response plan not documented",
			Severity:    "high",
			Suggestion:  "Create incident response procedures",
		})
	}

	// Check for access controls
	accessControlsExist := ce.checkAccessControls()
	if accessControlsExist {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Access Control",
			Description: "Access control procedures not documented",
			Severity:    "medium",
			Suggestion:  "Document access control procedures",
		})
	}

	// Check for monitoring
	monitoringConfigured := ce.checkMonitoringSystems()
	if monitoringConfigured {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Monitoring",
			Description: "Security monitoring not configured",
			Severity:    "medium",
			Suggestion:  "Implement security monitoring systems",
		})
	}

	// Check for change management
	changeManagementExists := ce.checkChangeManagement()
	if changeManagementExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Change Management",
			Description: "Change management procedures not documented",
			Severity:    "low",
			Suggestion:  "Document change management procedures",
		})
	}

	result.TotalChecks = 5
	result.Score = float64(result.PassedChecks) / float64(result.TotalChecks) * 100

	return mo.Ok(result)
}

// checkISO27001Compliance checks ISO 27001 compliance
func (ce *ComplianceEngine) checkISO27001Compliance(ctx context.Context) mo.Result[ComplianceResult] {
	result := ComplianceResult{
		Framework:   FrameworkISO27001,
		LastChecked: time.Now(),
	}

	// Check for Information Security Management System (ISMS)
	ismsExists := ce.checkISMS()
	if ismsExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Management System",
			Description: "ISMS not documented",
			Severity:    "high",
			Suggestion:  "Create Information Security Management System",
		})
	}

	// Check for risk assessment
	riskAssessmentExists := ce.checkRiskAssessment()
	if riskAssessmentExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Risk Management",
			Description: "Risk assessment process not documented",
			Severity:    "high",
			Suggestion:  "Create risk assessment and treatment procedures",
		})
	}

	// Check for asset management
	assetManagementExists := ce.checkAssetManagement()
	if assetManagementExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Asset Management",
			Description: "Asset management procedures not documented",
			Severity:    "medium",
			Suggestion:  "Create asset management procedures",
		})
	}

	// Check for human resource security
	hrSecurityExists := ce.checkHRSecurity()
	if hrSecurityExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Human Resources",
			Description: "HR security procedures not documented",
			Severity:    "medium",
			Suggestion:  "Create HR security procedures",
		})
	}

	// Check for business continuity
	businessContinuityExists := ce.checkBusinessContinuity()
	if businessContinuityExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Business Continuity",
			Description: "Business continuity plan not documented",
			Severity:    "medium",
			Suggestion:  "Create business continuity and disaster recovery plan",
		})
	}

	// Check for physical security
	physicalSecurityExists := ce.checkPhysicalSecurity()
	if physicalSecurityExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Physical Security",
			Description: "Physical security procedures not documented",
			Severity:    "low",
			Suggestion:  "Document physical security procedures",
		})
	}

	result.TotalChecks = 6
	result.Score = float64(result.PassedChecks) / float64(result.TotalChecks) * 100

	return mo.Ok(result)
}

// checkNISTCompliance checks NIST Cybersecurity Framework compliance
func (ce *ComplianceEngine) checkNISTCompliance(ctx context.Context) mo.Result[ComplianceResult] {
	result := ComplianceResult{
		Framework:   FrameworkNIST,
		LastChecked: time.Now(),
	}

	// Check for Identify function
	identifyExists := ce.checkNISTIdentify()
	if identifyExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Identify",
			Description: "Asset identification procedures not documented",
			Severity:    "high",
			Suggestion:  "Create asset identification and management procedures",
		})
	}

	// Check for Protect function
	protectExists := ce.checkNISTProtect()
	if protectExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Protect",
			Description: "Protective controls not implemented",
			Severity:    "high",
			Suggestion:  "Implement protective security controls",
		})
	}

	// Check for Detect function
	detectExists := ce.checkNISTDetect()
	if detectExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Detect",
			Description: "Detection capabilities not implemented",
			Severity:    "high",
			Suggestion:  "Implement security monitoring and detection",
		})
	}

	// Check for Respond function
	respondExists := ce.checkNISTRespond()
	if respondExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Respond",
			Description: "Response procedures not documented",
			Severity:    "high",
			Suggestion:  "Create incident response procedures",
		})
	}

	// Check for Recover function
	recoverExists := ce.checkNISTRecover()
	if recoverExists {
		result.PassedChecks++
	} else {
		result.FailedChecks = append(result.FailedChecks, ComplianceIssue{
			Category:    "Recover",
			Description: "Recovery procedures not documented",
			Severity:    "medium",
			Suggestion:  "Create recovery and restoration procedures",
		})
	}

	result.TotalChecks = 5
	result.Score = float64(result.PassedChecks) / float64(result.TotalChecks) * 100

	return mo.Ok(result)
}

// Helper methods for compliance checks

func (ce *ComplianceEngine) checkFileExists(filename string) bool {
	exists, _ := ce.fileSystem.Exists(filename)
	return exists
}

func (ce *ComplianceEngine) checkDPOConfigured() bool {
	// Check for DPO in security policy
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "data protection") ||
		strings.Contains(strings.ToLower(securityPolicy), "dpo")
}

func (ce *ComplianceEngine) checkDataBreachProcedures() bool {
	incidentResponseResult := ce.fileOps.ReadFile("incident-response.md")
	if incidentResponseResult.IsError() {
		return false
	}
	incidentResponse, _ := incidentResponseResult.Get()
	return strings.Contains(strings.ToLower(incidentResponse), "data breach") ||
		strings.Contains(strings.ToLower(incidentResponse), "breach notification")
}

func (ce *ComplianceEngine) checkDataRetentionPolicy() bool {
	privacyPolicyResult := ce.fileOps.ReadFile("privacy-policy.md")
	if privacyPolicyResult.IsError() {
		return false
	}
	privacyPolicy, _ := privacyPolicyResult.Get()
	return strings.Contains(strings.ToLower(privacyPolicy), "retention") ||
		strings.Contains(strings.ToLower(privacyPolicy), "data retention")
}

func (ce *ComplianceEngine) checkConsentManagement() bool {
	privacyPolicyResult := ce.fileOps.ReadFile("privacy-policy.md")
	if privacyPolicyResult.IsError() {
		return false
	}
	privacyPolicy, _ := privacyPolicyResult.Get()
	return strings.Contains(strings.ToLower(privacyPolicy), "consent") ||
		strings.Contains(strings.ToLower(privacyPolicy), "user consent")
}

func (ce *ComplianceEngine) checkAccessControls() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "access control") ||
		strings.Contains(strings.ToLower(securityPolicy), "user access")
}

func (ce *ComplianceEngine) checkMonitoringSystems() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "monitoring") ||
		strings.Contains(strings.ToLower(securityPolicy), "security monitoring")
}

func (ce *ComplianceEngine) checkChangeManagement() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "change management") ||
		strings.Contains(strings.ToLower(securityPolicy), "change control")
}

func (ce *ComplianceEngine) checkISMS() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "isms") ||
		strings.Contains(strings.ToLower(securityPolicy), "information security management")
}

func (ce *ComplianceEngine) checkRiskAssessment() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "risk assessment") ||
		strings.Contains(strings.ToLower(securityPolicy), "risk management")
}

func (ce *ComplianceEngine) checkAssetManagement() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "asset management") ||
		strings.Contains(strings.ToLower(securityPolicy), "asset inventory")
}

func (ce *ComplianceEngine) checkHRSecurity() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "human resources") ||
		strings.Contains(strings.ToLower(securityPolicy), "hr security")
}

func (ce *ComplianceEngine) checkBusinessContinuity() bool {
	incidentResponseResult := ce.fileOps.ReadFile("incident-response.md")
	if incidentResponseResult.IsError() {
		return false
	}
	incidentResponse, _ := incidentResponseResult.Get()
	return strings.Contains(strings.ToLower(incidentResponse), "business continuity") ||
		strings.Contains(strings.ToLower(incidentResponse), "disaster recovery")
}

func (ce *ComplianceEngine) checkPhysicalSecurity() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "physical security") ||
		strings.Contains(strings.ToLower(securityPolicy), "access control")
}

func (ce *ComplianceEngine) checkNISTIdentify() bool {
	return ce.checkFileExists("asset-inventory.md") || ce.checkFileExists("risk-assessment.md")
}

func (ce *ComplianceEngine) checkNISTProtect() bool {
	return ce.checkFileExists("security-policy.md") || ce.checkFileExists("access-control.md")
}

func (ce *ComplianceEngine) checkNISTDetect() bool {
	securityPolicyResult := ce.fileOps.ReadFile("security-policy.md")
	if securityPolicyResult.IsError() {
		return false
	}
	securityPolicy, _ := securityPolicyResult.Get()
	return strings.Contains(strings.ToLower(securityPolicy), "monitoring") ||
		strings.Contains(strings.ToLower(securityPolicy), "detection")
}

func (ce *ComplianceEngine) checkNISTRespond() bool {
	return ce.checkFileExists("incident-response.md")
}

func (ce *ComplianceEngine) checkNISTRecover() bool {
	incidentResponseResult := ce.fileOps.ReadFile("incident-response.md")
	if incidentResponseResult.IsError() {
		return false
	}
	incidentResponse, _ := incidentResponseResult.Get()
	return strings.Contains(strings.ToLower(incidentResponse), "recovery") ||
		strings.Contains(strings.ToLower(incidentResponse), "restoration")
}

// GenerateComplianceReport generates comprehensive compliance report
func (ce *ComplianceEngine) GenerateComplianceReport(ctx context.Context) mo.Result[ComplianceDashboard] {
	frameworks := []ComplianceFramework{
		FrameworkGDPR,
		FrameworkSOC2,
		FrameworkISO27001,
		FrameworkNIST,
	}

	dashboard := ComplianceDashboard{
		GeneratedAt:  time.Now(),
		Frameworks:   make([]FrameworkResult, 0),
		OverallScore: 0,
		TotalIssues:  0,
	}

	totalScore := 0.0
	for _, framework := range frameworks {
		result := ce.CheckCompliance(ctx, framework)
		if result.IsError() {
			continue
		}

		complianceResult, _ := result.Get()
		frameworkResult := FrameworkResult{
			Name:   string(framework),
			Score:  complianceResult.Score,
			Issues: len(complianceResult.FailedChecks),
			Status: ce.getComplianceStatus(complianceResult.Score),
		}

		dashboard.Frameworks = append(dashboard.Frameworks, frameworkResult)
		totalScore += complianceResult.Score
		dashboard.TotalIssues += len(complianceResult.FailedChecks)
	}

	if len(frameworks) > 0 {
		dashboard.OverallScore = totalScore / float64(len(frameworks))
	}

	return mo.Ok(dashboard)
}

// ComplianceDashboard represents overall compliance dashboard
type ComplianceDashboard struct {
	GeneratedAt  time.Time         `json:"generated_at"`
	Frameworks   []FrameworkResult `json:"frameworks"`
	OverallScore float64           `json:"overall_score"`
	TotalIssues  int               `json:"total_issues"`
}

// FrameworkResult represents framework-specific results
type FrameworkResult struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Issues int     `json:"issues"`
	Status string  `json:"status"`
}

func (ce *ComplianceEngine) getComplianceStatus(score float64) string {
	if score >= 90 {
		return "Excellent"
	} else if score >= 80 {
		return "Good"
	} else if score >= 70 {
		return "Fair"
	} else if score >= 60 {
		return "Poor"
	} else {
		return "Critical"
	}
}
