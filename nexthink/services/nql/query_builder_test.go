package nql

import (
	"strings"
	"testing"
)

func TestQueryBuilder_BasicQuery(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		DuringPast(7, Days).
		List("device.name").
		Build()

	expected := "devices during past 7d\n| list device.name"

	if query != expected {
		t.Errorf("Query mismatch:\nGot:  %s\nWant: %s", query, expected)
	}
}

func TestQueryBuilder_WithFilter(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		DuringPast(7, Days).
		WhereEquals("platform", "Windows").
		List("device.name").
		Build()

	if !strings.Contains(query, "| where platform == \"Windows\"") {
		t.Errorf("Query missing where clause: %s", query)
	}
}

func TestQueryBuilder_WithEventData(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		With("execution.crashes during past 7d").
		ComputeSum("total_crashes", "number_of_crashes").
		Build()

	if !strings.Contains(query, "| with execution.crashes during past 7d") {
		t.Errorf("Query missing with clause: %s", query)
	}

	if !strings.Contains(query, "| compute total_crashes = number_of_crashes.sum()") {
		t.Errorf("Query missing compute clause: %s", query)
	}
}

func TestQueryBuilder_Summarize(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		DuringPast(7, Days).
		SummarizeCount("device_count").
		SummarizeBy(FieldOSPlatform).
		Build()

	if !strings.Contains(query, "| summarize device_count = count()") {
		t.Errorf("Query missing summarize: %s", query)
	}

	if !strings.Contains(query, "by operating_system.platform") {
		t.Errorf("Query missing summarize by: %s", query)
	}
}

func TestQueryBuilder_Sort(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		List("device.name").
		SortDesc("device.name").
		Build()

	if !strings.Contains(query, "| sort device.name desc") {
		t.Errorf("Query missing sort clause: %s", query)
	}
}

func TestQueryBuilder_Limit(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		List("device.name").
		Limit(10).
		Build()

	if !strings.Contains(query, "| limit 10") {
		t.Errorf("Query missing limit clause: %s", query)
	}
}

func TestQueryBuilder_Comment(t *testing.T) {
	query := NewQueryBuilder().
		Comment("Test comment").
		FromDevices().
		List("device.name").
		Build()

	if !strings.Contains(query, "/* Test comment */") {
		t.Errorf("Query missing comment: %s", query)
	}
}

func TestQueryBuilder_WhereIn(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		WhereIn("platform", []string{"Windows", "macOS"}).
		List("device.name").
		Build()

	if !strings.Contains(query, `| where platform in ["Windows", "macOS"]`) {
		t.Errorf("Query missing where in clause: %s", query)
	}
}

func TestQueryBuilder_MultipleWhere(t *testing.T) {
	query := NewQueryBuilder().
		FromDevices().
		WhereEquals("platform", "Windows").
		WhereGreaterEqual("memory", "8GB").
		List("device.name").
		Build()

	if !strings.Contains(query, `platform == "Windows"`) {
		t.Errorf("Query missing first where: %s", query)
	}

	if !strings.Contains(query, "memory >= 8GB") {
		t.Errorf("Query missing second where: %s", query)
	}
}

