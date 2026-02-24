# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### NQL Service Enhancements

- **Query Builder**: Added fluent API for programmatic NQL query construction with type safety and validation
  - Method chaining for readable query construction
  - Support for all NQL constructs (table selection, time ranges, filters, aggregations)
  - Built-in validation with detailed error messages
  - IDE auto-completion support

- **Query Templates**: Added 18 pre-built query templates for common scenarios
  - Device health monitoring (crashes, memory usage, boot time)
  - User experience analysis (web errors, collaboration quality)
  - Application performance (error rates, crash analysis)
  - DEX score analysis (overall, by platform, low score users, component impact)
  - Network connectivity and performance monitoring
  - Workflow and remote action metrics

- **Result Set Processing**: Added type-safe result processing helpers
  - V1 and V2 result set wrappers with type-safe getters
  - Filter and map operations for data transformation
  - Row iteration utilities
  - V1 to V2 format conversion
  - JSON export capabilities

- **Export Workflow**: Added simplified workflow for large data exports
  - One-line export methods (ExportToCSV, ExportToJSON)
  - Progress tracking with customizable callbacks
  - Automatic polling and completion detection
  - Configurable timeouts and intervals
  - Human-readable result sizes and progress reporting

- **Data Model Constants**: Added comprehensive constants for type-safe query construction
  - 50+ table name constants
  - 100+ field name constants
  - 50+ value enumerations (platforms, hardware types, experience levels, etc.)
  - Time selection constants and helpers
  - Operator and function constants

- **Time Selection Helpers**: Added fluent API for time range specification
  - Predefined time constants (Past7Days, Past24Hours, etc.)
  - Relative and absolute time range builders
  - Time granularity constants for aggregations
  - High-resolution support for VDI data

- **Query Validation**: Added comprehensive client-side query validation
  - Syntax and structure validation
  - Operator compatibility checking
  - Comment balance validation
  - Detailed error reporting

- **Metadata Extraction**: Added helpers for execution and performance metrics
  - Query execution time tracking
  - Response duration and size metrics
  - Rate limit information extraction
  - Row count and status tracking

#### Examples

- Added QueryBuilder example demonstrating fluent query construction
- Added Templates example showing all 18 pre-built templates
- Added ResultSetProcessing example with type-safe data access patterns
- Added ExportWorkflow example with progress tracking
- Added ComprehensiveExample combining all enhancements

#### Documentation

- Added NQL Query Building Guide (complete guide to query builder)
- Added NQL Result Processing Guide (working with query results)
- Added NQL Export Workflow Guide (large data exports)
- Added NQL Templates Guide (pre-built query templates)
- Added NQL Best Practices Guide (optimization and patterns)
- Added NQL API Reference (complete reference documentation)
- Added NQL Enhancements summary document
- Updated README with comprehensive NQL enhancements section

### Changed

- Enhanced NQL service with developer-friendly features
- Updated README documentation to include new NQL capabilities

### Fixed

- N/A

## [1.1.0] - 2021-06-23

### Added

- Added x [@your_username](https://github.com/your_username)

### Changed

- Changed y [@your_username](https://github.com/your_username)

## [1.0.0] - 2021-06-20

### Added

- Inititated y [@your_username](https://github.com/your_username)
- Inititated z [@your_username](https://github.com/your_username)