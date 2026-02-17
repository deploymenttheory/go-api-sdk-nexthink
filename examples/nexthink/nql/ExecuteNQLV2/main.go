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

// This example demonstrates how to execute an NQL query using the V2 API.
//
// Execute NQL V2 is optimized for:
// - Relatively small requests at high frequency
// - Real-time queries with small result sets
// - Interactive dashboards
// - Frequent polling operations
//
// V2 returns data as objects (map[string]any) - cleaner structured data than V1.
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

	// Execute a pre-configured NQL query
	// Replace with your actual query ID from Nexthink admin
	request := &nql.ExecuteRequest{
		QueryID: "#your_query_id", // Must be pre-created in Nexthink
		// Platform: "windows", // Optional: filter by platform
	}

	result, resp, err := nxClient.NQL.ExecuteNQLV2(ctx, request)
	if err != nil {
		log.Fatalf("Failed to execute NQL query: %v", err)
	}

	// Display results
	fmt.Printf("\n=== NQL V2 Query Results ===\n")
	fmt.Printf("Query ID: %s\n", result.QueryID)
	fmt.Printf("Executed Query: %s\n", result.ExecutedQuery)
	fmt.Printf("Rows: %d\n", result.Rows)
	fmt.Printf("Execution Time: %s\n", result.ExecutionDateTime)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %v\n\n", resp.Duration)

	// Display data rows (limit to first 5 rows)
	fmt.Printf("Data (first 5 rows as objects):\n")
	maxRows := 5
	if len(result.Data) < maxRows {
		maxRows = len(result.Data)
	}

	for i := 0; i < maxRows; i++ {
		fmt.Printf("\nRow %d:\n", i+1)
		for key, value := range result.Data[i] {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	if len(result.Data) > 5 {
		fmt.Printf("\n... and %d more rows\n", len(result.Data)-5)
	}

	logger.Info("NQL V2 query executed successfully",
		zap.String("query_id", result.QueryID),
		zap.Int64("rows", result.Rows))

	fmt.Printf("\n✓ NQL V2 execution completed successfully!\n")
	fmt.Printf("\n💡 Tip: V2 returns structured objects, making it easier to work with the data\n")
}
