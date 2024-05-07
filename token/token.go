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

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdentifier(identifier string) TokenType {
	// Check if the identifier exists in the keyword map
	if token, ok := keywords[identifier]; ok {
		return token
	}
	return IDENTIFIERS
}
