# NQL Result Processing Guide

This guide explains how to work with NQL query results using the Result Set helpers.

## Table of Contents

- [Overview](#overview)
- [V2 Result Set](#v2-result-set)
- [V1 Result Set](#v1-result-set)
- [Type-Safe Access](#type-safe-access)
- [Data Transformation](#data-transformation)
- [Metadata Extraction](#metadata-extraction)
- [Best Practices](#best-practices)

## Overview

The SDK provides Result Set helpers that make it easier to work with NQL query responses. These helpers offer:

- Type-safe data access
- Iteration utilities
- Filtering and mapping operations
- Format conversion
- Metadata extraction

## V2 Result Set

V2 responses return data as an array of objects (map-based), which is easier to work with.

### Creating a Result Set

```go
// Execute query
result, apiResp, err := nqlService.ExecuteNQLV2(ctx, &nql.ExecuteRequest{
    QueryID: "#my_query",
})

// Create result set
resultSet := nql.NewV2ResultSet(result)
```

### Basic Information

```go
// Get row count
rowCount := resultSet.Rows()

// Get field names
fields := resultSet.Fields()

// Check if a field exists
hasName := resultSet.HasField("device.name")
```

### Accessing Data

```go
// Get a value by row and field
value, err := resultSet.Get(0, "device.name")

// Get an entire row
row, err := resultSet.GetRow(0)

// Iterate through all rows
err = resultSet.IterateRows(func(row int, data map[string]any) error {
    fmt.Printf("Row %d: %v\n", row, data)
    return nil
})
```

### Type-Safe Retrieval

```go
// Get string value
deviceName, err := resultSet.GetString(0, "device.name")

// Get integer value
crashCount, err := resultSet.GetInt(0, "total_crashes")

// Get float value
avgMemory, err := resultSet.GetFloat(0, "avg_memory")

// Get boolean value
isActive, err := resultSet.GetBool(0, "is_active")
```

### Filtering Results

```go
// Filter rows based on a condition
filtered := resultSet.Filter(func(row map[string]any) bool {
    if platform, ok := row["operating_system.platform"].(string); ok {
        return platform == "Windows"
    }
    return false
})

fmt.Printf("Found %d Windows devices\n", len(filtered))
```

### Mapping/Transforming Results

```go
// Transform each row
simplified := resultSet.Map(func(row map[string]any) map[string]any {
    return map[string]any{
        "name":     row["device.name"],
        "platform": row["operating_system.platform"],
        "status":   "active",
    }
})
```

### Converting to JSON

```go
// Convert result set to JSON
jsonData, err := resultSet.ToJSON()
if err != nil {
    log.Fatal(err)
}

// Save to file
os.WriteFile("results.json", jsonData, 0644)
```

## V1 Result Set

V1 responses return data as a 2D array with separate headers.

### Creating a Result Set

```go
// Execute V1 query
result, apiResp, err := nqlService.ExecuteNQLV1(ctx, &nql.ExecuteRequest{
    QueryID: "#my_query",
})

// Create result set
resultSet := nql.NewV1ResultSet(result)
```

### Basic Information

```go
// Get dimensions
rowCount := resultSet.Rows()
colCount := resultSet.Columns()

// Get headers
headers := resultSet.Headers
```

### Accessing Data

```go
// Get value by row and column index
value, err := resultSet.Get(0, 1)

// Get entire row
row, err := resultSet.GetRow(0)

// Find column index by name
colIdx, err := resultSet.FindColumnIndex("device.name")

// Get value by column name
value, err := resultSet.GetByColumnName(0, "device.name")
```

### Type-Safe Retrieval

```go
// Get string value
deviceName, err := resultSet.GetString(0, 1)

// Get integer value
crashCount, err := resultSet.GetInt(0, 2)

// Get float value
avgValue, err := resultSet.GetFloat(0, 3)
```

### Converting to V2 Format

```go
// Convert V1 result set to V2 format
v2Data := resultSet.ToV2Format()

// Or convert the entire result set
v2ResultSet := nql.ConvertV1ToV2(resultSet)
```

### Iterating Rows

```go
err = resultSet.IterateRows(func(row int, values []any) error {
    fmt.Printf("Row %d:\n", row)
    for i, value := range values {
        fmt.Printf("  %s: %v\n", resultSet.Headers[i], value)
    }
    return nil
})
```

## Type-Safe Access

### Handling Nil Values

```go
// Check for nil before using
value, err := resultSet.Get(0, "field_name")
if err != nil {
    return err
}

if value == nil {
    fmt.Println("Field is nil")
    return
}

// Type-safe methods return zero values for nil
strValue, _ := resultSet.GetString(0, "field_name") // returns ""
intValue, _ := resultSet.GetInt(0, "field_name")    // returns 0
```

### Type Conversion

```go
// GetInt handles multiple numeric types
value, err := resultSet.GetInt(0, "count")
// Handles: int, int32, int64, float64, string

// GetFloat handles multiple numeric types
value, err := resultSet.GetFloat(0, "average")
// Handles: float32, float64, int, int64, string
```

### Error Handling

```go
value, err := resultSet.GetString(0, "device.name")
if err != nil {
    // Handle errors:
    // - Row/column out of bounds
    // - Field not found
    // - Type mismatch
    log.Printf("Error getting value: %v", err)
    return
}
```

## Data Transformation

### Example: Extract Specific Fields

```go
type DeviceSummary struct {
    Name     string
    Platform string
    Crashes  int64
}

summaries := make([]DeviceSummary, 0)

resultSet.IterateRows(func(row int, data map[string]any) error {
    summary := DeviceSummary{}
    
    if name, err := resultSet.GetString(row, "device.name"); err == nil {
        summary.Name = name
    }
    
    if platform, err := resultSet.GetString(row, "operating_system.platform"); err == nil {
        summary.Platform = platform
    }
    
    if crashes, err := resultSet.GetInt(row, "total_crashes"); err == nil {
        summary.Crashes = crashes
    }
    
    summaries = append(summaries, summary)
    return nil
})
```

### Example: Aggregate Data

```go
// Calculate totals
totalCrashes := int64(0)
deviceCount := resultSet.Rows()

resultSet.IterateRows(func(row int, data map[string]any) error {
    if crashes, err := resultSet.GetInt(row, "total_crashes"); err == nil {
        totalCrashes += crashes
    }
    return nil
})

avgCrashes := float64(totalCrashes) / float64(deviceCount)
fmt.Printf("Average crashes per device: %.2f\n", avgCrashes)
```

### Example: Group By Field

```go
// Group devices by platform
groups := make(map[string][]map[string]any)

resultSet.IterateRows(func(row int, data map[string]any) error {
    if platform, ok := data["operating_system.platform"].(string); ok {
        groups[platform] = append(groups[platform], data)
    }
    return nil
})

for platform, devices := range groups {
    fmt.Printf("%s: %d devices\n", platform, len(devices))
}
```

## Metadata Extraction

### Basic Metadata

```go
metadata := nql.GetV2Metadata(result, apiResp)

fmt.Printf("Query ID: %s\n", metadata.QueryID)
fmt.Printf("Executed Query: %s\n", metadata.ExecutedQuery)
fmt.Printf("Rows Returned: %d\n", metadata.RowsReturned)
fmt.Printf("Response Time: %v\n", metadata.ResponseDuration)
fmt.Printf("Response Size: %d bytes\n", metadata.ResponseSize)
```

### Rate Limit Information

```go
rateLimitInfo := metadata.GetRateLimitInfo()

if rateLimitInfo != nil {
    if rateLimitInfo.Limit != "" {
        fmt.Printf("Rate Limit: %s\n", rateLimitInfo.Limit)
    }
    if rateLimitInfo.Remaining != "" {
        fmt.Printf("Remaining: %s\n", rateLimitInfo.Remaining)
    }
    if rateLimitInfo.Reset != "" {
        fmt.Printf("Reset: %s\n", rateLimitInfo.Reset)
    }
}
```

### Execution Time

```go
// Get execution time
execTime := metadata.ExecutionTime
fmt.Printf("Query executed at: %v\n", execTime)

// Time since execution
timeSince := metadata.TimeSinceExecution()
fmt.Printf("Time since execution: %v\n", timeSince)
```

### String Representation

```go
// Print metadata summary
fmt.Println(metadata.String())

// Output:
// ExecutionMetadata{
//   QueryID: #my_query
//   ExecutedQuery: devices during past 7d | list device.name
//   RowsReturned: 142
//   ExecutionTime: 2024-02-08T10:15:30Z
//   ResponseDuration: 1.234s
//   ResponseSize: 15432 bytes
//   ResponseStatus: 200
// }
```

## Best Practices

### 1. Use Type-Safe Methods

```go
// Good - type-safe with error handling
deviceName, err := resultSet.GetString(0, "device.name")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

// Avoid - requires manual type assertion
value, _ := resultSet.Get(0, "device.name")
deviceName := value.(string) // Can panic!
```

### 2. Check for Nil Values

```go
value, err := resultSet.Get(0, "optional_field")
if err != nil {
    return err
}

if value == nil {
    // Handle nil case
    return
}

// Use value
```

### 3. Use Iteration for Large Result Sets

```go
// Good - memory efficient
err = resultSet.IterateRows(func(row int, data map[string]any) error {
    // Process row
    return nil
})

// Avoid - loads all rows into memory
allRows := make([]map[string]any, resultSet.Rows())
for i := 0; i < resultSet.Rows(); i++ {
    allRows[i], _ = resultSet.GetRow(i)
}
```

### 4. Handle Errors Appropriately

```go
// Check bounds
if row >= resultSet.Rows() {
    return fmt.Errorf("row %d out of bounds", row)
}

// Check field existence (V2)
if !resultSet.HasField("field_name") {
    return fmt.Errorf("field not found")
}

// Handle type conversion errors
value, err := resultSet.GetInt(0, "field")
if err != nil {
    log.Printf("Could not convert to int: %v", err)
    // Use default or skip
}
```

### 5. Use V2 Format When Possible

```go
// V2 is easier to work with
result, _, err := nqlService.ExecuteNQLV2(ctx, req)

// If you must use V1, convert to V2
result, _, err := nqlService.ExecuteNQLV1(ctx, req)
resultSetV1 := nql.NewV1ResultSet(result)
resultSetV2 := nql.ConvertV1ToV2(resultSetV1)
```

## Next Steps

- See [NQL Query Building Guide](nql-query-building.md) for constructing queries
- See [NQL Export Workflow Guide](nql-export-workflow.md) for large data exports
- See [Examples](../../examples/nexthink/nql/ResultSetProcessing/main.go) for complete code samples
