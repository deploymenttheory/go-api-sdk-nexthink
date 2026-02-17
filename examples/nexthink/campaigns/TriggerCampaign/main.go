package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/campaigns"
	"go.uber.org/zap"
)

// This example demonstrates how to trigger a campaign to be sent to users.
//
// Campaigns allow you to:
// - Send interactive surveys to users
// - Collect user feedback
// - Gather information about user experience
//
// Request includes:
// - CampaignNqlID: The NQL ID of the campaign to send
// - UserSID: List of user Security IDs (SIDs) - max 10000
// - ExpiresInMinutes: Time before campaign expires (1-525600 minutes)
// - Parameters: Optional parameters to customize campaign content
//
// Response includes:
// - For successful user requests: RequestId (to retrieve status and answers later)
// - For failed user requests: Message explaining the failure reason
//
// Duplicate SIDs in the request are automatically filtered out.
//
// Prerequisites:
// - Campaign must be pre-created in Nexthink admin
// - Campaign must be configured with a NQL ID
// - You need the Security IDs (SIDs) for target users
// - Your API credentials must have permission to trigger campaigns

func main() {
	clientID := os.Getenv("NEXTHINK_CLIENT_ID")
	clientSecret := os.Getenv("NEXTHINK_CLIENT_SECRET")
	instance := os.Getenv("NEXTHINK_INSTANCE")
	region := os.Getenv("NEXTHINK_REGION")

	if clientID == "" || clientSecret == "" || instance == "" || region == "" {
		log.Fatal("NEXTHINK_CLIENT_ID, NEXTHINK_CLIENT_SECRET, NEXTHINK_INSTANCE, and NEXTHINK_REGION environment variables are required")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	nxClient, err := nexthink.NewClient(
		clientID,
		clientSecret,
		instance,
		region,
		client.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Trigger a campaign
	request := &campaigns.TriggerRequest{
		CampaignNqlId: "#user_satisfaction_survey", // Replace with your campaign's NQL ID
		UserSid: []string{
			"S-1-5-21-1234567890-1234567890-1234567890-1001",
			"S-1-5-21-1234567890-1234567890-1234567890-1002",
			"S-1-5-21-1234567890-1234567890-1234567890-1003",
		},
		ExpiresInMinutes: 10080, // 7 days
		Parameters: map[string]string{
			"department":    "Engineering",
			"survey_period": "Q4 2026",
			"contact_email": "feedback@example.com",
		},
	}

	result, resp, err := nxClient.Campaigns.TriggerCampaign(ctx, request)
	if err != nil {
		log.Fatalf("Failed to trigger campaign: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Campaign Triggered ===\n")
	fmt.Printf("Total Users: %d\n", len(request.UserSid))
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	// Display per-user results
	successCount := 0
	failCount := 0

	fmt.Printf("User Request Results:\n")
	for i, userRequest := range result.Requests {
		fmt.Printf("\n%d. User SID: %s\n", i+1, userRequest.UserSid)

		if userRequest.RequestId != "" {
			fmt.Printf("   ✓ Success - Request ID: %s\n", userRequest.RequestId)
			successCount++
		} else {
			fmt.Printf("   ✗ Failed - %s\n", userRequest.Message)
			failCount++
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Successful: %d\n", successCount)
	fmt.Printf("  Failed: %d\n", failCount)

	logger.Info("Campaign triggered",
		zap.String("campaign_id", request.CampaignNqlId),
		zap.Int("total_users", len(request.UserSid)),
		zap.Int("successful", successCount),
		zap.Int("failed", failCount))

	fmt.Printf("\n✓ Campaign trigger completed!\n")
	fmt.Printf("\n💡 Tips:\n")
	fmt.Printf("   - Use RequestId to retrieve campaign status and answers via NQL\n")
	fmt.Printf("   - Campaign will expire after %d minutes if not responded to\n", request.ExpiresInMinutes)
	fmt.Printf("   - Parameters can be used to customize campaign content dynamically\n")
}
