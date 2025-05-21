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
	case *tast.NullPtrExp:
		return llvmgen.Null(), nil
	case *tast.BoolExp:
		return llvmgen.LitBool(e.Value), nil
	case *tast.IntExp:
		return llvmgen.LitInt(e.Value), nil
	case *tast.DoubleExp:
		return llvmgen.LitDouble(e.Value), nil
	case *tast.NewArrExp:
		return cg.compileNewArrExp(e)
	case *tast.NewStructExp:
		return cg.compileNewStructExp(e)
	case *tast.IdentExp:
		return cg.compileIdentExp(e)
	case *tast.FuncExp:
		return cg.compileFuncExp(e)
	case *tast.ArrIndexExp:
		return cg.compileArrIndexExp(e)
	case *tast.FieldExp:
		return cg.compileFieldExp(e)
	case *tast.DerefExp:
		return cg.compileDerefExp(e)
	case *tast.StringExp:
		return cg.compileStringExp(e)
	case *tast.NegExp:
		return cg.compileNegExp(e)
	case *tast.NotExp:
		return cg.compileNotExp(e)
	case *tast.PostExp:
		return cg.compilePostExp(e)
	case *tast.PreExp:
		return cg.compilePreExp(e)
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
	default:
		return nil, fmt.Errorf(
			"compileExp: unhandled exp type %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}

func (cg *CodeGenerator) compileLExp(exp tast.Exp) (llvmgen.Reg, error) {
	switch e := exp.(type) {
	case *tast.ParenExp:
		return cg.compileLExp(e.Exp)
	case *tast.FuncExp:
		return cg.compileFuncLExp(e)
	case *tast.IdentExp:
		return cg.compileIdentLExp(e)
	case *tast.ArrIndexExp:
		return cg.compileArrIndexLExp(e)
	case *tast.FieldExp:
		return cg.compileFieldLExp(e)
	case *tast.DerefExp:
		return cg.compileDerefLExp(e)
	default:
		return "", fmt.Errorf(
			"compileLExp: expression is not assignable %T at %d:%d near '%s'",
			e, e.Line(), e.Col(), e.Text(),
		)
	}
}
