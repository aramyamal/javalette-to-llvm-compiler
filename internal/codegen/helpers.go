package codegen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) toLlvmRetType(typ tast.Type) llvmgen.Type {
	if _, isFieldProvider := typ.(tast.FieldProvider); isFieldProvider {
		return cg.toLlvmType(typ).Ptr()
	}
	return cg.toLlvmType(typ)
}

func (cg *CodeGenerator) toLlvmType(typ tast.Type) llvmgen.Type {

	switch t := typ.(type) {
	case *tast.StructType:
		structType, exists := cg.structs[t.Name]
		if exists {
			return structType
		}

		// register empty struct to stop recursion
		structType = llvmgen.StructDef(t.Name)
		cg.structs[t.Name] = structType

		fields := t.Fields()
		fieldLlvmTypes := make([]llvmgen.Type, len(fields))
		for i, fieldName := range fields {
			fieldInfo, ok := t.FieldInfo(fieldName)
			if !ok {
				panic(fmt.Sprintf(
					"unable to access struct field %s from %T",
					fieldName, t,
				))
			}
			fieldLlvmTypes[i] = cg.toLlvmType(fieldInfo.Type)
		}
		structType.Fields = fieldLlvmTypes
		return structType

	case *tast.TypedefType:
		return cg.toLlvmType(UnwrapTypedef(t))

	case *tast.PointerType:
		return cg.toLlvmType(t.Elem).Ptr()

	case *tast.ArrayType:
		elemType := cg.toLlvmType(t.Elem)
		name := arrayName(elemType)
		return llvmgen.StructDef(
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
	if strings.HasSuffix(name, "*") {
		name = name[:len(name)-1]
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

func UnwrapTypedef(t tast.Type) tast.Type {
	for {
		if td, ok := t.(*tast.TypedefType); ok {
			t = td.Aliased
		} else {
			return t
		}
	}
}
