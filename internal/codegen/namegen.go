package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

type NameGenerator struct {
	reg    int
	lab    int
	strIdx int
	strMap map[llvm.Global]llvm.LitString
	ptrMap map[string]int
}

func NewNameGenerator() *NameGenerator {
	return &NameGenerator{
		reg:    0,
		lab:    0,
		strIdx: 0,
		strMap: make(map[llvm.Global]llvm.LitString),
		ptrMap: make(map[string]int),
	}
}

func (ng *NameGenerator) addString(content string) (llvm.Global, int) {
	name := llvm.Global(fmt.Sprintf("s_%d", ng.strIdx))
	ng.strIdx++
	ng.strMap[name] = llvm.LitString(content)
	return llvm.Global(name), len(content) + 1
}

func (ng *NameGenerator) ptrName(name string) llvm.Reg {
	ptrCount := ng.ptrMap[name]
	if ptrCount == 0 {
		return llvm.Reg(fmt.Sprintf(".%s_ptr", name))
	}
	return llvm.Reg(fmt.Sprintf(".%s_%dptr", name, ptrCount))
}

func (ng *NameGenerator) resetPtrs() {
	ng.ptrMap = make(map[string]int)
}

func (ng *NameGenerator) resetStrings() {
	ng.strMap = make(map[llvm.Global]llvm.LitString)
}

func (ng *NameGenerator) nextReg() llvm.Reg {
	regName := fmt.Sprintf("t%d", ng.reg)
	ng.reg++
	return llvm.Reg(regName)
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
