package client

import (
	"testing"

	"resty.dev/v3"
)

func TestApplyHeaders_GlobalOnly(t *testing.T) {
	// Create a test client with global headers
	transport := &Transport{
		client: resty.New(),
		globalHeaders: map[string]string{
			"X-Global-Header-1": "global-value-1",
			"X-Global-Header-2": "global-value-2",
		},
	}

	// Create a request
	req := transport.client.R()

	// Apply headers with no per-request headers
	transport.applyHeaders(req, nil)

	// Verify global headers are set
	header := req.Header
	if got := header.Get("X-Global-Header-1"); got != "global-value-1" {
		t.Errorf("X-Global-Header-1 = %q, want %q", got, "global-value-1")
	}
	if got := header.Get("X-Global-Header-2"); got != "global-value-2" {
		t.Errorf("X-Global-Header-2 = %q, want %q", got, "global-value-2")
	}
}

func TestApplyHeaders_RequestOnly(t *testing.T) {
	// Create a test client with no global headers
	transport := &Transport{
		client:        resty.New(),
		globalHeaders: make(map[string]string),
	}

	// Create a request
	req := transport.client.R()

	// Apply headers with per-request headers
	requestHeaders := map[string]string{
		"X-Request-Header-1": "request-value-1",
		"X-Request-Header-2": "request-value-2",
	}
	transport.applyHeaders(req, requestHeaders)

	// Verify request headers are set
	header := req.Header
	if got := header.Get("X-Request-Header-1"); got != "request-value-1" {
		t.Errorf("X-Request-Header-1 = %q, want %q", got, "request-value-1")
	}
	if got := header.Get("X-Request-Header-2"); got != "request-value-2" {
		t.Errorf("X-Request-Header-2 = %q, want %q", got, "request-value-2")
	}
}

func TestApplyHeaders_Override(t *testing.T) {
	// Create a test client with global headers
	transport := &Transport{
		client: resty.New(),
		globalHeaders: map[string]string{
			"X-Shared-Header":  "global-value",
			"X-Global-Only":    "global-only-value",
			"Content-Type":     "application/json",
			"X-Another-Global": "another-global",
		},
	}

	// Create a request
	req := transport.client.R()

	// Apply headers with per-request headers that override some global headers
	requestHeaders := map[string]string{
		"X-Shared-Header": "request-value-overrides-global",
		"X-Request-Only":  "request-only-value",
		"Content-Type":    "text/csv",
	}
	transport.applyHeaders(req, requestHeaders)

	// Verify request headers override global headers
	header := req.Header
	if got := header.Get("X-Shared-Header"); got != "request-value-overrides-global" {
		t.Errorf("X-Shared-Header = %q, want %q (should be overridden)", got, "request-value-overrides-global")
	}
	if got := header.Get("Content-Type"); got != "text/csv" {
		t.Errorf("Content-Type = %q, want %q (should be overridden)", got, "text/csv")
	}

	// Verify global-only headers are still set
	if got := header.Get("X-Global-Only"); got != "global-only-value" {
		t.Errorf("X-Global-Only = %q, want %q", got, "global-only-value")
	}
	if got := header.Get("X-Another-Global"); got != "another-global" {
		t.Errorf("X-Another-Global = %q, want %q", got, "another-global")
	}

	// Verify request-only headers are set
	if got := header.Get("X-Request-Only"); got != "request-only-value" {
		t.Errorf("X-Request-Only = %q, want %q", got, "request-only-value")
	}
}

func TestApplyHeaders_EmptyValues(t *testing.T) {
	// Create a test client with global headers including empty values
	transport := &Transport{
		client: resty.New(),
		globalHeaders: map[string]string{
			"X-Valid-Header": "valid-value",
			"X-Empty-Global": "",
		},
	}

	// Create a request
	req := transport.client.R()

	// Apply headers with per-request headers including empty values
	requestHeaders := map[string]string{
		"X-Request-Header": "request-value",
		"X-Empty-Request":  "",
	}
	transport.applyHeaders(req, requestHeaders)

	// Verify valid headers are set
	header := req.Header
	if got := header.Get("X-Valid-Header"); got != "valid-value" {
		t.Errorf("X-Valid-Header = %q, want %q", got, "valid-value")
	}
	if got := header.Get("X-Request-Header"); got != "request-value" {
		t.Errorf("X-Request-Header = %q, want %q", got, "request-value")
	}

	// Verify empty headers are not set
	if got := header.Get("X-Empty-Global"); got != "" {
		t.Errorf("X-Empty-Global should not be set, got %q", got)
	}
	if got := header.Get("X-Empty-Request"); got != "" {
		t.Errorf("X-Empty-Request should not be set, got %q", got)
	}
}

