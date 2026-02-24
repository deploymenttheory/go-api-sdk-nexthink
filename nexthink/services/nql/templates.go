package nql

import "fmt"

// Query templates provide pre-built queries for common use cases
// Simplifies query construction for frequently used patterns

// =============================================================================
// Template Type
// =============================================================================

// Template represents a pre-built NQL query template
type Template struct {
	query string
	qb    *QueryBuilder
}

// Query returns the NQL query string (for creating in Nexthink admin)
func (t *Template) Query() string {
	return t.query
}

// QueryBuilder returns the QueryBuilder used to create this template
func (t *Template) QueryBuilder() *QueryBuilder {
	return t.qb
}

// ToRequest converts the template to an ExecuteRequest with the given QueryID
// The QueryID must be the ID of the query saved in Nexthink admin
func (t *Template) ToRequest(queryID string) *ExecuteRequest {
	return &ExecuteRequest{
		QueryID: queryID,
	}
}

// ToRequestWithPlatform converts to ExecuteRequest with platform parameter
func (t *Template) ToRequestWithPlatform(queryID, platform string) *ExecuteRequest {
	return &ExecuteRequest{
		QueryID:  queryID,
		Platform: platform,
	}
}

// newTemplate creates a new template from a QueryBuilder
func newTemplate(qb *QueryBuilder) *Template {
	return &Template{
		query: qb.Build(),
		qb:    qb,
	}
}

// =============================================================================
// Templates
// =============================================================================

// Templates provides access to pre-built query templates
type Templates struct{}

// NewTemplates creates a new templates instance
func NewTemplates() *Templates {
	return &Templates{}
}

// =============================================================================
// Device Queries
// =============================================================================

// DevicesWithCrashes returns devices with application crashes
// period: e.g., "during past 7d"
// binaryName: e.g., "outlook.exe" (supports wildcards like "outlook*")
func (t *Templates) DevicesWithCrashes(period, binaryName string) *Template {
	qb := NewQueryBuilder().
		FromDevices().
		With(fmt.Sprintf("execution.crashes %s", period)).
		ComputeSum("total_crashes", "number_of_crashes")
	
	if binaryName != "" {
		qb.WhereEquals(FieldBinaryName, binaryName)
	}
	
	qb.SortDesc("total_crashes")
	
	return newTemplate(qb)
}

// DevicesWithHighMemoryUsage returns devices with memory usage above threshold
// threshold: memory threshold in GB (e.g., 90)
// period: e.g., "during past 7d"
func (t *Templates) DevicesWithHighMemoryUsage(threshold int, period string) *Template {
	qb := NewQueryBuilder().
		FromDevices().
		During(period).
		Include(fmt.Sprintf("device_performance.events %s", period)).
		Compute("memory_usage_ratio", "event.system_drive_usage.avg()/event.system_drive_capacity.avg()*100").
		WhereGreaterEqual("memory_usage_ratio", fmt.Sprintf("%d", threshold)).
		List("device.name", "memory_usage_ratio").
		SortDesc("memory_usage_ratio")
	
	return newTemplate(qb)
}

// DevicesByPlatform returns device counts grouped by OS platform
func (t *Templates) DevicesByPlatform(period string) *Template {
	qb := NewQueryBuilder().
		FromDevices().
		During(period).
		SummarizeCount("device_count").
		SummarizeBy(FieldOSPlatform).
		SortDesc("device_count")
	
	return newTemplate(qb)
}

// DevicesWithSlowBootTime returns devices with boot time above threshold
// threshold: boot time threshold in seconds (e.g., 60)
// period: e.g., "during past 7d"
func (t *Templates) DevicesWithSlowBootTime(threshold int, period string) *Template {
	qb := NewQueryBuilder().
		FromDevices().
		With(fmt.Sprintf("session.logins %s", period)).
		ComputeAvg("avg_boot_time", FieldTimeUntilDesktopIsVisible).
		WhereGreaterEqual("avg_boot_time", fmt.Sprintf("%ds", threshold)).
		List("device.name", "avg_boot_time").
		SortDesc("avg_boot_time")
	
	return newTemplate(qb)
}

// =============================================================================
// User Queries
// =============================================================================

// UsersWithWebErrors returns users experiencing web errors
// period: e.g., "during past 7d"
// appName: application name (optional, empty string for all apps)
func (t *Templates) UsersWithWebErrors(period, appName string) *Template {
	qb := NewQueryBuilder().
		FromUsers().
		With(fmt.Sprintf("web.errors %s", period)).
		ComputeSum("total_errors", FieldNumberOfErrors)

	if appName != "" {
		qb.WhereEquals(FieldApplicationName, appName)
	}

	qb.List("user.name", "total_errors").
		SortDesc("total_errors")
	
	return newTemplate(qb)
}

