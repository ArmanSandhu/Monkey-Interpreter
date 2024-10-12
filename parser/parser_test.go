package parser

import (
	"testing"

	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/lexer"
)

func TestLetStatments(t *testing.T) {
	// Create test string
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	lxr := lexer.New(input)
	prsr := New(lxr)

	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)
	if program == nil {
		t.Fatalf("ParseProgram failed! Returned NIL!")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements contains the incorrect amount of statements! Expected 3 but returned %d instead!", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatment(t, statement, tt.expectedIdentifier) {
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
	// Create test string
	input := `
	return 5;
	return 10;
	return 993322;
	`

	lxr := lexer.New(input)
	prsr := New(lxr)

	program := prsr.ParseProgram()
	checkForParseErrors(t, prsr)
	if program == nil {
		t.Fatalf("ParseProgram failed! Returned NIL!")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements contains the incorrect amount of statements! Expected 3 but returned %d instead!", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not of type *ast.ReturnStatement. Instead received %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() is not 'return'. Instead received '%q'", returnStmt.TokenLiteral())
		}
	}
}
