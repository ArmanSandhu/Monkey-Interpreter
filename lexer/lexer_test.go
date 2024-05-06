package lexer

import (
	"testing"

	"github.com/armansandhu/monkey_interpreter/token"
)

// Test the NextToken function.
func TestNextToken(t *testing.T) {
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