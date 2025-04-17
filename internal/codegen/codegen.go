package codegen

import (
	"io"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/env"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

type CodeGenerator struct {
	env   *env.Environment[llvm.Reg]
	write *llvm.LLVMWriter
	ng    *NameGenerator
}

func NewCodeGenerator(w io.Writer) *CodeGenerator {
	env := env.NewEnvironment[llvm.Reg]()
	writer := llvm.NewLLVMWriter(w)
	nameGen := NewNameGenerator()
	return &CodeGenerator{env: env, write: writer, ng: nameGen}
}

func (cg *CodeGenerator) GenerateCode(prgm *tast.Prgm) error {
	// boilerplate std functions
	if err := cg.write.Declare(
		llvm.Void, "printInt", llvm.I32); err != nil {
		return err
	}
	if err := cg.write.Declare(
		llvm.Void, "printDouble", llvm.Double,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(
		llvm.Void, "printString", llvm.I8Ptr,
	); err != nil {
		return err
	}
	if err := cg.write.Declare(llvm.I32, "readInt"); err != nil {
		return err
	}
	if err := cg.write.Declare(llvm.Double, "readDouble"); err != nil {
		return err
	}
	cg.env.EnterContext()
	defer cg.env.ExitContext()

	for _, def := range prgm.Defs {
		cg.ng.resetReg()
		cg.ng.resetLab()
		cg.ng.resetPtrs()

		if err := cg.write.Newline(); err != nil {
			return err
		}
		if err := cg.compileDef(def); err != nil {
			return err
		}
		if err := cg.handleStrings(); err != nil {
			return err
		}
	}
	return nil
}

func (cg *CodeGenerator) emitVarAlloc(
	name string,
	typ llvm.Type,
	init ...llvm.Value,
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
		typ := llvm.Array(llvm.I8, len(str)+1)
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
