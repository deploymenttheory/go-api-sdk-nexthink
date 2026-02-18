package client

import (
	"context"
	"fmt"
	"io"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Get executes a GET request
func (t *Transport) Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "GET", path)
}

// Post executes a POST request with JSON body
func (t *Transport) Post(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "POST", path)
}

// PostWithQuery executes a POST request with both body and query parameters
func (t *Transport) PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	if body != nil {
		req.SetBody(body)
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "POST", path)
}

// PostForm executes a POST request with form-urlencoded data
func (t *Transport) PostForm(ctx context.Context, path string, formData map[string]string, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	if formData != nil {
		req.SetFormData(formData)
	}

	// Apply headers with precedence (global first, then per-request)
	// Note: Content-Type is handled automatically by resty for form data
	for k, v := range t.globalHeaders {
		if v != "" && k != "Content-Type" {
			req.SetHeader(k, v)
		}
	}
	for k, v := range headers {
		if v != "" && k != "Content-Type" {
			req.SetHeader(k, v)
		}
	}

	return t.executeRequest(req, "POST", path)
}

// PostMultipart executes a POST request with multipart form data and progress tracking
func (t *Transport) PostMultipart(ctx context.Context, path string, fileField string, fileName string, fileReader io.Reader, fileSize int64, formFields map[string]string, headers map[string]string, progressCallback interfaces.MultipartProgressCallback, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	if fileReader != nil && fileName != "" && fileField != "" {
		multipartField := &resty.MultipartField{
			Name:     fileField,
			FileName: fileName,
			Reader:   fileReader,
			FileSize: fileSize,
		}

		if progressCallback != nil {
			multipartField.ProgressCallback = func(progress resty.MultipartFieldProgress) {
				progressCallback(progress.Name, progress.FileName, progress.Written, progress.FileSize)
			}
		}

		req.SetMultipartFields(multipartField)
	}

	if len(formFields) > 0 {
		req.SetMultipartFormData(formFields)
	}

	// Apply headers with precedence (global first, then per-request)
	// Note: Content-Type is handled automatically by resty for multipart
	for k, v := range t.globalHeaders {
		if v != "" && k != "Content-Type" {
			req.SetHeader(k, v)
		}
	}
	for k, v := range headers {
		if v != "" && k != "Content-Type" {
			req.SetHeader(k, v)
		}
	}

	return t.executeRequest(req, "POST", path)
}

// Put executes a PUT request
func (t *Transport) Put(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "PUT", path)
}

// Patch executes a PATCH request
func (t *Transport) Patch(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "PATCH", path)
}

// Delete executes a DELETE request
func (t *Transport) Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "DELETE", path)
}

// DeleteWithBody executes a DELETE request with body (for bulk operations)
func (t *Transport) DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := t.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	t.applyHeaders(req, headers)

	return t.executeRequest(req, "DELETE", path)
}

// GetBytes performs a GET request and returns raw bytes without unmarshaling
// Use this for non-JSON responses like HTML, CSV, binary files, etc.
func (t *Transport) GetBytes(ctx context.Context, path string, queryParams map[string]string, headers map[string]string) (*interfaces.Response, []byte, error) {
	req := t.client.R().
		SetContext(ctx)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	t.applyHeaders(req, headers)

	t.logger.Debug("Executing bytes request",
		zap.String("method", "GET"),
		zap.String("path", path))

	resp, err := req.Get(path)
	clientResp := toInterfaceResponse(resp)
	if err != nil {
		t.logger.Error("Bytes request failed",
			zap.String("path", path),
			zap.Error(err))
		return clientResp, nil, fmt.Errorf("bytes request failed: %w", err)
	}

	if IsResponseError(clientResp) {
		return clientResp, nil, ParseErrorResponse(
			[]byte(resp.String()),
			resp.StatusCode(),
			resp.Status(),
			"GET",
			path,
			t.logger,
		)
	}

	body := []byte(resp.String())
	t.logger.Debug("Bytes request completed successfully",
		zap.String("path", path),
		zap.Int("status_code", resp.StatusCode()),
		zap.Int("content_length", len(body)))

	return clientResp, body, nil
}

// executeRequest is a centralized request executor that handles error processing
// Returns response metadata and error. Response is always non-nil for accessing headers.
func (t *Transport) executeRequest(req *resty.Request, method, path string) (*interfaces.Response, error) {
	t.logger.Debug("Executing API request",
		zap.String("method", method),
		zap.String("path", path))

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "PATCH":
		resp, err = req.Patch(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return toInterfaceResponse(nil), fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// Convert to interface response (always return response metadata)
	clientResp := toInterfaceResponse(resp)

	if err != nil {
		t.logger.Error("Request failed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Error(err))
		return clientResp, fmt.Errorf("request failed: %w", err)
	}

	if err := t.validateResponse(resp, method, path); err != nil {
		return clientResp, err
	}

	if IsResponseError(clientResp) {
		return clientResp, ParseErrorResponse(
			[]byte(resp.String()),
			resp.StatusCode(),
			resp.Status(),
			method,
			path,
			t.logger,
		)
	}

	t.logger.Debug("Request completed successfully",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", resp.StatusCode()))

	return clientResp, nil
}
