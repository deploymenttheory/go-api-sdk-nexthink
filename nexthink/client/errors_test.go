package client

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name        string
		apiError    *APIError
		wantContain []string
	}{
		{
			name: "error with code",
			apiError: &APIError{
				Code:       "INVALID_QUERY",
				Message:    "Query ID is malformed",
				StatusCode: 400,
				Status:     "400 Bad Request",
				Method:     "POST",
				Endpoint:   "/api/v1/nql/execute",
			},
			wantContain: []string{
				"Nexthink API error",
				"400",
				"Bad Request",
				"INVALID_QUERY",
				"POST",
				"/api/v1/nql/execute",
				"Query ID is malformed",
			},
		},
		{
			name: "error without code",
			apiError: &APIError{
				Message:    "Resource not found",
				StatusCode: 404,
				Status:     "404 Not Found",
				Method:     "GET",
				Endpoint:   "/api/v1/resource",
			},
			wantContain: []string{
				"Nexthink API error",
				"404",
				"Not Found",
				"GET",
				"/api/v1/resource",
				"Resource not found",
			},
		},
		{
			name: "error with details",
			apiError: &APIError{
				Code:       "VALIDATION_ERROR",
				Message:    "Validation failed",
				Details:    "Field 'query_id' must start with #",
				StatusCode: 422,
				Status:     "422 Unprocessable Entity",
				Method:     "POST",
				Endpoint:   "/api/v1/nql/execute",
			},
			wantContain: []string{
				"Nexthink API error",
				"422",
				"VALIDATION_ERROR",
				"Validation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiError.Error()
			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("Error() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

func TestParseErrorResponse_StructuredJSON(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name           string
		body           string
		statusCode     int
		status         string
		method         string
		endpoint       string
		wantCode       string
		wantMessage    string
		wantStatusCode int
	}{
		{
			name: "nested error object",
			body: `{
				"error": {
					"code": "INVALID_QUERY",
					"message": "The query ID is invalid",
					"details": "Query ID must start with #"
				}
			}`,
			statusCode:     400,
			status:         "400 Bad Request",
			method:         "POST",
			endpoint:       "/api/v1/nql/execute",
			wantCode:       "INVALID_QUERY",
			wantMessage:    "The query ID is invalid",
			wantStatusCode: 400,
		},
		{
			name: "top-level error fields",
			body: `{
				"message": "Authentication failed",
				"code": "INVALID_CREDENTIALS"
			}`,
			statusCode:     401,
			status:         "401 Unauthorized",
			method:         "GET",
			endpoint:       "/api/v1/data",
			wantCode:       "INVALID_CREDENTIALS",
			wantMessage:    "Authentication failed",
			wantStatusCode: 401,
		},
		{
			name: "message only",
			body: `{
				"message": "Resource not found"
			}`,
			statusCode:     404,
			status:         "404 Not Found",
			method:         "GET",
			endpoint:       "/api/v1/resource/123",
			wantMessage:    "Resource not found",
			wantStatusCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseErrorResponse(
				[]byte(tt.body),
				tt.statusCode,
				tt.status,
				tt.method,
				tt.endpoint,
				logger,
			)

			if err == nil {
				t.Fatal("ParseErrorResponse() returned nil error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("ParseErrorResponse() returned %T, want *APIError", err)
			}

			if apiErr.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", apiErr.Code, tt.wantCode)
			}

			if apiErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}

			if apiErr.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestParseErrorResponse_InvalidJSON(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name           string
		body           string
		statusCode     int
		wantMessage    string
		wantStatusCode int
	}{
		{
			name:           "plain text error",
			body:           "Something went wrong",
			statusCode:     500,
			wantMessage:    "Something went wrong",
			wantStatusCode: 500,
		},
		{
			name:           "HTML error page",
			body:           "<html><body>Error 404</body></html>",
			statusCode:     404,
			wantMessage:    "<html><body>Error 404</body></html>",
			wantStatusCode: 404,
		},
		{
			name:           "empty body uses default message",
			body:           "",
			statusCode:     503,
			wantMessage:    "Service temporarily unavailable. Retry might work",
			wantStatusCode: 503,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseErrorResponse(
				[]byte(tt.body),
				tt.statusCode,
				"",
				"GET",
				"/test",
				logger,
			)

			if err == nil {
				t.Fatal("ParseErrorResponse() returned nil error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("ParseErrorResponse() returned %T, want *APIError", err)
			}

			if apiErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}

			if apiErr.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestGetDefaultErrorMessage(t *testing.T) {
	tests := []struct {
		statusCode int
		want       string
	}{
		{StatusBadRequest, "The API request is invalid or malformed"},
		{StatusUnauthorized, "Authentication required or invalid credentials"},
		{StatusForbidden, "You are not allowed to perform the requested operation"},
		{StatusNotFound, "The requested resource was not found"},
		{StatusConflict, "The resource already exists"},
		{StatusUnprocessableEntity, "Validation error"},
		{StatusTooManyRequests, "Rate limit exceeded. Too many requests in a given time period"},
		{StatusInternalServerError, "Internal server error"},
		{StatusBadGateway, "Bad gateway"},
		{StatusServiceUnavailable, "Service temporarily unavailable. Retry might work"},
		{StatusGatewayTimeout, "The operation took too long to complete"},
		{999, "Unknown error"},
		{418, "Unknown error"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := getDefaultErrorMessage(tt.statusCode)
			if got != tt.want {
				t.Errorf("getDefaultErrorMessage(%d) = %q, want %q", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestIsBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: true,
		},
		{
			name: "404 error",
			err:  &APIError{StatusCode: 404},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBadRequest(tt.err)
			if got != tt.want {
				t.Errorf("IsBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 error",
			err:  &APIError{StatusCode: 401},
			want: true,
		},
		{
			name: "403 error",
			err:  &APIError{StatusCode: 403},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUnauthorized(tt.err)
			if got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "403 error",
			err:  &APIError{StatusCode: 403},
			want: true,
		},
		{
			name: "401 error",
			err:  &APIError{StatusCode: 401},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsForbidden(tt.err)
			if got != tt.want {
				t.Errorf("IsForbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 error",
			err:  &APIError{StatusCode: 404},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFound(tt.err)
			if got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsConflict(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "409 error",
			err:  &APIError{StatusCode: 409},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsConflict(tt.err)
			if got != tt.want {
				t.Errorf("IsConflict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "422 error",
			err:  &APIError{StatusCode: 422},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidationError(tt.err)
			if got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "429 error",
			err:  &APIError{StatusCode: 429},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRateLimited(tt.err)
			if got != tt.want {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "500 error",
			err:  &APIError{StatusCode: 500},
			want: true,
		},
		{
			name: "502 error",
			err:  &APIError{StatusCode: 502},
			want: true,
		},
		{
			name: "503 error",
			err:  &APIError{StatusCode: 503},
			want: true,
		},
		{
			name: "599 error",
			err:  &APIError{StatusCode: 599},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "600 error (out of 5xx range)",
			err:  &APIError{StatusCode: 600},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsServerError(tt.err)
			if got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTransient(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "503 error",
			err:  &APIError{StatusCode: 503},
			want: true,
		},
		{
			name: "504 error",
			err:  &APIError{StatusCode: 504},
			want: true,
		},
		{
			name: "500 error",
			err:  &APIError{StatusCode: 500},
			want: false,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTransient(tt.err)
			if got != tt.want {
				t.Errorf("IsTransient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "APIError with code",
			err:  &APIError{Code: "INVALID_QUERY"},
			want: "INVALID_QUERY",
		},
		{
			name: "APIError without code",
			err:  &APIError{Message: "Some error"},
			want: "",
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: "",
		},
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetErrorCode(tt.err)
			if got != tt.want {
				t.Errorf("GetErrorCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError_JSON_Marshalling(t *testing.T) {
	// Test that APIError can be marshalled and unmarshalled as JSON
	original := &APIError{
		Code:       "VALIDATION_ERROR",
		Message:    "Test error message",
		Details:    "Additional details",
		StatusCode: 422,
		Status:     "422 Unprocessable Entity",
		Endpoint:   "/api/v1/test",
		Method:     "POST",
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal APIError: %v", err)
	}

	// Unmarshal back
	var unmarshalled APIError
	err = json.Unmarshal(data, &unmarshalled)
	if err != nil {
		t.Fatalf("Failed to unmarshal APIError: %v", err)
	}

	// Compare (only exported fields will be in JSON)
	if unmarshalled.Code != original.Code {
		t.Errorf("Code = %q, want %q", unmarshalled.Code, original.Code)
	}
	if unmarshalled.Message != original.Message {
		t.Errorf("Message = %q, want %q", unmarshalled.Message, original.Message)
	}
	if unmarshalled.Details != original.Details {
		t.Errorf("Details = %q, want %q", unmarshalled.Details, original.Details)
	}
}

func TestErrorConstants(t *testing.T) {
	// Verify error constants have expected values
	constants := map[string]int{
		"StatusOK":                   StatusOK,
		"StatusCreated":              StatusCreated,
		"StatusBadRequest":           StatusBadRequest,
		"StatusUnauthorized":         StatusUnauthorized,
		"StatusForbidden":            StatusForbidden,
		"StatusNotFound":             StatusNotFound,
		"StatusConflict":             StatusConflict,
		"StatusUnprocessableEntity":  StatusUnprocessableEntity,
		"StatusTooManyRequests":      StatusTooManyRequests,
		"StatusInternalServerError":  StatusInternalServerError,
		"StatusBadGateway":           StatusBadGateway,
		"StatusServiceUnavailable":   StatusServiceUnavailable,
		"StatusGatewayTimeout":       StatusGatewayTimeout,
	}

	expected := map[string]int{
		"StatusOK":                   200,
		"StatusCreated":              201,
		"StatusBadRequest":           400,
		"StatusUnauthorized":         401,
		"StatusForbidden":            403,
		"StatusNotFound":             404,
		"StatusConflict":             409,
		"StatusUnprocessableEntity":  422,
		"StatusTooManyRequests":      429,
		"StatusInternalServerError":  500,
		"StatusBadGateway":           502,
		"StatusServiceUnavailable":   503,
		"StatusGatewayTimeout":       504,
	}

	for name, got := range constants {
		want := expected[name]
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}
}
