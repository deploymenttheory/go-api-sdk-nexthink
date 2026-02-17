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

// This example demonstrates how to create a Nexthink client with custom logging.
//
// Use this approach when:
// - You need detailed logging for debugging
// - You want to integrate with your existing logging infrastructure
// - You need to track API calls and responses
//
// The example shows:
// - Creating a zap logger (production and development options)
// - Passing the logger to the client
// - Automatic logging of requests, responses, and errors

func main() {
	clientID := os.Getenv("NEXTHINK_CLIENT_ID")
	clientSecret := os.Getenv("NEXTHINK_CLIENT_SECRET")
	instance := os.Getenv("NEXTHINK_INSTANCE")
	region := os.Getenv("NEXTHINK_REGION")

	if clientID == "" || clientSecret == "" || instance == "" || region == "" {
		log.Fatal("NEXTHINK_CLIENT_ID, NEXTHINK_CLIENT_SECRET, NEXTHINK_INSTANCE, and NEXTHINK_REGION environment variables are required")
	}

	// OPTION 1: Production logger (JSON format, appropriate log levels)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create production logger: %v", err)
	}
	defer logger.Sync()

	// OPTION 2: Development logger (human-readable, more verbose)
	// Uncomment to use development logger instead:
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	//     log.Fatalf("Failed to create development logger: %v", err)
	// }
	// defer logger.Sync()

	// Create client with custom logger
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

	// Make an API call - the client will automatically log the request/response
	ctx := context.Background()

	workflows, _, err := nxClient.Workflows.ListWorkflows(ctx)
	if err != nil {
		log.Fatalf("Failed to list workflows: %v", err)
	}

	fmt.Printf("\n=== Workflows Retrieved ===\n")
	fmt.Printf("Total workflows: %d\n\n", len(workflows))

	for i, workflow := range workflows {
		if i >= 3 { // Show only first 3
			fmt.Printf("... and %d more\n", len(workflows)-3)
			break
		}
		fmt.Printf("%d. %s\n", i+1, workflow.Name)
		fmt.Printf("   ID: %s\n", workflow.ID)
		fmt.Printf("   Status: %s\n", workflow.Status)
		fmt.Printf("   Trigger Methods: %v\n\n", workflow.TriggerMethods)
	}

	// Structured logging example
	logger.Info("Workflows retrieved successfully",
		zap.Int("count", len(workflows)),
		zap.String("instance", instance),
		zap.String("region", region))

	fmt.Printf("✓ Client with logger example completed!\n")
	fmt.Printf("\n💡 Check the console output above for structured log entries\n")
}
