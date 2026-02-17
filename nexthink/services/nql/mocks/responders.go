package mocks

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarcoal/httpmock"
)

// NQLMock provides mock responses for NQL API endpoints
type NQLMock struct {
	baseURL string
}

// NewNQLMock creates a new NQLMock instance
func NewNQLMock(baseURL string) *NQLMock {
	return &NQLMock{baseURL: baseURL}
}

// RegisterMocks registers all successful response mocks
func (m *NQLMock) RegisterMocks() {
	m.RegisterExecuteNQLV1Mock()
	m.RegisterExecuteNQLV2Mock()
	m.RegisterStartNQLExportMock()
	m.RegisterGetNQLExportStatusSubmittedMock()
	m.RegisterGetNQLExportStatusCompletedMock()
	m.RegisterGetNQLExportStatusErrorMock()
}

// RegisterErrorMocks registers all error response mocks
func (m *NQLMock) RegisterErrorMocks() {
	m.RegisterUnauthorizedErrorMock()
}

// RegisterExecuteNQLV1Mock registers the mock for ExecuteNQLV1
func (m *NQLMock) RegisterExecuteNQLV1Mock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/nql/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("validate_execute_v1_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterExecuteNQLV2Mock registers the mock for ExecuteNQLV2
func (m *NQLMock) RegisterExecuteNQLV2Mock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v2/nql/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("validate_execute_v2_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterStartNQLExportMock registers the mock for StartNQLExport
func (m *NQLMock) RegisterStartNQLExportMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/nql/export",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("validate_start_export.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterGetNQLExportStatusSubmittedMock registers the mock for GetNQLExportStatus (submitted)
func (m *NQLMock) RegisterGetNQLExportStatusSubmittedMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/nql/status/export-123-abc",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("validate_status_submitted.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterGetNQLExportStatusCompletedMock registers the mock for GetNQLExportStatus (completed)
func (m *NQLMock) RegisterGetNQLExportStatusCompletedMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/nql/status/export-456-def",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("validate_status_completed.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterGetNQLExportStatusErrorMock registers the mock for GetNQLExportStatus (error)
func (m *NQLMock) RegisterGetNQLExportStatusErrorMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/nql/status/export-789-ghi",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("validate_status_error.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterUnauthorizedErrorMock registers the mock for unauthorized errors
func (m *NQLMock) RegisterUnauthorizedErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/act/nql/execute/unauthorized",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_unauthorized.json")
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// loadMockData loads a mock JSON file from the mocks directory
func (m *NQLMock) loadMockData(filename string) []byte {
	_, currentFile, _, _ := runtime.Caller(0)
	mockDir := filepath.Dir(currentFile)
	mockPath := filepath.Join(mockDir, filename)

	data, err := os.ReadFile(mockPath)
	if err != nil {
		panic("Failed to load mock data from " + mockPath + ": " + err.Error())
	}
	return data
}
