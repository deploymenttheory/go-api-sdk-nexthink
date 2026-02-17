package workflows

import (
	"context"
	"net/http"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/workflows/mocks"
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

func TestExecuteV1_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewWorkflowsMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExecutionRequestV1{
		WorkflowID: "#password_reset",
		Devices:    []string{"device-001", "device-002"},
		Users:      []string{"S-1-5-21-1234567890-1234567890-1234567890-1001"},
		Params: map[string]string{
			"reason": "User forgot password",
		},
	}

	result, resp, err := service.ExecuteV1(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", result.RequestUUID)
	assert.Len(t, result.ExecutionsUUIDs, 2)
}

func TestExecuteV2_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewWorkflowsMock(baseURL)
	mockHandler.RegisterMocks()

	req := &ExecutionRequestV2{
		WorkflowID: "#device_diagnostics",
		Devices: []DeviceData{
			{
				Name: "DESKTOP-001",
				UID:  "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			},
		},
		Users: []UserData{
			{
				UPN: "user@example.com",
				SID: "S-1-5-21-1234567890-1234567890-1234567890-1001",
			},
		},
	}

	result, resp, err := service.ExecuteV2(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "d4e5f6g7-h8i9-0123-defg-h34567890123", result.RequestUUID)
	assert.Len(t, result.ExecutionsUUIDs, 1)
}

func TestExecuteV1_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *ExecutionRequestV1
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "execution request cannot be nil",
		},
		{
			name: "empty workflow ID",
			req: &ExecutionRequestV1{
				WorkflowID: "",
				Devices:    []string{"device-001"},
			},
			errMsg: "workflow ID cannot be empty",
		},
		{
			name: "no devices or users",
			req: &ExecutionRequestV1{
				WorkflowID: "#test_workflow",
			},
			errMsg: "at least one device or user must be provided",
		},
		{
			name: "too many devices",
			req: &ExecutionRequestV1{
				WorkflowID: "#test_workflow",
				Devices:    make([]string, 10001),
			},
			errMsg: "maximum 10000 devices allowed",
		},
		{
			name: "empty device ID",
			req: &ExecutionRequestV1{
				WorkflowID: "#test_workflow",
				Devices:    []string{"device-001", ""},
			},
			errMsg: "device at index 1 cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.ExecuteV1(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestExecuteV2_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *ExecutionRequestV2
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "execution request cannot be nil",
		},
		{
			name: "empty workflow ID",
			req: &ExecutionRequestV2{
				WorkflowID: "",
				Devices: []DeviceData{
					{Name: "DESKTOP-001"},
				},
			},
			errMsg: "workflow ID cannot be empty",
		},
		{
			name: "device with no identifiers",
			req: &ExecutionRequestV2{
				WorkflowID: "#test_workflow",
				Devices: []DeviceData{
					{},
				},
			},
			errMsg: "at least one device identifier",
		},
		{
			name: "invalid device UID",
			req: &ExecutionRequestV2{
				WorkflowID: "#test_workflow",
				Devices: []DeviceData{
					{UID: "invalid-uuid"},
				},
			},
			errMsg: "invalid device UID UUID format",
		},
		{
			name: "invalid user SID",
			req: &ExecutionRequestV2{
				WorkflowID: "#test_workflow",
				Users: []UserData{
					{SID: "invalid-sid"},
				},
			},
			errMsg: "invalid SID format",
		},
		{
			name: "invalid user UPN",
			req: &ExecutionRequestV2{
				WorkflowID: "#test_workflow",
				Users: []UserData{
					{UPN: "not-an-email"},
				},
			},
			errMsg: "invalid UPN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.ExecuteV2(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestListWorkflows_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewWorkflowsMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.ListWorkflows(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
	
	// Verify first workflow
	assert.Equal(t, "#password_reset", result[0].ID)
	assert.Equal(t, "Password Reset Workflow", result[0].Name)
	assert.Equal(t, "ACTIVE", result[0].Status)
	assert.Contains(t, result[0].TriggerMethods, "API")
	assert.Len(t, result[0].Versions, 2)
	assert.True(t, result[0].Versions[0].IsActive)
	
	// Verify inactive workflow
	assert.Equal(t, "#inactive_workflow", result[2].ID)
	assert.Equal(t, "INACTIVE", result[2].Status)
}

func TestGetWorkflowDetails_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewWorkflowsMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.GetWorkflowDetails(context.Background(), "#password_reset")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "#password_reset", result.ID)
	assert.Equal(t, "f6g7h8i9-j0k1-2345-fghi-j56789012345", result.UUID)
	assert.Equal(t, "Password Reset Workflow", result.Name)
	assert.Equal(t, "ACTIVE", result.Status)
	assert.Len(t, result.TriggerMethods, 3)
}

func TestGetWorkflowDetails_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		nqlID  string
		errMsg string
	}{
		{
			name:   "empty NQL ID",
			nqlID:  "",
			errMsg: "NQL ID cannot be empty",
		},
		{
			name:   "invalid NQL ID format - no hash",
			nqlID:  "invalid",
			errMsg: "NQL ID must start with #",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.GetWorkflowDetails(context.Background(), tt.nqlID)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestTriggerThinklet_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewWorkflowsMock(baseURL)
	mockHandler.RegisterMocks()

	workflowUUID := "f6a7b8c9-d0e1-2345-abcd-a56789012345"
	executionUUID := "b2c3d4e5-f6a7-8901-bcde-f12345678901"

	req := &ThinkletTriggerRequest{
		Parameters: map[string]string{
			"continue": "yes",
		},
	}

	result, resp, err := service.TriggerThinklet(context.Background(), workflowUUID, executionUUID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "a9b0c1d2-e3f4-5678-abcd-e89012345678", result.RequestUUID)
}

func TestTriggerThinklet_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	validUUID := "f6a7b8c9-d0e1-2345-abcd-a56789012345"

	tests := []struct {
		name          string
		workflowUUID  string
		executionUUID string
		errMsg        string
	}{
		{
			name:          "invalid workflow UUID",
			workflowUUID:  "invalid",
			executionUUID: validUUID,
			errMsg:        "invalid workflow UUID format",
		},
		{
			name:          "invalid execution UUID",
			workflowUUID:  validUUID,
			executionUUID: "invalid",
			errMsg:        "invalid execution UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.TriggerThinklet(context.Background(), tt.workflowUUID, tt.executionUUID, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
