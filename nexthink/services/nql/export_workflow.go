package nql

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Export workflow helpers simplify the asynchronous export process
// Provides high-level methods that handle the entire export lifecycle

// =============================================================================
// Export Options
// =============================================================================

// ExportOptions configures export workflow behavior
type ExportOptions struct {
	// Format specifies the export format ("csv" or "json")
	// Defaults to "csv" if not specified
	Format string
	
	// PollInterval is how often to check export status
	// Defaults to 5 seconds if not specified
	PollInterval time.Duration
	
	// Timeout is the maximum time to wait for export completion
	// Defaults to 10 minutes if not specified
	Timeout time.Duration
	
	// OnProgress is an optional callback for progress updates
	// Called each time the status is checked
	OnProgress func(status string, elapsedTime time.Duration)
	
	// OnStatusChange is an optional callback that fires when status changes
	// Called only when the export status transitions to a new state
	OnStatusChange func(oldStatus, newStatus string, elapsedTime time.Duration)
}

// DefaultExportOptions returns export options with sensible defaults
func DefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		Format:       ExportFormatCSV,
		PollInterval: 5 * time.Second,
		Timeout:      10 * time.Minute,
	}
}

// WithFormat sets the export format
func (opts *ExportOptions) WithFormat(format string) *ExportOptions {
	opts.Format = format
	return opts
}

// WithPollInterval sets the poll interval
func (opts *ExportOptions) WithPollInterval(interval time.Duration) *ExportOptions {
	opts.PollInterval = interval
	return opts
}

// WithTimeout sets the timeout
func (opts *ExportOptions) WithTimeout(timeout time.Duration) *ExportOptions {
	opts.Timeout = timeout
	return opts
}

// WithOnProgress sets the progress callback
func (opts *ExportOptions) WithOnProgress(callback func(status string, elapsedTime time.Duration)) *ExportOptions {
	opts.OnProgress = callback
	return opts
}

// WithOnStatusChange sets the status change callback
func (opts *ExportOptions) WithOnStatusChange(callback func(oldStatus, newStatus string, elapsedTime time.Duration)) *ExportOptions {
	opts.OnStatusChange = callback
	return opts
}

// =============================================================================
// Export Result
// =============================================================================

// ExportResult contains the complete result of an export operation
type ExportResult struct {
	// ExportID is the unique identifier for the export
	ExportID string
	
	// Data contains the exported data
	Data []byte
	
	// Format is the format of the data ("csv" or "json")
	Format string
	
	// Metadata contains execution metadata
	Metadata *ExportMetadata
	
	// TotalDuration is the total time from start to completion
	TotalDuration time.Duration
	
	// PollCount is the number of times status was polled
	PollCount int
}

// Size returns the size of the exported data in bytes
func (er *ExportResult) Size() int64 {
	return int64(len(er.Data))
}

