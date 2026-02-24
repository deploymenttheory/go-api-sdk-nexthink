package nql

import (
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

// Response metadata helpers for NQL query execution
// Provides detailed information about query execution and response characteristics

// =============================================================================
// Execution Metadata
// =============================================================================

// ExecutionMetadata provides detailed execution information
type ExecutionMetadata struct {
	// Query Information
	QueryID       string
	ExecutedQuery string
	
	// Result Information
	RowsReturned int64
	
	// Timing Information
	ExecutionTime    time.Time
	ResponseDuration time.Duration
	
	// Response Characteristics
	ResponseSize   int64
	ResponseStatus int
	
	// Headers (rate limits, etc.)
	Headers map[string][]string
}

// GetV1Metadata extracts metadata from V1 execution response
func GetV1Metadata(resp *ExecuteNQLV1Response, apiResp *interfaces.Response) *ExecutionMetadata {
	if resp == nil {
		return nil
	}
	
	metadata := &ExecutionMetadata{
		QueryID:       resp.QueryID,
		ExecutedQuery: resp.ExecutedQuery,
		RowsReturned:  resp.Rows,
	}
	
	// Parse execution datetime from V1 response
	if resp.ExecutionDateTime != nil {
		dt := resp.ExecutionDateTime
		metadata.ExecutionTime = time.Date(
			int(dt.Year), time.Month(dt.Month), int(dt.Day),
			int(dt.Hour), int(dt.Minute), int(dt.Second),
			0, time.UTC,
		)
	}
	
	// Add API response metadata if available
	if apiResp != nil {
		metadata.ResponseDuration = apiResp.Duration
		metadata.ResponseSize = apiResp.Size
		metadata.ResponseStatus = apiResp.StatusCode
		metadata.Headers = apiResp.Headers
	}
	
	return metadata
}

// GetV2Metadata extracts metadata from V2 execution response
func GetV2Metadata(resp *ExecuteNQLV2Response, apiResp *interfaces.Response) *ExecutionMetadata {
	if resp == nil {
		return nil
	}
	
	metadata := &ExecutionMetadata{
		QueryID:       resp.QueryID,
		ExecutedQuery: resp.ExecutedQuery,
		RowsReturned:  resp.Rows,
	}
	
	// Parse execution datetime from V2 response (ISO format)
	if resp.ExecutionDateTime != "" {
		// Try parsing as ISO 8601 format
		t, err := time.Parse(time.RFC3339, resp.ExecutionDateTime)
		if err == nil {
			metadata.ExecutionTime = t
		} else {
			// Try alternative format
			t, err = time.Parse("2006-01-02T15:04:05", resp.ExecutionDateTime)
			if err == nil {
				metadata.ExecutionTime = t
			}
		}
	}
	
	// Add API response metadata if available
	if apiResp != nil {
		metadata.ResponseDuration = apiResp.Duration
		metadata.ResponseSize = apiResp.Size
		metadata.ResponseStatus = apiResp.StatusCode
		metadata.Headers = apiResp.Headers
	}
	
	return metadata
}

// GetExportMetadata extracts metadata from export status response
func GetExportMetadata(resp *NQLExportStatusResponse, apiResp *interfaces.Response) *ExportMetadata {
	if resp == nil {
		return nil
	}
	
	metadata := &ExportMetadata{
		ExportID:         resp.ExportID,
		Status:           resp.Status,
		ResultsFileURL:   resp.ResultsFileURL,
		ErrorDescription: resp.ErrorDescription,
	}
	
	// Add API response metadata if available
	if apiResp != nil {
		metadata.ResponseDuration = apiResp.Duration
		metadata.ResponseStatus = apiResp.StatusCode
		metadata.ReceivedAt = apiResp.ReceivedAt
		metadata.Headers = apiResp.Headers
	}
	
	return metadata
}

// =============================================================================
// Export Metadata
// =============================================================================

// ExportMetadata provides detailed export operation information
type ExportMetadata struct {
	// Export Information
	ExportID         string
	Status           string
	ResultsFileURL   string
	ErrorDescription string
	
	// Timing Information
	ReceivedAt       time.Time
	ResponseDuration time.Duration
	
	// Response Characteristics
	ResponseStatus int
	
	// Headers
	Headers map[string][]string
}

// IsCompleted checks if the export is completed
func (em *ExportMetadata) IsCompleted() bool {
	return em.Status == ExportStatusCompleted
}

// IsError checks if the export encountered an error
func (em *ExportMetadata) IsError() bool {
	return em.Status == ExportStatusError
}

// IsInProgress checks if the export is still in progress
func (em *ExportMetadata) IsInProgress() bool {
	return em.Status == ExportStatusInProgress || em.Status == ExportStatusSubmitted
}

// =============================================================================
// Rate Limit Information
// =============================================================================

// RateLimitInfo contains rate limit information from response headers
type RateLimitInfo struct {
	Limit     string
	Remaining string
	Reset     string
	RetryAfter string
}

// GetRateLimitInfo extracts rate limit information from metadata
func (em *ExecutionMetadata) GetRateLimitInfo() *RateLimitInfo {
	if em.Headers == nil {
		return nil
	}
	
	info := &RateLimitInfo{}
	
	if limit := em.Headers["X-Rate-Limit"]; len(limit) > 0 {
		info.Limit = limit[0]
	}
	
	if remaining := em.Headers["X-Rate-Limit-Remaining"]; len(remaining) > 0 {
		info.Remaining = remaining[0]
	}
	
	if reset := em.Headers["X-Rate-Limit-Reset"]; len(reset) > 0 {
		info.Reset = reset[0]
	}
	
	if retryAfter := em.Headers["Retry-After"]; len(retryAfter) > 0 {
		info.RetryAfter = retryAfter[0]
	}
	
	return info
}

// =============================================================================
// Helper Methods
// =============================================================================

// TimeSinceExecution calculates the time since the query was executed
func (em *ExecutionMetadata) TimeSinceExecution() time.Duration {
	if em.ExecutionTime.IsZero() {
		return 0
	}
	return time.Since(em.ExecutionTime)
}

// String returns a string representation of the metadata
func (em *ExecutionMetadata) String() string {
	return formatMetadata(em)
}

// String returns a string representation of the export metadata
func (exm *ExportMetadata) String() string {
	return formatExportMetadata(exm)
}

// formatMetadata formats execution metadata as a string
func formatMetadata(em *ExecutionMetadata) string {
	if em == nil {
		return "<nil>"
	}
	
	result := "ExecutionMetadata{\n"
	result += "  QueryID: " + em.QueryID + "\n"
	result += "  ExecutedQuery: " + em.ExecutedQuery + "\n"
	result += "  RowsReturned: " + formatInt64(em.RowsReturned) + "\n"
	result += "  ExecutionTime: " + em.ExecutionTime.Format(time.RFC3339) + "\n"
	result += "  ResponseDuration: " + em.ResponseDuration.String() + "\n"
	result += "  ResponseSize: " + formatInt64(em.ResponseSize) + " bytes\n"
	result += "  ResponseStatus: " + formatInt(em.ResponseStatus) + "\n"
	result += "}"
	
	return result
}

// formatExportMetadata formats export metadata as a string
func formatExportMetadata(exm *ExportMetadata) string {
	if exm == nil {
		return "<nil>"
	}
	
	result := "ExportMetadata{\n"
	result += "  ExportID: " + exm.ExportID + "\n"
	result += "  Status: " + exm.Status + "\n"
	if exm.ResultsFileURL != "" {
		result += "  ResultsFileURL: " + exm.ResultsFileURL + "\n"
	}
	if exm.ErrorDescription != "" {
		result += "  ErrorDescription: " + exm.ErrorDescription + "\n"
	}
	result += "  ReceivedAt: " + exm.ReceivedAt.Format(time.RFC3339) + "\n"
	result += "  ResponseDuration: " + exm.ResponseDuration.String() + "\n"
	result += "  ResponseStatus: " + formatInt(exm.ResponseStatus) + "\n"
	result += "}"
	
	return result
}

// Helper functions for formatting
func formatInt64(i int64) string {
	return string(rune(i + '0'))
}

func formatInt(i int) string {
	return string(rune(i + '0'))
}
