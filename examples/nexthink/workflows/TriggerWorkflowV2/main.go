package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/workflows"
	"go.uber.org/zap"
)

// This example demonstrates how to trigger a workflow execution using V2 API with external identifiers.
//
// Trigger Workflow V2 uses external identifiers (more user-friendly than V1):
//
// For Devices, provide at least one of:
// - Name: The device name
// - UID: Globally unique device identifier (UUID)
// - CollectorUID: Nexthink Collector UUID
//
// For Users, provide at least one of:
// - SID: Security identifier
// - UPN: User principal name (email format)
// - UID: Globally unique user identifier (UUID)
//
// If multiple users or devices match the identifiers, the system triggers
// the workflow on the most recently active one (lastSeen).
//
// Prerequisites:
// - Workflow must be pre-created in Nexthink admin
// - Workflow must be API-enabled
// - Your API credentials must have permission to execute workflows

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

	// Trigger workflow with external identifiers (more user-friendly)
	request := &workflows.TriggerWorkflowV2Request{
		WorkflowID: "#your_workflow_id", // Must be pre-created in Nexthink
		Devices: []workflows.DeviceData{
			{
				Name: "DESKTOP-001",
				UID:  "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			},
			{
				Name: "LAPTOP-456",
			},
		},
		Users: []workflows.UserData{
			{
				UPN: "john.doe@example.com",
				SID: "S-1-5-21-1234567890-1234567890-1234567890-1001",
			},
			{
				UPN: "jane.smith@example.com",
			},
		},
		Params: map[string]string{
			"action":  "reset_password",
			"notify":  "true",
			"urgency": "high",
		},
	}

	result, resp, err := nxClient.Workflows.TriggerWorkflowV2(ctx, request)
	if err != nil {
		log.Fatalf("Failed to trigger workflow: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Workflow Triggered Successfully ===\n")
	fmt.Printf("Request UUID: %s\n", result.RequestUUID)
	fmt.Printf("Execution UUIDs: %d\n", len(result.ExecutionsUUIDs))
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	fmt.Printf("Individual Executions:\n")
	for i, execUUID := range result.ExecutionsUUIDs {
		fmt.Printf("  %d. %s\n", i+1, execUUID)
	}

	logger.Info("Workflow triggered successfully",
		zap.String("workflow_id", request.WorkflowID),
		zap.String("request_uuid", result.RequestUUID),
		zap.Int("executions", len(result.ExecutionsUUIDs)))

	fmt.Printf("\n✓ Workflow execution started!\n")
	fmt.Printf("\n💡 Advantages of V2:\n")
	fmt.Printf("   - Use human-readable identifiers (names, emails)\n")
	fmt.Printf("   - No need to lookup internal Collector IDs or SIDs\n")
	fmt.Printf("   - Automatically resolves to most recently active entity\n")
}
