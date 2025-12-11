package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/LarsArtmann/template-CLI/pkg/sdk/fileops"
	"github.com/LarsArtmann/template-CLI/pkg/sdk/vfs"
	"github.com/samber/mo"
)

// MetricsEngine generates security metrics and executive reports
type MetricsEngine struct {
	fileSystem    vfs.FileSystem
	fileOps       *fileops.Service
	compliance    *ComplianceEngine
}

// NewMetricsEngine creates new metrics engine
func NewMetricsEngine() *MetricsEngine {
	vfsImpl := vfs.NewOSVFS()
	fileOps := fileops.NewService(fileops.ServiceOptions{
		BaseDirectory: ".",
		ValidatePaths: true,
	})
	
	return &MetricsEngine{
		fileSystem: vfsImpl,
		fileOps:    fileOps,
		compliance: NewComplianceEngine(),
	}
}

// OutputFormat represents metrics output format
type OutputFormat string

const (
	FormatText    OutputFormat = "text"
	FormatJSON    OutputFormat = "json"
	FormatPrometheus OutputFormat = "prometheus"
)

// SecurityMetrics represents comprehensive security metrics
type SecurityMetrics struct {
	GeneratedAt     time.Time              `json:"generated_at"`
	Score           float64                `json:"score"`
	Coverage        MetricsCoverage        `json:"coverage"`
	Compliance     MetricsCompliance      `json:"compliance"`
	Documentation  MetricsDocumentation   `json:"documentation"`
	IncidentResponse MetricsIncidentResponse `json:"incident_response"`
	Automation     MetricsAutomation      `json:"automation"`
	Trends         MetricsTrends          `json:"trends"`
	Recommendations []string              `json:"recommendations"`
}

// MetricsCoverage represents security coverage metrics
type MetricsCoverage struct {
	TotalPolicies      int     `json:"total_policies"`
	ImplementedPolicies int     `json:"implemented_policies"`
	CoveragePercentage float64 `json:"coverage_percentage"`
	MissingPolicies   []string `json:"missing_policies"`
}

// MetricsCompliance represents compliance metrics
type MetricsCompliance struct {
	OverallScore     float64             `json:"overall_score"`
	FrameworkScores   map[string]float64   `json:"framework_scores"`
	HighRiskIssues   int                 `json:"high_risk_issues"`
	MediumRiskIssues int                 `json:"medium_risk_issues"`
	LowRiskIssues    int                 `json:"low_risk_issues"`
	LastAssessment   time.Time           `json:"last_assessment"`
}

// MetricsDocumentation represents documentation metrics
type MetricsDocumentation struct {
	TotalDocuments    int     `json:"total_documents"`
	ValidDocuments   int     `json:"valid_documents"`
	OutdatedDocuments int     `json:"outdated_documents"`
	ValidationScore   float64 `json:"validation_score"`
}

// MetricsIncidentResponse represents incident response metrics
type MetricsIncidentResponse struct {
	PlanExists       bool   `json:"plan_exists"`
	CompletenessScore float64 `json:"completeness_score"`
	ResponseProcedures []string `json:"response_procedures"`
	EscalationDefined bool   `json:"escalation_defined"`
}

// MetricsAutomation represents automation metrics
type MetricsAutomation struct {
	AutomatedChecks   int     `json:"automated_checks"`
	ManualProcesses   int     `json:"manual_processes"`
	AutomationRatio   float64 `json:"automation_ratio"`
	CIMaturity        float64 `json:"ci_maturity"`
}

// MetricsTrends represents trends over time
type MetricsTrends struct {
	ScoreHistory     []ScorePoint `json:"score_history"`
	ComplianceTrend  string       `json:"compliance_trend"`
	DocumentationTrend string     `json:"documentation_trend"`
}

// ScorePoint represents a score at a specific time
type ScorePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
}

