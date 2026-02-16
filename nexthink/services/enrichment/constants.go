package enrichment

// EndpointEnrichmentDataFields is the API endpoint for enriching fields
const EndpointEnrichmentDataFields = "/api/v1/enrichment/data/fields"

// Identification field names
const (
	IdentificationDeviceName  = "device/device/name"
	IdentificationDeviceUID   = "device/device/uid"
	IdentificationUserSID     = "user/user/sid"
	IdentificationUserUID     = "user/user/uid"
	IdentificationUserUPN     = "user/user/upn"
	IdentificationBinaryUID   = "binary/binary/uid"
	IdentificationPackageUID  = "package/package/uid"
)

// Device enrichment field names
const (
	FieldDeviceConfigurationTag           = "device/device/configuration_tag"
	FieldDeviceVirtualizationDesktopBroker = "device/device/virtualization/desktop_broker"
	FieldDeviceVirtualizationDesktopPool  = "device/device/virtualization/desktop_pool"
	FieldDeviceVirtualizationDiskImage    = "device/device/virtualization/disk_image"
	FieldDeviceVirtualizationEnvironment  = "device/device/virtualization/environment_name"
	FieldDeviceVirtualizationHostname     = "device/device/virtualization/hostname"
	FieldDeviceVirtualizationHypervisor   = "device/device/virtualization/hypervisor_name"
	FieldDeviceVirtualizationInstanceSize = "device/device/virtualization/instance_size"
	FieldDeviceVirtualizationLastUpdate   = "device/device/virtualization/last_update"
	FieldDeviceVirtualizationRegion       = "device/device/virtualization/region"
	FieldDeviceVirtualizationType         = "device/device/virtualization/type"
)

// User Entra ID (AD) enrichment field names
const (
	FieldUserADCity                = "user/user/ad/city"
	FieldUserADCountryCode         = "user/user/ad/country_code"
	FieldUserADDepartment          = "user/user/ad/department"
	FieldUserADDistinguishedName   = "user/user/ad/distinguished_name"
	FieldUserADEmailAddress        = "user/user/ad/email_address"
	FieldUserADFullName            = "user/user/ad/full_name"
	FieldUserADJobTitle            = "user/user/ad/job_title"
	FieldUserADLastUpdate          = "user/user/ad/last_update"
	FieldUserADOffice              = "user/user/ad/office"
	FieldUserADOrganizationalUnit  = "user/user/ad/organizational_unit"
	FieldUserADUsername            = "user/user/ad/username"
)

// Response status values
const (
	StatusSuccess        = "success"
	StatusPartialSuccess = "partial_success"
	StatusError          = "error"
)
