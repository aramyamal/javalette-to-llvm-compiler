package typechecker

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typedast"
)

func Typecheck(tree parser.IPrgmContext) (*typedast.Prgm, error) {
	prgm, ok := tree.(*parser.PrgmContext)
	if !ok {
		return nil, fmt.Errorf("expected *parser.ProgramContext, got %T", tree)
	}
	defs := prgm.AllDef()
	err := validateMainFunc(defs)
	if err != nil {
		return nil, err
	}

	env := NewEnvironment[typedast.Type]()

	if err := validateFuncSigns(env, defs); err != nil {
		return nil, err
	}

	typedDefs, err := checkDefs(env, prgm.AllDef())
	if err != nil {
		return nil, err
	}

	typedPrgm := typedast.NewPrgm(typedDefs)
	return typedPrgm, nil
}

func validateMainFunc(defs []parser.IDefContext) error {
	var mainFunc *parser.FuncDefContext
	for _, def := range defs {
		if funcDef, ok := def.(*parser.FuncDefContext); ok {
			if funcDef.Ident().GetText() == "main" {
				mainFunc = funcDef
				break
			}
		}
	}

	if mainFunc == nil {
		return fmt.Errorf("program has no entrypoint 'main'")
	}

	if typ, err := toAstType(mainFunc.Type_()); err != nil {
		return err
	} else if typ != typedast.Int {
		return fmt.Errorf("'main' entrypoint function does not have type int")
	}

	// TODO add check that main does not have params

	return nil
}

func validateFuncSigns(
	env *Environment[typedast.Type],
	defs []parser.IDefContext,
) error {
	for _, def := range defs {
		if funcDef, ok := def.(*parser.FuncDefContext); ok {
			name := funcDef.Ident().GetText()
			returnType, err := toAstType(funcDef.Type_())
			if err != nil {
				return err
			}

			params, err := extractParams(funcDef.AllArg())
			if err != nil {
				return err
			}
			if err := env.ExtendFunc(name, params, returnType); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractParams(args []parser.IArgContext) (map[string]typedast.Type, error) {
	_, params, err := extractArgs(args) // Ignore typedArgs slice
	return params, err
}

func checkDefs(
	env *Environment[typedast.Type],
	defs []parser.IDefContext,
) ([]typedast.Def, error) {

	var typedDefs []typedast.Def
	for _, def := range defs {
		typedDef, err := checkDef(env, def)
		if err != nil {
			return nil, err
		}
		typedDefs = append(typedDefs, typedDef)
	}
	return typedDefs, nil
}

func checkDef(
	env *Environment[typedast.Type],
	def parser.IDefContext,
) (typedast.Def, error) {
	line, col, text := extractPosData(def)
	switch d := def.(type) {
	case *parser.FuncDefContext:
		// TODO:
		// handle Ident by adding to func. context,
		// handle args by adding to environment,

		var typedStms []typedast.Stm
		for _, stm := range d.AllStm() {
			typedStm, err := checkStm(env, stm)
			if err != nil {
				return nil, err
			}
			typedStms = append(typedStms, typedStm)
		}
		typ, err := toAstType(d.Type_())
		if err != nil {
			return nil, err
		}
		typedArgs, err := toAstArgs(d.AllArg())
		if err != nil {
			return nil, err
		}
		return typedast.NewFuncDef(
			d.Ident().GetText(),
			typedArgs,
			typedStms,
			typ,
			line, col, text,
		), nil
	default:
		return nil, fmt.Errorf(
			"checkDef: unhandled def type %T at %d:%d near '%s'",
			d, line, col, text,
		)
	}
}

func toAstArgs(args []parser.IArgContext) ([]typedast.Arg, error) {
	typedArgs, _, err := extractArgs(args)
	return typedArgs, err
}

func checkStm(
	env *Environment[typedast.Type],
	stm parser.IStmContext,
) (typedast.Stm, error) {
	line, col, text := extractPosData(stm)
	switch s := stm.(type) {
	case *parser.ExpStmContext:
		inferredExp, err := inferExp(env, s.Exp())
		if err != nil {
			return nil, err
		}
		return typedast.NewExpStm(inferredExp, line, col, text), nil

	default:
		return nil, fmt.Errorf(
			"checkStm: unhandled stm type %T at %d:%d near '%s'",
			s, line, col, text,
		)
	}
}

func inferExp(
	env *Environment[typedast.Type],
	exp parser.IExpContext,
) (typedast.Exp, error) {
	line, col, text := extractPosData(exp)
	switch e := exp.(type) {
	case *parser.ParenExpContext:
		innerExp, err := inferExp(env, e.Exp())
		if err != nil {
			return nil, err
		}
		return typedast.NewParenExp(innerExp, innerExp.Type(), line, col, text),
			nil
	default:
		return nil, fmt.Errorf(
			"inferExp: unhandled exp type %T at %d:%d near '%s'",
			e, line, col, text,
		)
	}
}

func extractArgs(
	args []parser.IArgContext,
) ([]typedast.Arg, map[string]typedast.Type, error) {

	typedArgs := []typedast.Arg{}
	params := make(map[string]typedast.Type)

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

			if paramType == typedast.Void {
				return nil, nil, fmt.Errorf(
					"function definition parameter %s of type void at %d:%d",
					paramName, line, col,
				)
			}

			params[paramName] = paramType
			typedArgs = append(typedArgs, typedast.NewParamArg(
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
