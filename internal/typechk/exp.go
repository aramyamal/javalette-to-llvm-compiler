package typechk

import (
	"fmt"
	"strconv"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

func (tc *TypeChecker) inferExp(exp parser.IExpContext) (tast.Exp, error) {
	line, col, text := extractPosData(exp)
	switch e := exp.(type) {
	case *parser.ParenExpContext:
		return tc.inferParenExp(e, line, col, text)
	case *parser.BoolExpContext:
		return tc.inferBoolExp(e, line, col, text)
	case *parser.IntExpContext:
		return tc.inferIntExp(e, line, col, text)
	case *parser.DoubleExpContext:
		return tc.inferDoubleExp(e, line, col, text)
	case *parser.IdentExpContext:
		return tc.inferIdentExp(e, line, col, text)
	case *parser.FuncExpContext:
		return tc.inferFuncExp(e, line, col, text)
	case *parser.StringExpContext:
		return tc.inferStringExp(e, line, col, text)
	case *parser.NegExpContext:
		return tc.inferNegExp(e, line, col, text)
	case *parser.NotExpContext:
		return tc.inferNotExp(e, line, col, text)
	case *parser.PostExpContext:
		return tc.inferPostExp(e, line, col, text)
	case *parser.PreExpContext:
		return tc.inferPreExp(e, line, col, text)
	case *parser.MulExpContext:
		return tc.inferMulExp(e, line, col, text)
	case *parser.AddExpContext:
		return tc.inferAddExp(e, line, col, text)
	case *parser.CmpExpContext:
		return tc.inferCmpExp(e, line, col, text)
	case *parser.AndExpContext:
		return tc.inferAndExp(e, line, col, text)
	case *parser.OrExpContext:
		return tc.inferOrExp(e, line, col, text)
	case *parser.AssignExpContext:
		return tc.inferAssignExp(e, line, col, text)
	default:
		return nil, fmt.Errorf(
			"inferExp: unhandled exp type %T at %d:%d near '%s'",
			e, line, col, text,
		)
	}
}

func (tc *TypeChecker) inferParenExp(
	e *parser.ParenExpContext, line, col int, text string,
) (*tast.ParenExp, error) {
	innerExp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}
	return tast.NewParenExp(innerExp, innerExp.Type(), line, col, text), nil
}

func (tc *TypeChecker) inferBoolExp(
	e *parser.BoolExpContext, line, col int, text string,
) (*tast.BoolExp, error) {
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
}

func (tc *TypeChecker) inferIntExp(
	e *parser.IntExpContext, line, col int, text string,
) (*tast.IntExp, error) {
	value, err := strconv.Atoi(e.Integer().GetText())
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse integer '%s' at %d:%d: %w", text, line, col, err,
		)
	}
	return tast.NewIntExp(value, line, col, text), nil
}

func (tc *TypeChecker) inferDoubleExp(
	e *parser.DoubleExpContext, line, col int, text string,
) (*tast.DoubleExp, error) {
	value, err := strconv.ParseFloat(e.Double().GetText(), 64)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse double '%s' at %d:%d: %w", text, line, col, err,
		)
	}
	return tast.NewDoubleExp(value, line, col, text), nil
}

func (tc *TypeChecker) inferIdentExp(
	e *parser.IdentExpContext, line, col int, text string,
) (*tast.IdentExp, error) {
	varName := e.Ident().GetText()
	typ, ok := tc.env.LookupVar(varName)
	if !ok {
		return nil, fmt.Errorf(
			"trying to reference an undeclared variable '%s' at %d:%d",
			varName, line, col,
		)
	}
	return tast.NewIdentExp(varName, typ, line, col, text), nil
}

