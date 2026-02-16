package client

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Client represents the HTTP client for Nexthink API
type Client struct {
	client        *resty.Client
	logger        *zap.Logger
	authConfig    *AuthConfig
	tokenManager  *TokenManager
	BaseURL       string
	globalHeaders map[string]string
	userAgent     string
}

// NewClient creates a new Nexthink API client
func NewClient(clientID, clientSecret, instance, region string, options ...ClientOption) (*Client, error) {

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
	restyClient.SetRetryWaitTime(RetryWaitTime * time.Second)
	restyClient.SetRetryMaxWaitTime(RetryMaxWaitTime * time.Second)
	restyClient.SetHeader("User-Agent", userAgent)
	restyClient.SetHeader("Accept-Encoding", "gzip")

	// Create client instance
	client := &Client{
		client:        restyClient,
		logger:        logger,
		authConfig:    authConfig,
		BaseURL:       "", // Will be set by user via options
		globalHeaders: make(map[string]string),
		userAgent:     userAgent,
	}

	// Apply any additional options before auth setup
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	// Setup OAuth2 authentication
	tokenManager, err := SetupAuthentication(restyClient, authConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup authentication: %w", err)
	}
	client.tokenManager = tokenManager

	if client.BaseURL != "" {
		restyClient.SetBaseURL(client.BaseURL)
	}

	logger.Info("Nexthink API client created",
		zap.String("instance", instance),
		zap.String("region", region),
		zap.String("base_url", client.BaseURL))

	return client, nil
}

// GetHTTPClient returns the underlying resty client
func (c *Client) GetHTTPClient() *resty.Client {
	return c.client
}

// GetLogger returns the logger instance
func (c *Client) GetLogger() *zap.Logger {
	return c.logger
}

// GetTokenManager returns the token manager
func (c *Client) GetTokenManager() *TokenManager {
	return c.tokenManager
}

// RefreshToken manually refreshes the OAuth2 access token
func (c *Client) RefreshToken() error {
	_, err := c.tokenManager.RefreshToken()
	return err
}

// InvalidateToken invalidates the current token, forcing a refresh on next use
func (c *Client) InvalidateToken() {
	c.tokenManager.InvalidateToken()
}

// QueryBuilder creates a new query builder for constructing URL parameters
func (c *Client) QueryBuilder() interfaces.ServiceQueryBuilder {
	return NewQueryBuilder()
}
