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

	case *parser.StringExpContext:
		return tast.NewStringExp(e.String_().GetText(), line, col, text), nil

	case *parser.NegExpContext:
		typedExp, err := inferExp(env, e.Exp())
		if err != nil {
			return nil, err
		}
		typ := typedExp.Type()
		if !(typ == tast.Double || typ == tast.Int) {
			return nil, fmt.Errorf(
				"negation not defined for type %s at %d:%d near '%s'",
				typ.String(), line, col, text,
			)
		}
		return tast.NewNegExp(typedExp, typ, line, col, text), nil

	case *parser.NotExpContext:
		typedExp, err := inferExp(env, e.Exp())
		if err != nil {
			return nil, err
		}
		if typ := typedExp.Type(); typ != tast.Bool {
			return nil, fmt.Errorf(
				"'!' not defined for type bool at %d:%d near '%s'",
				line, col, text,
			)
		}
		return tast.NewNotExp(typedExp, line, col, text), nil

	case *parser.PostExpContext:
		varName := e.Ident().GetText()
		typ, ok := env.LookupVar(varName)
		if !ok {
			return nil, fmt.Errorf(
				"variable '%s' not found at %d:%d",
				varName, line, col,
			)
		}
		if typ != tast.Int { //&& typ != tast.Double {
			return nil, fmt.Errorf(
				"'++' or '--' operation can only be done on int or double at "+
					"%d:%d near '%s'",
				line, col, text,
			)
		}
		var op tast.Op
		switch e.IncDecOp().(type) {
		case *parser.IncContext:
			op = tast.OpInc
		case *parser.DecContext:
			op = tast.OpDec
		default:
			return nil, fmt.Errorf(
				"unhandled postfix operator type %T at %d:%d",
				e.IncDecOp(), line, col,
			)
		}
		return tast.NewPostExp(varName, op, typ, line, col, text), nil

	case *parser.PreExpContext:
		varName := e.Ident().GetText()
		typ, ok := env.LookupVar(varName)
		if !ok {
			return nil, fmt.Errorf(
				"variable '%s' not found at %d:%d",
				varName, line, col,
			)
		}
		if typ != tast.Int && typ != tast.Double {
			return nil, fmt.Errorf(
				"'++' or '--' operation can only be done on int or double at "+
					"%d:%d near '%s'",
				line, col, text,
			)
		}
		var op tast.Op
		switch e.IncDecOp().(type) {
		case *parser.IncContext:
			op = tast.OpInc
		case *parser.DecContext:
			op = tast.OpDec
		default:
			return nil, fmt.Errorf(
				"unhandled prefix operator type %T at %d:%d",
				e.IncDecOp(), line, col,
			)
		}
		return tast.NewPostExp(varName, op, typ, line, col, text), nil

	case *parser.MulExpContext:
		leftExp, err := inferExp(env, e.Exp(0))
		if err != nil {
			return nil, err
		}
		leftType := leftExp.Type()
		rightExp, err := inferExp(env, e.Exp(1))
		if err != nil {
			return nil, err
		}
		rightType := rightExp.Type()

		var op tast.Op
		switch e.MulOp().(type) {
		case *parser.MulContext:
			op = tast.OpMul
		case *parser.DivContext:
			op = tast.OpDiv
		case *parser.ModContext:
			op = tast.OpMod
			if rightType == tast.Double || leftType == tast.Double {
				return nil, fmt.Errorf(
					"%s-operation not allowed for bool at %d:%d near '%s'",
					op.String(), line, col, text,
				)
			}
		default:
			return nil, fmt.Errorf(
				"unhandled operator type %T at %d:%d",
				e.MulOp(), line, col,
			)
		}

		if leftType == tast.Bool || rightType == tast.Bool {
			return nil, fmt.Errorf(
				"%s-operation not allowed for bool at %d:%d near '%s'",
				op.String(), line, col, text,
			)
		}

		if leftType == tast.Void || rightType == tast.Void {
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

	case *parser.AddExpContext:

		leftExp, err := inferExp(env, e.Exp(0))
		if err != nil {
			return nil, err
		}
		leftType := leftExp.Type()
		rightExp, err := inferExp(env, e.Exp(1))
		if err != nil {
			return nil, err
		}
		rightType := rightExp.Type()

		var op tast.Op
		switch e.AddOp().(type) {
		case *parser.AddContext:
			op = tast.OpAdd
		case *parser.SubContext:
			op = tast.OpSub
		default:
			return nil, fmt.Errorf(
				"unhandled operator type %T at %d:%d",
				e.AddOp(), line, col,
			)
		}

		if leftType == tast.Bool || rightType == tast.Bool {
			return nil, fmt.Errorf(
				"%s-operation not allowed for bool at %d:%d near '%s'",
				op.String(), line, col, text,
			)
		}

		if leftType == tast.Void || rightType == tast.Void {
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

	case *parser.CmpExpContext:
		leftExp, err := inferExp(env, e.Exp(0))
		if err != nil {
			return nil, err
		}
		rightExp, err := inferExp(env, e.Exp(1))
		if err != nil {
			return nil, err
		}

		leftType := leftExp.Type()
		rightType := rightExp.Type()

		if leftType == tast.Void || rightType == tast.Void {
			return nil, fmt.Errorf(
				"comparison with void type not allowed at %d:%d near '%s'",
				line, col, text,
			)
		}

		var op tast.Op
		switch cmp := e.CmpOp().(type) {
		case *parser.LThContext:
			if leftType == tast.Bool || rightType == tast.Bool {
				return nil, fmt.Errorf(
					"number comparisons with bool not allowed at %d:%d near '%s'",
					line, col, text,
				)
			}
			op = tast.OpLt
		case *parser.GThContext:
			if leftType == tast.Bool || rightType == tast.Bool {
				return nil, fmt.Errorf(
					"number comparisons with bool not allowed at %d:%d near '%s'",
					line, col, text,
				)
			}
			op = tast.OpGt
		case *parser.LTEContext:
			if leftType == tast.Bool || rightType == tast.Bool {
				return nil, fmt.Errorf(
					"number comparisons with bool not allowed at %d:%d near '%s'",
					line, col, text,
				)
			}
			op = tast.OpLt
		case *parser.GTEContext:
			if leftType == tast.Bool || rightType == tast.Bool {
				return nil, fmt.Errorf(
					"number comparisons with bool not allowed at %d:%d near '%s'",
					line, col, text,
				)
			}
			op = tast.OpGt
		case *parser.EquContext:
			if (leftType == tast.Bool) != (rightType == tast.Bool) {
				return nil, fmt.Errorf(
					"equality comparison between bool and non-bool types not"+
						" allowed at %d:%d near '%s'",
					line, col, text,
				)
			}
			op = tast.OpEq
		case *parser.NEqContext:
			if (leftType == tast.Bool) != (rightType == tast.Bool) {
				return nil, fmt.Errorf(
					"inequality comparison between bool and non-bool types not"+
						" allowed at %d:%d near '%s'",
					line, col, text,
				)
			}
			op = tast.OpNe
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

	case *parser.AndExpContext:
		leftExp, err := inferExp(env, e.Exp(0))
		if err != nil {
			return nil, err
		}
		rightExp, err := inferExp(env, e.Exp(1))
		if err != nil {
			return nil, err
		}

		if leftExp.Type() != tast.Bool || rightExp.Type() != tast.Bool {
			return nil, fmt.Errorf(
				"AND (&&) operation can only occur between booleans at %d:%d "+
					"near '%s'",
				line, col, text,
			)
		}

		return tast.NewAndExp(leftExp, rightExp, line, col, text), nil

	case *parser.OrExpContext:
		leftExp, err := inferExp(env, e.Exp(0))
		if err != nil {
			return nil, err
		}
		rightExp, err := inferExp(env, e.Exp(1))
		if err != nil {
			return nil, err
		}

		if leftExp.Type() != tast.Bool || rightExp.Type() != tast.Bool {
			return nil, fmt.Errorf(
				"OR (||) operation can only occur between booleans at %d:%d "+
					"near '%s'",
				line, col, text,
			)
		}

		return tast.NewOrExp(leftExp, rightExp, line, col, text), nil

	case *parser.AssignExpContext:
		varName := e.Ident().GetText()
		varType, ok := env.LookupVar(varName)
		if !ok {
			return nil, fmt.Errorf(
				"trying to assign to undeclared variable '%s' at %d:%d",
				varName, line, col,
			)
		}

		expValue, err := inferExp(env, e.Exp())
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
