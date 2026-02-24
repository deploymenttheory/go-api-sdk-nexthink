# NQL Examples

Complete examples demonstrating the Nexthink Query Language (NQL) service.

## Quick Start

Examples are numbered in recommended learning order. Start with `01_ExecuteBasics` and progress through each example.

## Progressive Examples (Recommended Path)

### [01_ExecuteBasics](01_ExecuteBasics/main.go)
**Start here** - Learn the fundamentals:
- Creating a Nexthink client
- Executing queries with V1 and V2 APIs
- Understanding response format differences

**Key Concepts**: API basics, V1 vs V2, synchronous execution

---

### [02_QueryBuilder](02_QueryBuilder/main.go)
Build queries programmatically:
- Fluent API for query construction
- Type-safe method chaining
- Query validation
- Using constants for fields and tables

**Key Concepts**: Query Builder, validation, type safety

---

### [03_Templates](03_Templates/main.go)
Use pre-built query templates:
- Accessing the template library (18+ templates)
- Generating queries for common scenarios
- Converting templates to API requests with `ToRequest()`
- Template categories (device health, user experience, DEX scores)

**Key Concepts**: Templates, ToRequest pattern, common queries

---

### [04_TimeAndConstants](04_TimeAndConstants/main.go)
Master SDK helpers and constants:
- Time selection builders and predefined constants
- Data model constants (tables, fields, values)
- Operators and functions
- Time granularity for aggregations

**Key Concepts**: Constants, time helpers, type safety

---

### [05_ResultProcessing](05_ResultProcessing/main.go)
Process results safely:
- Using `ExecuteV2WithResultSet` convenience method
- Type-safe data access (GetString, GetInt, GetFloat)
- Filtering and transforming results
- Metadata extraction

**Key Concepts**: Result sets, type-safe access, metadata

---

### [06_ExportBasics](06_ExportBasics/main.go)
Handle large datasets manually:
- Starting export operations
- Polling for completion status
- Downloading exported data
- Understanding the 3-step export workflow

**Key Concepts**: Asynchronous exports, polling, manual workflow

---

### [07_ExportWorkflow](07_ExportWorkflow/main.go)
Simplified export automation:
- One-line export methods (`ExportToCSV`, `ExportToJSON`)
- Progress tracking with callbacks
- Custom export options
- Automatic polling and download

**Key Concepts**: Automated exports, progress tracking, convenience methods

---

### [08_Integration](08_Integration/main.go)
**Complete workflow** - Everything together:
- Combining templates, builders, and result processing
- Real-world monitoring scenario
- Best practices in action

**Key Concepts**: End-to-end workflow, production patterns

---

## Legacy Examples (For Reference)

These examples demonstrate the original low-level API:

- **[ExecuteNQLV1](ExecuteNQLV1/main.go)** - V1 API execution
- **[ExecuteNQLV2](ExecuteNQLV2/main.go)** - V2 API execution
- **[StartNQLExport](StartNQLExport/main.go)** - Start export operation
- **[GetNQLExportStatus](GetNQLExportStatus/main.go)** - Check export status
- **[WaitForNQLExport](WaitForNQLExport/main.go)** - Poll for export completion

**Note**: The numbered examples (01-08) are the recommended approach and demonstrate SDK enhancements.

## Running Examples

### Prerequisites

Set environment variables:
```bash
export NEXTHINK_CLIENT_ID="your-client-id"
export NEXTHINK_CLIENT_SECRET="your-client-secret"
export NEXTHINK_INSTANCE="your-instance"
export NEXTHINK_REGION="us"  # or "eu"
```

### For Query Execution Examples

You need a query created in Nexthink admin:
```bash
export NEXTHINK_QUERY_ID="#your_query_id"
```

### For Export Examples

Optionally set a separate query for exports:
```bash
export NEXTHINK_EXPORT_QUERY_ID="#your_large_query"
```

### Running an Example

```bash
cd 01_ExecuteBasics
go run main.go
```

## Learning Path

1. **Foundations** (Examples 01-02)
   - Start with basic execution
   - Learn query builder

2. **Construction** (Examples 03-04)
   - Use pre-built templates
   - Master helpers and constants

3. **Processing** (Example 05)
   - Type-safe result handling
   - Metadata extraction

4. **Advanced** (Examples 06-07)
   - Manual export workflow
   - Simplified export automation

5. **Integration** (Example 08)
   - Complete end-to-end workflow
   - Production-ready patterns

## Quick Reference

### Execute a Query
```go
resultSet, _, err := nqlService.ExecuteV2WithResultSet(ctx, &nql.ExecuteRequest{
    QueryID: "#my_query",
})
```

### Build a Query
```go
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    List("device.name").
    Build()
```

### Use a Template
```go
template := nql.NewTemplates().DevicesWithCrashes("during past 7d", "")
req := template.ToRequest("#device_crashes")
result, _, err := nqlService.ExecuteNQLV2(ctx, req)
```

### Export Large Data
```go
result, err := nqlService.ExportToCSV(ctx, "#large_query")
os.WriteFile("export.csv", result.Data, 0644)
```

## Documentation

- [NQL Query Building Guide](../../../docs/guides/nql-query-building.md)
- [NQL Result Processing Guide](../../../docs/guides/nql-result-processing.md)
- [NQL Export Workflow Guide](../../../docs/guides/nql-export-workflow.md)
- [NQL Templates Guide](../../../docs/guides/nql-templates.md)
- [NQL Best Practices](../../../docs/guides/nql-best-practices.md)
- [NQL API Reference](../../../docs/reference/nql-reference.md)
