package lexer

import (
	"github.com/armansandhu/monkey_interpreter/token"
)

// Lexer data structure:
// pos - the current position in the input, which points to the current character
// readpos - the current reading position in the input, which points to after the current character
// ch - the current character under examination
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

// This function takes in an input string and returns a Lexer struct
func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

// This helper function reads the next character and advances our position within the input string.
func (lexer *Lexer) readChar() {
	// check if we've reached the end of our input string
	if lexer.readPos >= len(lexer.input) {
		// set the current character to "NUL"
		lexer.ch = 0
	} else {
		// otherwise set the current character to the lexer's readPos
		lexer.ch = lexer.input[lexer.readPos]
	}
	// Update our position and readPosition
	lexer.pos = lexer.readPos
	lexer.readPos += 1
}

// This function returns a token based on the current character that we are looking at within the input
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	// Based on the current character return the appropriate token
	switch lexer.ch {
	case '=':
		tok = newToken(token.ASSIGN, lexer.ch)
	case '+':
		tok = newToken(token.PLUS, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	// Advance the pointers
	lexer.readChar()
	return tok
}

// This function returns a new Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}