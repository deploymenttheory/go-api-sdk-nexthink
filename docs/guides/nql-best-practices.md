# NQL Best Practices Guide

This guide provides best practices for using NQL with the Nexthink Go SDK.

## Table of Contents

- [Query Construction](#query-construction)
- [Performance Optimization](#performance-optimization)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Production Considerations](#production-considerations)
- [Common Pitfalls](#common-pitfalls)

## Query Construction

### Use the Query Builder

```go
// Good - type-safe, validated
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    List("device.name").
    Build()

// Avoid - error-prone, no validation
query := "devices during past 7d | list device.name"
```

### Use Constants for Fields

```go
// Good - autocomplete, refactor-safe
qb.WhereEquals(nql.FieldDeviceName, "server-01")
qb.List(nql.FieldOSPlatform, nql.FieldHardwareType)

// Avoid - typo-prone
qb.WhereEquals("device.name", "server-01")
qb.List("operating_system.platform", "hardware.type")
```

### Validate Before Execution

```go
// Always validate
query, err := qb.BuildAndValidate()
if err != nil {
    return fmt.Errorf("invalid query: %w", err)
}

// Or validate separately
if err := qb.Validate(); err != nil {
    return err
}
query := qb.Build()
```

### Add Comments for Complex Queries

```go
query := nql.NewQueryBuilder().
    Comment("Production device health analysis").
    Comment("Excludes test and virtual environments").
    FromDevices().
    // ... rest of query
    Build()
```

## Performance Optimization

### Choose Execute vs Export

```go
// Small result sets (< 1000 rows) - use Execute
result, _, err := nqlService.ExecuteNQLV2(ctx, &nql.ExecuteRequest{
    QueryID: "#small_query",
})

// Large result sets (> 1000 rows) - use Export
result, err := nqlService.ExportToCSV(ctx, "#large_query")
```

### Use Appropriate Time Ranges

```go
// Good - specific time range
qb.DuringPast(7, nql.Days)

// Avoid - unnecessarily long ranges
qb.DuringPast(365, nql.Days) // May time out or return too much data
```

### Limit Results for Interactive Queries

```go
// Good - limit for UI display
qb.List("device.name", "crash_count").
   SortDesc("crash_count").
   Limit(100)

// Avoid - returning all results when only displaying top items
qb.List("device.name", "crash_count").
   SortDesc("crash_count")
   // No limit - may return thousands of rows
```

### Use Summarize for Aggregations

```go
// Good - single aggregated value
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    SummarizeCount("total_devices").
    Build()

// Avoid - computing in application code
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    List("device.name").
    Build()
// Then counting rows in Go - inefficient
```

### Filter Early in Query

```go
// Good - filter before compute
qb.FromDevices().
   DuringPast(7, nql.Days).
   WhereEquals(nql.FieldOSPlatform, "Windows"). // Filter first
   With("execution.crashes during past 7d").
   ComputeSum("crashes", "number_of_crashes")

// Less efficient - compute then filter
qb.FromDevices().
   DuringPast(7, nql.Days).
   With("execution.crashes during past 7d").
   ComputeSum("crashes", "number_of_crashes").
   WhereEquals(nql.FieldOSPlatform, "Windows") // Filter later
```

## Error Handling

### Check All Errors

```go
// Execute query
result, apiResp, err := nqlService.ExecuteNQLV2(ctx, req)
if err != nil {
    // Handle error appropriately
    log.Printf("Query failed: %v", err)
    
    // Check for specific error types
    if client.IsRateLimited(err) {
        // Handle rate limiting
        time.Sleep(time.Minute)
        // Retry
    }
    
    return err
}

// Process result
resultSet := nql.NewV2ResultSet(result)
```

### Validate Inputs

```go
// Validate query ID format
if !strings.HasPrefix(queryID, "#") {
    return fmt.Errorf("query ID must start with #")
}

// Validate parameters
if period == "" {
    period = "during past 7d" // Default
}

// Validate thresholds
if threshold < 0 {
    return fmt.Errorf("threshold must be positive")
}
```

### Handle Empty Results

```go
resultSet := nql.NewV2ResultSet(result)

if resultSet.Rows() == 0 {
    log.Println("Query returned no results")
    return
}

// Process results
```

### Type-Safe Data Access

```go
// Good - check errors
deviceName, err := resultSet.GetString(0, "device.name")
if err != nil {
    log.Printf("Error getting device name: %v", err)
    return
}

// Avoid - can panic
value, _ := resultSet.Get(0, "device.name")
deviceName := value.(string) // Panics if not string
```

## Testing

### Test Query Construction

```go
func TestQueryBuilder(t *testing.T) {
    query := nql.NewQueryBuilder().
        FromDevices().
        DuringPast(7, nql.Days).
        List("device.name").
        Build()

    expected := "devices during past 7d\n| list device.name"
    
    if query != expected {
        t.Errorf("Query mismatch:\nGot: %s\nWant: %s", query, expected)
    }
}
```

### Test Query Validation

```go
func TestQueryValidation(t *testing.T) {
    // Invalid query - compute without with/include
    qb := nql.NewQueryBuilder().
        FromDevices().
        ComputeSum("crashes", "number_of_crashes")

    err := qb.Validate()
    if err == nil {
        t.Error("Expected validation error, got nil")
    }
}
```

### Mock Query Results

```go
func TestResultProcessing(t *testing.T) {
    // Create mock result
    mockResult := &nql.ExecuteNQLV2Response{
        QueryID: "#test_query",
        Rows:    2,
        Data: []map[string]any{
            {"device.name": "device-01", "crash_count": 5},
            {"device.name": "device-02", "crash_count": 3},
        },
    }

    resultSet := nql.NewV2ResultSet(mockResult)
    
    // Test processing
    if resultSet.Rows() != 2 {
        t.Errorf("Expected 2 rows, got %d", resultSet.Rows())
    }
}
```

## Production Considerations

### Use Environment Variables for Configuration

```go
// Good - configurable
queryID := os.Getenv("NEXTHINK_DEVICE_QUERY_ID")
if queryID == "" {
    queryID = "#default_device_query"
}

// Avoid - hardcoded
queryID := "#my_specific_query"
```

### Implement Retry Logic

```go
func executeWithRetry(ctx context.Context, nqlService *nql.Service, req *nql.ExecuteRequest, maxRetries int) (*nql.ExecuteNQLV2Response, error) {
    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        result, _, err := nqlService.ExecuteNQLV2(ctx, req)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Check if retryable
        if client.IsRateLimited(err) || client.IsTransient(err) {
            backoff := time.Duration(attempt+1) * 2 * time.Second
            time.Sleep(backoff)
            continue
        }
        
        // Non-retryable error
        return nil, err
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

### Log Query Execution

```go
logger.Info("Executing NQL query",
    zap.String("query_id", req.QueryID),
    zap.String("platform", req.Platform))

result, apiResp, err := nqlService.ExecuteNQLV2(ctx, req)

if err != nil {
    logger.Error("Query execution failed",
        zap.String("query_id", req.QueryID),
        zap.Error(err))
    return err
}

logger.Info("Query executed successfully",
    zap.String("query_id", req.QueryID),
    zap.Int64("rows", result.Rows),
    zap.Duration("duration", apiResp.Duration))
```

### Monitor Performance

```go
// Track query performance
metadata := nql.GetV2Metadata(result, apiResp)

// Alert on slow queries
if metadata.ResponseDuration > 5*time.Second {
    logger.Warn("Slow query detected",
        zap.String("query_id", metadata.QueryID),
        zap.Duration("duration", metadata.ResponseDuration))
}

// Track rate limits
rateLimitInfo := metadata.GetRateLimitInfo()
if rateLimitInfo != nil && rateLimitInfo.Remaining != "" {
    logger.Info("Rate limit status",
        zap.String("remaining", rateLimitInfo.Remaining))
}
```

### Handle Export Failures Gracefully

```go
opts := nql.DefaultExportOptions().
    WithTimeout(15 * time.Minute).
    WithOnProgress(func(status string, elapsed time.Duration) {
        if elapsed > 10*time.Minute {
            logger.Warn("Export taking longer than expected",
                zap.Duration("elapsed", elapsed),
                zap.String("status", status))
        }
    })

result, err := nqlService.ExportWorkflow(ctx, req, opts)
if err != nil {
    logger.Error("Export failed",
        zap.String("query_id", req.QueryID),
        zap.Error(err))
    
    // Implement fallback or alerting
    return err
}
```

## Common Pitfalls

### 1. Using List and Summarize Together

```go
// WRONG - cannot use both
query := nql.NewQueryBuilder().
    FromDevices().
    List("device.name").        // ❌
    SummarizeCount("total").    // ❌
    Build()

// CORRECT - choose one
query1 := nql.NewQueryBuilder().
    FromDevices().
    List("device.name").
    Build()

query2 := nql.NewQueryBuilder().
    FromDevices().
    SummarizeCount("total").
    Build()
```

### 2. Compute Without With/Include

```go
// WRONG - compute needs with or include
query := nql.NewQueryBuilder().
    FromDevices().
    ComputeSum("crashes", "number_of_crashes"). // ❌
    Build()

// CORRECT
query := nql.NewQueryBuilder().
    FromDevices().
    With("execution.crashes during past 7d").   // ✓
    ComputeSum("crashes", "number_of_crashes"). // ✓
    Build()
```

### 3. Forgetting Time Selection

```go
// WRONG - devices without time selection may return no data
query := nql.NewQueryBuilder().
    FromDevices().
    List("device.name").
    Build()

// CORRECT
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    List("device.name").
    Build()
```

### 4. Not Handling Nil Values

```go
// WRONG - can panic
value, _ := resultSet.Get(0, "field")
str := value.(string) // Panics if nil or wrong type

// CORRECT
str, err := resultSet.GetString(0, "field")
if err != nil {
    // Handle error
}
if str == "" {
    // Handle empty/nil
}
```

### 5. Loading Large Exports into Memory

```go
// WRONG - memory intensive
result, _ := nqlService.ExportToCSV(ctx, "#huge_query")
// result.Data is 500MB in memory

// Process all in memory
data := string(result.Data) // ❌

// CORRECT - save to disk immediately
result, _ := nqlService.ExportToCSV(ctx, "#huge_query")

// Save to disk
err := os.WriteFile("export.csv", result.Data, 0644)

// Free memory
result.Data = nil

// Process from disk if needed
```

### 6. Not Setting Appropriate Timeouts

```go
// WRONG - using default context for long operations
ctx := context.Background()
result, _ := nqlService.ExportWorkflow(ctx, req, opts)

// CORRECT - set timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
defer cancel()

result, _ := nqlService.ExportWorkflow(ctx, req, opts)
```

## Performance Tips

### 1. Use Specific Fields in List

```go
// Good - only get needed fields
qb.List("device.name", "operating_system.platform")

// Avoid - default fields may include unnecessary data
qb.Build() // Returns default fields
```

### 2. Use Where Filters Effectively

```go
// Good - filter to relevant subset
qb.WhereEquals("device.entity", "Production").
   WhereEquals(nql.FieldOSPlatform, "Windows")

// Avoid - filtering in application code
qb.Build() // Returns all data, filter in Go
```

### 3. Use Appropriate Aggregations

```go
// Good - aggregate on server
qb.SummarizeCount("total_devices").
   SummarizeAvg("avg_crashes", "crash_count")

// Avoid - downloading all data to aggregate locally
qb.List("device.name", "crash_count")
// Then computing average in Go
```

### 4. Batch Related Queries

```go
// Good - single query with multiple metrics
query := nql.NewQueryBuilder().
    FromDevices().
    Include("execution.crashes during past 7d").
    ComputeSum("crashes", "number_of_crashes").
    Include("web.errors during past 7d").
    ComputeSum("errors", "number_of_errors").
    Build()

// Avoid - multiple separate queries
// query1: crashes
// query2: errors
// Inefficient network usage
```

## Security Considerations

### Validate Query IDs

```go
// Validate format
if !strings.HasPrefix(queryID, "#") {
    return fmt.Errorf("invalid query ID format")
}

// Sanitize if accepting user input
queryID = strings.TrimSpace(queryID)
```

### Don't Log Sensitive Data

```go
// Good
logger.Info("Query executed",
    zap.String("query_id", result.QueryID),
    zap.Int64("rows", result.Rows))

// Avoid - may contain sensitive data
logger.Info("Query executed",
    zap.Any("full_result", result)) // ❌
```

### Use Context for Cancellation

```go
// Good - respect cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Listen for shutdown signal
go func() {
    <-shutdownChan
    cancel()
}()

result, _ := nqlService.ExecuteNQLV2(ctx, req)
```

## Next Steps

- See [NQL Query Building Guide](nql-query-building.md) for query construction
- See [NQL Result Processing Guide](nql-result-processing.md) for handling results
- See [NQL Templates Guide](nql-templates.md) for pre-built queries
- See [Examples](../../examples/nexthink/nql/) for complete code samples
