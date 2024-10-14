package evaluator

import (
	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/object"
)

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evaluateStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Evaluate(statement)
	}

	return result
}
