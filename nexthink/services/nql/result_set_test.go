package nql

import (
	"testing"
	"time"
)

func TestV2ResultSet_Basic(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		QueryID:           "#test_query",
		ExecutedQuery:     "devices during past 7d | list device.name",
		Rows:              2,
		ExecutionDateTime: "2024-02-08T10:15:30Z",
		Data: []map[string]any{
			{"device.name": "device-01", "crash_count": int64(5)},
			{"device.name": "device-02", "crash_count": int64(3)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	if resultSet.Rows() != 2 {
		t.Errorf("Expected 2 rows, got %d", resultSet.Rows())
	}

	if resultSet.QueryID != "#test_query" {
		t.Errorf("Expected QueryID #test_query, got: %s", resultSet.QueryID)
	}
}

func TestV2ResultSet_GetString(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"device.name": "device-01"},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	name, err := resultSet.GetString(0, "device.name")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}

	if name != "device-01" {
		t.Errorf("Expected device-01, got: %s", name)
	}
}

func TestV2ResultSet_GetInt(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"count": int64(42)},
			{"count": float64(100)},
			{"count": int(50)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	tests := []struct {
		row      int
		expected int64
	}{
		{0, 42},
		{1, 100},
		{2, 50},
	}

	for _, tt := range tests {
		value, err := resultSet.GetInt(tt.row, "count")
		if err != nil {
			t.Fatalf("GetInt(row=%d) failed: %v", tt.row, err)
		}

		if value != tt.expected {
			t.Errorf("Row %d: expected %d, got %d", tt.row, tt.expected, value)
		}
	}
}

func TestV2ResultSet_GetFloat(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"value": float64(42.5)},
			{"value": int64(100)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	// Float value
	val1, err := resultSet.GetFloat(0, "value")
	if err != nil {
		t.Fatalf("GetFloat failed: %v", err)
	}
	if val1 != 42.5 {
		t.Errorf("Expected 42.5, got: %f", val1)
	}

	// Int converted to float
	val2, err := resultSet.GetFloat(1, "value")
	if err != nil {
		t.Fatalf("GetFloat failed: %v", err)
	}
	if val2 != 100.0 {
		t.Errorf("Expected 100.0, got: %f", val2)
	}
}

func TestV2ResultSet_HasField(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"device.name": "device-01", "count": int64(5)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	if !resultSet.HasField("device.name") {
		t.Error("HasField(device.name) should return true")
	}

	if resultSet.HasField("nonexistent") {
		t.Error("HasField(nonexistent) should return false")
	}
}

func TestV2ResultSet_Fields(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"field1": "value1", "field2": "value2"},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	fields := resultSet.Fields()

	if len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fields))
	}
}

func TestV2ResultSet_Filter(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"platform": "Windows", "count": int64(10)},
			{"platform": "macOS", "count": int64(5)},
			{"platform": "Windows", "count": int64(15)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	// Filter to Windows only
	filtered := resultSet.Filter(func(row map[string]any) bool {
		platform, ok := row["platform"].(string)
		return ok && platform == "Windows"
	})

	if len(filtered) != 2 {
		t.Errorf("Expected 2 Windows devices, got %d", len(filtered))
	}
}

func TestV2ResultSet_Map(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"device.name": "device-01", "count": int64(5)},
			{"device.name": "device-02", "count": int64(3)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	// Transform to simplified format
	mapped := resultSet.Map(func(row map[string]any) map[string]any {
		return map[string]any{
			"name": row["device.name"],
		}
	})

	if len(mapped) != 2 {
		t.Errorf("Expected 2 mapped rows, got %d", len(mapped))
	}

	if mapped[0]["name"] != "device-01" {
		t.Errorf("Unexpected mapped value: %v", mapped[0])
	}
}

func TestV2ResultSet_IterateRows(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"device.name": "device-01"},
			{"device.name": "device-02"},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	count := 0
	err := resultSet.IterateRows(func(row int, data map[string]any) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateRows failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 iterations, got %d", count)
	}
}

