package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// AuthConfig holds authentication configuration for the Nexthink API
type AuthConfig struct {
	// ClientID is the OAuth2 client ID
	ClientID string

	// ClientSecret is the OAuth2 client secret
	ClientSecret string

	// Instance is the Nexthink instance name
	Instance string

	// Region is the Nexthink region (us, eu, pac, meta)
	Region string

	// TokenURL is the optional custom token endpoint URL
	// If not provided, it will be constructed from Instance and Region
	TokenURL string

	// Scope is the OAuth2 scope (defaults to service:integration)
	Scope string
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	TokenType   string `json:"token_type"`   // "Bearer"
	ExpiresIn   int    `json:"expires_in"`   // Token lifetime in seconds (900 = 15 minutes)
	AccessToken string `json:"access_token"` // The access token
	Scope       string `json:"scope"`        // "service:integration"
}

// TokenManager handles OAuth2 token lifecycle
type TokenManager struct {
	authConfig    *AuthConfig
	logger        *zap.Logger
	client        *resty.Client
	currentToken  *TokenResponse
	tokenExpiry   time.Time
	mu            sync.RWMutex
	refreshBuffer time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(authConfig *AuthConfig, client *resty.Client, logger *zap.Logger) *TokenManager {
	return &TokenManager{
		authConfig:    authConfig,
		logger:        logger,
		client:        client,
		refreshBuffer: TokenRefreshBuffer * time.Second,
	}
}

// Validate checks if the auth configuration is valid
func (a *AuthConfig) Validate() error {
	if a.ClientID == "" {
		return fmt.Errorf("client ID is required")
	}
	if a.ClientSecret == "" {
		return fmt.Errorf("client secret is required")
	}
	if a.Instance == "" {
		return fmt.Errorf("instance name is required")
	}
	if a.Region == "" {
		return fmt.Errorf("region is required")
	}

	// Validate region
	validRegions := map[string]bool{
		RegionUS:   true,
		RegionEU:   true,
		RegionPAC:  true,
		RegionMETA: true,
	}
	if !validRegions[a.Region] {
		return fmt.Errorf("invalid region: %s (must be one of: us, eu, pac, meta)", a.Region)
	}

	return nil
}

// GetTokenURL returns the token endpoint URL
func (a *AuthConfig) GetTokenURL() string {
	if a.TokenURL != "" {
		return a.TokenURL
	}
	return fmt.Sprintf(DefaultTokenURLTemplate, a.Instance, a.Region)
}

// GetScope returns the OAuth2 scope
func (a *AuthConfig) GetScope() string {
	if a.Scope != "" {
		return a.Scope
	}
	return ScopeServiceIntegration
}

// GenerateBasicAuth generates the Base64 encoded Basic auth string from clientId:clientSecret
func (a *AuthConfig) GenerateBasicAuth() string {
	credentials := fmt.Sprintf("%s:%s", a.ClientID, a.ClientSecret)
	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

// GetToken returns a valid access token, refreshing if necessary
func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.RLock()
	// Check if we have a valid token that won't expire soon
	if tm.currentToken != nil && time.Now().Add(tm.refreshBuffer).Before(tm.tokenExpiry) {
		token := tm.currentToken.AccessToken
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	// Need to refresh token
	return tm.RefreshToken()
}

// RefreshToken requests a new access token from the OAuth2 endpoint
func (tm *TokenManager) RefreshToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Double check in case another goroutine just refreshed
	if tm.currentToken != nil && time.Now().Add(tm.refreshBuffer).Before(tm.tokenExpiry) {
		return tm.currentToken.AccessToken, nil
	}

	tm.logger.Info("Requesting new OAuth2 access token",
		zap.String("instance", tm.authConfig.Instance),
		zap.String("region", tm.authConfig.Region))

	tokenURL := tm.authConfig.GetTokenURL()
	basicAuth := tm.authConfig.GenerateBasicAuth()
	scope := tm.authConfig.GetScope()

	// Create token request
	resp, err := tm.client.R().
		SetHeader("Content-Type", ContentTypeFormURLEncoded).
		SetHeader("Authorization", fmt.Sprintf("Basic %s", basicAuth)).
		SetFormData(map[string]string{
			"grant_type": GrantTypeClientCredentials,
			"scope":      scope,
		}).
		Post(tokenURL)

	if err != nil {
		tm.logger.Error("Failed to request access token",
			zap.Error(err),
			zap.String("token_url", tokenURL))
		return "", fmt.Errorf("failed to request access token: %w", err)
	}

	if resp.IsError() {
		tm.logger.Error("Token request failed",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
			zap.String("body", resp.String()))
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// Parse token response
	var tokenResp TokenResponse
	if err := json.Unmarshal([]byte(resp.String()), &tokenResp); err != nil {
		tm.logger.Error("Failed to parse token response",
			zap.Error(err),
			zap.String("body", resp.String()))
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	// Store token and calculate expiry
	tm.currentToken = &tokenResp
	tm.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	tm.logger.Info("Successfully obtained access token",
		zap.String("token_type", tokenResp.TokenType),
		zap.Int("expires_in", tokenResp.ExpiresIn),
		zap.Time("expires_at", tm.tokenExpiry))

	return tokenResp.AccessToken, nil
}

// InvalidateToken clears the current token, forcing a refresh on next use
func (tm *TokenManager) InvalidateToken() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.currentToken = nil
	tm.tokenExpiry = time.Time{}
	tm.logger.Info("Access token invalidated")
}

// SetupAuthentication configures the resty client with OAuth2 bearer token authentication
func SetupAuthentication(client *resty.Client, authConfig *AuthConfig, logger *zap.Logger) (*TokenManager, error) {
	if err := authConfig.Validate(); err != nil {
		logger.Error("Authentication validation failed", zap.Error(err))
		return nil, fmt.Errorf("authentication validation failed: %w", err)
	}

	// Create token manager
	tokenManager := NewTokenManager(authConfig, client, logger)

	// Get initial token
	token, err := tokenManager.RefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain initial access token: %w", err)
	}

	// Set up bearer token authentication
	client.SetAuthToken(token)

	// Add request middleware to ensure token is valid before each request
	client.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		token, err := tokenManager.GetToken()
		if err != nil {
			logger.Error("Failed to get valid token for request", zap.Error(err))
			return fmt.Errorf("failed to get valid token: %w", err)
		}
		req.SetAuthToken(token)
		return nil
	})

	logger.Info("OAuth2 authentication configured successfully",
		zap.String("instance", authConfig.Instance),
		zap.String("region", authConfig.Region),
		zap.String("scope", authConfig.GetScope()))

	return tokenManager, nil
}
