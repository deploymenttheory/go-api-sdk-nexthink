package remote_actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateExecutionRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *TriggerRemoteActionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &TriggerRemoteActionRequest{
				RemoteActionID:   "#test_action",
				Devices:          []string{"device-001"},
				ExpiresInMinutes: 1440,
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
			name: "empty remote action ID",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "",
				Devices:        []string{"device-001"},
			},
			wantErr: true,
			errMsg:  "remote action ID cannot be empty",
		},
		{
			name: "no devices",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "#test_action",
				Devices:        []string{},
			},
			wantErr: true,
			errMsg:  "at least 1 device is required",
		},
		{
			name: "empty device ID",
			req: &TriggerRemoteActionRequest{
				RemoteActionID: "#test_action",
				Devices:        []string{"device-001", "", "device-003"},
			},
			wantErr: true,
			errMsg:  "device at index 1 cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerRemoteActionRequest(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateTriggerInfo(t *testing.T) {
	tests := []struct {
		name    string
		info    *TriggerInfoRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid trigger info",
			info: &TriggerInfoRequest{
				ExternalSource:    "ServiceDesk",
				Reason:            "Test reason",
				ExternalReference: "TICKET-123",
			},
			wantErr: false,
		},
		{
			name: "reason too long",
			info: &TriggerInfoRequest{
				Reason: string(make([]byte, 501)),
			},
			wantErr: true,
			errMsg:  "reason cannot exceed 500 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerInfo(tt.info)
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
			nqlID:   "#test_action",
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
			nqlID:   "test_action",
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