func (tc *TypeChecker) inferFuncExp(
	e *parser.FuncExpContext, line, col int, text string,
) (*tast.FuncExp, error) {
	// check if func is defined before it is called and that call is correct
	funcName := e.Ident().GetText()
	// check if func signature in env
	sign, exists := tc.env.LookupFunc(funcName)
	if !exists {
		return nil, fmt.Errorf(
			"calling undefined function '%s' at %d:%d", funcName, line, col,
		)
	}
	expTypes := []types.Type{}
	typedExps := []tast.Exp{}

	// extract types in correct order
	paramTypes := make([]types.Type, 0, len(sign.ParamNames))
	for _, paramName := range sign.ParamNames {
		paramTypes = append(paramTypes, sign.Params[paramName])
	}

	for _, exp := range e.AllExp() {
		typedExp, err := tc.inferExp(exp)
		if err != nil {
			return nil, err
		}
		expTypes = append(expTypes, typedExp.Type())
		typedExps = append(typedExps, typedExp)
	}

	// check if number of arguments matches function signature
	if len(paramTypes) != len(expTypes) && len(sign.Params) > 0 {
		return nil, fmt.Errorf(
			"function '%s' called with wrong number of arguments at %d:%d",
			funcName, line, col,
		)
	}

	// verify and promote argument types
	for i := range paramTypes {
		expected := paramTypes[i]
		actual := expTypes[i]

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
		sign.Returns,
		line, col, text,
	), nil
}

func (tc *TypeChecker) inferStringExp(
	e *parser.StringExpContext, line, col int, text string,
) (*tast.StringExp, error) {
	stringWithQuotes := e.String_().GetText()
	// remove quote symbols
	stringWithoutQuotes := stringWithQuotes[1 : len(stringWithQuotes)-1]
	return tast.NewStringExp(stringWithoutQuotes, line, col, text), nil
}

func (tc *TypeChecker) inferNegExp(
	e *parser.NegExpContext, line, col int, text string,
) (*tast.NegExp, error) {
	typedExp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}
	typ := typedExp.Type()
	if !(typ == types.Double || typ == types.Int) {
		return nil, fmt.Errorf(
			"negation not defined for type %s at %d:%d near '%s'",
			typ.String(), line, col, text,
		)
	}
	return tast.NewNegExp(typedExp, typ, line, col, text), nil
}

func (tc *TypeChecker) inferNotExp(
	e *parser.NotExpContext, line, col int, text string,
) (*tast.NotExp, error) {
	typedExp, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}
	if typ := typedExp.Type(); typ != types.Bool {
		return nil, fmt.Errorf(
			"'!' not defined for type bool at %d:%d near '%s'", line, col, text,
		)
	}
	return tast.NewNotExp(typedExp, line, col, text), nil
}

func (tc *TypeChecker) inferPostExp(
	e *parser.PostExpContext, line, col int, text string,
) (*tast.PostExp, error) {
	varName := e.Ident().GetText()
	typ, ok := tc.env.LookupVar(varName)
	if !ok {
		return nil, fmt.Errorf(
			"variable '%s' not found at %d:%d", varName, line, col,
		)
	}
	if typ != types.Int { //&& typ != tast.Double {
		return nil, fmt.Errorf(
			// "'++' or '--' operation can only be done on int or double at "+
			"'++' or '--' operation can only be done on int at "+
				"%d:%d near '%s'", line, col, text,
		)
	}
	var op types.Op
	switch e.IncDecOp().(type) {
	case *parser.IncContext:
		op = types.OpInc
	case *parser.DecContext:
		op = types.OpDec
	default:
		return nil, fmt.Errorf(
			"unhandled postfix operator type %T at %d:%d",
			e.IncDecOp(), line, col,
		)
	}
	return tast.NewPostExp(varName, op, typ, line, col, text), nil
}

func (tc *TypeChecker) inferPreExp(
	e *parser.PreExpContext, line, col int, text string,
) (*tast.PreExp, error) {
	varName := e.Ident().GetText()
	typ, ok := tc.env.LookupVar(varName)
	if !ok {
		return nil, fmt.Errorf(
			"variable '%s' not found at %d:%d", varName, line, col,
		)
	}
	if typ != types.Int && typ != types.Double {
		return nil, fmt.Errorf(
			"'++' or '--' operation can only be done on int or double at "+
				"%d:%d near '%s'", line, col, text,
		)
	}
	var op types.Op
	switch e.IncDecOp().(type) {
	case *parser.IncContext:
		op = types.OpInc
	case *parser.DecContext:
		op = types.OpDec
	default:
		return nil, fmt.Errorf(
			"unhandled prefix operator type %T at %d:%d",
			e.IncDecOp(), line, col,
		)
	}
	return tast.NewPreExp(varName, op, typ, line, col, text), nil
}

