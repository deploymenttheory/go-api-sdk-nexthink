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

// This example demonstrates how to get detailed configuration for a specific remote action.
//
// Returns detailed configuration including:
// - Full remote action metadata
// - Targeting configuration (API/Manual/Workflow enabled)
// - Script details (platform support, inputs, outputs)
// - Execution settings (run-as, timeout)
//
// Use this to:
// - Get full configuration details before triggering
// - Check input/output parameters
// - Verify platform support and execution context

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

	// Get remote action details by NQL ID
	nqlID := "#clear_browser_cache" // Replace with your remote action's NQL ID

	action, resp, err := nxClient.RemoteActions.GetRemoteActionDetails(ctx, nqlID)
	if err != nil {
		log.Fatalf("Failed to get remote action details: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Remote Action Details ===\n")
	fmt.Printf("ID: %s\n", action.ID)
	fmt.Printf("UUID: %s\n", action.UUID)
	fmt.Printf("Name: %s\n", action.Name)
	fmt.Printf("Description: %s\n", action.Description)
	fmt.Printf("Origin: %s\n", action.Origin)
	fmt.Printf("Built-in Content Version: %s\n", action.BuiltInContentVersion)
	fmt.Printf("Purpose: %v\n", action.Purpose)
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	fmt.Printf("Targeting Configuration:\n")
	fmt.Printf("  API Enabled: %v\n", action.Targeting.APIEnabled)
	fmt.Printf("  Manual Enabled: %v\n", action.Targeting.ManualEnabled)
	fmt.Printf("  Workflow Enabled: %v\n", action.Targeting.WorkflowEnabled)
	fmt.Printf("  Manual Multiple Devices: %v\n\n", action.Targeting.ManualAllowMultipleDevices)

	fmt.Printf("Script Information:\n")
	fmt.Printf("  Run As: %s\n", action.ScriptInfo.RunAs)
	fmt.Printf("  Timeout: %d seconds\n", action.ScriptInfo.TimeoutSeconds)
	fmt.Printf("  Windows Script: %v\n", action.ScriptInfo.HasScriptWindows)
	fmt.Printf("  macOS Script: %v\n", action.ScriptInfo.HasScriptMacOS)
	fmt.Printf("  Execution Service: %s\n\n", action.ScriptInfo.ExecutionServiceDelegate)

	// Display input parameters
	if len(action.ScriptInfo.Inputs) > 0 {
		fmt.Printf("Input Parameters (%d):\n", len(action.ScriptInfo.Inputs))
		for i, input := range action.ScriptInfo.Inputs {
			fmt.Printf("\n  %d. %s\n", i+1, input.Name)
			fmt.Printf("     ID: %s\n", input.ID)
			fmt.Printf("     Description: %s\n", input.Description)
			fmt.Printf("     Windows: %v | macOS: %v\n", input.UsedByWindows, input.UsedByMacOS)
			fmt.Printf("     Allow Custom: %v\n", input.AllowCustomValue)
			if len(input.Options) > 0 {
				fmt.Printf("     Options: %v\n", input.Options)
			}
		}
		fmt.Printf("\n")
	}

	// Display output parameters
	if len(action.ScriptInfo.Outputs) > 0 {
		fmt.Printf("Output Parameters (%d):\n", len(action.ScriptInfo.Outputs))
		for i, output := range action.ScriptInfo.Outputs {
			fmt.Printf("\n  %d. %s\n", i+1, output.Name)
			fmt.Printf("     ID: %s\n", output.ID)
			fmt.Printf("     Type: %s\n", output.Type)
			fmt.Printf("     Description: %s\n", output.Description)
			fmt.Printf("     Windows: %v | macOS: %v\n", output.UsedByWindows, output.UsedByMacOS)
		}
		fmt.Printf("\n")
	}

	logger.Info("Remote action details retrieved",
		zap.String("action_id", action.ID),
		zap.String("name", action.Name),
		zap.Bool("api_enabled", action.Targeting.APIEnabled))

	fmt.Printf("✓ Remote action details retrieved successfully!\n")
}
