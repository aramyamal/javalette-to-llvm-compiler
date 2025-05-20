package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileNewArrExp(
	e *tast.NewArrExp,
) (llvmgen.Value, error) {

	// calculate array lengths
	var indices []llvmgen.Value
	for _, exp := range e.Exps {
		index, err := cg.compileExp(exp)
		if err != nil {
			return nil, err
		}
		indices = append(indices, index)
	}
	arrStructType, ok := toLlvmType(e.Type()).(*llvmgen.StructType)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in compileNewArrExp: "+
				"expected llvm struct type for array at %d:%d near %s",
			e.Line(), e.Col(), e.Text(),
		)
	}

	arrStructPtr, err := cg.allocArray(arrStructType, indices, 0)
	if err != nil {
		return nil, fmt.Errorf(
			"%w at %d:%d near %s", err, e.Line(), e.Col(), e.Text(),
		)
	}
	return arrStructPtr, nil
}

func (cg *CodeGenerator) compileArrIndexExp(
	e *tast.ArrIndexExp,
) (llvmgen.Reg, error) {

	arrPtr, err := cg.compileExp(e.Exp)
	if err != nil {
		return "", err
	}
	arrType := toLlvmType(e.Exp.Type())

	elemPtr, elemType, err := cg.emitArrElemPtr(arrPtr, arrType, e.IdxExps)
	if err != nil {
		return "", fmt.Errorf(
			"internal compiler error in compileArrIndexExp: %w at %d:%d near"+
				" %s", err, e.Line(), e.Col(), e.Text(),
		)
	}
	// if the element is an array/struct elemptr is a pointer to a pointer to
	// the array struct
	if structType, ok := elemType.(*llvmgen.StructType); ok {
		// arrPtr is pointer to actual struct
		arrPtr := cg.ng.nextReg()
		cg.write.Load(arrPtr, structType.Ptr(), structType.Ptr().Ptr(), elemPtr)
		return arrPtr, nil
	}

	// otherwise load the primitive type value
	elemValue := cg.ng.nextReg()
	cg.write.Load(elemValue, elemType, elemType.Ptr(), elemPtr)
	return elemValue, nil
}

func (cg *CodeGenerator) compileArrIndexLExp(
	e *tast.ArrIndexExp,
) (llvmgen.Reg, error) {

	arrPtr, err := cg.compileExp(e.Exp)
	if err != nil {
		return "", err
	}
	arrType := toLlvmType(e.Exp.Type())

	elemPtr, _, err := cg.emitArrElemPtr(arrPtr, arrType, e.IdxExps)
	if err != nil {
		return "", fmt.Errorf(
			"internal compiler error in compileArrIndexLExp: %w at %d:%d near"+
				" %s", err, e.Line(), e.Col(), e.Text(),
		)
	}
	return elemPtr, nil
}

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

func (cg *CodeGenerator) emitArrayTypeDecls(typ llvmgen.Type) error {
	structType, ok := typ.(*llvmgen.StructType)
	if !ok {
		return fmt.Errorf(
			"internal compiler error: emitArrayTypeDecls expected a StructType"+
				"but received %T", typ,
		)
	}

	if err := cg.emitTypeDecl(*structType); err != nil {
		return err
	}
	// if the second field is a pointer
	if ptrType, ok := structType.Fields[1].(llvmgen.PtrType); ok {
		// if pointer points to element that is a struct do recursive call
		if _, isStruct := ptrType.Elem.(*llvmgen.StructType); isStruct {
			return cg.emitArrayTypeDecls(ptrType.Elem)
		}
	}
	return nil
}

