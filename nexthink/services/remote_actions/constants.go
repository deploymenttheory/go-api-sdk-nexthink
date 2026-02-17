package remote_actions

const (
	// API Endpoints
	EndpointActExecute          = "/api/v1/act/execute"
	EndpointActRemoteActionList = "/api/v1/act/remote-action"
	EndpointActRemoteActionDetails = "/api/v1/act/remote-action/details"

	// Purpose enum values
	PurposeDataCollection = "DATA_COLLECTION"
	PurposeRemediation    = "REMEDIATION"

	// RunAs enum values
	RunAsLocalSystem         = "LOCAL_SYSTEM"
	RunAsInteractiveUser     = "INTERACTIVE_USER"
	RunAsDelegateToService   = "DELEGATE_TO_SERVICE"

	// Validation constraints
	MinDevices          = 1
	MaxDevices          = 10000
	MinExpiresInMinutes = 60
	MaxExpiresInMinutes = 10080
	MaxReasonLength     = 500
)
