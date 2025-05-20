package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileFieldExp(e *tast.FieldExp) (
	llvmgen.Value, error,
) {
	fieldPtr, fieldType, err := cg.emitFieldPtr(
		e, func() (llvmgen.Value, error) { return cg.compileExp(e.Exp) },
	)
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileFieldExp %w", err,
		)
	}
	fieldValue := cg.ng.nextReg()
	cg.write.Load(fieldValue, fieldType, fieldType.Ptr(), fieldPtr)
	return fieldValue, nil
}

func (cg *CodeGenerator) compileFieldLExp(e *tast.FieldExp) (
	llvmgen.Reg, error,
) {
	fieldPtr, _, err := cg.emitFieldPtr(
		e, func() (llvmgen.Value, error) { return cg.compileLExp(e.Exp) },
	)
	if err != nil {
		return "", fmt.Errorf(
			"internal compiler error in compileFieldLExp %w", err,
		)
	}
	return fieldPtr, nil
}

func (cg *CodeGenerator) emitFieldPtr(
	e *tast.FieldExp,
	emitStructPtr func() (llvmgen.Value, error),
) (llvmgen.Reg, llvmgen.Type, error) {
	basePtr, err := emitStructPtr()
	if err != nil {
		return "", nil, err
	}
	fieldProv, ok := e.Exp.Type().(tast.FieldProvider)
	if !ok {
		return "", nil, fmt.Errorf(
			"expected field provider type at %d:%d near %s",
			e.Line(), e.Col(), e.Text(),
		)
	}
	llvmType := toLlvmType(fieldProv)
	fieldInfo, ok := fieldProv.FieldInfo(e.Name)
	if !ok {
		return "", nil, fmt.Errorf(
			"expected field %s but is not available at %d:%d near %s",
			e.Name, e.Line(), e.Col(), e.Text(),
		)
	}
	fieldType := toLlvmType(fieldInfo.Type)
	fieldPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		fieldPtr, llvmType, llvmType.Ptr(), basePtr,
		llvmgen.LitInt(0), llvmgen.LitInt(fieldInfo.Idx),
	)
	return fieldPtr, fieldType, nil
}
