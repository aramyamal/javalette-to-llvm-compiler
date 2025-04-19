package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

func (cg *CodeGenerator) compileStm(stm tast.Stm) error {
	switch s := stm.(type) {
	case *tast.ExpStm:
		return cg.compileExpStm(s)
	case *tast.DeclsStm:
		return cg.compileDeclsStm(s)
	case *tast.ReturnStm:
		return cg.compileReturnStm(s)
	case *tast.VoidReturnStm:
		return cg.write.Ret(llvm.Void)
	case *tast.WhileStm:
		return cg.compileWhileStm(s)
	default:
		return fmt.Errorf(
			"compileStm: unhandled stm type %T at %d:%d near '%s'",
			s, s.Line(), s.Col(), s.Text(),
		)
	}
}

func (cg *CodeGenerator) compileExpStm(s *tast.ExpStm) error {
	if _, err := cg.compileExp(s.Exp); err != nil {
		return err
	}
	return nil
}

func (cg *CodeGenerator) compileDeclsStm(s *tast.DeclsStm) error {
	for _, item := range s.Items {
		llvmType := toLlvmType(item.Type())
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

	if err := cg.write.Label(endLab); err != nil {
		return err
	}

	return nil
}
