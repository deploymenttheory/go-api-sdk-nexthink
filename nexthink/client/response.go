package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// toInterfaceResponse converts a resty.Response to interfaces.Response
func toInterfaceResponse(resp *resty.Response) *interfaces.Response {
	if resp == nil {
		return &interfaces.Response{
			Headers: make(http.Header),
		}
	}

	return &interfaces.Response{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Headers:    resp.Header(),
		Body:       []byte(resp.String()),
		Duration:   resp.Duration(),
		ReceivedAt: resp.ReceivedAt(),
		Size:       resp.Size(),
	}
}

// Response helper functions for working with interfaces.Response

// IsResponseSuccess returns true if the response status code is 2xx
func IsResponseSuccess(resp *interfaces.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// IsResponseError returns true if the response status code is 4xx or 5xx
func IsResponseError(resp *interfaces.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 400
}

// GetResponseHeader returns a header value from the response by key (case-insensitive)
func GetResponseHeader(resp *interfaces.Response, key string) string {
	if resp == nil || resp.Headers == nil {
		return ""
	}
	return resp.Headers.Get(key)
}

// GetResponseHeaders returns all headers from the response
func GetResponseHeaders(resp *interfaces.Response) http.Header {
	if resp == nil {
		return make(http.Header)
	}
	return resp.Headers
}

// GetRateLimitHeaders extracts common rate limit headers from the response
// Returns: (limit, remaining, reset, retry-after)
func GetRateLimitHeaders(resp *interfaces.Response) (limit, remaining, reset, retryAfter string) {
	if resp == nil {
		return
	}
	return resp.Headers.Get("X-Rate-Limit"),
		resp.Headers.Get("X-Rate-Limit-Remaining"),
		resp.Headers.Get("X-Rate-Limit-Reset"),
		resp.Headers.Get("Retry-After")
}

// validateResponse validates the HTTP response before processing
// This includes checking for empty responses and validating Content-Type for JSON endpoints
func (t *Transport) validateResponse(resp *resty.Response, method, path string) error {
	// Handle empty responses (204 No Content, etc.)
	bodyLen := len(resp.String())
	if resp.Header().Get("Content-Length") == "0" || bodyLen == 0 {
		t.logger.Debug("Empty response received",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status_code", resp.StatusCode()))
		return nil
	}

	// For non-error responses with content, validate Content-Type is JSON
	// Skip validation for:
	// - Error responses (handled by error parser)
	// - Endpoints that explicitly return non-JSON (download endpoints, etc.)
	if !resp.IsError() && bodyLen > 0 {
		contentType := resp.Header().Get("Content-Type")

		// Allow responses without Content-Type header (some endpoints don't set it)
		if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
			t.logger.Warn("Unexpected Content-Type in response",
				zap.String("method", method),
				zap.String("path", path),
				zap.String("content_type", contentType),
				zap.String("expected", "application/json"))

			return fmt.Errorf("unexpected response Content-Type from %s %s: got %q, expected application/json",
				method, path, contentType)
		}
	}

	return nil
}
