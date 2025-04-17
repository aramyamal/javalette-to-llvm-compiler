package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

func (cg *CodeGenerator) compileStm(stm tast.Stm) error {
	switch s := stm.(type) {
	case *tast.ExpStm:
		if _, err := cg.compileExp(s.Exp); err != nil {
			return err
		}
		return nil

	case *tast.DeclsStm:
		for _, item := range s.Items {
			llvmType := toLlvmType(item.Type())
			switch i := item.(type) {
			case *tast.NoInitItem:
				if err := cg.emitVarAlloc(i.Id, llvmType); err != nil {
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

	case *tast.ReturnStm:
		reg, err := cg.compileExp(s.Exp)
		if err != nil {
			return err
		}
		if err := cg.write.Ret(toLlvmType(s.Type), reg); err != nil {
			return err
		}
		return nil

	case *tast.VoidReturnStm:
		return cg.write.Ret(llvm.Void)

	default:
		return fmt.Errorf(
			"compileStm: unhandled stm type %T at %d:%d near '%s'",
			s, s.Line(), s.Col(), s.Text(),
		)
	}
}
