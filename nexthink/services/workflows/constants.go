package workflows

const (
	// API Endpoints
	EndpointWorkflowsExecuteV1     = "/api/v1/workflows/execute"
	EndpointWorkflowsExecuteV2     = "/api/v2/workflows/execute"
	EndpointWorkflowsList          = "/api/v1/workflows"
	EndpointWorkflowsDetails       = "/api/v1/workflows/details"
	EndpointWorkflowsTriggerEvent  = "/api/v1/workflows/workflows/%s/execution/%s/trigger" // Format with workflowUUID, executionUUID

	// Workflow status values
	WorkflowStatusActive   = "ACTIVE"
	WorkflowStatusInactive = "INACTIVE"

	// Workflow dependency values
	DependencyUser         = "USER"
	DependencyDevice       = "DEVICE"
	DependencyUserAndDevice = "USER_AND_DEVICE"
	DependencyNone         = "NONE"

	// Trigger method values
	TriggerMethodAPI            = "API"
	TriggerMethodManual         = "MANUAL"
	TriggerMethodManualMultiple = "MANUAL_MULTIPLE"
	TriggerMethodScheduler      = "SCHEDULER"

	// Validation constraints
	MaxDevices = 10000
	MaxUsers   = 10000

	// UUID regex pattern
	UUIDPattern = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	
	// User SID regex pattern
	UserSIDPattern = "^S(-\\d+){2,10}$|^0$"
	
	// User UPN regex pattern (email format)
	UserUPNPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
)
