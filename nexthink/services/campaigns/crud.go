package campaigns

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

type (
	// CampaignsServiceInterface defines the interface for campaign operations
	//
	// Nexthink API docs: https://docs.nexthink.com/api/campaigns/trigger-a-campaign
	CampaignsServiceInterface interface {
		// TriggerCampaign triggers the sending of a campaign to specific users
		//
		// Triggers a campaign to be sent to a list of users identified by their SIDs.
		// The campaign must be pre-configured in Nexthink and identified by its NQL ID.
		// Parameters can be provided to customize the campaign content dynamically.
		//
		// The response includes:
		//  - For successful user requests: RequestId (used to retrieve status and answers later)
		//  - For failed user requests: Message explaining the failure reason
		//
		// Duplicate SIDs in the request are automatically filtered out from the response.
		//
		// Nexthink API docs: https://docs.nexthink.com/api/campaigns/trigger-a-campaign
		TriggerCampaign(ctx context.Context, req *TriggerRequest) (*TriggerSuccessResponse, *interfaces.Response, error)
	}

	// Service implements the CampaignsServiceInterface
	Service struct {
		client interfaces.HTTPClient
	}
)

var _ CampaignsServiceInterface = (*Service)(nil)

// NewService creates a new campaigns service instance
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// =============================================================================
// Trigger Campaign Operations
// =============================================================================

// TriggerCampaign triggers the sending of a campaign to specific users
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/euf/campaign/trigger
// https://docs.nexthink.com/api/campaigns/trigger-a-campaign
func (s *Service) TriggerCampaign(ctx context.Context, req *TriggerRequest) (*TriggerSuccessResponse, *interfaces.Response, error) {
	if err := ValidateTriggerRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointCampaignTrigger

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result TriggerSuccessResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