func TestV2ResultSet_ErrorHandling(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"device.name": "device-01"},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	// Out of bounds row
	_, err := resultSet.Get(10, "device.name")
	if err == nil {
		t.Error("Expected error for out of bounds row")
	}

	// Non-existent field
	_, err = resultSet.Get(0, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent field")
	}

	// Type mismatch
	mockResp2 := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"count": "not a number"},
		},
	}
	resultSet2 := NewV2ResultSet(mockResp2)

	_, err = resultSet2.GetInt(0, "count")
	if err == nil {
		t.Error("Expected error for type mismatch")
	}
}

func TestV1ResultSet_Basic(t *testing.T) {
	mockResp := &ExecuteNQLV1Response{
		QueryID:       "#test_query",
		ExecutedQuery: "devices during past 7d | list device.name",
		Rows:          2,
		Headers:       []string{"device.name", "count"},
		Data: [][]any{
			{"device-01", int64(5)},
			{"device-02", int64(3)},
		},
	}

	resultSet := NewV1ResultSet(mockResp)

	if resultSet.Rows() != 2 {
		t.Errorf("Expected 2 rows, got %d", resultSet.Rows())
	}

	if resultSet.Columns() != 2 {
		t.Errorf("Expected 2 columns, got %d", resultSet.Columns())
	}
}

func TestV1ResultSet_Get(t *testing.T) {
	mockResp := &ExecuteNQLV1Response{
		Headers: []string{"name", "count"},
		Data: [][]any{
			{"device-01", int64(5)},
		},
	}

	resultSet := NewV1ResultSet(mockResp)

	val, err := resultSet.Get(0, 0)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != "device-01" {
		t.Errorf("Expected device-01, got: %v", val)
	}
}

func TestV1ResultSet_FindColumnIndex(t *testing.T) {
	mockResp := &ExecuteNQLV1Response{
		Headers: []string{"name", "count", "platform"},
		Data:    [][]any{},
	}

	resultSet := NewV1ResultSet(mockResp)

	idx, err := resultSet.FindColumnIndex("count")
	if err != nil {
		t.Fatalf("FindColumnIndex failed: %v", err)
	}

	if idx != 1 {
		t.Errorf("Expected index 1, got: %d", idx)
	}

	// Non-existent column
	_, err = resultSet.FindColumnIndex("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent column")
	}
}

func TestV1ResultSet_ToV2Format(t *testing.T) {
	mockResp := &ExecuteNQLV1Response{
		Headers: []string{"name", "count"},
		Data: [][]any{
			{"device-01", int64(5)},
			{"device-02", int64(3)},
		},
	}

	resultSet := NewV1ResultSet(mockResp)

	v2Data := resultSet.ToV2Format()

	if len(v2Data) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(v2Data))
	}

	if v2Data[0]["name"] != "device-01" {
		t.Errorf("Unexpected conversion: %v", v2Data[0])
	}

	if v2Data[0]["count"] != int64(5) {
		t.Errorf("Unexpected conversion: %v", v2Data[0])
	}
}

func TestConvertV1ToV2(t *testing.T) {
	dt := &DateTime{
		Year:   2024,
		Month:  2,
		Day:    8,
		Hour:   10,
		Minute: 15,
		Second: 30,
	}

	mockResp := &ExecuteNQLV1Response{
		QueryID:           "#test",
		ExecutedQuery:     "devices",
		Rows:              1,
		ExecutionDateTime: dt,
		Headers:           []string{"name"},
		Data:              [][]any{{"device-01"}},
	}

	v1ResultSet := NewV1ResultSet(mockResp)
	v2ResultSet := ConvertV1ToV2(v1ResultSet)

	if v2ResultSet == nil {
		t.Fatal("ConvertV1ToV2 returned nil")
	}

	if v2ResultSet.Rows() != 1 {
		t.Errorf("Expected 1 row, got %d", v2ResultSet.Rows())
	}

	if v2ResultSet.QueryID != "#test" {
		t.Errorf("Expected QueryID #test, got: %s", v2ResultSet.QueryID)
	}

	// Verify datetime conversion
	execTime, err := v2ResultSet.ParseExecutionTime()
	if err != nil {
		t.Fatalf("ParseExecutionTime failed: %v", err)
	}

	expectedTime := time.Date(2024, 2, 8, 10, 15, 30, 0, time.UTC)
	if !execTime.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, execTime)
	}
}

