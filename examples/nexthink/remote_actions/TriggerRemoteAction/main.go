package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/remote_actions"
	"go.uber.org/zap"
)

// This example demonstrates how to trigger a remote action execution on devices.
//
// Remote Actions allow you to:
// - Execute remediation scripts on devices
// - Collect data from devices
// - Perform system maintenance tasks
//
// Request includes:
// - RemoteActionID: The NQL ID of the remote action to execute
// - Devices: List of Nexthink Collector IDs (1-10000)
// - Params: Optional script parameters as key-value pairs
// - ExpiresInMinutes: Expiration time if device doesn't come online (60-10080 minutes)
// - TriggerInfo: Optional metadata (external source, reason, reference)
//
// Returns a RequestID that can be used to query remote action executions in NQL.
//
// Prerequisites:
// - Remote action must be pre-created in Nexthink admin
// - Remote action must be API-enabled
// - You need the Collector IDs for target devices
// - Your API credentials must have permission to execute remote actions

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

	// Trigger a remote action
	request := &remote_actions.TriggerRemoteActionRequest{
		RemoteActionID: "#clear_browser_cache", // Replace with your remote action's NQL ID
		Devices: []string{
			"a1b2c3d4-e5f6-7890-abcd-ef1234567890", // Nexthink Collector ID
			"b2c3d4e5-f6a7-8901-bcde-f12345678901",
		},
		ExpiresInMinutes: 1440, // 24 hours
		Params: map[string]string{
			"browser":    "chrome",
			"clear_data": "cookies",
		},
		TriggerInfo: &remote_actions.TriggerInfoRequest{
			ExternalSource:    "ServiceDesk",
			Reason:            "User reported slow browser performance",
			ExternalReference: "TICKET-12345",
		},
	}

	result, resp, err := nxClient.RemoteActions.TriggerRemoteAction(ctx, request)
	if err != nil {
		log.Fatalf("Failed to trigger remote action: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Remote Action Triggered ===\n")
	fmt.Printf("Request ID: %s\n", result.RequestID)
	fmt.Printf("Expires In: %d minutes\n", result.ExpiresInMinutes)
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n", resp.Duration)

	logger.Info("Remote action triggered successfully",
		zap.String("remote_action_id", request.RemoteActionID),
		zap.String("request_id", result.RequestID),
		zap.Int("devices", len(request.Devices)))

	fmt.Printf("\n✓ Remote action execution initiated!\n")
	fmt.Printf("\n💡 Track execution status in NQL using:\n")
	fmt.Printf("   act.remote_actions.request_id = '%s'\n", result.RequestID)
	fmt.Printf("\n   The action will execute when the devices come online\n")
	fmt.Printf("   and will expire after %d minutes if devices remain offline\n", result.ExpiresInMinutes)
}
