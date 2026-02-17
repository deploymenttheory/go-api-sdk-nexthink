package client

import (
	"testing"
	"time"
)

func TestNewQueryBuilder(t *testing.T) {
	qb := NewQueryBuilder()
	if qb == nil {
		t.Fatal("NewQueryBuilder() returned nil")
	}
	if qb.params == nil {
		t.Fatal("QueryBuilder params map is nil")
	}
	if len(qb.params) != 0 {
		t.Errorf("NewQueryBuilder() should start with empty params, got %d params", len(qb.params))
	}
}

func TestQueryBuilder_AddString(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		wantAdded bool
	}{
		{
			name:      "add valid string",
			key:       "query",
			value:     "test",
			wantAdded: true,
		},
		{
			name:      "skip empty string",
			key:       "query",
			value:     "",
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddString(tt.key, tt.value)

			// Check fluent interface returns self
			if result != qb {
				t.Error("AddString() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.value {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.value)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddInt(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     int
		wantAdded bool
		wantValue string
	}{
		{
			name:      "add positive integer",
			key:       "limit",
			value:     100,
			wantAdded: true,
			wantValue: "100",
		},
		{
			name:      "skip zero",
			key:       "limit",
			value:     0,
			wantAdded: false,
		},
		{
			name:      "skip negative",
			key:       "limit",
			value:     -1,
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddInt(tt.key, tt.value)

			if result != qb {
				t.Error("AddInt() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.wantValue {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.wantValue)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddInt64(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     int64
		wantAdded bool
		wantValue string
	}{
		{
			name:      "add positive int64",
			key:       "timestamp",
			value:     1640000000,
			wantAdded: true,
			wantValue: "1640000000",
		},
		{
			name:      "skip zero",
			key:       "timestamp",
			value:     0,
			wantAdded: false,
		},
		{
			name:      "skip negative",
			key:       "timestamp",
			value:     -1,
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddInt64(tt.key, tt.value)

			if result != qb {
				t.Error("AddInt64() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.wantValue {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.wantValue)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddBool(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     bool
		wantValue string
	}{
		{
			name:      "add true",
			key:       "enabled",
			value:     true,
			wantValue: "true",
		},
		{
			name:      "add false",
			key:       "enabled",
			value:     false,
			wantValue: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddBool(tt.key, tt.value)

			if result != qb {
				t.Error("AddBool() should return self for fluent interface")
			}

			if !qb.Has(tt.key) {
				t.Errorf("Expected parameter %q to be added", tt.key)
			}
			if got := qb.Get(tt.key); got != tt.wantValue {
				t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.wantValue)
			}
		})
	}
}

func TestQueryBuilder_AddTime(t *testing.T) {
	testTime := time.Date(2021, 12, 20, 15, 30, 0, 0, time.UTC)
	expectedValue := testTime.Format(time.RFC3339)

	tests := []struct {
		name      string
		key       string
		value     time.Time
		wantAdded bool
		wantValue string
	}{
		{
			name:      "add valid time",
			key:       "created_at",
			value:     testTime,
			wantAdded: true,
			wantValue: expectedValue,
		},
		{
			name:      "skip zero time",
			key:       "created_at",
			value:     time.Time{},
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddTime(tt.key, tt.value)

			if result != qb {
				t.Error("AddTime() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.wantValue {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.wantValue)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddStringSlice(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		values    []string
		wantAdded bool
		wantValue string
	}{
		{
			name:      "add string slice",
			key:       "tags",
			values:    []string{"tag1", "tag2", "tag3"},
			wantAdded: true,
			wantValue: "tag1,tag2,tag3",
		},
		{
			name:      "filter empty strings",
			key:       "tags",
			values:    []string{"tag1", "", "tag2"},
			wantAdded: true,
			wantValue: "tag1,tag2",
		},
		{
			name:      "skip empty slice",
			key:       "tags",
			values:    []string{},
			wantAdded: false,
		},
		{
			name:      "skip all empty strings",
			key:       "tags",
			values:    []string{"", "", ""},
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddStringSlice(tt.key, tt.values)

			if result != qb {
				t.Error("AddStringSlice() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.wantValue {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.wantValue)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddIntSlice(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		values    []int
		wantAdded bool
		wantValue string
	}{
		{
			name:      "add int slice",
			key:       "ids",
			values:    []int{1, 2, 3},
			wantAdded: true,
			wantValue: "1,2,3",
		},
		{
			name:      "add int slice with zeros",
			key:       "ids",
			values:    []int{1, 0, 2},
			wantAdded: true,
			wantValue: "1,0,2",
		},
		{
			name:      "skip empty slice",
			key:       "ids",
			values:    []int{},
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddIntSlice(tt.key, tt.values)

			if result != qb {
				t.Error("AddIntSlice() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.wantValue {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.wantValue)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddCustom(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.AddCustom("custom_key", "custom_value")

	if result != qb {
		t.Error("AddCustom() should return self for fluent interface")
	}

	if !qb.Has("custom_key") {
		t.Error("Expected custom_key to be added")
	}
	if got := qb.Get("custom_key"); got != "custom_value" {
		t.Errorf("Get(custom_key) = %q, want %q", got, "custom_value")
	}
}

func TestQueryBuilder_AddIfNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		wantAdded bool
	}{
		{
			name:      "add non-empty value",
			key:       "filter",
			value:     "active",
			wantAdded: true,
		},
		{
			name:      "skip empty value",
			key:       "filter",
			value:     "",
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddIfNotEmpty(tt.key, tt.value)

			if result != qb {
				t.Error("AddIfNotEmpty() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.value {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.value)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_AddIfTrue(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		key       string
		value     string
		wantAdded bool
	}{
		{
			name:      "add when condition is true",
			condition: true,
			key:       "include_deleted",
			value:     "yes",
			wantAdded: true,
		},
		{
			name:      "skip when condition is false",
			condition: false,
			key:       "include_deleted",
			value:     "yes",
			wantAdded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			result := qb.AddIfTrue(tt.condition, tt.key, tt.value)

			if result != qb {
				t.Error("AddIfTrue() should return self for fluent interface")
			}

			if tt.wantAdded {
				if !qb.Has(tt.key) {
					t.Errorf("Expected parameter %q to be added", tt.key)
				}
				if got := qb.Get(tt.key); got != tt.value {
					t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.value)
				}
			} else {
				if qb.Has(tt.key) {
					t.Errorf("Expected parameter %q not to be added", tt.key)
				}
			}
		})
	}
}

func TestQueryBuilder_Merge(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("existing", "value1")

	other := map[string]string{
		"new_key":  "value2",
		"another":  "value3",
		"existing": "overwritten",
	}

	result := qb.Merge(other)

	if result != qb {
		t.Error("Merge() should return self for fluent interface")
	}

	if qb.Count() != 3 {
		t.Errorf("Expected 3 parameters after merge, got %d", qb.Count())
	}

	if got := qb.Get("existing"); got != "overwritten" {
		t.Errorf("Get(existing) = %q, want %q", got, "overwritten")
	}
	if got := qb.Get("new_key"); got != "value2" {
		t.Errorf("Get(new_key) = %q, want %q", got, "value2")
	}
	if got := qb.Get("another"); got != "value3" {
		t.Errorf("Get(another) = %q, want %q", got, "value3")
	}
}

func TestQueryBuilder_Remove(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key1", "value1")
	qb.AddString("key2", "value2")

	result := qb.Remove("key1")

	if result != qb {
		t.Error("Remove() should return self for fluent interface")
	}

	if qb.Has("key1") {
		t.Error("Expected key1 to be removed")
	}
	if !qb.Has("key2") {
		t.Error("Expected key2 to still exist")
	}

	// Removing non-existent key should not error
	qb.Remove("non_existent")
}

func TestQueryBuilder_HasAndGet(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("test", "value")

	if !qb.Has("test") {
		t.Error("Has(test) should return true")
	}
	if qb.Has("missing") {
		t.Error("Has(missing) should return false")
	}

	if got := qb.Get("test"); got != "value" {
		t.Errorf("Get(test) = %q, want %q", got, "value")
	}
	if got := qb.Get("missing"); got != "" {
		t.Errorf("Get(missing) = %q, want empty string", got)
	}
}

func TestQueryBuilder_Build(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key1", "value1")
	qb.AddInt("key2", 42)

	result := qb.Build()

	if len(result) != 2 {
		t.Errorf("Build() returned %d parameters, want 2", len(result))
	}
	if result["key1"] != "value1" {
		t.Errorf("Build()[key1] = %q, want %q", result["key1"], "value1")
	}
	if result["key2"] != "42" {
		t.Errorf("Build()[key2] = %q, want %q", result["key2"], "42")
	}

	// Verify returned map is a copy (modifying it shouldn't affect the builder)
	result["key3"] = "value3"
	if qb.Has("key3") {
		t.Error("Modifying Build() result should not affect the builder")
	}
}

func TestQueryBuilder_BuildString(t *testing.T) {
	tests := []struct {
		name string
		add  func(*QueryBuilder)
		want map[string]bool
	}{
		{
			name: "empty builder",
			add:  func(qb *QueryBuilder) {},
			want: map[string]bool{"": true},
		},
		{
			name: "single parameter",
			add: func(qb *QueryBuilder) {
				qb.AddString("key", "value")
			},
			want: map[string]bool{"key=value": true},
		},
		{
			name: "multiple parameters",
			add: func(qb *QueryBuilder) {
				qb.AddString("key1", "value1")
				qb.AddString("key2", "value2")
			},
			want: map[string]bool{
				"key1=value1&key2=value2": true,
				"key2=value2&key1=value1": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.add(qb)
			result := qb.BuildString()

			if !tt.want[result] {
				t.Errorf("BuildString() = %q, not in expected values %v", result, tt.want)
			}
		})
	}
}

func TestQueryBuilder_Clear(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddString("key1", "value1")
	qb.AddString("key2", "value2")

	if qb.Count() != 2 {
		t.Errorf("Expected 2 parameters before clear, got %d", qb.Count())
	}

	result := qb.Clear()

	if result != qb {
		t.Error("Clear() should return self for fluent interface")
	}

	if qb.Count() != 0 {
		t.Errorf("Expected 0 parameters after clear, got %d", qb.Count())
	}
	if !qb.IsEmpty() {
		t.Error("IsEmpty() should return true after clear")
	}
}

func TestQueryBuilder_CountAndIsEmpty(t *testing.T) {
	qb := NewQueryBuilder()

	if qb.Count() != 0 {
		t.Errorf("Count() = %d, want 0", qb.Count())
	}
	if !qb.IsEmpty() {
		t.Error("IsEmpty() should return true for new builder")
	}

	qb.AddString("key1", "value1")

	if qb.Count() != 1 {
		t.Errorf("Count() = %d, want 1", qb.Count())
	}
	if qb.IsEmpty() {
		t.Error("IsEmpty() should return false after adding parameter")
	}

	qb.AddString("key2", "value2")

	if qb.Count() != 2 {
		t.Errorf("Count() = %d, want 2", qb.Count())
	}
}

func TestQueryBuilder_FluentChaining(t *testing.T) {
	// Test that all methods can be chained fluently
	qb := NewQueryBuilder()
	result := qb.
		AddString("query", "test").
		AddInt("limit", 10).
		AddBool("active", true).
		AddStringSlice("tags", []string{"tag1", "tag2"}).
		AddCustom("custom", "value").
		Build()

	if len(result) != 5 {
		t.Errorf("Expected 5 parameters after chaining, got %d", len(result))
	}

	expected := map[string]string{
		"query":  "test",
		"limit":  "10",
		"active": "true",
		"tags":   "tag1,tag2",
		"custom": "value",
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("result[%q] = %q, want %q", k, result[k], v)
		}
	}
}