func TestV2ResultSet_NilHandling(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"field": nil},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	// GetString with nil should return empty string
	str, err := resultSet.GetString(0, "field")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	if str != "" {
		t.Errorf("Expected empty string for nil, got: %s", str)
	}

	// GetInt with nil should return 0
	num, err := resultSet.GetInt(0, "field")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	if num != 0 {
		t.Errorf("Expected 0 for nil, got: %d", num)
	}
}

func TestV1ResultSet_GetByColumnName(t *testing.T) {
	mockResp := &ExecuteNQLV1Response{
		Headers: []string{"name", "count"},
		Data: [][]any{
			{"device-01", int64(5)},
		},
	}

	resultSet := NewV1ResultSet(mockResp)

	val, err := resultSet.GetByColumnName(0, "name")
	if err != nil {
		t.Fatalf("GetByColumnName failed: %v", err)
	}

	if val != "device-01" {
		t.Errorf("Expected device-01, got: %v", val)
	}

	// Non-existent column
	_, err = resultSet.GetByColumnName(0, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent column")
	}
}

func TestV2ResultSet_ToJSON(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"device.name": "device-01", "count": int64(5)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	jsonData, err := resultSet.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("ToJSON returned empty data")
	}
}

func TestV2ResultSet_ParseExecutionTime(t *testing.T) {
	tests := []struct {
		name          string
		dateTimeStr   string
		shouldSucceed bool
	}{
		{
			name:          "RFC3339 format",
			dateTimeStr:   "2024-02-08T10:15:30Z",
			shouldSucceed: true,
		},
		{
			name:          "Alternative format",
			dateTimeStr:   "2024-02-08T10:15:30",
			shouldSucceed: true,
		},
		{
			name:          "Invalid format",
			dateTimeStr:   "invalid",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResp := &ExecuteNQLV2Response{
				ExecutionDateTime: tt.dateTimeStr,
			}

			resultSet := NewV2ResultSet(mockResp)

			_, err := resultSet.ParseExecutionTime()
			if tt.shouldSucceed && err != nil {
				t.Errorf("Expected success, got error: %v", err)
			}
			if !tt.shouldSucceed && err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestNewV2ResultSet_Nil(t *testing.T) {
	resultSet := NewV2ResultSet(nil)

	if resultSet != nil {
		t.Error("Expected nil result set for nil response")
	}
}

func TestNewV1ResultSet_Nil(t *testing.T) {
	resultSet := NewV1ResultSet(nil)

	if resultSet != nil {
		t.Error("Expected nil result set for nil response")
	}
}

func TestV2ResultSet_EmptyData(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Rows: 0,
		Data: []map[string]any{},
	}

	resultSet := NewV2ResultSet(mockResp)

	if resultSet.Rows() != 0 {
		t.Errorf("Expected 0 rows, got %d", resultSet.Rows())
	}

	fields := resultSet.Fields()
	if fields != nil {
		t.Errorf("Expected nil fields for empty data, got: %v", fields)
	}
}

func TestV2ResultSet_GetRow(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		Data: []map[string]any{
			{"name": "device-01", "count": int64(5)},
			{"name": "device-02", "count": int64(3)},
		},
	}

	resultSet := NewV2ResultSet(mockResp)

	row, err := resultSet.GetRow(0)
	if err != nil {
		t.Fatalf("GetRow failed: %v", err)
	}

	if row["name"] != "device-01" {
		t.Errorf("Unexpected row data: %v", row)
	}

	// Out of bounds
	_, err = resultSet.GetRow(10)
	if err == nil {
		t.Error("Expected error for out of bounds row")
	}
}

func TestV1ResultSet_IterateRows(t *testing.T) {
	mockResp := &ExecuteNQLV1Response{
		Headers: []string{"name", "count"},
		Data: [][]any{
			{"device-01", int64(5)},
			{"device-02", int64(3)},
		},
	}

	resultSet := NewV1ResultSet(mockResp)

	count := 0
	err := resultSet.IterateRows(func(row int, values []any) error {
		count++
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}
		return nil
	})

	if err != nil {
		t.Fatalf("IterateRows failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 iterations, got %d", count)
	}
}
