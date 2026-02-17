# Nexthink API Client Specification

## API Characteristics & Design Implications

### Authentication
- **Pattern**: OAuth2 Client Credentials Flow with separate token endpoint
- **Token Endpoint**: `https://{instance}-login.{region}.nexthink.cloud/oauth2/default/v1/token`
- **Token Lifetime**: 900 seconds (15 minutes)
- **Scope**: `service:integration`
- **Impact on Client**:
  - Requires `TokenManager` with automatic refresh logic
  - Needs middleware to inject current valid token into each request
  - Must handle token expiration during long-running operations
  - Refresh threshold set to 5 minutes before expiry to prevent mid-request expiration
  - Thread-safe token acquisition for concurrent requests

### Query Model
- **Pattern**: Pre-configured NQL queries referenced by ID
- **Query ID Format**: `#query_name` (must start with #)
- **Query Creation**: Queries must be created in Nexthink admin console first
- **No Ad-hoc Queries**: API does not accept arbitrary SQL/NQL strings
- **Impact on Client**:
  - Client-side validation for Query ID format (#-prefix)
  - No query builder for NQL strings
  - Query parameters validated against types (string, int, datetime)
  - Documentation emphasizes pre-configuration requirement

### Synchronous vs Asynchronous Patterns
- **Synchronous Execute (V1/V2)**:
  - Use case: Small result sets, high frequency, real-time queries
  - Returns results immediately in response body
  - V1: Returns headers array + data arrays
  - V2: Returns structured objects (map[string]any)
  
- **Asynchronous Export**:
  - Use case: Large data extracts, scheduled reports, bulk exports
  - Three-step process: StartExport → poll GetExportStatus → DownloadExport
  - Export status: SUBMITTED, IN_PROGRESS, COMPLETED, ERROR
  - Result delivered as S3 URL (external to Nexthink API)
  
- **Impact on Client**:
  - Two separate execution paths with different response types
  - Export requires polling helper (`WaitForExport`) with context timeout
  - Download uses standard `http.Client` (not SDK's transport) for S3 access
  - Status polling with exponential backoff to avoid API hammering

### Pagination
- **Pattern**: None for NQL results
- **Impact on Client**: No pagination helpers needed, results returned as complete sets or handled via export mechanism for large datasets

### Response Formats
- **Supported**: JSON, CSV
- **Content Negotiation**: Via `Accept` header
- **Impact on Client**:
  - Service methods specify Accept header per operation
  - CSV export returns byte stream for external processing

### Rate Limiting
- **Headers**: `X-Rate-Limit`, `X-Rate-Limit-Remaining`, `X-Rate-Limit-Reset`, `Retry-After`
- **Status Code**: 429 Too Many Requests
- **Impact on Client**:
  - Response wrapper exposes Headers for inspection
  - Helper function `GetRateLimitHeaders()` extracts rate limit info
  - Error type distinguishes rate limit errors
  - Retry logic respects `Retry-After` header

### Validation Requirements
- **UUID Format**: Workflows and campaigns require valid UUIDs
- **SID Format**: Windows Security IDs (S-1-5-...)
- **UPN Format**: User Principal Names (email format)
- **Query ID Format**: Must start with `#`
- **Batch Limits**:
  - Enrichment: 1-5000 operations
  - Workflows: max 10000 devices/users
  - Campaigns: max 10000 devices/users
  
- **Impact on Client**:
  - Regex validators for UUID, SID, UPN formats
  - Request validation before API calls to fail fast
  - Clear validation error messages for debugging

### Multi-Status Responses
- **Pattern**: Enrichment/Campaigns can return 207 Multi-Status
- **Scenarios**:
  - 200: All successful
  - 207: Partial success (some failed)
  - 400: All failed
  
- **Impact on Client**:
  - Response types handle multiple status codes
  - Partial success contains both successful and failed items
  - Error details per failed operation

### Resource Identification
- **Internal IDs**: Collector IDs (devices), SIDs (users)
- **External IDs**: Device names/UIDs, User UPNs/UIDs
- **NQL IDs**: `#query_name` format for queries, UUID for workflows
- **Impact on Client**:
  - V1 vs V2 endpoints for internal vs external IDs
  - Clear documentation on ID type requirements
  - Validation ensures correct ID format per endpoint

### Base URL Construction
- **Pattern**: `https://{instance}.api.{region}.nexthink.cloud/api/{version}/{service}`
- **Regions**: us, eu, pac, meta
- **Instance**: Customer-specific instance name
- **Impact on Client**:
  - Dynamic base URL construction from instance + region
  - Validation of region values
  - Token URL uses different subdomain: `{instance}-login.{region}.nexthink.cloud`

### Error Handling
- **Structured Errors**: API returns JSON error responses with details
- **Status Codes**:
  - 400: Validation failures
  - 401: Invalid authentication
  - 403: Permission denied
  - 404: Resource not found
  - 429: Rate limit exceeded
  - 500: Server error
  
- **Impact on Client**:
  - Custom error types per status code
  - Error responses parsed from JSON
  - Context preserved through error chain
  - Response object always returned (even on error) for header access
