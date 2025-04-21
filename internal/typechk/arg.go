package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func extractParams(
	args []parser.IArgContext,
) ([]string, map[string]tast.Type, error) {
	_, paramNames, params, err := extractArgs(args) // Ignore typedArgs slice
	return paramNames, params, err
}

func toAstArgs(args []parser.IArgContext) ([]tast.Arg, error) {
	typedArgs, _, _, err := extractArgs(args)
	return typedArgs, err
}

func extractArgs(
	args []parser.IArgContext,
) ([]tast.Arg, []string, map[string]tast.Type, error) {

	typedArgs := []tast.Arg{}
	paramNames := []string{}
	params := make(map[string]tast.Type)

	for _, arg := range args {
		line, col, text := extractPosData(arg)

		switch a := arg.(type) {
		case *parser.ParamArgContext:
			paramName := a.Ident().GetText()
			paramType, err := toTastType(a.Type_())
			if err != nil {
				return nil, nil, nil, err
			}

			if _, exists := params[paramName]; exists {
				return nil, nil, nil, fmt.Errorf(
					"duplicate function parameter name %s in function at %d:%d",
					paramName, line, col,
				)
			}

			if paramType == tast.Void {
				return nil, nil, nil, fmt.Errorf(
					"function definition parameter %s of type void at %d:%d",
					paramName, line, col,
				)
			}

			params[paramName] = paramType
			paramNames = append(paramNames, paramName)
			typedArgs = append(typedArgs, tast.NewParamArg(
				paramType, paramName, line, col, text,
			))

		default:
			return nil, nil, nil, fmt.Errorf(
				"unexpected argument type %T encountered at %d:%d",
				arg, line, col,
			)
		}
	}
	return typedArgs, paramNames, params, nil
}
