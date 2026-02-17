package nql

import (
	"fmt"
	"strings"
)

// ValidateExecuteRequest validates an NQL execute request
func ValidateExecuteRequest(req *ExecuteRequest) error {
	if req == nil {
		return fmt.Errorf("execute request cannot be nil")
	}

	if err := validateQueryID(req.QueryID); err != nil {
		return err
	}

	if req.Platform != "" && len(req.Platform) > MaxPlatformLength {
		return fmt.Errorf("platform exceeds maximum length of %d characters", MaxPlatformLength)
	}

	return nil
}

// ValidateExportRequest validates an NQL export request
func ValidateExportRequest(req *ExportRequest) error {
	if req == nil {
		return fmt.Errorf("export request cannot be nil")
	}

	if err := validateQueryID(req.QueryID); err != nil {
		return err
	}

	if req.Platform != "" && len(req.Platform) > MaxPlatformLength {
		return fmt.Errorf("platform exceeds maximum length of %d characters", MaxPlatformLength)
	}

	if req.Format != "" && req.Format != ExportFormatCSV && req.Format != ExportFormatJSON {
		return fmt.Errorf("format must be either 'csv' or 'json', got: %s", req.Format)
	}

	return nil
}

// ValidateExportID validates an export ID
func ValidateExportID(exportID string) error {
	if exportID == "" {
		return fmt.Errorf("export ID cannot be empty")
	}

	if len(exportID) > MaxQueryIDLength {
		return fmt.Errorf("export ID exceeds maximum length of %d characters", MaxQueryIDLength)
	}

	return nil
}

// validateQueryID validates a query ID
func validateQueryID(queryID string) error {
	if queryID == "" {
		return fmt.Errorf("query ID is required")
	}

	if !strings.HasPrefix(queryID, "#") {
		return fmt.Errorf("query ID must start with '#', got: %s", queryID)
	}

	if len(queryID) > MaxQueryIDLength {
		return fmt.Errorf("query ID exceeds maximum length of %d characters", MaxQueryIDLength)
	}

	return nil
}
