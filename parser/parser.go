package parser

import (
	"fmt"
	"strconv"

	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/lexer"
	"github.com/armansandhu/monkey_interpreter/token"
)

// Constant's Declaration
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -a or !a
	CALL        // fn(a)
)

// Parser data structure:
// l - pointer to an instance of the lexer
// currToken - pointer to the current token being processed
// peekToken - a pointer to the next tken that will be processed
// errors - a slice containing all the errors encountered as a part of the parsing process - debug only
// prefixParseFns - a map containing all the prefix parsing functions associated with a TokenType
// infixParseFns - a map containing all the infix parsing functions associated with a TokenType
type Parser struct {
	lxr *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// Define the two types of functions needed for implementing a Pratt Parser
// The prefix parse function will be used for cases where an expression has a prefix operator
// The infix parse function handles all other standard expression types such as 5 + 5
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// This function takes in a lexer struct, creates the parser.
func New(l *lexer.Lexer) *Parser {
	prsr := &Parser{lxr: l, errors: []string{}}

	// Two tokens are read so that the currToken and peekToken can be set.
	prsr.nextToken()
	prsr.nextToken()

	// Initialize the prefix parse map and register parsing functions for the associated tokens
	prsr.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	prsr.registerPrefix(token.IDENTIFIERS, prsr.parseIdentifier)
	prsr.registerPrefix(token.INT, prsr.parseIntegerLiteral)
	prsr.registerPrefix(token.BANG, prsr.parsePrefixExpression)
	prsr.registerPrefix(token.MINUS, prsr.parsePrefixExpression)

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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

// This function returns a statement based on encountering a RETURN token.
func (p *Parser) parseReturnStatement() ast.Statement {
	// Create a ReturnStatement struct
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	// Skip expression parsing for now
	if !p.expectPeek(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// This funciton returns an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// Create an ExpressionStatement struct
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

// This helper function helps add entries to our prefix parsing map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// This helper function helps add entries to our infix parsing map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// This function allows us to parse an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) noPrefixParseFnError(token token.TokenType) {
	msg := fmt.Sprintf("No Prefix Parse function found for %s found!", token)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse %q as an Integer!", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}
