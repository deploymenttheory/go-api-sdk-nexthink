package nql

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a test client with httpmock activated
func setupMockClient(t *testing.T) (*Service, string) {
	t.Helper()

	logger := zap.NewNop()
	baseURL := "https://test.api.us.nexthink.cloud"
	tokenURL := baseURL + "/api/v1/token"

	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	
	// Setup cleanup
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock the OAuth token endpoint
	httpmock.RegisterResponder("POST", tokenURL,
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "mock-access-token",
			"expires_in":   3600,
			"token_type":   "Bearer",
		}))

	// Create transport with the mocked HTTP transport
	transport, err := client.NewTransport("client-id", "client-secret", "test-instance", "us",
		client.WithLogger(logger),
		client.WithBaseURL(baseURL),
		client.WithCustomTokenURL(tokenURL),
		client.WithTransport(httpClient.Transport),
	)
	require.NoError(t, err)

	return NewService(transport), baseURL
}

func TestExecuteNQLV1_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExecuteRequest{
		QueryID: "#test_query",
	}

	result, resp, err := service.ExecuteNQLV1(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "#test_query", result.QueryID)
	assert.Equal(t, int64(2), result.Rows)
	assert.Len(t, result.Headers, 3)
	assert.Len(t, result.Data, 2)
	assert.NotNil(t, result.ExecutionDateTime)
}

func TestExecuteNQLV1_WithPlatform(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExecuteRequest{
		QueryID:  "#test_query",
		Platform: "windows",
	}

	result, resp, err := service.ExecuteNQLV1(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "#test_query", result.QueryID)
}

func TestExecuteNQLV2_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExecuteRequest{
		QueryID: "#test_query",
	}

	result, resp, err := service.ExecuteNQLV2(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "#test_query", result.QueryID)
	assert.Equal(t, int64(2), result.Rows)
	assert.Len(t, result.Data, 2)
	assert.NotEmpty(t, result.ExecutionDateTime)
	
	// V2 returns data as objects (map[string]any)
	assert.IsType(t, map[string]any{}, result.Data[0])
}

func TestExecuteNQLV2_WithPlatform(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExecuteRequest{
		QueryID:  "#test_query",
		Platform: "windows",
	}

	result, resp, err := service.ExecuteNQLV2(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "#test_query", result.QueryID)
}

func TestExecuteNQLV1_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *ExecuteRequest
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "execute request cannot be nil",
		},
		{
			name: "empty query ID",
			req: &ExecuteRequest{
				QueryID: "",
			},
			errMsg: "query ID is required",
		},
		{
			name: "invalid query ID format",
			req: &ExecuteRequest{
				QueryID: "invalid",
			},
			errMsg: "query ID must start with '#'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.ExecuteNQLV1(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestExecuteNQLV2_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *ExecuteRequest
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "execute request cannot be nil",
		},
		{
			name: "empty query ID",
			req: &ExecuteRequest{
				QueryID: "",
			},
			errMsg: "query ID is required",
		},
		{
			name: "invalid query ID format",
			req: &ExecuteRequest{
				QueryID: "invalid",
			},
			errMsg: "query ID must start with '#'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.ExecuteNQLV2(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestStartNQLExport_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExportRequest{
		QueryID: "#test_query",
		Format:  ExportFormatCSV,
	}

	result, resp, err := service.StartNQLExport(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ExportID)
}

func TestStartNQLExport_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *ExportRequest
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "export request cannot be nil",
		},
		{
			name: "empty query ID",
			req: &ExportRequest{
				QueryID: "",
			},
			errMsg: "query ID is required",
		},
		{
			name: "invalid format",
			req: &ExportRequest{
				QueryID: "#test_query",
				Format:  "invalid",
			},
			errMsg: "format must be either 'csv' or 'json'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.StartNQLExport(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestGetNQLExportStatus_Submitted(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.GetNQLExportStatus(context.Background(), "export-123-abc")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "SUBMITTED", result.Status)
	assert.Empty(t, result.ResultsFileURL)
}

func TestGetNQLExportStatus_Completed(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.GetNQLExportStatus(context.Background(), "export-456-def")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "COMPLETED", result.Status)
	assert.NotEmpty(t, result.ResultsFileURL)
}

func TestGetNQLExportStatus_Error(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.GetNQLExportStatus(context.Background(), "export-789-ghi")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "ERROR", result.Status)
	assert.NotEmpty(t, result.ErrorDescription)
}

func TestGetNQLExportStatus_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	_, _, err := service.GetNQLExportStatus(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "export ID cannot be empty")
}

func TestDownloadNQLExport_Success(t *testing.T) {
	t.Skip("DownloadNQLExport uses a separate HTTP client for S3 downloads which cannot be mocked in unit tests")
	
	service, _ := setupMockClient(t)

	// Mock S3 download URL
	s3URL := "https://s3.amazonaws.com/nexthink-exports/test.csv"
	csvData := "name,os,memory\ndevice1,Windows,100\ndevice2,macOS,150"

	httpmock.RegisterResponder("GET", s3URL,
		httpmock.NewStringResponder(200, csvData))

	data, err := service.DownloadNQLExport(context.Background(), s3URL)

	require.NoError(t, err)
	assert.Equal(t, []byte(csvData), data)
}

func TestDownloadNQLExport_EmptyURL(t *testing.T) {
	service, _ := setupMockClient(t)

	_, err := service.DownloadNQLExport(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "download URL cannot be empty")
}

func TestWaitForNQLExport_CompletesImmediately(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	result, err := service.WaitForNQLExport(context.Background(), "export-456-def", time.Second, 10*time.Second)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "COMPLETED", result.Status)
	assert.NotEmpty(t, result.ResultsFileURL)
}

func TestWaitForNQLExport_Timeout(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	_, err := service.WaitForNQLExport(context.Background(), "export-123-abc", 100*time.Millisecond, 500*time.Millisecond)

	require.Error(t, err)
	// Context deadline exceeded is expected when timeout occurs during polling
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestWaitForNQLExport_Error(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewNQLMock(baseURL)
	mockHandler.RegisterMocks()

	result, err := service.WaitForNQLExport(context.Background(), "export-789-ghi", time.Second, 10*time.Second)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ERROR", result.Status)
	assert.NotEmpty(t, result.ErrorDescription)
}

func TestWaitForNQLExport_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	_, err := service.WaitForNQLExport(context.Background(), "", time.Second, 10*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "export ID cannot be empty")
}

func TestIsTerminalStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"completed is terminal", ExportStatusCompleted, true},
		{"error is terminal", ExportStatusError, true},
		{"submitted is not terminal", ExportStatusSubmitted, false},
		{"in_progress is not terminal", ExportStatusInProgress, false},
		{"unknown status is not terminal", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTerminalStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}
