package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileForEachStm(s *tast.ForEachStm) error {
	cg.env.EnterContext()
	defer cg.env.ExitContext()

	arr, err := cg.compileExp(s.Exp)
	if err != nil {
		return err
	}

	// assert array struct type
	structType, ok := toLlvmType(s.Exp.Type()).(*llvmgen.StructType)
	if !ok {
		return fmt.Errorf(
			"internal compiler error in compileForEachStm: "+
				"expected llvm struct type for array at %d:%d near %s",
			s.Line(), s.Col(), s.Text(),
		)
	}
	// pointer to data field
	dataPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		dataPtr, structType, structType.Ptr(), arr,
		llvmgen.LitInt(0), llvmgen.LitInt(1),
	)
	// load data field
	dataFieldType := structType.Fields[1]
	dataArray := cg.ng.nextReg()
	ptrType, ok := dataFieldType.(llvmgen.PtrType)
	if !ok {
		return fmt.Errorf(
			"expected pointer type for array data field, got %s at %d:%d at %s",
			dataFieldType.String(), s.Col(), s.Line(), s.Text(),
		)
	}
	cg.write.Load(dataArray, ptrType, ptrType.Ptr(), dataPtr)

	// get array length
	lengthPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		lengthPtr, structType, structType.Ptr(), arr,
		llvmgen.LitInt(0), llvmgen.LitInt(0),
	)
	length := cg.ng.nextReg()
	cg.write.Load(length, llvmgen.I32, llvmgen.I32.Ptr(), lengthPtr)

	// declare/allocate variable
	elemType := toLlvmType(s.Type)
	// for arrays, the type is a pointer
	if _, isArray := s.Type.(*tast.ArrayType); isArray {
		elemType = elemType.Ptr()
	}
	variablePtr, err := cg.emitVarAlloc(s.Id, elemType, elemType.ZeroValue())
	if err != nil {
		return err
	}

	// create loop variable i
	idxVarName := cg.ng.nextTmpVar()
	idxPtr, err := cg.emitVarAlloc(
		idxVarName, llvmgen.I32, llvmgen.LitInt(0))
	if err != nil {
		return err
	}

	// create blocks for looping
	loopHead := cg.ng.nextLab()
	loopBody := cg.ng.nextLab()
	loopExit := cg.ng.nextLab()

	// branch to header
	cg.write.Br(loopHead)

	// in header, compare i < length
	cg.write.Block(loopHead)
	idxVal := cg.ng.nextReg()
	cg.write.Load(idxVal, llvmgen.I32, llvmgen.I32.Ptr(), idxPtr)

	cond := cg.ng.nextReg()
	cg.write.CmpLt(cond, llvmgen.I32, idxVal, length)
	cg.write.BrIf(llvmgen.I1, cond, loopBody, loopExit)

	// loop body
	cg.write.Block(loopBody)

	// get pointer to element at current index
	elemPtr := cg.ng.nextReg()
	cg.write.GetElementPtr(
		elemPtr, ptrType.Elem, ptrType, dataArray,
		idxVal,
	)

	variableValue := cg.ng.nextReg()
	if _, isStruct := elemType.(*llvmgen.StructType); isStruct {
		// for arrays/structs, elemPtr is a pointer to a pointer to the struct,
		// so load the pointer from elemPtr
		cg.write.Load(variableValue, ptrType, ptrType.Ptr(), elemPtr)
		cg.write.Store(ptrType, variableValue, ptrType.Ptr(), variablePtr)
	} else {
		// for primitive types, load the value
		cg.write.Load(variableValue, elemType, elemType.Ptr(), elemPtr)
		cg.write.Store(elemType, variableValue, elemType.Ptr(), variablePtr)
	}

	if err := cg.compileStm(s.Stm); err != nil {
		return err
	}

	// i++
	nextIdx := cg.ng.nextReg()
	cg.write.Add(nextIdx, llvmgen.I32, idxVal, llvmgen.LitInt(1))
	cg.write.Store(llvmgen.I32, nextIdx, llvmgen.I32.Ptr(), idxPtr)

	// branch to header
	cg.write.Br(loopHead)

	// set exit block
	cg.write.Block(loopExit)

	return nil
}
