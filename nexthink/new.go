package nexthink

import (
	"fmt"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/campaigns"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/enrichment"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/remote_actions"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/workflows"
	"go.uber.org/zap"
)

// Client is the main entry point for the Nexthink API SDK.
// It aggregates all service clients and provides a unified interface.
// Users should interact with the API exclusively through the provided service methods.
type Client struct {
	// transport is the internal HTTP transport layer (not exposed to users)
	transport *client.Transport

	// Services - users should only call methods on these services
	Campaigns     *campaigns.Service
	Enrichment    *enrichment.Service
	NQL           *nql.Service
	RemoteActions *remote_actions.Service
	Workflows     *workflows.Service
}

// NewClient creates a new Nexthink API client
//
// Parameters:
//   - clientID: The OAuth2 client ID
//   - clientSecret: The OAuth2 client secret
//   - instance: The Nexthink instance name
//   - region: The region (us, eu, pac, meta)
//   - options: Optional client configuration options
//
// Example:
//
//	client, err := nexthink.NewClient(
//	    "your-client-id",
//	    "your-client-secret",
//	    "your-instance",
//	    "us",
//	    nexthink.WithDebug(),
//	)
func NewClient(clientID, clientSecret, instance, region string, options ...client.ClientOption) (*Client, error) {
	// Create base HTTP transport
	transport, err := client.NewTransport(clientID, clientSecret, instance, region, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP transport: %w", err)
	}

	// Initialize service clients
	c := &Client{
		transport:     transport,
		Campaigns:     campaigns.NewService(transport),
		Enrichment:    enrichment.NewService(transport),
		NQL:           nql.NewService(transport),
		RemoteActions: remote_actions.NewService(transport),
		Workflows:     workflows.NewService(transport),
	}

	return c, nil
}

// NewClientFromEnv creates a new client using environment variables
//
// Required environment variables:
//   - NEXTHINK_CLIENT_ID: The OAuth2 client ID
//   - NEXTHINK_CLIENT_SECRET: The OAuth2 client secret
//   - NEXTHINK_INSTANCE: The Nexthink instance name
//   - NEXTHINK_REGION: The region (us, eu, pac, meta)
//
// Optional environment variables:
//   - NEXTHINK_BASE_URL: Custom base URL (overrides default)
//
// Example:
//
//	client, err := nexthink.NewClientFromEnv()
func NewClientFromEnv(options ...client.ClientOption) (*Client, error) {
	clientID := os.Getenv("NEXTHINK_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("NEXTHINK_CLIENT_ID environment variable is required")
	}

	clientSecret := os.Getenv("NEXTHINK_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("NEXTHINK_CLIENT_SECRET environment variable is required")
	}

	instance := os.Getenv("NEXTHINK_INSTANCE")
	if instance == "" {
		return nil, fmt.Errorf("NEXTHINK_INSTANCE environment variable is required")
	}

	region := os.Getenv("NEXTHINK_REGION")
	if region == "" {
		return nil, fmt.Errorf("NEXTHINK_REGION environment variable is required")
	}

	// Check for optional environment variables and append to options
	if baseURL := os.Getenv("NEXTHINK_BASE_URL"); baseURL != "" {
		options = append(options, client.WithBaseURL(baseURL))
	}

	return NewClient(clientID, clientSecret, instance, region, options...)
}

// GetLogger returns the configured zap logger instance.
// Use this to add custom logging within your application using the same logger.
//
// Returns:
//   - *zap.Logger: The configured logger instance
func (c *Client) GetLogger() *zap.Logger {
	return c.transport.GetLogger()
}

// GetTokenManager returns the token manager instance for advanced token operations.
// This allows access to low-level token management functionality when needed.
//
// Returns:
//   - *client.TokenManager: The token manager instance
func (c *Client) GetTokenManager() *client.TokenManager {
	return c.transport.GetTokenManager()
}

// RefreshToken manually refreshes the OAuth2 access token.
// This can be useful when you need to explicitly refresh the token.
//
// Returns:
//   - error: Any error encountered during token refresh
func (c *Client) RefreshToken() error {
	return c.transport.RefreshToken()
}

// InvalidateToken invalidates the current token, forcing a refresh on next use.
// This is useful when you know the token is no longer valid.
func (c *Client) InvalidateToken() {
	c.transport.InvalidateToken()
}
