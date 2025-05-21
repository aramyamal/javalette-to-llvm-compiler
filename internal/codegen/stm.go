package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
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
		return cg.write.Ret(llvmgen.Void)
	case *tast.ForEachStm:
		return cg.compileForEachStm(s)
	case *tast.WhileStm:
		return cg.compileWhileStm(s)
	case *tast.BlockStm:
		return cg.compileBlockStm(s)
	case *tast.IfStm:
		return cg.compileIfStm(s)
	case *tast.BlankStm:
		return nil
	default:
		return fmt.Errorf(
			"compileStm: unhandled stm type %T at %d:%d near '%s'",
			s, s.Line(), s.Col(), s.Text(),
		)
	}
}
