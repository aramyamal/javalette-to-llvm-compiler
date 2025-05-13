package codegen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func toLlvmType(typ tast.Type) llvmgen.Type {

	switch t := typ.(type) {
	case *tast.ArrayType:
		elemType := toLlvmType(t.Elem)
		name := arrayName(elemType)
		return llvmgen.TypeDef(
			name,           // generated name
			llvmgen.I32,    // length field
			elemType.Ptr(), // pointer to data
		)
	}

	switch typ {
	case tast.Int:
		return llvmgen.I32
	case tast.Bool:
		return llvmgen.I1
	case tast.Double:
		return llvmgen.Double
	case tast.String:
		return llvmgen.I8.Ptr()
	case tast.Void:
		return llvmgen.Void
	default:
		panic(fmt.Sprintf(
			"Conversion of type %s to LLVM not supported",
			typ.String(),
		))
	}
}

func arrayName(elem llvmgen.Type) string {
	arrayRe := regexp.MustCompile(`^arrayof_(.+)_(\d+)D$`)
	name := elem.String()
	if strings.HasPrefix(name, "%") {
		name = name[1:]
	}
	if matches := arrayRe.FindStringSubmatch(name); matches != nil {
		base := matches[1]
		dimStr := matches[2]
		if n, err := strconv.Atoi(dimStr); err == nil {
			return fmt.Sprintf("arrayof_%s_%dD", base, n+1)
		}
	}
	return "arrayof_" + name + "_1D"
}
