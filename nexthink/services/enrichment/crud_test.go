package enrichment

import (
	"context"
	"net/http"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/enrichment/mocks"
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

func TestEnrichFields_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewEnrichmentMock(baseURL)
	mockHandler.RegisterMocks()

	req := &EnrichmentRequest{
		Domain: "configuration",
		Enrichments: []Enrichment{
			{
				Identification: []Identification{
					{
						Name:  IdentificationDeviceName,
						Value: "DESKTOP-001",
					},
				},
				Fields: []Field{
					{
						Name:  FieldDeviceConfigurationTag,
						Value: "Production",
					},
					{
						Name:  FieldDeviceVirtualizationType,
						Value: "physical",
					},
				},
			},
			{
				Identification: []Identification{
					{
						Name:  IdentificationDeviceUID,
						Value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
					},
				},
				Fields: []Field{
					{
						Name:  FieldDeviceConfigurationTag,
						Value: "Development",
					},
				},
			},
		},
	}

	result, resp, err := service.EnrichFields(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestEnrichFields_ValidationError(t *testing.T) {
	service, _ := setupMockClient(t)

	tests := []struct {
		name   string
		req    *EnrichmentRequest
		errMsg string
	}{
		{
			name:   "nil request",
			req:    nil,
			errMsg: "enrichment request cannot be nil",
		},
		{
			name: "empty domain",
			req: &EnrichmentRequest{
				Domain: "",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  IdentificationDeviceName,
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{
							{
								Name:  FieldDeviceConfigurationTag,
								Value: "Production",
							},
						},
					},
				},
			},
			errMsg: "domain is required",
		},
		{
			name: "empty enrichments",
			req: &EnrichmentRequest{
				Domain:      "configuration",
				Enrichments: []Enrichment{},
			},
			errMsg: "enrichments is required",
		},
		{
			name: "too many enrichments",
			req: &EnrichmentRequest{
				Domain:      "configuration",
				Enrichments: make([]Enrichment, 5001),
			},
			errMsg: "enrichments cannot contain more than 5000 items",
		},
		{
			name: "invalid identification count",
			req: &EnrichmentRequest{
				Domain: "configuration",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{},
						Fields: []Field{
							{
								Name:  FieldDeviceConfigurationTag,
								Value: "Production",
							},
						},
					},
				},
			},
			errMsg: "identification must contain exactly 1 item",
		},
		{
			name: "empty identification name",
			req: &EnrichmentRequest{
				Domain: "configuration",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  "",
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{
							{
								Name:  FieldDeviceConfigurationTag,
								Value: "Production",
							},
						},
					},
				},
			},
			errMsg: "identification[0].name is required",
		},
		{
			name: "invalid identification name",
			req: &EnrichmentRequest{
				Domain: "configuration",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  "invalid/field/name",
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{
							{
								Name:  FieldDeviceConfigurationTag,
								Value: "Production",
							},
						},
					},
				},
			},
			errMsg: "identification[0].name has invalid value",
		},
		{
			name: "empty fields",
			req: &EnrichmentRequest{
				Domain: "configuration",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  IdentificationDeviceName,
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{},
					},
				},
			},
			errMsg: "fields is required",
		},
		{
			name: "empty field name",
			req: &EnrichmentRequest{
				Domain: "configuration",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  IdentificationDeviceName,
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{
							{
								Name:  "",
								Value: "Production",
							},
						},
					},
				},
			},
			errMsg: "fields[0].name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.EnrichFields(context.Background(), tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestEnrichFields_AllFieldTypes(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := mocks.NewEnrichmentMock(baseURL)
	mockHandler.RegisterMocks()

	req := &EnrichmentRequest{
		Domain: "configuration",
		Enrichments: []Enrichment{
			{
				Identification: []Identification{
					{
						Name:  IdentificationDeviceName,
						Value: "DESKTOP-001",
					},
				},
				Fields: []Field{
					{
						Name:  FieldDeviceConfigurationTag,
						Value: "Production",
					},
					{
						Name:  FieldDeviceVirtualizationLastUpdate,
						Value: 1609459200,
					},
					{
						Name:  FieldDeviceVirtualizationType,
						Value: "vmware",
					},
				},
			},
		},
	}

	result, resp, err := service.EnrichFields(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}