func (cg *CodeGenerator) allocArray(
	arrStructType *llvmgen.StructType,
	dims []llvmgen.Value,
	level int,
) (llvmgen.Value, error) {

	// declare @calloc if not already declared before
	if err := cg.emitFuncDecl(
		llvmgen.I8.Ptr(), "calloc", llvmgen.I64, llvmgen.I64,
	); err != nil {
		return nil, err
	}

	// emit the type declarations of the array wrappers if not already
	// emitted before
	if err := cg.emitArrayTypeDecls(arrStructType); err != nil {
		return nil, err
	}

	// handle the empty dims case: allocate an empty array struct
	if len(dims) == 0 {
		// allocate array struct on heap
		structSize, _ := cg.emitSizeOf(arrStructType)
		arrStructPtr, _ := cg.emitCalloc(
			llvmgen.LitInt(1), structSize, arrStructType,
		)

		// set length field (field 0) to 0
		lenFieldPtr := cg.ng.nextReg()
		cg.write.GetElementPtr(
			lenFieldPtr, arrStructType, arrStructType.Ptr(), arrStructPtr,
			llvmgen.LitInt(0), llvmgen.LitInt(0),
		)
		cg.write.Store(llvmgen.I32, llvmgen.LitInt(0), llvmgen.I32.Ptr(), lenFieldPtr)

		// set pointer field (field 1) to null
		ptrFieldPtr := cg.ng.nextReg()
		cg.write.GetElementPtr(
			ptrFieldPtr, arrStructType, arrStructType.Ptr(), arrStructPtr,
			llvmgen.LitInt(0), llvmgen.LitInt(1),
		)
		ptrType, ok := arrStructType.Fields[1].(llvmgen.PtrType)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error in allocArray: expected pointer type for"+
					"array data field (field 1), but got %s",
				arrStructType.Fields[1].String(),
			)
		}
		cg.write.Store(ptrType, llvmgen.Null(), ptrType.Ptr(), ptrFieldPtr)

		return arrStructPtr, nil
	}

	// get element type which is pointer to the next array struct or primitive
	ptrType, ok := arrStructType.Fields[1].(llvmgen.PtrType)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in allocArray: expected pointer type for"+
				"array data field (field 1), but got %s",
			arrStructType.Fields[1].String(),
		)
	}
	elemType := ptrType.Elem

	// emit length for this dimension in I64 to work with calloc
	lengthReg := cg.ng.nextReg()
	cg.write.ZExt(lengthReg, llvmgen.I32, dims[level], llvmgen.I64)

	// compute size element in bytes
	elemSize, _ := cg.emitSizeOf(elemType)

	// allocate data array
	dataTypedPtr, _ := cg.emitCalloc(lengthReg, elemSize, elemType)

	// allocate array struct itself on heap
	structSize, _ := cg.emitSizeOf(arrStructType)
	arrStructPtr, _ := cg.emitCalloc(
		llvmgen.LitInt(1), structSize, arrStructType,
	)

	// set length field (field 0)
	lenFieldPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		lenFieldPtr, arrStructType, arrStructType.Ptr(), arrStructPtr,
		llvmgen.LitInt(0), llvmgen.LitInt(0),
	)
	cg.write.Store(llvmgen.I32, dims[level], llvmgen.I32.Ptr(), lenFieldPtr)

	// set pointer field (field 1)
	ptrFieldPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		ptrFieldPtr, arrStructType, arrStructType.Ptr(), arrStructPtr,
		llvmgen.LitInt(0), llvmgen.LitInt(1),
	)
	cg.write.Store(
		elemType.Ptr(), dataTypedPtr, elemType.Ptr().Ptr(), ptrFieldPtr,
	)

	// if this is not the innermost dimension, recursively allocate inner arrays
	if level+1 < len(dims) {
		// for (i = 0; i < dims[level]; ++i) {
		//     data[i] = allocArray(nextStruct, dims, level+1)
		// }

		// create loop variable i
		idxVarName := cg.ng.nextTmpVar()
		idxPtr, err := cg.emitVarAlloc(
			idxVarName, llvmgen.I32, llvmgen.LitInt(0))
		if err != nil {
			return nil, err
		}

		// create blocks for looping
		loopHead := cg.ng.nextLab()
		loopBody := cg.ng.nextLab()
		loopExit := cg.ng.nextLab()

		// branch to header
		cg.write.Br(loopHead)

		// in header, compare i < dims[level]
		cg.write.Block(loopHead)

		idxVal := cg.ng.nextReg()
		cg.write.Load(idxVal, llvmgen.I32, llvmgen.I32.Ptr(), idxPtr)

		cond := cg.ng.nextReg()
		cg.write.CmpLt(cond, llvmgen.I32, idxVal, dims[level])
		cg.write.BrIf(llvmgen.I1, cond, loopBody, loopExit)

		// loop body
		cg.write.Block(loopBody)
		elemPtr := cg.ng.nextReg()
		cg.write.GetElementPtr(
			elemPtr, elemType, elemType.Ptr(), dataTypedPtr,
			idxVal,
		)
		// recursively allocate next dimension
		elemStruct, ok := elemType.(*llvmgen.StructType)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error at allocArray:" +
					"could not typecast element type to struct",
			)
		}
		innerArr, err := cg.allocArray(elemStruct, dims, level+1)
		if err != nil {
			return nil, err
		}
		// store the allocated inner array to elemPtr
		cg.write.Store(
			elemStruct.Ptr(), innerArr, elemStruct.Ptr().Ptr(), elemPtr,
		)

		// i++
		nextIdx := cg.ng.nextReg()
		cg.write.Add(nextIdx, llvmgen.I32, idxVal, llvmgen.LitInt(1))
		cg.write.Store(llvmgen.I32, nextIdx, llvmgen.I32.Ptr(), idxPtr)

		// branch to header
		cg.write.Br(loopHead)

		// set exit block
		cg.write.Block(loopExit)
	}

	return arrStructPtr, nil
}

