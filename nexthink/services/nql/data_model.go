package nql

// Data model constants for NQL queries
// Provides type-safe access to table names, namespaces, and common fields

// =============================================================================
// Tables (shortcuts)
// =============================================================================

const (
	// Inventory Object Tables
	TableDevices      = "devices"      // device.devices shortcut
	TableUsers        = "users"        // user.users shortcut
	TableApplications = "applications" // application.applications shortcut
	TableBinaries     = "binaries"     // binary.binaries shortcut
	TableCampaigns    = "campaigns"    // campaign.campaigns shortcut
	TablePackages     = "packages"     // package.packages shortcut
)

// =============================================================================
// Namespaces
// =============================================================================

const (
	NamespaceApplication        = "application"
	NamespaceBinary             = "binary"
	NamespaceCampaign           = "campaign"
	NamespaceCollaboration      = "collaboration"
	NamespaceConnection         = "connection"
	NamespaceConnectivity       = "connectivity"
	NamespaceDevice             = "device"
	NamespaceDevicePerformance  = "device_performance"
	NamespaceDex                = "dex"
	NamespaceExecution          = "execution"
	NamespacePackage            = "package"
	NamespaceRemoteAction       = "remote_action"
	NamespaceSession            = "session"
	NamespaceUser               = "user"
	NamespaceWeb                = "web"
	NamespaceWorkflow           = "workflow"
)

// =============================================================================
// Event Tables (namespace.table format)
// =============================================================================

const (
	// Collaboration Events
	TableCollaborationSessions = "collaboration.sessions"
	
	// Connection Events
	TableConnectionEvents = "connection.events"
	
	// Connectivity Events
	TableConnectivityEvents = "connectivity.events"
	
	// Device Performance Events
	TableDevicePerformanceEvents       = "device_performance.events"
	TableDevicePerformanceBoots        = "device_performance.boots"
	TableDevicePerformanceHardResets   = "device_performance.hard_resets"
	TableDevicePerformanceSystemCrashes = "device_performance.system_crashes"
	
	// DEX Events
	TableDexScores            = "dex.scores"
	TableDexApplicationScores = "dex.application_scores"
	
	// Execution Events
	TableExecutionEvents  = "execution.events"
	TableExecutionCrashes = "execution.crashes"
	
	// Remote Action Events
	TableRemoteActionExecutions = "remote_action.executions"
	
	// Session Events
	TableSessionEvents = "session.events"
	TableSessionLogins = "session.logins"
	TableSessionVDIEvents = "session.vdi_events"
	
	// Web Events
	TableWebEvents      = "web.events"
	TableWebPageViews   = "web.page_views"
	TableWebErrors      = "web.errors"
	TableWebTransactions = "web.transactions"
	
	// Workflow Events
	TableWorkflowExecutions = "workflow.executions"
)

// =============================================================================
// Common Fields (by object/namespace)
// =============================================================================

// Device Fields
const (
	FieldDeviceName                 = "device.name"
	FieldDeviceEntity               = "device.entity"
	FieldDeviceLastSeen             = "device.last_seen"
	FieldDeviceDaysSinceLastSeen    = "device.days_since_last_seen"
	FieldDeviceCollectorUID         = "device.collector.uid"
	FieldDeviceOrganizationEntity   = "device.organization.entity"
	FieldDevicePublicIPCountry      = "device.public_ip.country"
	FieldDevicePublicIPISP          = "device.public_ip.isp"
)

// Operating System Fields
const (
	FieldOSName                 = "operating_system.name"
	FieldOSPlatform             = "operating_system.platform"
	FieldOSVersion              = "operating_system.version"
	FieldOSLastUpdate           = "operating_system.last_update"
)

// Hardware Fields
const (
	FieldHardwareType          = "hardware.type"
	FieldHardwareManufacturer  = "hardware.manufacturer"
	FieldHardwareModel         = "hardware.model"
	FieldHardwareMemory        = "hardware.memory"
	FieldHardwareProcessor     = "hardware.processor"
)

// User Fields
const (
	FieldUserName     = "user.name"
	FieldUsername     = "username"
	FieldUserType     = "user.type"
	FieldUserSID      = "user.sid"
	FieldUserEntity   = "user.entity"
)

// Application Fields
const (
	FieldApplicationName    = "application.name"
	FieldApplicationVersion = "application.version"
	FieldApplicationVendor  = "application.vendor"
)

// Binary Fields
const (
	FieldBinaryName         = "binary.name"
	FieldBinaryVersion      = "binary.version"
	FieldBinaryPlatform     = "binary.platform"
	FieldBinaryArchitecture = "binary.architecture"
	FieldBinarySize         = "binary.size"
	FieldBinaryMD5Hash      = "binary.md5_hash"
)

