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

// Precedence Table - associates token types with their precedence
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

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
	prsr.registerPrefix(token.TRUE, prsr.parseBoolean)
	prsr.registerPrefix(token.FALSE, prsr.parseBoolean)
	prsr.registerPrefix(token.LPAREN, prsr.parseGroupedExpression)
	prsr.registerPrefix(token.IF, prsr.parseIfExpression)
	prsr.registerPrefix(token.FUNCTION, prsr.parseFunctionLiteral)

	// Initialize the infix parse map and register parsing functions for all the infix operators
	prsr.infixParseFns = make(map[token.TokenType]infixParseFn)
	prsr.registerInfix(token.PLUS, prsr.parseInfixExpression)
	prsr.registerInfix(token.MINUS, prsr.parseInfixExpression)
	prsr.registerInfix(token.SLASH, prsr.parseInfixExpression)
	prsr.registerInfix(token.ASTERISK, prsr.parseInfixExpression)
	prsr.registerInfix(token.EQ, prsr.parseInfixExpression)
	prsr.registerInfix(token.NOT_EQ, prsr.parseInfixExpression)
	prsr.registerInfix(token.LT, prsr.parseInfixExpression)
	prsr.registerInfix(token.GT, prsr.parseInfixExpression)
	prsr.registerInfix(token.LPAREN, prsr.parseCallExpresssion)

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
func (p *Parser) parseLetStatement() *ast.LetStatement {
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

	stmt.Value = p.parseExpression(LOWEST)

	// Skip expression parsing for now
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// This function returns a statement based on encountering a RETURN token.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// Create a ReturnStatement struct
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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
	// Check if there is an appropriate parsing function associated with the currToken.Type at the prefix position.
	prefix := p.prefixParseFns[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	// We're going to try and find the appropriate infix parse function for the next token as long as the precedence is higher
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

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

// This helper function returns the precedence associated with the token type of p.peekToken. If no precedence is found we return with the lowest precedence.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// This helper function returns the precedence associated with the token type of p.currToken. If no precedence is found we return with the lowest precedence.
func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// This method parses a boolean
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

// This method parses an if expression
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	return block
}

// This method parses a function literal
func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpresssion(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currToken, Function: function}
	expression.Arguments = p.parseCallArguments()
	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return arguments
	}

	p.nextToken()
	arguments = append(arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		arguments = append(arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return arguments
}
