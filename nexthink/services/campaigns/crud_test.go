package campaigns

import (
	"context"
	"net/http"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/campaigns/mocks"
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

func TestTriggerCampaign_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewCampaignsMock(baseURL)
	mockHandler.RegisterTriggerCampaignSuccessMock()

	req := &TriggerRequest{
		CampaignNqlId: "#security_awareness",
		UserSid: []string{
			"S-1-5-21-1234567890-1234567890-1234567890-1001",
			"S-1-5-21-1234567890-1234567890-1234567890-1002",
		},
		ExpiresInMinutes: 1440,
		Parameters: map[string]string{
			"department": "IT",
			"location":   "New York",
		},
	}

	result, resp, err := service.TriggerCampaign(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Len(t, result.Requests, 2)
	
	// Verify first request
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", result.Requests[0].RequestId)
	assert.Equal(t, "S-1-5-21-1234567890-1234567890-1234567890-1001", result.Requests[0].UserSid)
	assert.Empty(t, result.Requests[0].Message)
	
	// Verify second request
	assert.Equal(t, "b2c3d4e5-f6a7-8901-bcde-f12345678901", result.Requests[1].RequestId)
	assert.Equal(t, "S-1-5-21-1234567890-1234567890-1234567890-1002", result.Requests[1].UserSid)
}

func TestTriggerCampaign_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *TriggerRequest
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "trigger request cannot be nil",
		},
		{
			name: "empty campaign NQL ID",
			req: &TriggerRequest{
				CampaignNqlId: "",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1440,
			},
			errMsg: "campaignNqlId is required",
		},
		{
			name: "empty userSid list",
			req: &TriggerRequest{
				CampaignNqlId:    "#security_awareness",
				UserSid:          []string{},
				ExpiresInMinutes: 1440,
			},
			errMsg: "userSid is required",
		},
		{
			name: "too many userSids",
			req: &TriggerRequest{
				CampaignNqlId:    "#security_awareness",
				UserSid:          make([]string, 10001),
				ExpiresInMinutes: 1440,
			},
			errMsg: "userSid cannot contain more than 10000 SIDs",
		},
		{
			name: "empty userSid in list",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
					"",
				},
				ExpiresInMinutes: 1440,
			},
			errMsg: "userSid[1] cannot be empty",
		},
		{
			name: "expiresInMinutes too low",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 0,
			},
			errMsg: "expiresInMinutes must be at least 1",
		},
		{
			name: "expiresInMinutes too high",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 525601,
			},
			errMsg: "expiresInMinutes cannot exceed 525600",
		},
		{
			name: "too many parameters",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1440,
				Parameters:       make(map[string]string, 31),
			},
			errMsg: "parameters cannot contain more than 30 items",
		},
	}

	// Fill the too many parameters test case
	for _, tt := range tests {
		if tt.name == "too many parameters" {
			for i := 0; i < 31; i++ {
				tt.req.Parameters[string(rune('a'+i))] = "value"
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.TriggerCampaign(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestTriggerCampaign_MinimalRequest(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewCampaignsMock(baseURL)
	mockHandler.RegisterMocks()

	req := &TriggerRequest{
		CampaignNqlId: "#test_campaign",
		UserSid: []string{
			"S-1-5-21-1234567890-1234567890-1234567890-1001",
		},
		ExpiresInMinutes: 60,
	}

	result, resp, err := service.TriggerCampaign(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Len(t, result.Requests, 2)
}

func TestTriggerCampaign_WithParameters(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewCampaignsMock(baseURL)
	mockHandler.RegisterMocks()

	req := &TriggerRequest{
		CampaignNqlId: "#test_campaign",
		UserSid: []string{
			"S-1-5-21-1234567890-1234567890-1234567890-1001",
		},
		ExpiresInMinutes: 1440,
		Parameters: map[string]string{
			"param1": "value1",
			"param2": "value2",
			"param3": "value3",
		},
	}

	result, resp, err := service.TriggerCampaign(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Len(t, result.Requests, 2)
}
