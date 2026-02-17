package campaigns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateTriggerRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *TriggerRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1440,
			},
			wantErr: false,
		},
		{
			name: "valid request with parameters",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1440,
				Parameters: map[string]string{
					"department": "IT",
					"location":   "New York",
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "trigger request cannot be nil",
		},
		{
			name: "empty campaign NQL ID",
			req: &TriggerRequest{
				CampaignNqlId: "",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1440,
			},
			wantErr: true,
			errMsg:  "campaignNqlId is required",
		},
		{
			name: "empty userSid list",
			req: &TriggerRequest{
				CampaignNqlId:    "#security_awareness",
				UserSid:          []string{},
				ExpiresInMinutes: 1440,
			},
			wantErr: true,
			errMsg:  "userSid is required",
		},
		{
			name: "too many userSids",
			req: &TriggerRequest{
				CampaignNqlId:    "#security_awareness",
				UserSid:          make([]string, 10001),
				ExpiresInMinutes: 1440,
			},
			wantErr: true,
			errMsg:  "userSid cannot contain more than 10000 SIDs",
		},
		{
			name: "empty userSid in list",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
					"",
				},
				ExpiresInMinutes: 1440,
			},
			wantErr: true,
			errMsg:  "userSid[1] cannot be empty",
		},
		{
			name: "expiresInMinutes zero",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 0,
			},
			wantErr: true,
			errMsg:  "expiresInMinutes must be at least 1",
		},
		{
			name: "expiresInMinutes negative",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: -1,
			},
			wantErr: true,
			errMsg:  "expiresInMinutes must be at least 1",
		},
		{
			name: "expiresInMinutes too high",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 525601,
			},
			wantErr: true,
			errMsg:  "expiresInMinutes cannot exceed 525600",
		},
		{
			name: "too many parameters",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1440,
				Parameters:       make(map[string]string, 31),
			},
			wantErr: true,
			errMsg:  "parameters cannot contain more than 30 items",
		},
		{
			name: "min expiresInMinutes",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 1,
			},
			wantErr: false,
		},
		{
			name: "max expiresInMinutes",
			req: &TriggerRequest{
				CampaignNqlId: "#security_awareness",
				UserSid: []string{
					"S-1-5-21-1234567890-1234567890-1234567890-1001",
				},
				ExpiresInMinutes: 525600,
			},
			wantErr: false,
		},
	}

	// Fill the too many parameters test case
	for _, tt := range tests {
		if tt.name == "too many parameters" {
			for i := 0; i < 31; i++ {
				tt.req.Parameters[string(rune('a'+i))] = "value"
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerRequest(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
