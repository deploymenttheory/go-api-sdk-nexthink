package campaigns

import "fmt"

// ValidateTriggerRequest validates a campaign trigger request
func ValidateTriggerRequest(req *TriggerRequest) error {
	if req == nil {
		return fmt.Errorf("trigger request cannot be nil")
	}

	if req.CampaignNqlId == "" {
		return fmt.Errorf("campaignNqlId is required")
	}

	if len(req.UserSid) == 0 {
		return fmt.Errorf("userSid is required and must contain at least one SID")
	}

	if len(req.UserSid) > 10000 {
		return fmt.Errorf("userSid cannot contain more than 10000 SIDs (got %d)", len(req.UserSid))
	}

	for i, sid := range req.UserSid {
		if sid == "" {
			return fmt.Errorf("userSid[%d] cannot be empty", i)
		}
	}

	if req.ExpiresInMinutes < 1 {
		return fmt.Errorf("expiresInMinutes must be at least 1 (got %d)", req.ExpiresInMinutes)
	}

	if req.ExpiresInMinutes > 525600 {
		return fmt.Errorf("expiresInMinutes cannot exceed 525600 (got %d)", req.ExpiresInMinutes)
	}

	if len(req.Parameters) > 30 {
		return fmt.Errorf("parameters cannot contain more than 30 items (got %d)", len(req.Parameters))
	}

	return nil
}
