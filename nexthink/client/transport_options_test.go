package client

import (
	"crypto/tls"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
)

func TestWithBaseURL(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			baseURL: "https://myinstance.api.us.nexthink.cloud",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "not a url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithBaseURL(tt.baseURL),
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && client != nil && client.BaseURL != tt.baseURL {
				t.Errorf("BaseURL = %s, want %s", client.BaseURL, tt.baseURL)
			}
		})
	}
}

func TestWithCustomTokenURL(t *testing.T) {
	logger := zaptest.NewLogger(t)

	customTokenURL := "https://custom-auth.example.com/oauth2/token"

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithCustomTokenURL(customTokenURL),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if client.authConfig.TokenURL != customTokenURL {
		t.Errorf("TokenURL = %s, want %s", client.authConfig.TokenURL, customTokenURL)
	}
}

func TestWithScope(t *testing.T) {
	logger := zaptest.NewLogger(t)

	customScope := "custom:scope"

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithScope(customScope),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if client.authConfig.Scope != customScope {
		t.Errorf("Scope = %s, want %s", client.authConfig.Scope, customScope)
	}
}

func TestWithTimeout(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "valid 30 second timeout",
			timeout: 30 * time.Second,
			wantErr: false,
		},
		{
			name:    "valid 60 second timeout",
			timeout: 60 * time.Second,
			wantErr: false,
		},
		{
			name:    "invalid negative timeout",
			timeout: -1 * time.Second,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithTimeout(tt.timeout),
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithRetryCount(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{
			name:    "valid retry count 3",
			count:   3,
			wantErr: false,
		},
		{
			name:    "valid retry count 5",
			count:   5,
			wantErr: false,
		},
		{
			name:    "invalid negative count",
			count:   -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithRetryCount(tt.count),
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithRetryWaitTime(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name     string
		waitTime time.Duration
	}{
		{
			name:     "1 second wait",
			waitTime: 1 * time.Second,
		},
		{
			name:     "5 second wait",
			waitTime: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithRetryWaitTime(tt.waitTime),
			)

			if err != nil {
				t.Fatalf("NewTransport() error = %v, want nil", err)
			}

			if client == nil {
				t.Fatal("NewTransport() returned nil client")
			}
		})
	}
}

func TestWithRetryMaxWaitTime(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name        string
		maxWaitTime time.Duration
	}{
		{
			name:        "5 second max wait",
			maxWaitTime: 5 * time.Second,
		},
		{
			name:        "30 second max wait",
			maxWaitTime: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithRetryMaxWaitTime(tt.maxWaitTime),
			)

			if err != nil {
				t.Fatalf("NewTransport() error = %v, want nil", err)
			}

			if client == nil {
				t.Fatal("NewTransport() returned nil client")
			}
		})
	}
}

func TestWithRetryConfiguration_Combined(t *testing.T) {
	logger := zaptest.NewLogger(t)

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithRetryCount(5),
		WithRetryWaitTime(3*time.Second),
		WithRetryMaxWaitTime(30*time.Second),
	)

	if err != nil {
		t.Fatalf("NewTransport() with combined retry options error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewTransport() returned nil client")
	}
}

func TestWithDebug(t *testing.T) {
	logger := zaptest.NewLogger(t)

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithDebug(),
	)

	if err != nil {
		t.Fatalf("NewTransport() with debug error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewTransport() returned nil client")
	}
}

func TestWithUserAgent(t *testing.T) {
	logger := zaptest.NewLogger(t)

	customUA := "CustomUserAgent/1.0"

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithUserAgent(customUA),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if client.userAgent != customUA {
		t.Errorf("userAgent = %s, want %s", client.userAgent, customUA)
	}
}

func TestWithCustomAgent(t *testing.T) {
	logger := zaptest.NewLogger(t)

	customAgent := "MyApp/2.0"

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithCustomAgent(customAgent),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	// Should contain the custom agent
	if !strings.Contains(client.userAgent, customAgent) {
		t.Errorf("userAgent = %s, want to contain %s", client.userAgent, customAgent)
	}
}

func TestWithGlobalHeader(t *testing.T) {
	logger := zaptest.NewLogger(t)

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithGlobalHeader("X-Custom-Header", "custom-value"),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if client.globalHeaders["X-Custom-Header"] != "custom-value" {
		t.Errorf("globalHeaders[X-Custom-Header] = %s, want %s",
			client.globalHeaders["X-Custom-Header"], "custom-value")
	}
}

func TestWithGlobalHeaders(t *testing.T) {
	logger := zaptest.NewLogger(t)

	headers := map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	}

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithGlobalHeaders(headers),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	for k, v := range headers {
		if client.globalHeaders[k] != v {
			t.Errorf("globalHeaders[%s] = %s, want %s", k, client.globalHeaders[k], v)
		}
	}
}

func TestWithProxy(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name     string
		proxyURL string
		wantErr  bool
	}{
		{
			name:     "valid HTTP proxy",
			proxyURL: "http://proxy.example.com:8080",
			wantErr:  false,
		},
		{
			name:     "invalid proxy URL",
			proxyURL: "not a url",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithProxy(tt.proxyURL),
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithTLSClientConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTLSClientConfig(tlsConfig),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewTransport() returned nil client")
	}
}

func TestWithInsecureSkipVerify(t *testing.T) {
	logger := zaptest.NewLogger(t)

	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithInsecureSkipVerify(),
	)

	if err != nil {
		t.Fatalf("NewTransport() error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewTransport() returned nil client")
	}
}

func TestWithMinTLSVersion(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		minVersion uint16
	}{
		{
			name:       "TLS 1.2",
			minVersion: tls.VersionTLS12,
		},
		{
			name:       "TLS 1.3",
			minVersion: tls.VersionTLS13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewTransport(
				"test-id",
				"test-secret",
				"test-instance",
				RegionUS,
				WithLogger(logger),
				WithMinTLSVersion(tt.minVersion),
			)

			if err != nil {
				t.Fatalf("NewTransport() error = %v, want nil", err)
			}

			if client == nil {
				t.Fatal("NewTransport() returned nil client")
			}
		})
	}
}

func TestMultipleOptions(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Test combining multiple options
	client, err := NewTransport(
		"test-id",
		"test-secret",
		"test-instance",
		RegionUS,
		WithLogger(logger),
		WithTimeout(30*time.Second),
		WithRetryCount(3),
		WithRetryWaitTime(2*time.Second),
		WithRetryMaxWaitTime(20*time.Second),
		WithUserAgent("CustomUA/1.0"),
		WithGlobalHeader("X-Test", "test-value"),
		WithDebug(),
	)

	if err != nil {
		t.Fatalf("NewTransport() with multiple options error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewTransport() returned nil client")
	}

	// Verify some of the options were applied
	if client.userAgent != "CustomUA/1.0" {
		t.Errorf("userAgent = %s, want CustomUA/1.0", client.userAgent)
	}

	if client.globalHeaders["X-Test"] != "test-value" {
		t.Errorf("globalHeaders[X-Test] = %s, want test-value", client.globalHeaders["X-Test"])
	}
}
