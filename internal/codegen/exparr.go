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
	arrayStruct, ok := toLlvmType(e.Type()).(*llvmgen.StructType)
	if !ok {
		return nil, fmt.Errorf(
			"internal compiler error in compileNewArrExp: "+
				"expected llvm struct type for array at %d:%d near %s",
			e.Line(), e.Col(), e.Text(),
		)
	}

	// declare @calloc if not already declared before
	if err := cg.emitFuncDecl(
		llvmgen.I8.Ptr(), "calloc", llvmgen.I64, llvmgen.I64,
	); err != nil {
		return nil, err
	}

	// emit the type declarations of the array wrappers if not already
	// emitted before
	if err := cg.emitArrayTypeDecls(arrayStruct); err != nil {
		return nil, err
	}

	return cg.allocArray(*arrayStruct, indices, 0)
}

func (cg *CodeGenerator) compileArrIndexExp(
	e *tast.ArrIndexExp,
) (llvmgen.Value, error) {

	arrayPtr, err := cg.compileExp(e.Exp)
	if err != nil {
		return nil, err
	}

	currentPtr := arrayPtr
	currentType := toLlvmType(e.Exp.Type())

	for i, idxExp := range e.IdxExps {
		idxValue, err := cg.compileExp(idxExp)
		if err != nil {
			return nil, err
		}

		structType, ok := currentType.(*llvmgen.StructType)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error: expected struct type for array at "+
					"dimension %d, got %s",
				i+1, currentType.String(),
			)
		}

		// pointer to data field
		dataPtr := cg.ng.nextReg()
		if err := cg.write.GetElementPtr(
			dataPtr, structType, structType.Ptr(), currentPtr,
			llvmgen.LitInt(0), llvmgen.LitInt(1),
		); err != nil {
			return nil, err
		}

		// load data field
		dataFieldType := structType.Fields[1]
		dataArray := cg.ng.nextReg()
		ptrType, ok := dataFieldType.(llvmgen.PtrType)
		if !ok {
			return nil, fmt.Errorf("expected pointer type for array data field, got %s", dataFieldType.String())
		}
		if err := cg.write.Load(
			dataArray, ptrType, ptrType, dataPtr,
		); err != nil {
			return nil, err
		}

		// get pointer to element at current index
		elementType, ok := dataFieldType.(llvmgen.PtrType)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error: expected pointer type for data field"+
					" at dimension %d, got %s",
				i+1, dataFieldType.String(),
			)
		}

		elemPtr := cg.ng.nextReg()
		if err := cg.write.GetElementPtr(
			elemPtr, ptrType.Elem, ptrType, dataArray, idxValue,
		); err != nil {
			return nil, err
		}

		// if this is the last dimension, load the value
		if i == len(e.IdxExps)-1 {
			elemValue := cg.ng.nextReg()
			if err := cg.write.Load(
				elemValue, elementType.Elem, elementType.Elem.Ptr(),
				elemPtr,
			); err != nil {
				return nil, err
			}
			return elemValue, nil
		}

		// otherwise load the next array struct pointer
		nextArrayPtr := cg.ng.nextReg()
		if err := cg.write.Load(
			nextArrayPtr, elementType, elementType, elemPtr,
		); err != nil {
			return nil, err
		}

		// update for next iteration
		currentPtr = nextArrayPtr
		currentType = elementType.Elem

	}
	return nil, fmt.Errorf(
		"internal compiler error: no index expressions in array access",
	)
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
	arrayStruct llvmgen.StructType,
	dims []llvmgen.Value,
	level int,
) (llvmgen.Value, error) {
	// get element type which is pointer to the next array struct or primitive
	ptrType, ok := arrayStruct.Fields[1].(llvmgen.PtrType)
	if !ok {
		return nil, fmt.Errorf("")
	}
	elemType := ptrType.Elem

	// emit length for this dimension in I64 to work with calloc
	lengthReg := cg.ng.nextReg()
	if err := cg.write.ZExt(
		lengthReg,
		llvmgen.I32,
		dims[level],
		llvmgen.I64,
	); err != nil {
		return nil, err
	}

	// compute size element in bytes
	elemSize, err := cg.emitSizeOf(elemType)
	if err != nil {
		return nil, err
	}

	// allocate zero intialized memeory with calloc for the data
	dataRaw := cg.ng.nextReg()
	if err := cg.write.Call(
		dataRaw,
		llvmgen.I8.Ptr(),
		"calloc",
		llvmgen.Arg(llvmgen.I64, lengthReg),
		llvmgen.Arg(llvmgen.I64, elemSize),
	); err != nil {
		return nil, err
	}

	// bitcast the I8 pointer from calloc to correct pointer type
	dataTypedPtr := cg.ng.nextReg()
	if err := cg.write.Bitcast(
		dataTypedPtr,
		llvmgen.I8.Ptr(),
		dataRaw,
		elemType.Ptr(),
	); err != nil {
		return nil, err
	}

	// allocate array struct itself on heap
	structSize, err := cg.emitSizeOf(&arrayStruct)
	if err != nil {
		return nil, err
	}
	arrStructRaw := cg.ng.nextReg()
	if err := cg.write.Call(
		arrStructRaw,
		llvmgen.I8.Ptr(),
		"calloc",
		llvmgen.Arg(llvmgen.I64, llvmgen.LitInt(1)),
		llvmgen.Arg(llvmgen.I64, structSize),
	); err != nil {
		return nil, err
	}
	arrStructPtr := cg.ng.nextReg()
	if err := cg.write.Bitcast(
		arrStructPtr,
		llvmgen.I8.Ptr(),
		arrStructRaw,
		(&arrayStruct).Ptr(),
	); err != nil {
		return nil, err
	}

	// set length field (field 0)
	lenFieldPtr := cg.ng.nextReg()
	if err := cg.write.GetElementPtr(
		lenFieldPtr,
		&arrayStruct,
		(&arrayStruct).Ptr(),
		arrStructPtr,
		llvmgen.LitInt(0), llvmgen.LitInt(0),
	); err != nil {
		return nil, err
	}
	if err := cg.write.Store(
		llvmgen.I32,
		dims[level],
		llvmgen.I32.Ptr(),
		lenFieldPtr,
	); err != nil {
		return nil, err
	}

	// set pointer field (field 1)
	ptrFieldPtr := cg.ng.nextReg()
	if err := cg.write.GetElementPtr(
		ptrFieldPtr,
		&arrayStruct,
		(&arrayStruct).Ptr(),
		arrStructPtr,
		llvmgen.LitInt(0), llvmgen.LitInt(1),
	); err != nil {
		return nil, err
	}
	if err := cg.write.Store(
		elemType.Ptr(),
		dataTypedPtr,
		elemType.Ptr().Ptr(),
		ptrFieldPtr,
	); err != nil {
		return nil, err
	}

	// if this is not the innermost dimension, recursively allocate inner arrays
	if level+1 < len(dims) {
		// for (i = 0; i < dims[level]; ++i) {
		//     data[i] = allocArray(nextStruct, dims, level+1)
		// }

		// create loop variable i
		idxVarName := cg.ng.nextTmpVar()
		if err := cg.emitVarAlloc(
			idxVarName, llvmgen.I32, llvmgen.LitInt(0),
		); err != nil {
			return nil, err
		}
		idxPtr, ok := cg.env.LookupVar(idxVarName)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error: could not load tmp variable from" +
					"environment used in multi-dim array allocation",
			)
		}

		// create blocks for looping
		loopHead := cg.ng.nextLab()
		loopBody := cg.ng.nextLab()
		loopExit := cg.ng.nextLab()

		// branch to header
		if err := cg.write.Br(loopHead); err != nil {
			return nil, err
		}

		// in header, compare i < dims[level]
		if err := cg.write.Block(loopHead); err != nil {
			return nil, err
		}
		idxVal := cg.ng.nextReg()
		if err := cg.write.Load(
			idxVal, llvmgen.I32, llvmgen.I32.Ptr(), idxPtr,
		); err != nil {
			return nil, err
		}
		cond := cg.ng.nextReg()
		if err := cg.write.CmpLt(
			cond,
			llvmgen.I32,
			idxVal,
			dims[level],
		); err != nil {
			return nil, err
		}
		if err := cg.write.BrIf(
			llvmgen.I1, cond, loopBody, loopExit,
		); err != nil {
			return nil, err
		}

		// loop body
		cg.write.Block(loopBody)
		elemPtr := cg.ng.nextReg()
		if err := cg.write.GetElementPtr(
			elemPtr, elemType, elemType.Ptr(), dataTypedPtr, idxVal,
		); err != nil {
			return nil, err
		}
		// recursively allocate next dimension
		elemStruct, ok := elemType.(*llvmgen.StructType)
		if !ok {
			return nil, fmt.Errorf(
				"internal compiler error at allocArray:" +
					"could not typecast element type to struct",
			)
		}
		innerArr, err := cg.allocArray(*elemStruct, dims, level+1)
		if err != nil {
			return nil, err
		}
		// store the allocated inner array to elemPtr
		if err := cg.write.Store(
			elemStruct.Ptr(), innerArr,
			elemStruct.Ptr().Ptr(), elemPtr,
		); err != nil {
			return nil, err
		}

		// i++
		nextIdx := cg.ng.nextReg()
		if err := cg.write.Add(
			nextIdx, llvmgen.I32, idxVal, llvmgen.LitInt(1),
		); err != nil {
			return nil, err
		}
		if err := cg.write.Store(
			llvmgen.I32, nextIdx, llvmgen.I32.Ptr(), idxPtr,
		); err != nil {
			return nil, err
		}
		// branch to header
		if err := cg.write.Br(loopHead); err != nil {
			return nil, err
		}

		// set exit block
		if err := cg.write.Block(loopExit); err != nil {
			return nil, err
		}
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
	if err := cg.write.GetElementPtr(
		sizePtrReg,
		typ,
		typ.Ptr(),
		llvmgen.Null(),
		llvmgen.LitInt(1),
	); err != nil {
		return nil, err
	}

	// convert to int
	if err := cg.write.PtrToInt(
		sizeReg,
		typ.Ptr(),
		llvmgen.I64,
		sizePtrReg,
	); err != nil {
		return nil, err
	}

	return sizeReg, nil
}
