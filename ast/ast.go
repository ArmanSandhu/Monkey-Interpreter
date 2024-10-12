package ast

import (
	"github.com/armansandhu/monkey_interpreter/token"
)

// The base interface which all nodes in our Abstract Syntax Tree (AST) must implement.
// Note: The TokenLiteral method is only meant to aid in debugging and testing.
type Node interface {
	TokenLiteral() string
}

// An extenstion of the Node interface. A Statement does not produce a value.
type Statement interface {
	Node
	statementNode()
}

// An extenstion of the Node interface. An Expression does produce a value.
type Expression interface {
	Node
	expressionNode()
}

// This struct represents the root node of our AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // This will be the LET Token
	Name  *Identifier
	Value Expression
}

type Identifier struct {
	Token token.Token // This will be the IDENTIFIERS Token
	Value string
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
