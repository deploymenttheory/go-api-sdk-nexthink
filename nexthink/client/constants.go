package client

const (
	// DefaultTokenURL is the OAuth2 token endpoint template
	// Replace {instance} and {region} in runtime
	DefaultTokenURLTemplate = "https://%s-login.%s.nexthink.cloud/oauth2/default/v1/token"

	// DefaultAPIVersion is the API version
	DefaultAPIVersion = "v1"

	// UserAgentBase is the base user agent string prefix
	UserAgentBase = "go-api-sdk-nexthink"

	// DefaultTimeout is the default HTTP client timeout in seconds
	DefaultTimeout = 120

	// MaxRetries is the maximum number of retries for failed requests
	MaxRetries = 3

	// RetryWaitTime is the wait time between retries in seconds
	RetryWaitTime = 2

	// RetryMaxWaitTime is the maximum wait time between retries in seconds
	RetryMaxWaitTime = 10

	// TokenLifetime is the token lifetime in seconds (15 minutes)
	TokenLifetime = 900

	// TokenRefreshBuffer is the buffer time before token expiry to refresh (2 minutes)
	TokenRefreshBuffer = 120
)

// OAuth2 constants
const (
	GrantTypeClientCredentials = "client_credentials"
	ScopeServiceIntegration    = "service:integration"
)

// Response format constants
const (
	FormatJSON = "json"
)

// HTTP headers
const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	AcceptJSON                = "application/json"
)

// Region constants
const (
	RegionUS   = "us"   // United States
	RegionEU   = "eu"   // European Union
	RegionPAC  = "pac"  // Asia-Pacific
	RegionMETA = "meta" // Middle East, Turkey and Africa
)
