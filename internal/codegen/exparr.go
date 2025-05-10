package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileNewArrExp(
	e *tast.NewArrExp,
) (llvmgen.Value, error) {
	return nil, fmt.Errorf("compileNewArrExp: not yet implemented")
}

func (cg *CodeGenerator) compileArrIndexExp(
	e *tast.ArrIndexExp,
) (llvmgen.Value, error) {
	return nil, fmt.Errorf("compileArrIndexExp: not yet implemented")
}

func (cg *CodeGenerator) compileArrPostExp(
	e *tast.ArrPostExp,
) (llvmgen.Value, error) {
	return nil, fmt.Errorf("compileArrPostExp: not yet implemented")
}

func (cg *CodeGenerator) compileArrPreExp(
	e *tast.ArrPreExp,
) (llvmgen.Value, error) {
	return nil, fmt.Errorf("compileArrPreExp: not yet implemented")
}

func (cg *CodeGenerator) compileArrAssignExp(
	e *tast.ArrAssignExp,
) (llvmgen.Value, error) {
	return nil, fmt.Errorf("compileArrAssignExp: not yet implemented")
}