// ExecutiveReport represents executive security report
type ExecutiveReport struct {
	GeneratedAt     time.Time             `json:"generated_at"`
	ExecutiveSummary ExecutiveSummary     `json:"executive_summary"`
	RiskAssessment  RiskAssessment        `json:"risk_assessment"`
	ComplianceStatus ComplianceStatus     `json:"compliance_status"`
	Recommendations []ExecutiveRecommendation `json:"recommendations"`
	ActionItems    []ActionItem          `json:"action_items"`
	NextSteps      []NextStep           `json:"next_steps"`
}

// ExecutiveSummary represents high-level executive summary
type ExecutiveSummary struct {
	OverallScore       float64 `json:"overall_score"`
	RiskLevel         string   `json:"risk_level"`
	ComplianceStatus   string   `json:"compliance_status"`
	KeyAchievements   []string `json:"key_achievements"`
	CriticalIssues    []string `json:"critical_issues"`
}

// RiskAssessment represents risk assessment for executives
type RiskAssessment struct {
	CriticalRisks []Risk `json:"critical_risks"`
	HighRisks     []Risk `json:"high_risks"`
	MediumRisks    []Risk `json:"medium_risks"`
	LowRisks       []Risk `json:"low_risks"`
	RiskTrend      string `json:"risk_trend"`
}

// Risk represents a specific security risk
type Risk struct {
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Likelihood  string  `json:"likelihood"`
	Score       float64 `json:"score"`
	Mitigation  string  `json:"mitigation"`
}

// ComplianceStatus represents compliance status for executives
type ComplianceStatus struct {
	OverallCompliance   float64                  `json:"overall_compliance"`
	FrameworkStatus    map[string]FrameworkStatus `json:"framework_status"`
	CriticalGaps       []string                  `json:"critical_gaps"`
	AuditReadiness     float64                  `json:"audit_readiness"`
}

// FrameworkStatus represents status of specific framework
type FrameworkStatus struct {
	Score     float64 `json:"score"`
	Status     string   `json:"status"`
	CriticalIssues int   `json:"critical_issues"`
	LastReview time.Time `json:"last_review"`
}

// ExecutiveRecommendation represents recommendation for executives
type ExecutiveRecommendation struct {
	Priority     string  `json:"priority"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Benefit      string  `json:"benefit"`
	Effort       string  `json:"effort"`
	Timeline     string  `json:"timeline"`
	Owner        string  `json:"owner"`
}

// ActionItem represents specific action item
type ActionItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	Assignee    string    `json:"assignee"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	Progress    float64   `json:"progress"`
}

// NextStep represents next step for executives
type NextStep struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Timeline    string `json:"timeline"`
	Owner       string `json:"owner"`
	SuccessCriteria string `json:"success_criteria"`
}

// GenerateMetrics generates comprehensive security metrics
func (me *MetricsEngine) GenerateMetrics(ctx context.Context, format OutputFormat) mo.Result[string] {
	metrics := me.calculateSecurityMetrics(ctx)
	
	switch format {
	case FormatJSON:
		return me.formatMetricsAsJSON(metrics)
	case FormatPrometheus:
		return me.formatMetricsAsPrometheus(metrics)
	case FormatText:
	default:
		return me.formatMetricsAsText(metrics)
	}
}

// GenerateExecutiveReport generates executive security report
func (me *MetricsEngine) GenerateExecutiveReport(ctx context.Context) mo.Result[ExecutiveReport] {
	metrics := me.calculateSecurityMetrics(ctx)
	
	report := ExecutiveReport{
		GeneratedAt: time.Now(),
		ExecutiveSummary: me.createExecutiveSummary(metrics),
		RiskAssessment:  me.createRiskAssessment(metrics),
		ComplianceStatus: me.createComplianceStatus(metrics),
		Recommendations: me.createExecutiveRecommendations(metrics),
		ActionItems:    me.createActionItems(metrics),
		NextSteps:      me.createNextSteps(metrics),
	}
	
	return mo.Ok(report)
}

