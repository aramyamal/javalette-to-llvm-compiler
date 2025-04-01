package typechk

import (
	"fmt"
	"strconv"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func inferExp(
	env *Environment[tast.Type],
	exp parser.IExpContext,
) (tast.Exp, error) {
	line, col, text := extractPosData(exp)
	switch e := exp.(type) {
	case *parser.ParenExpContext:
		innerExp, err := inferExp(env, e.Exp())
		if err != nil {
			return nil, err
		}
		return tast.NewParenExp(innerExp, innerExp.Type(), line, col, text),
			nil
	case *parser.BoolExpContext:
		switch t := e.BoolLit().(type) {
		case *parser.FalseLitContext:
			return tast.NewBoolExp(false, line, col, text), nil
		case *parser.TrueLitContext:
			return tast.NewBoolExp(true, line, col, text), nil
		default:
			return nil, fmt.Errorf(
				"checkExp: unhandled bool literal type %T at %d:%d near '%s'",
				t, line, col, text,
			)
		}
	case *parser.IntExpContext:
		value, err := strconv.Atoi(e.Integer().GetText())
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse integer '%s' at %d:%d: %w",
				text, line, col, err,
			)
		}
		return tast.NewIntExp(value, line, col, text), nil

	case *parser.DoubleExpContext:
		value, err := strconv.ParseFloat(e.Double().GetText(), 64)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse double '%s' at %d:%d: %w",
				text, line, col, err,
			)
		}
		return tast.NewDoubleExp(value, line, col, text), nil

	case *parser.IdentExpContext:
		varName := e.Ident().GetText()
		typ, ok := env.LookupVar(varName)
		if !ok {
			return nil, fmt.Errorf(
				"trying to reference an undeclared variable '%s' at %d:%d",
				varName, line, col,
			)
		}
		return tast.NewIdentExp(varName, typ, line, col, text), nil

	case *parser.FuncExpContext:
		// check if func is defined before it is called and that call is correct
		funcName := e.Ident().GetText()
		// check if func signature in env
		sign, exists := env.LookupFunc(funcName)
		if !exists {
			return nil, fmt.Errorf(
				"calling undefined function '%s' at %d:%d",
				funcName, line, col,
			)
		}
		types := []tast.Type{}
		typedExps := []tast.Exp{}

		// extract types in correct order
		paramTypes := make([]tast.Type, 0, len(sign.paramNames))
		for _, paramName := range sign.paramNames {
			paramTypes = append(paramTypes, sign.params[paramName])
		}

		for _, exp := range e.AllExp() {
			typedExp, err := inferExp(env, exp)
			if err != nil {
				return nil, err
			}
			types = append(types, typedExp.Type())
			typedExps = append(typedExps, typedExp)
		}

		// check if number of arguments matches function signature
		if len(paramTypes) != len(types) && len(sign.params) > 0 {
			return nil, fmt.Errorf(
				"function '%s' called with wrong number of arguments at %d:%d",
				funcName, line, col,
			)
		}

		// verify and promote argument types
		for i := range paramTypes {
			expected := paramTypes[i]
			actual := types[i]

			if !isConvertible(expected, actual) {
				return nil, fmt.Errorf(
					"argument %d of function '%s' has incompatible type. "+
						"Expected %s but got %s at %d:%d",
					i+1, funcName, expected, actual, line, col,
				)
			}

			// promote expression if needed
			typedExps[i] = promoteExp(typedExps[i], expected)
		}

		return tast.NewFuncExp(
			funcName,
			typedExps,
			sign.returns,
			line, col, text,
		), nil

	default:
		return nil, fmt.Errorf(
			"inferExp: unhandled exp type %T at %d:%d near '%s'",
			e, line, col, text,
		)
	}
}

func promoteExp(exp tast.Exp, typ tast.Type) tast.Exp {
	if exp.Type() == tast.Int && typ == tast.Double {
		return tast.NewIntToDoubleExp(exp)
	}
	return exp
}
