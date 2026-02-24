package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

// Example 05: Result Processing
//
// This example demonstrates type-safe result processing:
// - Using ExecuteV2WithResultSet convenience method
// - Type-safe data access (GetString, GetInt, GetFloat)
// - Iterating through results
// - Filtering and transforming data
// - Extracting metadata
//
// Result sets provide type-safe access and prevent common errors like:
// - Type assertion panics
// - Index out of bounds
// - Nil pointer dereferencing

func main() {
	client, err := nexthink.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	nqlService := client.NQL
	ctx := context.Background()

	queryID := os.Getenv("NEXTHINK_QUERY_ID")
	if queryID == "" {
		log.Fatal("NEXTHINK_QUERY_ID environment variable required")
	}

	fmt.Println("=== Example 05: Result Processing ===\n")

	// =========================================================================
	// Execute with Result Set (New Method)
	// =========================================================================
	fmt.Println("1. Execute with Result Set")
	fmt.Println("--------------------------")

	// New convenience method returns V2ResultSet directly
	resultSet, apiResp, err := nqlService.ExecuteV2WithResultSet(ctx, &nql.ExecuteRequest{
		QueryID: queryID,
	})
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	fmt.Printf("✓ Query executed successfully\n")
	fmt.Printf("  Rows: %d\n", resultSet.Rows())
	fmt.Printf("  Response time: %v\n\n", apiResp.Duration)

	// =========================================================================
	// Type-Safe Data Access
	// =========================================================================
	fmt.Println("2. Type-Safe Data Access")
	fmt.Println("------------------------")

	if resultSet.Rows() > 0 {
		// Get available fields
		fields := resultSet.Fields()
		fmt.Printf("Available fields: %v\n\n", fields)

		// Type-safe getters with error handling
		if resultSet.HasField("device.name") {
			deviceName, err := resultSet.GetString(0, "device.name")
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				fmt.Printf("Device name (string): %s\n", deviceName)
			}
		}

		if resultSet.HasField("total_crashes") {
			crashes, err := resultSet.GetInt(0, "total_crashes")
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				fmt.Printf("Crashes (int64): %d\n", crashes)
			}
		}
		fmt.Println()
	}

	// =========================================================================
	// Iterating Through Results
	// =========================================================================
	fmt.Println("3. Iterating Through Results")
	fmt.Println("----------------------------")

	count := 0
	err = resultSet.IterateRows(func(row int, data map[string]any) error {
		if count < 3 {
			fmt.Printf("Row %d:\n", row+1)
			for field, value := range data {
				fmt.Printf("  %s: %v\n", field, value)
			}
			fmt.Println()
			count++
		}
		return nil
	})
	if err != nil {
		log.Printf("Iteration error: %v", err)
	}

	// =========================================================================
	// Filtering Results
	// =========================================================================
	fmt.Println("4. Filtering Results")
	fmt.Println("--------------------")

	// Filter to Windows devices only
	filtered := resultSet.Filter(func(row map[string]any) bool {
		if platform, ok := row[nql.FieldOSPlatform].(string); ok {
			return platform == nql.PlatformWindows
		}
		return false
	})

	fmt.Printf("Filtered to Windows devices: %d rows\n", len(filtered))
	if len(filtered) > 0 {
		fmt.Println("First filtered result:")
		for key, value := range filtered[0] {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	fmt.Println()

	// =========================================================================
	// Transforming Results
	// =========================================================================
	fmt.Println("5. Transforming Results")
	fmt.Println("-----------------------")

	// Extract only specific fields
	simplified := resultSet.Map(func(row map[string]any) map[string]any {
		return map[string]any{
			"name":     row["device.name"],
			"platform": row[nql.FieldOSPlatform],
			"active":   true,
		}
	})

	fmt.Printf("Transformed %d rows\n", len(simplified))
	if len(simplified) > 0 {
		fmt.Printf("Example transformed row: %v\n", simplified[0])
	}
	fmt.Println()

	// =========================================================================
	// Metadata Extraction
	// =========================================================================
	fmt.Println("6. Metadata Extraction")
	fmt.Println("----------------------")

	// Get metadata from raw response
	rawResult, apiResp, _ := nqlService.ExecuteNQLV2(ctx, &nql.ExecuteRequest{
		QueryID: queryID,
	})

	metadata := nql.GetV2Metadata(rawResult, apiResp)
	if metadata != nil {
		fmt.Printf("Query ID: %s\n", metadata.QueryID)
		fmt.Printf("Rows: %d\n", metadata.RowsReturned)
		fmt.Printf("Response time: %v\n", metadata.ResponseDuration)
		fmt.Printf("Response size: %d bytes\n", metadata.ResponseSize)

		// Rate limit information
		rateLimitInfo := metadata.GetRateLimitInfo()
		if rateLimitInfo != nil && rateLimitInfo.Remaining != "" {
			fmt.Printf("Rate limit remaining: %s\n", rateLimitInfo.Remaining)
		}
	}
	fmt.Println()

	// =========================================================================
	// JSON Export
	// =========================================================================
	fmt.Println("7. Converting to JSON")
	fmt.Println("---------------------")

	jsonData, err := resultSet.ToJSON()
	if err != nil {
		log.Printf("JSON conversion failed: %v", err)
	} else {
		// Save to file
		err = os.WriteFile("results.json", jsonData, 0644)
		if err != nil {
			log.Printf("Save failed: %v", err)
		} else {
			fmt.Printf("✓ Saved %d bytes to results.json\n", len(jsonData))
		}
	}

	fmt.Println("\n=== Next: See 06_ExportBasics for large data exports ===")
}
