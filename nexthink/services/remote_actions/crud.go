package remote_actions

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

type (
	// RemoteActionsServiceInterface defines the interface for remote actions operations
	//
	// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api
	RemoteActionsServiceInterface interface {
		// TriggerRemoteAction triggers the execution of a remote action for a set of devices
		//
		// Triggers a remote action to execute on specified devices identified by their Collector IDs.
		// The remote action must be pre-configured in Nexthink and API-enabled.
		//
		// Request includes:
		//  - RemoteActionID: The NQL ID of the remote action to execute
		//  - Devices: List of Nexthink Collector IDs (1-10000)
		//  - Params: Optional script parameters as key-value pairs
		//  - ExpiresInMinutes: Expiration time if device doesn't come online (60-10080 minutes)
		//  - TriggerInfo: Optional metadata (external source, reason, reference)
		//
		// Returns a RequestID that can be used to query remote action executions in NQL.
		//
		// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api#trigger-a-remote-action
		TriggerRemoteAction(ctx context.Context, req *TriggerRemoteActionRequest) (*TriggerRemoteActionResponse, *interfaces.Response, error)

		// ListRemoteActions retrieves all remote actions with their configuration information
		//
		// Returns a list of all available remote actions including:
		//  - ID, UUID, Name, Description
		//  - Purpose (DATA_COLLECTION, REMEDIATION)
		//  - Targeting configuration (API/Manual/Workflow enabled)
		//  - Script information (inputs, outputs, timeout, run-as)
		//
		// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api#list-remote-actions
		ListRemoteActions(ctx context.Context) ([]RemoteAction, *interfaces.Response, error)

		// GetRemoteActionDetails retrieves the configuration of a specific remote action by NQL ID
		//
		// Returns detailed configuration for a specific remote action including:
		//  - Full targeting configuration
		//  - Script details (platform support, inputs, outputs)
		//  - Execution settings
		//
		// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api#get-remote-action-details
		GetRemoteActionDetails(ctx context.Context, nqlID string) (*RemoteAction, *interfaces.Response, error)
	}

	// Service implements the RemoteActionsServiceInterface
	Service struct {
		client interfaces.HTTPClient
	}
)

var _ RemoteActionsServiceInterface = (*Service)(nil)

// NewService creates a new remote actions service instance
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// =============================================================================
// Trigger Remote Action Operations
// =============================================================================

// TriggerRemoteAction triggers the execution of a remote action for a set of devices
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/act/execute
// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api#trigger-a-remote-action
func (s *Service) TriggerRemoteAction(ctx context.Context, req *TriggerRemoteActionRequest) (*TriggerRemoteActionResponse, *interfaces.Response, error) {
	if err := ValidateTriggerRemoteActionRequest(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointActExecute

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result TriggerRemoteActionResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// =============================================================================
// List Remote Actions Operations
// =============================================================================

// ListRemoteActions retrieves all remote actions with their configuration information
// URL: GET https://instance.api.region.nexthink.cloud/api/v1/act/remote-action
// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api#list-remote-actions
func (s *Service) ListRemoteActions(ctx context.Context) ([]RemoteAction, *interfaces.Response, error) {
	endpoint := EndpointActRemoteActionList

	headers := map[string]string{
		"Accept": "application/json",
	}

	var result []RemoteAction
	resp, err := s.client.Get(ctx, endpoint, nil, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// =============================================================================
// Get Remote Action Details Operations
// =============================================================================

// GetRemoteActionDetails retrieves the configuration of a specific remote action by NQL ID
// URL: GET https://instance.api.region.nexthink.cloud/api/v1/act/remote-action/details?nql-id={nqlId}
// Nexthink API docs: https://docs.nexthink.com/api/remote-actions/remote-actions-api#get-remote-action-details
func (s *Service) GetRemoteActionDetails(ctx context.Context, nqlID string) (*RemoteAction, *interfaces.Response, error) {
	if err := ValidateNqlID(nqlID); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointActRemoteActionDetails

	queryParams := map[string]string{
		"nql-id": nqlID,
	}

	headers := map[string]string{
		"Accept": "application/json",
	}

	var result RemoteAction
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
