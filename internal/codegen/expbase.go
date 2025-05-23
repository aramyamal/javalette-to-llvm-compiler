package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileStringExp(e *tast.StringExp) (
	llvmgen.Value, error,
) {
	des := cg.ng.nextReg()
	glbVar, strLen, alreadyWritten := cg.ng.getOrAddString(e.Value)
	if !alreadyWritten {
		typ := llvmgen.Array(llvmgen.I8, strLen)
		cg.write.InternalConstant(glbVar, typ, llvmgen.LitString(e.Value))
	}
	cg.write.GetElementPtr(des, llvmgen.Array(llvmgen.I8, strLen),
		llvmgen.Array(llvmgen.I8, strLen).Ptr(), glbVar,
		llvmgen.LitInt(0), llvmgen.LitInt(0),
	)
	return des, nil
}

func (cg *CodeGenerator) compileIdentExp(e *tast.IdentExp) (
	llvmgen.Value, error,
) {
	ptrReg, err := cg.compileIdentLExp(e)
	if err != nil {
		return nil, err
	}
	typ := cg.toLlvmType(e.Type())

	if _, isStruct := typ.(*llvmgen.StructType); isStruct {
		// for arrays and structs, load a pointer to them
		des := cg.ng.nextReg()
		ptrType := typ.Ptr()
		cg.write.Load(des, ptrType, ptrType.Ptr(), ptrReg)
		return des, nil
	}

	// load the value in the pointer
	des := cg.ng.nextReg()
	cg.write.Load(des, typ, typ.Ptr(), ptrReg)
	return des, nil
}

func (cg *CodeGenerator) compileIdentLExp(e *tast.IdentExp) (
	llvmgen.Reg, error,
) {
	reg, ok := cg.env.LookupVar(e.Id)
	if !ok {
		return "", fmt.Errorf(
			"internal compiler error: undefined variable '%s' encountered"+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Id, e.Line(), e.Col(), e.Text(),
		)
	}
	return reg, nil
}

func (cg *CodeGenerator) compileFuncExp(e *tast.FuncExp) (
	llvmgen.Value, error,
) {
	args, err := cg.emitFuncArgs(e.Exps)
	if err != nil {
		return nil, err
	}

	des := cg.ng.nextReg()
	cg.write.Call(des, cg.toLlvmRetType(e.Type()), llvmgen.Global(e.Id), args...)
	return des, nil
}

func (cg *CodeGenerator) compileFuncLExp(e *tast.FuncExp) (
	llvmgen.Reg, error,
) {
	typ := cg.toLlvmType(e.Type())
	if _, isFieldProvider := typ.(tast.FieldProvider); !isFieldProvider {
		return "", fmt.Errorf(
			"internal compiler error in compileFuncLExp: return type %s"+
				"is not assignable at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			typ.String(), e.Line(), e.Col(), e.Text(),
		)
	}

	args, err := cg.emitFuncArgs(e.Exps)
	if err != nil {
		return "", err
	}

	des := cg.ng.nextReg()
	cg.write.Call(des, cg.toLlvmRetType(e.Type()), llvmgen.Global(e.Id), args...)
	return des, nil
}

func (cg *CodeGenerator) emitFuncArgs(exps []tast.Exp) (
	[]llvmgen.FuncArg, error,
) {
	var args []llvmgen.FuncArg
	for _, exp := range exps {
		value, err := cg.compileExp(exp)
		if err != nil {
			return nil, err
		}
		args = append(args, llvmgen.Arg(cg.toLlvmRetType(exp.Type()), value))
	}
	return args, nil
}

