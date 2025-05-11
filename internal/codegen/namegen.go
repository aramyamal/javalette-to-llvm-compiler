package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

// NameGenerator generates unique names for LLVM registers, labels, and string
// literals during code generation.
type NameGenerator struct {
	reg    int
	lab    int
	strIdx int
	strMap map[llvmgen.Global]llvmgen.LitString
	strRev map[string]llvmgen.Global
	ptrMap map[string]int
}

// NewNameGenerator returns a new instance of NameGenerator with all counters
// and maps initialized.
func NewNameGenerator() *NameGenerator {
	return &NameGenerator{
		reg:    0,
		lab:    0,
		strIdx: 0,
		strMap: make(map[llvmgen.Global]llvmgen.LitString),
		strRev: make(map[string]llvmgen.Global),
		ptrMap: make(map[string]int),
	}
}

func (ng *NameGenerator) getOrAddString(content string) (llvmgen.Global, int, bool) {
	if name, ok := ng.strRev[content]; ok {
		return name, len(content) + 1, true // string already existed
	}
	name := llvmgen.Global(fmt.Sprintf("s_%d", ng.strIdx))
	ng.strIdx++
	ng.strMap[name] = llvmgen.LitString(content)
	ng.strRev[content] = name
	return name, len(content) + 1, false // new string
}

func (ng *NameGenerator) ptrName(name string) llvmgen.Reg {
	ptrCount := ng.ptrMap[name]
	ng.ptrMap[name] = ptrCount + 1
	if ptrCount == 0 {
		return llvmgen.Reg(fmt.Sprintf(".%s_p", name))
	}
	return llvmgen.Reg(fmt.Sprintf(".%s_p%d", name, ptrCount))
}

func (ng *NameGenerator) resetPtrs() {
	ng.ptrMap = make(map[string]int)
}

func (ng *NameGenerator) nextReg() llvmgen.Reg {
	regName := fmt.Sprintf("t%d", ng.reg)
	ng.reg++
	return llvmgen.Reg(regName)
}

func (ng *NameGenerator) nextLab() string {
	labName := fmt.Sprintf("l%d", ng.lab)
	ng.lab++
	return labName
}

func (ng *NameGenerator) resetReg() {
	ng.reg = 0
}

func (ng *NameGenerator) resetLab() {
	ng.lab = 0
}
