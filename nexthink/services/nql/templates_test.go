package nql

import (
	"strings"
	"testing"
)

func TestTemplates_DevicesWithCrashes(t *testing.T) {
	templates := NewTemplates()

	template := templates.DevicesWithCrashes("during past 7d", "outlook.exe")

	if template == nil {
		t.Fatal("Template is nil")
	}

	query := template.Query()

	// Verify query contains expected components
	if !strings.Contains(query, "devices") {
		t.Errorf("Query missing devices table: %s", query)
	}

	if !strings.Contains(query, "execution.crashes") {
		t.Errorf("Query missing crashes clause: %s", query)
	}

	if !strings.Contains(query, "outlook.exe") {
		t.Errorf("Query missing binary filter: %s", query)
	}

	if !strings.Contains(query, "total_crashes") {
		t.Errorf("Query missing computed field: %s", query)
	}
}

func TestTemplates_DevicesByPlatform(t *testing.T) {
	templates := NewTemplates()

	template := templates.DevicesByPlatform("during past 7d")

	query := template.Query()

	if !strings.Contains(query, "summarize device_count = count()") {
		t.Errorf("Query missing summarize: %s", query)
	}

	if !strings.Contains(query, "operating_system.platform") {
		t.Errorf("Query missing group by platform: %s", query)
	}
}

func TestTemplates_UsersWithWebErrors(t *testing.T) {
	templates := NewTemplates()

	// Test with app name
	template1 := templates.UsersWithWebErrors("during past 7d", "Salesforce")
	query1 := template1.Query()

	if !strings.Contains(query1, "Salesforce") {
		t.Errorf("Query missing app filter: %s", query1)
	}

	// Test without app name
	template2 := templates.UsersWithWebErrors("during past 7d", "")
	query2 := template2.Query()

	if strings.Contains(query2, "application.name") {
		t.Errorf("Query should not have app filter: %s", query2)
	}
}

func TestTemplates_DEXScoreByPlatform(t *testing.T) {
	templates := NewTemplates()

	template := templates.DEXScoreByPlatform("during past 24h")
	query := template.Query()

	if !strings.Contains(query, "dex.scores") {
		t.Errorf("Query missing dex.scores: %s", query)
	}

	if !strings.Contains(query, "operating_system.platform") {
		t.Errorf("Query missing platform grouping: %s", query)
	}
}

func TestTemplates_ToRequest(t *testing.T) {
	templates := NewTemplates()

	template := templates.DevicesWithCrashes("during past 7d", "")

	// Convert to request
	req := template.ToRequest("#test_query")

	if req.QueryID != "#test_query" {
		t.Errorf("Expected QueryID #test_query, got: %s", req.QueryID)
	}

	if req.Platform != "" {
		t.Errorf("Expected empty platform, got: %s", req.Platform)
	}
}

func TestTemplates_ToRequestWithPlatform(t *testing.T) {
	templates := NewTemplates()

	template := templates.DevicesByPlatform("during past 7d")

	// Convert to request with platform
	req := template.ToRequestWithPlatform("#test_query", "Windows")

	if req.QueryID != "#test_query" {
		t.Errorf("Expected QueryID #test_query, got: %s", req.QueryID)
	}

	if req.Platform != "Windows" {
		t.Errorf("Expected platform Windows, got: %s", req.Platform)
	}
}

func TestTemplates_QueryBuilder(t *testing.T) {
	templates := NewTemplates()

	template := templates.DevicesWithCrashes("during past 7d", "chrome.exe")

	// Access underlying query builder
	qb := template.QueryBuilder()

	if qb == nil {
		t.Fatal("QueryBuilder is nil")
	}

	// Should be able to build same query
	query := qb.Build()
	if query != template.Query() {
		t.Error("QueryBuilder.Build() doesn't match Template.Query()")
	}
}

func TestTemplates_AllTemplates(t *testing.T) {
	templates := NewTemplates()

	allTemplates := templates.GetAllTemplates()

	if len(allTemplates) < 15 {
		t.Errorf("Expected at least 15 templates, got: %d", len(allTemplates))
	}

	// Verify some expected templates exist
	expectedTemplates := []string{
		"DevicesWithCrashes",
		"DevicesByPlatform",
		"UsersWithWebErrors",
		"DEXScoreByPlatform",
		"OverallDEXScore",
	}

	for _, expected := range expectedTemplates {
		found := false
		for _, tmpl := range allTemplates {
			if tmpl == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template not found: %s", expected)
		}
	}
}

func TestTemplates_DevicesWithHighMemoryUsage(t *testing.T) {
	templates := NewTemplates()

	template := templates.DevicesWithHighMemoryUsage(90, "during past 7d")
	query := template.Query()

	if !strings.Contains(query, "memory_usage_ratio") {
		t.Errorf("Query missing memory ratio computation: %s", query)
	}

	if !strings.Contains(query, ">= 90") {
		t.Errorf("Query missing threshold filter: %s", query)
	}
}

func TestTemplates_TopCrashingApplications(t *testing.T) {
	templates := NewTemplates()

	template := templates.TopCrashingApplications("during past 7d", 10)
	query := template.Query()

	if !strings.Contains(query, "execution.crashes") {
		t.Errorf("Query missing crashes table: %s", query)
	}

	if !strings.Contains(query, "limit 10") {
		t.Errorf("Query missing limit: %s", query)
	}

	if !strings.Contains(query, "crash_count") {
		t.Errorf("Query missing crash count: %s", query)
	}
}

func TestTemplates_WorkflowExecutionSuccess(t *testing.T) {
	templates := NewTemplates()

	template := templates.WorkflowExecutionSuccess("during past 30d")
	query := template.Query()

	if !strings.Contains(query, "workflow.executions") {
		t.Errorf("Query missing workflow executions table: %s", query)
	}

	if !strings.Contains(query, "status == \"success\"") {
		t.Errorf("Query missing status filter: %s", query)
	}
}

func TestTemplate_Methods(t *testing.T) {
	templates := NewTemplates()
	template := templates.DevicesWithCrashes("during past 7d", "test.exe")

	// Test Query() returns non-empty string
	query := template.Query()
	if query == "" {
		t.Error("Query() returned empty string")
	}

	// Test QueryBuilder() returns non-nil
	qb := template.QueryBuilder()
	if qb == nil {
		t.Error("QueryBuilder() returned nil")
	}

	// Test ToRequest() creates valid request
	req := template.ToRequest("#test")
	if req.QueryID != "#test" {
		t.Errorf("ToRequest() created invalid request: %+v", req)
	}
}
