package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
	"go.uber.org/zap"
)

// This example demonstrates the complete NQL export workflow:
// 1. Start an export
// 2. Wait for it to complete (with automatic polling)
// 3. Download the results
//
// WaitForNQLExport is a convenience method that:
// - Polls GetNQLExportStatus() at regular intervals
// - Returns when export reaches COMPLETED or ERROR status
// - Handles timeouts gracefully
//
// Recommended settings:
// - Poll interval: 5-10 seconds
// - Timeout: 5-10 minutes

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

	// Step 1: Start the export
	fmt.Printf("\n=== Starting NQL Export ===\n")
	exportRequest := &nql.ExportRequest{
		QueryID: "#your_large_query_id", // Must be pre-created in Nexthink
		Format:  nql.ExportFormatCSV,
	}

	startResult, _, err := nxClient.NQL.StartNQLExport(ctx, exportRequest)
	if err != nil {
		log.Fatalf("Failed to start export: %v", err)
	}

	fmt.Printf("Export ID: %s\n", startResult.ExportID)
	fmt.Printf("Initial Status: %s\n", startResult.Status)

	// Step 2: Wait for export to complete
	fmt.Printf("\n=== Waiting for Export to Complete ===\n")
	fmt.Printf("Polling every 5 seconds with 10 minute timeout...\n\n")

	pollInterval := 5 * time.Second
	timeout := 10 * time.Minute

	statusResult, err := nxClient.NQL.WaitForNQLExport(ctx, startResult.ExportID, pollInterval, timeout)
	if err != nil {
		log.Fatalf("Failed to wait for export: %v", err)
	}

	fmt.Printf("Final Status: %s\n", statusResult.Status)

	if statusResult.Status == "COMPLETED" {
		fmt.Printf("✓ Export completed successfully!\n")
		fmt.Printf("Download URL: %s\n\n", statusResult.ResultsFileURL)

		// Step 3: Download the results
		fmt.Printf("=== Downloading Export Results ===\n")
		data, err := nxClient.NQL.DownloadNQLExport(ctx, statusResult.ResultsFileURL)
		if err != nil {
			log.Fatalf("Failed to download export: %v", err)
		}

		fmt.Printf("Downloaded: %d bytes\n", len(data))

		// Optionally save to file
		outputFile := "nql_export_results.csv"
		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			log.Fatalf("Failed to save export to file: %v", err)
		}

		fmt.Printf("Saved to: %s\n", outputFile)

		logger.Info("NQL export completed and downloaded",
			zap.String("export_id", startResult.ExportID),
			zap.Int("size_bytes", len(data)))

		fmt.Printf("\n✓ Complete NQL export workflow finished successfully!\n")
	} else if statusResult.Status == "ERROR" {
		fmt.Printf("✗ Export failed\n")
		fmt.Printf("Error: %s\n", statusResult.ErrorDescription)
		log.Fatalf("Export failed")
	}
}
