package evaluator

import (
	"fmt"

	"github.com/armansandhu/monkey_interpreter/ast"
	"github.com/armansandhu/monkey_interpreter/object"
)

// constant values for Booleans and Nulls
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Evaluate(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Evaluate(node.Right, env)
		if isError(right) {
			return right
		}
		return evaluatePrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Evaluate(node.Left, env)
		if isError(left) {
			return left
		}
		right := Evaluate(node.Right, env)
		if isError(right) {
			return right
		}
		return evaluateInfixExpression(left, node.Operator, right)
	case *ast.BlockStatement:
		return evaluateBlockStatement(node, env)
	case *ast.IfExpression:
		return evaluateIfExpression(node, env)
	case *ast.ReturnStatement:
		value := Evaluate(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.LetStatement:
		value := Evaluate(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
	case *ast.Identifier:
		return evaluateIdentifier(node, env)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		return &object.Function{Parameters: parameters, Body: body}
	case *ast.CallExpression:
		function := Evaluate(node.Function, env)
		if isError(function) {
			return function
		}
		arguments := evaluateExpressions(node.Arguments, env)
		if len(arguments) == 1 && isError(arguments[0]) {
			return arguments[0]
		}
		return applyFunction(function, arguments)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}
	return nil
}

func evaluateProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Evaluate(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(right)
	case "-":
		return evaluateMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown Operator: %s%s", operator, right.Type())
	}
}

func evaluateBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evaluateMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown Operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluateInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evaluateIntegerInfixExpression(left, operator, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("Type Mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evaluateStringInfixExpression(left, operator, right)
	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateStringInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	if operator != "+" {
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return &object.String{Value: leftValue + rightValue}
}

func evaluateIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Evaluate(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Evaluate(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Evaluate(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evaluateBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Evaluate(statement, env)

		if result != nil {
			returnType := result.Type()
			if returnType == object.RETURN_VALUE_OBJ || returnType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evaluateIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(i.Value); ok {
		return value
	}

	if builtin, ok := builtins[i.Value]; ok {
		return builtin
	}

	return newError("Identifier Not Found: " + i.Value)
}

func evaluateExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range expressions {
		evaluated := Evaluate(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(function object.Object, arguments []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, arguments)
		evaluated := Evaluate(fn.Body, extendedEnv)
		return unWrapReturnValue(evaluated)
	case *object.BuiltIn:
		return fn.Function(arguments...)
	default:
		return newError("Object is not a Function! Received a '%s'", function.Type())
	}
}

func extendFunctionEnv(function *object.Function, arguments []object.Object) *object.Environment {
	environment := object.NewEnclosedEnvironment(function.Env)

	for index, parameter := range function.Parameters {
		environment.Set(parameter.Value, arguments[index])
	}

	return environment
}

func unWrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
