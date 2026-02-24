# Go SDK for Nexthink API

[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/go-api-sdk-nexthink)](https://goreportcard.com/report/github.com/deploymenttheory/go-api-sdk-nexthink)
[![GoDoc](https://pkg.go.dev/badge/github.com/deploymenttheory/go-api-sdk-nexthink)](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-nexthink)
[![License](https://img.shields.io/github/license/deploymenttheory/go-api-sdk-nexthink)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/deploymenttheory/go-api-sdk-nexthink)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/deploymenttheory/go-api-sdk-nexthink)](https://github.com/deploymenttheory/go-api-sdk-nexthink/releases)
[![codecov](https://codecov.io/gh/deploymenttheory/go-api-sdk-nexthink/graph/badge.svg)](https://codecov.io/gh/deploymenttheory/go-api-sdk-nexthink)
[![Tests](https://github.com/deploymenttheory/go-api-sdk-nexthink/workflows/Tests/badge.svg)](https://github.com/deploymenttheory/go-api-sdk-nexthink/actions)
![Status: Alpha](https://img.shields.io/badge/status-alpha-orange)

A community Go client library for the [Nexthink API](https://doc.nexthink.com/Documentation/Nexthink/latest/APIAndIntegrations/IntroducingNexthinkAPI).


## Quick Start

Get started quickly with the SDK using the **[Quick Start Guide](docs/guides/quick-start.md)**, which includes:
- Installation instructions
- Your first API call
- Common operations (NQL queries, workflows, remote actions)
- Error handling patterns
- Response metadata access
- Links to configuration guides for production use

## Examples

The [examples directory](examples/nexthink/) contains complete working examples for all SDK features:

### Client Configuration
- [Basic Client Creation](examples/nexthink/_build_client/new_client/main.go)
- [Client with Structured Logging](examples/nexthink/_build_client/new_client_with_logger/main.go)

### NQL (Nexthink Query Language)

#### Basic Operations
- [Execute NQL Query V1](examples/nexthink/nql/ExecuteNQLV1/main.go)
- [Execute NQL Query V2](examples/nexthink/nql/ExecuteNQLV2/main.go)
- [Start NQL Export](examples/nexthink/nql/StartNQLExport/main.go)
- [Get NQL Export Status](examples/nexthink/nql/GetNQLExportStatus/main.go)
- [Wait for NQL Export (Full Workflow)](examples/nexthink/nql/WaitForNQLExport/main.go)

#### Advanced Features
- [Query Builder](examples/nexthink/nql/QueryBuilder/main.go) - Fluent API for building NQL queries
- [Templates](examples/nexthink/nql/Templates/main.go) - Pre-built queries for common scenarios
- [Result Set Processing](examples/nexthink/nql/ResultSetProcessing/main.go) - Type-safe data access and transformations
- [Export Workflow](examples/nexthink/nql/ExportWorkflow/main.go) - Simplified large data exports
- [Comprehensive Example](examples/nexthink/nql/ComprehensiveExample/main.go) - Full-featured example combining all enhancements

### Workflows
- [Trigger Workflow V1](examples/nexthink/workflows/TriggerWorkflowV1/main.go)
- [Trigger Workflow V2](examples/nexthink/workflows/TriggerWorkflowV2/main.go)
- [List Workflows](examples/nexthink/workflows/ListWorkflows/main.go)
- [Get Workflow Details](examples/nexthink/workflows/GetWorkflowDetails/main.go)

### Remote Actions
- [Trigger Remote Action](examples/nexthink/remote_actions/TriggerRemoteAction/main.go)
- [List Remote Actions](examples/nexthink/remote_actions/ListRemoteActions/main.go)
- [Get Remote Action Details](examples/nexthink/remote_actions/GetRemoteActionDetails/main.go)

### Enrichment
- [Enrich Fields](examples/nexthink/enrichment/EnrichFields/main.go)

### Campaigns
- [Trigger Campaign](examples/nexthink/campaigns/TriggerCampaign/main.go)

Each example includes a complete `main.go` with comments explaining the code.


## SDK Services

### Core Query and Data Services

- **NQL (Nexthink Query Language)**: Execute NQL queries, start and monitor data exports
  - V1 and V2 query execution
  - Asynchronous export operations with status tracking
  - Automatic waiting and polling for export completion
  - **Query Builder**: Fluent API for programmatic query construction
  - **Query Templates**: Pre-built queries for common scenarios
  - **Result Sets**: Type-safe data access and transformation helpers
  - **Export Workflow**: Simplified large data exports with progress tracking
  - **Data Model Constants**: Type-safe constants for tables, fields, and values
  - **Query Validation**: Client-side validation before API execution
  - **Metadata Extraction**: Detailed execution and performance metrics

### Automation and Orchestration

- **Workflows**: Trigger and manage Nexthink workflows
  - V1 and V2 workflow execution
  - List available workflows
  - Get detailed workflow information

- **Remote Actions**: Execute remote actions on endpoints
  - Trigger remote actions with parameters
  - List available remote actions
  - Get remote action details and schemas

### Data Enrichment and User Engagement

- **Enrichment**: Enrich Nexthink data with custom fields
  - Add custom device and user metadata
  - Update existing enrichment data

- **Campaigns**: Send targeted campaigns to users
  - Trigger campaigns with custom parameters
  - Multi-user campaign execution
  - Track campaign request status per user

## NQL Enhancements

The SDK includes comprehensive enhancements for working with NQL (Nexthink Query Language):

### Query Builder

Build NQL queries programmatically with a type-safe fluent API:

```go
import "github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/services/nql"

// Build a query using method chaining
query := nql.NewQueryBuilder().
    FromDevices().
    DuringPast(7, nql.Days).
    With("execution.crashes during past 7d").
    WhereEquals("binary.name", "outlook.exe").
    ComputeSum("total_crashes", "number_of_crashes").
    List("device.name", "total_crashes").
    SortDesc("total_crashes").
    Limit(20).
    Build()

// Validate before use
if err := qb.Validate(); err != nil {
    log.Fatal(err)
}
```

### Query Templates

Use pre-built query templates for common scenarios:

```go
templates := nql.NewTemplates()

// Get devices with crashes
query := templates.DevicesWithCrashes("during past 7d", "outlook.exe")

// Analyze DEX scores by platform
query := templates.DEXScoreByPlatform("during past 24h")

// Find users with low DEX scores
query := templates.UsersWithLowDEXScore(50, "during past 24h")
```

### Result Set Processing

Process query results with type-safe helpers:

```go
// Execute query
result, apiResp, err := nqlService.ExecuteNQLV2(ctx, req)

// Create result set
resultSet := nql.NewV2ResultSet(result)

// Type-safe data access
deviceName, err := resultSet.GetString(0, "device.name")
crashCount, err := resultSet.GetInt(0, "total_crashes")

// Iterate through results
resultSet.IterateRows(func(row int, data map[string]any) error {
    fmt.Printf("Device: %v\n", data)
    return nil
})

// Filter and transform
windowsDevices := resultSet.Filter(func(row map[string]any) bool {
    platform, _ := row["operating_system.platform"].(string)
    return platform == "Windows"
})

// Convert to JSON
jsonData, _ := resultSet.ToJSON()
```

### Export Workflow

Simplified workflow for large data exports:

```go
// Simple CSV export
result, err := nqlService.ExportToCSV(ctx, "#large_query")
os.WriteFile("export.csv", result.Data, 0644)

// Export with progress tracking
opts := nql.DefaultExportOptions().
    WithPollInterval(5 * time.Second).
    WithTimeout(15 * time.Minute).
    WithOnProgress(func(status string, elapsed time.Duration) {
        fmt.Printf("[%v] %s\n", elapsed, status)
    })

result, err := nqlService.ExportWorkflow(ctx, req, opts)
fmt.Printf("Exported %s in %v\n", result.SizeFormatted(), result.TotalDuration)
```

### Data Model Constants

Use constants for type-safe query construction:

```go
// Table constants
nql.TableDevices              // "devices"
nql.TableExecutionCrashes     // "execution.crashes"
nql.TableWebErrors            // "web.errors"

// Field constants
nql.FieldDeviceName           // "device.name"
nql.FieldOSPlatform           // "operating_system.platform"
nql.FieldHardwareType         // "hardware.type"

// Value constants
nql.PlatformWindows           // "Windows"
nql.HardwareTypeLaptop        // "laptop"
nql.ExperienceLevelGood       // "good"

// Time selection helpers
nql.Past7Days                 // "during past 7d"
nql.Past24Hours               // "during past 24h"
```

### Comprehensive Documentation

- **[NQL Query Building Guide](docs/guides/nql-query-building.md)** - Complete guide to the query builder
- **[NQL Result Processing Guide](docs/guides/nql-result-processing.md)** - Working with query results
- **[NQL Export Workflow Guide](docs/guides/nql-export-workflow.md)** - Large data exports
- **[NQL Templates Guide](docs/guides/nql-templates.md)** - Pre-built query templates
- **[NQL Best Practices Guide](docs/guides/nql-best-practices.md)** - Optimization and patterns
- **[NQL API Reference](docs/reference/nql-reference.md)** - Complete API reference

## HTTP Client Configuration

The SDK includes a powerful HTTP client with production-ready configuration options:

- **[Authentication](docs/guides/authentication.md)** - OAuth2 token management with automatic refresh
- **[Timeouts & Retries](docs/guides/timeouts-retries.md)** - Configurable timeouts and automatic retry logic
- **[TLS/SSL Configuration](docs/guides/tls-configuration.md)** - Custom certificates, mutual TLS, and security settings
- **[Proxy Support](docs/guides/proxy.md)** - HTTP/HTTPS/SOCKS5 proxy configuration
- **[Custom Headers](docs/guides/custom-headers.md)** - Global and per-request header management
- **[Structured Logging](docs/guides/logging.md)** - Integration with zap for production logging
- **[OpenTelemetry Tracing](docs/guides/opentelemetry.md)** - Distributed tracing and observability
- **[Debug Mode](docs/guides/debugging.md)** - Detailed request/response inspection

## Configuration Options

The SDK client supports extensive configuration through functional options. Below is the complete list of available configuration options grouped by category.

### Basic Configuration

```go
client.WithAPIVersion("v1")              // Set API version
client.WithBaseURL("https://...")        // Custom base URL
client.WithTimeout(30*time.Second)       // Request timeout
client.WithRetryCount(3)                 // Number of retry attempts
```

### TLS/Security

```go
client.WithMinTLSVersion(tls.VersionTLS12)                    // Minimum TLS version
client.WithTLSClientConfig(tlsConfig)                         // Custom TLS configuration
client.WithRootCertificates("/path/to/ca.pem")                // Custom CA certificates
client.WithRootCertificateFromString(caPEM)                   // CA certificate from string
client.WithClientCertificate("/path/cert.pem", "/path/key.pem") // Client certificate (mTLS)
client.WithClientCertificateFromString(certPEM, keyPEM)       // Client cert from string
client.WithInsecureSkipVerify()                               // Skip cert verification (dev only!)
```

### Network

```go
client.WithProxy("http://proxy:8080")    // HTTP/HTTPS/SOCKS5 proxy
client.WithTransport(customTransport)    // Custom HTTP transport
```

### Headers

```go
client.WithUserAgent("MyApp/1.0")                      // Set User-Agent header
client.WithCustomAgent("MyApp", "1.0")                 // Custom agent with version
client.WithGlobalHeader("X-Custom-Header", "value")    // Add single global header
client.WithGlobalHeaders(map[string]string{...})       // Add multiple global headers
```

### Observability

```go
client.WithLogger(zapLogger)            // Structured logging with zap
client.WithTracing(otelConfig)          // OpenTelemetry distributed tracing
client.WithDebug()                      // Enable debug mode (dev only!)
```

### Example: Production Configuration

```go
import (
    "crypto/tls"
    "time"
    "go.uber.org/zap"
    "github.com/deploymenttheory/go-api-sdk-nexthink/nexthink/client"
)

logger, _ := zap.NewProduction()

// Initialize client with OAuth2 credentials
apiClient, err := client.NewClient(
    "your-client-id",
    "your-client-secret",
    "your-instance",      // e.g., "company"
    "us",                 // Region: "us" or "eu"
    client.WithTimeout(30*time.Second),
    client.WithRetryCount(3),
    client.WithLogger(logger),
    client.WithMinTLSVersion(tls.VersionTLS12),
    client.WithGlobalHeader("X-Application-Name", "MyITApp"),
)
if err != nil {
    logger.Fatal("Failed to create Nexthink client", zap.Error(err))
}

// Access services
nqlService := apiClient.NQL()
workflowsService := apiClient.Workflows()
remoteActionsService := apiClient.RemoteActions()
```

See the [configuration guides](docs/guides/) for detailed documentation on each option.


## Documentation

- [Nexthink API Documentation](https://doc.nexthink.com/Documentation/Nexthink/latest/APIAndIntegrations/IntroducingNexthinkAPI)
- [GoDoc](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-nexthink)

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/deploymenttheory/go-api-sdk-nexthink/issues)
- **Documentation**: [API Docs](https://doc.nexthink.com/Documentation/Nexthink/latest/APIAndIntegrations/IntroducingNexthinkAPI)

## Disclaimer

This is an unofficial SDK and is not affiliated with or endorsed by Nexthink.
