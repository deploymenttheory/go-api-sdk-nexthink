package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
)

// This example demonstrates the most basic way to create a Nexthink client.
//
// IMPORTANT SECURITY NOTE:
// This example shows both environment variable (recommended) and hardcoded
// credential approaches. Always use environment variables in real code!
//
// Use this approach when:
// - You want the simplest possible client setup
// - You don't need custom configuration
// - You're getting started with the SDK
//
// The client uses sensible defaults:
// - 120 second timeout
// - 3 retries with exponential backoff
// - Production-level logging

func main() {
	// OPTION 1: From environment variables (RECOMMENDED)
	// This is the recommended approach - never hardcode credentials!
	clientID := os.Getenv("NEXTHINK_CLIENT_ID")
	clientSecret := os.Getenv("NEXTHINK_CLIENT_SECRET")
	instance := os.Getenv("NEXTHINK_INSTANCE")
	region := os.Getenv("NEXTHINK_REGION") // us, eu, etc.

	if clientID == "" || clientSecret == "" || instance == "" || region == "" {
		log.Fatal("NEXTHINK_CLIENT_ID, NEXTHINK_CLIENT_SECRET, NEXTHINK_INSTANCE, and NEXTHINK_REGION environment variables are required")
	}

	// OPTION 2: Hardcoded (NOT RECOMMENDED - for demonstration only)
	// Never do this in production! Only for local testing/learning.
	// clientID := "your-client-id-here"        // ⚠️ DON'T DO THIS IN REAL CODE!
	// clientSecret := "your-client-secret-here" // ⚠️ DON'T DO THIS IN REAL CODE!
	// instance := "your-instance-name"          // ⚠️ DON'T DO THIS IN REAL CODE!
	// region := "us"                            // ⚠️ DON'T DO THIS IN REAL CODE!

	// Create the simplest possible client
	client, err := nexthink.NewClient(clientID, clientSecret, instance, region)
	if err != nil {
		log.Fatalf("Failed to create Nexthink client: %v", err)
	}

	// Use the client to make a simple API call - list workflows
	ctx := context.Background()

	workflows, resp, err := client.Workflows.ListWorkflows(ctx)
	if err != nil {
		log.Fatalf("Failed to list workflows: %v", err)
	}

	// Display results
	fmt.Printf("✓ Client created successfully\n\n")
	fmt.Printf("Workflows:\n")
	fmt.Printf("  Count: %d\n", len(workflows))
	fmt.Printf("  Status Code: %d\n", resp.StatusCode)
	fmt.Printf("  Request Duration: %v\n", resp.Duration)

	if len(workflows) > 0 {
		fmt.Printf("\nFirst Workflow:\n")
		fmt.Printf("  ID: %s\n", workflows[0].ID)
		fmt.Printf("  Name: %s\n", workflows[0].Name)
		fmt.Printf("  Status: %s\n", workflows[0].Status)
	}

	fmt.Printf("\n✓ Basic client example completed successfully!\n")
	fmt.Printf("\n💡 Security Reminder:\n")
	fmt.Printf("   Always use environment variables for credentials\n")
	fmt.Printf("   Never hardcode credentials in your source code\n")
	fmt.Printf("   Add .env files to .gitignore\n")
}
