package nql

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

type (
	// NQLServiceInterface defines the interface for NQL operations
	//
	// Nexthink API docs: https://docs.nexthink.com/api/nql
	NQLServiceInterface interface {
		// ExecuteNQLV1 executes an NQL query synchronously using API V1
		//
		// Executes a pre-configured NQL query and returns the results immediately.
		// Optimized for relatively small requests at high frequency.
		//
		// The query must be pre-created in Nexthink admin (admin/NQL API queries)
		// and identified by its Query ID (format: #query_name).
		//
		// V1 Response contains:
		//  - QueryID: Identifier of the executed query
		//  - ExecutedQuery: Final query with replaced parameters
		//  - Rows: Number of rows returned
		//  - ExecutionDateTime: DateTime object with execution timestamp
		//  - Headers: Ordered list of column headers
		//  - Data: Array of arrays containing row data
		//
		// Use this for:
		//  - Real-time queries with small result sets
		//  - Interactive dashboards
		//  - Frequent polling operations
		//
		// Nexthink API docs: https://docs.nexthink.com/api/nql/execute-an-nql
		ExecuteNQLV1(ctx context.Context, req *ExecuteRequest) (*ExecuteNQLV1Response, *interfaces.Response, error)

		// ExecuteNQLV2 executes an NQL query synchronously using API V2
		//
		// Executes a pre-configured NQL query and returns the results immediately.
		// Optimized for relatively small requests at high frequency.
		//
		// The query must be pre-created in Nexthink admin (admin/NQL API queries)
		// and identified by its Query ID (format: #query_name).
		//
		// V2 Response contains:
		//  - QueryID: Identifier of the executed query
		//  - ExecutedQuery: Final query with replaced parameters
		//  - Rows: Number of rows returned
		//  - ExecutionDateTime: ISO format string with execution timestamp
		//  - Data: Array of objects (map[string]any) with key-value pairs
		//
		// V2 provides cleaner structured data compared to V1.
		//
		// Use this for:
		//  - Real-time queries with small result sets
		//  - Interactive dashboards
		//  - Frequent polling operations
		//
		// Nexthink API docs: https://docs.nexthink.com/api/nql/execute-an-nql
		ExecuteNQLV2(ctx context.Context, req *ExecuteRequest) (*ExecuteNQLV2Response, *interfaces.Response, error)

		// StartNQLExport starts an asynchronous NQL export
		//
		// Initiates an export for a pre-configured NQL query.
		// Optimized for large queries at low frequency.
		//
		// The query must be pre-created in Nexthink admin (admin/NQL API queries)
		// and identified by its Query ID (format: #query_name).
		//
		// This is an asynchronous operation that returns an exportID immediately.
		// Use GetNQLExportStatus() to check completion and get the download URL.
		//
		// Use this for:
		//  - Large data extracts
		//  - Scheduled reports
		//  - Bulk data exports
		//
		// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql
		StartNQLExport(ctx context.Context, req *ExportRequest) (*StartNQLExportResponse, *interfaces.Response, error)

		// GetNQLExportStatus checks the status of an export operation
		//
		// Retrieves the current status of an export initiated with StartNQLExport().
		//
		// Status values:
		//  - SUBMITTED: Export is queued
		//  - IN_PROGRESS: Export is currently running
		//  - COMPLETED: Export is ready (ResultsFileURL will be available)
		//  - ERROR: Export failed (ErrorDescription will contain error details)
		//
		// When status is COMPLETED, the response includes a ResultsFileURL (S3 URL)
		// that can be used with DownloadNQLExport().
		//
		// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql#status-of-an-export
		GetNQLExportStatus(ctx context.Context, exportID string) (*NQLExportStatusResponse, *interfaces.Response, error)

		// DownloadNQLExport downloads a completed export from the S3 URL
		//
		// Downloads the export data from the S3 URL provided in GetNQLExportStatus().
		// Only works when export status is COMPLETED.
		//
		// The download URL typically expires after a certain time period.
		//
		// Returns the raw export data as bytes (CSV or JSON format).
		DownloadNQLExport(ctx context.Context, downloadURL string) ([]byte, error)

		// WaitForNQLExport polls the export status until it completes or fails
		//
		// This is a convenience method that polls GetNQLExportStatus() at regular intervals
		// until the export reaches a terminal state (COMPLETED or ERROR).
		//
		// Parameters:
		//  - ctx: Context for cancellation
		//  - exportID: The export ID from StartNQLExport()
		//  - pollInterval: How often to check status (recommended: 5-10 seconds)
		//  - timeout: Maximum time to wait (recommended: 5-10 minutes)
		//
		// Returns the final status response when export completes or an error if it fails.
		WaitForNQLExport(ctx context.Context, exportID string, pollInterval, timeout time.Duration) (*NQLExportStatusResponse, error)
	}

	// Service implements the NQLServiceInterface
	Service struct {
		client interfaces.HTTPClient
	}
)

