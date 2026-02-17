package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

type testResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func setupTestClient(t *testing.T, baseURL string) *Transport {
	logger := zaptest.NewLogger(t)

	transport := &Transport{
		client:        resty.New().SetBaseURL(baseURL),
		logger:        logger,
		BaseURL:       baseURL,
		globalHeaders: make(map[string]string),
		userAgent:     "test-agent",
	}

	return transport
}

func TestGet_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify path
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}

		// Verify query params
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected query param limit=10, got %s", r.URL.Query().Get("limit"))
		}

		// Verify headers
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Errorf("Expected header X-Test-Header=test-value, got %s", r.Header.Get("X-Test-Header"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{
			ID:      "123",
			Message: "success",
		})
	}))
	defer server.Close()

	// Create client
	client := setupTestClient(t, server.URL)

	// Execute request
	var result testResponse
	resp, err := client.Get(
		context.Background(),
		"/test",
		map[string]string{"limit": "10"},
		map[string]string{"X-Test-Header": "test-value"},
		&result,
	)

	// Verify response
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Get() response is nil")
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if result.ID != "123" {
		t.Errorf("ID = %q, want %q", result.ID, "123")
	}

	if result.Message != "success" {
		t.Errorf("Message = %q, want %q", result.Message, "success")
	}
}

