package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func extractParams(args []parser.IArgContext) (map[string]tast.Type, error) {
	_, params, err := extractArgs(args) // Ignore typedArgs slice
	return params, err
}

func toAstArgs(args []parser.IArgContext) ([]tast.Arg, error) {
	typedArgs, _, err := extractArgs(args)
	return typedArgs, err
}

func extractArgs(
	args []parser.IArgContext,
) ([]tast.Arg, map[string]tast.Type, error) {

	typedArgs := []tast.Arg{}
	params := make(map[string]tast.Type)

	for _, arg := range args {
		line, col, text := extractPosData(arg)

		switch a := arg.(type) {
		case *parser.ParamArgContext:
			paramName := a.Ident().GetText()
			paramType, err := toAstType(a.Type_())
			if err != nil {
				return nil, nil, err
			}

			if _, exists := params[paramName]; exists {
				return nil, nil, fmt.Errorf(
					"duplicate function parameter name %s in function at %d:%d",
					paramName, line, col,
				)
			}

			if paramType == tast.Void {
				return nil, nil, fmt.Errorf(
					"function definition parameter %s of type void at %d:%d",
					paramName, line, col,
				)
			}

			params[paramName] = paramType
			typedArgs = append(typedArgs, tast.NewParamArg(
				paramType, paramName, line, col, text,
			))

		default:
			return nil, nil, fmt.Errorf(
				"unexpected argument type %T encountered at %d:%d",
				arg, line, col,
			)
		}
	}
	return typedArgs, params, nil
}
