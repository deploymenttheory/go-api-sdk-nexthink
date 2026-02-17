package nql

// ExecuteRequest represents an NQL query execution request.
// Optimized for relatively small requests at high frequency.
type ExecuteRequest struct {
	// QueryID is the NQL query identifier (must be pre-created in Nexthink admin)
	// Format: #query_name
	QueryID string `json:"queryId"`

	// Platform optionally specifies the platform for the query
	Platform string `json:"platform,omitempty"`
}

// ExecuteNQLV1Response represents the response from an NQL execute V1 operation
// API docs: https://docs.nexthink.com/api/nql/execute-an-nql
type ExecuteNQLV1Response struct {
	// QueryID is the identifier of the executed query
	QueryID string `json:"queryId,omitempty"`

	// ExecutedQuery is the final query executed with replaced parameters
	ExecutedQuery string `json:"executedQuery,omitempty"`

	// Rows is the number of rows returned
	Rows int64 `json:"rows,omitempty"`

	// ExecutionDateTime is the date and time of execution
	ExecutionDateTime *DateTime `json:"executionDateTime,omitempty"`

	// Headers is the ordered list with the headers of the returned fields
	Headers []string `json:"headers,omitempty"`

	// Data is the list of rows with the data returned by the query execution
	Data [][]any `json:"data,omitempty"`
}

// ExecuteNQLV2Response represents the response from an NQL execute V2 operation
// API docs: https://docs.nexthink.com/api/nql/execute-an-nql
type ExecuteNQLV2Response struct {
	// QueryID is the identifier of the executed query
	QueryID string `json:"queryId,omitempty"`

	// ExecutedQuery is the final query executed with replaced parameters
	ExecutedQuery string `json:"executedQuery,omitempty"`

	// Rows is the number of rows returned
	Rows int64 `json:"rows,omitempty"`

	// ExecutionDateTime is the date and time of execution in ISO format
	ExecutionDateTime string `json:"executionDateTime,omitempty"`

	// Data is the list of rows with the data returned by the query execution
	// Each row is an object with key-value pairs
	Data []map[string]any `json:"data,omitempty"`
}

// DateTime represents a date and time object from the Nexthink API
type DateTime struct {
	Year   int64 `json:"year,omitempty"`
	Month  int64 `json:"month,omitempty"`
	Day    int64 `json:"day,omitempty"`
	Hour   int64 `json:"hour,omitempty"`
	Minute int64 `json:"minute,omitempty"`
	Second int64 `json:"second,omitempty"`
}

// ExportRequest represents an NQL export request.
// Optimized for large queries at low frequency.
// This is an asynchronous operation - use GetExportStatus to check completion.
type ExportRequest struct {
	// QueryID is the NQL query identifier (must be pre-created in Nexthink admin)
	// Format: #query_name
	QueryID string `json:"queryId"`

	// Platform optionally specifies the platform for the query
	Platform string `json:"platform,omitempty"`

	// Format specifies the export format (csv or json)
	// Defaults to csv if not specified
	Format string `json:"format,omitempty"`
}

// StartNQLExportResponse represents the initial response from starting an export
type StartNQLExportResponse struct {
	// ExportID is the unique identifier for this export operation
	// Use this ID to check status and download the results
	ExportID string `json:"exportId"`

	// Status is the current status of the export
	Status string `json:"status"`

	// Message provides additional information about the export
	Message string `json:"message,omitempty"`
}

// NQLExportStatusResponse represents the status of an export operation
type NQLExportStatusResponse struct {
	// ExportID is the unique identifier for this export operation (not present in Python lib but useful)
	ExportID string `json:"exportId,omitempty"`

	// Status is the current status of the export
	// Possible values: SUBMITTED, IN_PROGRESS, COMPLETED, ERROR
	Status string `json:"status"`

	// ResultsFileURL is the S3 URL to download the export (available when status is COMPLETED)
	ResultsFileURL string `json:"resultsFileUrl,omitempty"`

	// ErrorDescription provides error details when status is ERROR
	ErrorDescription string `json:"errorDescription,omitempty"`
}

// ErrorResponse represents an error response from the NQL API
type ErrorResponse struct {
	// Error contains the error details
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	// Code is the error code
	Code string `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Details provides additional error context
	Details map[string]any `json:"details,omitempty"`
}