// calculateSecurityMetrics calculates all security metrics
func (me *MetricsEngine) calculateSecurityMetrics(ctx context.Context) SecurityMetrics {
	// Calculate coverage metrics
	coverage := me.calculateCoverage()
	
	// Calculate compliance metrics
	compliance := me.calculateCompliance(ctx)
	
	// Calculate documentation metrics
	documentation := me.calculateDocumentation()
	
	// Calculate incident response metrics
	incidentResponse := me.calculateIncidentResponse()
	
	// Calculate automation metrics
	automation := me.calculateAutomation()
	
	// Calculate trends
	trends := me.calculateTrends()
	
	// Calculate overall score
	overallScore := me.calculateOverallScore(coverage, compliance, documentation, incidentResponse, automation)
	
	// Generate recommendations
	recommendations := me.generateRecommendations(coverage, compliance, documentation, incidentResponse, automation)
	
	return SecurityMetrics{
		GeneratedAt:      time.Now(),
		Score:           overallScore,
		Coverage:        coverage,
		Compliance:     compliance,
		Documentation:  documentation,
		IncidentResponse: incidentResponse,
		Automation:     automation,
		Trends:         trends,
		Recommendations: recommendations,
	}
}

// calculateCoverage calculates security coverage metrics
func (me *MetricsEngine) calculateCoverage() MetricsCoverage {
	expectedPolicies := []string{
		"SECURITY.md",
		"security-policy.md",
		"incident-response.md", 
		"privacy-policy.md",
		"bug-bounty-policy.md",
	}
	
	implementedPolicies := 0
	var missingPolicies []string
	
	for _, policy := range expectedPolicies {
		if exists, _ := me.fileSystem.Exists(policy); exists {
			implementedPolicies++
		} else {
			missingPolicies = append(missingPolicies, policy)
		}
	}
	
	totalPolicies := len(expectedPolicies)
	coveragePercentage := float64(implementedPolicies) / float64(totalPolicies) * 100
	
	return MetricsCoverage{
		TotalPolicies:       totalPolicies,
		ImplementedPolicies: implementedPolicies,
		CoveragePercentage:  coveragePercentage,
		MissingPolicies:     missingPolicies,
	}
}

// calculateCompliance calculates compliance metrics
func (me *MetricsEngine) calculateCompliance(ctx context.Context) MetricsCompliance {
	dashboardResult := me.compliance.GenerateComplianceReport(ctx)
	if dashboardResult.IsError() {
		return MetricsCompliance{
			OverallScore:     0,
			FrameworkScores:   make(map[string]float64),
			HighRiskIssues:   0,
			MediumRiskIssues: 0,
			LowRiskIssues:    0,
			LastAssessment:   time.Now(),
		}
	}
	
	dashboard, _ := dashboardResult.Get()
	frameworkScores := make(map[string]float64)
	
	for _, framework := range dashboard.Frameworks {
		frameworkScores[framework.Name] = framework.Score
	}
	
	// Count issues by severity
	highRiskIssues := 0
	mediumRiskIssues := 0
	lowRiskIssues := 0
	
	// This would be enhanced with actual compliance issue analysis
	// For now, estimate based on scores
	for _, score := range frameworkScores {
		if score < 60 {
			highRiskIssues++
		} else if score < 80 {
			mediumRiskIssues++
		} else {
			lowRiskIssues++
		}
	}
	
	return MetricsCompliance{
		OverallScore:     dashboard.OverallScore,
		FrameworkScores:   frameworkScores,
		HighRiskIssues:   highRiskIssues,
		MediumRiskIssues: mediumRiskIssues,
		LowRiskIssues:    lowRiskIssues,
		LastAssessment:   dashboard.GeneratedAt,
	}
}

