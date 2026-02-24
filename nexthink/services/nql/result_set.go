package nql

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Result processing helpers for NQL query responses
// Provides type-safe access to V1 and V2 response data with conversion utilities

// =============================================================================
// V1 Result Set
// =============================================================================

// V1ResultSet provides typed access to V1 response data
type V1ResultSet struct {
	QueryID           string
	ExecutedQuery     string
	RowCount          int64
	ExecutionDateTime *DateTime
	Headers           []string
	data              [][]any
}

// NewV1ResultSet creates a new V1 result set from a response
func NewV1ResultSet(resp *ExecuteNQLV1Response) *V1ResultSet {
	if resp == nil {
		return nil
	}
	
	return &V1ResultSet{
		QueryID:           resp.QueryID,
		ExecutedQuery:     resp.ExecutedQuery,
		RowCount:          resp.Rows,
		ExecutionDateTime: resp.ExecutionDateTime,
		Headers:           resp.Headers,
		data:              resp.Data,
	}
}

// Rows returns the number of rows
func (rs *V1ResultSet) Rows() int {
	if rs.data == nil {
		return 0
	}
	return len(rs.data)
}

// Columns returns the number of columns
func (rs *V1ResultSet) Columns() int {
	if rs.Headers == nil {
		return 0
	}
	return len(rs.Headers)
}

// Get retrieves a cell value with bounds checking
func (rs *V1ResultSet) Get(row, col int) (any, error) {
	if row < 0 || row >= rs.Rows() {
		return nil, fmt.Errorf("row index %d out of bounds (0-%d)", row, rs.Rows()-1)
	}
	
	if col < 0 || col >= rs.Columns() {
		return nil, fmt.Errorf("column index %d out of bounds (0-%d)", col, rs.Columns()-1)
	}
	
	return rs.data[row][col], nil
}

// GetString retrieves a string value with type checking
func (rs *V1ResultSet) GetString(row, col int) (string, error) {
	val, err := rs.Get(row, col)
	if err != nil {
		return "", err
	}
	
	if val == nil {
		return "", nil
	}
	
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value at [%d][%d] is not a string: %T", row, col, val)
	}
	
	return str, nil
}

// GetInt retrieves an int64 value with type checking and conversion
func (rs *V1ResultSet) GetInt(row, col int) (int64, error) {
	val, err := rs.Get(row, col)
	if err != nil {
		return 0, err
	}
	
	if val == nil {
		return 0, nil
	}
	
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("value at [%d][%d] cannot be converted to int64: %T", row, col, val)
	}
}

// GetFloat retrieves a float64 value with type checking and conversion
func (rs *V1ResultSet) GetFloat(row, col int) (float64, error) {
	val, err := rs.Get(row, col)
	if err != nil {
		return 0, err
	}
	
	if val == nil {
		return 0, nil
	}
	
	switch v := val.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("value at [%d][%d] cannot be converted to float64: %T", row, col, val)
	}
}

// GetBool retrieves a boolean value with type checking
func (rs *V1ResultSet) GetBool(row, col int) (bool, error) {
	val, err := rs.Get(row, col)
	if err != nil {
		return false, err
	}
	
	if val == nil {
		return false, nil
	}
	
	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("value at [%d][%d] is not a boolean: %T", row, col, val)
	}
	
	return b, nil
}

// GetRow retrieves an entire row
func (rs *V1ResultSet) GetRow(row int) ([]any, error) {
	if row < 0 || row >= rs.Rows() {
		return nil, fmt.Errorf("row index %d out of bounds (0-%d)", row, rs.Rows()-1)
	}
	
	return rs.data[row], nil
}

// IterateRows iterates over all rows with a callback function
func (rs *V1ResultSet) IterateRows(fn func(row int, values []any) error) error {
	for i, rowData := range rs.data {
		if err := fn(i, rowData); err != nil {
			return fmt.Errorf("error processing row %d: %w", i, err)
		}
	}
	return nil
}

// ToV2Format converts V1 response to V2 format (map-based)
func (rs *V1ResultSet) ToV2Format() []map[string]any {
	if rs.data == nil || rs.Headers == nil {
		return nil
	}
	
	result := make([]map[string]any, 0, len(rs.data))
	
	for _, row := range rs.data {
		rowMap := make(map[string]any)
		for i, header := range rs.Headers {
			if i < len(row) {
				rowMap[header] = row[i]
			}
		}
		result = append(result, rowMap)
	}
	
	return result
}

// ToJSON converts the result set to JSON
func (rs *V1ResultSet) ToJSON() ([]byte, error) {
	return json.Marshal(rs.ToV2Format())
}

// FindColumnIndex finds the index of a column by name
func (rs *V1ResultSet) FindColumnIndex(columnName string) (int, error) {
	for i, header := range rs.Headers {
		if header == columnName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column '%s' not found", columnName)
}

// GetByColumnName retrieves a value by row index and column name
func (rs *V1ResultSet) GetByColumnName(row int, columnName string) (any, error) {
	colIdx, err := rs.FindColumnIndex(columnName)
	if err != nil {
		return nil, err
	}
	return rs.Get(row, colIdx)
}

// =============================================================================
// V2 Result Set
// =============================================================================

// V2ResultSet provides typed access to V2 response data
type V2ResultSet struct {
	QueryID           string
	ExecutedQuery     string
	RowCount          int64
	ExecutionDateTime string
	data              []map[string]any
}

// NewV2ResultSet creates a new V2 result set from a response
func NewV2ResultSet(resp *ExecuteNQLV2Response) *V2ResultSet {
	if resp == nil {
		return nil
	}
	
	return &V2ResultSet{
		QueryID:           resp.QueryID,
		ExecutedQuery:     resp.ExecutedQuery,
		RowCount:          resp.Rows,
		ExecutionDateTime: resp.ExecutionDateTime,
		data:              resp.Data,
	}
}

// Rows returns the number of rows
func (rs *V2ResultSet) Rows() int {
	if rs.data == nil {
		return 0
	}
	return len(rs.data)
}

// Get retrieves a field value by row index and field name
func (rs *V2ResultSet) Get(row int, field string) (any, error) {
	if row < 0 || row >= rs.Rows() {
		return nil, fmt.Errorf("row index %d out of bounds (0-%d)", row, rs.Rows()-1)
	}
	
	val, exists := rs.data[row][field]
	if !exists {
		return nil, fmt.Errorf("field '%s' not found in row %d", field, row)
	}
	
	return val, nil
}

// GetString retrieves a string value with type checking
func (rs *V2ResultSet) GetString(row int, field string) (string, error) {
	val, err := rs.Get(row, field)
	if err != nil {
		return "", err
	}
	
	if val == nil {
		return "", nil
	}
	
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("field '%s' at row %d is not a string: %T", field, row, val)
	}
	
	return str, nil
}

