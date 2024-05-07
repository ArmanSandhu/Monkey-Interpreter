package lexer

import (
	"testing"

	"github.com/armansandhu/monkey_interpreter/token"
)

// Test the NextToken function.
func TestNextTokenBasic(t *testing.T) {
	// Create test string
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		token := lexer.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - TokenType Wrong! Expected=%q, Got=%q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Token Literal Wrong! Expected=%q, Got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}
}

func TestNextTokenExtended(t *testing.T) {
	// Create test string
	input := `let six = 6;
	let eight = 8;
	
	let add = fn(a, b) {
		a + b
	};
	
	let result = add(six, eight);
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIERS, "six"},
		{token.ASSIGN, "="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "eight"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIERS, "a"},
		{token.COMMA, ","},
		{token.IDENTIFIERS, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIERS, "a"},
		{token.PLUS, "+"},
		{token.IDENTIFIERS, "b"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIERS, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIERS, "six"},
		{token.COMMA, ","},
		{token.IDENTIFIERS, "eight"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		token := lexer.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - TokenType Wrong! Expected=%q, Got=%q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Token Literal Wrong! Expected=%q, Got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}
}

func TestNextTokenExtendedOperators(t *testing.T) {
	// Create test string
	input := `let six = 6;
	let eight = 8;
	
	let add = fn(a, b) {
		a + b
	};
	
	let result = add(six, eight);
	!-/*5;
	5 < 10 > 5;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIERS, "six"},
		{token.ASSIGN, "="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "eight"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIERS, "a"},
		{token.COMMA, ","},
		{token.IDENTIFIERS, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIERS, "a"},
		{token.PLUS, "+"},
		{token.IDENTIFIERS, "b"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIERS, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIERS, "six"},
		{token.COMMA, ","},
		{token.IDENTIFIERS, "eight"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		token := lexer.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - TokenType Wrong! Expected=%q, Got=%q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Token Literal Wrong! Expected=%q, Got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}
}

func TestNextTokenExtendedKeywords(t *testing.T) {
	// Create test string
	input := `let six = 6;
	let eight = 8;
	
	let add = fn(a, b) {
		a + b
	};
	
	let result = add(six, eight);
	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIERS, "six"},
		{token.ASSIGN, "="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "eight"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIERS, "a"},
		{token.COMMA, ","},
		{token.IDENTIFIERS, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIERS, "a"},
		{token.PLUS, "+"},
		{token.IDENTIFIERS, "b"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIERS, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIERS, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIERS, "six"},
		{token.COMMA, ","},
		{token.IDENTIFIERS, "eight"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		token := lexer.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - TokenType Wrong! Expected=%q, Got=%q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Token Literal Wrong! Expected=%q, Got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}
}
