# NQL Query Templates Guide

This guide explains how to use pre-built query templates for common use cases.

## Table of Contents

- [Overview](#overview)
- [Available Templates](#available-templates)
- [Device Templates](#device-templates)
- [User Templates](#user-templates)
- [Application Templates](#application-templates)
- [Performance Templates](#performance-templates)
- [DEX Score Templates](#dex-score-templates)
- [Monitoring Templates](#monitoring-templates)
- [Workflow Templates](#workflow-templates)
- [Customizing Templates](#customizing-templates)

## Overview

Query templates provide pre-built NQL queries for common monitoring and analysis scenarios. They save time and ensure best practices are followed.

### Basic Usage

```go
import "github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"

// Create templates instance
templates := nql.NewTemplates()

// Generate a query
query := templates.DevicesWithCrashes("during past 7d", "outlook.exe")

// Use the query as reference or save it to Nexthink admin
fmt.Println(query)
```

## Available Templates

List all available templates:

```go
allTemplates := templates.GetAllTemplates()
for _, tmpl := range allTemplates {
    fmt.Println(tmpl)
}
```

## Device Templates

### Devices with Crashes

Find devices experiencing application crashes:

```go
// All devices with crashes
query := templates.DevicesWithCrashes("during past 7d", "")

// Specific application
query := templates.DevicesWithCrashes("during past 7d", "outlook.exe")

// With wildcard
query := templates.DevicesWithCrashes("during past 30d", "chrome*")
```

### Devices with High Memory Usage

Find devices with high memory utilization:

```go
// Devices with > 90% memory usage
query := templates.DevicesWithHighMemoryUsage(90, "during past 7d")

// Devices with > 95% memory usage
query := templates.DevicesWithHighMemoryUsage(95, "during past 24h")
```

### Devices by Platform

Count devices grouped by operating system:

```go
query := templates.DevicesByPlatform("during past 7d")

// Returns counts by Windows, macOS, Linux
```

### Devices with Slow Boot Time

Find devices with boot times exceeding a threshold:

```go
// Devices with boot time > 60 seconds
query := templates.DevicesWithSlowBootTime(60, "during past 7d")

// Devices with boot time > 90 seconds
query := templates.DevicesWithSlowBootTime(90, "during past 30d")
```

### Devices with System Crashes

Find devices with system crashes above threshold:

```go
// Devices with 3+ system crashes
query := templates.DevicesWithSystemCrashes(3, "during past 7d")
```

## User Templates

### Users with Web Errors

Find users experiencing web application errors:

```go
// All users with web errors
query := templates.UsersWithWebErrors("during past 7d", "")

// Specific application
query := templates.UsersWithWebErrors("during past 24h", "Salesforce")
```

### Users with Poor Collaboration Quality

Find users with poor audio/video quality:

```go
query := templates.UsersWithPoorCollaborationQuality("during past 24h")

// Returns users with poor call quality
```

## Application Templates

### Applications with High Error Rate

Find applications with high web error rates:

```go
// Applications with high error rates
// minPageViews: minimum threshold to avoid false positives
query := templates.ApplicationsWithHighErrorRate(100, "during past 60min")
```

### Top Crashing Applications

Find applications with the most crashes:

```go
// Top 10 crashing applications
query := templates.TopCrashingApplications("during past 7d", 10)

// Top 20 crashing applications
query := templates.TopCrashingApplications("during past 30d", 20)
```

### Web Page Load Performance

Analyze page load performance:

```go
// All applications
query := templates.WebPageLoadPerformance("", "during past 7d")

// Specific application
query := templates.WebPageLoadPerformance("Confluence", "during past 7d")
```

## Performance Templates

### Network Connectivity Issues

Find devices with connectivity problems:

```go
query := templates.NetworkConnectivityIssues("during past 7d")

// Returns devices with poor WiFi signal or high noise
```

### Binaries with High Crash Rate

Find binaries causing crashes:

```go
// Binaries with 20+ crashes
query := templates.BinariesWithHighCrashRate(20, "during past 7d")
```

## DEX Score Templates

### Overall DEX Score

Calculate overall DEX score for population:

```go
query := templates.OverallDEXScore("during past 24h")
```

### DEX Score by Platform

Break down DEX scores by operating system:

```go
query := templates.DEXScoreByPlatform("during past 24h")

// Returns scores for Windows, macOS, Linux
```

### Users with Low DEX Score

Find users with DEX scores below threshold:

```go
// Users with DEX score < 50
query := templates.UsersWithLowDEXScore(50, "during past 24h")

// Users with DEX score < 70
query := templates.UsersWithLowDEXScore(70, "during past 24h")
```

### DEX Score Impact by Component

Analyze DEX score impact by component:

```go
// Logon speed impact
query := templates.DEXScoreImpactByComponent("logon_speed", "during past 24h")

// Boot speed impact
query := templates.DEXScoreImpactByComponent("boot_speed", "during past 24h")

// Software reliability impact
query := templates.DEXScoreImpactByComponent("software_reliability", "during past 24h")
```

## Monitoring Templates

### Workflow Execution Success

Track successful workflow executions:

```go
query := templates.WorkflowExecutionSuccess("during past 30d")

// Returns success counts by workflow name
```

### Remote Action Cost Savings

Estimate cost savings from remote actions:

```go
// $20 saved per successful execution
query := templates.RemoteActionSavingsEstimate(20, "during past 30d")

// $50 saved per execution
query := templates.RemoteActionSavingsEstimate(50, "during past 90d")
```

## Customizing Templates

### Modify Template Output

```go
// Get template query
query := templates.DevicesWithCrashes("during past 7d", "outlook.exe")

// Parse and modify if needed
// Or use as a reference to create your own variant
```

### Create Your Own Templates

```go
// Create a helper function
func MyCustomTemplate(threshold int, period string) string {
    return nql.NewQueryBuilder().
        FromDevices().
        During(period).
        Include(fmt.Sprintf("execution.crashes %s", period)).
        ComputeSum("crashes", "number_of_crashes").
        WhereGreaterEqual("crashes", fmt.Sprintf("%d", threshold)).
        List("device.name", "crashes").
        SortDesc("crashes").
        Build()
}

// Use it
query := MyCustomTemplate(10, "during past 7d")
```

### Combine Templates with Query Builder

```go
// Start with a template
baseQuery := templates.DevicesWithCrashes("during past 7d", "")

// Parse and extend (or build from scratch using similar pattern)
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    With("execution.crashes during past 7d").
    ComputeSum("total_crashes", "number_of_crashes").
    // Add custom filters
    WhereEquals("device.entity", "Production").
    WhereNotIn(nql.FieldHardwareType, []string{"virtual"}).
    List("device.name", "total_crashes", "device.entity").
    SortDesc("total_crashes").
    Build()
```

## Using Templates with the API

Templates generate query strings that must be saved in Nexthink admin before use:

```go
// Step 1: Generate query from template
query := templates.DevicesWithCrashes("during past 7d", "outlook.exe")

// Step 2: Save query in Nexthink admin
// - Go to Admin > NQL API queries
// - Create new query
// - Paste the generated query
// - Save with a Query ID (e.g., #devices_with_outlook_crashes)

// Step 3: Execute using the Query ID
result, _, err := nqlService.ExecuteNQLV2(ctx, &nql.ExecuteRequest{
    QueryID: "#devices_with_outlook_crashes",
})
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-api-sdk-nexthink/nexthink"
    "github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"
)

func main() {
    // Create templates
    templates := nql.NewTemplates()

    // Generate queries for different use cases
    fmt.Println("=== Device Health Queries ===\n")

    // 1. Crash analysis
    crashQuery := templates.DevicesWithCrashes("during past 7d", "")
    fmt.Println("Crash Analysis Query:")
    fmt.Println(crashQuery)
    fmt.Println()

    // 2. Performance monitoring
    perfQuery := templates.DevicesWithHighMemoryUsage(90, "during past 7d")
    fmt.Println("Performance Query:")
    fmt.Println(perfQuery)
    fmt.Println()

    // 3. DEX score analysis
    dexQuery := templates.DEXScoreByPlatform("during past 24h")
    fmt.Println("DEX Score Query:")
    fmt.Println(dexQuery)
    fmt.Println()

    // Save these queries to Nexthink admin and execute them
}
```

## Next Steps

- See [NQL Query Building Guide](nql-query-building.md) for custom query construction
- See [NQL Result Processing Guide](nql-result-processing.md) for handling results
- See [Examples](../../examples/nexthink/nql/Templates/main.go) for complete code samples
