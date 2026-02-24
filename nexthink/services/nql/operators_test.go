package nql

import (
	"testing"
)

func TestComparisonOperators(t *testing.T) {
	// Test that operator constants have the correct values
	if string(OpEquals) != "==" {
		t.Errorf("OpEquals should be ==, got: %s", string(OpEquals))
	}

	if string(OpNotEquals) != "!=" {
		t.Errorf("OpNotEquals should be !=, got: %s", string(OpNotEquals))
	}

	if string(OpGreater) != ">" {
		t.Errorf("OpGreater should be >, got: %s", string(OpGreater))
	}

	if string(OpLess) != "<" {
		t.Errorf("OpLess should be <, got: %s", string(OpLess))
	}

	if string(OpGreaterEqual) != ">=" {
		t.Errorf("OpGreaterEqual should be >=, got: %s", string(OpGreaterEqual))
	}

	if string(OpLessEqual) != "<=" {
		t.Errorf("OpLessEqual should be <=, got: %s", string(OpLessEqual))
	}
}

func TestListOperators(t *testing.T) {
	if string(OpIn) != "in" {
		t.Errorf("OpIn should be in, got: %s", string(OpIn))
	}

	if string(OpNotIn) != "!in" {
		t.Errorf("OpNotIn should be !in, got: %s", string(OpNotIn))
	}

	if string(OpContains) != "contains" {
		t.Errorf("OpContains should be contains, got: %s", string(OpContains))
	}
}

func TestLogicalOperators(t *testing.T) {
	if string(LogicalAnd) != "and" {
		t.Errorf("LogicalAnd should be and, got: %s", string(LogicalAnd))
	}

	if string(LogicalOr) != "or" {
		t.Errorf("LogicalOr should be or, got: %s", string(LogicalOr))
	}
}

func TestArithmeticOperators(t *testing.T) {
	if string(ArithmeticAdd) != "+" {
		t.Errorf("ArithmeticAdd should be +, got: %s", string(ArithmeticAdd))
	}

	if string(ArithmeticSubtract) != "-" {
		t.Errorf("ArithmeticSubtract should be -, got: %s", string(ArithmeticSubtract))
	}

	if string(ArithmeticMultiply) != "*" {
		t.Errorf("ArithmeticMultiply should be *, got: %s", string(ArithmeticMultiply))
	}

	if string(ArithmeticDivide) != "/" {
		t.Errorf("ArithmeticDivide should be /, got: %s", string(ArithmeticDivide))
	}
}

func TestAggregateFunctions(t *testing.T) {
	// Test basic aggregate functions
	if string(FuncCount) != "count" {
		t.Errorf("FuncCount should be count, got: %s", string(FuncCount))
	}

	if string(FuncSum) != "sum" {
		t.Errorf("FuncSum should be sum, got: %s", string(FuncSum))
	}

	if string(FuncAvg) != "avg" {
		t.Errorf("FuncAvg should be avg, got: %s", string(FuncAvg))
	}

	if string(FuncMin) != "min" {
		t.Errorf("FuncMin should be min, got: %s", string(FuncMin))
	}

	if string(FuncMax) != "max" {
		t.Errorf("FuncMax should be max, got: %s", string(FuncMax))
	}
}


func TestOperatorString(t *testing.T) {
	op := OpEquals
	if op.String() != "==" {
		t.Errorf("Operator.String() should return ==, got: %s", op.String())
	}
}

func TestLogicalOperatorString(t *testing.T) {
	op := LogicalAnd
	if op.String() != "and" {
		t.Errorf("LogicalOperator.String() should return and, got: %s", op.String())
	}
}

func TestAggregateFuncString(t *testing.T) {
	fn := FuncCount
	if fn.String() != "count" {
		t.Errorf("AggregateFunc.String() should return count, got: %s", fn.String())
	}
}
