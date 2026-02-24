package nql

// Operators and functions for NQL queries
// Provides type-safe access to NQL operators, comparison operators, and aggregate functions

// =============================================================================
// Comparison Operators
// =============================================================================

// Operator represents an NQL comparison operator
type Operator string

const (
	// Equality Operators
	OpEquals    Operator = "=="
	OpEqualsAlt Operator = "="  // Alternative syntax for ==
	OpNotEquals Operator = "!="
	
	// Relational Operators
	OpGreater      Operator = ">"
	OpLess         Operator = "<"
	OpGreaterEqual Operator = ">="
	OpLessEqual    Operator = "<="
	
	// Membership Operators
	OpIn         Operator = "in"
	OpNotIn      Operator = "!in"
	OpContains   Operator = "contains"
	OpNotContains Operator = "!contains"
)

// String returns the string representation of the operator
func (o Operator) String() string {
	return string(o)
}

// =============================================================================
// Logical/Bitwise Operators
// =============================================================================

// LogicalOperator represents an NQL logical operator
type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "and"
	LogicalOr  LogicalOperator = "or"
)

// String returns the string representation of the logical operator
func (lo LogicalOperator) String() string {
	return string(lo)
}

// =============================================================================
// Arithmetic Operators
// =============================================================================

// ArithmeticOperator represents an NQL arithmetic operator
type ArithmeticOperator string

const (
	ArithmeticAdd      ArithmeticOperator = "+"
	ArithmeticSubtract ArithmeticOperator = "-"
	ArithmeticMultiply ArithmeticOperator = "*"
	ArithmeticDivide   ArithmeticOperator = "/"
)

// String returns the string representation of the arithmetic operator
func (ao ArithmeticOperator) String() string {
	return string(ao)
}

// =============================================================================
// Aggregate Functions
// =============================================================================

// AggregateFunc represents an NQL aggregate function
type AggregateFunc string

const (
	// Basic Aggregates
	FuncSum   AggregateFunc = "sum"
	FuncAvg   AggregateFunc = "avg"
	FuncCount AggregateFunc = "count"
	FuncMin   AggregateFunc = "min"
	FuncMax   AggregateFunc = "max"
	FuncLast  AggregateFunc = "last"
	
	// Conditional Aggregates
	FuncCountIf AggregateFunc = "countif"
	FuncSumIf   AggregateFunc = "sumif"
	
	// Percentile Functions
	FuncP95 AggregateFunc = "p95"
	FuncP05 AggregateFunc = "p05"
)

// String returns the string representation of the aggregate function
func (af AggregateFunc) String() string {
	return string(af)
}

// =============================================================================
// Datetime Functions
// =============================================================================

// DateTimeFunc represents an NQL datetime function
type DateTimeFunc string

const (
	FuncTimeElapsed DateTimeFunc = "time_elapsed"
	FuncHour        DateTimeFunc = "hour"
	FuncDay         DateTimeFunc = "day"
	FuncDayOfWeek   DateTimeFunc = "day_of_week"
)

// String returns the string representation of the datetime function
func (dtf DateTimeFunc) String() string {
	return string(dtf)
}

// =============================================================================
// Format Functions
// =============================================================================

// FormatFunc represents the NQL as() format function
type FormatFunc string

const (
	FormatEnergy   FormatFunc = "energy"
	FormatWeight   FormatFunc = "weight"
	FormatCurrency FormatFunc = "currency"
	FormatPercent  FormatFunc = "percent"
	FormatBitrate  FormatFunc = "bitrate"
)

// String returns the string representation of the format function
func (ff FormatFunc) String() string {
	return string(ff)
}

// Currency codes for format function
const (
	CurrencyCAD = "CAD"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyGBP = "GBP"
	CurrencyCHF = "CHF"
)

// =============================================================================
// Sort Directions
// =============================================================================

// SortDirection represents the sort direction
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// String returns the string representation of the sort direction
func (sd SortDirection) String() string {
	return string(sd)
}

// =============================================================================
// Aggregated Metrics (field suffixes)
// =============================================================================

// These are suffixes that can be appended to metric fields
const (
	MetricSuffixAvg   = ".avg"
	MetricSuffixSum   = ".sum"
	MetricSuffixCount = ".count"
	MetricSuffixMin   = ".min"
	MetricSuffixMax   = ".max"
)

// =============================================================================
// Helper Functions
// =============================================================================

// IsComparisonOperator checks if the operator is a valid comparison operator
func IsComparisonOperator(op Operator) bool {
	switch op {
	case OpEquals, OpEqualsAlt, OpNotEquals,
		OpGreater, OpLess, OpGreaterEqual, OpLessEqual,
		OpIn, OpNotIn, OpContains, OpNotContains:
		return true
	default:
		return false
	}
}

// IsArithmeticOperator checks if the operator is a valid arithmetic operator
func IsArithmeticOperator(op ArithmeticOperator) bool {
	switch op {
	case ArithmeticAdd, ArithmeticSubtract, ArithmeticMultiply, ArithmeticDivide:
		return true
	default:
		return false
	}
}

// IsLogicalOperator checks if the operator is a valid logical operator
func IsLogicalOperator(op LogicalOperator) bool {
	switch op {
	case LogicalAnd, LogicalOr:
		return true
	default:
		return false
	}
}

// IsAggregateFunction checks if the function is a valid aggregate function
func IsAggregateFunction(fn AggregateFunc) bool {
	switch fn {
	case FuncSum, FuncAvg, FuncCount, FuncMin, FuncMax, FuncLast,
		FuncCountIf, FuncSumIf, FuncP95, FuncP05:
		return true
	default:
		return false
	}
}

// IsDateTimeFunction checks if the function is a valid datetime function
func IsDateTimeFunction(fn DateTimeFunc) bool {
	switch fn {
	case FuncTimeElapsed, FuncHour, FuncDay, FuncDayOfWeek:
		return true
	default:
		return false
	}
}

// IsFormatFunction checks if the function is a valid format function
func IsFormatFunction(fn FormatFunc) bool {
	switch fn {
	case FormatEnergy, FormatWeight, FormatCurrency, FormatPercent, FormatBitrate:
		return true
	default:
		return false
	}
}
