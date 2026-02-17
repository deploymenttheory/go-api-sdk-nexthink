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
type ClientOption func(*Transport) error

// WithBaseURL sets a custom base URL for the API client
func WithBaseURL(baseURL string) ClientOption {
	return func(t *Transport) error {
		if err := ValidateBaseURL(baseURL); err != nil {
			return fmt.Errorf("invalid base URL: %w", err)
		}
		t.BaseURL = baseURL
		t.logger.Info("Base URL configured", zap.String("base_url", baseURL))
		return nil
	}
}

// WithCustomTokenURL sets a custom OAuth2 token endpoint URL
func WithCustomTokenURL(tokenURL string) ClientOption {
	return func(t *Transport) error {
		t.authConfig.TokenURL = tokenURL
		t.logger.Info("Custom token URL configured", zap.String("token_url", tokenURL))
		return nil
	}
}

// WithScope sets a custom OAuth2 scope
func WithScope(scope string) ClientOption {
	return func(t *Transport) error {
		t.authConfig.Scope = scope
		t.logger.Info("OAuth2 scope configured", zap.String("scope", scope))
		return nil
	}
}

// WithTimeout sets a custom timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(t *Transport) error {
		timeoutSeconds := int(timeout.Seconds())
		if err := ValidateTimeout(timeoutSeconds); err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}
		t.client.SetTimeout(timeout)
		t.logger.Info("HTTP timeout configured", zap.Duration("timeout", timeout))
		return nil
	}
}

// WithRetryCount sets the number of retries for failed requests
func WithRetryCount(count int) ClientOption {
	return func(t *Transport) error {
		if err := ValidateRetryCount(count); err != nil {
			return fmt.Errorf("invalid retry count: %w", err)
		}
		t.client.SetRetryCount(count)
		t.logger.Info("Retry count configured", zap.Int("retry_count", count))
		return nil
	}
}

// WithRetryWaitTime sets the default wait time between retry attempts
// This is the initial/minimum wait time before the first retry
func WithRetryWaitTime(waitTime time.Duration) ClientOption {
	return func(t *Transport) error {
		t.client.SetRetryWaitTime(waitTime)
		t.logger.Info("Retry wait time configured", zap.Duration("wait_time", waitTime))
		return nil
	}
}

// WithRetryMaxWaitTime sets the maximum wait time between retry attempts
// The wait time increases exponentially with each retry up to this maximum
func WithRetryMaxWaitTime(maxWaitTime time.Duration) ClientOption {
	return func(t *Transport) error {
		t.client.SetRetryMaxWaitTime(maxWaitTime)
		t.logger.Info("Retry max wait time configured", zap.Duration("max_wait_time", maxWaitTime))
		return nil
	}
}

// WithLogger sets a custom logger for the client
func WithLogger(logger *zap.Logger) ClientOption {
	return func(t *Transport) error {
		t.logger = logger
		t.logger.Info("Custom logger configured")
		return nil
	}
}

// WithDebug enables debug mode which logs request and response details
func WithDebug() ClientOption {
	return func(t *Transport) error {
		t.client.SetDebug(true)
		t.logger.Info("Debug mode enabled")
		return nil
	}
}

// WithUserAgent sets a custom user agent string
func WithUserAgent(userAgent string) ClientOption {
	return func(t *Transport) error {
		t.client.SetHeader("User-Agent", userAgent)
		t.userAgent = userAgent
		t.logger.Info("User agent configured", zap.String("user_agent", userAgent))
		return nil
	}
}

// WithCustomAgent allows appending a custom identifier to the default user agent
// Format: "go-api-sdk-nexthink/1.0.0; <customAgent>; gzip"
func WithCustomAgent(customAgent string) ClientOption {
	return func(t *Transport) error {
		enhancedUA := fmt.Sprintf("%s/%s; %s; gzip", UserAgentBase, Version, customAgent)
		t.client.SetHeader("User-Agent", enhancedUA)
		t.userAgent = enhancedUA
		t.logger.Info("Custom agent configured", zap.String("user_agent", enhancedUA))
		return nil
	}
}

