package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

// Example 01: Basic NQL Query Execution
//
// This example demonstrates the fundamentals of executing NQL queries:
// - Creating a Nexthink client
// - Executing queries with V1 and V2 APIs
// - Understanding the response format differences
//
// Prerequisites:
// - NEXTHINK_CLIENT_ID, NEXTHINK_CLIENT_SECRET environment variables
// - NEXTHINK_INSTANCE, NEXTHINK_REGION environment variables
// - NEXTHINK_QUERY_ID environment variable with a query ID (e.g., #my_query)
//
// The query must be pre-created in Nexthink admin (Admin > NQL API queries)

func main() {
	// Create Nexthink client from environment variables
	client, err := nexthink.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	nqlService := client.NQL
	ctx := context.Background()

	// Get query ID from environment
	queryID := os.Getenv("NEXTHINK_QUERY_ID")
	if queryID == "" {
		log.Fatal("NEXTHINK_QUERY_ID environment variable required")
	}

	fmt.Println("=== Example 01: Basic NQL Query Execution ===\n")

	// =========================================================================
	// V2 API (Recommended)
	// =========================================================================
	fmt.Println("1. Executing with V2 API (returns map-based data)")
	fmt.Println("------------------------------------------------")

	resultV2, apiResp, err := nqlService.ExecuteNQLV2(ctx, &nql.ExecuteRequest{
		QueryID: queryID,
	})
	if err != nil {
		log.Fatalf("V2 execution failed: %v", err)
	}

	fmt.Printf("✓ Query executed successfully\n")
	fmt.Printf("  Query ID: %s\n", resultV2.QueryID)
	fmt.Printf("  Rows returned: %d\n", resultV2.Rows)
	fmt.Printf("  Response time: %v\n", apiResp.Duration)
	fmt.Printf("  Execution time: %s\n\n", resultV2.ExecutionDateTime)

	// V2 returns data as array of maps (easier to work with)
	if len(resultV2.Data) > 0 {
		fmt.Println("First row (V2 format):")
		for key, value := range resultV2.Data[0] {
			fmt.Printf("  %s: %v\n", key, value)
		}
		fmt.Println()
	}

	// =========================================================================
	// V1 API (For compatibility)
	// =========================================================================
	fmt.Println("2. Executing with V1 API (returns array-based data)")
	fmt.Println("----------------------------------------------------")

	resultV1, apiResp, err := nqlService.ExecuteNQLV1(ctx, &nql.ExecuteRequest{
		QueryID: queryID,
	})
	if err != nil {
		log.Fatalf("V1 execution failed: %v", err)
	}

	fmt.Printf("✓ Query executed successfully\n")
	fmt.Printf("  Query ID: %s\n", resultV1.QueryID)
	fmt.Printf("  Rows returned: %d\n", resultV1.Rows)
	fmt.Printf("  Response time: %v\n\n", apiResp.Duration)

	// V1 returns headers separately from data (2D array format)
	fmt.Printf("Headers: %v\n", resultV1.Headers)
	if len(resultV1.Data) > 0 {
		fmt.Printf("First row: %v\n\n", resultV1.Data[0])
	}

	// =========================================================================
	// Key Differences
	// =========================================================================
	fmt.Println("3. V1 vs V2 Comparison")
	fmt.Println("----------------------")
	fmt.Println("V1:")
	fmt.Println("  • Headers: []string")
	fmt.Println("  • Data: [][]any (2D array)")
	fmt.Println("  • ExecutionDateTime: DateTime object")
	fmt.Println("  • Access by index: data[row][col]")
	fmt.Println()
	fmt.Println("V2 (Recommended):")
	fmt.Println("  • Data: []map[string]any (array of objects)")
	fmt.Println("  • ExecutionDateTime: ISO string")
	fmt.Println("  • Access by key: data[row][\"field_name\"]")
	fmt.Println("  • Easier to work with")
	fmt.Println()

	fmt.Println("=== Next: See 02_QueryBuilder for programmatic query construction ===")
}