func (cg *CodeGenerator) compileNegExp(e *tast.NegExp) (llvmgen.Value, error) {
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	llvmType := cg.toLlvmType(e.Type())
	switch llvmType {
	case llvmgen.I32:
		err = cg.write.Sub(des, llvmType, llvmgen.LitInt(0), value)
	case llvmgen.Double:
		err = cg.write.Sub(des, llvmType, llvmgen.LitDouble(0.0), value)
	default:
		return nil, fmt.Errorf(
			"internal compiler error: unable to negate expression "+
				"during code generation at %d:%d near '%s'. "+
				"This should have been caught during type checking.",
			e.Line(), e.Col(), e.Text(),
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileNegExp: %w at %d:%d near %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
}

func (cg *CodeGenerator) compileNotExp(e *tast.NotExp) (llvmgen.Value, error) {
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	cg.write.Xor(des, llvmgen.I1, value, llvmgen.LitBool(true))
	return des, nil
}

func (cg *CodeGenerator) compilePostExp(e *tast.PostExp) (
	llvmgen.Value, error,
) {
	ptrReg, err := cg.compileLExp(e.Exp)
	if err != nil {
		return nil, err
	}
	orig := cg.ng.nextReg()
	typ := cg.toLlvmType(e.Type())
	cg.write.Load(orig, typ, typ.Ptr(), ptrReg)
	incrm := cg.ng.nextReg()

	switch e.Op {
	case tast.OpInc:
		err = cg.write.Add(incrm, typ, orig, llvmgen.LitInt(1))
	case tast.OpDec:
		err = cg.write.Sub(incrm, typ, orig, llvmgen.LitInt(1))
	default:
		return nil, fmt.Errorf(
			"compilePostExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)

	}
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compilePostExp: %w at %d:%d near %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}
	cg.write.Store(typ, incrm, typ.Ptr(), ptrReg)
	return orig, nil
}

func (cg *CodeGenerator) compilePreExp(e *tast.PreExp) (llvmgen.Value, error) {
	ptrReg, err := cg.compileLExp(e.Exp)
	if err != nil {
		return nil, err
	}
	orig := cg.ng.nextReg()
	typ := cg.toLlvmType(e.Type())
	cg.write.Load(orig, typ, typ.Ptr(), ptrReg)
	incrm := cg.ng.nextReg()

	switch e.Op {
	case tast.OpInc:
		err = cg.write.Add(incrm, typ, orig, llvmgen.LitInt(1))
	case tast.OpDec:
		err = cg.write.Sub(incrm, typ, orig, llvmgen.LitInt(1))
	default:
		return nil, fmt.Errorf(
			"compileExp->PostExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)

	}
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compilePreExp: %w at %d:%d near %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}
	cg.write.Store(typ, incrm, typ.Ptr(), ptrReg)
	return incrm, nil
}

func (cg *CodeGenerator) compileMulExp(e *tast.MulExp) (llvmgen.Value, error) {
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	switch e.Op {
	case tast.OpMul:
		err = cg.write.Mul(des, cg.toLlvmType(e.Type()), lhs, rhs)
	case tast.OpDiv:
		err = cg.write.Div(des, cg.toLlvmType(e.Type()), lhs, rhs)
	case tast.OpMod:
		err = cg.write.Rem(des, cg.toLlvmType(e.Type()), lhs, rhs)
	default:
		return nil, fmt.Errorf(
			"compileExp->MulExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileMulExp: %w at %d:%d near %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
}

func (cg *CodeGenerator) compileAddExp(e *tast.AddExp) (llvmgen.Value, error) {
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	switch e.Op {
	case tast.OpAdd:
		err = cg.write.Add(des, cg.toLlvmType(e.Type()), lhs, rhs)
	case tast.OpSub:
		err = cg.write.Sub(des, cg.toLlvmType(e.Type()), lhs, rhs)
	default:
		return nil, fmt.Errorf(
			"compileExp->AddExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileAddExp: %w at %d:%d near %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
}

func (cg *CodeGenerator) compileCmpExp(
	e *tast.CmpExp,
) (llvmgen.Value, error) {
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	des := cg.ng.nextReg()
	typ := cg.toLlvmType(e.LeftExp.Type())

	switch e.Op {
	case tast.OpLt:
		err = cg.write.CmpLt(des, typ, lhs, rhs)
	case tast.OpGt:
		err = cg.write.CmpGt(des, typ, lhs, rhs)
	case tast.OpLe:
		err = cg.write.CmpLe(des, typ, lhs, rhs)
	case tast.OpGe:
		err = cg.write.CmpGe(des, typ, lhs, rhs)
	case tast.OpEq:
		err = cg.write.CmpEq(des, typ, lhs, rhs)
	case tast.OpNe:
		err = cg.write.CmpNe(des, typ, lhs, rhs)
	default:
		return nil, fmt.Errorf(
			"compileExp->CmpExp: unhandled op type '%v' at %d:%d near '%s'",
			e.Op.Name(), e.Line(), e.Col(), e.Text(),
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileCmpExp: %w at %d:%d near %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}
	return des, nil
}

func (cg *CodeGenerator) compileAndExp(e *tast.AndExp) (llvmgen.Value, error) {
	llvmType := cg.toLlvmType(e.Type())
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	falseLab := cg.ng.nextLab()
	evalLab := cg.ng.nextLab()
	endLab := cg.ng.nextLab()

	if err := cg.write.BrIf(llvmType, lhs, evalLab, falseLab); err != nil {
		return nil, err
	}

	cg.write.Label(evalLab)
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	cg.write.Br(endLab)

	cg.write.Label(falseLab)
	cg.write.Br(endLab)

	des := cg.ng.nextReg()
	cg.write.Label(endLab)
	if err := cg.write.Phi(
		des, llvmType,
		llvmgen.Phi(llvmgen.LitBool(false), falseLab),
		llvmgen.Phi(rhs, evalLab),
	); err != nil {
		return nil, err
	}
	return des, nil
}

func (cg *CodeGenerator) compileOrExp(e *tast.OrExp) (llvmgen.Value, error) {
	llvmType := cg.toLlvmType(e.Type())
	lhs, err := cg.compileExp(e.LeftExp)
	if err != nil {
		return nil, err
	}
	trueLab := cg.ng.nextLab()
	evalLab := cg.ng.nextLab()
	endLab := cg.ng.nextLab()

	if err := cg.write.BrIf(llvmType, lhs, trueLab, evalLab); err != nil {
		return nil, err
	}

	cg.write.Label(evalLab)
	rhs, err := cg.compileExp(e.RightExp)
	if err != nil {
		return nil, err
	}
	cg.write.Br(endLab)

	cg.write.Label(trueLab)
	cg.write.Br(endLab)

	des := cg.ng.nextReg()
	cg.write.Label(endLab)
	if err := cg.write.Phi(
		des, llvmType,
		llvmgen.Phi(llvmgen.LitBool(true), trueLab), llvmgen.Phi(rhs, evalLab),
	); err != nil {
		return nil, err
	}
	return des, nil
}

func (cg *CodeGenerator) compileAssignExp(
	e *tast.AssignExp,
) (llvmgen.Value, error) {
	lhsPtr, err := cg.compileLExp(e.ExpLhs)
	if err != nil {
		return nil, err
	}
	value, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	typ := cg.toLlvmRetType(e.Type())
	cg.write.Store(typ, value, typ.Ptr(), lhsPtr)
	return value, nil
}
