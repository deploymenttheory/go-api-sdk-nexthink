package nql

import (
	"fmt"
	"strings"
)

// Query builder provides a fluent API for constructing NQL queries programmatically
// Enables type-safe, IDE-friendly query construction with validation

// =============================================================================
// QueryBuilder
// =============================================================================

// QueryBuilder provides a fluent API for building NQL queries
type QueryBuilder struct {
	table            string
	timeSelection    string
	withClauses      []string
	includeClauses   []string
	computeClauses   []string
	whereClauses     []string
	listFields       []string
	sortField        string
	sortDirection    SortDirection
	limitValue       int
	summarizeClauses []string
	summarizeBy      []string
	comments         []string
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// =============================================================================
// Table Selection
// =============================================================================

// From sets the table/namespace for the query
// Example: From("devices") or From("execution.crashes")
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.table = table
	return qb
}

// FromDevices is a shorthand for From(TableDevices)
func (qb *QueryBuilder) FromDevices() *QueryBuilder {
	return qb.From(TableDevices)
}

// FromUsers is a shorthand for From(TableUsers)
func (qb *QueryBuilder) FromUsers() *QueryBuilder {
	return qb.From(TableUsers)
}

// FromApplications is a shorthand for From(TableApplications)
func (qb *QueryBuilder) FromApplications() *QueryBuilder {
	return qb.From(TableApplications)
}

// FromBinaries is a shorthand for From(TableBinaries)
func (qb *QueryBuilder) FromBinaries() *QueryBuilder {
	return qb.From(TableBinaries)
}

// =============================================================================
// Time Selection
// =============================================================================

// During sets the time selection clause
// Example: During("past 7d") or During(Past7Days)
func (qb *QueryBuilder) During(timeSelection string) *QueryBuilder {
	qb.timeSelection = timeSelection
	return qb
}

// DuringPast sets a "during past" time selection
// Example: DuringPast(7, Days)
func (qb *QueryBuilder) DuringPast(value int, unit TimeUnit) *QueryBuilder {
	qb.timeSelection = fmt.Sprintf("during past %d%s", value, unit)
	return qb
}

// On sets an "on" time selection for a specific date
// Example: On("Feb 8, 2024")
func (qb *QueryBuilder) On(date string) *QueryBuilder {
	qb.timeSelection = fmt.Sprintf("on %s", date)
	return qb
}

// FromTo sets a "from ... to ..." time selection
// Example: FromTo("2024-01-01", "2024-01-31")
func (qb *QueryBuilder) FromTo(from, to string) *QueryBuilder {
	qb.timeSelection = fmt.Sprintf("from %s to %s", from, to)
	return qb
}

// =============================================================================
// With Clause (joins with filtering)
// =============================================================================

// With adds a "with" clause (joins event table, filters to objects with events)
// Example: With("web.errors during past 7d")
func (qb *QueryBuilder) With(clause string) *QueryBuilder {
	qb.withClauses = append(qb.withClauses, clause)
	return qb
}

// WithTable adds a "with" clause for a table with time selection
// Example: WithTable(TableWebErrors, "during past 7d")
func (qb *QueryBuilder) WithTable(table, timeSelection string) *QueryBuilder {
	clause := table
	if timeSelection != "" {
		clause += " " + timeSelection
	}
	qb.withClauses = append(qb.withClauses, clause)
	return qb
}

// =============================================================================
// Include Clause (joins without filtering)
// =============================================================================

// Include adds an "include" clause (joins event table, keeps all objects)
// Example: Include("execution.crashes during past 7d")
func (qb *QueryBuilder) Include(clause string) *QueryBuilder {
	qb.includeClauses = append(qb.includeClauses, clause)
	return qb
}

// IncludeTable adds an "include" clause for a table with time selection
// Example: IncludeTable(TableExecutionCrashes, "during past 7d")
func (qb *QueryBuilder) IncludeTable(table, timeSelection string) *QueryBuilder {
	clause := table
	if timeSelection != "" {
		clause += " " + timeSelection
	}
	qb.includeClauses = append(qb.includeClauses, clause)
	return qb
}

// =============================================================================
// Compute Clause
// =============================================================================

// Compute adds a compute clause
// Example: Compute("total_crashes", "count()")
func (qb *QueryBuilder) Compute(alias, expression string) *QueryBuilder {
	qb.computeClauses = append(qb.computeClauses, fmt.Sprintf("%s = %s", alias, expression))
	return qb
}

// ComputeCount adds a count() compute clause
// Example: ComputeCount("total_count")
func (qb *QueryBuilder) ComputeCount(alias string) *QueryBuilder {
	return qb.Compute(alias, "count()")
}

// ComputeSum adds a sum() compute clause
// Example: ComputeSum("total_crashes", "number_of_crashes")
func (qb *QueryBuilder) ComputeSum(alias, field string) *QueryBuilder {
	return qb.Compute(alias, fmt.Sprintf("%s.sum()", field))
}

