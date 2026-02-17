package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
	"go.uber.org/zap"
)

// This example demonstrates how to start an asynchronous NQL export.
//
// Export NQL is optimized for:
// - Large queries at low frequency
// - Large data extracts
// - Scheduled reports
// - Bulk data exports
//
// This is an asynchronous operation that:
// 1. Returns an exportID immediately
// 2. Processes the query in the background
// 3. Provides a download URL when complete
//
// Use GetNQLExportStatus() to check completion status.
//
// Prerequisites:
// - NQL query must be pre-created in Nexthink admin (Content Management > NQL API queries)
// - Query must have a Query ID (format: #query_name)
// - Your API credentials must have permission to execute NQL queries

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

	// Start an NQL export
	// Replace with your actual query ID from Nexthink admin
	request := &nql.ExportRequest{
		QueryID: "#your_large_query_id", // Must be pre-created in Nexthink
		Format:  nql.ExportFormatCSV,    // or nql.ExportFormatJSON
		// Platform: "windows", // Optional: filter by platform
	}

	result, resp, err := nxClient.NQL.StartNQLExport(ctx, request)
	if err != nil {
		log.Fatalf("Failed to start NQL export: %v", err)
	}

	// Display results
	fmt.Printf("\n=== NQL Export Started ===\n")
	fmt.Printf("Export ID: %s\n", result.ExportID)
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n", resp.Duration)

	logger.Info("NQL export started successfully",
		zap.String("export_id", result.ExportID),
		zap.String("status", result.Status))

	fmt.Printf("\n✓ Export initiated successfully!\n")
	fmt.Printf("\n💡 Next steps:\n")
	fmt.Printf("   1. Use GetNQLExportStatus() to check the export status\n")
	fmt.Printf("   2. When status is COMPLETED, use DownloadNQLExport() to get the data\n")
	fmt.Printf("   3. Or use WaitForNQLExport() to wait for completion automatically\n")
	fmt.Printf("\n   Export ID to check: %s\n", result.ExportID)
}
