// Package codegen provides tools for generating llvm code by traversing a typed
// abstract syntax tree (TAST) created using the tast package. The main entry
// point is the CodeGenerator type, LLVM code using the llvmgen package.
package codegen

import (
	"io"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

// CodeGenerator generates LLVM code from the typed abstract syntax tree (TAST)
// of a Javalette program, as produced by typechk.TypeChecker.
type CodeGenerator struct {
	env   *env.Environment[llvmgen.Reg]
	write *llvmgen.Writer
	ng    *NameGenerator
}

// NewCodeGenerator creates and returns a new CodeGenerator instance that writes
// to w.
func NewCodeGenerator(w io.Writer) *CodeGenerator {
	env := env.NewEnvironment[llvmgen.Reg]()
	writer := llvmgen.NewWriter(w)
	nameGen := NewNameGenerator()
	return &CodeGenerator{env: env, write: writer, ng: nameGen}
}

// GenerateCode performs LLVM code generation for the given TAST prgm
// representing a Javalette program. The input shoud be a pointer to the root of
// the TAST (*tast.Prgm). If an error is encountered during traversal,
// GenerateCode returns it.
func (cg *CodeGenerator) GenerateCode(prgm *tast.Prgm) error {
	// boilerplate std functions
	if err := cg.write.Declare(
		llvmgen.Void, "printInt", llvmgen.I32); err != nil {
		return err
	}
	if err := cg.write.Declare(
		llvmgen.Void, "printDouble", llvmgen.Double,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(
		llvmgen.Void, "printString", llvmgen.I8Ptr,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(llvmgen.I32, "readInt"); err != nil {
		return err
	}
	if err := cg.write.Declare(llvmgen.Double, "readDouble"); err != nil {
		return err
	}
	cg.env.EnterContext()
	defer cg.env.ExitContext()

	for _, def := range prgm.Defs {
		cg.env.EnterContext()

		cg.ng.resetReg()
		cg.ng.resetLab()
		cg.ng.resetPtrs()

		if err := cg.write.Newline(); err != nil {
			return err
		}
		if err := cg.compileDef(def); err != nil {
			return err
		}

		cg.env.ExitContext()
	}

	if err := cg.write.WriteAll(); err != nil {
		return err
	}

	return nil
}

func (cg *CodeGenerator) emitVarAlloc(
	name string,
	typ llvmgen.Type,
	init ...llvmgen.Value,
) error {
	varPtr := cg.ng.ptrName(name)
	cg.env.ExtendVar(name, varPtr)
	if err := cg.write.Alloca(varPtr, typ); err != nil {
		return err
	}
	if len(init) > 0 && init[0] != nil {
		if err := cg.write.Store(typ, init[0], varPtr); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CodeGenerator) handleStrings() error {
	for name, str := range cg.ng.strMap {
		typ := llvmgen.Array(llvmgen.I8, len(str)+1)
		if err := cg.write.Newline(); err != nil {
			return err
		}
		if err := cg.write.InternalConstant(name, typ, str); err != nil {
			return err
		}
	}
	cg.ng.resetStrings()
	return nil
}
