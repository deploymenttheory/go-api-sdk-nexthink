package nql

import (
	"fmt"
	"regexp"
	"strings"
)

// Query validation provides comprehensive validation of NQL queries
// Validates syntax, operators, functions, and common patterns

// =============================================================================
// QueryValidator
// =============================================================================

// QueryValidator provides comprehensive query validation
type QueryValidator struct{}

// NewQueryValidator creates a new query validator
func NewQueryValidator() *QueryValidator {
	return &QueryValidator{}
}

// =============================================================================
// Full Query Validation
// =============================================================================

// ValidateQuery performs comprehensive query validation
func (qv *QueryValidator) ValidateQuery(query string) error {
	if query == "" {
		return fmt.Errorf("query cannot be empty")
	}
	
	// Check for balanced comment blocks
	if err := qv.ValidateComments(query); err != nil {
		return err
	}
	
	// Check for table specification
	if err := qv.ValidateTableSelection(query); err != nil {
		return err
	}
	
	return nil
}

// =============================================================================
// Component Validation
// =============================================================================

// ValidateTableName validates a table name
func (qv *QueryValidator) ValidateTableName(table string) error {
	if table == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	
	// Check if it's a known table or namespace.table format
	if strings.Contains(table, ".") {
		parts := strings.Split(table, ".")
		if len(parts) != 2 {
			return fmt.Errorf("invalid table format: %s (expected namespace.table)", table)
		}
		
		if parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid table format: %s (namespace and table cannot be empty)", table)
		}
	}
	
	return nil
}

// ValidateTableSelection validates that a query has a table selection
func (qv *QueryValidator) ValidateTableSelection(query string) error {
	// Remove comments first
	query = qv.removeComments(query)
	
	// Get the first non-empty line
	lines := strings.Split(query, "\n")
	var firstLine string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "|") {
			firstLine = trimmed
			break
		}
	}
	
	if firstLine == "" {
		return fmt.Errorf("query must start with a table selection")
	}
	
	// Very basic check - should start with a table name
	if strings.HasPrefix(firstLine, "|") {
		return fmt.Errorf("query must start with a table selection, not a clause")
	}
	
	return nil
}

// ValidateTimeSelection validates a time selection clause
func (qv *QueryValidator) ValidateTimeSelection(selection string) error {
	if selection == "" {
		return nil // Time selection is optional
	}
	
	// Check for valid time selection keywords
	validKeywords := []string{"during past", "from", "to", "on", "ago"}
	hasValidKeyword := false
	
	for _, keyword := range validKeywords {
		if strings.Contains(strings.ToLower(selection), keyword) {
			hasValidKeyword = true
			break
		}
	}
	
	if !hasValidKeyword {
		return fmt.Errorf("invalid time selection: %s (must contain 'during past', 'from...to', 'on', or 'ago')", selection)
	}
	
	return nil
}

// ValidateWhereClause validates a where clause
func (qv *QueryValidator) ValidateWhereClause(clause string) error {
	if clause == "" {
		return fmt.Errorf("where clause cannot be empty")
	}
	
	// Check for at least one operator
	hasOperator := false
	operators := []string{"==", "=", "!=", ">", "<", ">=", "<=", "in", "!in", "contains", "!contains"}
	
	for _, op := range operators {
		if strings.Contains(clause, op) {
			hasOperator = true
			break
		}
	}
	
	if !hasOperator {
		return fmt.Errorf("where clause must contain a comparison operator")
	}
	
	return nil
}

// ValidateOperatorUsage validates operator usage with field types
func (qv *QueryValidator) ValidateOperatorUsage(field string, operator Operator, value string) error {
	// Basic validation - could be expanded based on field type
	if field == "" {
		return fmt.Errorf("field name cannot be empty")
	}
	
	if !IsComparisonOperator(operator) {
		return fmt.Errorf("invalid operator: %s", operator)
	}
	
	// Validate operator compatibility
	switch operator {
	case OpIn, OpNotIn:
		// Should have brackets
		if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
			return fmt.Errorf("in/!in operator requires array syntax: [value1, value2]")
		}
	case OpContains, OpNotContains:
		// Typically used with array fields
		// Could add more specific validation here
	}
	
	return nil
}

// ValidateComments validates comment syntax
func (qv *QueryValidator) ValidateComments(query string) error {
	// Count opening and closing comment markers
	openCount := strings.Count(query, "/*")
	closeCount := strings.Count(query, "*/")
	
	if openCount != closeCount {
		return fmt.Errorf("unbalanced comment blocks (found %d /* and %d */)", openCount, closeCount)
	}
	
	// Check for invalid comment placements (simplified check)
	// NQL doesn't allow comments between | and keyword
	invalidPattern := regexp.MustCompile(`\|\s*/\*.*?\*/\s*\w+`)
	if invalidPattern.MatchString(query) {
		return fmt.Errorf("invalid comment placement: comments cannot appear between | and statement keyword")
	}
	
	return nil
}

// =============================================================================
// Function Validation
// =============================================================================

