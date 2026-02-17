package client

import (
	"fmt"
	"strings"
)

// ValidateTransportConfig validates the transport configuration parameters
func ValidateTransportConfig(clientID, clientSecret, instance, region string) error {
	if clientID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	if clientSecret == "" {
		return fmt.Errorf("client secret cannot be empty")
	}

	if instance == "" {
		return fmt.Errorf("instance cannot be empty")
	}

	if region == "" {
		return fmt.Errorf("region cannot be empty")
	}

	// Validate region is one of the supported regions
	validRegions := map[string]bool{
		RegionUS:   true,
		RegionEU:   true,
		RegionPAC:  true,
		RegionMETA: true,
	}

	if !validRegions[region] {
		return fmt.Errorf("invalid region '%s': must be one of: us, eu, pac, meta", region)
	}

	if strings.Contains(instance, " ") {
		return fmt.Errorf("instance name cannot contain spaces")
	}

	if len(instance) > 100 {
		return fmt.Errorf("instance name too long (max 100 characters)")
	}

	return nil
}

// ValidateBaseURL validates a base URL format
func ValidateBaseURL(baseURL string) error {
	if baseURL == "" {
		return fmt.Errorf("base URL cannot be empty")
	}

	if !strings.HasPrefix(baseURL, "https://") && !strings.HasPrefix(baseURL, "http://") {
		return fmt.Errorf("base URL must start with http:// or https://")
	}

	if strings.HasSuffix(baseURL, "/") {
		return fmt.Errorf("base URL should not end with a trailing slash")
	}

	return nil
}

// ValidateTimeout validates timeout value
func ValidateTimeout(timeout int) error {
	if timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	if timeout > 3600 {
		return fmt.Errorf("timeout too large (max 3600 seconds)")
	}

	return nil
}

// ValidateRetryCount validates retry count
func ValidateRetryCount(retryCount int) error {
	if retryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}

	if retryCount > 10 {
		return fmt.Errorf("retry count too large (max 10)")
	}

	return nil
}

// ValidateProxyURL validates proxy URL format
func ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil // Empty is valid (no proxy)
	}

	if !strings.HasPrefix(proxyURL, "http://") &&
		!strings.HasPrefix(proxyURL, "https://") &&
		!strings.HasPrefix(proxyURL, "socks5://") {
		return fmt.Errorf("proxy URL must start with http://, https://, or socks5://")
	}

	return nil
}