// GetInt retrieves an int64 value with type checking and conversion
func (rs *V2ResultSet) GetInt(row int, field string) (int64, error) {
	val, err := rs.Get(row, field)
	if err != nil {
		return 0, err
	}
	
	if val == nil {
		return 0, nil
	}
	
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("field '%s' at row %d cannot be converted to int64: %T", field, row, val)
	}
}

// GetFloat retrieves a float64 value with type checking and conversion
func (rs *V2ResultSet) GetFloat(row int, field string) (float64, error) {
	val, err := rs.Get(row, field)
	if err != nil {
		return 0, err
	}
	
	if val == nil {
		return 0, nil
	}
	
	switch v := val.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("field '%s' at row %d cannot be converted to float64: %T", field, row, val)
	}
}

// GetBool retrieves a boolean value with type checking
func (rs *V2ResultSet) GetBool(row int, field string) (bool, error) {
	val, err := rs.Get(row, field)
	if err != nil {
		return false, err
	}
	
	if val == nil {
		return false, nil
	}
	
	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("field '%s' at row %d is not a boolean: %T", field, row, val)
	}
	
	return b, nil
}

// GetRow retrieves an entire row
func (rs *V2ResultSet) GetRow(row int) (map[string]any, error) {
	if row < 0 || row >= rs.Rows() {
		return nil, fmt.Errorf("row index %d out of bounds (0-%d)", row, rs.Rows()-1)
	}
	
	return rs.data[row], nil
}

// Fields returns all field names from the first row
func (rs *V2ResultSet) Fields() []string {
	if rs.Rows() == 0 {
		return nil
	}
	
	fields := make([]string, 0, len(rs.data[0]))
	for field := range rs.data[0] {
		fields = append(fields, field)
	}
	
	return fields
}

// HasField checks if a field exists in the result set
func (rs *V2ResultSet) HasField(field string) bool {
	if rs.Rows() == 0 {
		return false
	}
	
	_, exists := rs.data[0][field]
	return exists
}

// Filter returns rows matching a predicate
func (rs *V2ResultSet) Filter(fn func(row map[string]any) bool) []map[string]any {
	if rs.data == nil {
		return nil
	}
	
	result := make([]map[string]any, 0)
	for _, row := range rs.data {
		if fn(row) {
			result = append(result, row)
		}
	}
	
	return result
}

// Map transforms each row using a mapping function
func (rs *V2ResultSet) Map(fn func(row map[string]any) map[string]any) []map[string]any {
	if rs.data == nil {
		return nil
	}
	
	result := make([]map[string]any, 0, len(rs.data))
	for _, row := range rs.data {
		result = append(result, fn(row))
	}
	
	return result
}

// IterateRows iterates over all rows with a callback function
func (rs *V2ResultSet) IterateRows(fn func(row int, data map[string]any) error) error {
	for i, rowData := range rs.data {
		if err := fn(i, rowData); err != nil {
			return fmt.Errorf("error processing row %d: %w", i, err)
		}
	}
	return nil
}

// ToJSON converts the result set to JSON
func (rs *V2ResultSet) ToJSON() ([]byte, error) {
	return json.Marshal(rs.data)
}

// ParseExecutionTime parses the ISO format execution datetime to time.Time
func (rs *V2ResultSet) ParseExecutionTime() (time.Time, error) {
	if rs.ExecutionDateTime == "" {
		return time.Time{}, fmt.Errorf("execution datetime is empty")
	}
	
	// Try parsing as ISO 8601 format
	t, err := time.Parse(time.RFC3339, rs.ExecutionDateTime)
	if err != nil {
		// Try alternative format
		t, err = time.Parse("2006-01-02T15:04:05", rs.ExecutionDateTime)
	}
	
	return t, err
}

// =============================================================================
// Helper Functions
// =============================================================================

// ConvertV1ToV2 converts a V1 result set to V2 format
func ConvertV1ToV2(v1 *V1ResultSet) *V2ResultSet {
	if v1 == nil {
		return nil
	}
	
	v2Data := v1.ToV2Format()
	
	// Convert execution datetime
	execTime := ""
	if v1.ExecutionDateTime != nil {
		dt := v1.ExecutionDateTime
		t := time.Date(
			int(dt.Year), time.Month(dt.Month), int(dt.Day),
			int(dt.Hour), int(dt.Minute), int(dt.Second),
			0, time.UTC,
		)
		execTime = t.Format(time.RFC3339)
	}
	
	return &V2ResultSet{
		QueryID:           v1.QueryID,
		ExecutedQuery:     v1.ExecutedQuery,
		RowCount:          v1.RowCount,
		ExecutionDateTime: execTime,
		data:              v2Data,
	}
}
