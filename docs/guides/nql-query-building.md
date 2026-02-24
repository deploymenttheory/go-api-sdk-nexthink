# NQL Query Building Guide

This guide explains how to use the NQL Query Builder to construct Nexthink Query Language queries programmatically.

## Table of Contents

- [Overview](#overview)
- [Basic Usage](#basic-usage)
- [Query Components](#query-components)
- [Advanced Features](#advanced-features)
- [Validation](#validation)
- [Best Practices](#best-practices)
- [Examples](#examples)

## Overview

The NQL Query Builder provides a fluent API for constructing NQL queries in a type-safe, IDE-friendly manner. Instead of manually writing query strings, you can use method chaining to build queries programmatically.

### Benefits

- **Type Safety**: Catch errors at compile time
- **IDE Support**: Auto-completion and inline documentation
- **Validation**: Built-in query validation
- **Maintainability**: Easier to read and modify
- **Testability**: Easier to unit test query construction

## Basic Usage

```go
import "github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"

// Create a new query builder
qb := nql.NewQueryBuilder()

// Build a simple query
query := qb.
    FromDevices().
    DuringPast(7, nql.Days).
    List("device.name", "operating_system.name").
    Limit(10).
    Build()

// Output:
// devices during past 7d
// | list device.name, operating_system.name
// | limit 10
```

## Query Components

### Table Selection

Select the source table for your query:

```go
// Using shortcuts
qb.FromDevices()        // devices
qb.FromUsers()          // users
qb.FromApplications()   // applications
qb.FromBinaries()       // binaries

// Using constants
qb.From(nql.TableDevices)
qb.From(nql.TableExecutionCrashes)

// Using full table path
qb.From("execution.crashes")
```

### Time Selection

Specify the time frame for your query:

```go
// During past
qb.DuringPast(7, nql.Days)
qb.DuringPast(24, nql.Hours)
qb.DuringPast(30, nql.Minutes)

// Predefined constants
qb.During(nql.Past7Days)
qb.During(nql.Past24Hours)

// Specific date
qb.On("Feb 8, 2024")
qb.On("2024-02-08")

// Date range
qb.FromTo("2024-01-01", "2024-01-31")
```

### With and Include Clauses

Join event tables with object tables:

```go
// With clause (filters to objects with events)
qb.With("web.errors during past 7d")
qb.WithTable(nql.TableWebErrors, "during past 7d")

// Include clause (keeps all objects)
qb.Include("execution.crashes during past 7d")
qb.IncludeTable(nql.TableExecutionCrashes, "during past 7d")
```

### Compute Clause

Calculate metrics from events:

```go
// Generic compute
qb.Compute("total_crashes", "number_of_crashes.sum()")

// Convenience methods
qb.ComputeCount("device_count")
qb.ComputeSum("total_crashes", "number_of_crashes")
qb.ComputeAvg("avg_memory", "free_memory")
qb.ComputeMax("max_cpu", "cpu_usage")
qb.ComputeMin("min_memory", "free_memory")
qb.ComputeLast("last_boot", "boot_time")
```

### Where Clause

Filter results based on conditions:

```go
// Basic conditions
qb.WhereEquals("binary.name", "outlook.exe")
qb.WhereNotEquals("hardware.type", "virtual")
qb.WhereGreater("total_crashes", "5")
qb.WhereLess("free_memory", "1GB")
qb.WhereGreaterEqual("crash_count", "3")
qb.WhereLessEqual("boot_time", "60s")

// List conditions
qb.WhereIn("operating_system.platform", []string{"Windows", "macOS"})
qb.WhereNotIn("hardware.type", []string{"virtual", "null"})

// Array conditions
qb.WhereContains("tags", "VDI")
qb.WhereNotContains("tags", "test")

// Custom conditions
qb.Where("device.last_seen > 7d ago")
qb.Where("process_visibility == foreground")
```

### List Clause

Select specific fields to return:

```go
qb.List("device.name", "operating_system.name", "total_crashes")

// Using constants
qb.List(
    nql.FieldDeviceName,
    nql.FieldOSName,
    nql.FieldHardwareType,
)
```

### Sort Clause

Order results:

```go
qb.SortDesc("total_crashes")
qb.SortAsc("device.name")
qb.Sort("size", nql.SortDesc)
```

### Limit Clause

Restrict the number of rows:

```go
qb.Limit(100)
```

### Summarize Clause

Aggregate data:

```go
// Basic summarize
qb.SummarizeCount("total_devices")
qb.SummarizeSum("total_crashes", "number_of_crashes")
qb.SummarizeAvg("avg_score", "dex_score")

// Summarize with grouping
qb.SummarizeCount("device_count").
   SummarizeBy("operating_system.platform", "hardware.type")

// Summarize with time bucketing
qb.SummarizeCount("crash_count").
   SummarizeByTime(nql.Granularity1Day)

// Using constants
qb.SummarizeBy(nql.FieldOSPlatform)
```

### Comments

Add documentation to your queries:

```go
qb.Comment("Find devices with high crash rates in production")
```

## Advanced Features

### Conditional Query Building

Build queries dynamically based on runtime conditions:

```go
qb := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days)

// Add filters conditionally
if platform != "" {
    qb.WhereEquals(nql.FieldOSPlatform, platform)
}

if minCrashes > 0 {
    qb.With("execution.crashes during past 7d").
       ComputeSum("crashes", "number_of_crashes").
       WhereGreaterEqual("crashes", fmt.Sprintf("%d", minCrashes))
}

query := qb.Build()
```

### Multiple Where Clauses

Chain multiple where conditions (treated as AND):

```go
qb.WhereEquals("operating_system.platform", "Windows").
   WhereEquals("hardware.type", "laptop").
   WhereGreater("crash_count", "5")

// Equivalent to:
// where operating_system.platform == "Windows"
// where hardware.type == "laptop"
// where crash_count > 5
```

### Multiple Compute Clauses

Calculate multiple metrics:

```go
qb.ComputeSum("total_crashes", "number_of_crashes").
   ComputeAvg("avg_memory", "free_memory").
   ComputeMax("max_cpu", "cpu_usage")
```

## Validation

Validate queries before execution:

```go
qb := nql.NewQueryBuilder().
    FromDevices().
    With("execution.crashes during past 7d").
    ComputeSum("crashes", "number_of_crashes")

// Validate the query
if err := qb.Validate(); err != nil {
    log.Printf("Query validation failed: %v", err)
    return
}

// Build after validation
query := qb.Build()

// Or combine both
query, err := qb.BuildAndValidate()
if err != nil {
    log.Printf("Query validation failed: %v", err)
    return
}
```

### Validation Rules

The validator checks for:

- Table selection is present
- Compute requires with or include clause
- Cannot use both list and summarize
- Proper operator usage

## Best Practices

### 1. Use Constants

Prefer constants over hardcoded strings:

```go
// Good
qb.WhereEquals(nql.FieldDeviceName, "server-01")
qb.From(nql.TableDevices)

// Avoid
qb.WhereEquals("device.name", "server-01")
qb.From("devices")
```

### 2. Validate Before Execution

Always validate queries before sending to the API:

```go
query, err := qb.BuildAndValidate()
if err != nil {
    return fmt.Errorf("invalid query: %w", err)
}
```

### 3. Use Comments for Complex Queries

Document complex queries:

```go
qb.Comment("Find production devices with crash rates exceeding threshold").
   Comment("Excludes virtual machines and test environments")
```

### 4. Break Down Complex Queries

For readability, build queries in stages:

```go
qb := nql.NewQueryBuilder()

// Stage 1: Base selection
qb.FromDevices().DuringPast(7, nql.Days)

// Stage 2: Add crash analysis
qb.With("execution.crashes during past 7d").
   ComputeSum("total_crashes", "number_of_crashes")

// Stage 3: Add filters
qb.WhereEquals(nql.FieldOSPlatform, "Windows").
   WhereGreaterEqual("total_crashes", "5")

// Stage 4: Format output
qb.List(nql.FieldDeviceName, "total_crashes").
   SortDesc("total_crashes").
   Limit(20)

query := qb.Build()
```

### 5. Reuse Builder Patterns

Create helper functions for common patterns:

```go
func devicesWithMetric(metric string, threshold int) *nql.QueryBuilder {
    return nql.NewQueryBuilder().
        FromDevices().
        DuringPast(7, nql.Days).
        WhereGreaterEqual(metric, fmt.Sprintf("%d", threshold))
}

// Usage
query := devicesWithMetric("crash_count", 5).
    List("device.name", "crash_count").
    Build()
```

## Examples

### Example 1: Simple Device Query

```go
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    List("device.name", "operating_system.name").
    Limit(10).
    Build()
```

### Example 2: Crash Analysis

```go
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(30, nql.Days).
    With("execution.crashes during past 30d").
    WhereEquals("binary.name", "outlook.exe").
    ComputeSum("total_crashes", "number_of_crashes").
    WhereGreaterEqual("total_crashes", "5").
    List("device.name", "total_crashes").
    SortDesc("total_crashes").
    Build()
```

### Example 3: DEX Score by Platform

```go
query := nql.NewQueryBuilder().
    FromDevices().
    During(nql.Past24Hours).
    Include("dex.scores during past 24h").
    ComputeAvg("dex_per_device", "value").
    Where("dex_per_device != NULL").
    SummarizeAvg("dex_score", "dex_per_device").
    SummarizeBy(nql.FieldOSPlatform).
    SortDesc("dex_score").
    Build()
```

### Example 4: Time-Series Aggregation

```go
query := nql.NewQueryBuilder().
    From(nql.TableExecutionCrashes).
    DuringPast(30, nql.Days).
    SummarizeCount("crash_count").
    Summarize("devices_affected", "device.count()").
    SummarizeByTime(nql.Granularity1Day).
    SortAsc("start_time").
    Build()
```

### Example 5: Complex Multi-Condition Query

```go
query := nql.NewQueryBuilder().
    Comment("Production devices with performance issues").
    FromDevices().
    DuringPast(7, nql.Days).
    WhereIn(nql.FieldOSPlatform, []string{"Windows", "macOS"}).
    WhereEquals("device.entity", "Production").
    WhereNotEquals(nql.FieldHardwareType, nql.HardwareTypeVirtual).
    With("device_performance.events during past 7d").
    ComputeAvg("avg_cpu", "cpu_usage").
    ComputeAvg("avg_memory", "memory_usage").
    WhereGreater("avg_cpu", "80").
    List(
        nql.FieldDeviceName,
        nql.FieldOSName,
        "avg_cpu",
        "avg_memory",
    ).
    SortDesc("avg_cpu").
    Limit(50).
    Build()
```

## Next Steps

- See [NQL Result Processing Guide](nql-result-processing.md) for handling query results
- See [NQL Templates Guide](nql-templates.md) for pre-built query templates
- See [NQL Best Practices](nql-best-practices.md) for optimization tips
