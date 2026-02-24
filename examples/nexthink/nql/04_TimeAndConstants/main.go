package main

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

// Example 04: Time Selection and Constants
//
// This example demonstrates using SDK helpers for type-safe query construction:
// - Time selection builders and constants
// - Data model constants (tables, fields, values)
// - Operators and functions
//
// These helpers provide:
// - Type safety (catch errors at compile time)
// - IDE auto-completion
// - Consistency across queries
// - Less typing, fewer errors

func main() {
	fmt.Println("=== Example 04: Time Selection and Constants ===\n")

	// =========================================================================
	// Time Selection - Predefined Constants
	// =========================================================================
	fmt.Println("1. Predefined Time Constants")
	fmt.Println("----------------------------")

	fmt.Printf("Past 7 days:    %s\n", nql.Past7Days)
	fmt.Printf("Past 24 hours:  %s\n", nql.Past24Hours)
	fmt.Printf("Past 30 days:   %s\n", nql.Past30Days)
	fmt.Printf("Yesterday:      %s\n", nql.Yesterday)
	fmt.Println()

	// Using in queries
	query1 := nql.NewQueryBuilder().
		FromDevices().
		During(nql.Past7Days).
		List("device.name").
		Build()

	fmt.Printf("Query with constant:\n%s\n\n", query1)

	// =========================================================================
	// Time Selection - Builder API
	// =========================================================================
	fmt.Println("2. Time Selection Builder")
	fmt.Println("-------------------------")

	// During past
	time1 := nql.NewTimeSelection().
		DuringPast(7, nql.Days).
		Build()
	fmt.Printf("During past: %s\n", time1)

	// Date range
	time2 := nql.NewTimeSelection().
		From("2024-01-01").
		To("2024-01-31").
		Build()
	fmt.Printf("Date range: %s\n", time2)

	// Relative range
	time3 := nql.NewTimeSelection().
		FromRelative(21, nql.Days).
		ToRelative(13, nql.Days).
		Build()
	fmt.Printf("Relative range: %s\n", time3)

	// Specific date
	time4 := nql.NewTimeSelection().
		On("Feb 8, 2024").
		Build()
	fmt.Printf("Specific date: %s\n\n", time4)

	// =========================================================================
	// Data Model Constants - Tables
	// =========================================================================
	fmt.Println("3. Table Constants")
	fmt.Println("------------------")

	fmt.Printf("Devices table:    %s\n", nql.TableDevices)
	fmt.Printf("Users table:      %s\n", nql.TableUsers)
	fmt.Printf("Crashes table:    %s\n", nql.TableExecutionCrashes)
	fmt.Printf("Web errors:       %s\n", nql.TableWebErrors)
	fmt.Printf("DEX scores:       %s\n", nql.TableDexScores)
	fmt.Println()

	// Using in queries
	query2 := nql.NewQueryBuilder().
		From(nql.TableExecutionCrashes).
		DuringPast(7, nql.Days).
		SummarizeCount("crash_count").
		Build()

	fmt.Printf("Query with table constant:\n%s\n\n", query2)

	// =========================================================================
	// Data Model Constants - Fields
	// =========================================================================
	fmt.Println("4. Field Constants")
	fmt.Println("------------------")

	fmt.Printf("Device name:      %s\n", nql.FieldDeviceName)
	fmt.Printf("OS platform:      %s\n", nql.FieldOSPlatform)
	fmt.Printf("Hardware type:    %s\n", nql.FieldHardwareType)
	fmt.Printf("User name:        %s\n", nql.FieldUserName)
	fmt.Printf("Binary name:      %s\n", nql.FieldBinaryName)
	fmt.Println()

	// Using in queries
	query3 := nql.NewQueryBuilder().
		FromDevices().
		DuringPast(7, nql.Days).
		WhereEquals(nql.FieldOSPlatform, nql.PlatformWindows).
		List(nql.FieldDeviceName, nql.FieldOSName, nql.FieldHardwareType).
		Build()

	fmt.Printf("Query with field constants:\n%s\n\n", query3)

	// =========================================================================
	// Data Model Constants - Values
	// =========================================================================
	fmt.Println("5. Value Constants")
	fmt.Println("------------------")

	fmt.Printf("Windows platform:       %s\n", nql.PlatformWindows)
	fmt.Printf("macOS platform:         %s\n", nql.PlatformMacOS)
	fmt.Printf("Laptop hardware:        %s\n", nql.HardwareTypeLaptop)
	fmt.Printf("Desktop hardware:       %s\n", nql.HardwareTypeDesktop)
	fmt.Printf("Good experience:        %s\n", nql.ExperienceLevelGood)
	fmt.Printf("Frustrating experience: %s\n", nql.ExperienceLevelFrustrating)
	fmt.Println()

	// =========================================================================
	// Operators
	// =========================================================================
	fmt.Println("6. Operator Constants")
	fmt.Println("---------------------")

	fmt.Printf("Equals:           %s\n", nql.OpEquals)
	fmt.Printf("Not equals:       %s\n", nql.OpNotEquals)
	fmt.Printf("Greater than:     %s\n", nql.OpGreater)
	fmt.Printf("Less than:        %s\n", nql.OpLess)
	fmt.Printf("In list:          %s\n", nql.OpIn)
	fmt.Printf("Contains:         %s\n", nql.OpContains)
	fmt.Println()

	// =========================================================================
	// Complete Example with Constants
	// =========================================================================
	fmt.Println("7. Complete Example Using Constants")
	fmt.Println("------------------------------------")

	query4 := nql.NewQueryBuilder().
		Comment("Query using all constant types").
		From(nql.TableDevices).
		During(nql.Past7Days).
		WhereEquals(nql.FieldOSPlatform, nql.PlatformWindows).
		WhereIn(nql.FieldHardwareType, []string{
			nql.HardwareTypeLaptop,
			nql.HardwareTypeDesktop,
		}).
		List(
			nql.FieldDeviceName,
			nql.FieldOSName,
			nql.FieldHardwareType,
			nql.FieldHardwareManufacturer,
		).
		SortAsc(nql.FieldDeviceName).
		Limit(20).
		Build()

	fmt.Printf("%s\n\n", query4)

	// =========================================================================
	// Time Granularity (for aggregations)
	// =========================================================================
	fmt.Println("8. Time Granularity for Aggregations")
	fmt.Println("-------------------------------------")

	fmt.Printf("15 minutes:  %s\n", nql.Granularity15Min)
	fmt.Printf("1 hour:      %s\n", nql.Granularity1Hour)
	fmt.Printf("1 day:       %s\n", nql.Granularity1Day)
	fmt.Printf("7 days:      %s\n", nql.Granularity7Days)
	fmt.Println()

	query5 := nql.NewQueryBuilder().
		From(nql.TableExecutionCrashes).
		DuringPast(30, nql.Days).
		SummarizeCount("crashes").
		SummarizeByTime(nql.Granularity1Day).
		Build()

	fmt.Printf("Time-bucketed aggregation:\n%s\n\n", query5)

	// =========================================================================
	// Helper Functions
	// =========================================================================
	fmt.Println("9. Time Helper Functions")
	fmt.Println("------------------------")

	// Format current time for NQL
	now := time.Now()
	fmt.Printf("Current date:     %s\n", nql.FormatDate(now))
	fmt.Printf("Current datetime: %s\n", nql.FormatDateTime(now))
	fmt.Printf("Short format:     %s\n", nql.FormatDateShort(now))
	fmt.Println()

	fmt.Println("=== Next: See 05_ResultProcessing for handling query results ===")
}
