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

// Example 07: Export Workflow (Simplified)
//
// This example demonstrates the simplified export workflow:
// - One-line export methods (ExportToCSV, ExportToJSON)
// - Progress tracking with callbacks
// - Customizable options
//
// The workflow methods handle everything automatically:
// - Starting the export
// - Polling for completion
// - Downloading results
//
// Compare this to Example 06 to see the simplification!

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

	fmt.Println("=== Example 07: Export Workflow (Simplified) ===\n")

	// =========================================================================
	// Simple CSV Export (One Line!)
	// =========================================================================
	fmt.Println("1. Simple CSV Export")
	fmt.Println("--------------------")

	result, err := nqlService.ExportToCSV(ctx, queryID)
	if err != nil {
		log.Fatalf("Export failed: %v", err)
	}

	fmt.Printf("✓ Export completed!\n")
	fmt.Printf("  Export ID: %s\n", result.ExportID)
	fmt.Printf("  Size: %s\n", result.SizeFormatted())
	fmt.Printf("  Duration: %v\n", result.TotalDuration)
	fmt.Printf("  Polls made: %d\n\n", result.PollCount)

	// Save to file
	os.WriteFile("export_simple.csv", result.Data, 0644)
	fmt.Println("  Saved to: export_simple.csv\n")

	// =========================================================================
	// Export with Progress Tracking
	// =========================================================================
	fmt.Println("2. Export with Progress Tracking")
	fmt.Println("---------------------------------")

	result2, err := nqlService.ExportWithProgress(
		ctx,
		queryID,
		nql.ExportFormatJSON,
		func(status string) {
			fmt.Printf("  Status: %s\n", status)
		},
	)
	if err != nil {
		log.Printf("Export failed: %v", err)
	} else {
		fmt.Printf("\n✓ Export completed: %s\n\n", result2.SizeFormatted())
		os.WriteFile("export_progress.json", result2.Data, 0644)
	}

	// =========================================================================
	// Export with Custom Options
	// =========================================================================
	fmt.Println("3. Export with Custom Options")
	fmt.Println("------------------------------")

	opts := nql.DefaultExportOptions().
		WithFormat(nql.ExportFormatCSV).
		WithPollInterval(3 * time.Second).
		WithTimeout(15 * time.Minute).
		WithOnProgress(func(status string, elapsed time.Duration) {
			fmt.Printf("  [%v] %s\n", elapsed.Round(time.Second), status)
		}).
		WithOnStatusChange(func(oldStatus, newStatus string, elapsed time.Duration) {
			fmt.Printf("  [%v] Status: %s → %s\n",
				elapsed.Round(time.Second), oldStatus, newStatus)
		})

	result3, err := nqlService.ExportWorkflow(ctx, &nql.ExportRequest{
		QueryID: queryID,
		Format:  nql.ExportFormatCSV,
	}, opts)

	if err != nil {
		log.Printf("Export failed: %v", err)
	} else {
		fmt.Printf("\n✓ Export completed!\n")
		fmt.Printf("  Size: %s\n", result3.SizeFormatted())
		fmt.Printf("  Total duration: %v\n", result3.TotalDuration)
		fmt.Printf("  Polls: %d\n\n", result3.PollCount)

		filename := fmt.Sprintf("export_custom_%s.csv", time.Now().Format("20060102"))
		os.WriteFile(filename, result3.Data, 0644)
		fmt.Printf("  Saved to: %s\n", filename)
	}

	fmt.Println("\n=== Next: See 08_Integration for complete end-to-end workflow ===")
}
