package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
)

type NameGenerator struct {
	reg int
	lab int
}

func (ng *NameGenerator) nextReg() llvm.Var {
	ng.reg++
	return llvm.Var(fmt.Sprintf("t%d", ng.reg))
}

func (ng *NameGenerator) nextLab() string {
	ng.lab++
	return fmt.Sprintf("l%d", ng.lab)
}

func (ng *NameGenerator) currentReg() llvm.Var {
	return llvm.Var(fmt.Sprintf("t%d", ng.reg))
}

func (ng *NameGenerator) currentLab() string {
	return fmt.Sprintf("l%d", ng.lab)
}

func (ng *NameGenerator) resetReg() {
	ng.reg = 0
}

func (ng *NameGenerator) resetLab() {
	ng.lab = 0
}
