package enrichment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEnrichmentRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *EnrichmentRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &EnrichmentRequest{
				Domain: "configuration",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  IdentificationDeviceName,
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{
							{
								Name:  FieldDeviceConfigurationTag,
								Value: "Production",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "enrichment request cannot be nil",
		},
		{
			name: "empty domain",
			req: &EnrichmentRequest{
				Domain: "",
				Enrichments: []Enrichment{
					{
						Identification: []Identification{
							{
								Name:  IdentificationDeviceName,
								Value: "DESKTOP-001",
							},
						},
						Fields: []Field{
							{
								Name:  FieldDeviceConfigurationTag,
								Value: "Production",
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "domain is required",
		},
		{
			name: "empty enrichments",
			req: &EnrichmentRequest{
				Domain:      "configuration",
				Enrichments: []Enrichment{},
			},
			wantErr: true,
			errMsg:  "enrichments is required",
		},
		{
			name: "too many enrichments",
			req: &EnrichmentRequest{
				Domain:      "configuration",
				Enrichments: make([]Enrichment, 5001),
			},
			wantErr: true,
			errMsg:  "enrichments cannot contain more than 5000 items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnrichmentRequest(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateEnrichment(t *testing.T) {
	tests := []struct {
		name       string
		enrichment *Enrichment
		index      int
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid enrichment",
			enrichment: &Enrichment{
				Identification: []Identification{
					{
						Name:  IdentificationDeviceName,
						Value: "DESKTOP-001",
					},
				},
				Fields: []Field{
					{
						Name:  FieldDeviceConfigurationTag,
						Value: "Production",
					},
				},
			},
			index:   0,
			wantErr: false,
		},
		{
			name: "no identification",
			enrichment: &Enrichment{
				Identification: []Identification{},
				Fields: []Field{
					{
						Name:  FieldDeviceConfigurationTag,
						Value: "Production",
					},
				},
			},
			index:   0,
			wantErr: true,
			errMsg:  "identification must contain exactly 1 item",
		},
		{
			name: "multiple identifications",
			enrichment: &Enrichment{
				Identification: []Identification{
					{
						Name:  IdentificationDeviceName,
						Value: "DESKTOP-001",
					},
					{
						Name:  IdentificationDeviceUID,
						Value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
					},
				},
				Fields: []Field{
					{
						Name:  FieldDeviceConfigurationTag,
						Value: "Production",
					},
				},
			},
			index:   0,
			wantErr: true,
			errMsg:  "identification must contain exactly 1 item",
		},
		{
			name: "empty fields",
			enrichment: &Enrichment{
				Identification: []Identification{
					{
						Name:  IdentificationDeviceName,
						Value: "DESKTOP-001",
					},
				},
				Fields: []Field{},
			},
			index:   0,
			wantErr: true,
			errMsg:  "fields is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnrichment(tt.index, tt.enrichment)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateIdentification(t *testing.T) {
	tests := []struct {
		name       string
		id         *Identification
		index      int
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid device name identification",
			id: &Identification{
				Name:  IdentificationDeviceName,
				Value: "DESKTOP-001",
			},
			index:   0,
			wantErr: false,
		},
		{
			name: "valid user SID identification",
			id: &Identification{
				Name:  IdentificationUserSID,
				Value: "S-1-5-21-1234567890-1234567890-1234567890-1001",
			},
			index:   0,
			wantErr: false,
		},
		{
			name: "empty name",
			id: &Identification{
				Name:  "",
				Value: "DESKTOP-001",
			},
			index:   0,
			wantErr: true,
			errMsg:  "identification[0].name is required",
		},
		{
			name: "empty value",
			id: &Identification{
				Name:  IdentificationDeviceName,
				Value: "",
			},
			index:   0,
			wantErr: true,
			errMsg:  "identification[0].value is required",
		},
		{
			name: "invalid name",
			id: &Identification{
				Name:  "invalid/field/name",
				Value: "DESKTOP-001",
			},
			index:   0,
			wantErr: true,
			errMsg:  "identification[0].name has invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIdentification(tt.index, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateField(t *testing.T) {
	tests := []struct {
		name             string
		field            *Field
		enrichmentIndex  int
		fieldIndex       int
		wantErr          bool
		errMsg           string
	}{
		{
			name: "valid field with string value",
			field: &Field{
				Name:  FieldDeviceConfigurationTag,
				Value: "Production",
			},
			enrichmentIndex: 0,
			fieldIndex:      0,
			wantErr:         false,
		},
		{
			name: "valid field with integer value",
			field: &Field{
				Name:  FieldDeviceVirtualizationLastUpdate,
				Value: 1609459200,
			},
			enrichmentIndex: 0,
			fieldIndex:      0,
			wantErr:         false,
		},
		{
			name: "empty field name",
			field: &Field{
				Name:  "",
				Value: "Production",
			},
			enrichmentIndex: 0,
			fieldIndex:      0,
			wantErr:         true,
			errMsg:          "fields[0].name is required",
		},
		{
			name: "nil field value",
			field: &Field{
				Name:  FieldDeviceConfigurationTag,
				Value: nil,
			},
			enrichmentIndex: 0,
			fieldIndex:      0,
			wantErr:         true,
			errMsg:          "fields[0].value is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateField(tt.enrichmentIndex, tt.fieldIndex, tt.field)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