// Web Fields
const (
	FieldPageLoadTimeOverall = "page_load_time.overall"
	FieldPageLoadTimeBackend = "page_load_time.backend"
	FieldPageLoadTimeClient  = "page_load_time.client"
	FieldPageLoadTimeNetwork = "page_load_time.network"
	FieldExperienceLevel     = "experience_level"
	FieldNumberOfPageViews   = "number_of_page_views"
	FieldNumberOfErrors      = "number_of_errors"
)

// Execution Fields
const (
	FieldNumberOfCrashes     = "number_of_crashes"
	FieldNumberOfFreezes     = "number_of_freezes"
	FieldExecutionDuration   = "execution_duration"
	FieldProcessVisibility   = "process_visibility"
)

// Session Fields
const (
	FieldTimeUntilDesktopIsVisible = "time_until_desktop_is_visible"
	FieldTimeUntilDesktopIsReady   = "time_until_desktop_is_ready"
	FieldLogonType                 = "logon_type"
)

// DEX Score Fields
const (
	FieldDexScoreValue                      = "value"
	FieldDexEndpointValue                   = "endpoint.value"
	FieldDexCollaborationValue              = "collaboration.value"
	FieldDexLogonSpeedValue                 = "endpoint.logon_speed_value"
	FieldDexBootSpeedValue                  = "endpoint.boot_speed_value"
	FieldDexSoftwareReliabilityValue        = "endpoint.software_reliability_value"
	FieldDexVirtualSessionLagValue          = "endpoint.virtual_session_lag_value"
	FieldDexLogonSpeedScoreImpact           = "endpoint.logon_speed_score_impact"
	FieldDexBootSpeedScoreImpact            = "endpoint.boot_speed_score_impact"
	FieldDexSoftwareReliabilityScoreImpact  = "endpoint.software_reliability_score_impact"
)

// Context Fields
const (
	FieldContextLocationCountry = "context.location.country"
	FieldContextLocationType    = "context.location.type"
	FieldContextLocationState   = "context.location.state"
	FieldContextDevicePlatform  = "context.device_platform"
)

// Connectivity Fields
const (
	FieldConnectionType            = "connection_type"
	FieldWifiSignalStrength        = "wifi.signal_strength"
	FieldWifiReceiveRate           = "wifi.receive_rate"
	FieldWifiNoiseLevel            = "wifi.noise_level"
	FieldPrimaryPhysicalAdapterType = "primary_physical_adapter.type"
)

// =============================================================================
// Common Values / Enumerations
// =============================================================================

// Operating System Platforms
const (
	PlatformWindows = "Windows"
	PlatformMacOS   = "macOS"
	PlatformLinux   = "Linux"
)

// Hardware Types
const (
	HardwareTypeLaptop   = "laptop"
	HardwareTypeDesktop  = "desktop"
	HardwareTypeVirtual  = "virtual"
	HardwareTypeServer   = "server"
)

// User Types
const (
	UserTypeLocalUser  = "LOCAL_USER"
	UserTypeLocalAdmin = "LOCAL_ADMIN"
	UserTypeDomainUser = "DOMAIN_USER"
)

// Experience Levels (Web)
const (
	ExperienceLevelGood        = "good"
	ExperienceLevelFrustrating = "frustrating"
	ExperienceLevelTolerable   = "tolerable"
)

// Process Visibility
const (
	ProcessVisibilityForeground = "foreground"
	ProcessVisibilityBackground = "background"
)

// Connection Types
const (
	ConnectionTypeEthernet = "Ethernet"
	ConnectionTypeWifi     = "Wi-Fi"
	ConnectionTypeCellular = "Cellular"
)

// Collaboration Quality
const (
	QualityPoor = "poor"
	QualityGood = "good"
)

// DEX Application Score Node Types
const (
	NodeTypeApplication    = "application"
	NodeTypePageLoads      = "page_loads"
	NodeTypeTransactions   = "transactions"
	NodeTypeWebReliability = "web_reliability"
	NodeTypeCrashes        = "crashes"
	NodeTypeFreezes        = "freezes"
)

// Workflow/Remote Action Status
const (
	ExecutionStatusSuccess = "success"
	ExecutionStatusFailure = "failure"
	ExecutionStatusPending = "pending"
)

// Workflow Trigger Methods
const (
	TriggerMethodManual    = "manual"
	TriggerMethodScheduled = "scheduled"
	TriggerMethodAPI       = "api"
)

// Remote Action Purpose
const (
	PurposeRemediation = "remediation"
	PurposeDiagnostic  = "diagnostic"
)

// Boot Types
const (
	BootTypeFastStartup = "fast_startup"
	BootTypeColdBoot    = "cold_boot"
)

// Location Types
const (
	LocationTypeRemote = "Remote"
	LocationTypeOffice = "Office"
)
