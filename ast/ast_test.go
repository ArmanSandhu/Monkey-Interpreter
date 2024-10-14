package ast

import (
	"testing"

	"github.com/armansandhu/monkey_interpreter/token"
)

func TestString(t *testing.T) {
	// Construct a valid AST by hand
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIERS, Literal: "firstVar"},
					Value: "firstVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIERS, Literal: "secondVar"},
					Value: "secondVar",
				},
			},
		},
	}

	if program.String() != "let firstVar = secondVar;" {
		t.Errorf("Incorrect Program String detected! Program String found was '%q'", program.String())
	}
}
