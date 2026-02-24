package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

// Example 03: Query Templates
//
// This example demonstrates using pre-built query templates:
// - Accessing template library
// - Generating queries from templates
// - Converting templates to API requests
// - Executing template-based queries
//
// Templates provide ready-to-use queries for common scenarios like:
// - Device health monitoring
// - User experience analysis
// - Application performance
// - DEX score tracking

func main() {
	client, err := nexthink.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	nqlService := client.NQL
	ctx := context.Background()
	templates := nql.NewTemplates()

	fmt.Println("=== Example 03: Query Templates ===\n")

	// =========================================================================
	// Exploring Available Templates
	// =========================================================================
	fmt.Println("1. Available Templates")
	fmt.Println("----------------------")

	allTemplates := templates.GetAllTemplates()
	fmt.Printf("Found %d pre-built templates:\n", len(allTemplates))
	for i, name := range allTemplates[:5] {
		fmt.Printf("  %d. %s\n", i+1, name)
	}
	fmt.Printf("  ... and %d more\n\n", len(allTemplates)-5)

	// =========================================================================
	// Device Health Templates
	// =========================================================================
	fmt.Println("2. Device Health Monitoring")
	fmt.Println("---------------------------")

	// Template for devices with crashes
	crashTemplate := templates.DevicesWithCrashes("during past 7d", "outlook.exe")
	fmt.Println("Devices with Outlook crashes:")
	fmt.Printf("%s\n\n", crashTemplate.Query())

	// Template for high memory usage
	memoryTemplate := templates.DevicesWithHighMemoryUsage(90, "during past 7d")
	fmt.Println("Devices with >90% memory usage:")
	fmt.Printf("%s\n\n", memoryTemplate.Query())

	// =========================================================================
	// User Experience Templates
	// =========================================================================
	fmt.Println("3. User Experience Analysis")
	fmt.Println("---------------------------")

	// Users with web errors
	webErrorTemplate := templates.UsersWithWebErrors("during past 7d", "Salesforce")
	fmt.Println("Users with Salesforce errors:")
	fmt.Printf("%s\n\n", webErrorTemplate.Query())

	// =========================================================================
	// DEX Score Templates
	// =========================================================================
	fmt.Println("4. DEX Score Analysis")
	fmt.Println("---------------------")

	// Overall DEX score
	dexTemplate := templates.DEXScoreByPlatform("during past 24h")
	fmt.Println("DEX scores by platform:")
	fmt.Printf("%s\n\n", dexTemplate.Query())

	// =========================================================================
	// Using Templates with API
	// =========================================================================
	fmt.Println("5. Template to API Request")
	fmt.Println("--------------------------")

	// Templates return *Template objects with helper methods:
	// - .Query() - get the NQL query string
	// - .ToRequest(queryID) - convert to ExecuteRequest
	// - .QueryBuilder() - access the underlying builder

	queryID := os.Getenv("NEXTHINK_QUERY_ID")
	if queryID == "" {
		fmt.Println("To execute a template query:")
		fmt.Println("1. Copy template query to Nexthink admin")
		fmt.Println("2. Save with Query ID (e.g., #devices_with_crashes)")
		fmt.Println("3. Convert template to request:")
		fmt.Println()
		fmt.Println("   template := templates.DevicesWithCrashes(...)")
		fmt.Println("   req := template.ToRequest(\"#devices_with_crashes\")")
		fmt.Println("   result, _, err := nqlService.ExecuteNQLV2(ctx, req)")
	} else {
		fmt.Printf("Executing template query with ID: %s\n", queryID)

		// Get a template
		template := templates.DevicesByPlatform("during past 7d")

		// Convert to ExecuteRequest
		req := template.ToRequest(queryID)

		// Execute the query
		result, _, err := nqlService.ExecuteNQLV2(ctx, req)
		if err != nil {
			log.Printf("Execution failed: %v", err)
		} else {
			fmt.Printf("✓ Template query executed successfully\n")
			fmt.Printf("  Rows returned: %d\n", result.Rows)
			
			// Show first result
			if len(result.Data) > 0 {
				fmt.Println("\nFirst result:")
				for key, value := range result.Data[0] {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
		}
	}

	// =========================================================================
	// Template Categories
	// =========================================================================
	fmt.Println("\n6. Template Categories")
	fmt.Println("----------------------")
	fmt.Println("Device Health:")
	fmt.Println("  • DevicesWithCrashes")
	fmt.Println("  • DevicesWithHighMemoryUsage")
	fmt.Println("  • DevicesWithSlowBootTime")
	fmt.Println()
	fmt.Println("User Experience:")
	fmt.Println("  • UsersWithWebErrors")
	fmt.Println("  • UsersWithPoorCollaborationQuality")
	fmt.Println()
	fmt.Println("DEX Scores:")
	fmt.Println("  • OverallDEXScore")
	fmt.Println("  • DEXScoreByPlatform")
	fmt.Println("  • UsersWithLowDEXScore")
	fmt.Println()
	fmt.Println("Performance:")
	fmt.Println("  • WebPageLoadPerformance")
	fmt.Println("  • NetworkConnectivityIssues")

	fmt.Println("\n=== Next: See 04_TimeAndConstants for helpers and constants ===")
}