func (tc *TypeChecker) inferMulExp(
	e *parser.MulExpContext, line, col int, text string,
) (*tast.MulExp, error) {
	leftExp, err := tc.inferExp(e.Exp(0))
	if err != nil {
		return nil, err
	}
	leftType := leftExp.Type()
	rightExp, err := tc.inferExp(e.Exp(1))
	if err != nil {
		return nil, err
	}
	rightType := rightExp.Type()

	var op types.Op
	switch e.MulOp().(type) {
	case *parser.MulContext:
		op = types.OpMul
	case *parser.DivContext:
		op = types.OpDiv
	case *parser.ModContext:
		op = types.OpMod
		if rightType == types.Double || leftType == types.Double {
			return nil, fmt.Errorf(
				"%s-operation not allowed for bool at %d:%d near '%s'",
				op.String(), line, col, text,
			)
		}
	default:
		return nil, fmt.Errorf(
			"unhandled operator type %T at %d:%d", e.MulOp(), line, col,
		)
	}

	if leftType == types.Bool || rightType == types.Bool {
		return nil, fmt.Errorf(
			"%s-operation not allowed for bool at %d:%d near '%s'",
			op.String(), line, col, text,
		)
	}

	if leftType == types.Void || rightType == types.Void {
		return nil, fmt.Errorf(
			"%s-operation not allowed for void at %d:%d near '%s'",
			op.String(), line, col, text,
		)
	}

	typ, err := dominantType(leftType, rightType)
	if err != nil {
		return nil, err
	}

	return tast.NewMulExp(
		promoteExp(leftExp, typ),
		promoteExp(rightExp, typ),
		op, typ, line, col, text), nil
}

func (tc *TypeChecker) inferAddExp(
	e *parser.AddExpContext, line, col int, text string,
) (*tast.AddExp, error) {
	leftExp, err := tc.inferExp(e.Exp(0))
	if err != nil {
		return nil, err
	}
	leftType := leftExp.Type()
	rightExp, err := tc.inferExp(e.Exp(1))
	if err != nil {
		return nil, err
	}
	rightType := rightExp.Type()

	var op types.Op
	switch e.AddOp().(type) {
	case *parser.AddContext:
		op = types.OpAdd
	case *parser.SubContext:
		op = types.OpSub
	default:
		return nil, fmt.Errorf(
			"unhandled operator type %T at %d:%d", e.AddOp(), line, col,
		)
	}

	if leftType == types.Bool || rightType == types.Bool {
		return nil, fmt.Errorf(
			"%s-operation not allowed for bool at %d:%d near '%s'",
			op.String(), line, col, text,
		)
	}

	if leftType == types.Void || rightType == types.Void {
		return nil, fmt.Errorf(
			"%s-operation not allowed for void at %d:%d near '%s'",
			op.String(), line, col, text,
		)
	}

	typ, err := dominantType(leftType, rightType)
	if err != nil {
		return nil, err
	}

	return tast.NewAddExp(
		promoteExp(leftExp, typ),
		promoteExp(rightExp, typ),
		op, typ, line, col, text), nil
}

