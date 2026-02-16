package enrichment

// EnrichmentRequest represents the request body for enriching fields
type EnrichmentRequest struct {
	// Enrichments is the list of enrichments (1-5000 items)
	Enrichments []Enrichment `json:"enrichments" validate:"required,min=1,max=5000,dive"`

	// Domain is the domain for which the enrichment applies (for information and tracking purposes)
	Domain string `json:"domain" validate:"required,min=1"`
}

// Enrichment represents a single enrichment operation
type Enrichment struct {
	// Identification is the list of fields to identify the object (exactly 1 item)
	Identification []Identification `json:"identification" validate:"required,len=1,dive"`

	// Fields is the list of fields to be enriched (min 1 item)
	Fields []Field `json:"fields" validate:"required,min=1,dive"`
}

// Identification represents the identification information for an object
type Identification struct {
	// Name is the field name used to identify the object
	// Valid values: device/device/name, device/device/uid, user/user/sid, user/user/uid,
	// user/user/upn, binary/binary/uid, package/package/uid
	Name string `json:"name" validate:"required"`

	// Value is the value used to identify the object
	Value string `json:"value" validate:"required"`
}

// Field represents a field to be enriched with its value
type Field struct {
	// Name is the name of the field to be enriched
	Name string `json:"name" validate:"required"`

	// Value is the desired value for the enrichment
	// Can be string, integer, or date
	Value any `json:"value" validate:"required"`
}

// =============================================================================
// Response Models
// =============================================================================

// SuccessResponse represents a successful enrichment response
type SuccessResponse struct {
	// Status is the status of the request ("success")
	Status string `json:"status"`
}

// PartialSuccessResponse represents a partial success response
type PartialSuccessResponse struct {
	// Status is the status of the request ("partial_success")
	Status string `json:"status"`

	// Errors contains the list of individual errors for objects with errors
	Errors []IndividualObjectError `json:"errors"`
}

// BadRequestResponse represents an error response when all objects contain errors
type BadRequestResponse struct {
	// Status is the status of the request ("error")
	Status string `json:"status"`

	// Errors contains the list of individual errors for all objects
	Errors []IndividualObjectError `json:"errors"`
}

// IndividualObjectError represents an error for a specific object
type IndividualObjectError struct {
	// Identification is the field used to identify the object
	Identification []Identification `json:"identification"`

	// Errors contains the list of errors for this object
	Errors []Error `json:"errors"`
}

// Error represents a single error with message and code
type Error struct {
	// Message is the descriptive error message
	Message string `json:"message"`

	// Code is the internal error code
	Code string `json:"code"`
}

// ForbiddenResponse represents a forbidden error response
type ForbiddenResponse struct {
	// Message is the error message when no permissions
	Message string `json:"message"`
}
