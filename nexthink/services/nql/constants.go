package nql

const (
	// API Endpoints - Execute
	EndpointNqlExecuteV1 = "/api/v1/nql/execute"
	EndpointNqlExecuteV2 = "/api/v2/nql/execute"

	// API Endpoints - Export (V1 only)
	EndpointNqlExport = "/api/v1/nql/export"
	EndpointNqlStatus = "/api/v1/nql/status" // + /{exportId}

	// Export Status Values (matching Nexthink API)
	ExportStatusSubmitted  = "SUBMITTED"
	ExportStatusInProgress = "IN_PROGRESS"
	ExportStatusCompleted  = "COMPLETED"
	ExportStatusError      = "ERROR"

	// Export Format Values
	ExportFormatCSV  = "csv"
	ExportFormatJSON = "json"

	// Validation Constraints
	MaxQueryIDLength  = 256
	MaxPlatformLength = 50
	MaxExportIDLength = 256
)
