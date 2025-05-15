package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileFieldExp(e *tast.FieldExp) (
	llvmgen.Value, error,
) {
	expValue, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}
	// assert field provider
	fieldProv, ok := e.Exp.Type().(tast.FieldProvider)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in compileFieldExp: "+
				"expected field provider type at %d:%d near %s",
			e.Line(), e.Col(), e.Text(),
		)
	}
	llvmType := toLlvmType(fieldProv)
	fieldInfo, ok := fieldProv.FieldInfo(e.Name)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in compileFieldExp: "+
				"expected field %s but is not available at %d:%d near %s",
			e.Name, e.Line(), e.Col(), e.Text(),
		)
	}
	fieldType := toLlvmType(fieldInfo.Type)
	// pointer to field
	fieldPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		fieldPtr, llvmType, llvmType.Ptr(), expValue,
		llvmgen.LitInt(0), llvmgen.LitInt(fieldInfo.Idx),
	)
	// load data field
	fieldValue := cg.ng.nextReg()
	cg.write.Load(fieldValue, fieldType, fieldType.Ptr(), fieldPtr)

	return fieldValue, nil
}

