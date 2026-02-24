# NQL SDK Reference

Complete reference for the NQL service enhancements in the Nexthink Go SDK.

## Table of Contents

- [Query Builder API](#query-builder-api)
- [Time Selection](#time-selection)
- [Result Sets](#result-sets)
- [Templates](#templates)
- [Export Workflow](#export-workflow)
- [Data Model Constants](#data-model-constants)
- [Operators](#operators)
- [Functions](#functions)
- [Validation](#validation)
- [Metadata](#metadata)

## Query Builder API

### QueryBuilder Methods

| Method | Description | Example |
|--------|-------------|---------|
| `From(table)` | Set source table | `From("devices")` |
| `FromDevices()` | Select devices table | `FromDevices()` |
| `FromUsers()` | Select users table | `FromUsers()` |
| `DuringPast(value, unit)` | Set time range | `DuringPast(7, Days)` |
| `On(date)` | Select specific date | `On("Feb 8, 2024")` |
| `FromTo(from, to)` | Set date range | `FromTo("2024-01-01", "2024-01-31")` |
| `With(clause)` | Add with clause | `With("web.errors during past 7d")` |
| `Include(clause)` | Add include clause | `Include("execution.crashes during past 7d")` |
| `Compute(alias, expr)` | Add compute | `Compute("total", "count()")` |
| `ComputeSum(alias, field)` | Compute sum | `ComputeSum("total", "crashes")` |
| `ComputeAvg(alias, field)` | Compute average | `ComputeAvg("avg", "memory")` |
| `Where(condition)` | Add filter | `Where("name == 'test'")` |
| `WhereEquals(field, value)` | Equals filter | `WhereEquals("platform", "Windows")` |
| `WhereIn(field, values)` | In filter | `WhereIn("type", []string{"laptop"})` |
| `List(fields...)` | Select fields | `List("name", "platform")` |
| `Sort(field, direction)` | Sort results | `Sort("name", SortAsc)` |
| `SortDesc(field)` | Sort descending | `SortDesc("crashes")` |
| `Limit(n)` | Limit rows | `Limit(100)` |
| `Summarize(alias, expr)` | Add summarize | `Summarize("total", "count()")` |
| `SummarizeBy(fields...)` | Group by | `SummarizeBy("platform")` |
| `Build()` | Build query string | `Build()` |
| `Validate()` | Validate query | `Validate()` |

## Time Selection

### TimeSelection Methods

| Method | Description | Example |
|--------|-------------|---------|
| `DuringPast(value, unit)` | Past period | `DuringPast(7, Days)` |
| `From(date)` | Start date | `From("2024-01-01")` |
| `To(date)` | End date | `To("2024-01-31")` |
| `On(date)` | Specific date | `On("Feb 8, 2024")` |
| `FromRelative(value, unit)` | Relative start | `FromRelative(21, Days)` |
| `ToRelative(value, unit)` | Relative end | `ToRelative(13, Days)` |
| `ByHighResolution()` | 30s resolution | `ByHighResolution()` |

### Time Constants

| Constant | Value |
|----------|-------|
| `Past15Minutes` | `"during past 15min"` |
| `Past30Minutes` | `"during past 30min"` |
| `Past1Hour` | `"during past 1h"` |
| `Past24Hours` | `"during past 24h"` |
| `Past7Days` | `"during past 7d"` |
| `Past30Days` | `"during past 30d"` |
| `Yesterday` | `"from 1d ago to 1d ago"` |

### Time Units

| Constant | Value |
|----------|-------|
| `Minutes` | `"min"` |
| `Hours` | `"h"` |
| `Days` | `"d"` |

## Result Sets

### V2ResultSet Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `Rows()` | `int` | Number of rows |
| `Fields()` | `[]string` | Field names |
| `Get(row, field)` | `any, error` | Get value |
| `GetString(row, field)` | `string, error` | Get string |
| `GetInt(row, field)` | `int64, error` | Get integer |
| `GetFloat(row, field)` | `float64, error` | Get float |
| `GetBool(row, field)` | `bool, error` | Get boolean |
| `GetRow(row)` | `map[string]any, error` | Get row |
| `HasField(field)` | `bool` | Check field exists |
| `Filter(fn)` | `[]map[string]any` | Filter rows |
| `Map(fn)` | `[]map[string]any` | Transform rows |
| `IterateRows(fn)` | `error` | Iterate with callback |
| `ToJSON()` | `[]byte, error` | Convert to JSON |

### V1ResultSet Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `Rows()` | `int` | Number of rows |
| `Columns()` | `int` | Number of columns |
| `Get(row, col)` | `any, error` | Get value |
| `GetString(row, col)` | `string, error` | Get string |
| `GetInt(row, col)` | `int64, error` | Get integer |
| `GetRow(row)` | `[]any, error` | Get row |
| `FindColumnIndex(name)` | `int, error` | Find column |
| `GetByColumnName(row, name)` | `any, error` | Get by name |
| `ToV2Format()` | `[]map[string]any` | Convert to V2 |
| `IterateRows(fn)` | `error` | Iterate with callback |

## Templates

### Available Templates

| Template | Description |
|----------|-------------|
| `DevicesWithCrashes` | Devices with application crashes |
| `DevicesWithHighMemoryUsage` | High memory usage devices |
| `DevicesByPlatform` | Device counts by OS |
| `DevicesWithSlowBootTime` | Slow boot time devices |
| `UsersWithWebErrors` | Users with web errors |
| `UsersWithPoorCollaborationQuality` | Poor call quality users |
| `ApplicationsWithHighErrorRate` | High error rate apps |
| `TopCrashingApplications` | Most crashed applications |
| `WebPageLoadPerformance` | Page load metrics |
| `NetworkConnectivityIssues` | Connectivity problems |
| `OverallDEXScore` | Overall DEX score |
| `DEXScoreByPlatform` | DEX by OS platform |
| `UsersWithLowDEXScore` | Low DEX score users |
| `DEXScoreImpactByComponent` | Score impact analysis |
| `DevicesWithSystemCrashes` | System crash monitoring |
| `BinariesWithHighCrashRate` | High crash binaries |
| `WorkflowExecutionSuccess` | Workflow success metrics |
| `RemoteActionSavingsEstimate` | Cost savings estimate |

## Export Workflow

### ExportOptions Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Format` | `string` | `"csv"` | Export format |
| `PollInterval` | `time.Duration` | `5s` | Poll frequency |
| `Timeout` | `time.Duration` | `10m` | Max wait time |
| `OnProgress` | `func` | `nil` | Progress callback |
| `OnStatusChange` | `func` | `nil` | Status change callback |

### Export Methods

| Method | Description |
|--------|-------------|
| `ExportWorkflow(ctx, req, opts)` | Complete export workflow |
| `ExportToCSV(ctx, queryID)` | Simple CSV export |
| `ExportToJSON(ctx, queryID)` | Simple JSON export |
| `ExportWithProgress(ctx, queryID, format, fn)` | Export with progress |
| `IsExportReady(ctx, exportID)` | Check if ready |
| `GetExportProgress(ctx, exportID)` | Get progress string |

## Data Model Constants

### Tables

| Constant | Value |
|----------|-------|
| `TableDevices` | `"devices"` |
| `TableUsers` | `"users"` |
| `TableApplications` | `"applications"` |
| `TableBinaries` | `"binaries"` |
| `TableExecutionCrashes` | `"execution.crashes"` |
| `TableWebPageViews` | `"web.page_views"` |
| `TableWebErrors` | `"web.errors"` |
| `TableDexScores` | `"dex.scores"` |

### Common Fields

| Category | Constants |
|----------|-----------|
| Device | `FieldDeviceName`, `FieldDeviceEntity` |
| OS | `FieldOSName`, `FieldOSPlatform`, `FieldOSVersion` |
| Hardware | `FieldHardwareType`, `FieldHardwareManufacturer` |
| User | `FieldUserName`, `FieldUsername` |
| Application | `FieldApplicationName`, `FieldApplicationVersion` |
| Binary | `FieldBinaryName`, `FieldBinaryVersion` |

### Common Values

| Category | Constants |
|----------|-----------|
| Platforms | `PlatformWindows`, `PlatformMacOS`, `PlatformLinux` |
| Hardware Types | `HardwareTypeLaptop`, `HardwareTypeDesktop`, `HardwareTypeVirtual` |
| User Types | `UserTypeLocalUser`, `UserTypeLocalAdmin` |
| Experience Levels | `ExperienceLevelGood`, `ExperienceLevelFrustrating` |

## Operators

### Comparison Operators

| Operator | Type | Description |
|----------|------|-------------|
| `OpEquals` | `==` | Equals |
| `OpNotEquals` | `!=` | Not equals |
| `OpGreater` | `>` | Greater than |
| `OpLess` | `<` | Less than |
| `OpGreaterEqual` | `>=` | Greater or equal |
| `OpLessEqual` | `<=` | Less or equal |
| `OpIn` | `in` | In list |
| `OpNotIn` | `!in` | Not in list |
| `OpContains` | `contains` | Contains |
| `OpNotContains` | `!contains` | Not contains |

### Logical Operators

| Operator | Description |
|----------|-------------|
| `LogicalAnd` | AND condition |
| `LogicalOr` | OR condition |

### Sort Directions

| Constant | Value |
|----------|-------|
| `SortAsc` | `"asc"` |
| `SortDesc` | `"desc"` |

## Functions

### Aggregate Functions

| Function | Description |
|----------|-------------|
| `FuncSum` | Sum of values |
| `FuncAvg` | Average value |
| `FuncCount` | Count records |
| `FuncMin` | Minimum value |
| `FuncMax` | Maximum value |
| `FuncLast` | Last value |
| `FuncCountIf` | Conditional count |
| `FuncSumIf` | Conditional sum |
| `FuncP95` | 95th percentile |
| `FuncP05` | 5th percentile |

### DateTime Functions

| Function | Description |
|----------|-------------|
| `FuncTimeElapsed` | Time since event |
| `FuncHour` | Extract hour |
| `FuncDay` | Extract day |
| `FuncDayOfWeek` | Day of week |

### Format Functions

| Function | Description |
|----------|-------------|
| `FormatEnergy` | Energy units |
| `FormatWeight` | Weight units |
| `FormatCurrency` | Currency format |
| `FormatPercent` | Percentage |
| `FormatBitrate` | Bitrate units |

## Validation

### QueryValidator Methods

| Method | Description |
|--------|-------------|
| `ValidateQuery(query)` | Full validation |
| `ValidateTableName(table)` | Table format |
| `ValidateTimeSelection(sel)` | Time syntax |
| `ValidateWhereClause(clause)` | Where syntax |
| `ValidateComments(query)` | Comment balance |
| `ValidateOperatorUsage(field, op, value)` | Operator compatibility |

### Validation Functions

| Function | Description |
|----------|-------------|
| `ValidateNQLQuery(query)` | Validate query string |
| `ValidateNQLQueryDetailed(query)` | Get all errors |

## Metadata

### ExecutionMetadata Fields

| Field | Type | Description |
|-------|------|-------------|
| `QueryID` | `string` | Query identifier |
| `ExecutedQuery` | `string` | Actual query executed |
| `RowsReturned` | `int64` | Rows in result |
| `ExecutionTime` | `time.Time` | When executed |
| `ResponseDuration` | `time.Duration` | Response time |
| `ResponseSize` | `int64` | Response bytes |
| `ResponseStatus` | `int` | HTTP status |
| `Headers` | `map[string][]string` | Response headers |

### Metadata Functions

| Function | Description |
|----------|-------------|
| `GetV1Metadata(resp, apiResp)` | Extract V1 metadata |
| `GetV2Metadata(resp, apiResp)` | Extract V2 metadata |
| `GetExportMetadata(resp, apiResp)` | Extract export metadata |

### ExecutionMetadata Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `GetRateLimitInfo()` | `*RateLimitInfo` | Rate limit details |
| `TimeSinceExecution()` | `time.Duration` | Time since query |
| `String()` | `string` | Format as string |

## Constants Reference

### Export Status

| Constant | Value |
|----------|-------|
| `ExportStatusSubmitted` | `"SUBMITTED"` |
| `ExportStatusInProgress` | `"IN_PROGRESS"` |
| `ExportStatusCompleted` | `"COMPLETED"` |
| `ExportStatusError` | `"ERROR"` |

### Export Format

| Constant | Value |
|----------|-------|
| `ExportFormatCSV` | `"csv"` |
| `ExportFormatJSON` | `"json"` |

### Time Granularity

| Constant | Value |
|----------|-------|
| `Granularity15Min` | `"15 min"` |
| `Granularity30Min` | `"30 min"` |
| `Granularity1Hour` | `"1 h"` |
| `Granularity1Day` | `"1 d"` |
| `Granularity7Days` | `"7 d"` |

## Quick Reference

### Complete Query Example

```go
query := nql.NewQueryBuilder().
    Comment("Device crash analysis").
    FromDevices().
    DuringPast(7, nql.Days).
    With("execution.crashes during past 7d").
    WhereEquals(nql.FieldOSPlatform, nql.PlatformWindows).
    WhereIn(nql.FieldHardwareType, []string{nql.HardwareTypeLaptop, nql.HardwareTypeDesktop}).
    ComputeSum("total_crashes", nql.FieldNumberOfCrashes).
    WhereGreaterEqual("total_crashes", "5").
    List(
        nql.FieldDeviceName,
        nql.FieldOSName,
        "total_crashes",
    ).
    SortDesc("total_crashes").
    Limit(20).
    Build()
```

### Complete Export Example

```go
opts := nql.DefaultExportOptions().
    WithFormat(nql.ExportFormatCSV).
    WithPollInterval(5 * time.Second).
    WithTimeout(15 * time.Minute).
    WithOnProgress(func(status string, elapsed time.Duration) {
        log.Printf("[%v] %s", elapsed, status)
    })

result, err := nqlService.ExportWorkflow(ctx, &nql.ExportRequest{
    QueryID: "#large_query",
}, opts)

if err != nil {
    log.Fatal(err)
}

os.WriteFile("export.csv", result.Data, 0644)
fmt.Printf("Saved %s\n", result.SizeFormatted())
```

### Complete Result Processing Example

```go
result, apiResp, err := nqlService.ExecuteNQLV2(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Create result set
resultSet := nql.NewV2ResultSet(result)

// Extract metadata
metadata := nql.GetV2Metadata(result, apiResp)
fmt.Printf("Executed in: %v\n", metadata.ResponseDuration)

// Process results
resultSet.IterateRows(func(row int, data map[string]any) error {
    deviceName, _ := resultSet.GetString(row, "device.name")
    fmt.Printf("Device: %s\n", deviceName)
    return nil
})

// Filter results
windowsDevices := resultSet.Filter(func(row map[string]any) bool {
    platform, _ := row[nql.FieldOSPlatform].(string)
    return platform == nql.PlatformWindows
})

// Convert to JSON
jsonData, _ := resultSet.ToJSON()
os.WriteFile("results.json", jsonData, 0644)
```

## Next Steps

- [NQL Query Building Guide](../guides/nql-query-building.md)
- [NQL Result Processing Guide](../guides/nql-result-processing.md)
- [NQL Export Workflow Guide](../guides/nql-export-workflow.md)
- [NQL Templates Guide](../guides/nql-templates.md)
- [NQL Best Practices Guide](../guides/nql-best-practices.md)
