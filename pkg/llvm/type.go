package llvm

import "fmt"

type Type interface {
	String() string
	alignment() int
	ZeroValue() Value
}

type PrimitiveType int

const (
	I32 PrimitiveType = iota
	Double
	I8
	I8Ptr
	I1
	Void
)

func (t PrimitiveType) String() string {
	switch t {
	case I32:
		return "i32"
	case Double:
		return "double"
	case I8:
		return "i8"
	case I8Ptr:
		return "i8*"
	case I1:
		return "i1"
	case Void:
		return "void"
	default:
		panic(fmt.Sprintf("unsupported type: %d", t))
	}
}

func (t PrimitiveType) alignment() int {
	switch t {
	case I32:
		return 4
	case Double:
		return 8
	case I8:
		return 1
	case I8Ptr:
		return 1
	case I1:
		return 1
	case Void:
		panic("void type does not have alignment")
	default:
		panic(fmt.Sprintf("unsupported alignment of type: %d", t))
	}
}

func (t PrimitiveType) ZeroValue() Value {
	switch t {
	case I32:
		return LitInt(0)
	case Double:
		return LitDouble(0.0)
	case I8:
		return LitInt(0)
	case I8Ptr:
		panic("I8 pointer type does not have zero value")
	case I1:
		return LitBool(false)
	case Void:
		panic("void type does not have zero value")
	default:
		panic(fmt.Sprintf("unsupported zero value of type: %d", t))
	}
}

type ArrayType struct {
	typ        Type
	dimensions []int
}

func Array(llvmType Type, dimensions ...int) ArrayType {
	return ArrayType{
		typ:        llvmType,
		dimensions: dimensions,
	}
}

func (t ArrayType) String() string {
	typeStr := t.typ.String()
	for i := len(t.dimensions) - 1; i >= 0; i-- {
		typeStr = fmt.Sprintf("[%d x %s]", t.dimensions[i], typeStr)
	}
	return typeStr
}

func (t ArrayType) alignment() int {
	return t.typ.alignment()
}

func (t ArrayType) ZeroValue() Value {
	panic("zero value of array not yet implemented")
}
