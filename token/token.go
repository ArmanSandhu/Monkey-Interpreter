package token

// This allows us to use many different values as TokenTypes
type TokenType string

// Define the constants we'll need for the Monkey programming language
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENTIFIERS = "IDENTIFIERS"
	INT         = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	LBRACE = "{"
	RPAREN = ")"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// Token data structure
type Token struct {
	Type    TokenType
	Literal string
}
