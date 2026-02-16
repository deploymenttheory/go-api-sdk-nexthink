package client

import (
	"crypto/tls"
	"fmt"
	"maps"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ClientOption is a function type for configuring the Client
type ClientOption func(*Client) error

// WithBaseURL sets a custom base URL for the API client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.BaseURL = baseURL
		c.logger.Info("Base URL configured", zap.String("base_url", baseURL))
		return nil
	}
}

// WithCustomTokenURL sets a custom OAuth2 token endpoint URL
func WithCustomTokenURL(tokenURL string) ClientOption {
	return func(c *Client) error {
		c.authConfig.TokenURL = tokenURL
		c.logger.Info("Custom token URL configured", zap.String("token_url", tokenURL))
		return nil
	}
}

// WithScope sets a custom OAuth2 scope
func WithScope(scope string) ClientOption {
	return func(c *Client) error {
		c.authConfig.Scope = scope
		c.logger.Info("OAuth2 scope configured", zap.String("scope", scope))
		return nil
	}
}

// WithTimeout sets a custom timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.client.SetTimeout(timeout)
		c.logger.Info("HTTP timeout configured", zap.Duration("timeout", timeout))
		return nil
	}
}

// WithRetryCount sets the number of retries for failed requests
func WithRetryCount(count int) ClientOption {
	return func(c *Client) error {
		c.client.SetRetryCount(count)
		c.logger.Info("Retry count configured", zap.Int("retry_count", count))
		return nil
	}
}

// WithRetryWaitTime sets the default wait time between retry attempts
// This is the initial/minimum wait time before the first retry
func WithRetryWaitTime(waitTime time.Duration) ClientOption {
	return func(c *Client) error {
		c.client.SetRetryWaitTime(waitTime)
		c.logger.Info("Retry wait time configured", zap.Duration("wait_time", waitTime))
		return nil
	}
}

// WithRetryMaxWaitTime sets the maximum wait time between retry attempts
// The wait time increases exponentially with each retry up to this maximum
func WithRetryMaxWaitTime(maxWaitTime time.Duration) ClientOption {
	return func(c *Client) error {
		c.client.SetRetryMaxWaitTime(maxWaitTime)
		c.logger.Info("Retry max wait time configured", zap.Duration("max_wait_time", maxWaitTime))
		return nil
	}
}

// WithLogger sets a custom logger for the client
func WithLogger(logger *zap.Logger) ClientOption {
	return func(c *Client) error {
		c.logger = logger
		c.logger.Info("Custom logger configured")
		return nil
	}
}

// WithDebug enables debug mode which logs request and response details
func WithDebug() ClientOption {
	return func(c *Client) error {
		c.client.SetDebug(true)
		c.logger.Info("Debug mode enabled")
		return nil
	}
}

// WithUserAgent sets a custom user agent string
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.client.SetHeader("User-Agent", userAgent)
		c.userAgent = userAgent
		c.logger.Info("User agent configured", zap.String("user_agent", userAgent))
		return nil
	}
}

// WithCustomAgent allows appending a custom identifier to the default user agent
// Format: "go-api-sdk-nexthink/1.0.0; <customAgent>; gzip"
func WithCustomAgent(customAgent string) ClientOption {
	return func(c *Client) error {
		enhancedUA := fmt.Sprintf("%s/%s; %s; gzip", UserAgentBase, Version, customAgent)
		c.client.SetHeader("User-Agent", enhancedUA)
		c.userAgent = enhancedUA
		c.logger.Info("Custom agent configured", zap.String("user_agent", enhancedUA))
		return nil
	}
}

// WithGlobalHeader sets a global header that will be included in all requests
// Per-request headers will override global headers with the same key
func WithGlobalHeader(key, value string) ClientOption {
	return func(c *Client) error {
		c.globalHeaders[key] = value
		c.logger.Info("Global header configured", zap.String("key", key), zap.String("value", value))
		return nil
	}
}

// WithGlobalHeaders sets multiple global headers at once
func WithGlobalHeaders(headers map[string]string) ClientOption {
	return func(c *Client) error {
		maps.Copy(c.globalHeaders, headers)
		c.logger.Info("Multiple global headers configured", zap.Int("count", len(headers)))
		return nil
	}
}

// WithProxy sets an HTTP proxy for all requests
// Example: "http://proxy.company.com:8080" or "socks5://127.0.0.1:1080"
func WithProxy(proxyURL string) ClientOption {
	return func(c *Client) error {
		c.client.SetProxy(proxyURL)
		c.logger.Info("Proxy configured", zap.String("proxy", proxyURL))
		return nil
	}
}

