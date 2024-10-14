package evaluator

import (
	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/object"
)

// constant values for Booleans and Nulls
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