func TestGet_EmptyQueryParams(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify empty query params are filtered
		if r.URL.Query().Get("empty") != "" {
			t.Error("Empty query param should not be sent")
		}

		// Verify non-empty param is sent
		if r.URL.Query().Get("valid") != "value" {
			t.Errorf("Expected query param valid=value, got %s", r.URL.Query().Get("valid"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	_, err := client.Get(
		context.Background(),
		"/test",
		map[string]string{"empty": "", "valid": "value"},
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
}

func TestPost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Decode body
		var received testResponse
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}

		if received.Message != "test message" {
			t.Errorf("Expected message 'test message', got %q", received.Message)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(testResponse{
			ID:      "456",
			Message: "created",
		})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	requestBody := testResponse{Message: "test message"}
	var result testResponse
	resp, err := client.Post(
		context.Background(),
		"/test",
		requestBody,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("Post() error = %v, want nil", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("StatusCode = %d, want 201", resp.StatusCode)
	}

	if result.ID != "456" {
		t.Errorf("ID = %q, want %q", result.ID, "456")
	}
}

func TestPost_NilBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no body sent
		if r.ContentLength > 0 {
			t.Error("Expected no body for nil body parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	_, err := client.Post(
		context.Background(),
		"/test",
		nil,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("Post() error = %v, want nil", err)
	}
}

func TestPostWithQuery_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify query params
		if r.URL.Query().Get("action") != "create" {
			t.Errorf("Expected query param action=create, got %s", r.URL.Query().Get("action"))
		}

		// Verify body
		var received testResponse
		json.NewDecoder(r.Body).Decode(&received)
		if received.Message != "test" {
			t.Errorf("Expected body message 'test', got %q", received.Message)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "789"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	_, err := client.PostWithQuery(
		context.Background(),
		"/test",
		map[string]string{"action": "create"},
		testResponse{Message: "test"},
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("PostWithQuery() error = %v, want nil", err)
	}

	if result.ID != "789" {
		t.Errorf("ID = %q, want %q", result.ID, "789")
	}
}

func TestPut_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "updated"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := client.Put(
		context.Background(),
		"/test/123",
		testResponse{Message: "update"},
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("Put() error = %v, want nil", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if result.ID != "updated" {
		t.Errorf("ID = %q, want %q", result.ID, "updated")
	}
}

func TestPatch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "patched"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := client.Patch(
		context.Background(),
		"/test/123",
		map[string]string{"field": "value"},
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("Patch() error = %v, want nil", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if result.ID != "patched" {
		t.Errorf("ID = %q, want %q", result.ID, "patched")
	}
}

func TestDelete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		// Verify query params
		if r.URL.Query().Get("confirm") != "true" {
			t.Errorf("Expected query param confirm=true, got %s", r.URL.Query().Get("confirm"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{Message: "deleted"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := client.Delete(
		context.Background(),
		"/test/123",
		map[string]string{"confirm": "true"},
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if result.Message != "deleted" {
		t.Errorf("Message = %q, want %q", result.Message, "deleted")
	}
}

func TestDeleteWithBody_Success(t *testing.T) {
	// Test that DeleteWithBody method works correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{Message: "bulk deleted"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	deleteBody := map[string][]string{
		"ids": {"123", "456"},
	}
	var result testResponse
	resp, err := client.DeleteWithBody(
		context.Background(),
		"/test/bulk",
		deleteBody,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("DeleteWithBody() error = %v, want nil", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if result.Message != "bulk deleted" {
		t.Errorf("Message = %q, want %q", result.Message, "bulk deleted")
	}
}

func TestDeleteWithBody_NilBody(t *testing.T) {
	// Test DeleteWithBody with nil body
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{Message: "deleted"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := client.DeleteWithBody(
		context.Background(),
		"/test/delete",
		nil,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("DeleteWithBody() error = %v, want nil", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
}

func TestPostForm_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify Content-Type is form-urlencoded
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/x-www-form-urlencoded") {
			t.Errorf("Expected Content-Type to contain application/x-www-form-urlencoded, got %s", contentType)
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}

		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}

		if r.FormValue("scope") != "service:integration" {
			t.Errorf("Expected scope=service:integration, got %s", r.FormValue("scope"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{
			ID:      "form-123",
			Message: "form submitted",
		})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	formData := map[string]string{
		"grant_type": "client_credentials",
		"scope":      "service:integration",
	}

	var result testResponse
	resp, err := client.PostForm(
		context.Background(),
		"/oauth2/token",
		formData,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("PostForm() error = %v, want nil", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if result.ID != "form-123" {
		t.Errorf("ID = %q, want %q", result.ID, "form-123")
	}

	if result.Message != "form submitted" {
		t.Errorf("Message = %q, want %q", result.Message, "form submitted")
	}
}

func TestPostForm_NilFormData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	_, err := client.PostForm(
		context.Background(),
		"/test",
		nil,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("PostForm() error = %v, want nil", err)
	}
}

func TestPostMultipart_Success(t *testing.T) {
	expectedFileName := "test.txt"
	expectedFileContent := "test file content"
	expectedFieldName := "uploadedFile"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify Content-Type is multipart
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			t.Errorf("Expected Content-Type to contain multipart/form-data, got %s", contentType)
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}

		// Verify form field
		if r.FormValue("description") != "test description" {
			t.Errorf("Expected description='test description', got %s", r.FormValue("description"))
		}

		// Verify file
		file, header, err := r.FormFile(expectedFieldName)
		if err != nil {
			t.Errorf("Failed to get file: %v", err)
		}
		defer file.Close()

		if header.Filename != expectedFileName {
			t.Errorf("Expected filename %q, got %q", expectedFileName, header.Filename)
		}

		fileContent, _ := io.ReadAll(file)
		if string(fileContent) != expectedFileContent {
			t.Errorf("Expected file content %q, got %q", expectedFileContent, string(fileContent))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(testResponse{
			ID:      "upload-123",
			Message: "file uploaded",
		})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	fileContent := strings.NewReader(expectedFileContent)
	formFields := map[string]string{
		"description": "test description",
	}

	progressCalled := false
	progressCallback := func(fieldName, fileName string, bytesWritten, totalBytes int64) {
		progressCalled = true
		if fieldName != expectedFieldName {
			t.Errorf("Progress callback fieldName = %q, want %q", fieldName, expectedFieldName)
		}
		if fileName != expectedFileName {
			t.Errorf("Progress callback fileName = %q, want %q", fileName, expectedFileName)
		}
	}

	var result testResponse
	resp, err := client.PostMultipart(
		context.Background(),
		"/test/upload",
		expectedFieldName,
		expectedFileName,
		fileContent,
		int64(len(expectedFileContent)),
		formFields,
		nil,
		progressCallback,
		&result,
	)

	if err != nil {
		t.Fatalf("PostMultipart() error = %v, want nil", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("StatusCode = %d, want 201", resp.StatusCode)
	}

	if result.ID != "upload-123" {
		t.Errorf("ID = %q, want %q", result.ID, "upload-123")
	}

	if result.Message != "file uploaded" {
		t.Errorf("Message = %q, want %q", result.Message, "file uploaded")
	}

	if !progressCalled {
		t.Error("Progress callback was not called")
	}
}

func TestPostMultipart_NilCallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	fileContent := strings.NewReader("test")
	var result testResponse
	_, err := client.PostMultipart(
		context.Background(),
		"/test/upload",
		"file",
		"test.txt",
		fileContent,
		4,
		nil,
		nil,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("PostMultipart() error = %v, want nil", err)
	}
}

func TestPostMultipart_EmptyFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	_, err := client.PostMultipart(
		context.Background(),
		"/test/upload",
		"",
		"",
		nil,
		0,
		nil,
		nil,
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("PostMultipart() error = %v, want nil", err)
	}
}

func TestGetBytes_Success(t *testing.T) {
	csvData := "id,name\n1,test1\n2,test2"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify query params
		if r.URL.Query().Get("format") != "csv" {
			t.Errorf("Expected query param format=csv, got %s", r.URL.Query().Get("format"))
		}

		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	resp, data, err := client.GetBytes(
		context.Background(),
		"/test/export",
		map[string]string{"format": "csv"},
		nil,
	)

	if err != nil {
		t.Fatalf("GetBytes() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("GetBytes() response is nil")
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if string(data) != csvData {
		t.Errorf("CSV data = %q, want %q", string(data), csvData)
	}
}

func TestGetBytes_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Resource not found",
		})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	resp, data, err := client.GetBytes(
		context.Background(),
		"/test/not-found",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("GetBytes() error = nil, want error")
	}

	if resp == nil {
		t.Fatal("GetBytes() response is nil, should return metadata even on error")
	}

	if data != nil {
		t.Errorf("GetBytes() data should be nil on error, got %v", data)
	}

	// Verify it's an APIError
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("Expected *APIError, got %T", err)
	}

	if apiErr != nil && apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}

func TestExecuteRequest_UnsupportedMethod(t *testing.T) {
	client := setupTestClient(t, "http://localhost")

	req := client.client.R()
	_, err := client.executeRequest(req, "UNSUPPORTED", "/test")

	if err == nil {
		t.Fatal("executeRequest() error = nil, want error for unsupported method")
	}

	if err.Error() != "unsupported HTTP method: UNSUPPORTED" {
		t.Errorf("Error message = %q, want %q", err.Error(), "unsupported HTTP method: UNSUPPORTED")
	}
}

func TestRequest_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Invalid request",
			"code":    "INVALID_QUERY",
		})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := client.Get(
		context.Background(),
		"/test",
		nil,
		nil,
		&result,
	)

	if err == nil {
		t.Fatal("Get() error = nil, want error")
	}

	// Response metadata should still be available
	if resp == nil {
		t.Fatal("Get() response is nil, should return metadata even on error")
	}

	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}

	// Verify error is APIError
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.Message != "Invalid request" {
		t.Errorf("Error message = %q, want %q", apiErr.Message, "Invalid request")
	}

	if apiErr.Code != "INVALID_QUERY" {
		t.Errorf("Error code = %q, want %q", apiErr.Code, "INVALID_QUERY")
	}
}

func TestRequest_WithGlobalHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify global header
		if r.Header.Get("X-Global") != "global-value" {
			t.Errorf("Expected global header X-Global=global-value, got %s", r.Header.Get("X-Global"))
		}

		// Verify request header overrides
		if r.Header.Get("X-Override") != "request-value" {
			t.Errorf("Expected overridden header X-Override=request-value, got %s", r.Header.Get("X-Override"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)
	client.globalHeaders["X-Global"] = "global-value"
	client.globalHeaders["X-Override"] = "global-override"

	var result testResponse
	_, err := client.Get(
		context.Background(),
		"/test",
		nil,
		map[string]string{"X-Override": "request-value"},
		&result,
	)

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
}

func TestRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should never be reached due to context cancellation
		t.Error("Request should not reach server due to cancelled context")
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result testResponse
	_, err := client.Get(ctx, "/test", nil, nil, &result)

	if err == nil {
		t.Fatal("Get() error = nil, want error for cancelled context")
	}
}

func TestRequest_InvalidContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return HTML instead of JSON
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Not JSON</body></html>"))
	}))
	defer server.Close()

	client := setupTestClient(t, server.URL)

	var result testResponse
	_, err := client.Get(
		context.Background(),
		"/test",
		nil,
		nil,
		&result,
	)

	if err == nil {
		t.Fatal("Get() error = nil, want error for invalid content type")
	}

	// Should contain content-type error message
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}
