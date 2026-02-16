package interfaces

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Response represents HTTP response metadata that can be returned alongside errors
// This allows callers to access response headers (rate limits, retry-after, etc.) even on errors
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
	Duration   time.Duration // Time taken for the request
	ReceivedAt time.Time     // When the response was received
	Size       int64         // Response body size in bytes
}

// MultipartProgressCallback is a callback function for multipart upload progress
type MultipartProgressCallback func(fieldName string, fileName string, bytesWritten int64, totalBytes int64)

// HTTPClient interface that services will use
// This breaks import cycles by providing a contract without implementation
type HTTPClient interface {
	// Get executes a GET request and unmarshals the JSON response into the result parameter.
	// Query parameters and headers are applied if provided.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Get(
		ctx context.Context, // request context
		path string, // API endpoint path
		queryParams map[string]string, // URL query parameters
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// Post executes a POST request with a JSON body.
	// The body is marshaled to JSON and the response is unmarshaled into the result parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Post(
		ctx context.Context, // request context
		path string, // API endpoint path
		body any, // request body to marshal as JSON
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// PostWithQuery executes a POST request with both query parameters and a JSON body.
	// The body is marshaled to JSON and the response is unmarshaled into the result parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	PostWithQuery(
		ctx context.Context, // request context
		path string, // API endpoint path
		queryParams map[string]string, // URL query parameters
		body any, // request body to marshal as JSON
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// Put executes a PUT request with a JSON body.
	// The body is marshaled to JSON and the response is unmarshaled into the result parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Put(
		ctx context.Context, // request context
		path string, // API endpoint path
		body any, // request body to marshal as JSON
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// Patch executes a PATCH request with a JSON body.
	// The body is marshaled to JSON and the response is unmarshaled into the result parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Patch(
		ctx context.Context, // request context
		path string, // API endpoint path
		body any, // request body to marshal as JSON
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// Delete executes a DELETE request and unmarshals the JSON response into the result parameter.
	// Query parameters and headers are applied if provided.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Delete(
		ctx context.Context, // request context
		path string, // API endpoint path
		queryParams map[string]string, // URL query parameters
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// DeleteWithBody executes a DELETE request with a JSON body (for bulk operations).
	// The body is marshaled to JSON and the response is unmarshaled into the result parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	DeleteWithBody(
		ctx context.Context, // request context
		path string, // API endpoint path
		body any, // request body to marshal as JSON
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// PostForm executes a POST request with form-urlencoded data.
	// The Content-Type header is automatically set to application/x-www-form-urlencoded.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	PostForm(
		ctx context.Context, // request context
		path string, // API endpoint path
		formData map[string]string, // form fields as key-value pairs
		headers map[string]string, // HTTP headers
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// PostMultipart executes a POST request with multipart/form-data encoding, typically for file uploads.
	// The Content-Type header is automatically set to multipart/form-data with a boundary.
	// Progress tracking is supported via the optional progressCallback parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	PostMultipart(
		ctx context.Context, // request context
		path string, // API endpoint path
		fileField string, // form field name for the file
		fileName string, // name of the file being uploaded
		fileReader io.Reader, // reader for file content
		fileSize int64, // size of the file in bytes
		formFields map[string]string, // additional form fields
		headers map[string]string, // HTTP headers
		progressCallback MultipartProgressCallback, // optional progress callback
		result any, // pointer to unmarshal response into
	) (*Response, error)

	// GetBytes performs a GET request and returns raw bytes without unmarshaling.
	// Use this for non-JSON responses like CSV, HTML, binary files, etc.
	GetBytes(
		ctx context.Context, // request context
		path string, // API endpoint path
		queryParams map[string]string, // URL query parameters
		headers map[string]string, // HTTP headers
	) (*Response, []byte, error)

	// GetLogger returns the configured zap logger instance.
	GetLogger() *zap.Logger
}

// ServiceQueryBuilder defines the query builder contract for services.
// Provides a fluent interface for constructing URL query parameters.
type ServiceQueryBuilder interface {
	// AddString adds a string parameter if the value is not empty.
	// Returns the builder for method chaining.
	AddString(key, value string) QueryBuilder

	// AddInt adds an integer parameter if the value is greater than 0.
	// The integer is converted to a string representation.
	// Returns the builder for method chaining.
	AddInt(key string, value int) QueryBuilder

	// AddInt64 adds an int64 parameter if the value is greater than 0.
	// The int64 is converted to a string representation.
	// Returns the builder for method chaining.
	AddInt64(key string, value int64) QueryBuilder

	// AddBool adds a boolean parameter.
	// The boolean is converted to "true" or "false" string representation.
	// Returns the builder for method chaining.
	AddBool(key string, value bool) QueryBuilder

	// AddTime adds a time parameter in RFC3339 format if the time is not zero.
	// Returns the builder for method chaining.
	AddTime(key string, value time.Time) QueryBuilder

	// AddStringSlice adds a string slice parameter as comma-separated values.
	// Empty string values within the slice are skipped.
	// Returns the builder for method chaining.
	AddStringSlice(key string, values []string) QueryBuilder

	// AddIntSlice adds an integer slice parameter as comma-separated values.
	// Returns the builder for method chaining.
	AddIntSlice(key string, values []int) QueryBuilder

	// AddCustom adds a custom parameter with any value without validation.
	// Use this when you need to add a parameter regardless of its value.
	// Returns the builder for method chaining.
	AddCustom(key, value string) QueryBuilder

	// AddIfNotEmpty adds a parameter only if the value is not empty.
	// Functionally equivalent to AddString.
	// Returns the builder for method chaining.
	AddIfNotEmpty(key, value string) QueryBuilder

	// AddIfTrue adds a parameter only if the condition is true.
	// Returns the builder for method chaining.
	AddIfTrue(condition bool, key, value string) QueryBuilder

	// Merge copies all parameters from another map into this builder.
	// Existing parameters with the same keys will be overwritten.
	// Returns the builder for method chaining.
	Merge(other map[string]string) QueryBuilder

	// Remove deletes a parameter from the builder.
	// Returns the builder for method chaining.
	Remove(key string) QueryBuilder

	// Has checks if a parameter exists in the builder.
	Has(key string) bool

	// Get retrieves the value of a parameter.
	// Returns an empty string if the parameter does not exist.
	Get(key string) string

	// Build returns a copy of the query parameters as a map.
	// The returned map is a copy to prevent external modification.
	Build() map[string]string

	// BuildString returns the query parameters as a URL-encoded string.
	// Parameters are joined with "&" separators in key=value format.
	// Returns an empty string if no parameters are set.
	BuildString() string

	// Clear removes all parameters from the builder.
	// Returns the builder for method chaining.
	Clear() QueryBuilder

	// Count returns the number of parameters currently in the builder.
	Count() int

	// IsEmpty returns true if no parameters are set in the builder.
	IsEmpty() bool
}

// QueryBuilder interface for method chaining.
// Embeds ServiceQueryBuilder to provide the same functionality.
type QueryBuilder interface {
	ServiceQueryBuilder
}
