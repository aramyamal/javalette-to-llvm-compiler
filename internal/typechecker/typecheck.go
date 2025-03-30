package typechecker

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parsing"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typedast"
)

func Typecheck(tree parsing.IPrgmContext) (*typedast.Prgm, error) {
	prgm, ok := tree.(*parsing.PrgmContext)
	if !ok {
		return nil, fmt.Errorf("expected *parsing.ProgramContext, got %T", tree)
	}

	env := NewEnvironment[typedast.Type]()

	var typedDefs []typedast.Def
	for _, def := range prgm.AllDef() {
		typedDef, err := checkDef(env, def)
		if err != nil {
			return nil, err
		}
		typedDefs = append(typedDefs, typedDef)
	}

	typedPrgm := typedast.NewPrgm(typedDefs)
	return typedPrgm, nil
}

func checkDef(
	env *Environment[typedast.Type],
	def parsing.IDefContext,
) (typedast.Def, error) {
	line, col, text := extractData(def)
	switch d := def.(type) {
	case *parsing.FuncDefContext:
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
		return typedast.NewFuncDef(
			d.Ident().GetText(),
			[]typedast.Arg{}, // fix this
			typedStms,
			0, // fix this, supposed to be Type enum calculated from d._Type()
			line, col, text,
		), nil
	default:
		return nil, fmt.Errorf(
			"checkDef: unhandled def type %T at %d:%d near '%s'",
			d, line, col, text,
		)
	}
}

func checkStm(
	env *Environment[typedast.Type],
	stm parsing.IStmContext,
) (typedast.Stm, error) {
	line, col, text := extractData(stm)
	switch s := stm.(type) {
	case *parsing.ExpStmContext:
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
	exp parsing.IExpContext,
) (typedast.Exp, error) {
	line, col, text := extractData(exp)
	switch e := exp.(type) {
	case *parsing.ParenExpContext:
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

func extractData(pr antlr.ParserRuleContext) (int, int, string) {
	return pr.GetStart().GetLine(), pr.GetStart().GetColumn(), pr.GetText()
}
