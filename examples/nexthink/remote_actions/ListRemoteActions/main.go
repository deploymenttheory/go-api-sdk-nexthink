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

// This example demonstrates how to list all remote actions with their configurations.
//
// Returns a list of all available remote actions including:
// - ID, UUID, Name, Description
// - Purpose (DATA_COLLECTION, REMEDIATION)
// - Targeting configuration (API/Manual/Workflow enabled)
// - Script information (inputs, outputs, timeout, run-as)
//
// Use this to:
// - Discover available remote actions in your Nexthink environment
// - Check which remote actions are API-enabled
// - View input/output parameters for remote actions

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

	// List all remote actions
	actions, resp, err := nxClient.RemoteActions.ListRemoteActions(ctx)
	if err != nil {
		log.Fatalf("Failed to list remote actions: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Available Remote Actions ===\n")
	fmt.Printf("Total: %d remote actions\n", len(actions))
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	// Count by purpose and targeting
	remediationCount := 0
	dataCollectionCount := 0
	apiEnabledCount := 0

	for _, action := range actions {
		for _, purpose := range action.Purpose {
			if purpose == "REMEDIATION" {
				remediationCount++
			}
			if purpose == "DATA_COLLECTION" {
				dataCollectionCount++
			}
		}
		if action.Targeting.APIEnabled {
			apiEnabledCount++
		}
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("  Remediation Actions: %d\n", remediationCount)
	fmt.Printf("  Data Collection Actions: %d\n", dataCollectionCount)
	fmt.Printf("  API-Enabled: %d\n\n", apiEnabledCount)

	// Display detailed information for API-enabled actions
	fmt.Printf("API-Enabled Remote Actions:\n")
	apiCount := 0
	for _, action := range actions {
		if !action.Targeting.APIEnabled {
			continue
		}
		apiCount++

		fmt.Printf("\n%d. %s\n", apiCount, action.Name)
		fmt.Printf("   ID: %s\n", action.ID)
		fmt.Printf("   UUID: %s\n", action.UUID)
		fmt.Printf("   Origin: %s\n", action.Origin)
		fmt.Printf("   Purpose: %v\n", action.Purpose)
		fmt.Printf("   Description: %s\n", action.Description)
		fmt.Printf("\n   Script Info:\n")
		fmt.Printf("     Run As: %s\n", action.ScriptInfo.RunAs)
		fmt.Printf("     Timeout: %d seconds\n", action.ScriptInfo.TimeoutSeconds)
		fmt.Printf("     Windows: %v\n", action.ScriptInfo.HasScriptWindows)
		fmt.Printf("     macOS: %v\n", action.ScriptInfo.HasScriptMacOS)
		fmt.Printf("     Inputs: %d\n", len(action.ScriptInfo.Inputs))
		fmt.Printf("     Outputs: %d\n", len(action.ScriptInfo.Outputs))

		// Show input parameters
		if len(action.ScriptInfo.Inputs) > 0 {
			fmt.Printf("\n   Input Parameters:\n")
			for _, input := range action.ScriptInfo.Inputs {
				fmt.Printf("     - %s: %s\n", input.Name, input.Description)
				if len(input.Options) > 0 {
					fmt.Printf("       Options: %v\n", input.Options)
				}
			}
		}
	}

	logger.Info("Remote actions listed successfully",
		zap.Int("total", len(actions)),
		zap.Int("api_enabled", apiEnabledCount))

	fmt.Printf("\n✓ Remote actions list retrieved successfully!\n")
}