// ComputeAvg adds an avg() compute clause
// Example: ComputeAvg("avg_memory", "free_memory")
func (qb *QueryBuilder) ComputeAvg(alias, field string) *QueryBuilder {
	return qb.Compute(alias, fmt.Sprintf("%s.avg()", field))
}

// ComputeMax adds a max() compute clause
func (qb *QueryBuilder) ComputeMax(alias, field string) *QueryBuilder {
	return qb.Compute(alias, fmt.Sprintf("%s.max()", field))
}

// ComputeMin adds a min() compute clause
func (qb *QueryBuilder) ComputeMin(alias, field string) *QueryBuilder {
	return qb.Compute(alias, fmt.Sprintf("%s.min()", field))
}

// ComputeLast adds a last() compute clause
func (qb *QueryBuilder) ComputeLast(alias, field string) *QueryBuilder {
	return qb.Compute(alias, fmt.Sprintf("%s.last()", field))
}

// =============================================================================
// Where Clause
// =============================================================================

// Where adds a where clause
// Example: Where("binary.name == \"outlook.exe\"")
func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.whereClauses = append(qb.whereClauses, condition)
	return qb
}

// WhereEquals adds a where clause with equals operator
// Example: WhereEquals("binary.name", "outlook.exe")
func (qb *QueryBuilder) WhereEquals(field, value string) *QueryBuilder {
	// Quote string values
	quotedValue := quoteValue(value)
	return qb.Where(fmt.Sprintf("%s == %s", field, quotedValue))
}

// WhereNotEquals adds a where clause with not equals operator
func (qb *QueryBuilder) WhereNotEquals(field, value string) *QueryBuilder {
	quotedValue := quoteValue(value)
	return qb.Where(fmt.Sprintf("%s != %s", field, quotedValue))
}

// WhereGreater adds a where clause with greater than operator
func (qb *QueryBuilder) WhereGreater(field, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s > %s", field, value))
}

// WhereLess adds a where clause with less than operator
func (qb *QueryBuilder) WhereLess(field, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s < %s", field, value))
}

// WhereGreaterEqual adds a where clause with >= operator
func (qb *QueryBuilder) WhereGreaterEqual(field, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s >= %s", field, value))
}

// WhereLessEqual adds a where clause with <= operator
func (qb *QueryBuilder) WhereLessEqual(field, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s <= %s", field, value))
}

// WhereIn adds a where clause with in operator
// Example: WhereIn("hardware.type", []string{"laptop", "desktop"})
func (qb *QueryBuilder) WhereIn(field string, values []string) *QueryBuilder {
	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = quoteValue(v)
	}
	return qb.Where(fmt.Sprintf("%s in [%s]", field, strings.Join(quotedValues, ", ")))
}

// WhereNotIn adds a where clause with !in operator
func (qb *QueryBuilder) WhereNotIn(field string, values []string) *QueryBuilder {
	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = quoteValue(v)
	}
	return qb.Where(fmt.Sprintf("%s !in [%s]", field, strings.Join(quotedValues, ", ")))
}

// WhereContains adds a where clause with contains operator
// Example: WhereContains("tags", "VDI")
func (qb *QueryBuilder) WhereContains(field, value string) *QueryBuilder {
	quotedValue := quoteValue(value)
	return qb.Where(fmt.Sprintf("%s contains %s", field, quotedValue))
}

// WhereNotContains adds a where clause with !contains operator
func (qb *QueryBuilder) WhereNotContains(field, value string) *QueryBuilder {
	quotedValue := quoteValue(value)
	return qb.Where(fmt.Sprintf("%s !contains %s", field, quotedValue))
}

// =============================================================================
// List Clause
// =============================================================================

// List adds fields to the list clause
// Example: List("name", "type", "version")
func (qb *QueryBuilder) List(fields ...string) *QueryBuilder {
	qb.listFields = append(qb.listFields, fields...)
	return qb
}

// =============================================================================
// Sort Clause
// =============================================================================

// Sort adds a sort clause
// Example: Sort("size", SortDesc)
func (qb *QueryBuilder) Sort(field string, direction SortDirection) *QueryBuilder {
	qb.sortField = field
	qb.sortDirection = direction
	return qb
}

// SortAsc adds an ascending sort clause
func (qb *QueryBuilder) SortAsc(field string) *QueryBuilder {
	return qb.Sort(field, SortAsc)
}

// SortDesc adds a descending sort clause
func (qb *QueryBuilder) SortDesc(field string) *QueryBuilder {
	return qb.Sort(field, SortDesc)
}

// =============================================================================
// Limit Clause
// =============================================================================

// Limit adds a limit clause
// Example: Limit(100)
func (qb *QueryBuilder) Limit(value int) *QueryBuilder {
	qb.limitValue = value
	return qb
}

// =============================================================================
// Summarize Clause
// =============================================================================

