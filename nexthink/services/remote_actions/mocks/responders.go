package mocks

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarcoal/httpmock"
)

// RemoteActionsMock provides mock responses for Remote Actions API endpoints
type RemoteActionsMock struct {
	baseURL string
}

// NewRemoteActionsMock creates a new RemoteActionsMock instance
func NewRemoteActionsMock(baseURL string) *RemoteActionsMock {
	return &RemoteActionsMock{baseURL: baseURL}
}

// RegisterMocks registers all successful response mocks
func (m *RemoteActionsMock) RegisterMocks() {
	m.RegisterExecuteMock()
	m.RegisterListRemoteActionsMock()
	m.RegisterGetRemoteActionDetailsMock()
}

// RegisterErrorMocks registers all error response mocks
func (m *RemoteActionsMock) RegisterErrorMocks() {
	m.RegisterUnauthorizedErrorMock()
	m.RegisterValidationErrorMock()
}

// RegisterExecuteMock registers the mock for Execute
func (m *RemoteActionsMock) RegisterExecuteMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/act/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("execute_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterListRemoteActionsMock registers the mock for ListRemoteActions
func (m *RemoteActionsMock) RegisterListRemoteActionsMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/act/remote-action",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("list_remote_actions_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterGetRemoteActionDetailsMock registers the mock for GetRemoteActionDetails
func (m *RemoteActionsMock) RegisterGetRemoteActionDetailsMock() {
	httpmock.RegisterResponder(
		"GET",
		m.baseURL+"/api/v1/act/remote-action/details",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("get_details_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterUnauthorizedErrorMock registers a 401 unauthorized error mock
func (m *RemoteActionsMock) RegisterUnauthorizedErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/act/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_unauthorized.json")
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterValidationErrorMock registers a 400 validation error mock
func (m *RemoteActionsMock) RegisterValidationErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/act/execute",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_validation.json")
			resp := httpmock.NewBytesResponse(400, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// loadMockData loads mock JSON data from file
func (m *RemoteActionsMock) loadMockData(filename string) []byte {
	_, currentFile, _, _ := runtime.Caller(0)
	mockDir := filepath.Dir(currentFile)
	mockFile := filepath.Join(mockDir, filename)

	data, err := os.ReadFile(mockFile)
	if err != nil {
		panic("Failed to load mock data: " + err.Error())
	}

	return data
}
