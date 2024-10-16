package evaluator

import "github.com/armansandhu/monkey_interpreter/object"

var builtins = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Incorrect number of arguments detected! Only needed 1 but instead received %d!", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("Argument to `len` is not supported! Instead received an %s!", args[0].Type())
			}
		},
	},
}
