package mocks

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarcoal/httpmock"
)

// CampaignsMock provides mock responses for Campaigns API endpoints
type CampaignsMock struct {
	baseURL string
}

// NewCampaignsMock creates a new CampaignsMock instance
func NewCampaignsMock(baseURL string) *CampaignsMock {
	return &CampaignsMock{baseURL: baseURL}
}

// RegisterMocks registers all successful response mocks
func (m *CampaignsMock) RegisterMocks() {
	m.RegisterTriggerCampaignSuccessMock()
	m.RegisterTriggerCampaignPartialSuccessMock()
}

// RegisterErrorMocks registers all error response mocks
func (m *CampaignsMock) RegisterErrorMocks() {
	m.RegisterUnauthorizedErrorMock()
	m.RegisterValidationErrorMock()
}

// RegisterTriggerCampaignSuccessMock registers the mock for TriggerCampaign with all successful
func (m *CampaignsMock) RegisterTriggerCampaignSuccessMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/euf/campaign/trigger",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("trigger_campaign_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterTriggerCampaignPartialSuccessMock registers the mock for TriggerCampaign with some failures
func (m *CampaignsMock) RegisterTriggerCampaignPartialSuccessMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/euf/campaign/trigger",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("trigger_campaign_partial_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterUnauthorizedErrorMock registers a 401 unauthorized error mock
func (m *CampaignsMock) RegisterUnauthorizedErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/euf/campaign/trigger",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_unauthorized.json")
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterValidationErrorMock registers a 400 validation error mock
func (m *CampaignsMock) RegisterValidationErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/euf/campaign/trigger",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_validation.json")
			resp := httpmock.NewBytesResponse(400, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// loadMockData loads mock JSON data from file
func (m *CampaignsMock) loadMockData(filename string) []byte {
	_, currentFile, _, _ := runtime.Caller(0)
	mockDir := filepath.Dir(currentFile)
	mockFile := filepath.Join(mockDir, filename)

	data, err := os.ReadFile(mockFile)
	if err != nil {
		panic("Failed to load mock data: " + err.Error())
	}

	return data
}
