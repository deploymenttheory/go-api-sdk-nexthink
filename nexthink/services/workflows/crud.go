package workflows

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/interfaces"
)

type (
	// WorkflowsServiceInterface defines the interface for workflows operations
	//
	// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow
	WorkflowsServiceInterface interface {
		// ExecuteV1 triggers a workflow execution using internal IDs (Collector IDs, SIDs)
		//
		// Triggers a workflow using Nexthink internal identifiers:
		//  - Devices: Nexthink Collector IDs (max 10000)
		//  - Users: Security IDs (SID) (max 10000)
		//
		// Returns RequestUUID and ExecutionsUUIDs to track execution via NQL.
		//
		// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#get-api-v1-workflows-details
		ExecuteV1(ctx context.Context, req *ExecutionRequestV1) (*ExecutionResponse, *interfaces.Response, error)

		// ExecuteV2 triggers a workflow execution using external identifiers
		//
		// Triggers a workflow using external identifiers:
		//  - Devices: Names, UIDs, or Collector UIDs (max 10000)
		//  - Users: SID, UPN, or UID (max 10000)
		//
		// Returns RequestUUID and ExecutionsUUIDs to track execution via NQL.
		//
		// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#trigger-a-workflow-v2
		ExecuteV2(ctx context.Context, req *ExecutionRequestV2) (*ExecutionResponse, *interfaces.Response, error)

		// ListWorkflows retrieves all workflows with their configurations
		//
		// Returns a list of all workflows including:
		//  - ID, UUID, Name, Description, Status
		//  - Available trigger methods
		//  - Version information
		//
		// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#list-workflows
		ListWorkflows(ctx context.Context) ([]Workflow, *interfaces.Response, error)

		// GetWorkflowDetails retrieves the configuration of a specific workflow by NQL ID
		//
		// Returns detailed configuration for a specific workflow including:
		//  - Full workflow metadata
		//  - Available trigger methods
		//  - Version history
		//
		// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#get-workflow
		GetWorkflowDetails(ctx context.Context, nqlID string) (*Workflow, *interfaces.Response, error)

		// TriggerThinklet triggers a waiting workflow execution via a thinklet
		//
		// When a workflow execution is waiting for a thinklet trigger, use this endpoint
		// to send the trigger event with optional parameters.
		//
		// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#trigger-wait-for-event
		TriggerThinklet(ctx context.Context, workflowUUID, executionUUID string, req *ThinkletTriggerRequest) (*ThinkletTriggerResponse, *interfaces.Response, error)
	}

	// Service implements the WorkflowsServiceInterface
	Service struct {
		client interfaces.HTTPClient
	}
)

var _ WorkflowsServiceInterface = (*Service)(nil)

// NewService creates a new workflows service instance
func NewService(client interfaces.HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// =============================================================================
// Execute Workflow Operations
// =============================================================================

// ExecuteV1 triggers a workflow execution using internal IDs (Collector IDs, SIDs)
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/workflows/execute
// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#trigger-a-workflow-v1
func (s *Service) ExecuteV1(ctx context.Context, req *ExecutionRequestV1) (*ExecutionResponse, *interfaces.Response, error) {
	if err := ValidateExecutionRequestV1(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointWorkflowsExecuteV1

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result ExecutionResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ExecuteV2 triggers a workflow execution using external identifiers
// URL: POST https://instance.api.region.nexthink.cloud/api/v2/workflows/execute
// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#trigger-a-workflow-v2
func (s *Service) ExecuteV2(ctx context.Context, req *ExecutionRequestV2) (*ExecutionResponse, *interfaces.Response, error) {
	if err := ValidateExecutionRequestV2(req); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointWorkflowsExecuteV2

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var result ExecutionResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// =============================================================================
// List Workflows Operations
// =============================================================================

// ListWorkflows retrieves all workflows with their configurations
// URL: GET https://instance.api.region.nexthink.cloud/api/v1/workflows
// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#list-workflows
func (s *Service) ListWorkflows(ctx context.Context) ([]Workflow, *interfaces.Response, error) {
	endpoint := EndpointWorkflowsList

	headers := map[string]string{
		"Accept": "application/json",
	}

	var result []Workflow
	resp, err := s.client.Get(ctx, endpoint, nil, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// =============================================================================
// Get Workflow Details Operations
// =============================================================================

// GetWorkflowDetails retrieves the configuration of a specific workflow by NQL ID
// URL: GET https://instance.api.region.nexthink.cloud/api/v1/workflows/details?nql-id={nqlId}
// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#get-api-v1-workflows-details
func (s *Service) GetWorkflowDetails(ctx context.Context, nqlID string) (*Workflow, *interfaces.Response, error) {
	if err := ValidateNqlID(nqlID); err != nil {
		return nil, nil, err
	}

	endpoint := EndpointWorkflowsDetails

	queryParams := map[string]string{
		"nql-id": nqlID,
	}

	headers := map[string]string{
		"Accept": "application/json",
	}

	var result Workflow
	resp, err := s.client.Get(ctx, endpoint, queryParams, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// =============================================================================
// Trigger Thinklet Operations
// =============================================================================

// TriggerThinklet triggers a waiting workflow execution via a thinklet
// URL: POST https://instance.api.region.nexthink.cloud/api/v1/workflows/workflows/{workflowUuid}/execution/{executionUuid}/trigger
// Nexthink API docs: https://docs.nexthink.com/api/workflows/trigger-a-workflow#post-api-v1-workflows-workflows-workflowuuid-execution-executionuuid-trigger
func (s *Service) TriggerThinklet(ctx context.Context, workflowUUID, executionUUID string, req *ThinkletTriggerRequest) (*ThinkletTriggerResponse, *interfaces.Response, error) {
	if err := ValidateUUID(workflowUUID, "workflow"); err != nil {
		return nil, nil, err
	}

	if err := ValidateUUID(executionUUID, "execution"); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf(EndpointWorkflowsTriggerEvent, workflowUUID, executionUUID)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// If req is nil, create empty request
	if req == nil {
		req = &ThinkletTriggerRequest{}
	}

	var result ThinkletTriggerResponse
	resp, err := s.client.Post(ctx, endpoint, req, headers, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
