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
	m.RegisterExecuteV1Mock()
	m.RegisterExecuteV2Mock()
	m.RegisterStartExportMock()
	m.RegisterGetExportStatusSubmittedMock()
	m.RegisterGetExportStatusCompletedMock()
	m.RegisterGetExportStatusErrorMock()
}

// RegisterErrorMocks registers all error response mocks
func (m *NQLMock) RegisterErrorMocks() {
	m.RegisterUnauthorizedErrorMock()
}

// RegisterExecuteV1Mock registers the mock for ExecuteV1
func (m *NQLMock) RegisterExecuteV1Mock() {
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

// RegisterExecuteV2Mock registers the mock for ExecuteV2
func (m *NQLMock) RegisterExecuteV2Mock() {
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

// RegisterStartExportMock registers the mock for StartExport
func (m *NQLMock) RegisterStartExportMock() {
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

// RegisterGetExportStatusSubmittedMock registers the mock for GetExportStatus (submitted)
func (m *NQLMock) RegisterGetExportStatusSubmittedMock() {
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

// RegisterGetExportStatusCompletedMock registers the mock for GetExportStatus (completed)
func (m *NQLMock) RegisterGetExportStatusCompletedMock() {
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

// RegisterGetExportStatusErrorMock registers the mock for GetExportStatus (error)
func (m *NQLMock) RegisterGetExportStatusErrorMock() {
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
