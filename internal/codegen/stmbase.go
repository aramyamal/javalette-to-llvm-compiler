package codegen

import (
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (cg *CodeGenerator) compileExpStm(s *tast.ExpStm) error {
	if _, err := cg.compileExp(s.Exp); err != nil {
		return err
	}
	return nil
}

func (cg *CodeGenerator) compileDeclsStm(s *tast.DeclsStm) error {
	for _, item := range s.Items {
		llvmType := toLlvmType(item.Type())

		// for arrays, the type is a pointer
		if _, isArray := item.Type().(*tast.ArrayType); isArray {
			llvmType = llvmType.Ptr()
		}

		switch i := item.(type) {
		case *tast.NoInitItem:
			if err := cg.emitVarAlloc(
				i.Id, llvmType,
				llvmType.ZeroValue(),
			); err != nil {
				return err
			}
		case *tast.InitItem:
			value, err := cg.compileExp(i.Exp)
			if err != nil {
				return err
			}

			if err := cg.emitVarAlloc(i.Id, llvmType, value); err != nil {
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
	if err := cg.write.Ret(toLlvmType(s.Type), reg); err != nil {
		return err
	}
	return nil

}

func (cg *CodeGenerator) compileWhileStm(s *tast.WhileStm) error {
	conditionLab := cg.ng.nextLab()
	bodyLab := cg.ng.nextLab()
	endLab := cg.ng.nextLab()
	if err := cg.write.Br(conditionLab); err != nil {
		return err
	}
	if err := cg.write.Label(conditionLab); err != nil {
		return err
	}
	des, err := cg.compileExp(s.Exp)
	if err != nil {
		return err
	}
	llvmType := toLlvmType(s.Exp.Type())
	if err := cg.write.BrIf(llvmType, des, bodyLab, endLab); err != nil {
		return err
	}

	if err := cg.write.Label(bodyLab); err != nil {
		return err
	}
	if err := cg.compileStm(s.Stm); err != nil {
		return err
	}
	if err := cg.write.Br(conditionLab); err != nil {
		return err
	}

	return cg.write.Label(endLab)
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
	if err := cg.write.BrIf(llvmType, cond, thenLabel, elseLabel); err != nil {
		return err
	}

	// then
	if err := cg.write.Label(thenLabel); err != nil {
		return err
	}
	if err := cg.compileStm(s.ThenStm); err != nil {
		return err
	}
	thenReturns := tast.GuaranteesReturn(s.ThenStm)
	if !thenReturns {
		if err := cg.write.Br(endLabel); err != nil {
			return err
		}
	}

	// else
	if err := cg.write.Label(elseLabel); err != nil {
		return err
	}
	elseReturns := false
	if s.ElseStm != nil {
		if err := cg.compileStm(s.ElseStm); err != nil {
			return err
		}
		elseReturns = tast.GuaranteesReturn(s.ElseStm)
	}
	if !elseReturns {
		if err := cg.write.Br(endLabel); err != nil {
			return err
		}
	}

	// only emit the end label if at least one branch does not return
	if !thenReturns || !elseReturns {
		if err := cg.write.Label(endLabel); err != nil {
			return err
		}
	}

	return nil
}