// Summarize adds a summarize clause
// Example: Summarize("total_devices", "count()")
func (qb *QueryBuilder) Summarize(alias, expression string) *QueryBuilder {
	qb.summarizeClauses = append(qb.summarizeClauses, fmt.Sprintf("%s = %s", alias, expression))
	return qb
}

// SummarizeCount adds a count() summarize clause
func (qb *QueryBuilder) SummarizeCount(alias string) *QueryBuilder {
	return qb.Summarize(alias, "count()")
}

// SummarizeSum adds a sum() summarize clause
func (qb *QueryBuilder) SummarizeSum(alias, field string) *QueryBuilder {
	return qb.Summarize(alias, fmt.Sprintf("%s.sum()", field))
}

// SummarizeAvg adds an avg() summarize clause
func (qb *QueryBuilder) SummarizeAvg(alias, field string) *QueryBuilder {
	return qb.Summarize(alias, fmt.Sprintf("%s.avg()", field))
}

// SummarizeBy adds grouping fields for summarize by
// Example: SummarizeBy("entity", "platform")
func (qb *QueryBuilder) SummarizeBy(fields ...string) *QueryBuilder {
	qb.summarizeBy = append(qb.summarizeBy, fields...)
	return qb
}

// SummarizeByTime adds a time-based grouping for summarize by
// Example: SummarizeByTime(Granularity1Day)
func (qb *QueryBuilder) SummarizeByTime(granularity TimeGranularity) *QueryBuilder {
	qb.summarizeBy = append(qb.summarizeBy, string(granularity))
	return qb
}

// =============================================================================
// Comments
// =============================================================================

// Comment adds a comment to the query
// Example: Comment("This query finds devices with crashes")
func (qb *QueryBuilder) Comment(text string) *QueryBuilder {
	qb.comments = append(qb.comments, text)
	return qb
}

// =============================================================================
// Build
// =============================================================================

// Build constructs the final NQL query string
func (qb *QueryBuilder) Build() string {
	var parts []string
	
	// Add comments at the beginning
	for _, comment := range qb.comments {
		parts = append(parts, fmt.Sprintf("/* %s */", comment))
	}
	
	// Table + time selection
	if qb.table != "" {
		clause := qb.table
		if qb.timeSelection != "" {
			clause += " " + qb.timeSelection
		}
		parts = append(parts, clause)
	}
	
	// With clauses
	for _, with := range qb.withClauses {
		parts = append(parts, "| with "+with)
	}
	
	// Include clauses
	for _, include := range qb.includeClauses {
		parts = append(parts, "| include "+include)
	}
	
	// Compute clauses
	for _, compute := range qb.computeClauses {
		parts = append(parts, "| compute "+compute)
	}
	
	// Where clauses
	for _, where := range qb.whereClauses {
		parts = append(parts, "| where "+where)
	}
	
	// List clause
	if len(qb.listFields) > 0 {
		parts = append(parts, "| list "+strings.Join(qb.listFields, ", "))
	}
	
	// Summarize clause
	if len(qb.summarizeClauses) > 0 {
		summarize := "| summarize " + strings.Join(qb.summarizeClauses, ", ")
		if len(qb.summarizeBy) > 0 {
			summarize += " by " + strings.Join(qb.summarizeBy, ", ")
		}
		parts = append(parts, summarize)
	}
	
	// Sort clause
	if qb.sortField != "" {
		parts = append(parts, fmt.Sprintf("| sort %s %s", qb.sortField, qb.sortDirection))
	}
	
	// Limit clause
	if qb.limitValue > 0 {
		parts = append(parts, fmt.Sprintf("| limit %d", qb.limitValue))
	}
	
	return strings.Join(parts, "\n")
}

// String returns the query string (alias for Build)
func (qb *QueryBuilder) String() string {
	return qb.Build()
}

// =============================================================================
// Helper Functions
// =============================================================================

// quoteValue quotes a string value for NQL queries
func quoteValue(value string) string {
	// If already quoted, return as-is
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return value
	}
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return value
	}
	
	// Quote the value
	return fmt.Sprintf("\"%s\"", value)
}

// =============================================================================
// Validation
// =============================================================================

// Validate performs basic validation on the query
func (qb *QueryBuilder) Validate() error {
	if qb.table == "" {
		return fmt.Errorf("table selection is required (use From())")
	}
	
	// Can't have both list and summarize
	if len(qb.listFields) > 0 && len(qb.summarizeClauses) > 0 {
		return fmt.Errorf("cannot use both list and summarize in the same query")
	}
	
	// Compute requires with or include
	if len(qb.computeClauses) > 0 && len(qb.withClauses) == 0 && len(qb.includeClauses) == 0 {
		return fmt.Errorf("compute clause requires a with or include clause")
	}
	
	return nil
}

// BuildAndValidate builds the query and validates it
func (qb *QueryBuilder) BuildAndValidate() (string, error) {
	if err := qb.Validate(); err != nil {
		return "", err
	}
	return qb.Build(), nil
}
