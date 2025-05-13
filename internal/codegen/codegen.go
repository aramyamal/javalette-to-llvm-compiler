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
	env         *env.Environment[llvmgen.Reg]
	write       *llvmgen.Writer
	ng          *NameGenerator
	declTypes   map[string]struct{}
	declGlobals map[string]struct{}
}

// NewCodeGenerator creates and returns a new CodeGenerator instance that writes
// to w.
func NewCodeGenerator(w io.Writer) *CodeGenerator {
	env := env.NewEnvironment[llvmgen.Reg]()
	writer := llvmgen.NewWriter(w)
	nameGen := NewNameGenerator()
	return &CodeGenerator{
		env:         env,
		write:       writer,
		ng:          nameGen,
		declTypes:   make(map[string]struct{}),
		declGlobals: make(map[string]struct{}),
	}
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
	if err := cg.emitFuncDecl(
		llvmgen.Void, "printDouble", llvmgen.Double,
	); err != nil {
		return err
	}
	if err := cg.emitFuncDecl(
		llvmgen.Void, "printString", llvmgen.I8.Ptr(),
	); err != nil {
		return err
	}
	if err := cg.emitFuncDecl(llvmgen.I32, "readInt"); err != nil {
		return err
	}
	if err := cg.emitFuncDecl(llvmgen.Double, "readDouble"); err != nil {
		return err
	}
	cg.env.EnterContext()
	defer cg.env.ExitContext()

	for _, def := range prgm.Defs {
		cg.env.EnterContext()

		cg.ng.resetNames()

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

func (cg *CodeGenerator) addGlobal(name string) bool {
	if _, ok := cg.declGlobals[name]; !ok {
		cg.declGlobals[name] = struct{}{}
		return false
	}
	return true
}

func (cg *CodeGenerator) emitTypeDecl(structType llvmgen.StructType) error {
	if _, ok := cg.declGlobals[structType.Name]; ok {
		return nil
	}
	cg.declTypes[structType.Name] = struct{}{}
	if err := cg.write.TypeDef(structType); err != nil {
		return nil
	}
	return nil
}

func (cg *CodeGenerator) emitFuncDecl(
	returns llvmgen.Type,
	funcName string,
	inputs ...llvmgen.Type,
) error {
	if _, ok := cg.declGlobals[funcName]; ok {
		return nil
	}
	cg.declGlobals[funcName] = struct{}{}
	if err := cg.write.Declare(
		returns,
		llvmgen.Global(funcName),
		inputs...,
	); err != nil {
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
		if err := cg.write.Store(typ, init[0], typ.Ptr(), varPtr); err != nil {
			return err
		}
	}
	return nil
}