// WithGlobalHeader sets a global header that will be included in all requests
// Per-request headers will override global headers with the same key
func WithGlobalHeader(key, value string) ClientOption {
	return func(t *Transport) error {
		t.globalHeaders[key] = value
		t.logger.Info("Global header configured", zap.String("key", key), zap.String("value", value))
		return nil
	}
}

// WithGlobalHeaders sets multiple global headers at once
func WithGlobalHeaders(headers map[string]string) ClientOption {
	return func(t *Transport) error {
		maps.Copy(t.globalHeaders, headers)
		t.logger.Info("Multiple global headers configured", zap.Int("count", len(headers)))
		return nil
	}
}

// WithProxy sets an HTTP proxy for all requests
// Example: "http://proxy.company.com:8080" or "socks5://127.0.0.1:1080"
func WithProxy(proxyURL string) ClientOption {
	return func(t *Transport) error {
		if err := ValidateProxyURL(proxyURL); err != nil {
			return fmt.Errorf("invalid proxy URL: %w", err)
		}
		t.client.SetProxy(proxyURL)
		t.logger.Info("Proxy configured", zap.String("proxy", proxyURL))
		return nil
	}
}

// WithTLSClientConfig sets custom TLS configuration
// Use this for custom certificate validation, minimum TLS version, etc.
func WithTLSClientConfig(tlsConfig *tls.Config) ClientOption {
	return func(t *Transport) error {
		t.client.SetTLSClientConfig(tlsConfig)
		t.logger.Info("TLS client config configured",
			zap.Uint16("min_version", tlsConfig.MinVersion),
			zap.Bool("insecure_skip_verify", tlsConfig.InsecureSkipVerify))
		return nil
	}
}

// WithClientCertificate sets a client certificate for mutual TLS authentication
// Loads certificate from PEM-encoded files
func WithClientCertificate(certFile, keyFile string) ClientOption {
	return func(t *Transport) error {
		t.client.SetCertificateFromFile(certFile, keyFile)
		t.logger.Info("Client certificate configured",
			zap.String("cert_file", certFile),
			zap.String("key_file", keyFile))
		return nil
	}
}

// WithClientCertificateFromString sets a client certificate from PEM-encoded strings
func WithClientCertificateFromString(certPEM, keyPEM string) ClientOption {
	return func(t *Transport) error {
		t.client.SetCertificateFromString(certPEM, keyPEM)
		t.logger.Info("Client certificate configured from string")
		return nil
	}
}

// WithRootCertificates adds custom root CA certificates for server validation
// Useful for private CAs or self-signed certificates
func WithRootCertificates(pemFilePaths ...string) ClientOption {
	return func(t *Transport) error {
		t.client.SetClientRootCertificates(pemFilePaths...)
		t.logger.Info("Root certificates configured", zap.Int("count", len(pemFilePaths)))
		return nil
	}
}

// WithRootCertificateFromString adds a custom root CA certificate from PEM string
func WithRootCertificateFromString(pemContent string) ClientOption {
	return func(t *Transport) error {
		t.client.SetClientRootCertificateFromString(pemContent)
		t.logger.Info("Root certificate configured from string")
		return nil
	}
}

// WithTransport sets a custom HTTP transport (http.RoundTripper)
// Use this for advanced transport customization beyond TLS/proxy
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(t *Transport) error {
		t.client.SetTransport(transport)
		t.logger.Info("Custom transport configured")
		return nil
	}
}

// WithInsecureSkipVerify disables TLS certificate verification (USE WITH CAUTION)
// This should ONLY be used for testing/development with self-signed certificates
func WithInsecureSkipVerify() ClientOption {
	return func(t *Transport) error {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		t.client.SetTLSClientConfig(tlsConfig)
		t.logger.Warn("TLS certificate verification DISABLED - use only for testing")
		return nil
	}
}

// WithMinTLSVersion sets the minimum TLS version for connections
// Common values: tls.VersionTLS12, tls.VersionTLS13
func WithMinTLSVersion(minVersion uint16) ClientOption {
	return func(t *Transport) error {
		tlsConfig := &tls.Config{
			MinVersion: minVersion,
		}
		t.client.SetTLSClientConfig(tlsConfig)

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

		t.logger.Info("Minimum TLS version configured",
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
	return func(t *Transport) error {
		return t.EnableTracing(config)
	}
}
