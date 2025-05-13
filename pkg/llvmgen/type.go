package llvmgen

import (
	"fmt"
)

type Type interface {
	String() string
	alignment() int
	ZeroValue() Value
	Ptr() PtrType
}

type PrimitiveType int

const (
	I32 PrimitiveType = iota
	I64
	Double
	I8
	I1
	Void
)

func (t PrimitiveType) String() string {
	switch t {
	case I32:
		return "i32"
	case I64:
		return "i64"
	case Double:
		return "double"
	case I8:
		return "i8"
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
	case I64:
		return 8
	case Double:
		return 8
	case I8:
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
	case I64:
		return LitInt(0)
	case Double:
		return LitDouble(0.0)
	case I8:
		return LitInt(0)
	case I1:
		return LitBool(false)
	case Void:
		panic("void type does not have zero value")
	default:
		panic(fmt.Sprintf("unsupported zero value of type: %d", t))
	}
}

func (t PrimitiveType) Size() int {
	switch t {
	case I32:
		return 4
	case I64:
		return 8
	case Double:
		return 8
	case I8:
		return 1
	case I1:
		return 1
	case Void:
		panic("void type does not have size")
	default:
		panic(fmt.Sprintf("unsupported size of type: %d", t))
	}
}

func (t PrimitiveType) Ptr() PtrType {
	return ptr(t)
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

func (t ArrayType) Ptr() PtrType {
	return ptr(t)
}

type StructType struct {
	Name   string
	Fields []Type
}

func TypeDef(name string, fields ...Type) *StructType {
	return &StructType{Name: name, Fields: fields}
}

func (t StructType) String() string {
	return "%" + t.Name
}

// the alignment should be the maximum alignment of the fields
func (t *StructType) alignment() int {
	maxAlign := 1
	for _, f := range t.Fields {
		if a := f.alignment(); a > maxAlign {
			maxAlign = a
		}
	}
	return maxAlign
}

func (t *StructType) ZeroValue() Value {
	panic("zero value for struct type not yet implemented")
}

func (t *StructType) Size() int {
	panic("size calculation for struct type not yet implemented")
}

func (t *StructType) Ptr() PtrType {
	return ptr(t)
}

type PtrType struct {
	Elem Type
}

func (p PtrType) String() string {
	return p.Elem.String() + "*"
}

func (p PtrType) alignment() int {
	return 1 // pointer alignment (platform dependent, but 1 is safe for LLVM IR)
}

func (p PtrType) ZeroValue() Value {
	return Null()
}

func (p PtrType) Size() int {
	return 8
}

func (p PtrType) Ptr() PtrType {
	return ptr(p)
}

func ptr(elem Type) PtrType {
	return PtrType{Elem: elem}
}

// compile time check for implementation
var _ Type = PrimitiveType(0)
var _ Type = ArrayType{}
var _ Type = &StructType{}
var _ Type = PtrType{}
