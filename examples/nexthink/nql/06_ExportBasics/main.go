package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

// Example 06: Export Basics (Manual Workflow)
//
// This example demonstrates the manual export workflow for large datasets:
// - Starting an export operation
// - Polling for completion status
// - Downloading the exported data
//
// Use exports for:
// - Large result sets (thousands of rows)
// - Scheduled reports
// - Bulk data extraction
// - Historical analysis
//
// Export workflow:
// 1. Start → returns exportID
// 2. Poll status → wait for COMPLETED
// 3. Download → retrieve data from S3

func main() {
	client, err := nexthink.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	nqlService := client.NQL
	ctx := context.Background()

	queryID := os.Getenv("NEXTHINK_EXPORT_QUERY_ID")
	if queryID == "" {
		queryID = os.Getenv("NEXTHINK_QUERY_ID")
	}
	if queryID == "" {
		log.Fatal("NEXTHINK_QUERY_ID or NEXTHINK_EXPORT_QUERY_ID environment variable required")
	}

	fmt.Println("=== Example 06: Export Basics (Manual Workflow) ===\n")

	// =========================================================================
	// Step 1: Start Export
	// =========================================================================
	fmt.Println("Step 1: Starting Export")
	fmt.Println("-----------------------")

	startResp, _, err := nqlService.StartNQLExport(ctx, &nql.ExportRequest{
		QueryID: queryID,
		Format:  nql.ExportFormatCSV,
	})
	if err != nil {
		log.Fatalf("Failed to start export: %v", err)
	}

	exportID := startResp.ExportID
	fmt.Printf("✓ Export started\n")
	fmt.Printf("  Export ID: %s\n", exportID)
	fmt.Printf("  Status: %s\n\n", startResp.Status)

	// =========================================================================
	// Step 2: Poll for Completion
	// =========================================================================
	fmt.Println("Step 2: Polling for Completion")
	fmt.Println("-------------------------------")

	pollInterval := 5 * time.Second
	maxPolls := 30
	var status *nql.NQLExportStatusResponse

	for attempt := 1; attempt <= maxPolls; attempt++ {
		time.Sleep(pollInterval)

		status, _, err = nqlService.GetNQLExportStatus(ctx, exportID)
		if err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}

		fmt.Printf("Poll %d/%d: %s\n", attempt, maxPolls, status.Status)

		// Check if terminal status
		if status.Status == nql.ExportStatusCompleted {
			fmt.Println("✓ Export completed!")
			break
		} else if status.Status == nql.ExportStatusError {
			log.Fatalf("Export failed: %s", status.ErrorDescription)
		}
	}
	fmt.Println()

	// =========================================================================
	// Step 3: Download Export
	// =========================================================================
	fmt.Println("Step 3: Downloading Export Data")
	fmt.Println("--------------------------------")

	if status.Status != nql.ExportStatusCompleted {
		log.Fatal("Export did not complete in time")
	}

	if status.ResultsFileURL == "" {
		log.Fatal("No download URL provided")
	}

	fmt.Printf("Downloading from: %s\n", status.ResultsFileURL)

	data, err := nqlService.DownloadNQLExport(ctx, status.ResultsFileURL)
	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	fmt.Printf("✓ Downloaded %d bytes\n", len(data))

	// Save to file
	filename := fmt.Sprintf("export_%s.csv", exportID)
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
	} else {
		fmt.Printf("✓ Saved to: %s\n", filename)
	}
	fmt.Println()

	// =========================================================================
	// Alternative: Using WaitForNQLExport Helper
	// =========================================================================
	fmt.Println("4. Using WaitForNQLExport Helper")
	fmt.Println("---------------------------------")
	fmt.Println("For simpler polling, use WaitForNQLExport:")
	fmt.Println()
	fmt.Println("  startResp, _, _ := nqlService.StartNQLExport(ctx, req)")
	fmt.Println("  status, _ := nqlService.WaitForNQLExport(ctx, startResp.ExportID, 5*time.Second, 10*time.Minute)")
	fmt.Println("  data, _ := nqlService.DownloadNQLExport(ctx, status.ResultsFileURL)")
	fmt.Println()

	fmt.Println("=== Next: See 07_ExportWorkflow for fully automated exports ===")
}
