package client

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap/zaptest"
)

func TestNewTransport_ValidationErrors(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		instance     string
		region       string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "empty client ID",
			clientID:     "",
			clientSecret: "test-secret",
			instance:     "test-instance",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "client ID cannot be empty",
		},
		{
			name:         "empty client secret",
			clientID:     "test-id",
			clientSecret: "",
			instance:     "test-instance",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "client secret cannot be empty",
		},
		{
			name:         "empty instance",
			clientID:     "test-id",
			clientSecret: "test-secret",
			instance:     "",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "instance cannot be empty",
		},
		{
			name:         "empty region",
			clientID:     "test-id",
			clientSecret: "test-secret",
			instance:     "test-instance",
			region:       "",
			wantErr:      true,
			errContains:  "region cannot be empty",
		},
		{
			name:         "invalid region",
			clientID:     "test-id",
			clientSecret: "test-secret",
			instance:     "test-instance",
			region:       "invalid-region",
			wantErr:      true,
			errContains:  "invalid region",
		},
		{
			name:         "instance with spaces",
			clientID:     "test-id",
			clientSecret: "test-secret",
			instance:     "test instance",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "instance name cannot contain spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransport(tt.clientID, tt.clientSecret, tt.instance, tt.region)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("NewTransport() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestNewTransport_OptionValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		option      ClientOption
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid timeout",
			option:      WithTimeout(-1 * time.Second),
			wantErr:     true,
			errContains: "timeout must be greater than 0",
		},
		{
			name:        "invalid retry count",
			option:      WithRetryCount(-1),
			wantErr:     true,
			errContains: "retry count cannot be negative",
		},
		{
			name:        "invalid base URL",
			option:      WithBaseURL("not-a-url"),
			wantErr:     true,
			errContains: "invalid base URL",
		},
		{
			name:        "invalid proxy URL",
			option:      WithProxy("not-a-proxy"),
			wantErr:     true,
			errContains: "invalid proxy URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				tt.option,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("NewTransport() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestNewTransport_Success(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock successful token response
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if transport == nil {
		t.Fatal("NewTransport() returned nil transport")
	}

	// Verify basic fields are set
	if transport.BaseURL == "" {
		t.Error("BaseURL should be set")
	}

	if transport.authConfig == nil {
		t.Error("authConfig should be set")
	}

	if transport.tokenManager == nil {
		t.Error("tokenManager should be set")
	}

	if transport.client == nil {
		t.Error("client should be set")
	}

	if transport.logger == nil {
		t.Error("logger should be set")
	}
}

func TestNewTransport_DefaultBaseURL(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	tests := []struct {
		name         string
		instance     string
		region       string
		expectedURL  string
	}{
		{
			name:        "US region",
			instance:    "myinstance",
			region:      RegionUS,
			expectedURL: "https://myinstance.api.us.nexthink.cloud",
		},
		{
			name:        "EU region",
			instance:    "myinstance",
			region:      RegionEU,
			expectedURL: "https://myinstance.api.eu.nexthink.cloud",
		},
		{
			name:        "PAC region",
			instance:    "myinstance",
			region:      RegionPAC,
			expectedURL: "https://myinstance.api.pac.nexthink.cloud",
		},
		{
			name:        "META region",
			instance:    "myinstance",
			region:      RegionMETA,
			expectedURL: "https://myinstance.api.meta.nexthink.cloud",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock token endpoint for each region
			tokenURL := "https://" + tt.instance + "-login." + tt.region + ".nexthink.cloud/oauth2/default/v1/token"
			httpmock.RegisterResponder("POST", tokenURL,
				httpmock.NewJsonResponderOrPanic(200, map[string]any{
					"access_token": "test-token",
					"expires_in":   900,
					"token_type":   "Bearer",
					"scope":        "service:integration",
				}))

			logger := zaptest.NewLogger(t)

			transport, err := NewTransport(
				"test-id",
				"test-secret",
				tt.instance,
				tt.region,
				WithLogger(logger),
				WithTransport(httpClient.Transport),
			)

			if err != nil {
				t.Fatalf("NewTransport() error = %v, want nil", err)
			}

			if transport.BaseURL != tt.expectedURL {
				t.Errorf("BaseURL = %q, want %q", transport.BaseURL, tt.expectedURL)
			}

			httpmock.Reset()
		})
	}
}

func TestNewTransport_CustomBaseURL(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)
	customURL := "https://custom.api.nexthink.com"

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithBaseURL(customURL),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if transport.BaseURL != customURL {
		t.Errorf("BaseURL = %q, want %q", transport.BaseURL, customURL)
	}
}

func TestNewTransport_DefaultUserAgent(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	// Verify user agent contains expected components
	if !strings.Contains(transport.userAgent, UserAgentBase) {
		t.Errorf("userAgent = %q, want to contain %q", transport.userAgent, UserAgentBase)
	}

	if !strings.Contains(transport.userAgent, "gzip") {
		t.Errorf("userAgent = %q, want to contain 'gzip'", transport.userAgent)
	}
}

