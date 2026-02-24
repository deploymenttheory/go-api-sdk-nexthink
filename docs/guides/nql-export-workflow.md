# NQL Export Workflow Guide

This guide explains how to use the simplified export workflow for large NQL queries.

## Table of Contents

- [Overview](#overview)
- [When to Use Exports](#when-to-use-exports)
- [Simple Export Methods](#simple-export-methods)
- [Export with Options](#export-with-options)
- [Progress Tracking](#progress-tracking)
- [Manual Export Workflow](#manual-export-workflow)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Overview

NQL exports are designed for large queries that return substantial amounts of data. The export workflow is asynchronous:

1. **Start Export**: Initiate the export operation
2. **Poll Status**: Check status until completion
3. **Download Data**: Retrieve the exported data from S3

The SDK provides simplified methods that handle this entire workflow automatically.

## When to Use Exports

Use exports for:

- **Large Result Sets**: Queries returning thousands of rows
- **Scheduled Reports**: Automated data extraction
- **Bulk Data Analysis**: Comprehensive data dumps
- **Historical Analysis**: Long time-range queries

Use synchronous Execute for:

- **Real-Time Queries**: Small result sets needed immediately
- **Interactive Dashboards**: Live data updates
- **API Integrations**: Frequent small queries

## Simple Export Methods

### CSV Export

```go
// Export to CSV (simplest method)
result, err := nqlService.ExportToCSV(ctx, "#my_large_query")
if err != nil {
    log.Fatal(err)
}

// Save to file
os.WriteFile("export.csv", result.Data, 0644)

// Check results
fmt.Printf("Export ID: %s\n", result.ExportID)
fmt.Printf("Size: %s\n", result.SizeFormatted())
fmt.Printf("Duration: %v\n", result.TotalDuration)
fmt.Printf("Polls: %d\n", result.PollCount)
```

### JSON Export

```go
// Export to JSON
result, err := nqlService.ExportToJSON(ctx, "#my_large_query")
if err != nil {
    log.Fatal(err)
}

// Save to file
os.WriteFile("export.json", result.Data, 0644)
```

### Export with Progress

```go
// Export with simple progress callback
result, err := nqlService.ExportWithProgress(
    ctx,
    "#my_large_query",
    nql.ExportFormatCSV,
    func(status string) {
        fmt.Printf("Status: %s\n", status)
    },
)
```

## Export with Options

For more control, use `ExportWorkflow` with custom options:

### Basic Options

```go
opts := nql.DefaultExportOptions().
    WithFormat(nql.ExportFormatCSV).
    WithPollInterval(5 * time.Second).
    WithTimeout(10 * time.Minute)

result, err := nqlService.ExportWorkflow(ctx, &nql.ExportRequest{
    QueryID: "#my_query",
}, opts)
```

### Available Options

```go
type ExportOptions struct {
    // Format: "csv" or "json"
    Format string
    
    // How often to check status (default: 5s)
    PollInterval time.Duration
    
    // Maximum wait time (default: 10m)
    Timeout time.Duration
    
    // Progress callback
    OnProgress func(status string, elapsedTime time.Duration)
    
    // Status change callback
    OnStatusChange func(oldStatus, newStatus string, elapsedTime time.Duration)
}
```

## Progress Tracking

### Simple Progress Callback

```go
opts := nql.DefaultExportOptions().
    WithOnProgress(func(status string, elapsed time.Duration) {
        fmt.Printf("[%v] %s\n", elapsed.Round(time.Second), status)
    })

result, err := nqlService.ExportWorkflow(ctx, &nql.ExportRequest{
    QueryID: "#my_query",
}, opts)
```

Output:
```
[5s] SUBMITTED
[10s] IN_PROGRESS
[15s] IN_PROGRESS
[20s] COMPLETED
```

### Status Change Tracking

```go
opts := nql.DefaultExportOptions().
    WithOnStatusChange(func(oldStatus, newStatus string, elapsed time.Duration) {
        log.Printf("[%v] Status: %s → %s", elapsed, oldStatus, newStatus)
    })
```

### Combined Progress Tracking

```go
opts := nql.DefaultExportOptions().
    WithOnProgress(func(status string, elapsed time.Duration) {
        // Called on every poll
        fmt.Printf(".")
    }).
    WithOnStatusChange(func(oldStatus, newStatus string, elapsed time.Duration) {
        // Called only when status changes
        fmt.Printf("\n[%v] %s → %s\n", elapsed, oldStatus, newStatus)
    })
```

## Manual Export Workflow

For maximum control, use the low-level methods:

```go
// Step 1: Start export
startResp, _, err := nqlService.StartNQLExport(ctx, &nql.ExportRequest{
    QueryID: "#my_query",
    Format:  nql.ExportFormatCSV,
})
if err != nil {
    log.Fatal(err)
}

exportID := startResp.ExportID
fmt.Printf("Export started: %s\n", exportID)

// Step 2: Poll for completion
ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

timeout := time.After(10 * time.Minute)

var status *nql.NQLExportStatusResponse

for {
    select {
    case <-timeout:
        log.Fatal("Export timed out")
        
    case <-ticker.C:
        status, _, err = nqlService.GetNQLExportStatus(ctx, exportID)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("Status: %s\n", status.Status)
        
        if status.Status == nql.ExportStatusCompleted {
            goto download
        } else if status.Status == nql.ExportStatusError {
            log.Fatalf("Export failed: %s", status.ErrorDescription)
        }
    }
}

download:
// Step 3: Download data
fmt.Printf("Downloading from: %s\n", status.ResultsFileURL)
data, err := nqlService.DownloadNQLExport(ctx, status.ResultsFileURL)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Downloaded %d bytes\n", len(data))
os.WriteFile("export.csv", data, 0644)
```

## Error Handling

### Export Failures

```go
result, err := nqlService.ExportWorkflow(ctx, req, opts)
if err != nil {
    // Check error type
    if strings.Contains(err.Error(), "timeout") {
        log.Println("Export timed out - increase timeout option")
    } else if strings.Contains(err.Error(), "export failed") {
        log.Println("Export failed on server side")
    } else {
        log.Printf("Unknown error: %v", err)
    }
    return
}
```

### Status Checking

```go
// Check if export is ready
isReady, err := nqlService.IsExportReady(ctx, exportID)
if err != nil {
    log.Printf("Failed to check status: %v", err)
    return
}

if !isReady {
    fmt.Println("Export not ready yet")
    return
}

// Get human-readable progress
progress, err := nqlService.GetExportProgress(ctx, exportID)
fmt.Printf("Progress: %s\n", progress)
```

### Context Cancellation

```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

result, err := nqlService.ExportWorkflow(ctx, req, opts)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Export cancelled due to context timeout")
    } else if ctx.Err() == context.Canceled {
        log.Println("Export cancelled by user")
    } else {
        log.Printf("Export failed: %v", err)
    }
}
```

## Best Practices

### 1. Choose Appropriate Poll Interval

```go
// For small exports (< 1 minute expected)
opts.WithPollInterval(3 * time.Second)

// For medium exports (1-5 minutes expected)
opts.WithPollInterval(5 * time.Second)

// For large exports (> 5 minutes expected)
opts.WithPollInterval(10 * time.Second)
```

### 2. Set Reasonable Timeouts

```go
// Short queries
opts.WithTimeout(5 * time.Minute)

// Medium queries
opts.WithTimeout(10 * time.Minute)

// Very large queries
opts.WithTimeout(30 * time.Minute)
```

### 3. Use Progress Callbacks for Long Operations

```go
lastUpdate := time.Now()

opts := nql.DefaultExportOptions().
    WithOnProgress(func(status string, elapsed time.Duration) {
        if time.Since(lastUpdate) > 30*time.Second {
            log.Printf("Still running: %s (elapsed: %v)", status, elapsed)
            lastUpdate = time.Now()
        }
    })
```

### 4. Handle Large Files Appropriately

```go
result, err := nqlService.ExportWorkflow(ctx, req, opts)
if err != nil {
    return err
}

// Check size before loading into memory
fmt.Printf("Export size: %s\n", result.SizeFormatted())

if result.Size() > 100*1024*1024 { // > 100MB
    log.Println("Warning: Large export file")
    
    // Consider streaming to disk instead of keeping in memory
    err = os.WriteFile("large_export.csv", result.Data, 0644)
    result.Data = nil // Free memory
}
```

### 5. Save Exports with Timestamps

```go
result, err := nqlService.ExportWorkflow(ctx, req, opts)
if err != nil {
    return err
}

// Generate filename with timestamp
timestamp := time.Now().Format("20060102_150405")
filename := fmt.Sprintf("export_%s_%s.%s", 
    result.ExportID, 
    timestamp, 
    result.Format)

os.WriteFile(filename, result.Data, 0644)
```

### 6. Monitor Export Performance

```go
pollCount := 0
totalWaitTime := time.Duration(0)

opts := nql.DefaultExportOptions().
    WithOnProgress(func(status string, elapsed time.Duration) {
        pollCount++
        totalWaitTime = elapsed
    })

result, err := nqlService.ExportWorkflow(ctx, req, opts)

// Log performance metrics
log.Printf("Export completed: polls=%d, wait=%v, size=%s",
    pollCount, totalWaitTime, result.SizeFormatted())
```

## Complete Example

```go
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

func main() {
    // Create client
    client, err := nexthink.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    nqlService := client.NQL

    // Configure export options
    opts := nql.DefaultExportOptions().
        WithFormat(nql.ExportFormatCSV).
        WithPollInterval(5 * time.Second).
        WithTimeout(15 * time.Minute).
        WithOnProgress(func(status string, elapsed time.Duration) {
            fmt.Printf("[%v] %s\n", elapsed.Round(time.Second), status)
        }).
        WithOnStatusChange(func(old, new string, elapsed time.Duration) {
            log.Printf("Status changed: %s → %s", old, new)
        })

    // Execute export workflow
    result, err := nqlService.ExportWorkflow(ctx, &nql.ExportRequest{
        QueryID: "#large_device_query",
        Format:  nql.ExportFormatCSV,
    }, opts)

    if err != nil {
        log.Fatalf("Export failed: %v", err)
    }

    // Save results
    filename := fmt.Sprintf("export_%s.csv", time.Now().Format("20060102"))
    err = os.WriteFile(filename, result.Data, 0644)
    if err != nil {
        log.Fatal(err)
    }

    // Log summary
    fmt.Printf("\n✓ Export completed successfully!\n")
    fmt.Printf("  Export ID: %s\n", result.ExportID)
    fmt.Printf("  Size: %s\n", result.SizeFormatted())
    fmt.Printf("  Duration: %v\n", result.TotalDuration)
    fmt.Printf("  Polls: %d\n", result.PollCount)
    fmt.Printf("  Saved to: %s\n", filename)
}
```

## Next Steps

- See [NQL Query Building Guide](nql-query-building.md) for constructing queries
- See [NQL Result Processing Guide](nql-result-processing.md) for processing results
- See [Examples](../../examples/nexthink/nql/ExportWorkflow/main.go) for complete code samples