func (cg *CodeGenerator) emitSizeOf(typ llvmgen.Type) (llvmgen.Value, error) {

	// use known size for primitive types
	if primType, ok := typ.(llvmgen.PrimitiveType); ok {
		return llvmgen.LitInt(primType.Size()), nil
	}
	// for other types, use getelementptr and ptrtoint trick

	sizeReg := cg.ng.nextReg()
	// emit ptr that with null as base pointer which gives address just
	// past the first element, that is size of the type
	sizePtrReg := cg.ng.nextReg()
	cg.write.GetElementPtr(
		sizePtrReg, typ, typ.Ptr(), llvmgen.Null(),
		llvmgen.LitInt(1),
	)

	// convert to int
	cg.write.PtrToInt(sizeReg, typ.Ptr(), llvmgen.I64, sizePtrReg)

	return sizeReg, nil
}

func (cg *CodeGenerator) emitCalloc(
	numElems llvmgen.Value,
	elemSize llvmgen.Value,
	resultType llvmgen.Type,
) (llvmgen.Value, error) {
	// allocate zero intialized memeory with calloc for the data
	raw := cg.ng.nextReg()
	cg.write.Call(
		raw, llvmgen.I8.Ptr(), "calloc",
		llvmgen.Arg(llvmgen.I64, numElems),
		llvmgen.Arg(llvmgen.I64, elemSize),
	)

	// bitcast the I8 pointer from calloc to correct pointer type
	typed := cg.ng.nextReg()
	cg.write.Bitcast(typed, llvmgen.I8.Ptr(), raw, resultType.Ptr())

	return typed, nil
}

func (cg *CodeGenerator) emitUninitStruct(
	typ tast.Type,
) (llvmgen.Value, error) {
	llvmType := toLlvmType(typ)

	// handle array types
	if llvmStructType, isStruct := llvmType.(*llvmgen.StructType); isStruct {
		if _, isTastArray := typ.(*tast.ArrayType); isTastArray {
			// lllocate an empty array (length 0, data null)
			arrPtr, err := cg.allocArray(llvmStructType, []llvmgen.Value{}, 0)
			if err != nil {
				return nil, fmt.Errorf(
					"emitUninitStruct: failed to allocate empty array: %w", err)
			}
			return arrPtr, nil
		}

		return nil, fmt.Errorf(
			"emitUninitStruct: unsupported struct type (type: %s)",
			typ.String(),
		)
	}

	// not a struct/array type
	return nil, fmt.Errorf(
		"emitUninitStruct: not a struct or array type (type: %s)",
		typ.String(),
	)
}

func (cg *CodeGenerator) emitArrElemPtr(
	arrPtr llvmgen.Value,
	arrType llvmgen.Type,
	idxExps []tast.Exp,
) (elemPtr llvmgen.Reg, elemElemType llvmgen.Type, err error) {
	currentPtr := arrPtr
	currentType := arrType

	for i, idxExp := range idxExps {
		idxValue, err := cg.compileExp(idxExp)
		if err != nil {
			return "", nil, err
		}

		structType, ok := currentType.(*llvmgen.StructType)
		if !ok {
			return "", nil, fmt.Errorf(
				"expected struct type for array at dimension %d, got %s",
				i+1, currentType.String(),
			)
		}

		// pointer to data field
		dataPtr := cg.ng.nextReg()
		cg.write.GetElementPtr(
			dataPtr, structType, structType.Ptr(), currentPtr,
			llvmgen.LitInt(0), llvmgen.LitInt(1),
		)

		// load data field
		dataFieldType := structType.Fields[1]
		dataArray := cg.ng.nextReg()
		ptrType, ok := dataFieldType.(llvmgen.PtrType)
		if !ok {
			return "", nil, fmt.Errorf(
				"expected pointer type for array data field, got %s",
				dataFieldType.String(),
			)
		}
		cg.write.Load(dataArray, ptrType, ptrType.Ptr(), dataPtr)

		// get pointer to element at current index
		elementType, ok := dataFieldType.(llvmgen.PtrType)
		if !ok {
			return "", nil, fmt.Errorf(
				"expected pointer type for data field at dimension %d, got %s",
				i+1, dataFieldType.String(),
			)
		}

		elemPtr := cg.ng.nextReg()
		cg.write.GetElementPtr(
			elemPtr, ptrType.Elem, ptrType, dataArray,
			idxValue,
		)

		// If this is the last dimension, return the pointer and type
		if i == len(idxExps)-1 {
			return elemPtr, elementType.Elem, nil
		}

		// otherwise load the next array struct pointer
		nextArrayPtr := cg.ng.nextReg()
		cg.write.Load(nextArrayPtr, elementType, elementType.Ptr(), elemPtr)

		// update for next iteration
		currentPtr = nextArrayPtr
		currentType = elementType.Elem
	}
	return "", nil, fmt.Errorf("no index expressions in array access")
}