func TestQueryBuilder_Validation_MissingTable(t *testing.T) {
	qb := NewQueryBuilder().
		List("device.name")

	err := qb.Validate()
	if err == nil {
		t.Error("Expected validation error for missing table, got nil")
	}

	if !strings.Contains(err.Error(), "table selection is required") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestQueryBuilder_Validation_ComputeWithoutWith(t *testing.T) {
	qb := NewQueryBuilder().
		FromDevices().
		ComputeSum("crashes", "number_of_crashes")

	err := qb.Validate()
	if err == nil {
		t.Error("Expected validation error for compute without with/include, got nil")
	}

	if !strings.Contains(err.Error(), "compute clause requires") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestQueryBuilder_Validation_ListAndSummarize(t *testing.T) {
	qb := NewQueryBuilder().
		FromDevices().
		List("device.name").
		SummarizeCount("total")

	err := qb.Validate()
	if err == nil {
		t.Error("Expected validation error for list + summarize, got nil")
	}

	if !strings.Contains(err.Error(), "cannot use both list and summarize") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestQueryBuilder_BuildAndValidate(t *testing.T) {
	qb := NewQueryBuilder().
		FromDevices().
		DuringPast(7, Days).
		List("device.name")

	query, err := qb.BuildAndValidate()
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}

	if query == "" {
		t.Error("BuildAndValidate returned empty query")
	}
}

func TestQueryBuilder_Shortcuts(t *testing.T) {
	tests := []struct {
		name     string
		builder  *QueryBuilder
		expected string
	}{
		{
			name:     "FromDevices",
			builder:  NewQueryBuilder().FromDevices(),
			expected: "devices",
		},
		{
			name:     "FromUsers",
			builder:  NewQueryBuilder().FromUsers(),
			expected: "users",
		},
		{
			name:     "FromApplications",
			builder:  NewQueryBuilder().FromApplications(),
			expected: "applications",
		},
		{
			name:     "FromBinaries",
			builder:  NewQueryBuilder().FromBinaries(),
			expected: "binaries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.builder.Build()
			if !strings.HasPrefix(query, tt.expected) {
				t.Errorf("Expected query to start with %s, got: %s", tt.expected, query)
			}
		})
	}
}

func TestQueryBuilder_ComputeHelpers(t *testing.T) {
	tests := []struct {
		name     string
		builder  *QueryBuilder
		contains string
	}{
		{
			name: "ComputeCount",
			builder: NewQueryBuilder().FromDevices().
				With("execution.crashes during past 7d").
				ComputeCount("total"),
			contains: "total = count()",
		},
		{
			name: "ComputeSum",
			builder: NewQueryBuilder().FromDevices().
				With("execution.crashes during past 7d").
				ComputeSum("total", "crashes"),
			contains: "total = crashes.sum()",
		},
		{
			name: "ComputeAvg",
			builder: NewQueryBuilder().FromDevices().
				With("device_performance.events during past 7d").
				ComputeAvg("avg_memory", "memory"),
			contains: "avg_memory = memory.avg()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.builder.Build()
			if !strings.Contains(query, tt.contains) {
				t.Errorf("Expected query to contain %s, got: %s", tt.contains, query)
			}
		})
	}
}

func TestQueryBuilder_TimeSelection(t *testing.T) {
	tests := []struct {
		name     string
		builder  *QueryBuilder
		expected string
	}{
		{
			name: "DuringPast",
			builder: NewQueryBuilder().FromDevices().
				DuringPast(7, Days),
			expected: "devices during past 7d",
		},
		{
			name: "On",
			builder: NewQueryBuilder().FromDevices().
				On("2024-02-08"),
			expected: "devices on 2024-02-08",
		},
		{
			name: "FromTo",
			builder: NewQueryBuilder().FromDevices().
				FromTo("2024-01-01", "2024-01-31"),
			expected: "devices from 2024-01-01 to 2024-01-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.builder.Build()
			if !strings.HasPrefix(query, tt.expected) {
				t.Errorf("Expected query to start with %s, got: %s", tt.expected, query)
			}
		})
	}
}

func TestQueryBuilder_ComplexQuery(t *testing.T) {
	query := NewQueryBuilder().
		Comment("Complex query test").
		FromDevices().
		DuringPast(7, Days).
		With("execution.crashes during past 7d").
		WhereEquals(FieldOSPlatform, PlatformWindows).
		WhereIn(FieldHardwareType, []string{HardwareTypeLaptop, HardwareTypeDesktop}).
		ComputeSum("crashes", "number_of_crashes").
		WhereGreaterEqual("crashes", "3").
		List(FieldDeviceName, "crashes").
		SortDesc("crashes").
		Limit(20).
		Build()

	// Verify all components are present
	components := []string{
		"/* Complex query test */",
		"devices during past 7d",
		"| with execution.crashes during past 7d",
		"platform == \"Windows\"",
		"hardware.type in [\"laptop\", \"desktop\"]",
		"| compute crashes = number_of_crashes.sum()",
		"crashes >= 3",
		"| list device.name, crashes",
		"| sort crashes desc",
		"| limit 20",
	}

	for _, component := range components {
		if !strings.Contains(query, component) {
			t.Errorf("Query missing component: %s\nFull query: %s", component, query)
		}
	}
}
