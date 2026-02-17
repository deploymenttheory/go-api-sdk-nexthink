package workflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateTriggerWorkflowV1Request(t *testing.T) {
	tests := []struct {
		name    string
		req     *TriggerWorkflowV1Request
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with devices",
			req: &TriggerWorkflowV1Request{
				WorkflowID: "#test_workflow",
				Devices:    []string{"device-001"},
			},
			wantErr: false,
		},
		{
			name: "valid request with users",
			req: &TriggerWorkflowV1Request{
				WorkflowID: "#test_workflow",
				Users:      []string{"S-1-5-21-1234567890-1234567890-1234567890-1001"},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "execution request cannot be nil",
		},
		{
			name: "empty workflow ID",
			req: &TriggerWorkflowV1Request{
				WorkflowID: "",
				Devices:    []string{"device-001"},
			},
			wantErr: true,
			errMsg:  "workflow ID cannot be empty",
		},
		{
			name: "no devices or users",
			req: &TriggerWorkflowV1Request{
				WorkflowID: "#test_workflow",
			},
			wantErr: true,
			errMsg:  "at least one device or user must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerWorkflowV1Request(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateTriggerWorkflowV2Request(t *testing.T) {
	tests := []struct {
		name    string
		req     *TriggerWorkflowV2Request
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with devices",
			req: &TriggerWorkflowV2Request{
				WorkflowID: "#test_workflow",
				Devices: []DeviceData{
					{Name: "DESKTOP-001"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid request with users",
			req: &TriggerWorkflowV2Request{
				WorkflowID: "#test_workflow",
				Users: []UserData{
					{SID: "S-1-5-21-1234567890-1234567890-1234567890-1001"},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "execution request cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerWorkflowV2Request(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateDeviceData(t *testing.T) {
	tests := []struct {
		name    string
		device  *DeviceData
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid device with name",
			device:  &DeviceData{Name: "DESKTOP-001"},
			wantErr: false,
		},
		{
			name:    "valid device with UID",
			device:  &DeviceData{UID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"},
			wantErr: false,
		},
		{
			name:    "valid device with collector UID",
			device:  &DeviceData{CollectorUID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"},
			wantErr: false,
		},
		{
			name:    "device with no identifiers",
			device:  &DeviceData{},
			wantErr: true,
			errMsg:  "at least one device identifier",
		},
		{
			name:    "invalid UID format",
			device:  &DeviceData{UID: "invalid-uuid"},
			wantErr: true,
			errMsg:  "invalid device UID UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDeviceData(tt.device)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserData(t *testing.T) {
	tests := []struct {
		name    string
		user    *UserData
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid user with UID",
			user:    &UserData{UID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"},
			wantErr: false,
		},
		{
			name:    "valid user with UPN",
			user:    &UserData{UPN: "user@example.com"},
			wantErr: false,
		},
		{
			name:    "valid user with SID",
			user:    &UserData{SID: "S-1-5-21-1234567890-1234567890-1234567890-1001"},
			wantErr: false,
		},
		{
			name:    "user with no identifiers",
			user:    &UserData{},
			wantErr: true,
			errMsg:  "at least one user identifier",
		},
		{
			name:    "invalid UID format",
			user:    &UserData{UID: "invalid-uuid"},
			wantErr: true,
			errMsg:  "invalid user UID UUID format",
		},
		{
			name:    "invalid UPN format",
			user:    &UserData{UPN: "not-an-email"},
			wantErr: true,
			errMsg:  "invalid UPN format",
		},
		{
			name:    "invalid SID format",
			user:    &UserData{SID: "invalid-sid"},
			wantErr: true,
			errMsg:  "invalid SID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserData(tt.user)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateNqlID(t *testing.T) {
	tests := []struct {
		name    string
		nqlID   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid NQL ID",
			nqlID:   "#test_workflow",
			wantErr: false,
		},
		{
			name:    "empty NQL ID",
			nqlID:   "",
			wantErr: true,
			errMsg:  "NQL ID cannot be empty",
		},
		{
			name:    "missing hash prefix",
			nqlID:   "test_workflow",
			wantErr: true,
			errMsg:  "NQL ID must start with #",
		},
		{
			name:    "only hash",
			nqlID:   "#",
			wantErr: true,
			errMsg:  "NQL ID must be at least 2 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNqlID(tt.nqlID)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name      string
		uuid      string
		fieldName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid UUID",
			uuid:      "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			fieldName: "test",
			wantErr:   false,
		},
		{
			name:      "empty UUID",
			uuid:      "",
			fieldName: "test",
			wantErr:   true,
			errMsg:    "test cannot be empty",
		},
		{
			name:      "invalid UUID format",
			uuid:      "invalid-uuid",
			fieldName: "test",
			wantErr:   true,
			errMsg:    "invalid test UUID format",
		},
		{
			name:      "UUID without dashes",
			uuid:      "a1b2c3d4e5f67890abcdef1234567890",
			fieldName: "test",
			wantErr:   true,
			errMsg:    "invalid test UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.uuid, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
