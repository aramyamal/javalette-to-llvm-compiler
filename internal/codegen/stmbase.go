package codegen

import (
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileExpStm(s *tast.ExpStm) error {
	if _, err := cg.compileExp(s.Exp); err != nil {
		return err
	}
	return nil
}

func (cg *CodeGenerator) compileDeclsStm(s *tast.DeclsStm) error {
	for _, item := range s.Items {
		typ := item.Type()
		llvmType := toLlvmType(typ)

		// for arrays, the type is a pointer
		if _, isArray := item.Type().(*tast.ArrayType); isArray {
			llvmType = llvmType.Ptr()
		}

		switch i := item.(type) {
		case *tast.NoInitItem:
			var initValue llvmgen.Value
			var err error

			// handle array initialiation separately
			if _, isArray := typ.(*tast.ArrayType); isArray {
				initValue, err = cg.emitUninitStruct(typ)
				if err != nil {
					return err
				}
			} else {
				initValue = llvmType.ZeroValue()
			}

			if _, err := cg.emitVarAlloc(i.Id, llvmType, initValue); err != nil {
				return err
			}
		case *tast.InitItem:
			value, err := cg.compileExp(i.Exp)
			if err != nil {
				return err
			}

			if _, err := cg.emitVarAlloc(i.Id, llvmType, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cg *CodeGenerator) compileReturnStm(s *tast.ReturnStm) error {
	reg, err := cg.compileExp(s.Exp)
	if err != nil {
		return err
	}

	cg.write.Ret(toLlvmRetType(s.Type), reg)
	return nil
}

func (cg *CodeGenerator) compileWhileStm(s *tast.WhileStm) error {
	conditionLab := cg.ng.nextLab()
	bodyLab := cg.ng.nextLab()
	endLab := cg.ng.nextLab()
	cg.write.Br(conditionLab)
	cg.write.Label(conditionLab)
	des, err := cg.compileExp(s.Exp)
	if err != nil {
		return err
	}
	llvmType := toLlvmType(s.Exp.Type())
	cg.write.BrIf(llvmType, des, bodyLab, endLab)

	cg.write.Label(bodyLab)
	if err := cg.compileStm(s.Stm); err != nil {
		return err
	}
	cg.write.Br(conditionLab)

	cg.write.Label(endLab)
	return nil
}

func (cg *CodeGenerator) compileBlockStm(s *tast.BlockStm) error {
	cg.env.EnterContext()
	defer cg.env.ExitContext()
	for _, stm := range s.Stms {
		if err := cg.compileStm(stm); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CodeGenerator) compileIfStm(s *tast.IfStm) error {
	thenLabel := cg.ng.nextLab()
	elseLabel := cg.ng.nextLab()
	endLabel := cg.ng.nextLab()

	cond, err := cg.compileExp(s.Exp)
	if err != nil {
		return err
	}

	llvmType := toLlvmType(s.Exp.Type())
	cg.write.BrIf(llvmType, cond, thenLabel, elseLabel)

	// then
	cg.write.Label(thenLabel)
	if err := cg.compileStm(s.ThenStm); err != nil {
		return err
	}
	thenReturns := tast.GuaranteesReturn(s.ThenStm)
	if !thenReturns {
		cg.write.Br(endLabel)
	}

	// else
	cg.write.Label(elseLabel)
	elseReturns := false
	if s.ElseStm != nil {
		if err := cg.compileStm(s.ElseStm); err != nil {
			return err
		}
		elseReturns = tast.GuaranteesReturn(s.ElseStm)
	}
	if !elseReturns {
		cg.write.Br(endLabel)
	}

	// only emit the end label if at least one branch does not return
	if !thenReturns || !elseReturns {
		cg.write.Label(endLabel)
	}

	return nil
}
