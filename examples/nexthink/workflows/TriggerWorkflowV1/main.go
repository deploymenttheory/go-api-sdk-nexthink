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

// This example demonstrates how to trigger a workflow execution using V1 API with internal IDs.
//
// Trigger Workflow V1 uses Nexthink internal identifiers:
// - Devices: Nexthink Collector IDs (UUID format) - max 10000
// - Users: Security IDs (SID format) - max 10000
//
// Returns RequestUUID and ExecutionsUUIDs to track execution via NQL queries.
//
// Prerequisites:
// - Workflow must be pre-created in Nexthink admin
// - Workflow must be API-enabled
// - You need the Collector IDs (for devices) or SIDs (for users)
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

	// Trigger workflow with internal IDs
	request := &workflows.TriggerWorkflowV1Request{
		WorkflowID: "#your_workflow_id", // Must be pre-created in Nexthink
		Devices: []string{
			"a1b2c3d4-e5f6-7890-abcd-ef1234567890", // Nexthink Collector ID
			"b2c3d4e5-f6a7-8901-bcde-f12345678901",
		},
		Users: []string{
			"S-1-5-21-1234567890-1234567890-1234567890-1001", // Security ID (SID)
		},
		Params: map[string]string{
			"reason": "Scheduled maintenance",
		},
	}

	result, resp, err := nxClient.Workflows.TriggerWorkflowV1(ctx, request)
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
	fmt.Printf("\n💡 Track execution status in NQL using:\n")
	fmt.Printf("   workflow.executions.request_id = '%s'\n", result.RequestUUID)
	fmt.Printf("   workflow.executions.execution_id in (%v)\n", result.ExecutionsUUIDs)
}
