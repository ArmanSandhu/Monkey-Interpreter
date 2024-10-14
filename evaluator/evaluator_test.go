package evaluator

import (
	"testing"
	"github.com/armansandhu/monkey_interpreter/object"
	"github.com/armansandhu/monkey_interpreter/parser"
	"github.com/armansandhu/monkey_interpreter/lexer"
)

func TestEvalIntegerExpression(t *testing.T)) {
	tests := []struct {
		input		string
		expected 	int64
	}{
		{"5", 5}
		{"10", 10}
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	lxr := lexer.New(tt.input)
	prsr := New(lxr)
	program := prsr.ParseProgram()

	return Eval(program)
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