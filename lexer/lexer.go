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

	// Skip any whitespace
	lexer.skipWhiteSpace()

	// Based on the current character return the appropriate token
	switch lexer.ch {
	case '=':
		if lexer.peekChar() == '=' {
			firstChar := lexer.ch
			lexer.readChar()
			newLiteral := string(firstChar) + string(lexer.ch)
			tok = token.Token{Type: token.EQ, Literal: newLiteral}
		} else {
			tok = newToken(token.ASSIGN, lexer.ch)
		}
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
	case '-':
		tok = newToken(token.MINUS, lexer.ch)
	case '*':
		tok = newToken(token.ASTERISK, lexer.ch)
	case '/':
		tok = newToken(token.SLASH, lexer.ch)
	case '!':
		if lexer.peekChar() == '=' {
			firstChar := lexer.ch
			lexer.readChar()
			newLiteral := string(firstChar) + string(lexer.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: newLiteral}
		} else {
			tok = newToken(token.BANG, lexer.ch)
		}
	case '<':
		tok = newToken(token.LT, lexer.ch)
	case '>':
		tok = newToken(token.RT, lexer.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(lexer.ch) {
			tok.Literal = lexer.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}

	// Advance the pointers
	lexer.readChar()
	return tok
}

// This helper function returns a new Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// This helper function checks if a byte is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// This helper function checks if a byte is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// This function reads in an identifier and advances the lexer's position until a non-letter character is encountered
func (lexer *Lexer) readIdentifier() string {
	currPos := lexer.pos
	// If the current character is a letter move forward until a non letter character has been found
	for isLetter(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[currPos:lexer.pos]
}

// This function reads a number and advances the lexer's positions until a non-digit is encountered
func (lexer *Lexer) readNumber() string {
	currPos := lexer.pos
	// If the current character is a digit move forward until a non digit has been found
	for isDigit(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[currPos:lexer.pos]
}

// This function skips any existing whitespace
func (lexer *Lexer) skipWhiteSpace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

// This helper function will check for two character operators such as == and !=
func (lexer *Lexer) peekChar() byte {
	// check if we've reached the end of our input string
	if lexer.readPos >= len(lexer.input) {
		return 0
	} else {
		// otherwise return the next pos of the input
		return lexer.input[lexer.readPos]
	}
}
