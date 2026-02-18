package client

import (
	"fmt"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Transport represents the HTTP transport layer for Nexthink API.
// It provides methods for making HTTP requests to the Nexthink API with built-in
// authentication, retry logic, and request/response logging.
// This is an internal component - users should use nexthink.NewClient() instead.
type Transport struct {
	client        *resty.Client
	logger        *zap.Logger
	authConfig    *AuthConfig
	tokenManager  *TokenManager
	BaseURL       string
	globalHeaders map[string]string
	userAgent     string
}

// NewTransport creates a new Nexthink API transport.
// This is an internal function - users should use nexthink.NewClient() instead.
func NewTransport(clientID, clientSecret, instance, region string, options ...ClientOption) (*Transport, error) {

	if err := ValidateTransportConfig(clientID, clientSecret, instance, region); err != nil {
		return nil, fmt.Errorf("invalid transport configuration: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	authConfig := &AuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Instance:     instance,
		Region:       region,
		Scope:        ScopeServiceIntegration,
	}

	// Format: "go-api-sdk-nexthink/1.0.0; gzip"
	userAgent := fmt.Sprintf("%s/%s; gzip", UserAgentBase, Version)

	// Create resty client
	restyClient := resty.New()
	restyClient.SetTimeout(DefaultTimeout * time.Second)
	restyClient.SetRetryCount(MaxRetries)
	restyClient.SetRetryWaitTime(time.Duration(RetryWaitTime) * time.Second)
	restyClient.SetRetryMaxWaitTime(time.Duration(RetryMaxWaitTime) * time.Second)
	restyClient.SetHeader("User-Agent", userAgent)
	restyClient.SetHeader("Accept-Encoding", "gzip")

	// Auto-detect proxy from environment variables (can be overridden via options)
	// Checks http_proxy, HTTP_PROXY, https_proxy, HTTPS_PROXY
	if proxyURL := os.Getenv("https_proxy"); proxyURL != "" {
		restyClient.SetProxy(proxyURL)
		logger.Info("Auto-detected HTTPS proxy from environment", zap.String("proxy", proxyURL))
	} else if proxyURL := os.Getenv("HTTPS_PROXY"); proxyURL != "" {
		restyClient.SetProxy(proxyURL)
		logger.Info("Auto-detected HTTPS proxy from environment", zap.String("proxy", proxyURL))
	} else if proxyURL := os.Getenv("http_proxy"); proxyURL != "" {
		restyClient.SetProxy(proxyURL)
		logger.Info("Auto-detected HTTP proxy from environment", zap.String("proxy", proxyURL))
	} else if proxyURL := os.Getenv("HTTP_PROXY"); proxyURL != "" {
		restyClient.SetProxy(proxyURL)
		logger.Info("Auto-detected HTTP proxy from environment", zap.String("proxy", proxyURL))
	}

	// Construct default BaseURL if not provided via options
	// Format: https://{instance}.api.{region}.nexthink.cloud
	defaultBaseURL := fmt.Sprintf("https://%s.api.%s.nexthink.cloud", instance, region)

	// Create transport instance
	transport := &Transport{
		client:        restyClient,
		logger:        logger,
		authConfig:    authConfig,
		BaseURL:       defaultBaseURL, // Default BaseURL, can be overridden via options
		globalHeaders: make(map[string]string),
		userAgent:     userAgent,
	}

	// Apply any additional options before auth setup
	for _, option := range options {
		if err := option(transport); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	// Setup OAuth2 authentication
	tokenManager, err := SetupAuthentication(restyClient, authConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup authentication: %w", err)
	}
	transport.tokenManager = tokenManager

	restyClient.SetBaseURL(transport.BaseURL)

	logger.Info("Nexthink API transport created",
		zap.String("instance", instance),
		zap.String("region", region),
		zap.String("base_url", transport.BaseURL))

	return transport, nil
}

// GetHTTPClient returns the underlying resty client
func (t *Transport) GetHTTPClient() *resty.Client {
	return t.client
}

// GetLogger returns the logger instance
func (t *Transport) GetLogger() *zap.Logger {
	return t.logger
}

// GetTokenManager returns the token manager
func (t *Transport) GetTokenManager() *TokenManager {
	return t.tokenManager
}

// RefreshToken manually refreshes the OAuth2 access token
func (t *Transport) RefreshToken() error {
	_, err := t.tokenManager.RefreshToken()
	return err
}

// InvalidateToken invalidates the current token, forcing a refresh on next use
func (t *Transport) InvalidateToken() {
	t.tokenManager.InvalidateToken()
}

// QueryBuilder creates a new query builder for constructing URL parameters
func (t *Transport) QueryBuilder() interfaces.ServiceQueryBuilder {
	return NewQueryBuilder()
}