func TestTransport_GetHTTPClient(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	client := transport.GetHTTPClient()
	if client == nil {
		t.Error("GetHTTPClient() returned nil")
	}

	if client != transport.client {
		t.Error("GetHTTPClient() did not return the internal client")
	}
}

func TestTransport_GetLogger(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	retrievedLogger := transport.GetLogger()
	if retrievedLogger == nil {
		t.Error("GetLogger() returned nil")
	}

	if retrievedLogger != logger {
		t.Error("GetLogger() did not return the configured logger")
	}
}

func TestTransport_GetTokenManager(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	tokenManager := transport.GetTokenManager()
	if tokenManager == nil {
		t.Error("GetTokenManager() returned nil")
	}

	if tokenManager != transport.tokenManager {
		t.Error("GetTokenManager() did not return the internal token manager")
	}
}

func TestTransport_RefreshToken(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint for both initial and refresh calls
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "refreshed-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	// Call RefreshToken
	err = transport.RefreshToken()
	if err != nil {
		t.Errorf("RefreshToken() error = %v, want nil", err)
	}

	// Verify the token was refreshed by checking call count
	info := httpmock.GetCallCountInfo()
	callCount := info["POST https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token"]
	if callCount < 2 {
		t.Errorf("Token endpoint was called %d times, want at least 2 (initial + refresh)", callCount)
	}
}

func TestTransport_InvalidateToken(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	// Call InvalidateToken - should not panic
	transport.InvalidateToken()

	// Verify token is invalidated by checking that GetToken will refresh
	token, err := transport.tokenManager.GetToken()
	if err != nil {
		t.Errorf("GetToken() after invalidation error = %v", err)
	}

	if token == "" {
		t.Error("GetToken() returned empty token after invalidation")
	}
}

func TestTransport_QueryBuilder(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	qb := transport.QueryBuilder()
	if qb == nil {
		t.Error("QueryBuilder() returned nil")
	}

	// Verify it's a functional query builder
	qb.AddString("test", "value")
	if !qb.Has("test") {
		t.Error("QueryBuilder() did not return a functional query builder")
	}
}

func TestNewTransport_GlobalHeaders(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithGlobalHeader("X-Custom-Header", "custom-value"),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	if transport.globalHeaders["X-Custom-Header"] != "custom-value" {
		t.Errorf("globalHeaders[X-Custom-Header] = %q, want %q",
			transport.globalHeaders["X-Custom-Header"], "custom-value")
	}
}

func TestNewTransport_AuthConfig(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	logger := zaptest.NewLogger(t)

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	// Verify authConfig was populated correctly
	if transport.authConfig.ClientID != "test-id" {
		t.Errorf("authConfig.ClientID = %q, want %q", transport.authConfig.ClientID, "test-id")
	}

	if transport.authConfig.ClientSecret != "test-secret" {
		t.Errorf("authConfig.ClientSecret = %q, want %q", transport.authConfig.ClientSecret, "test-secret")
	}

	if transport.authConfig.Instance != "test-instance" {
		t.Errorf("authConfig.Instance = %q, want %q", transport.authConfig.Instance, "test-instance")
	}

	if transport.authConfig.Region != RegionUS {
		t.Errorf("authConfig.Region = %q, want %q", transport.authConfig.Region, RegionUS)
	}

	if transport.authConfig.Scope != ScopeServiceIntegration {
		t.Errorf("authConfig.Scope = %q, want %q", transport.authConfig.Scope, ScopeServiceIntegration)
	}
}

func TestNewTransport_CustomScope(t *testing.T) {
	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mock token endpoint
	httpmock.RegisterResponder("POST", "https://test-instance-login.us.nexthink.cloud/oauth2/default/v1/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "custom:scope",
		}))

	logger := zaptest.NewLogger(t)
	customScope := "custom:scope"

	transport, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithScope(customScope),
		WithTransport(httpClient.Transport),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v", err)
	}

	if transport.authConfig.Scope != customScope {
		t.Errorf("authConfig.Scope = %q, want %q", transport.authConfig.Scope, customScope)
	}
}

func TestTransport_Struct(t *testing.T) {
	// Test that Transport struct has all expected fields
	transport := &Transport{
		client:        nil,
		logger:        nil,
		authConfig:    nil,
		tokenManager:  nil,
		BaseURL:       "https://test.example.com",
		globalHeaders: make(map[string]string),
		userAgent:     "test-agent",
	}

	if transport.BaseURL != "https://test.example.com" {
		t.Errorf("BaseURL = %q, want %q", transport.BaseURL, "https://test.example.com")
	}

	if transport.userAgent != "test-agent" {
		t.Errorf("userAgent = %q, want %q", transport.userAgent, "test-agent")
	}

	if transport.globalHeaders == nil {
		t.Error("globalHeaders should not be nil")
	}
}