// SizeFormatted returns a human-readable size string
func (er *ExportResult) SizeFormatted() string {
	bytes := er.Size()
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// =============================================================================
// Export Workflow
// =============================================================================

// ExportWorkflow executes the complete export workflow
// This is the main convenience method that handles:
// 1. Starting the export
// 2. Polling for completion
// 3. Downloading the result
func (s *Service) ExportWorkflow(ctx context.Context, req *ExportRequest, opts *ExportOptions) (*ExportResult, error) {
	// Use default options if none provided
	if opts == nil {
		opts = DefaultExportOptions()
	}
	
	// Apply defaults for unset options
	if opts.Format == "" {
		opts.Format = ExportFormatCSV
	}
	if opts.PollInterval <= 0 {
		opts.PollInterval = 5 * time.Second
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 10 * time.Minute
	}
	
	// Set format in request
	if req.Format == "" {
		req.Format = opts.Format
	}
	
	// Start timing
	startTime := time.Now()
	
	// Step 1: Start the export
	s.client.GetLogger().Info("Starting NQL export",
		zap.String("query_id", req.QueryID),
		zap.String("format", req.Format))
	
	startResp, _, err := s.StartNQLExport(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start export: %w", err)
	}
	
	exportID := startResp.ExportID
	s.client.GetLogger().Info("Export started",
		zap.String("export_id", exportID),
		zap.String("initial_status", startResp.Status))
	
	// Step 2: Wait for completion with progress callbacks
	lastStatus := startResp.Status
	pollCount := 0
	
	finalStatus, err := s.waitForExportWithCallbacks(ctx, exportID, opts, &lastStatus, &pollCount, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for export: %w", err)
	}
	
	// Step 3: Download the result
	if finalStatus.ResultsFileURL == "" {
		return nil, fmt.Errorf("export completed but no download URL provided")
	}
	
	s.client.GetLogger().Info("Downloading export data",
		zap.String("export_id", exportID))
	
	data, err := s.DownloadNQLExport(ctx, finalStatus.ResultsFileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download export: %w", err)
	}
	
	totalDuration := time.Since(startTime)
	
	s.client.GetLogger().Info("Export workflow completed successfully",
		zap.String("export_id", exportID),
		zap.Int64("data_size", int64(len(data))),
		zap.Duration("total_duration", totalDuration),
		zap.Int("poll_count", pollCount))
	
	return &ExportResult{
		ExportID:      exportID,
		Data:          data,
		Format:        req.Format,
		TotalDuration: totalDuration,
		PollCount:     pollCount,
	}, nil
}

// waitForExportWithCallbacks polls for export completion with progress callbacks
func (s *Service) waitForExportWithCallbacks(
	ctx context.Context,
	exportID string,
	opts *ExportOptions,
	lastStatus *string,
	pollCount *int,
	startTime time.Time,
) (*NQLExportStatusResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()
	
	ticker := time.NewTicker(opts.PollInterval)
	defer ticker.Stop()
	
	// Check initial status
	status, _, err := s.GetNQLExportStatus(timeoutCtx, exportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial export status: %w", err)
	}
	*pollCount++
	
	// Initial progress callback
	elapsedTime := time.Since(startTime)
	if opts.OnProgress != nil {
		opts.OnProgress(status.Status, elapsedTime)
	}
	
	// Check for status change
	if opts.OnStatusChange != nil && status.Status != *lastStatus {
		opts.OnStatusChange(*lastStatus, status.Status, elapsedTime)
		*lastStatus = status.Status
	}
	
	// Check if already completed
	if isTerminalStatus(status.Status) {
		if status.Status == ExportStatusError {
			return nil, fmt.Errorf("export failed: %s", status.ErrorDescription)
		}
		return status, nil
	}
	
	// Poll until terminal status
	for {
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("timeout waiting for export to complete after %v (status: %s)", opts.Timeout, status.Status)
			
		case <-ticker.C:
			status, _, err = s.GetNQLExportStatus(timeoutCtx, exportID)
			if err != nil {
				return nil, fmt.Errorf("failed to get export status: %w", err)
			}
			*pollCount++
			
			elapsedTime := time.Since(startTime)
			
			// Progress callback
			if opts.OnProgress != nil {
				opts.OnProgress(status.Status, elapsedTime)
			}
			
			// Status change callback
			if opts.OnStatusChange != nil && status.Status != *lastStatus {
				opts.OnStatusChange(*lastStatus, status.Status, elapsedTime)
				*lastStatus = status.Status
			}
			
			// Check for terminal status
			if isTerminalStatus(status.Status) {
				if status.Status == ExportStatusError {
					return nil, fmt.Errorf("export failed: %s", status.ErrorDescription)
				}
				return status, nil
			}
		}
	}
}

// =============================================================================
// Simplified Export Methods
// =============================================================================

// ExportToCSV is a convenience method that exports a query to CSV format
func (s *Service) ExportToCSV(ctx context.Context, queryID string) (*ExportResult, error) {
	return s.ExportWorkflow(ctx, &ExportRequest{
		QueryID: queryID,
		Format:  ExportFormatCSV,
	}, nil)
}

// ExportToJSON is a convenience method that exports a query to JSON format
func (s *Service) ExportToJSON(ctx context.Context, queryID string) (*ExportResult, error) {
	return s.ExportWorkflow(ctx, &ExportRequest{
		QueryID: queryID,
		Format:  ExportFormatJSON,
	}, nil)
}

// ExportWithProgress exports a query with a simple progress callback
func (s *Service) ExportWithProgress(ctx context.Context, queryID, format string, progressFn func(status string)) (*ExportResult, error) {
	opts := DefaultExportOptions().
		WithFormat(format).
		WithOnProgress(func(status string, elapsed time.Duration) {
			if progressFn != nil {
				progressFn(status)
			}
		})
	
	return s.ExportWorkflow(ctx, &ExportRequest{
		QueryID: queryID,
		Format:  format,
	}, opts)
}

// =============================================================================
// Export Status Helpers
// =============================================================================

// IsExportReady checks if an export is ready for download
func (s *Service) IsExportReady(ctx context.Context, exportID string) (bool, error) {
	status, _, err := s.GetNQLExportStatus(ctx, exportID)
	if err != nil {
		return false, err
	}
	
	return status.Status == ExportStatusCompleted, nil
}

// GetExportProgress returns human-readable progress information
func (s *Service) GetExportProgress(ctx context.Context, exportID string) (string, error) {
	status, _, err := s.GetNQLExportStatus(ctx, exportID)
	if err != nil {
		return "", err
	}
	
	switch status.Status {
	case ExportStatusSubmitted:
		return "Export queued and waiting to start", nil
	case ExportStatusInProgress:
		return "Export in progress", nil
	case ExportStatusCompleted:
		return "Export completed successfully", nil
	case ExportStatusError:
		return fmt.Sprintf("Export failed: %s", status.ErrorDescription), nil
	default:
		return fmt.Sprintf("Unknown status: %s", status.Status), nil
	}
}
