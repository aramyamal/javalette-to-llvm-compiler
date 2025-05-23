package typechk

import (
	"fmt"
	"slices"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) checkDefs(defs []parser.IDefContext) ([]tast.Def, error) {

	var typedDefs []tast.Def
	for _, def := range defs {
		typedDef, err := tc.checkDef(def)
		if err != nil {
			return nil, err
		}
		typedDefs = append(typedDefs, typedDef)
		tc.env.SetReturnType(tast.Unknown)
	}
	return typedDefs, nil
}

func (tc *TypeChecker) checkDef(def parser.IDefContext) (tast.Def, error) {
	tc.env.EnterContext()
	line, col, text := extractPosData(def)
	switch d := def.(type) {
	case *parser.FuncDefContext:
		return tc.checkFuncDef(d, line, col, text)
	case *parser.StructDefContext:
		return tc.checkStructDef(d, line, col, text)
	case *parser.TypedefDefContext:
		return tc.checkTypedefDef(d, line, col, text)
	default:
		return nil, fmt.Errorf(
			"checkDef: unhandled def type %T at %d:%d near '%s'",
			d, line, col, text,
		)
	}
}

func (tc *TypeChecker) checkFuncDef(
	d *parser.FuncDefContext, line, col int, text string,
) (*tast.FuncDef, error) {
	_, params, err := tc.extractParams(d.AllArg())
	if err != nil {
		return nil, err
	}

	for varName, typ := range params {
		ok := tc.env.ExtendVar(varName, typ)
		if !ok {
			return nil, fmt.Errorf(
				"duplicate parameter name '%s' in function '%s' at %d:%d",
				varName, d.Ident().GetText(), line, col,
			)
		}
	}

	typ, err := tc.toTastType(d.Type_())
	if err != nil {
		return nil, err
	}
	tc.env.SetReturnType(typ)

	var typedStms []tast.Stm
	for _, stm := range d.AllStm() {
		typedStm, err := tc.checkStm(stm)
		if err != nil {
			return nil, err
		}
		typedStms = append(typedStms, typedStm)
	}

	hasReturn := slices.ContainsFunc(typedStms, tast.GuaranteesReturn)

	if typ != tast.Void && !hasReturn {
		return nil, fmt.Errorf(
			"function '%s' at %d:%d does not have a return", text, line, col,
		)
	}

	if typ == tast.Void && !hasReturn {
		var voidReturn *tast.VoidReturnStm
		if len(typedStms) > 0 {
			lastStm := typedStms[len(typedStms)-1]
			voidReturn = tast.NewVoidReturnStm(
				lastStm.Line(),
				lastStm.Col(),
				lastStm.Text(),
			)
		} else {
			voidReturn = tast.NewVoidReturnStm(line, col, text)
		}
		typedStms = append(typedStms, voidReturn)
	}

	typedArgs, err := tc.toTastArgs(d.AllArg())
	if err != nil {
		return nil, err
	}
	tc.env.ExitContext()
	return tast.NewFuncDef(
		d.Ident().GetText(),
		typedArgs,
		typedStms,
		typ,
		line, col, text,
	), nil
}

func (tc *TypeChecker) checkStructDef(
	d *parser.StructDefContext, line, col int, text string,
) (*tast.StructDef, error) {
	structName := d.Ident().GetText()
	envStruct, ok := tc.env.LookupStruct(structName)
	if !ok {
		return nil, fmt.Errorf(
			"error typechecking struct %s at %d:%d near %s",
			structName, line, col, text,
		)
	}
	structType, ok := envStruct.(*tast.StructType)
	if !ok {
		return nil, fmt.Errorf(
			"error typechecking struct, did not retrieve expected type %T, but"+
				" got %T instead", structType, envStruct,
		)
	}
	return tast.NewStructDef(structType, line, col, text), nil
}

func (tc *TypeChecker) checkTypedefDef(
	d *parser.TypedefDefContext, line, col int, text string,
) (*tast.TypedefDef, error) {
	alias := d.Type_(1).GetText()
	aliasedType, err := tc.toTastType(d.Type_(0))
	if err != nil {
		return nil, fmt.Errorf(
			"cannot create typedef '%s' at %d:%d near %s: %w",
			alias, line, col, text, err,
		)
	}
	return tast.NewTypedefDef(alias, aliasedType, line, col, text), nil
}
