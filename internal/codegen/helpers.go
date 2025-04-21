package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func toLlvmType(typ tast.Type) llvmgen.Type {
	switch typ {
	case tast.Int:
		return llvmgen.I32
	case tast.Bool:
		return llvmgen.I1
	case tast.Double:
		return llvmgen.Double
	case tast.String:
		return llvmgen.I8Ptr
	case tast.Void:
		return llvmgen.Void
	default:
		panic(fmt.Sprintf(
			"Conversion of type %s to LLVM not supported",
			typ.String(),
		))
	}
}
