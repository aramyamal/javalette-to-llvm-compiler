package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvm"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

func toLlvmType(typ types.Type) llvm.Type {
	switch typ {
	case types.Int:
		return llvm.I32
	case types.Bool:
		return llvm.I1
	case types.Double:
		return llvm.Double
	case types.String:
		return llvm.I8Ptr
	case types.Void:
		return llvm.Void
	default:
		panic(fmt.Sprintf(
			"Conversion of type %s to LLVM not supported",
			typ.String(),
		))
	}
}