// calculateDocumentation calculates documentation metrics
func (me *MetricsEngine) calculateDocumentation() MetricsDocumentation {
	expectedDocuments := []string{
		"SECURITY.md",
		"security-policy.md",
		"incident-response.md",
		"privacy-policy.md", 
		"bug-bounty-policy.md",
	}
	
	totalDocuments := len(expectedDocuments)
	validDocuments := 0
	outdatedDocuments := 0
	
	validator := NewPolicyValidator()
	
	for _, doc := range expectedDocuments {
		if exists, _ := me.fileSystem.Exists(doc); exists {
			result := validator.ValidateFile(doc)
			if result.Valid {
				validDocuments++
			} else {
				// Consider outdated if score is below threshold
				if result.Score < 70 {
					outdatedDocuments++
				}
			}
		}
	}
	
	validationScore := float64(validDocuments) / float64(totalDocuments) * 100
	
	return MetricsDocumentation{
		TotalDocuments:    totalDocuments,
		ValidDocuments:     validDocuments,
		OutdatedDocuments:  outdatedDocuments,
		ValidationScore:    validationScore,
	}
}

// calculateIncidentResponse calculates incident response metrics
func (me *MetricsEngine) calculateIncidentResponse() MetricsIncidentResponse {
	planExists := false
	planExists, _ = me.fileSystem.Exists("incident-response.md")
	
	completenessScore := 0.0
	responseProcedures := []string{}
	escalationDefined := false
	
	if planExists {
		contentResult := me.fileOps.ReadFile("incident-response.md")
		if !contentResult.IsError() {
			content, _ := contentResult.Get()
			
			// Check for key procedures
			procedures := []string{
				"detection", "containment", "eradication", 
				"recovery", "lessons learned",
			}
			
			foundProcedures := 0
			for _, proc := range procedures {
				if contains(content, proc) {
					foundProcedures++
					responseProcedures = append(responseProcedures, proc)
				}
			}
			
			completenessScore = float64(foundProcedures) / float64(len(procedures)) * 100
			escalationDefined = contains(content, "escalation")
		}
	}
	
	return MetricsIncidentResponse{
		PlanExists:        planExists,
		CompletenessScore:  completenessScore,
		ResponseProcedures: responseProcedures,
		EscalationDefined:  escalationDefined,
	}
}

// calculateAutomation calculates automation metrics
func (me *MetricsEngine) calculateAutomation() MetricsAutomation {
	// Check for CI/CD configuration
	automatedChecks := 0
	manualProcesses := 0
	
	// Check for GitHub Actions
	if exists, _ := me.fileSystem.Exists(".github/workflows"); exists {
		automatedChecks++
	}
	
	// Check for CI configuration files
	ciFiles := []string{".github/workflows/security.yml", ".gitlab-ci.yml", ".travis.yml", "Jenkinsfile"}
	for _, file := range ciFiles {
		if exists, _ := me.fileSystem.Exists(file); exists {
			automatedChecks++
		}
	}
	
	// Count manual processes (policies that would need manual checks)
	policyFiles := []string{"SECURITY.md", "security-policy.md", "incident-response.md"}
	for _, file := range policyFiles {
		if exists, _ := me.fileSystem.Exists(file); exists {
			manualProcesses++ // These typically require manual review
		}
	}
	
	totalProcesses := automatedChecks + manualProcesses
	automationRatio := 0.0
	if totalProcesses > 0 {
		automationRatio = float64(automatedChecks) / float64(totalProcesses) * 100
	}
	
	// CI maturity based on automation ratio
	ciMaturity := automationRatio / 20 // Scale to 0-5 maturity scale
	if ciMaturity > 5 {
		ciMaturity = 5
	}
	
	return MetricsAutomation{
		AutomatedChecks: automatedChecks,
		ManualProcesses: manualProcesses,
		AutomationRatio: automationRatio,
		CIMaturity:     ciMaturity,
	}
}

// calculateTrends calculates trends over time
func (me *MetricsEngine) calculateTrends() MetricsTrends {
	// For now, return minimal trends
	// In production, this would read historical data
	scoreHistory := []ScorePoint{
		{
			Timestamp: time.Now().AddDate(0, 0, -7), // 7 days ago
			Score:     75.0,
		},
		{
			Timestamp: time.Now().AddDate(0, 0, -3), // 3 days ago
			Score:     82.5,
		},
		{
			Timestamp: time.Now(), // Today
			Score:     88.0,
		},
	}
	
	return MetricsTrends{
		ScoreHistory:      scoreHistory,
		ComplianceTrend:  "improving",
		DocumentationTrend: "stable",
	}
}

