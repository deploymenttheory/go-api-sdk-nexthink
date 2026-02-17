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
		// ExecuteV1 executes an NQL query synchronously using API V1
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
		ExecuteV1(ctx context.Context, req *ExecuteRequest) (*ExecuteV1Response, *interfaces.Response, error)

		// ExecuteV2 executes an NQL query synchronously using API V2
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
		ExecuteV2(ctx context.Context, req *ExecuteRequest) (*ExecuteV2Response, *interfaces.Response, error)

		// StartExport starts an asynchronous NQL export
		//
		// Initiates an export for a pre-configured NQL query.
		// Optimized for large queries at low frequency.
		//
		// The query must be pre-created in Nexthink admin (admin/NQL API queries)
		// and identified by its Query ID (format: #query_name).
		//
		// This is an asynchronous operation that returns an exportID immediately.
		// Use GetExportStatus() to check completion and get the download URL.
		//
		// Use this for:
		//  - Large data extracts
		//  - Scheduled reports
		//  - Bulk data exports
		//
		// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql
		StartExport(ctx context.Context, req *ExportRequest) (*ExportResponse, *interfaces.Response, error)

		// GetExportStatus checks the status of an export operation
		//
		// Retrieves the current status of an export initiated with StartExport().
		//
		// Status values:
		//  - SUBMITTED: Export is queued
		//  - IN_PROGRESS: Export is currently running
		//  - COMPLETED: Export is ready (ResultsFileURL will be available)
		//  - ERROR: Export failed (ErrorDescription will contain error details)
		//
		// When status is COMPLETED, the response includes a ResultsFileURL (S3 URL)
		// that can be used with DownloadExport().
		//
		// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql#status-of-an-export
		GetExportStatus(ctx context.Context, exportID string) (*ExportStatusResponse, *interfaces.Response, error)

		// DownloadExport downloads a completed export from the S3 URL
		//
		// Downloads the export data from the S3 URL provided in GetExportStatus().
		// Only works when export status is COMPLETED.
		//
		// The download URL typically expires after a certain time period.
		//
		// Returns the raw export data as bytes (CSV or JSON format).
		DownloadExport(ctx context.Context, downloadURL string) ([]byte, error)

		// WaitForExport polls the export status until it completes or fails
		//
		// This is a convenience method that polls GetExportStatus() at regular intervals
		// until the export reaches a terminal state (COMPLETED or ERROR).
		//
		// Parameters:
		//  - ctx: Context for cancellation
		//  - exportID: The export ID from StartExport()
		//  - pollInterval: How often to check status (recommended: 5-10 seconds)
		//  - timeout: Maximum time to wait (recommended: 5-10 minutes)
		//
		// Returns the final status response when export completes or an error if it fails.
		WaitForExport(ctx context.Context, exportID string, pollInterval, timeout time.Duration) (*ExportStatusResponse, error)
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

// ExecuteV1 executes an NQL query synchronously using API V1
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/nql/execute
// Nexthink API docs: https://docs.nexthink.com/api/nql/execute-an-nql
func (s *Service) ExecuteV1(ctx context.Context, req *ExecuteRequest) (*ExecuteV1Response, *interfaces.Response, error) {
	if err := ValidateExecuteRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointNqlExecuteV1

	headers := map[string]string{
		"Accept":       "application/json, text/csv",
		"Content-Type": "application/json",
	}

	var result ExecuteV1Response
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ExecuteV2 executes an NQL query synchronously using API V2
// URL: POST https://instance.api.region.nexthink.cloud/api/v2/nql/execute
// Nexthink API docs: https://docs.nexthink.com/api/nql/execute-an-nql
func (s *Service) ExecuteV2(ctx context.Context, req *ExecuteRequest) (*ExecuteV2Response, *interfaces.Response, error) {
	if err := ValidateExecuteRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointNqlExecuteV2

	headers := map[string]string{
		"Accept":       "application/json, text/csv",
		"Content-Type": "application/json",
	}

	var result ExecuteV2Response
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// =============================================================================
// NQL Export Operations
// =============================================================================

// StartExport starts an asynchronous NQL export
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/nql/export
// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql
func (s *Service) StartExport(ctx context.Context, req *ExportRequest) (*ExportResponse, *interfaces.Response, error) {
	if err := ValidateExportRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointNqlExport

	headers := map[string]string{
		"Accept":       "application/json, text/csv",
		"Content-Type": "application/json",
	}

	var result ExportResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetExportStatus checks the status of an export operation
// URL: GET https://instance.api.region.nexthink.cloud/api/v1/nql/status/{exportId}
// Nexthink API docs: https://docs.nexthink.com/api/nql/export-an-nql#status-of-an-export
func (s *Service) GetExportStatus(ctx context.Context, exportID string) (*ExportStatusResponse, *interfaces.Response, error) {
	if err := ValidateExportID(exportID); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", EndpointNqlStatus, exportID)

	headers := map[string]string{
		"Accept": "application/json, text/csv",
	}

	var result ExportStatusResponse
	resp, err := s.client.Get(ctx, endpoint, nil, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DownloadExport downloads a completed export from an S3 pre-signed URL.
//
// Note: This method uses a standard HTTP client (not the SDK transport) because:
//   - S3 URLs are external to the Nexthink API
//   - They don't require Nexthink authentication
//   - They're pre-signed with temporary credentials from AWS
//   - The download is a simple GET request to AWS S3
//
// The HTTP client is configured with a 5-minute timeout for large downloads.
func (s *Service) DownloadExport(ctx context.Context, downloadURL string) ([]byte, error) {
	if downloadURL == "" {
		return nil, fmt.Errorf("download URL cannot be empty")
	}

	// Create an HTTP request for the S3 download
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	// Use a standard HTTP client for S3 downloads (external to Nexthink API)
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

// WaitForExport polls the export status until it completes or fails
func (s *Service) WaitForExport(ctx context.Context, exportID string, pollInterval, timeout time.Duration) (*ExportStatusResponse, error) {
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

	status, _, err := s.GetExportStatus(timeoutCtx, exportID)
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
			status, _, err = s.GetExportStatus(timeoutCtx, exportID)
			if err != nil {
				return nil, fmt.Errorf("failed to get export status: %w", err)
			}

			if isTerminalStatus(status.Status) {
				return status, nil
			}
		}
	}
}
