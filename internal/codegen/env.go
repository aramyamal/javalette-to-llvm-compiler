package codegen

import "github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"

type CodegenContext map[string]llvmgen.Reg

type CodegenEnv struct {
	contexts []CodegenContext
}

func (e *CodegenEnv) EnterContext() {
	e.contexts = append(e.contexts, make(CodegenContext))
}

func (e *CodegenEnv) ExitContext() {
	e.contexts = e.contexts[:len(e.contexts)-1]
}

func (e *CodegenEnv) LookupVar(name string) (llvmgen.Reg, bool) {
	for i := len(e.contexts) - 1; i >= 0; i-- {
		if reg, ok := e.contexts[i][name]; ok {
			return reg, true
		}
	}
	return "", false
}

func (e *CodegenEnv) AddVar(name string, reg llvmgen.Reg) {
	e.contexts[len(e.contexts)-1][name] = reg
}

func (e *CodegenEnv) ExtendVar(varName string, value llvmgen.Reg) bool {
	if len(e.contexts) == 0 {
		return false
	}
	ctx, ok := e.Peek()
	if !ok {
		return false
	}
	if _, ok := (*ctx)[varName]; ok {
		return false
	}
	e.contexts[len(e.contexts)-1][varName] = value
	return true
}

func (e CodegenEnv) Peek() (*CodegenContext, bool) {
	if len(e.contexts) == 0 {
		return nil, false
	}
	return &e.contexts[len(e.contexts)-1], true
}

func NewCodegenEnv() *CodegenEnv {
	return &CodegenEnv{
		contexts: []CodegenContext{make(CodegenContext)},
	}
}