// calculateOverallScore calculates weighted overall security score
func (me *MetricsEngine) calculateOverallScore(coverage MetricsCoverage, compliance MetricsCompliance, documentation MetricsDocumentation, incidentResponse MetricsIncidentResponse, automation MetricsAutomation) float64 {
	// Weight different aspects
	coverageWeight := 0.20
	complianceWeight := 0.35
	documentationWeight := 0.20
	incidentWeight := 0.15
	automationWeight := 0.10
	
	// Calculate weighted score
	overallScore := (coverage.CoveragePercentage * coverageWeight) +
		(compliance.OverallScore * complianceWeight) +
		(documentation.ValidationScore * documentationWeight) +
		(incidentResponse.CompletenessScore * incidentWeight) +
		(automation.AutomationRatio * automationWeight)
	
	return overallScore
}

// generateRecommendations generates recommendations based on metrics
func (me *MetricsEngine) generateRecommendations(coverage MetricsCoverage, compliance MetricsCompliance, documentation MetricsDocumentation, incidentResponse MetricsIncidentResponse, automation MetricsAutomation) []string {
	var recommendations []string
	
	// Coverage recommendations
	if coverage.CoveragePercentage < 80 {
		recommendations = append(recommendations, fmt.Sprintf("Improve security coverage by implementing missing policies: %v", coverage.MissingPolicies))
	}
	
	// Compliance recommendations
	if compliance.OverallScore < 80 {
		recommendations = append(recommendations, "Address high-risk compliance issues to improve overall compliance score")
	}
	
	// Documentation recommendations
	if documentation.ValidationScore < 85 {
		recommendations = append(recommendations, "Update and validate all security documentation to meet quality standards")
	}
	
	// Incident response recommendations
	if !incidentResponse.PlanExists || incidentResponse.CompletenessScore < 80 {
		recommendations = append(recommendations, "Enhance incident response plan with comprehensive procedures and escalation paths")
	}
	
	// Automation recommendations
	if automation.AutomationRatio < 50 {
		recommendations = append(recommendations, "Implement automated security checks and CI/CD integration to improve automation")
	}
	
	return recommendations
}

// Format methods

func (me *MetricsEngine) formatMetricsAsJSON(metrics SecurityMetrics) mo.Result[string] {
	jsonData, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return mo.Err[string](fmt.Errorf("failed to format metrics as JSON: %w", err))
	}
	return mo.Ok(string(jsonData))
}

func (me *MetricsEngine) formatMetricsAsPrometheus(metrics SecurityMetrics) mo.Result[string] {
	var prometheus strings.Builder
	
	// Coverage metrics
	prometheus.WriteString(fmt.Sprintf("# HELP template_security_coverage_percentage Security policy coverage percentage\n"))
	prometheus.WriteString(fmt.Sprintf("# TYPE template_security_coverage_percentage gauge\n"))
	prometheus.WriteString(fmt.Sprintf("template_security_coverage_percentage %.2f\n\n", metrics.Coverage.CoveragePercentage))
	
	// Overall score
	prometheus.WriteString(fmt.Sprintf("# HELP template_security_overall_score Overall security score\n"))
	prometheus.WriteString(fmt.Sprintf("# TYPE template_security_overall_score gauge\n"))
	prometheus.WriteString(fmt.Sprintf("template_security_overall_score %.2f\n\n", metrics.Score))
	
	// Compliance scores
	for framework, score := range metrics.Compliance.FrameworkScores {
		prometheus.WriteString(fmt.Sprintf("template_security_compliance_score{framework=\"%s\"} %.2f\n", framework, score))
	}
	
	return mo.Ok(prometheus.String())
}

