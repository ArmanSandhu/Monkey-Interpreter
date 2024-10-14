package parser

import (
	"fmt"
	"testing"

	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/lexer"
)

func TestLetStatments(t *testing.T) {
	// Create test stuct
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)

		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contains the incorrect amount of statements! Expected 1 but returned %d instead!", len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatment(t, statement, tt.expectedIdentifier) {
			return
		}

		value := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, tt.expectedValue) {
			return
		}
	}
}

func checkForParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has encountered %d errors!", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser Error: %q", msg)
	}

	t.FailNow()
}

func testLetStatment(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral is not 'let'. Instead received '%q'", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt is not of type *ast.LetStatement. Instead received %T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s'. Instead received '%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() is not '%s'. Instead received '%s'", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	// Create test struct
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)

		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contains the incorrect amount of statements! Expected 1 but returned %d instead!", len(program.Statements))
		}

		statement := program.Statements[0]

		returnStmt, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not of type *ast.ReturnStatement. Instead received %T", statement)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() is not 'return'. Instead received '%q'", returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	// Create test string
	input := "foobar;"

	lxr := lexer.New(input)
	prsr := New(lxr)
	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.Statement[0] is not of type ast.Expression! Instead received '%T'", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression is not of type *ast.Expression! Instead received '%T'", statement.Expression)
	}
	if identifier.Value != "foobar" {
		t.Fatalf("Identifer Value is not %s. Instead received '%s'", "foobar", identifier.Value)
	}
	if identifier.TokenLiteral() != "foobar" {
		t.Fatalf("Identifer TokenLiteral() is not %s. Instead received '%s'", "foobar", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	// Create test string
	input := "7;"

	lxr := lexer.New(input)
	prsr := New(lxr)
	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.Statement[0] is not of type ast.Expression! Instead received '%T'", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expression is not of type *ast.Expression! Instead received '%T'", statement.Expression)
	}
	if literal.Value != 7 {
		t.Fatalf("Literal Value is not %d. Instead received '%d'", 7, literal.Value)
	}
	if literal.TokenLiteral() != "7" {
		t.Fatalf("Literal TokenLiteral() is not %s. Instead received '%s'", "7", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!foobar", "!", "foobar"},
		{"-foobar", "-", "foobar"},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)
		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		if len(program.Statements) != 1 {
			t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Program.Statement[0] is not of type ast.Expression! Instead received '%T'", program.Statements[0])
		}

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expression is not of type *ast.Expression! Instead received '%T'", statement.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("Expression operator is not '%s'. Instead received '%s'", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not of type *ast.IntegerLiteral! Instead received '%T'", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d! Instead received '%d'", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral is not %d! Instead received '%s'", value, integer.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"7 + 7;", 7, "+", 7},
		{"7 - 7;", 7, "-", 7},
		{"7 * 7;", 7, "*", 7},
		{"7 / 7;", 7, "/", 7},
		{"7 > 7;", 7, ">", 7},
		{"7 < 7;", 7, "<", 7},
		{"7 == 7;", 7, "==", 7},
		{"7 != 7;", 7, "!=", 7},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)
		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		if len(program.Statements) != 1 {
			t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Program.Statement[0] is not of type ast.Expression! Instead received '%T'", program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)
		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Incorrect parsing detected!. Expected %q but instead received '%q'", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression is not of type *ast.Identifier! Instead received '%T'", exp)
		return false
	}

	if identifier.Value != value {
		t.Fatalf("Identifer Value is not %s. Instead received '%s'", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Fatalf("Identifer TokenLiteral() is not %s. Instead received '%s'", value, identifier.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("Type of Expression cannot be handled! Received type '%T'", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	operatorExpression, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Expression is not of type *ast.InfixExpression! Instead received '%T(%s)'", exp, exp)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Fatalf("Expression operator is not '%s'. Instead received '%s'", operator, operatorExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input        string
		expectedBool bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)
		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		if len(program.Statements) != 1 {
			t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Program.Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", program.Statements[0])
		}

		booleanFlag, ok := statement.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("Expression is not of type *ast.Boolean! Instead received '%T'", statement.Expression)
		}

		if booleanFlag.Value != tt.expectedBool {
			t.Fatalf("Boolean Value is not %t. Instead received '%t'", tt.expectedBool, booleanFlag.Value)
		}
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	booleanFlag, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expression is not of type *ast.Boolean! Instead received '%T'", exp)
		return false
	}

	if booleanFlag.Value != value {
		t.Fatalf("Boolean Value is not %t. Instead received '%t'", value, booleanFlag.Value)
		return false
	}

	if booleanFlag.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf("Boolean TokenLiteral() is not %t. Instead received '%s'", value, booleanFlag.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	lxr := lexer.New(input)
	prsr := New(lxr)
	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", program.Statements[0])
	}

	exp, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression is not of type *ast.IfExpression! Instead received '%T'", statement.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Incorrect amount of Consequences detected! Needed 1 but received '%d'", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("Expression Alternative Statements was not nil. Instead received '%v'", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	lxr := lexer.New(input)
	prsr := New(lxr)
	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", program.Statements[0])
	}

	exp, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression is not of type *ast.IfExpression! Instead received '%T'", statement.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Incorrect amount of Consequences detected! Needed 1 but received '%d'", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("Incorrect amount of Alternatives detected! Needed 1 but received '%d'", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	lxr := lexer.New(input)
	prsr := New(lxr)
	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements! Expected 1 but got '%d'", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expression is not of type *ast.FunctionLiterals! Instead received '%T'", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("Incorrect amount of function literal parameters found! Expected 2 but receieved '%d'", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Function Body does not have enough statements! Expected 1 but got '%d'", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Function.Body.Statement[0] is not of type ast.ExpressionStatement! Instead received '%T'", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		lxr := lexer.New(tt.input)
		prsr := New(lxr)
		program := prsr.ParseProgram()
		checkForParseErrors(t, prsr)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf("Incorrect amount of function literal parameters found! Expected %d but receieved '%d'", len(tt.expectedParams), len(function.Parameters))
		}

		for i, identifier := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}
