package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

func TestIsResponseSuccess(t *testing.T) {
	tests := []struct {
		name     string
		resp     *interfaces.Response
		expected bool
	}{
		{
			name:     "nil response",
			resp:     nil,
			expected: false,
		},
		{
			name:     "200 OK",
			resp:     &interfaces.Response{StatusCode: 200},
			expected: true,
		},
		{
			name:     "201 Created",
			resp:     &interfaces.Response{StatusCode: 201},
			expected: true,
		},
		{
			name:     "207 Multi-Status",
			resp:     &interfaces.Response{StatusCode: 207},
			expected: true,
		},
		{
			name:     "299 edge case",
			resp:     &interfaces.Response{StatusCode: 299},
			expected: true,
		},
		{
			name:     "400 Bad Request",
			resp:     &interfaces.Response{StatusCode: 400},
			expected: false,
		},
		{
			name:     "404 Not Found",
			resp:     &interfaces.Response{StatusCode: 404},
			expected: false,
		},
		{
			name:     "500 Internal Server Error",
			resp:     &interfaces.Response{StatusCode: 500},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsResponseSuccess(tt.resp)
			if result != tt.expected {
				t.Errorf("IsResponseSuccess() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsResponseError(t *testing.T) {
	tests := []struct {
		name     string
		resp     *interfaces.Response
		expected bool
	}{
		{
			name:     "nil response",
			resp:     nil,
			expected: false,
		},
		{
			name:     "200 OK",
			resp:     &interfaces.Response{StatusCode: 200},
			expected: false,
		},
		{
			name:     "400 Bad Request",
			resp:     &interfaces.Response{StatusCode: 400},
			expected: true,
		},
		{
			name:     "404 Not Found",
			resp:     &interfaces.Response{StatusCode: 404},
			expected: true,
		},
		{
			name:     "429 Too Many Requests",
			resp:     &interfaces.Response{StatusCode: 429},
			expected: true,
		},
		{
			name:     "500 Internal Server Error",
			resp:     &interfaces.Response{StatusCode: 500},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsResponseError(tt.resp)
			if result != tt.expected {
				t.Errorf("IsResponseError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetResponseHeader(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Custom-Header", "test-value")

	tests := []struct {
		name     string
		resp     *interfaces.Response
		key      string
		expected string
	}{
		{
			name:     "nil response",
			resp:     nil,
			key:      "Content-Type",
			expected: "",
		},
		{
			name:     "nil headers",
			resp:     &interfaces.Response{},
			key:      "Content-Type",
			expected: "",
		},
		{
			name:     "existing header",
			resp:     &interfaces.Response{Headers: headers},
			key:      "Content-Type",
			expected: "application/json",
		},
		{
			name:     "case insensitive header",
			resp:     &interfaces.Response{Headers: headers},
			key:      "content-type",
			expected: "application/json",
		},
		{
			name:     "missing header",
			resp:     &interfaces.Response{Headers: headers},
			key:      "Missing-Header",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetResponseHeader(tt.resp, tt.key)
			if result != tt.expected {
				t.Errorf("GetResponseHeader() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetRateLimitHeaders(t *testing.T) {
	tests := []struct {
		name              string
		resp              *interfaces.Response
		expectedLimit     string
		expectedRemaining string
		expectedReset     string
		expectedRetry     string
	}{
		{
			name:              "nil response",
			resp:              nil,
			expectedLimit:     "",
			expectedRemaining: "",
			expectedReset:     "",
			expectedRetry:     "",
		},
		{
			name: "rate limit headers present",
			resp: &interfaces.Response{
				Headers: http.Header{
					"X-Rate-Limit":           []string{"1000"},
					"X-Rate-Limit-Remaining": []string{"950"},
					"X-Rate-Limit-Reset":     []string{"1640000000"},
					"Retry-After":            []string{"60"},
				},
			},
			expectedLimit:     "1000",
			expectedRemaining: "950",
			expectedReset:     "1640000000",
			expectedRetry:     "60",
		},
		{
			name: "partial rate limit headers",
			resp: &interfaces.Response{
				Headers: http.Header{
					"X-Rate-Limit":           []string{"1000"},
					"X-Rate-Limit-Remaining": []string{"950"},
				},
			},
			expectedLimit:     "1000",
			expectedRemaining: "950",
			expectedReset:     "",
			expectedRetry:     "",
		},
		{
			name: "no rate limit headers",
			resp: &interfaces.Response{
				Headers: make(http.Header),
			},
			expectedLimit:     "",
			expectedRemaining: "",
			expectedReset:     "",
			expectedRetry:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, remaining, reset, retry := GetRateLimitHeaders(tt.resp)
			if limit != tt.expectedLimit {
				t.Errorf("limit = %v, want %v", limit, tt.expectedLimit)
			}
			if remaining != tt.expectedRemaining {
				t.Errorf("remaining = %v, want %v", remaining, tt.expectedRemaining)
			}
			if reset != tt.expectedReset {
				t.Errorf("reset = %v, want %v", reset, tt.expectedReset)
			}
			if retry != tt.expectedRetry {
				t.Errorf("retry = %v, want %v", retry, tt.expectedRetry)
			}
		})
	}
}

func TestGetResponseHeaders(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Custom-Header", "test-value")

	tests := []struct {
		name     string
		resp     *interfaces.Response
		expected int
	}{
		{
			name:     "nil response",
			resp:     nil,
			expected: 0,
		},
		{
			name:     "response with headers",
			resp:     &interfaces.Response{Headers: headers},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetResponseHeaders(tt.resp)
			if len(result) != tt.expected {
				t.Errorf("GetResponseHeaders() returned %d headers, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestResponseStructFields(t *testing.T) {
	// Test that Response struct can hold all expected fields
	resp := &interfaces.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    make(http.Header),
		Body:       []byte("test body"),
		Duration:   100 * time.Millisecond,
		ReceivedAt: time.Now(),
		Size:       9,
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %v, want 200", resp.StatusCode)
	}
	if resp.Status != "200 OK" {
		t.Errorf("Status = %v, want '200 OK'", resp.Status)
	}
	if resp.Duration != 100*time.Millisecond {
		t.Errorf("Duration = %v, want 100ms", resp.Duration)
	}
	if resp.Size != 9 {
		t.Errorf("Size = %v, want 9", resp.Size)
	}
	if string(resp.Body) != "test body" {
		t.Errorf("Body = %v, want 'test body'", string(resp.Body))
	}
}
