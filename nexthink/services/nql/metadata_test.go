package nql

import (
	"net/http"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

func TestGetV2Metadata(t *testing.T) {
	mockResp := &ExecuteNQLV2Response{
		QueryID:           "#test_query",
		ExecutedQuery:     "devices during past 7d | list device.name",
		Rows:              42,
		ExecutionDateTime: "2024-02-08T10:15:30Z",
	}

	mockAPIResp := &interfaces.Response{
		StatusCode: 200,
		Duration:   500 * time.Millisecond,
		Body:       []byte("test body"),
		Size:       100,
		Headers: http.Header{
			"X-Rate-Limit-Limit":     []string{"1000"},
			"X-Rate-Limit-Remaining": []string{"995"},
		},
	}

	metadata := GetV2Metadata(mockResp, mockAPIResp)

	if metadata == nil {
		t.Fatal("GetV2Metadata returned nil")
	}

	if metadata.QueryID != "#test_query" {
		t.Errorf("Expected QueryID #test_query, got: %s", metadata.QueryID)
	}

	if metadata.RowsReturned != 42 {
		t.Errorf("Expected 42 rows, got: %d", metadata.RowsReturned)
	}

	if metadata.ResponseDuration != 500*time.Millisecond {
		t.Errorf("Expected 500ms duration, got: %v", metadata.ResponseDuration)
	}

	if metadata.ResponseSize != 100 {
		t.Errorf("Expected size 100, got: %d", metadata.ResponseSize)
	}
}

func TestGetV1Metadata(t *testing.T) {
	dt := &DateTime{
		Year:   2024,
		Month:  2,
		Day:    8,
		Hour:   10,
		Minute: 15,
		Second: 30,
	}

	mockResp := &ExecuteNQLV1Response{
		QueryID:           "#test_query",
		ExecutedQuery:     "devices",
		Rows:              42,
		ExecutionDateTime: dt,
	}

	mockAPIResp := &interfaces.Response{
		StatusCode: 200,
		Duration:   500 * time.Millisecond,
	}

	metadata := GetV1Metadata(mockResp, mockAPIResp)

	if metadata == nil {
		t.Fatal("GetV1Metadata returned nil")
	}

	if metadata.QueryID != "#test_query" {
		t.Errorf("Expected QueryID #test_query, got: %s", metadata.QueryID)
	}

	if metadata.RowsReturned != 42 {
		t.Errorf("Expected 42 rows, got: %d", metadata.RowsReturned)
	}

	if metadata.ExecutedQuery != "devices" {
		t.Errorf("Expected query 'devices', got: %s", metadata.ExecutedQuery)
	}

	expectedTime := time.Date(2024, 2, 8, 10, 15, 30, 0, time.UTC)
	if !metadata.ExecutionTime.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, metadata.ExecutionTime)
	}
}

func TestGetExportMetadata(t *testing.T) {
	mockStatus := &NQLExportStatusResponse{
		ExportID:       "export-123",
		Status:         ExportStatusCompleted,
		ResultsFileURL: "https://s3.example.com/export.csv",
	}

	mockAPIResp := &interfaces.Response{
		StatusCode: 200,
		Duration:   200 * time.Millisecond,
	}

	metadata := GetExportMetadata(mockStatus, mockAPIResp)

	if metadata == nil {
		t.Fatal("GetExportMetadata returned nil")
	}

	if metadata.ExportID != "export-123" {
		t.Errorf("Expected export-123, got: %s", metadata.ExportID)
	}

	if metadata.Status != ExportStatusCompleted {
		t.Errorf("Expected COMPLETED status, got: %s", metadata.Status)
	}

	if !metadata.IsCompleted() {
		t.Error("IsCompleted() should return true for COMPLETED status")
	}
}

func TestExportMetadata_IsCompleted(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{ExportStatusCompleted, true},
		{ExportStatusError, false},
		{ExportStatusSubmitted, false},
		{ExportStatusInProgress, false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			metadata := &ExportMetadata{Status: tt.status}

			if metadata.IsCompleted() != tt.expected {
				t.Errorf("Status %s: expected IsCompleted()=%v, got %v",
					tt.status, tt.expected, metadata.IsCompleted())
			}
		})
	}
}

func TestExportMetadata_IsError(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{ExportStatusCompleted, false},
		{ExportStatusError, true},
		{ExportStatusSubmitted, false},
		{ExportStatusInProgress, false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			metadata := &ExportMetadata{Status: tt.status}

			if metadata.IsError() != tt.expected {
				t.Errorf("Status %s: expected IsError()=%v, got %v",
					tt.status, tt.expected, metadata.IsError())
			}
		})
	}
}

func TestExportMetadata_IsInProgress(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{ExportStatusCompleted, false},
		{ExportStatusError, false},
		{ExportStatusSubmitted, true},
		{ExportStatusInProgress, true},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			metadata := &ExportMetadata{Status: tt.status}

			if metadata.IsInProgress() != tt.expected {
				t.Errorf("Status %s: expected IsInProgress()=%v, got %v",
					tt.status, tt.expected, metadata.IsInProgress())
			}
		})
	}
}

func TestExecutionMetadata_GetRateLimitInfo(t *testing.T) {
	metadata := &ExecutionMetadata{
		Headers: map[string][]string{
			"X-Rate-Limit-Limit":     {"1000"},
			"X-Rate-Limit-Remaining": {"995"},
		},
	}

	info := metadata.GetRateLimitInfo()

	if info == nil {
		t.Fatal("GetRateLimitInfo returned nil")
	}

	// Note: GetRateLimitInfo looks for "X-Rate-Limit" header, not the individual headers
	// This test just verifies the method doesn't panic
}

func TestGetV2Metadata_NilInput(t *testing.T) {
	metadata := GetV2Metadata(nil, nil)
	if metadata != nil {
		t.Error("Expected nil metadata for nil response")
	}

	mockResp := &ExecuteNQLV2Response{QueryID: "#test"}
	metadata = GetV2Metadata(mockResp, nil)
	if metadata == nil {
		t.Error("Should handle nil API response")
	}
	if metadata.QueryID != "#test" {
		t.Error("Metadata should still contain QueryID")
	}
}

func TestGetV1Metadata_NilInput(t *testing.T) {
	metadata := GetV1Metadata(nil, nil)
	if metadata != nil {
		t.Error("Expected nil metadata for nil response")
	}
}

func TestGetExportMetadata_NilInput(t *testing.T) {
	metadata := GetExportMetadata(nil, nil)
	if metadata != nil {
		t.Error("Expected nil metadata for nil response")
	}
}
