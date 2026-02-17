package workflows

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateExecutionRequestV1 validates an ExecutionRequestV1
func ValidateExecutionRequestV1(req *ExecutionRequestV1) error {
	if req == nil {
		return fmt.Errorf("execution request cannot be nil")
	}

	if req.WorkflowID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}

	// At least one of devices or users must be provided
	if len(req.Devices) == 0 && len(req.Users) == 0 {
		return fmt.Errorf("at least one device or user must be provided")
	}

	// Validate devices if provided
	if len(req.Devices) > 0 {
		if len(req.Devices) > MaxDevices {
			return fmt.Errorf("maximum %d devices allowed", MaxDevices)
		}

		for i, device := range req.Devices {
			if device == "" {
				return fmt.Errorf("device at index %d cannot be empty", i)
			}
		}
	}

	// Validate users if provided
	if len(req.Users) > 0 {
		if len(req.Users) > MaxUsers {
			return fmt.Errorf("maximum %d users allowed", MaxUsers)
		}

		for i, user := range req.Users {
			if user == "" {
				return fmt.Errorf("user at index %d cannot be empty", i)
			}
		}
	}

	return nil
}

// ValidateExecutionRequestV2 validates an ExecutionRequestV2
func ValidateExecutionRequestV2(req *ExecutionRequestV2) error {
	if req == nil {
		return fmt.Errorf("execution request cannot be nil")
	}

	if req.WorkflowID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}

	// At least one of devices or users must be provided
	if len(req.Devices) == 0 && len(req.Users) == 0 {
		return fmt.Errorf("at least one device or user must be provided")
	}

	// Validate devices if provided
	if len(req.Devices) > 0 {
		if len(req.Devices) > MaxDevices {
			return fmt.Errorf("maximum %d devices allowed", MaxDevices)
		}

		for i, device := range req.Devices {
			if err := ValidateDeviceData(&device); err != nil {
				return fmt.Errorf("invalid device at index %d: %w", i, err)
			}
		}
	}

	// Validate users if provided
	if len(req.Users) > 0 {
		if len(req.Users) > MaxUsers {
			return fmt.Errorf("maximum %d users allowed", MaxUsers)
		}

		for i, user := range req.Users {
			if err := ValidateUserData(&user); err != nil {
				return fmt.Errorf("invalid user at index %d: %w", i, err)
			}
		}
	}

	return nil
}

// ValidateDeviceData validates DeviceData
func ValidateDeviceData(device *DeviceData) error {
	// At least one identifier must be provided
	if device.Name == "" && device.UID == "" && device.CollectorUID == "" {
		return fmt.Errorf("at least one device identifier (name, uid, or collectorUid) must be provided")
	}

	// Validate UID format if provided
	if device.UID != "" {
		if err := ValidateUUID(device.UID, "device UID"); err != nil {
			return err
		}
	}

	// Validate CollectorUID format if provided
	if device.CollectorUID != "" {
		if err := ValidateUUID(device.CollectorUID, "collector UID"); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUserData validates UserData
func ValidateUserData(user *UserData) error {
	// At least one identifier must be provided
	if user.UID == "" && user.UPN == "" && user.SID == "" {
		return fmt.Errorf("at least one user identifier (uid, upn, or sid) must be provided")
	}

	// Validate UID format if provided
	if user.UID != "" {
		if err := ValidateUUID(user.UID, "user UID"); err != nil {
			return err
		}
	}

	// Validate UPN format if provided (basic email validation)
	if user.UPN != "" {
		matched, err := regexp.MatchString(UserUPNPattern, user.UPN)
		if err != nil {
			return fmt.Errorf("failed to validate UPN format: %w", err)
		}
		if !matched {
			return fmt.Errorf("invalid UPN format (expected email format)")
		}
	}

	// Validate SID format if provided
	if user.SID != "" {
		matched, err := regexp.MatchString(UserSIDPattern, user.SID)
		if err != nil {
			return fmt.Errorf("failed to validate SID format: %w", err)
		}
		if !matched {
			return fmt.Errorf("invalid SID format (expected S-* format)")
		}
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

// ValidateUUID validates a UUID format
func ValidateUUID(uuid, fieldName string) error {
	if uuid == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	matched, err := regexp.MatchString(UUIDPattern, uuid)
	if err != nil {
		return fmt.Errorf("failed to validate %s UUID format: %w", fieldName, err)
	}

	if !matched {
		return fmt.Errorf("invalid %s UUID format", fieldName)
	}

	return nil
}