// ValidateAggregateFunction validates an aggregate function usage
func (qv *QueryValidator) ValidateAggregateFunction(fn AggregateFunc, field string) error {
	if !IsAggregateFunction(fn) {
		return fmt.Errorf("invalid aggregate function: %s", fn)
	}
	
	// Some functions require a field
	requiresField := map[AggregateFunc]bool{
		FuncSum:     true,
		FuncAvg:     true,
		FuncMin:     true,
		FuncMax:     true,
		FuncLast:    true,
		FuncSumIf:   true,
		FuncCount:   false, // count() can be used without a field
		FuncCountIf: false,
		FuncP95:     true,
		FuncP05:     true,
	}
	
	if requiresField[fn] && field == "" {
		return fmt.Errorf("function %s() requires a field", fn)
	}
	
	return nil
}

// ValidateDateTimeFunction validates a datetime function usage
func (qv *QueryValidator) ValidateDateTimeFunction(fn DateTimeFunc, field string) error {
	if !IsDateTimeFunction(fn) {
		return fmt.Errorf("invalid datetime function: %s", fn)
	}
	
	if field == "" {
		return fmt.Errorf("datetime function %s() requires a field", fn)
	}
	
	return nil
}

// =============================================================================
// Structure Validation
// =============================================================================

// ValidateComputeRequirement validates that compute clauses have required context
func (qv *QueryValidator) ValidateComputeRequirement(query string) error {
	// If query has compute but no with/include, it's invalid
	hasCompute := strings.Contains(query, "| compute")
	hasWith := strings.Contains(query, "| with")
	hasInclude := strings.Contains(query, "| include")
	
	if hasCompute && !hasWith && !hasInclude {
		return fmt.Errorf("compute clause requires a with or include clause")
	}
	
	return nil
}

// ValidateListSummarizeConflict validates that list and summarize aren't both used
func (qv *QueryValidator) ValidateListSummarizeConflict(query string) error {
	hasList := strings.Contains(query, "| list")
	hasSummarize := strings.Contains(query, "| summarize")
	
	if hasList && hasSummarize {
		return fmt.Errorf("cannot use both list and summarize in the same query")
	}
	
	return nil
}

// =============================================================================
// Field Validation
// =============================================================================

// ValidateFieldName validates a field name format
func (qv *QueryValidator) ValidateFieldName(field string) error {
	if field == "" {
		return fmt.Errorf("field name cannot be empty")
	}
	
	// Basic format check: should contain alphanumeric, dots, underscores
	validField := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\.]*$`)
	if !validField.MatchString(field) {
		return fmt.Errorf("invalid field name format: %s", field)
	}
	
	return nil
}

// =============================================================================
// Helper Methods
// =============================================================================

// removeComments removes comment blocks from a query
func (qv *QueryValidator) removeComments(query string) string {
	// Simple regex to remove /* ... */ comments
	commentPattern := regexp.MustCompile(`/\*.*?\*/`)
	return commentPattern.ReplaceAllString(query, "")
}

// ExtractClauses extracts all clauses from a query
func (qv *QueryValidator) ExtractClauses(query string) map[string][]string {
	clauses := make(map[string][]string)
	
	// Remove comments first
	query = qv.removeComments(query)
	
	// Split by pipe and analyze each part
	parts := strings.Split(query, "|")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Extract clause type (first word)
		words := strings.Fields(part)
		if len(words) == 0 {
			continue
		}
		
		clauseType := strings.ToLower(words[0])
		clauses[clauseType] = append(clauses[clauseType], part)
	}
	
	return clauses
}

// =============================================================================
// Validation Rules
// =============================================================================

// ValidationRule represents a validation rule
type ValidationRule struct {
	Name     string
	Validate func(query string) error
}

// GetValidationRules returns all validation rules
func (qv *QueryValidator) GetValidationRules() []ValidationRule {
	return []ValidationRule{
		{
			Name:     "TableSelection",
			Validate: qv.ValidateTableSelection,
		},
		{
			Name:     "Comments",
			Validate: qv.ValidateComments,
		},
		{
			Name:     "ComputeRequirement",
			Validate: qv.ValidateComputeRequirement,
		},
		{
			Name:     "ListSummarizeConflict",
			Validate: qv.ValidateListSummarizeConflict,
		},
	}
}

// ValidateWithRules validates a query against all rules
func (qv *QueryValidator) ValidateWithRules(query string) []error {
	var errors []error
	
	rules := qv.GetValidationRules()
	for _, rule := range rules {
		if err := rule.Validate(query); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", rule.Name, err))
		}
	}
	
	return errors
}

// =============================================================================
// Public Validation Functions
// =============================================================================

// ValidateNQLQuery validates an NQL query (convenience function)
func ValidateNQLQuery(query string) error {
	validator := NewQueryValidator()
	return validator.ValidateQuery(query)
}

// ValidateNQLQueryDetailed validates an NQL query and returns all validation errors
func ValidateNQLQueryDetailed(query string) []error {
	validator := NewQueryValidator()
	return validator.ValidateWithRules(query)
}
