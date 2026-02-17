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

// This example demonstrates how to check the status of an NQL export.
//
// Export Status Values:
// - SUBMITTED: Export is queued
// - IN_PROGRESS: Export is currently running
// - COMPLETED: Export is ready (ResultsFileURL will be available)
// - ERROR: Export failed (ErrorDescription will contain error details)
//
// When status is COMPLETED, the response includes a ResultsFileURL (S3 URL)
// that can be used with DownloadNQLExport().

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

	// Check export status
	// Replace with the export ID from StartNQLExport()
	exportID := "your-export-id-here"

	result, resp, err := nxClient.NQL.GetNQLExportStatus(ctx, exportID)
	if err != nil {
		log.Fatalf("Failed to get export status: %v", err)
	}

	// Display results
	fmt.Printf("\n=== NQL Export Status ===\n")
	fmt.Printf("Export ID: %s\n", exportID)
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	// Display status-specific information
	switch result.Status {
	case "COMPLETED":
		fmt.Printf("✓ Export completed!\n")
		fmt.Printf("Download URL: %s\n", result.ResultsFileURL)
		fmt.Printf("\n💡 Use DownloadNQLExport() to download the results\n")

	case "ERROR":
		fmt.Printf("✗ Export failed\n")
		fmt.Printf("Error: %s\n", result.ErrorDescription)

	case "SUBMITTED", "IN_PROGRESS":
		fmt.Printf("⏳ Export still processing...\n")
		fmt.Printf("\n💡 Poll again in a few seconds or use WaitForNQLExport()\n")

	default:
		fmt.Printf("Unknown status: %s\n", result.Status)
	}

	logger.Info("Export status checked",
		zap.String("export_id", exportID),
		zap.String("status", result.Status))

	fmt.Printf("\n✓ Status check completed!\n")
}
