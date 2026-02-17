package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"go.uber.org/zap"
)

// This example demonstrates how to get detailed configuration for a specific workflow.
//
// Returns detailed configuration including:
// - Full workflow metadata (ID, UUID, Name, Description)
// - Status (ACTIVE/INACTIVE)
// - Available trigger methods
// - Version history
//
// Use this to:
// - Get detailed information about a specific workflow
// - Check workflow configuration before triggering
// - Verify trigger methods and versions

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

	// Get workflow details by NQL ID
	nqlID := "#your_workflow_id" // Replace with your workflow's NQL ID

	workflow, resp, err := nxClient.Workflows.GetWorkflowDetails(ctx, nqlID)
	if err != nil {
		log.Fatalf("Failed to get workflow details: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Workflow Details ===\n")
	fmt.Printf("ID: %s\n", workflow.ID)
	fmt.Printf("UUID: %s\n", workflow.UUID)
	fmt.Printf("Name: %s\n", workflow.Name)
	fmt.Printf("Description: %s\n", workflow.Description)
	fmt.Printf("Status: %s\n", workflow.Status)
	fmt.Printf("Last Update: %s\n", workflow.LastUpdateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	fmt.Printf("Trigger Methods:\n")
	for _, method := range workflow.TriggerMethods {
		fmt.Printf("  - %s\n", method)
	}

	fmt.Printf("\nVersions (%d):\n", len(workflow.Versions))
	for _, version := range workflow.Versions {
		activeStatus := ""
		if version.IsActive {
			activeStatus = " (ACTIVE)"
		}
		fmt.Printf("  - Version %d%s\n", version.VersionNumber, activeStatus)
	}

	logger.Info("Workflow details retrieved",
		zap.String("workflow_id", workflow.ID),
		zap.String("name", workflow.Name),
		zap.String("status", workflow.Status))

	fmt.Printf("\n✓ Workflow details retrieved successfully!\n")
}
