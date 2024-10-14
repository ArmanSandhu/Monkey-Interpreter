package evaluator

import (
	"testing"

	"github.com/armansandhu/monkey_interpreter/lexer"
	"github.com/armansandhu/monkey_interpreter/object"
	"github.com/armansandhu/monkey_interpreter/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!7", false},
		{"!!true", true},
		{"!!false", false},
		{"!!7", true},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEvaluate(input string) object.Object {
	lxr := lexer.New(input)
	prsr := parser.New(lxr)
	program := prsr.ParseProgram()

	return Evaluate(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not of type Integer! Instead received '%T' (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Object has the incorrect value! Expected '%d' but receieved '%d'", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not of type Boolean! Instead received '%T' (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Object has the incorrect value! Expected '%t' but receieved '%t'", expected, result.Value)
		return false
	}

	return true
}
