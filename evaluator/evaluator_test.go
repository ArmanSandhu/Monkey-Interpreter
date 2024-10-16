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
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
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

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
			if (10 > 1) {
			return 10;
			}
			return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{
			"5 + true;",
			"Type Mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type Mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Unknown Operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
			if (10 > 1) {
			return true + false;
			}
			return 1;
			}
			`,
			"Unknown Operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier Not Found: foobar",
		},
		{
			`"Hello" - "World"`,
			"Unknown Operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)

		errorObject, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Object is not of type Error! Instead received '%T' (%+v)", evaluated, evaluated)
			continue
		}

		if errorObject.Message != tt.expectedMsg {
			t.Errorf("Object has the incorrect error message! Expected '%s' but receieved '%s'", tt.expectedMsg, errorObject.Message)
		}
	}
}

func TestLetStatments(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEvaluate(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	evaluated := testEvaluate(input)
	function, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf("Object is not of type Function! Instead received '%T' (%+v)", evaluated, evaluated)
	}

	if len(function.Parameters) != 1 {
		t.Errorf("Function has the incorrect amount of parameters! Needed 1 but receieved '%d'", len(function.Parameters))
	}

	if function.Parameters[0].String() != "x" {
		t.Errorf("Function has an incorrect parameter! Needed 'x' but receieved '%q'", function.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if function.Body.String() != expectedBody {
		t.Errorf("Function has an incorrect body! Needed %q but receieved '%q'", expectedBody, function.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEvaluate(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`

	evaluated := testEvaluate(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not of type String! Instead received '%T' (%+v)", evaluated, evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("String has the incorrect value! It should be 'hello world'. Instead received '%q'", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello" + " " + "world"`

	evaluated := testEvaluate(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not of type String! Instead received '%T' (%+v)", evaluated, evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("String has the incorrect value! It should be 'hello world'. Instead received '%q'", str.Value)
	}
}

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "Argument to `len` is not supported! Instead received an INTEGER!"},
		{`len("one", "two")`, "Incorrect number of arguments detected! Only needed 1 but instead received 2!"},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errorObject, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Object is not of type Error! Instead received '%T' (%+v)", evaluated, evaluated)
				continue
			}

			if errorObject.Message != tt.expected {
				t.Errorf("Object has the incorrect error message! Expected '%s' but receieved '%s'", tt.expected, errorObject.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEvaluate(input)
	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Object is not of type Array! Instead received '%T' (%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("Incorrect amount of Array elements detected! Expected 3 but receieved '%d'", len(array.Elements))
	}

	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 4)
	testIntegerObject(t, array.Elements[2], 6)
}

func testEvaluate(input string) object.Object {
	lxr := lexer.New(input)
	prsr := parser.New(lxr)
	program := prsr.ParseProgram()
	env := object.NewEnvironment()

	return Evaluate(program, env)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object is not of type Null! Instead received '%T' (%+v)", obj, obj)
		return false
	}
	return true
}
