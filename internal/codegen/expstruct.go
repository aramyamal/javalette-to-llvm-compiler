package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileNewStructExp(
	e *tast.NewStructExp,
) (llvmgen.Value, error) {
	ptrType, ok := cg.toLlvmType(e.Type()).(llvmgen.PtrType)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in compileNewStructExp: "+
				"expected llvm pointer type at %d:%d near %s "+
				"but got type %T",
			e.Line(), e.Col(), e.Text(), cg.toLlvmType(e.Type()),
		)
	}

	structType, ok := ptrType.Elem.(*llvmgen.StructType)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in compileNewStructExp: "+
				"expected llvm pointer to struct type at %d:%d near %s",
			e.Line(), e.Col(), e.Text(),
		)
	}

	structSize, err := cg.emitSizeOf(structType)
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileNewStructExp: %w at %d:%d at %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}

	structPtr, err := cg.emitCalloc(llvmgen.LitInt(1), structSize, structType)
	if err != nil {
		return nil, fmt.Errorf(
			"internal compiler error in compileNewStructExp: %w at %d:%d at %s",
			err, e.Line(), e.Col(), e.Text(),
		)
	}

	return structPtr, nil
}

func (cg *CodeGenerator) compileDerefExp(
	e *tast.DerefExp,
) (llvmgen.Value, error) {
	fieldPtr, err := cg.compileDerefLExp(e)
	if err != nil {
		return nil, fmt.Errorf("compileDerefExp: %w", err)
	}

	typ := cg.toLlvmRetType(e.Type())

	value := cg.ng.nextReg()
	cg.write.Load(value, typ, typ.Ptr(), fieldPtr)
	return value, nil
}

func (cg *CodeGenerator) compileDerefLExp(
	e *tast.DerefExp,
) (llvmgen.Reg, error) {
	structPtr, err := cg.compileExp(e.Exp)
	if err != nil {
		return "", err
	}

	llvmType := cg.toLlvmType(e.Exp.Type())
	structPtrType, ok1 := llvmType.(llvmgen.PtrType)
	if !ok1 {
		return "", fmt.Errorf(
			"compileDerefLExp: expected struct pointer type at %d:%d near '%s' (got %T)",
			e.Line(), e.Col(), e.Text(), llvmType,
		)
	}
	structType, ok2 := structPtrType.Elem.(*llvmgen.StructType)
	if !ok2 {
		return "", fmt.Errorf(
			"compileDerefLExp: expected struct type at %d:%d near '%s' (got %T)",
			e.Line(), e.Col(), e.Text(), structPtrType.Elem,
		)
	}

	ptrTastType, ok1 := e.Exp.Type().(*tast.PointerType)
	if !ok1 {
		return "", fmt.Errorf(
			"compileDerefLExp: expected pointer TAST type at %d:%d near '%s'",
			e.Line(), e.Col(), e.Text(),
		)
	}
	fieldProv, ok2 := ptrTastType.Elem.(tast.FieldProvider)
	if !ok2 {
		return "", fmt.Errorf(
			"compileDerefLExp: type does not provide fields at %d:%d near '%s'",
			e.Line(), e.Col(), e.Text(),
		)
	}
	fieldInfo, ok := fieldProv.FieldInfo(e.Name)
	if !ok {
		return "", fmt.Errorf(
			"compileDerefLExp: field %s not found at %d:%d near '%s'",
			e.Name, e.Line(), e.Col(), e.Text(),
		)
	}

	fieldPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		fieldPtr, structType, structType.Ptr(), structPtr,
		llvmgen.LitInt(0), llvmgen.LitInt(fieldInfo.Idx),
	)

	return fieldPtr, nil
}
