package enrichment

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

type (
	// EnrichmentServiceInterface defines the interface for enrichment operations
	//
	// Nexthink API docs: https://docs.nexthink.com/api/enrichment/enrich-fields-for-given-objects
	EnrichmentServiceInterface interface {
		// EnrichFields enriches fields for given objects
		//
		// Enrichment can be done for:
		//  - Manual custom fields (any object)
		//  - Virtualization fields (devices only)
		//  - Configuration_tag (devices only)
		//  - Organization field (users only)
		//  - Entra ID fields (users only)
		//
		// The request can contain 1-5000 enrichment operations.
		// Each enrichment identifies an object (device, user, binary, or package) and specifies
		// the fields to enrich with their desired values.
		//
		// Response types:
		//  - 200 OK: All objects processed successfully (SuccessResponse)
		//  - 207 Multi-Status: Some objects processed, others failed (PartialSuccessResponse)
		//  - 400 Bad Request: All objects failed (BadRequestResponse)
		//  - 401 Unauthorized: Invalid authentication
		//  - 403 Forbidden: No permission to trigger enrichment
		//
		// Nexthink API docs: https://docs.nexthink.com/api/enrichment/enrich-fields-for-given-objects
		EnrichFields(ctx context.Context, req *EnrichmentRequest) (any, *interfaces.Response, error)
	}

	// Service implements the EnrichmentServiceInterface
	Service struct {
		client interfaces.HTTPClient
	}
)

var _ EnrichmentServiceInterface = (*Service)(nil)

// NewService creates a new enrichment service instance
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// =============================================================================
// Enrich Fields Operations
// =============================================================================

// EnrichFields enriches fields for given objects
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/enrichment/data/fields
// https://docs.nexthink.com/api/enrichment/enrich-fields-for-given-objects
func (s *Service) EnrichFields(ctx context.Context, req *EnrichmentRequest) (any, *interfaces.Response, error) {
	if err := ValidateEnrichmentRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointEnrichmentDataFields

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result any
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)

	// Return the response regardless of error for status code checking
	return result, resp, err
}
