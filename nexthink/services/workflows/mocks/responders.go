package mocks

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarcoal/httpmock"
)

// WorkflowsMock provides mock responses for Workflows API endpoints
type WorkflowsMock struct {
	baseURL string
}

// NewWorkflowsMock creates a new WorkflowsMock instance
func NewWorkflowsMock(baseURL string) *WorkflowsMock {
	return &WorkflowsMock{baseURL: baseURL}
}

// RegisterMocks registers all successful response mocks
func (m *WorkflowsMock) RegisterMocks() {
	m.RegisterExecuteV1Mock()
	m.RegisterExecuteV2Mock()
	m.RegisterListWorkflowsMock()
	m.RegisterGetWorkflowDetailsMock()
	m.RegisterTriggerThinkletMock()
}

// RegisterErrorMocks registers all error response mocks
func (m *WorkflowsMock) RegisterErrorMocks() {
	m.RegisterUnauthorizedErrorMock()
	m.RegisterValidationErrorMock()
}

// RegisterExecuteV1Mock registers the mock for ExecuteV1
func (m *WorkflowsMock) RegisterExecuteV1Mock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/workflows/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("execute_v1_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterExecuteV2Mock registers the mock for ExecuteV2
func (m *WorkflowsMock) RegisterExecuteV2Mock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v2/workflows/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("execute_v2_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterListWorkflowsMock registers the mock for ListWorkflows
func (m *WorkflowsMock) RegisterListWorkflowsMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/workflows",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("list_workflows_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterGetWorkflowDetailsMock registers the mock for GetWorkflowDetails
func (m *WorkflowsMock) RegisterGetWorkflowDetailsMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/workflows/details",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("get_details_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterTriggerThinkletMock registers the mock for TriggerThinklet
func (m *WorkflowsMock) RegisterTriggerThinkletMock() {
	httpmock.RegisterResponder(
		"POST",
		"=~^"+m.baseURL+"/api/v1/workflows/workflows/[^/]+/execution/[^/]+/trigger$",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("trigger_thinklet_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterUnauthorizedErrorMock registers a 401 unauthorized error mock
func (m *WorkflowsMock) RegisterUnauthorizedErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/workflows/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_unauthorized.json")
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterValidationErrorMock registers a 400 validation error mock
func (m *WorkflowsMock) RegisterValidationErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/workflows/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_validation.json")
			resp := httpmock.NewBytesResponse(400, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// loadMockData loads mock JSON data from file
func (m *WorkflowsMock) loadMockData(filename string) []byte {
	_, currentFile, _, _ := runtime.Caller(0)
	mockDir := filepath.Dir(currentFile)
	mockFile := filepath.Join(mockDir, filename)

	data, err := os.ReadFile(mockFile)
	if err != nil {
		panic("Failed to load mock data: " + err.Error())
	}

	return data
}