// UsersWithPoorCollaborationQuality returns users with poor collaboration quality
func (t *Templates) UsersWithPoorCollaborationQuality(period string) *Template {
	qb := NewQueryBuilder().
		FromUsers().
		With(fmt.Sprintf("collaboration.sessions %s", period)).
		Where("session.audio.quality == poor or session.video.quality == poor").
		ComputeCount("poor_sessions").
		List("user.name", "poor_sessions").
		SortDesc("poor_sessions")
	
	return newTemplate(qb)
}

// =============================================================================
// Application Queries
// =============================================================================

// ApplicationsWithHighErrorRate returns applications with high web error rates
// minPageViews: minimum page views threshold to avoid false positives
// period: e.g., "during past 60min"
func (t *Templates) ApplicationsWithHighErrorRate(minPageViews int, period string) *Template {
	qb := NewQueryBuilder().
		From(TableApplications).
		With(fmt.Sprintf("web.page_views %s", period)).
		Where("is_soft_navigation = false").
		ComputeSum("total_page_views", "number_of_page_views").
		With(fmt.Sprintf("web.errors %s", period)).
		ComputeSum("error_count", "error.number_of_errors").
		Summarize("error_ratio", "error_count.sum() * 100 / total_page_views.sum()").
		SummarizeBy(FieldApplicationName).
		WhereGreater("total_page_views", fmt.Sprintf("%d", minPageViews)).
		SortDesc("error_ratio")
	
	return newTemplate(qb)
}

// TopCrashingApplications returns applications with the most crashes
func (t *Templates) TopCrashingApplications(period string, limit int) *Template {
	qb := NewQueryBuilder().
		From(TableExecutionCrashes).
		During(period).
		SummarizeCount("crash_count").
		Summarize("device_count", "device.count()").
		SummarizeBy(FieldBinaryName).
		SortDesc("crash_count").
		Limit(limit)
	
	return newTemplate(qb)
}

// =============================================================================
// Performance Queries
// =============================================================================

// WebPageLoadPerformance returns average page load times by application
func (t *Templates) WebPageLoadPerformance(appName, period string) *Template {
	qb := NewQueryBuilder().
		From(TableWebPageViews).
		During(period)

	if appName != "" {
		qb.WhereEquals(FieldApplicationName, appName)
	}

	qb.Summarize("backend_time", "page_load_time.backend.avg()").
		Summarize("network_time", "page_load_time.network.avg()").
		Summarize("client_time", "page_load_time.client.avg()").
		SummarizeBy(FieldApplicationName).
		SortDesc("backend_time")
	
	return newTemplate(qb)
}

// NetworkConnectivityIssues returns devices with connectivity issues
func (t *Templates) NetworkConnectivityIssues(period string) *Template {
	qb := NewQueryBuilder().
		From(TableConnectivityEvents).
		During(period).
		Where("wifi.signal_strength.avg <= -67 or wifi.noise_level.avg >= -80").
		SummarizeAvg("avg_signal", "wifi.signal_strength").
		SummarizeAvg("avg_noise", "wifi.noise_level").
		SummarizeBy("device.name").
		SortAsc("avg_signal")
	
	return newTemplate(qb)
}

// =============================================================================
// DEX Score Queries
// =============================================================================

// OverallDEXScore returns the overall DEX score for a population
func (t *Templates) OverallDEXScore(period string) *Template {
	qb := NewQueryBuilder().
		FromUsers().
		Include(fmt.Sprintf("dex.scores %s", period)).
		ComputeAvg("dex_per_user", "value").
		Where("dex_per_user != NULL").
		SummarizeAvg("overall_dex", "dex_per_user")
	
	return newTemplate(qb)
}

// DEXScoreByPlatform returns DEX scores grouped by OS platform
func (t *Templates) DEXScoreByPlatform(period string) *Template {
	qb := NewQueryBuilder().
		FromDevices().
		Include(fmt.Sprintf("dex.scores %s", period)).
		ComputeAvg("dex_per_device", "value").
		Where("dex_per_device != NULL").
		SummarizeAvg("dex_score", "dex_per_device").
		SummarizeBy(FieldOSPlatform).
		SortDesc("dex_score")
	
	return newTemplate(qb)
}

// UsersWithLowDEXScore returns users with DEX scores below threshold
// threshold: DEX score threshold (e.g., 50)
func (t *Templates) UsersWithLowDEXScore(threshold int, period string) *Template {
	qb := NewQueryBuilder().
		FromUsers().
		With(fmt.Sprintf("dex.scores %s", period)).
		ComputeAvg("user_dex", "value").
		WhereLess("user_dex", fmt.Sprintf("%d", threshold)).
		List("user.name", "user_dex").
		SortAsc("user_dex")
	
	return newTemplate(qb)
}

