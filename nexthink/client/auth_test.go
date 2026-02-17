package client

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

func TestAuthConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AuthConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Instance:     "test-instance",
				Region:       RegionUS,
			},
			wantErr: false,
		},
		{
			name: "valid config with custom token URL",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Instance:     "test-instance",
				Region:       RegionEU,
				TokenURL:     "https://custom-token-url.example.com/oauth2/token",
			},
			wantErr: false,
		},
		{
			name: "valid config with custom scope",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Instance:     "test-instance",
				Region:       RegionPAC,
				Scope:        "custom:scope",
			},
			wantErr: false,
		},
		{
			name: "empty client ID",
			config: &AuthConfig{
				ClientID:     "",
				ClientSecret: "test-client-secret",
				Instance:     "test-instance",
				Region:       RegionUS,
			},
			wantErr: true,
			errMsg:  "client ID is required",
		},
		{
			name: "empty client secret",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "",
				Instance:     "test-instance",
				Region:       RegionUS,
			},
			wantErr: true,
			errMsg:  "client secret is required",
		},
		{
			name: "empty instance",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Instance:     "",
				Region:       RegionUS,
			},
			wantErr: true,
			errMsg:  "instance name is required",
		},
		{
			name: "empty region",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Instance:     "test-instance",
				Region:       "",
			},
			wantErr: true,
			errMsg:  "region is required",
		},
		{
			name: "invalid region",
			config: &AuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Instance:     "test-instance",
				Region:       "invalid",
			},
			wantErr: true,
			errMsg:  "invalid region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("AuthConfig.Validate() error message = %q, want to contain %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestAuthConfig_GetTokenURL(t *testing.T) {
	tests := []struct {
		name   string
		config *AuthConfig
		want   string
	}{
		{
			name: "US region default URL",
			config: &AuthConfig{
				Instance: "myinstance",
				Region:   RegionUS,
			},
			want: "https://myinstance-login.us.nexthink.cloud/oauth2/default/v1/token",
		},
		{
			name: "EU region default URL",
			config: &AuthConfig{
				Instance: "myinstance",
				Region:   RegionEU,
			},
			want: "https://myinstance-login.eu.nexthink.cloud/oauth2/default/v1/token",
		},
		{
			name: "PAC region default URL",
			config: &AuthConfig{
				Instance: "myinstance",
				Region:   RegionPAC,
			},
			want: "https://myinstance-login.pac.nexthink.cloud/oauth2/default/v1/token",
		},
		{
			name: "META region default URL",
			config: &AuthConfig{
				Instance: "myinstance",
				Region:   RegionMETA,
			},
			want: "https://myinstance-login.meta.nexthink.cloud/oauth2/default/v1/token",
		},
		{
			name: "custom token URL overrides default",
			config: &AuthConfig{
				Instance: "myinstance",
				Region:   RegionUS,
				TokenURL: "https://custom-url.example.com/token",
			},
			want: "https://custom-url.example.com/token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetTokenURL()
			if got != tt.want {
				t.Errorf("GetTokenURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAuthConfig_GetScope(t *testing.T) {
	tests := []struct {
		name   string
		config *AuthConfig
		want   string
	}{
		{
			name: "default scope",
			config: &AuthConfig{
				Scope: "",
			},
			want: ScopeServiceIntegration,
		},
		{
			name: "custom scope",
			config: &AuthConfig{
				Scope: "custom:scope",
			},
			want: "custom:scope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetScope()
			if got != tt.want {
				t.Errorf("GetScope() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAuthConfig_GenerateBasicAuth(t *testing.T) {
	tests := []struct {
		name         string
		config       *AuthConfig
		wantDecoded  string
	}{
		{
			name: "basic auth generation",
			config: &AuthConfig{
				ClientID:     "my-client-id",
				ClientSecret: "my-client-secret",
			},
			wantDecoded: "my-client-id:my-client-secret",
		},
		{
			name: "basic auth with special characters",
			config: &AuthConfig{
				ClientID:     "client@example.com",
				ClientSecret: "p@ssw0rd!#$",
			},
			wantDecoded: "client@example.com:p@ssw0rd!#$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GenerateBasicAuth()
			
			// Verify it's valid base64
			decoded, err := base64.StdEncoding.DecodeString(got)
			if err != nil {
				t.Fatalf("GenerateBasicAuth() produced invalid base64: %v", err)
			}
			
			if string(decoded) != tt.wantDecoded {
				t.Errorf("GenerateBasicAuth() decoded = %q, want %q", string(decoded), tt.wantDecoded)
			}
		})
	}
}

func TestNewTokenManager(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Instance:     "test-instance",
		Region:       RegionUS,
	}

	tm := NewTokenManager(authConfig, client, logger)

	if tm == nil {
		t.Fatal("NewTokenManager() returned nil")
	}
	if tm.authConfig != authConfig {
		t.Error("TokenManager authConfig not set correctly")
	}
	if tm.logger != logger {
		t.Error("TokenManager logger not set correctly")
	}
	if tm.client != client {
		t.Error("TokenManager client not set correctly")
	}
	if tm.refreshBuffer != TokenRefreshBuffer*time.Second {
		t.Errorf("TokenManager refreshBuffer = %v, want %v", tm.refreshBuffer, TokenRefreshBuffer*time.Second)
	}
}

func TestTokenManager_RefreshToken(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	authConfig := &AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Instance:     "test-instance",
		Region:       RegionUS,
	}

	tokenURL := authConfig.GetTokenURL()

	// Mock successful token response
	httpmock.RegisterResponder("POST", tokenURL,
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "test-access-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	// Configure client to use mocked transport
	client.SetTransport(httpClient.Transport)

	tm := NewTokenManager(authConfig, client, logger)

	token, err := tm.RefreshToken()
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if token != "test-access-token" {
		t.Errorf("RefreshToken() token = %q, want %q", token, "test-access-token")
	}

	// Verify token is stored
	tm.mu.RLock()
	if tm.currentToken == nil {
		t.Error("currentToken not set after refresh")
	}
	if tm.currentToken.AccessToken != "test-access-token" {
		t.Errorf("currentToken.AccessToken = %q, want %q", tm.currentToken.AccessToken, "test-access-token")
	}
	tm.mu.RUnlock()
}

func TestTokenManager_GetToken_CachedToken(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Instance:     "test-instance",
		Region:       RegionUS,
	}

	tm := NewTokenManager(authConfig, client, logger)

	// Manually set a valid cached token
	tm.mu.Lock()
	tm.currentToken = &TokenResponse{
		AccessToken: "cached-token",
		ExpiresIn:   900,
		TokenType:   "Bearer",
		Scope:       "service:integration",
	}
	tm.tokenExpiry = time.Now().Add(10 * time.Minute) // Valid for 10 more minutes
	tm.mu.Unlock()

	token, err := tm.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	if token != "cached-token" {
		t.Errorf("GetToken() = %q, want %q (cached)", token, "cached-token")
	}
}

func TestTokenManager_GetToken_ExpiredToken(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	authConfig := &AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Instance:     "test-instance",
		Region:       RegionUS,
	}

	tokenURL := authConfig.GetTokenURL()

	// Mock successful token response
	httpmock.RegisterResponder("POST", tokenURL,
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "refreshed-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	// Configure client to use mocked transport
	client.SetTransport(httpClient.Transport)

	tm := NewTokenManager(authConfig, client, logger)

	// Set an expired token
	tm.mu.Lock()
	tm.currentToken = &TokenResponse{
		AccessToken: "expired-token",
		ExpiresIn:   900,
		TokenType:   "Bearer",
		Scope:       "service:integration",
	}
	tm.tokenExpiry = time.Now().Add(-1 * time.Minute) // Expired 1 minute ago
	tm.mu.Unlock()

	token, err := tm.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	if token != "refreshed-token" {
		t.Errorf("GetToken() = %q, want %q (refreshed)", token, "refreshed-token")
	}
}

func TestSetupAuthentication_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	// Create a custom HTTP client and activate httpmock on it
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	authConfig := &AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Instance:     "test-instance",
		Region:       RegionUS,
	}

	tokenURL := authConfig.GetTokenURL()

	// Mock successful token response
	httpmock.RegisterResponder("POST", tokenURL,
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "initial-token",
			"expires_in":   900,
			"token_type":   "Bearer",
			"scope":        "service:integration",
		}))

	// Configure client to use mocked transport
	client.SetTransport(httpClient.Transport)

	tokenManager, err := SetupAuthentication(client, authConfig, logger)
	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v", err)
	}

	if tokenManager == nil {
		t.Fatal("SetupAuthentication() returned nil TokenManager")
	}
}

