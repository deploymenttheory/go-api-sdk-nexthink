package workflows

import "time"

// TriggerWorkflowV1Request represents a V1 workflow trigger request
type TriggerWorkflowV1Request struct {
	// WorkflowID is the ID of the workflow to execute
	WorkflowID string `json:"workflowId"`

	// Devices are Nexthink Collector IDs of devices (max 10000)
	// Note: If devices are included, users are optional
	Devices []string `json:"devices,omitempty"`

	// Users are security IDs of users (max 10000)
	// Note: If users are included, devices are optional
	Users []string `json:"users,omitempty"`

	// Params are optional parameters to send to the workflow
	Params map[string]string `json:"params,omitempty"`
}

// TriggerWorkflowV2Request represents a V2 workflow trigger request with external identifiers
type TriggerWorkflowV2Request struct {
	// WorkflowID is the ID of the workflow to execute
	WorkflowID string `json:"workflowId"`

	// Devices are device identifiers (Collector UIDs, names, or UIDs) - max 10000
	Devices []DeviceData `json:"devices,omitempty"`

	// Users are user identifiers (SID, UPN, or UID) - max 10000
	Users []UserData `json:"users,omitempty"`

	// Params are optional parameters to send to the workflow
	Params map[string]string `json:"params,omitempty"`
}

// DeviceData represents device identification data
type DeviceData struct {
	// Name is the device name
	Name string `json:"name,omitempty"`

	// UID is the globally unique device identifier (UUID format)
	UID string `json:"uid,omitempty"`

	// CollectorUID is the Nexthink Collector UUID of the device
	CollectorUID string `json:"collectorUid,omitempty"`
}

// UserData represents user identification data
type UserData struct {
	// UID is the globally unique user identifier (UUID format)
	UID string `json:"uid,omitempty"`

	// UPN is the user's principal name (email format)
	UPN string `json:"upn,omitempty"`

	// SID is the security identifier of the user
	SID string `json:"sid,omitempty"`
}

// TriggerWorkflowResponse represents the response from triggering a workflow
type TriggerWorkflowResponse struct {
	// RequestUUID is the request ID
	// Use this ID to query workflow executions in NQL (workflow.executions.request_id)
	RequestUUID string `json:"requestUuid"`

	// ExecutionsUUIDs is the list of execution IDs for each object targeted
	// Query using workflow.executions.execution_id
	ExecutionsUUIDs []string `json:"executionsUuids"`
}

// ThinkletTriggerRequest represents a request to trigger a waiting workflow execution
type ThinkletTriggerRequest struct {
	// Parameters are optional parameters to send to the thinklet
	Parameters map[string]string `json:"parameters,omitempty"`
}

// ThinkletTriggerResponse represents the response from triggering a thinklet
type ThinkletTriggerResponse struct {
	// RequestUUID is the request ID
	RequestUUID string `json:"requestUuid"`
}

// Workflow represents a workflow configuration
type Workflow struct {
	// ID is the workflow ID
	ID string `json:"id"`

	// UUID is the workflow UUID
	UUID string `json:"uuid"`

	// Name is the workflow name
	Name string `json:"name"`

	// Description is the workflow description
	Description string `json:"description"`

	// Status is the workflow status (ACTIVE, INACTIVE)
	Status string `json:"status"`

	// LastUpdateTime is the last update timestamp
	LastUpdateTime time.Time `json:"lastUpdateTime"`

	// TriggerMethods are the available trigger methods
	TriggerMethods []string `json:"triggerMethods"`

	// Versions are the workflow versions
	Versions []WorkflowVersion `json:"versions"`
}

// WorkflowVersion represents a workflow version
type WorkflowVersion struct {
	// VersionNumber is the version number
	VersionNumber int `json:"versionNumber,omitempty"`

	// IsActive indicates if this version is active
	IsActive bool `json:"isActive,omitempty"`
}

// ErrorResponse represents an error response from the Workflows API
type ErrorResponse struct {
	// Code is the error code
	Code string `json:"code"`

	// Details provides error details
	Details string `json:"details"`
}
