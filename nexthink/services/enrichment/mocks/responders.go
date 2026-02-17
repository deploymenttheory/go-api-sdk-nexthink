package mocks

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarcoal/httpmock"
)

// EnrichmentMock provides mock responses for Enrichment API endpoints
type EnrichmentMock struct {
	baseURL string
}

// NewEnrichmentMock creates a new EnrichmentMock instance
func NewEnrichmentMock(baseURL string) *EnrichmentMock {
	return &EnrichmentMock{baseURL: baseURL}
}

// RegisterMocks registers all successful response mocks
func (m *EnrichmentMock) RegisterMocks() {
	m.RegisterEnrichFieldsSuccessMock()
	m.RegisterEnrichFieldsPartialSuccessMock()
}

// RegisterErrorMocks registers all error response mocks
func (m *EnrichmentMock) RegisterErrorMocks() {
	m.RegisterUnauthorizedErrorMock()
	m.RegisterValidationErrorMock()
}

// RegisterEnrichFieldsSuccessMock registers the mock for EnrichFields with success status
func (m *EnrichmentMock) RegisterEnrichFieldsSuccessMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/enrichment/data/fields",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("enrich_fields_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterEnrichFieldsPartialSuccessMock registers the mock for EnrichFields with partial success
func (m *EnrichmentMock) RegisterEnrichFieldsPartialSuccessMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/enrichment/data/fields",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("enrich_fields_partial_success.json")
			resp := httpmock.NewBytesResponse(200, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterUnauthorizedErrorMock registers a 401 unauthorized error mock
func (m *EnrichmentMock) RegisterUnauthorizedErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/enrichment/data/fields",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_unauthorized.json")
			resp := httpmock.NewBytesResponse(401, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// RegisterValidationErrorMock registers a 400 validation error mock
func (m *EnrichmentMock) RegisterValidationErrorMock() {
	httpmock.RegisterResponder(
		"POST",
		m.baseURL+"/api/v1/enrichment/data/fields",
		func(req *http.Request) (*http.Response, error) {
			mockData := m.loadMockData("error_validation.json")
			resp := httpmock.NewBytesResponse(400, mockData)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// loadMockData loads mock JSON data from file
func (m *EnrichmentMock) loadMockData(filename string) []byte {
	_, currentFile, _, _ := runtime.Caller(0)
	mockDir := filepath.Dir(currentFile)
	mockFile := filepath.Join(mockDir, filename)

	data, err := os.ReadFile(mockFile)
	if err != nil {
		panic("Failed to load mock data: " + err.Error())
	}

	return data
}
