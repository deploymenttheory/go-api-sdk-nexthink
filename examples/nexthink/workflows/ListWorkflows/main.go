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

// This example demonstrates how to retrieve all workflows with their configurations.
//
// Returns a list of all workflows including:
// - ID, UUID, Name, Description, Status
// - Available trigger methods (API, Manual, Scheduler, etc.)
// - Version information
//
// Use this to:
// - Discover available workflows in your Nexthink environment
// - Check which workflows are API-enabled
// - Verify workflow status (ACTIVE/INACTIVE)

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

	// List all workflows
	workflows, resp, err := nxClient.Workflows.ListWorkflows(ctx)
	if err != nil {
		log.Fatalf("Failed to list workflows: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Available Workflows ===\n")
	fmt.Printf("Total: %d workflows\n", len(workflows))
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	// Count by status
	activeCount := 0
	apiEnabledCount := 0

	for _, workflow := range workflows {
		if workflow.Status == "ACTIVE" {
			activeCount++
		}
		for _, method := range workflow.TriggerMethods {
			if method == "API" {
				apiEnabledCount++
				break
			}
		}
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("  Active: %d\n", activeCount)
	fmt.Printf("  API-Enabled: %d\n\n", apiEnabledCount)

	// Display detailed information for each workflow
	fmt.Printf("Workflows:\n")
	for i, workflow := range workflows {
		fmt.Printf("\n%d. %s\n", i+1, workflow.Name)
		fmt.Printf("   ID: %s\n", workflow.ID)
		fmt.Printf("   UUID: %s\n", workflow.UUID)
		fmt.Printf("   Status: %s\n", workflow.Status)
		fmt.Printf("   Description: %s\n", workflow.Description)
		fmt.Printf("   Trigger Methods: %v\n", workflow.TriggerMethods)
		fmt.Printf("   Last Update: %s\n", workflow.LastUpdateTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Versions: %d\n", len(workflow.Versions))
	}

	logger.Info("Workflows listed successfully",
		zap.Int("total", len(workflows)),
		zap.Int("active", activeCount),
		zap.Int("api_enabled", apiEnabledCount))

	fmt.Printf("\n✓ Workflow list retrieved successfully!\n")
}
