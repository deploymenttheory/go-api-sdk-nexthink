package nql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateExecuteRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *ExecuteRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with query ID only",
			req: &ExecuteRequest{
				QueryID: "#test_query",
			},
			wantErr: false,
		},
		{
			name: "valid request with platform",
			req: &ExecuteRequest{
				QueryID:  "#test_query",
				Platform: "windows",
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "execute request cannot be nil",
		},
		{
			name: "empty query ID",
			req: &ExecuteRequest{
				QueryID: "",
			},
			wantErr: true,
			errMsg:  "query ID is required",
		},
		{
			name: "query ID without hash prefix",
			req: &ExecuteRequest{
				QueryID: "test_query",
			},
			wantErr: true,
			errMsg:  "query ID must start with '#'",
		},
		{
			name: "query ID too long",
			req: &ExecuteRequest{
				QueryID: "#" + strings.Repeat("a", MaxQueryIDLength),
			},
			wantErr: true,
			errMsg:  "query ID exceeds maximum length",
		},
		{
			name: "platform too long",
			req: &ExecuteRequest{
				QueryID:  "#test_query",
				Platform: strings.Repeat("a", MaxPlatformLength+1),
			},
			wantErr: true,
			errMsg:  "platform exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExecuteRequest(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateExportRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *ExportRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with query ID only",
			req: &ExportRequest{
				QueryID: "#test_query",
			},
			wantErr: false,
		},
		{
			name: "valid request with all fields",
			req: &ExportRequest{
				QueryID:  "#test_query",
				Platform: "windows",
				Format:   ExportFormatCSV,
			},
			wantErr: false,
		},
		{
			name: "valid request with JSON format",
			req: &ExportRequest{
				QueryID: "#test_query",
				Format:  ExportFormatJSON,
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "export request cannot be nil",
		},
		{
			name: "empty query ID",
			req: &ExportRequest{
				QueryID: "",
			},
			wantErr: true,
			errMsg:  "query ID is required",
		},
		{
			name: "query ID without hash prefix",
			req: &ExportRequest{
				QueryID: "test_query",
			},
			wantErr: true,
			errMsg:  "query ID must start with '#'",
		},
		{
			name: "invalid format",
			req: &ExportRequest{
				QueryID: "#test_query",
				Format:  "xml",
			},
			wantErr: true,
			errMsg:  "format must be either 'csv' or 'json'",
		},
		{
			name: "query ID too long",
			req: &ExportRequest{
				QueryID: "#" + strings.Repeat("a", MaxQueryIDLength),
			},
			wantErr: true,
			errMsg:  "query ID exceeds maximum length",
		},
		{
			name: "platform too long",
			req: &ExportRequest{
				QueryID:  "#test_query",
				Platform: strings.Repeat("a", MaxPlatformLength+1),
			},
			wantErr: true,
			errMsg:  "platform exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExportRequest(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateExportID(t *testing.T) {
	tests := []struct {
		name     string
		exportID string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid export ID",
			exportID: "abc123-def456-ghi789",
			wantErr:  false,
		},
		{
			name:     "empty export ID",
			exportID: "",
			wantErr:  true,
			errMsg:   "export ID cannot be empty",
		},
		{
			name:     "export ID too long",
			exportID: strings.Repeat("a", MaxQueryIDLength+1),
			wantErr:  true,
			errMsg:   "export ID exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExportID(tt.exportID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateQueryID(t *testing.T) {
	tests := []struct {
		name    string
		queryID string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid query ID",
			queryID: "#test_query",
			wantErr: false,
		},
		{
			name:    "valid query ID with underscores",
			queryID: "#get_pilot_collector_devices",
			wantErr: false,
		},
		{
			name:    "valid query ID with numbers",
			queryID: "#query123",
			wantErr: false,
		},
		{
			name:    "empty query ID",
			queryID: "",
			wantErr: true,
			errMsg:  "query ID is required",
		},
		{
			name:    "query ID without hash",
			queryID: "test_query",
			wantErr: true,
			errMsg:  "query ID must start with '#'",
		},
		{
			name:    "query ID too long",
			queryID: "#" + strings.Repeat("a", MaxQueryIDLength),
			wantErr: true,
			errMsg:  "query ID exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateQueryID(tt.queryID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