var _ NQLServiceInterface = (*Service)(nil)

// NewService creates a new NQL service instance
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// =============================================================================
// NQL Execute Operations
// =============================================================================

// ExecuteNQLV1 executes an NQL query synchronously using API V1
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/nql/execute
// Nexthink API docs: https://docs.nexthink.com/api/nql/execute-an-nql
func (s *Service) ExecuteNQLV1(ctx context.Context, req *ExecuteRequest) (*ExecuteNQLV1Response, *interfaces.Response, error) {
	if err := ValidateExecuteRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointNqlExecuteV1

	headers := map[string]string{
		"Accept":       "application/json, text/csv",
		"Content-Type": "application/json",
	}

	var result ExecuteNQLV1Response
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ExecuteNQLV2 executes an NQL query synchronously using API V2
// URL: POST https://instance.api.region.nexthink.cloud/api/v2/nql/execute
// Nexthink API docs: https://docs.nexthink.com/api/nql/execute-an-nql
func (s *Service) ExecuteNQLV2(ctx context.Context, req *ExecuteRequest) (*ExecuteNQLV2Response, *interfaces.Response, error) {
	if err := ValidateExecuteRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointNqlExecuteV2

	headers := map[string]string{
		"Accept":       "application/json, text/csv",
		"Content-Type": "application/json",
	}

	var result ExecuteNQLV2Response
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// =============================================================================
// NQL Export Operations
// =============================================================================

// StartNQLExport starts an asynchronous NQL export
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/nql/export
// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql
func (s *Service) StartNQLExport(ctx context.Context, req *ExportRequest) (*StartNQLExportResponse, *interfaces.Response, error) {
	if err := ValidateExportRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointNqlExport

	headers := map[string]string{
		"Accept":       "application/json, text/csv",
		"Content-Type": "application/json",
	}

	var result StartNQLExportResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetNQLExportStatus checks the status of an export operation
// URL: GET https://instance.api.region.nexthink.cloud/api/v1/nql/status/{exportId}
// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql#status-of-an-export
func (s *Service) GetNQLExportStatus(ctx context.Context, exportID string) (*NQLExportStatusResponse, *interfaces.Response, error) {
	if err := ValidateExportID(exportID); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", EndpointNqlStatus, exportID)

	headers := map[string]string{
		"Accept": "application/json, text/csv",
	}

	var result NQLExportStatusResponse
	resp, err := s.client.Get(ctx, endpoint, nil, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DownloadNQLExport downloads a completed export from an S3 pre-signed URL.
//
// Note: This method uses a standard HTTP client (not the SDK transport) because:
//   - S3 URLs are external to the Nexthink API
//   - They don't require Nexthink authentication
//   - They're pre-signed with temporary credentials from AWS
//   - The download is a simple GET request to AWS S3
//
// The HTTP client is configured with a 5-minute timeout for large downloads.
func (s *Service) DownloadNQLExport(ctx context.Context, downloadURL string) ([]byte, error) {
	if downloadURL == "" {
		return nil, fmt.Errorf("download URL cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download export: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read export data: %w", err)
	}

	return data, nil
}

// WaitForNQLExport polls the export status until it completes or fails
func (s *Service) WaitForNQLExport(ctx context.Context, exportID string, pollInterval, timeout time.Duration) (*NQLExportStatusResponse, error) {
	if err := ValidateExportID(exportID); err != nil {
		return nil, err
	}

	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}

	if timeout <= 0 {
		timeout = 10 * time.Minute
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	status, _, err := s.GetNQLExportStatus(timeoutCtx, exportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial export status: %w", err)
	}

	if isTerminalStatus(status.Status) {
		return status, nil
	}

	for {
		select {
		case <-timeoutCtx.Done():
			return status, fmt.Errorf("timeout waiting for export to complete after %v", timeout)

		case <-ticker.C:
			status, _, err = s.GetNQLExportStatus(timeoutCtx, exportID)
			if err != nil {
				return nil, fmt.Errorf("failed to get export status: %w", err)
			}

			if isTerminalStatus(status.Status) {
				return status, nil
			}
		}
	}
}