func (tc *TypeChecker) inferCmpExp(
	e *parser.CmpExpContext, line, col int, text string,
) (*tast.CmpExp, error) {
	leftExp, err := tc.inferExp(e.Exp(0))
	if err != nil {
		return nil, err
	}
	rightExp, err := tc.inferExp(e.Exp(1))
	if err != nil {
		return nil, err
	}

	leftType := leftExp.Type()
	rightType := rightExp.Type()

	if leftType == types.Void || rightType == types.Void {
		return nil, fmt.Errorf(
			"comparison with void type not allowed at %d:%d near '%s'",
			line, col, text,
		)
	}

	var op types.Op
	switch cmp := e.CmpOp().(type) {
	case *parser.LThContext:
		if leftType == types.Bool || rightType == types.Bool {
			return nil, fmt.Errorf(
				"number comparisons with bool not allowed at %d:%d near '%s'",
				line, col, text,
			)
		}
		op = types.OpLt
	case *parser.GThContext:
		if leftType == types.Bool || rightType == types.Bool {
			return nil, fmt.Errorf(
				"number comparisons with bool not allowed at %d:%d near '%s'",
				line, col, text,
			)
		}
		op = types.OpGt
	case *parser.LTEContext:
		if leftType == types.Bool || rightType == types.Bool {
			return nil, fmt.Errorf(
				"number comparisons with bool not allowed at %d:%d near '%s'",
				line, col, text,
			)
		}
		op = types.OpLe
	case *parser.GTEContext:
		if leftType == types.Bool || rightType == types.Bool {
			return nil, fmt.Errorf(
				"number comparisons with bool not allowed at %d:%d near '%s'",
				line, col, text,
			)
		}
		op = types.OpGe
	case *parser.EquContext:
		if (leftType == types.Bool) != (rightType == types.Bool) {
			return nil, fmt.Errorf(
				"equality comparison between bool and non-bool types not"+
					" allowed at %d:%d near '%s'", line, col, text,
			)
		}
		op = types.OpEq
	case *parser.NEqContext:
		if (leftType == types.Bool) != (rightType == types.Bool) {
			return nil, fmt.Errorf(
				"inequality comparison between bool and non-bool types not"+
					" allowed at %d:%d near '%s'",
				line, col, text,
			)
		}
		op = types.OpNe
	default:
		return nil, fmt.Errorf(
			"unhandled comparison operator type %T at %d:%d",
			cmp, line, col,
		)
	}

	// Get dominant type for proper promotion
	domType, err := dominantType(leftType, rightType)
	if err != nil {
		return nil, err
	}

	return tast.NewCmpExp(
		promoteExp(leftExp, domType),
		promoteExp(rightExp, domType),
		op,
		line, col, text,
	), nil
}

func (tc *TypeChecker) inferAndExp(
	e *parser.AndExpContext, line, col int, text string,
) (*tast.AndExp, error) {
	leftExp, err := tc.inferExp(e.Exp(0))
	if err != nil {
		return nil, err
	}
	rightExp, err := tc.inferExp(e.Exp(1))
	if err != nil {
		return nil, err
	}
	if leftExp.Type() != types.Bool || rightExp.Type() != types.Bool {
		return nil, fmt.Errorf(
			"AND (&&) operation can only occur between booleans at %d:%d "+
				"near '%s'", line, col, text,
		)
	}
	return tast.NewAndExp(leftExp, rightExp, line, col, text), nil
}

func (tc *TypeChecker) inferOrExp(
	e *parser.OrExpContext, line, col int, text string,
) (*tast.OrExp, error) {
	leftExp, err := tc.inferExp(e.Exp(0))
	if err != nil {
		return nil, err
	}
	rightExp, err := tc.inferExp(e.Exp(1))
	if err != nil {
		return nil, err
	}

	if leftExp.Type() != types.Bool || rightExp.Type() != types.Bool {
		return nil, fmt.Errorf(
			"OR (||) operation can only occur between booleans at %d:%d "+
				"near '%s'", line, col, text,
		)
	}
	return tast.NewOrExp(leftExp, rightExp, line, col, text), nil
}

func (tc *TypeChecker) inferAssignExp(
	e *parser.AssignExpContext, line, col int, text string,
) (*tast.AssignExp, error) {
	varName := e.Ident().GetText()
	varType, ok := tc.env.LookupVar(varName)
	if !ok {
		return nil, fmt.Errorf(
			"trying to assign to undeclared variable '%s' at %d:%d",
			varName, line, col,
		)
	}

	expValue, err := tc.inferExp(e.Exp())
	if err != nil {
		return nil, err
	}

	if !isConvertible(varType, expValue.Type()) {
		return nil, fmt.Errorf(
			"illegal implicit conversion in assignment. Expected %s, "+
				"but got %s at %d:%d near '%s'",
			varType, expValue.Type(), line, col, text,
		)
	}

	return tast.NewAssignExp(
		varName,
		promoteExp(expValue, varType),
		varType,
		line, col, text,
	), nil
}

func promoteExp(exp tast.Exp, typ types.Type) tast.Exp {
	if exp.Type() == types.Int && typ == types.Double {
		return tast.NewIntToDoubleExp(exp)
	}
	return exp
}