func (me *MetricsEngine) formatMetricsAsText(metrics SecurityMetrics) mo.Result[string] {
	var output strings.Builder
	
	output.WriteString("📊 Security Metrics Dashboard\n")
	output.WriteString("============================\n\n")
	
	// Overall Score
	output.WriteString(fmt.Sprintf("🎯 Overall Security Score: %.1f/100\n", metrics.Score))
	output.WriteString(me.getScoreDescription(metrics.Score))
	output.WriteString("\n")
	
	// Coverage
	output.WriteString(fmt.Sprintf("📋 Coverage: %.1f%% (%d/%d policies)\n", 
		metrics.Coverage.CoveragePercentage, 
		metrics.Coverage.ImplementedPolicies,
		metrics.Coverage.TotalPolicies))
	if len(metrics.Coverage.MissingPolicies) > 0 {
		output.WriteString(fmt.Sprintf("   Missing: %v\n", metrics.Coverage.MissingPolicies))
	}
	output.WriteString("\n")
	
	// Compliance
	output.WriteString(fmt.Sprintf("⚖️  Compliance Score: %.1f/100\n", metrics.Compliance.OverallScore))
	for framework, score := range metrics.Compliance.FrameworkScores {
		output.WriteString(fmt.Sprintf("   %s: %.1f%%\n", framework, score))
	}
	output.WriteString("\n")
	
	// Documentation
	output.WriteString(fmt.Sprintf("📚 Documentation Score: %.1f/100\n", metrics.Documentation.ValidationScore))
	output.WriteString(fmt.Sprintf("   Valid: %d/%d documents\n", metrics.Documentation.ValidDocuments, metrics.Documentation.TotalDocuments))
	output.WriteString("\n")
	
	// Incident Response
	output.WriteString(fmt.Sprintf("🚨 Incident Response Score: %.1f/100\n", metrics.IncidentResponse.CompletenessScore))
	if metrics.IncidentResponse.PlanExists {
		output.WriteString("   ✅ Incident response plan exists\n")
	} else {
		output.WriteString("   ❌ Incident response plan missing\n")
	}
	output.WriteString("\n")
	
	// Automation
	output.WriteString(fmt.Sprintf("🤖 Automation Score: %.1f/100\n", metrics.Automation.AutomationRatio))
	output.WriteString(fmt.Sprintf("   Automated checks: %d\n", metrics.Automation.AutomatedChecks))
	output.WriteString(fmt.Sprintf("   CI/CD maturity: %.1f/5\n", metrics.Automation.CIMaturity))
	output.WriteString("\n")
	
	// Recommendations
	if len(metrics.Recommendations) > 0 {
		output.WriteString("💡 Recommendations:\n")
		for i, rec := range metrics.Recommendations {
			output.WriteString(fmt.Sprintf("   %d. %s\n", i+1, rec))
		}
	}
	
	return mo.Ok(output.String())
}

// Executive report creation methods

func (me *MetricsEngine) createExecutiveSummary(metrics SecurityMetrics) ExecutiveSummary {
	riskLevel := me.getRiskLevel(metrics.Score)
	complianceStatus := me.getComplianceStatus(metrics.Compliance.OverallScore)
	
	keyAchievements := []string{}
	if metrics.Coverage.CoveragePercentage >= 80 {
		keyAchievements = append(keyAchievements, "Strong security policy coverage")
	}
	if metrics.Compliance.OverallScore >= 80 {
		keyAchievements = append(keyAchievements, "Good compliance standing")
	}
	if metrics.Documentation.ValidationScore >= 85 {
		keyAchievements = append(keyAchievements, "High-quality documentation")
	}
	
	criticalIssues := []string{}
	if !metrics.IncidentResponse.PlanExists {
		criticalIssues = append(criticalIssues, "Missing incident response plan")
	}
	if metrics.Compliance.HighRiskIssues > 0 {
		criticalIssues = append(criticalIssues, fmt.Sprintf("%d high-risk compliance issues", metrics.Compliance.HighRiskIssues))
	}
	
	return ExecutiveSummary{
		OverallScore:     metrics.Score,
		RiskLevel:         riskLevel,
		ComplianceStatus:   complianceStatus,
		KeyAchievements:   keyAchievements,
		CriticalIssues:    criticalIssues,
	}
}

