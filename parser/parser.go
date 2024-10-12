package parser

import (
	"fmt"

	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/lexer"
	"github.com/armansandhu/monkey_interpreter/token"
)

// Parser data structure:
// l - pointer to an instance of the lexer
// currToken - pointer to the current token being processed
// peekToken - a pointer to the next tken that will be processed
// errors - a slice containing all the errors encountered as a part of the parsing process - debug only
type Parser struct {
	lxr *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string
}

// This function takes in a lexer struct, creates the parser.
func New(l *lexer.Lexer) *Parser {
	prsr := &Parser{lxr: l, errors: []string{}}

	// Two tokens are read so that the currToken and peekToken can be set.
	prsr.nextToken()
	prsr.nextToken()

	return prsr
}

// This helper function advances the currToken and peekToken pointers.
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lxr.NextToken()
}

// This function is responsible for creating our AST.
func (p *Parser) ParseProgram() *ast.Program {
	// Create the 'root' of the AST
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Iterate over every token in the input until the EOF token is encountered
	for !p.currTokenIs(token.EOF) {
		// Parse a statement
		stmt := p.parseStatement()

		// If the statemnent returned is not nil, then we've parsed something
		// Append it to the slice of statements belonging to this node
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// Advance the currToken and peekToken pointers
		p.nextToken()
	}

	return program
}

// This function returns a statement depending on the different tokens.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// This function returns a statement based on encountering a LET token.
func (p *Parser) parseLetStatement() ast.Statement {
	// Create a LetStatement struct
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENTIFIERS) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	// Skip expression parsing for now
	if !p.expectPeek(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) currTokenIs(token token.TokenType) bool {
	return p.currToken.Type == token
}

func (p *Parser) peekTokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

// This method is one of the assertion functions needed for parsing to happen correctly.
// It enforces the correctness of the order of tokens by checking the type of the next token.
func (p *Parser) expectPeek(token token.TokenType) bool {
	if p.peekTokenIs(token) {
		p.nextToken()
		return true
	} else {
		p.peekError(token)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(token token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be '%s', instead received '%s'!", token, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
