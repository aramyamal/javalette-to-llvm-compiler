package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileExp(exp tast.Exp) (llvmgen.Value, error) {
	switch e := exp.(type) {
	case *tast.ParenExp:
		return cg.compileExp(e.Exp)
	case *tast.BoolExp:
		return llvmgen.LitBool(e.Value), nil
	case *tast.IntExp:
		return llvmgen.LitInt(e.Value), nil
	case *tast.DoubleExp:
		return llvmgen.LitDouble(e.Value), nil
	case *tast.NewArrExp:
		return cg.compileNewArrExp(e)
	case *tast.IdentExp:
		return cg.compileIdentExp(e)
	case *tast.FuncExp:
		return cg.compileFuncExp(e)
	case *tast.ArrIndexExp:
		return cg.compileArrIndexExp(e)
	case *tast.FieldExp:
		return cg.compileFieldExp(e)
	case *tast.StringExp:
		return cg.compileStringExp(e)
	case *tast.NegExp:
		return cg.compileNegExp(e)
	case *tast.NotExp:
		return cg.compileNotExp(e)
	case *tast.PostExp:
		return cg.compilePostExp(e)
	case *tast.ArrPostExp:
		return cg.compileArrPostExp(e)
	case *tast.PreExp:
		return cg.compilePreExp(e)
	case *tast.ArrPreExp:
		return cg.compileArrPreExp(e)
	case *tast.MulExp:
		return cg.compileMulExp(e)
	case *tast.AddExp:
		return cg.compileAddExp(e)
	case *tast.CmpExp:
		return cg.compileCmpExp(e)
	case *tast.AndExp:
		return cg.compileAndExp(e)
	case *tast.OrExp:
		return cg.compileOrExp(e)
	case *tast.AssignExp:
		return cg.compileAssignExp(e)
	case *tast.ArrAssignExp:
		return cg.compileArrAssignExp(e)
	default:
		return nil, fmt.Errorf(
			"compileExp: unhandled exp type %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}
