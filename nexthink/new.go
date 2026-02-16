package nexthink

import (
	"fmt"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/campaigns"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/enrichment"
)

// Client is the main entry point for the Nexthink API SDK
// It aggregates all service clients and provides a unified interface
type Client struct {
	*client.Client

	// Services
	Campaigns  *campaigns.Service
	Enrichment *enrichment.Service
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
	// Create base HTTP client
	httpClient, err := client.NewClient(clientID, clientSecret, instance, region, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Initialize service clients
	c := &Client{
		Client:     httpClient,
		Campaigns:  campaigns.NewService(httpClient),
		Enrichment: enrichment.NewService(httpClient),
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
