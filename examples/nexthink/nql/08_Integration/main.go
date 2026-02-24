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

// Example 08: Integration - Complete Workflow
//
// This example demonstrates a complete end-to-end workflow combining all SDK features:
// - Using templates to generate queries
// - Executing with result sets
// - Processing results with type-safe methods
// - Extracting metadata for monitoring
// - Building custom queries
//
// Scenario: Monitor device health across the organization

func main() {
	client, err := nexthink.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	nqlService := client.NQL
	ctx := context.Background()
	templates := nql.NewTemplates()

	fmt.Println("=== Example 08: Integration - Device Health Monitoring ===\n")

	// =========================================================================
	// Scenario: Identify devices needing attention
	// =========================================================================

	// Step 1: Use template to generate crash analysis query
	fmt.Println("Step 1: Generate Query from Template")
	fmt.Println("-------------------------------------")

	crashTemplate := templates.DevicesWithCrashes("during past 7d", "")
	fmt.Println("Generated query for crash analysis:")
	fmt.Println(crashTemplate.Query())
	fmt.Println()
	fmt.Println("Note: Copy this query to Nexthink admin and save with ID #device_crashes")
	fmt.Println()

	// Step 2: Build a custom query using QueryBuilder
	fmt.Println("Step 2: Build Custom Query")
	fmt.Println("--------------------------")

	customQuery := nql.NewQueryBuilder().
		Comment("Production Windows laptops with crashes").
		From(nql.TableDevices).
		During(nql.Past7Days).
		WhereEquals(nql.FieldOSPlatform, nql.PlatformWindows).
		WhereEquals(nql.FieldHardwareType, nql.HardwareTypeLaptop).
		WhereEquals("device.entity", "Production").
		With("execution.crashes during past 7d").
		ComputeSum("total_crashes", "number_of_crashes").
		WhereGreaterEqual("total_crashes", "3").
		List(
			nql.FieldDeviceName,
			nql.FieldOSName,
			"total_crashes",
		).
		SortDesc("total_crashes").
		Limit(50)

	// Validate before using
	if err := customQuery.Validate(); err != nil {
		log.Fatalf("Query validation failed: %v", err)
	}

	fmt.Println("✓ Custom query validated")
	fmt.Println(customQuery.Build())
	fmt.Println()

	// Step 3: Execute queries if credentials are available
	queryID1 := os.Getenv("NEXTHINK_QUERY_ID")
	if queryID1 == "" {
		fmt.Println("Step 3: Query Execution")
		fmt.Println("-----------------------")
		fmt.Println("Note: Set NEXTHINK_QUERY_ID to execute queries")
		fmt.Println("Workflow:")
		fmt.Println("  1. Save query in Nexthink admin with ID")
		fmt.Println("  2. Execute: resultSet, _, err := nqlService.ExecuteV2WithResultSet(ctx, req)")
		fmt.Println("  3. Process: deviceName, _ := resultSet.GetString(0, \"device.name\")")
		fmt.Println()
	} else {
		fmt.Println("Step 3: Execute and Process Results")
		fmt.Println("------------------------------------")

		// Execute with result set convenience method
		resultSet, apiResp, err := nqlService.ExecuteV2WithResultSet(ctx, &nql.ExecuteRequest{
			QueryID: queryID1,
		})
		if err != nil {
			log.Printf("Execution failed: %v", err)
		} else {
			// Extract metadata
			fmt.Printf("✓ Query executed in %v\n", apiResp.Duration)
			fmt.Printf("  Rows: %d\n", resultSet.Rows())

			// Process first few results
			if resultSet.Rows() > 0 {
				fmt.Println("\nDevices needing attention:")

				count := 0
				resultSet.IterateRows(func(row int, data map[string]any) error {
					if count < 5 {
						name, _ := resultSet.GetString(row, "device.name")
						fmt.Printf("  • %s\n", name)
						count++
					}
					return nil
				})

				// Filter Windows devices
				windowsDevices := resultSet.Filter(func(row map[string]any) bool {
					platform, ok := row[nql.FieldOSPlatform].(string)
					return ok && platform == nql.PlatformWindows
				})

				fmt.Printf("\nWindows devices: %d of %d\n", len(windowsDevices), resultSet.Rows())
			}
		}
		fmt.Println()
	}

	// =========================================================================
	// Step 4: Large data export if needed
	// =========================================================================
	fmt.Println("Step 4: Large Data Export (if needed)")
	fmt.Println("--------------------------------------")

	exportQueryID := os.Getenv("NEXTHINK_EXPORT_QUERY_ID")
	if exportQueryID != "" {
		fmt.Println("Starting export with progress tracking...")

		opts := nql.DefaultExportOptions().
			WithFormat(nql.ExportFormatCSV).
			WithOnProgress(func(status string, elapsed time.Duration) {
				if elapsed%5*time.Second == 0 {
					fmt.Printf("  [%v] %s\n", elapsed.Round(time.Second), status)
				}
			})

		exportResult, err := nqlService.ExportWorkflow(ctx, &nql.ExportRequest{
			QueryID: exportQueryID,
		}, opts)

		if err != nil {
			log.Printf("Export failed: %v", err)
		} else {
			fmt.Printf("\n✓ Export completed: %s in %v\n",
				exportResult.SizeFormatted(),
				exportResult.TotalDuration)

			timestamp := time.Now().Format("20060102_150405")
			filename := fmt.Sprintf("device_health_%s.csv", timestamp)
			os.WriteFile(filename, exportResult.Data, 0644)
			fmt.Printf("  Saved to: %s\n", filename)
		}
	} else {
		fmt.Println("Note: Set NEXTHINK_EXPORT_QUERY_ID for large data exports")
	}

	// =========================================================================
	// Summary
	// =========================================================================
	fmt.Println("\n=== Workflow Summary ===")
	fmt.Println("This example demonstrated:")
	fmt.Println("  ✓ Template-based query generation")
	fmt.Println("  ✓ Custom query building with validation")
	fmt.Println("  ✓ Type-safe result processing")
	fmt.Println("  ✓ Metadata extraction")
	fmt.Println("  ✓ Data filtering and transformation")
	fmt.Println("  ✓ Simplified export workflow")
	fmt.Println()
	fmt.Println("For production use:")
	fmt.Println("  1. Create queries in Nexthink admin")
	fmt.Println("  2. Use ExecuteV2WithResultSet for type safety")
	fmt.Println("  3. Monitor with metadata extraction")
	fmt.Println("  4. Export large datasets with ExportWorkflow")
}