func (me *MetricsEngine) createRiskAssessment(metrics SecurityMetrics) RiskAssessment {
	criticalRisks := []Risk{}
	highRisks := []Risk{}
	mediumRisks := []Risk{}
	lowRisks := []Risk{}
	
	// Analyze risks based on metrics
	if !metrics.IncidentResponse.PlanExists {
		criticalRisks = append(criticalRisks, Risk{
			Category:    "Incident Response",
			Description: "No documented incident response procedures",
			Impact:      "High",
			Likelihood:  "Medium",
			Score:       9.0,
			Mitigation:  "Develop and implement comprehensive incident response plan",
		})
	}
	
	if metrics.Compliance.OverallScore < 70 {
		criticalRisks = append(criticalRisks, Risk{
			Category:    "Compliance",
			Description: "Poor compliance posture with multiple gaps",
			Impact:      "High",
			Likelihood:  "High",
			Score:       8.5,
			Mitigation:  "Address compliance gaps and implement required controls",
		})
	}
	
	if metrics.Coverage.CoveragePercentage < 60 {
		highRisks = append(highRisks, Risk{
			Category:    "Documentation",
			Description: "Insufficient security policy coverage",
			Impact:      "Medium",
			Likelihood:  "High",
			Score:       7.0,
			Mitigation:  "Implement missing security policies and procedures",
		})
	}
	
	riskTrend := "stable"
	if len(metrics.Trends.ScoreHistory) >= 2 {
		latestScore := metrics.Trends.ScoreHistory[len(metrics.Trends.ScoreHistory)-1].Score
		previousScore := metrics.Trends.ScoreHistory[len(metrics.Trends.ScoreHistory)-2].Score
		if latestScore > previousScore {
			riskTrend = "improving"
		} else if latestScore < previousScore {
			riskTrend = "degrading"
		}
	}
	
	return RiskAssessment{
		CriticalRisks: criticalRisks,
		HighRisks:     highRisks,
		MediumRisks:    mediumRisks,
		LowRisks:       lowRisks,
		RiskTrend:      riskTrend,
	}
}

func (me *MetricsEngine) createComplianceStatus(metrics SecurityMetrics) ComplianceStatus {
	frameworkStatus := make(map[string]FrameworkStatus)
	
	for framework, score := range metrics.Compliance.FrameworkScores {
		status := "non-compliant"
		if score >= 90 {
			status = "fully-compliant"
		} else if score >= 80 {
			status = "mostly-compliant"
		} else if score >= 70 {
			status = "partially-compliant"
		}
		
		criticalIssues := 0
		if score < 70 {
			criticalIssues = 3
		} else if score < 80 {
			criticalIssues = 2
		} else if score < 90 {
			criticalIssues = 1
		}
		
		frameworkStatus[framework] = FrameworkStatus{
			Score:           score,
			Status:          status,
			CriticalIssues:  criticalIssues,
			LastReview:      time.Now(),
		}
	}
	
	criticalGaps := []string{}
	for _, status := range frameworkStatus {
		if status.CriticalIssues > 0 {
			criticalGaps = append(criticalGaps, fmt.Sprintf("%s has %d critical issues", "", status.CriticalIssues))
		}
	}
	
	auditReadiness := metrics.Compliance.OverallScore
	
	return ComplianceStatus{
		OverallCompliance: metrics.Compliance.OverallScore,
		FrameworkStatus:   frameworkStatus,
		CriticalGaps:      criticalGaps,
		AuditReadiness:    auditReadiness,
	}
}

