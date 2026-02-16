package client

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// HTTP Status Codes
const (
	// Success status codes
	StatusOK      = 200 // Successful GET requests
	StatusCreated = 201 // Successful POST requests

	// Client error status codes
	StatusBadRequest          = 400 // Bad request, invalid arguments
	StatusUnauthorized        = 401 // Authentication required, invalid credentials
	StatusForbidden           = 403 // Forbidden operation
	StatusNotFound            = 404 // Resource not found
	StatusConflict            = 409 // Resource already exists
	StatusUnprocessableEntity = 422 // Validation errors
	StatusTooManyRequests     = 429 // Rate limit exceeded

	// Server error status codes
	StatusInternalServerError = 500 // Server-side error
	StatusBadGateway          = 502 // Gateway error
	StatusServiceUnavailable  = 503 // Service temporarily unavailable
	StatusGatewayTimeout      = 504 // Deadline exceeded
)

// APIError represents an error response from the Nexthink API
type APIError struct {
	Code    string `json:"code,omitempty"`    // Error code if provided
	Message string `json:"message,omitempty"` // Error message
	Details string `json:"details,omitempty"` // Additional error details

	// HTTP response details
	StatusCode int    // HTTP status code
	Status     string // HTTP status text
	Endpoint   string // API endpoint that returned the error
	Method     string // HTTP method used
}

// genericErrorResponse represents a generic API error response wrapper
type genericErrorResponse struct {
	Error   *APIError `json:"error,omitempty"`
	Message string    `json:"message,omitempty"`
	Code    string    `json:"code,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("Nexthink API error (%d %s) [%s] at %s %s: %s",
			e.StatusCode, e.Status, e.Code, e.Method, e.Endpoint, e.Message)
	}
	return fmt.Sprintf("Nexthink API error (%d %s) at %s %s: %s",
		e.StatusCode, e.Status, e.Method, e.Endpoint, e.Message)
}

// ParseErrorResponse parses an error response from the API
func ParseErrorResponse(body []byte, statusCode int, status, method, endpoint string, logger *zap.Logger) error {
	apiError := &APIError{
		StatusCode: statusCode,
		Status:     status,
		Endpoint:   endpoint,
		Method:     method,
	}

	// Try to parse as structured error response
	var errResp genericErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		// Check if error object exists
		if errResp.Error != nil {
			apiError.Code = errResp.Error.Code
			apiError.Message = errResp.Error.Message
			apiError.Details = errResp.Error.Details
		} else {
			// Try top-level message and code
			if errResp.Message != "" {
				apiError.Message = errResp.Message
			}
			if errResp.Code != "" {
				apiError.Code = errResp.Code
			}
		}

		if apiError.Code != "" || apiError.Message != "" {
			logger.Error("API error response",
				zap.Int("status_code", statusCode),
				zap.String("status", status),
				zap.String("method", method),
				zap.String("endpoint", endpoint),
				zap.String("error_code", apiError.Code),
				zap.String("message", apiError.Message))
			return apiError
		}
	}

	// If JSON parsing fails or doesn't match expected format, use raw body as message
	apiError.Message = string(body)
	if apiError.Message == "" {
		apiError.Message = getDefaultErrorMessage(statusCode)
	}

	logger.Error("API error response",
		zap.Int("status_code", statusCode),
		zap.String("status", status),
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.String("message", apiError.Message))

	return apiError
}

// getDefaultErrorMessage returns a default error message based on status code
func getDefaultErrorMessage(statusCode int) string {
	switch statusCode {
	case StatusBadRequest:
		return "The API request is invalid or malformed"
	case StatusUnauthorized:
		return "Authentication required or invalid credentials"
	case StatusForbidden:
		return "You are not allowed to perform the requested operation"
	case StatusNotFound:
		return "The requested resource was not found"
	case StatusConflict:
		return "The resource already exists"
	case StatusUnprocessableEntity:
		return "Validation error"
	case StatusTooManyRequests:
		return "Rate limit exceeded. Too many requests in a given time period"
	case StatusInternalServerError:
		return "Internal server error"
	case StatusBadGateway:
		return "Bad gateway"
	case StatusServiceUnavailable:
		return "Service temporarily unavailable. Retry might work"
	case StatusGatewayTimeout:
		return "The operation took too long to complete"
	default:
		return "Unknown error"
	}
}

// Error type check helpers

// IsBadRequest checks if the error is a bad request error (400)
func IsBadRequest(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusBadRequest
	}
	return false
}

// IsUnauthorized checks if the error is an authentication error (401)
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusUnauthorized
	}
	return false
}

// IsForbidden checks if the error is a forbidden error (403)
func IsForbidden(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusForbidden
	}
	return false
}

// IsNotFound checks if the error is a not found error (404)
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusNotFound
	}
	return false
}

// IsConflict checks if the error is a conflict error (409) - resource already exists
func IsConflict(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusConflict
	}
	return false
}

// IsValidationError checks if the error is a validation/unprocessable entity error (422)
func IsValidationError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusUnprocessableEntity
	}
	return false
}

// IsRateLimited checks if the error is a rate limit error (429)
func IsRateLimited(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusTooManyRequests
	}
	return false
}

// IsServerError checks if the error is a server error (5xx)
func IsServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}

// IsTransient checks if the error is transient and can be retried
func IsTransient(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusServiceUnavailable ||
			apiErr.StatusCode == StatusGatewayTimeout
	}
	return false
}

// GetErrorCode returns the error code from the error
func GetErrorCode(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code
	}
	return ""
}