// WithTLSClientConfig sets custom TLS configuration
// Use this for custom certificate validation, minimum TLS version, etc.
func WithTLSClientConfig(tlsConfig *tls.Config) ClientOption {
	return func(c *Client) error {
		c.client.SetTLSClientConfig(tlsConfig)
		c.logger.Info("TLS client config configured",
			zap.Uint16("min_version", tlsConfig.MinVersion),
			zap.Bool("insecure_skip_verify", tlsConfig.InsecureSkipVerify))
		return nil
	}
}

// WithClientCertificate sets a client certificate for mutual TLS authentication
// Loads certificate from PEM-encoded files
func WithClientCertificate(certFile, keyFile string) ClientOption {
	return func(c *Client) error {
		c.client.SetCertificateFromFile(certFile, keyFile)
		c.logger.Info("Client certificate configured",
			zap.String("cert_file", certFile),
			zap.String("key_file", keyFile))
		return nil
	}
}

// WithClientCertificateFromString sets a client certificate from PEM-encoded strings
func WithClientCertificateFromString(certPEM, keyPEM string) ClientOption {
	return func(c *Client) error {
		c.client.SetCertificateFromString(certPEM, keyPEM)
		c.logger.Info("Client certificate configured from string")
		return nil
	}
}

// WithRootCertificates adds custom root CA certificates for server validation
// Useful for private CAs or self-signed certificates
func WithRootCertificates(pemFilePaths ...string) ClientOption {
	return func(c *Client) error {
		c.client.SetClientRootCertificates(pemFilePaths...)
		c.logger.Info("Root certificates configured", zap.Int("count", len(pemFilePaths)))
		return nil
	}
}

// WithRootCertificateFromString adds a custom root CA certificate from PEM string
func WithRootCertificateFromString(pemContent string) ClientOption {
	return func(c *Client) error {
		c.client.SetClientRootCertificateFromString(pemContent)
		c.logger.Info("Root certificate configured from string")
		return nil
	}
}

// WithTransport sets a custom HTTP transport (http.RoundTripper)
// Use this for advanced transport customization beyond TLS/proxy
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *Client) error {
		c.client.SetTransport(transport)
		c.logger.Info("Custom transport configured")
		return nil
	}
}

// WithInsecureSkipVerify disables TLS certificate verification (USE WITH CAUTION)
// This should ONLY be used for testing/development with self-signed certificates
func WithInsecureSkipVerify() ClientOption {
	return func(c *Client) error {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		c.client.SetTLSClientConfig(tlsConfig)
		c.logger.Warn("TLS certificate verification DISABLED - use only for testing")
		return nil
	}
}

// WithMinTLSVersion sets the minimum TLS version for connections
// Common values: tls.VersionTLS12, tls.VersionTLS13
func WithMinTLSVersion(minVersion uint16) ClientOption {
	return func(c *Client) error {
		tlsConfig := &tls.Config{
			MinVersion: minVersion,
		}
		c.client.SetTLSClientConfig(tlsConfig)

		versionName := "unknown"
		switch minVersion {
		case tls.VersionTLS10:
			versionName = "TLS 1.0"
		case tls.VersionTLS11:
			versionName = "TLS 1.1"
		case tls.VersionTLS12:
			versionName = "TLS 1.2"
		case tls.VersionTLS13:
			versionName = "TLS 1.3"
		}

		c.logger.Info("Minimum TLS version configured",
			zap.String("version", versionName),
			zap.Uint16("version_code", minVersion))
		return nil
	}
}

// WithTracing enables OpenTelemetry tracing for all HTTP requests.
// This wraps the HTTP client transport with automatic instrumentation.
//
// Example usage:
//
//	client, err := client.NewClient(clientID, clientSecret, instance, region,
//	    client.WithTracing(nil), // Uses default config with global tracer provider
//	)
//
// Or with custom configuration:
//
//	otelConfig := &client.OTelConfig{
//	    TracerProvider: myTracerProvider,
//	    ServiceName:    "my-nexthink-client",
//	}
//	client, err := client.NewClient(clientID, clientSecret, instance, region,
//	    client.WithTracing(otelConfig),
//	)
//
// The instrumentation automatically captures:
// - HTTP method, URL, status code
// - Request/response timing
// - Error details
// - All OpenTelemetry semantic conventions for HTTP
func WithTracing(config *OTelConfig) ClientOption {
	return func(c *Client) error {
		return c.EnableTracing(config)
	}
}
