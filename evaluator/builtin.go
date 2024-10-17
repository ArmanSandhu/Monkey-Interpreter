package evaluator

import (
	"fmt"

	"github.com/armansandhu/monkey_interpreter/object"
)

var builtins = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Incorrect number of arguments detected! Only needed 1 but instead received %d!", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("Argument to `len` is not supported! Instead received an %s!", args[0].Type())
			}
		},
	},
	"first": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Incorrect number of arguments detected! Only needed 1 but instead received %d!", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `first` must be ARRAY! Instead received an %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Incorrect number of arguments detected! Only needed 1 but instead received %d!", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `last` must be ARRAY! Instead received an %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[len(array.Elements)-1]
			}

			return NULL
		},
	},
	"rest": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Incorrect number of arguments detected! Only needed 1 but instead received %d!", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `rest` must be ARRAY! Instead received an %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				newArray := make([]object.Object, length-1, length-1)
				copy(newArray, array.Elements[1:length])
				return &object.Array{Elements: newArray}
			}

			return NULL
		},
	},
	"push": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Incorrect number of arguments detected! Only needed 2 but instead received %d!", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("Argument to `push` must be ARRAY! Instead received an %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)

			newArray := make([]object.Object, length+1, length+1)
			copy(newArray, array.Elements)
			newArray[length] = args[1]
			return &object.Array{Elements: newArray}
		},
	},
	"puts": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
