package remote_actions

import (
	"context"
	"net/http"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/remote_actions/mocks"
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

func TestTriggerRemoteAction_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewRemoteActionsMock(baseURL)
	mockHandler.RegisterMocks()

	req := &TriggerRemoteActionRequest{
		RemoteActionID:   "#clear_browser_cache",
		Devices:          []string{"device-001", "device-002"},
		ExpiresInMinutes: 1440,
		Params: map[string]string{
			"browser": "chrome",
		},
		TriggerInfo: &TriggerInfoRequest{
			ExternalSource:    "ServiceDesk",
			Reason:            "User reported slow browser performance",
			ExternalReference: "TICKET-12345",
		},
	}

	result, resp, err := service.TriggerRemoteAction(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.RequestID)
	assert.Equal(t, 1440, result.ExpiresInMinutes)
}

func TestTriggerRemoteAction_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *TriggerRemoteActionRequest
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "execution request cannot be nil",
		},
		{
			name: "empty remote action ID",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "",
				Devices:        []string{"device-001"},
			},
			errMsg: "remote action ID cannot be empty",
		},
		{
			name: "no devices",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "#test_action",
				Devices:        []string{},
			},
			errMsg: "at least 1 device is required",
		},
		{
			name: "too many devices",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "#test_action",
				Devices:        make([]string, 10001),
			},
			errMsg: "maximum 10000 devices allowed",
		},
		{
			name: "invalid expires in minutes - too low",
			req: &TriggerRemoteActionRequest{
				RemoteActionID:   "#test_action",
				Devices:          []string{"device-001"},
				ExpiresInMinutes: 30,
			},
			errMsg: "expiresInMinutes must be at least 60",
		},
		{
			name: "invalid expires in minutes - too high",
			req: &TriggerRemoteActionRequest{
				RemoteActionID:   "#test_action",
				Devices:          []string{"device-001"},
				ExpiresInMinutes: 20000,
			},
			errMsg: "expiresInMinutes cannot exceed 10080",
		},
		{
			name: "trigger info reason too long",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "#test_action",
				Devices:        []string{"device-001"},
				TriggerInfo: &TriggerInfoRequest{
					Reason: string(make([]byte, 501)),
				},
			},
			errMsg: "reason cannot exceed 500 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.TriggerRemoteAction(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestListRemoteActions_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewRemoteActionsMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.ListRemoteActions(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
	
	// Verify first remote action
	assert.Equal(t, "#clear_browser_cache", result[0].ID)
	assert.Equal(t, "Clear Browser Cache", result[0].Name)
	assert.Equal(t, "NEXTHINK", result[0].Origin)
	assert.True(t, result[0].Targeting.APIEnabled)
	assert.Equal(t, "LOCAL_SYSTEM", result[0].ScriptInfo.RunAs)
	assert.Len(t, result[0].ScriptInfo.Inputs, 1)
	assert.Len(t, result[0].ScriptInfo.Outputs, 1)
	
	// Verify second remote action
	assert.Equal(t, "#collect_logs", result[1].ID)
	assert.Equal(t, "Collect System Logs", result[1].Name)
	assert.Equal(t, "CUSTOM", result[1].Origin)
	assert.False(t, result[1].Targeting.ManualEnabled)
}

func TestGetRemoteActionDetails_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewRemoteActionsMock(baseURL)
	mockHandler.RegisterMocks()

	result, resp, err := service.GetRemoteActionDetails(context.Background(), "#clear_browser_cache")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, "#clear_browser_cache", result.ID)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", result.UUID)
	assert.Equal(t, "Clear Browser Cache", result.Name)
	assert.Contains(t, result.Purpose, "REMEDIATION")
	assert.True(t, result.ScriptInfo.HasScriptWindows)
	assert.True(t, result.ScriptInfo.HasScriptMacOS)
}

func TestGetRemoteActionDetails_ValidationError(t *testing.T) {
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
		{
			name:   "invalid NQL ID format - too short",
			nqlID:  "#",
			errMsg: "NQL ID must be at least 2 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.GetRemoteActionDetails(context.Background(), tt.nqlID)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