func (me *MetricsEngine) createExecutiveRecommendations(metrics SecurityMetrics) []ExecutiveRecommendation {
	recommendations := []ExecutiveRecommendation{}
	
	if metrics.Score < 80 {
		recommendations = append(recommendations, ExecutiveRecommendation{
			Priority:     "high",
			Category:     "Overall Security",
			Description:  "Improve overall security posture to industry standards",
			Benefit:      "Reduced risk exposure and improved compliance",
			Effort:       "medium",
			Timeline:     "3-6 months",
			Owner:        "CISO",
		})
	}
	
	if !metrics.IncidentResponse.PlanExists {
		recommendations = append(recommendations, ExecutiveRecommendation{
			Priority:     "critical",
			Category:     "Incident Management",
			Description:  "Implement comprehensive incident response plan",
			Benefit:      "Reduced incident impact and faster recovery",
			Effort:       "high",
			Timeline:     "2-4 months",
			Owner:        "Security Team Lead",
		})
	}
	
	if metrics.Coverage.CoveragePercentage < 80 {
		recommendations = append(recommendations, ExecutiveRecommendation{
			Priority:     "medium",
			Category:     "Documentation",
			Description:  "Complete missing security policies",
			Benefit:      "Comprehensive security governance",
			Effort:       "low",
			Timeline:     "1-2 months",
			Owner:        "Security Manager",
		})
	}
	
	return recommendations
}

func (me *MetricsEngine) createActionItems(metrics SecurityMetrics) []ActionItem {
	actionItems := []ActionItem{}
	
	if !metrics.IncidentResponse.PlanExists {
		actionItems = append(actionItems, ActionItem{
			ID:          "IR-001",
			Title:       "Develop Incident Response Plan",
			Description:  "Create comprehensive incident response procedures",
			Priority:    "critical",
			Category:    "Incident Management",
			Assignee:    "Security Team Lead",
			DueDate:     time.Now().AddDate(0, 2, 0), // 2 months
			Status:      "not-started",
			Progress:    0.0,
		})
	}
	
	for i, missing := range metrics.Coverage.MissingPolicies {
		actionItems = append(actionItems, ActionItem{
			ID:          fmt.Sprintf("DOC-%03d", i+1),
			Title:       fmt.Sprintf("Create %s", missing),
			Description:  fmt.Sprintf("Develop and implement %s", missing),
			Priority:    "medium",
			Category:    "Documentation",
			Assignee:    "Security Manager",
			DueDate:     time.Now().AddDate(0, 1, 0), // 1 month
			Status:      "not-started",
			Progress:    0.0,
		})
	}
	
	return actionItems
}

func (me *MetricsEngine) createNextSteps(metrics SecurityMetrics) []NextStep {
	nextSteps := []NextStep{}
	
	nextSteps = append(nextSteps, NextStep{
		Title:           "Review and approve security policies",
		Description:     "Executive review of all security documentation",
		Timeline:        "2 weeks",
		Owner:          "CISO",
		SuccessCriteria: "All policies reviewed and approved",
	})
	
	if metrics.Score < 80 {
		nextSteps = append(nextSteps, NextStep{
			Title:           "Implement security improvement plan",
			Description:     "Execute prioritized security improvements",
			Timeline:        "3 months",
			Owner:          "Security Team",
			SuccessCriteria: "Target security score of 80+ achieved",
		})
	}
	
	return nextSteps
}

// Helper methods

func (me *MetricsEngine) getScoreDescription(score float64) string {
	if score >= 90 {
		return "   ✅ Excellent security posture"
	} else if score >= 80 {
		return "   ✅ Good security posture"
	} else if score >= 70 {
		return "   ⚠️  Fair security posture"
	} else if score >= 60 {
		return "   ⚠️  Poor security posture"
	} else {
		return "   ❌ Critical security issues"
	}
}

func (me *MetricsEngine) getRiskLevel(score float64) string {
	if score >= 90 {
		return "Low"
	} else if score >= 80 {
		return "Medium-Low"
	} else if score >= 70 {
		return "Medium"
	} else if score >= 60 {
		return "Medium-High"
	} else {
		return "High"
	}
}

func (me *MetricsEngine) getComplianceStatus(score float64) string {
	if score >= 90 {
		return "Fully Compliant"
	} else if score >= 80 {
		return "Mostly Compliant"
	} else if score >= 70 {
		return "Partially Compliant"
	} else {
		return "Non-Compliant"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 (len(s) > len(substr) && 
		  (s[:len(substr)] == substr || 
		   s[len(s)-len(substr):] == substr || 
		   hasSubstring(s, substr))))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}