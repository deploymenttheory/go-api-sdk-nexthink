package campaigns

// TriggerRequest represents the request body for triggering a campaign
type TriggerRequest struct {
	// CampaignNqlId is the ID of the campaign to send
	CampaignNqlId string `json:"campaignNqlId" validate:"required,min=1"`

	// UserSid is the list of SIDs of users that the campaign should be sent to (1-10000 items)
	UserSid []string `json:"userSid" validate:"required,min=1,max=10000,dive,required"`

	// ExpiresInMinutes is the number of minutes before the campaign response expires (1-525600)
	// Starting from the current time. The expiration date is set at API call time.
	ExpiresInMinutes int `json:"expiresInMinutes" validate:"required,min=1,max=525600"`

	// Parameters are key-value pairs for parameters within the campaign to be replaced (max 30 items)
	// The provided keys must match exactly the IDs of all parameters of the campaign
	Parameters map[string]string `json:"parameters,omitempty" validate:"omitempty,max=30"`
}

// TriggerResponseDetails represents details of a campaign trigger request for a user
type TriggerResponseDetails struct {
	// RequestId is the ID of the request created for the user (if successful)
	RequestId string `json:"requestId,omitempty"`

	// UserSid is the SID of the user
	UserSid string `json:"userSid"`

	// Message is the reason why a request could not be created for that user SID (if failed)
	Message string `json:"message,omitempty"`
}

// TriggerSuccessResponse represents the successful response from triggering a campaign
type TriggerSuccessResponse struct {
	// Requests contains identifiers of the requests created for each user SID
	// Duplicate SIDs in the request are filtered out from the response list
	Requests []TriggerResponseDetails `json:"requests"`
}

// TriggerErrorResponse represents an error response from the Campaigns API
type TriggerErrorResponse struct {
	// Code is the error code returned to the client
	Code string `json:"code"`

	// Message is the error message returned to the client
	Message string `json:"message"`
}

