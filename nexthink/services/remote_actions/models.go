package remote_actions

// TriggerRemoteActionRequest represents the request to trigger a remote action
type TriggerRemoteActionRequest struct {
	// RemoteActionID is the ID of the remote action to execute
	RemoteActionID string `json:"remoteActionId"`

	// Params are parameters to send to the script (optional)
	Params map[string]string `json:"params,omitempty"`

	// Devices are Nexthink Collector IDs of the devices (1-10000)
	Devices []string `json:"devices"`

	// ExpiresInMinutes is time before execution expires if device doesn't come online (60-10080)
	ExpiresInMinutes int `json:"expiresInMinutes,omitempty"`

	// TriggerInfo contains additional trigger metadata (optional)
	TriggerInfo *TriggerInfoRequest `json:"triggerInfo,omitempty"`
}

// TriggerInfoRequest contains metadata about the trigger source
type TriggerInfoRequest struct {
	// ExternalSource is the external application/tool name
	ExternalSource string `json:"externalSource,omitempty"`

	// Reason is the reason behind triggering the action (max 500 chars)
	Reason string `json:"reason,omitempty"`

	// ExternalReference is the external ticket reference ID
	ExternalReference string `json:"externalReference,omitempty"`
}

// TriggerRemoteActionResponse represents the response from triggering a remote action
type TriggerRemoteActionResponse struct {
	// RequestID is the Nexthink ID of the request created
	// Use this ID to query remote action executions in NQL
	RequestID string `json:"requestId"`

	// ExpiresInMinutes is the expiration time for the execution
	ExpiresInMinutes int `json:"expiresInMinutes,omitempty"`
}

// RemoteAction represents a remote action configuration
type RemoteAction struct {
	// ID is the remote action ID
	ID string `json:"id"`

	// UUID is the remote action UUID
	UUID string `json:"uuid"`

	// Name is the remote action name
	Name string `json:"name"`

	// Description is the remote action description
	Description string `json:"description"`

	// Origin indicates the source of the remote action
	Origin string `json:"origin"`

	// BuiltInContentVersion is the version of built-in content
	BuiltInContentVersion string `json:"builtInContentVersion"`

	// Purpose is the list of purposes for this remote action
	Purpose []string `json:"purpose"`

	// Targeting contains the targeting configuration
	Targeting Targeting `json:"targeting"`

	// ScriptInfo contains script execution details
	ScriptInfo ScriptInfo `json:"scriptInfo"`
}

// Targeting represents the targeting configuration for a remote action
type Targeting struct {
	// APIEnabled indicates if the remote action can be triggered via API
	APIEnabled bool `json:"apiEnabled"`

	// ManualEnabled indicates if manual triggering is enabled
	ManualEnabled bool `json:"manualEnabled"`

	// WorkflowEnabled indicates if workflow triggering is enabled
	WorkflowEnabled bool `json:"workflowEnabled"`

	// ManualAllowMultipleDevices indicates if multiple device selection is allowed
	ManualAllowMultipleDevices bool `json:"manualAllowMultipleDevices"`
}

// ScriptInfo contains script execution details
type ScriptInfo struct {
	// ExecutionServiceDelegate is the execution service delegate
	ExecutionServiceDelegate string `json:"executionServiceDelegate"`

	// RunAs specifies the execution context
	RunAs string `json:"runAs"`

	// TimeoutSeconds is the script timeout in seconds
	TimeoutSeconds int `json:"timeoutSeconds"`

	// HasScriptWindows indicates if Windows script is available
	HasScriptWindows bool `json:"hasScriptWindows"`

	// HasScriptMacOS indicates if macOS script is available
	HasScriptMacOS bool `json:"hasScriptMacOs"`

	// Inputs is the list of script input parameters
	Inputs []Input `json:"inputs"`

	// Outputs is the list of script output parameters
	Outputs []Output `json:"outputs"`
}

// Input represents a script input parameter
type Input struct {
	// ID is the input parameter ID
	ID string `json:"id"`

	// Name is the input parameter name
	Name string `json:"name"`

	// Description is the input parameter description
	Description string `json:"description"`

	// UsedByWindows indicates if used by Windows script
	UsedByWindows bool `json:"usedByWindows"`

	// UsedByMacOS indicates if used by macOS script
	UsedByMacOS bool `json:"usedByMacOs"`

	// Options are the available option values
	Options []string `json:"options"`

	// AllowCustomValue indicates if custom values are allowed
	AllowCustomValue bool `json:"allowCustomValue"`
}

// Output represents a script output parameter
type Output struct {
	// ID is the output parameter ID
	ID string `json:"id"`

	// Name is the output parameter name
	Name string `json:"name"`

	// Type is the output parameter type
	Type string `json:"type"`

	// Description is the output parameter description
	Description string `json:"description"`

	// UsedByWindows indicates if used by Windows script
	UsedByWindows bool `json:"usedByWindows"`

	// UsedByMacOS indicates if used by macOS script
	UsedByMacOS bool `json:"usedByMacOs"`
}

// ErrorResponse represents an error response from the Remote Actions API
type ErrorResponse struct {
	// Code is the error code
	Code string `json:"code"`

	// Message is the error message
	Message string `json:"message"`
}