// DEXScoreImpactByComponent returns DEX score impact by component
func (t *Templates) DEXScoreImpactByComponent(component, period string) *Template {
	impactField := "endpoint.logon_speed_score_impact" // default

	switch component {
	case "boot_speed":
		impactField = "endpoint.boot_speed_score_impact"
	case "software_reliability":
		impactField = "endpoint.software_reliability_score_impact"
	case "virtual_session_lag":
		impactField = "endpoint.virtual_session_lag_score_impact"
	}

	qb := NewQueryBuilder().
		FromUsers().
		Include(fmt.Sprintf("dex.scores %s", period)).
		ComputeAvg("impact_per_user", impactField).
		ComputeAvg("dex_per_user", "value").
		Where("dex_per_user != NULL").
		Summarize("total_impact", "(impact_per_user.avg()*countif(impact_per_user != NULL))/countif(dex_per_user != NULL)")
	
	return newTemplate(qb)
}

// =============================================================================
// Monitoring / Alert Queries
// =============================================================================

// DevicesWithSystemCrashes returns devices with system crashes above threshold
// threshold: number of crashes (e.g., 3)
// period: e.g., "during past 7d"
func (t *Templates) DevicesWithSystemCrashes(threshold int, period string) *Template {
	qb := NewQueryBuilder().
		FromDevices().
		With(fmt.Sprintf("device_performance.system_crashes %s", period)).
		ComputeSum("crash_count", "number_of_system_crashes").
		WhereGreaterEqual("crash_count", fmt.Sprintf("%d", threshold)).
		List("device.name", "crash_count").
		SortDesc("crash_count")
	
	return newTemplate(qb)
}

// BinariesWithHighCrashRate returns binaries with crash counts above threshold
// threshold: minimum crash count
func (t *Templates) BinariesWithHighCrashRate(threshold int, period string) *Template {
	qb := NewQueryBuilder().
		From(TableExecutionCrashes).
		During(period).
		SummarizeCount("total_crashes").
		Summarize("devices_affected", "device.count()").
		SummarizeBy(FieldBinaryName).
		WhereGreaterEqual("total_crashes", fmt.Sprintf("%d", threshold)).
		SortDesc("total_crashes")
	
	return newTemplate(qb)
}

// =============================================================================
// Workflow / Remote Action Queries
// =============================================================================

// WorkflowExecutionSuccess returns successful workflow executions
func (t *Templates) WorkflowExecutionSuccess(period string) *Template {
	qb := NewQueryBuilder().
		From(TableWorkflowExecutions).
		During(period).
		WhereEquals("status", ExecutionStatusSuccess).
		SummarizeSum("executions", "number_of_executions").
		SummarizeBy("workflow.name").
		SortDesc("executions")
	
	return newTemplate(qb)
}

// RemoteActionSavingsEstimate returns estimated cost savings from remote actions
// costPerExecution: estimated cost saved per successful execution
func (t *Templates) RemoteActionSavingsEstimate(costPerExecution int, period string) *Template {
	qb := NewQueryBuilder().
		From(TableRemoteActionExecutions).
		During(period).
		WhereEquals("status", ExecutionStatusSuccess).
		WhereEquals("purpose", PurposeRemediation).
		Summarize("amt_saved", fmt.Sprintf("(number_of_executions.sum()) * (%d)", costPerExecution)).
		SummarizeBy("remote_action.name").
		SortDesc("amt_saved")
	
	return newTemplate(qb)
}

// =============================================================================
// Helper Methods
// =============================================================================

// GetAllTemplates returns a list of all available template names
func (t *Templates) GetAllTemplates() []string {
	return []string{
		"DevicesWithCrashes",
		"DevicesWithHighMemoryUsage",
		"DevicesByPlatform",
		"DevicesWithSlowBootTime",
		"UsersWithWebErrors",
		"UsersWithPoorCollaborationQuality",
		"ApplicationsWithHighErrorRate",
		"TopCrashingApplications",
		"WebPageLoadPerformance",
		"NetworkConnectivityIssues",
		"OverallDEXScore",
		"DEXScoreByPlatform",
		"UsersWithLowDEXScore",
		"DEXScoreImpactByComponent",
		"DevicesWithSystemCrashes",
		"BinariesWithHighCrashRate",
		"WorkflowExecutionSuccess",
		"RemoteActionSavingsEstimate",
	}
}