func TestSetupAuthentication_InvalidConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	tests := []struct {
		name       string
		authConfig *AuthConfig
		wantErr    bool
	}{
		{
			name: "empty client ID",
			authConfig: &AuthConfig{
				ClientID:     "",
				ClientSecret: "test-secret",
				Instance:     "test-instance",
				Region:       RegionUS,
			},
			wantErr: true,
		},
		{
			name: "empty client secret",
			authConfig: &AuthConfig{
				ClientID:     "test-id",
				ClientSecret: "",
				Instance:     "test-instance",
				Region:       RegionUS,
			},
			wantErr: true,
		},
		{
			name: "invalid region",
			authConfig: &AuthConfig{
				ClientID:     "test-id",
				ClientSecret: "test-secret",
				Instance:     "test-instance",
				Region:       "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SetupAuthentication(client, tt.authConfig, logger)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetupAuthentication() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthConfig_ValidRegions(t *testing.T) {
	validRegions := []string{RegionUS, RegionEU, RegionPAC, RegionMETA}

	for _, region := range validRegions {
		t.Run(region, func(t *testing.T) {
			config := &AuthConfig{
				ClientID:     "test-id",
				ClientSecret: "test-secret",
				Instance:     "test-instance",
				Region:       region,
			}

			err := config.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v for valid region %q", err, region)
			}
		})
	}
}

func TestAuthConfig_LongCredentials(t *testing.T) {
	// Test with very long client ID and secret (should still be valid)
	longID := strings.Repeat("a", 1000)
	longSecret := strings.Repeat("b", 1000)

	config := &AuthConfig{
		ClientID:     longID,
		ClientSecret: longSecret,
		Instance:     "test-instance",
		Region:       RegionUS,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() with long credentials error = %v, want nil", err)
	}

	// Verify BasicAuth can handle long credentials
	basicAuth := config.GenerateBasicAuth()
	decoded, err := base64.StdEncoding.DecodeString(basicAuth)
	if err != nil {
		t.Fatalf("GenerateBasicAuth() produced invalid base64 for long credentials: %v", err)
	}

	expected := longID + ":" + longSecret
	if string(decoded) != expected {
		t.Error("GenerateBasicAuth() failed to encode long credentials correctly")
	}
}

func TestAuthConfig_SpecialCharacters(t *testing.T) {
	// Test with special characters in credentials
	specialChars := []struct {
		name   string
		id     string
		secret string
	}{
		{"dashes", "client-id-with-dashes", "secret-with-dashes"},
		{"underscores", "client_id_with_underscores", "secret_with_underscores"},
		{"dots", "client.id.with.dots", "secret.with.dots"},
		{"mixed", "client-id_123.test", "secret_456-test.com"},
		{"symbols", "client@example.com", "p@ssw0rd!#$%"},
	}

	for _, tt := range specialChars {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				ClientID:     tt.id,
				ClientSecret: tt.secret,
				Instance:     "test-instance",
				Region:       RegionUS,
			}

			err := config.Validate()
			if err != nil {
				t.Errorf("Validate() with special characters error = %v, want nil", err)
			}

			basicAuth := config.GenerateBasicAuth()
			decoded, err := base64.StdEncoding.DecodeString(basicAuth)
			if err != nil {
				t.Fatalf("GenerateBasicAuth() produced invalid base64: %v", err)
			}

			expected := tt.id + ":" + tt.secret
			if string(decoded) != expected {
				t.Errorf("GenerateBasicAuth() decoded = %q, want %q", string(decoded), expected)
			}
		})
	}
}
