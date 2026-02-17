package remote_actions

import (
	"fmt"
	"strings"
)

// ValidateTriggerRemoteActionRequest validates a TriggerRemoteActionRequest
func ValidateTriggerRemoteActionRequest(req *TriggerRemoteActionRequest) error {
	if req == nil {
		return fmt.Errorf("execution request cannot be nil")
	}

	if req.RemoteActionID == "" {
		return fmt.Errorf("remote action ID cannot be empty")
	}

	if len(req.Devices) < MinDevices {
		return fmt.Errorf("at least %d device is required", MinDevices)
	}

	if len(req.Devices) > MaxDevices {
		return fmt.Errorf("maximum %d devices allowed", MaxDevices)
	}

	// Validate each device ID is not empty
	for i, device := range req.Devices {
		if device == "" {
			return fmt.Errorf("device at index %d cannot be empty", i)
		}
	}

	// Validate expiresInMinutes if provided
	if req.ExpiresInMinutes != 0 {
		if req.ExpiresInMinutes < MinExpiresInMinutes {
			return fmt.Errorf("expiresInMinutes must be at least %d", MinExpiresInMinutes)
		}
		if req.ExpiresInMinutes > MaxExpiresInMinutes {
			return fmt.Errorf("expiresInMinutes cannot exceed %d", MaxExpiresInMinutes)
		}
	}

	// Validate TriggerInfo if provided
	if req.TriggerInfo != nil {
		if err := ValidateTriggerInfo(req.TriggerInfo); err != nil {
			return fmt.Errorf("invalid trigger info: %w", err)
		}
	}

	return nil
}

// ValidateTriggerInfo validates TriggerInfoRequest
func ValidateTriggerInfo(info *TriggerInfoRequest) error {
	if info.Reason != "" && len(info.Reason) > MaxReasonLength {
		return fmt.Errorf("reason cannot exceed %d characters", MaxReasonLength)
	}

	return nil
}

// ValidateNqlID validates an NQL ID format
func ValidateNqlID(nqlID string) error {
	if nqlID == "" {
		return fmt.Errorf("NQL ID cannot be empty")
	}

	if !strings.HasPrefix(nqlID, "#") {
		return fmt.Errorf("NQL ID must start with #")
	}

	if len(nqlID) < 2 {
		return fmt.Errorf("NQL ID must be at least 2 characters")
	}

	return nil
}
