package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

// Example 02: Query Builder
//
// This example demonstrates programmatic query construction using the fluent Query Builder API:
// - Building queries with method chaining
// - Query validation
// - Using the builder with actual API execution
//
// The Query Builder helps you construct queries programmatically with:
// - Type safety
// - IDE auto-completion
// - Built-in validation
// - Readable, maintainable code

func main() {
	client, err := nexthink.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	nqlService := client.NQL
	ctx := context.Background()

	fmt.Println("=== Example 02: Query Builder ===\n")

	// =========================================================================
	// Basic Query Construction
	// =========================================================================
	fmt.Println("1. Simple Query")
	fmt.Println("---------------")

	query1 := nql.NewQueryBuilder().
		FromDevices().
		DuringPast(7, nql.Days).
		List("device.name", "operating_system.name").
		Limit(10).
		Build()

	fmt.Printf("Generated NQL:\n%s\n\n", query1)

	// =========================================================================
	// Query with Filters
	// =========================================================================
	fmt.Println("2. Query with Filters")
	fmt.Println("---------------------")

	query2 := nql.NewQueryBuilder().
		FromDevices().
		DuringPast(7, nql.Days).
		WhereEquals(nql.FieldOSPlatform, nql.PlatformWindows).
		WhereIn(nql.FieldHardwareType, []string{nql.HardwareTypeLaptop, nql.HardwareTypeDesktop}).
		List(nql.FieldDeviceName, nql.FieldOSName, nql.FieldHardwareType).
		SortAsc(nql.FieldDeviceName).
		Build()

	fmt.Printf("Generated NQL:\n%s\n\n", query2)

	// =========================================================================
	// Query with Event Data
	// =========================================================================
	fmt.Println("3. Query with Event Data (Crashes)")
	fmt.Println("----------------------------------")

	query3 := nql.NewQueryBuilder().
		FromDevices().
		DuringPast(7, nql.Days).
		With("execution.crashes during past 7d").
		ComputeSum("total_crashes", "number_of_crashes").
		WhereGreaterEqual("total_crashes", "3").
		List("device.name", "total_crashes").
		SortDesc("total_crashes").
		Build()

	fmt.Printf("Generated NQL:\n%s\n\n", query3)

	// =========================================================================
	// Query with Aggregation
	// =========================================================================
	fmt.Println("4. Aggregation Query")
	fmt.Println("--------------------")

	query4 := nql.NewQueryBuilder().
		FromDevices().
		DuringPast(7, nql.Days).
		SummarizeCount("device_count").
		SummarizeBy(nql.FieldOSPlatform).
		SortDesc("device_count").
		Build()

	fmt.Printf("Generated NQL:\n%s\n\n", query4)

	// =========================================================================
	// Query Validation
	// =========================================================================
	fmt.Println("5. Query Validation")
	fmt.Println("-------------------")

	validBuilder := nql.NewQueryBuilder().
		FromDevices().
		With("execution.crashes during past 7d").
		ComputeSum("crashes", "number_of_crashes")

	if err := validBuilder.Validate(); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Query is valid")
	}

	// Invalid query example (compute without with/include)
	invalidBuilder := nql.NewQueryBuilder().
		FromDevices().
		ComputeSum("crashes", "number_of_crashes")

	if err := invalidBuilder.Validate(); err != nil {
		fmt.Printf("✓ Caught invalid query: %v\n", err)
	}
	fmt.Println()

	// =========================================================================
	// Using with API Execution
	// =========================================================================
	fmt.Println("6. Executing with Query Builder")
	fmt.Println("--------------------------------")

	queryID := os.Getenv("NEXTHINK_QUERY_ID")
	if queryID == "" {
		fmt.Println("Note: NEXTHINK_QUERY_ID not set")
		fmt.Println("To execute queries:")
		fmt.Println("1. Copy generated query to Nexthink admin")
		fmt.Println("2. Save with Query ID (e.g., #devices_list)")
		fmt.Println("3. Set NEXTHINK_QUERY_ID environment variable")
		fmt.Println("4. Use ExecuteQueryBuilder method:")
		fmt.Println()
		fmt.Println("   qb := nql.NewQueryBuilder().FromDevices()...")
		fmt.Println("   result, _, err := nqlService.ExecuteQueryBuilder(ctx, queryID, qb)")
	} else {
		// Demonstrate ExecuteQueryBuilder
		qb := nql.NewQueryBuilder().
			FromDevices().
			DuringPast(7, nql.Days).
			List("device.name")

		resultSet, _, err := nqlService.ExecuteQueryBuilder(ctx, queryID, qb)
		if err != nil {
			log.Printf("Execution failed: %v", err)
		} else {
			fmt.Printf("✓ Query executed successfully\n")
			fmt.Printf("  Rows: %d\n", resultSet.Rows())
		}
	}

	fmt.Println("\n=== Next: See 03_Templates for pre-built queries ===")
}