func TestApplyHeaders_NilMaps(t *testing.T) {
	// Create a test client with nil global headers
	transport := &Transport{
		client:        resty.New(),
		globalHeaders: nil,
	}

	// Create a request
	req := transport.client.R()

	// Apply headers with nil request headers - should not panic
	transport.applyHeaders(req, nil)

	// No headers should be set
	if len(req.Header) > 0 {
		t.Errorf("Expected no headers to be set, got %d headers", len(req.Header))
	}
}

func TestApplyHeaders_EmptyMaps(t *testing.T) {
	// Create a test client with empty global headers
	transport := &Transport{
		client:        resty.New(),
		globalHeaders: make(map[string]string),
	}

	// Create a request
	req := transport.client.R()

	// Apply headers with empty request headers
	requestHeaders := make(map[string]string)
	transport.applyHeaders(req, requestHeaders)

	// No headers should be set
	if len(req.Header) > 0 {
		t.Errorf("Expected no headers to be set, got %d headers", len(req.Header))
	}
}

func TestApplyHeaders_Precedence(t *testing.T) {
	// Test that demonstrates the precedence order explicitly
	tests := []struct {
		name           string
		globalHeaders  map[string]string
		requestHeaders map[string]string
		checkHeader    string
		expectedValue  string
		description    string
	}{
		{
			name: "request overrides global",
			globalHeaders: map[string]string{
				"Authorization": "Bearer global-token",
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer request-token",
			},
			checkHeader:   "Authorization",
			expectedValue: "Bearer request-token",
			description:   "Per-request Authorization should override global",
		},
		{
			name: "global only",
			globalHeaders: map[string]string{
				"X-API-Key": "global-api-key",
			},
			requestHeaders: map[string]string{},
			checkHeader:    "X-API-Key",
			expectedValue:  "global-api-key",
			description:    "Global header should be used when no request header provided",
		},
		{
			name:          "request only",
			globalHeaders: map[string]string{},
			requestHeaders: map[string]string{
				"X-Request-ID": "req-123",
			},
			checkHeader:   "X-Request-ID",
			expectedValue: "req-123",
			description:   "Request header should be used when no global header exists",
		},
		{
			name: "empty request value doesn't override",
			globalHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			requestHeaders: map[string]string{
				"Content-Type": "",
			},
			checkHeader:   "Content-Type",
			expectedValue: "application/json",
			description:   "Empty request header value should not override global",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &Transport{
				client:        resty.New(),
				globalHeaders: tt.globalHeaders,
			}

			req := transport.client.R()
			transport.applyHeaders(req, tt.requestHeaders)

			got := req.Header.Get(tt.checkHeader)
			if got != tt.expectedValue {
				t.Errorf("%s: %s = %q, want %q", tt.description, tt.checkHeader, got, tt.expectedValue)
			}
		})
	}
}

func TestApplyHeaders_CaseInsensitivity(t *testing.T) {
	// HTTP headers are case-insensitive, verify this works correctly
	transport := &Transport{
		client: resty.New(),
		globalHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	req := transport.client.R()
	transport.applyHeaders(req, nil)

	// HTTP headers are canonicalized by net/http, so we can check different cases
	if got := req.Header.Get("content-type"); got != "application/json" {
		t.Errorf("Header case-insensitive check failed: got %q, want %q", got, "application/json")
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Header case-insensitive check failed: got %q, want %q", got, "application/json")
	}
	if got := req.Header.Get("CONTENT-TYPE"); got != "application/json" {
		t.Errorf("Header case-insensitive check failed: got %q, want %q", got, "application/json")
	}
}

func TestApplyHeaders_CommonHeaders(t *testing.T) {
	// Test common HTTP headers that might be used
	transport := &Transport{
		client: resty.New(),
		globalHeaders: map[string]string{
			"User-Agent":    "NexthinkSDK/1.0",
			"Accept":        "application/json",
			"Authorization": "Bearer global-token",
		},
	}

	req := transport.client.R()

	requestHeaders := map[string]string{
		"X-Request-ID":  "12345",
		"Authorization": "Bearer request-token",
		"Content-Type":  "application/json",
	}

	transport.applyHeaders(req, requestHeaders)

	// Check all headers are properly set
	expectedHeaders := map[string]string{
		"User-Agent":    "NexthinkSDK/1.0",
		"Accept":        "application/json",
		"Authorization": "Bearer request-token",
		"X-Request-ID":  "12345",
		"Content-Type":  "application/json",
	}

	for key, expected := range expectedHeaders {
		if got := req.Header.Get(key); got != expected {
			t.Errorf("Header %q = %q, want %q", key, got, expected)
		}
	}
}
