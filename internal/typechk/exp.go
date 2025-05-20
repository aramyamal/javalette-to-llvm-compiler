package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
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
	case *parser.NewArrExpContext:
		return tc.inferNewArrExp(e, line, col, text)
	case *parser.IdentExpContext:
		return tc.inferIdentExp(e, line, col, text)
	case *parser.FuncExpContext:
		return tc.inferFuncExp(e, line, col, text)
	case *parser.ArrIndexExpContext:
		return tc.inferArrIndexExp(e, line, col, text)
	case *parser.FieldExpContext:
		return tc.inferFieldExp(e, line, col, text)
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
