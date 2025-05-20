package tast

import "strings"

// Type represents a Javalette type in the typed abstract syntax tree (TAST).
type Type interface {
	String() string
	isTastType()
}

type FieldInfo struct {
	Type Type
	Idx  int
}

type FieldProvider interface {
	Type
	FieldInfo(name string) (*FieldInfo, bool)
	Fields() []string
}

type BaseType int

const (
	Unknown BaseType = iota
	Int
	Double
	Bool
	String
	Void
)

func (b BaseType) String() string {
	return [...]string{
		"Unknown",
		"Int",
		"Double",
		"Bool",
		"String",
		"Void",
	}[b]
}
func (b BaseType) isTastType() {}

type ArrayType struct {
	Elem Type
}

func (a *ArrayType) String() string {
	return a.Elem.String() + "[]"
}
func (b ArrayType) isTastType() {}

func (a *ArrayType) FieldInfo(name string) (*FieldInfo, bool) {
	if name == "length" {
		return &FieldInfo{Type: Int, Idx: 0}, true
	}
	return nil, false
}

func (a *ArrayType) Fields() []string {
	return []string{"length"}
}

func Array(elem Type) *ArrayType {
	return &ArrayType{Elem: elem}
}

// check that ArrayType implements FieldProvider
var _ FieldProvider = (*ArrayType)(nil)

type StructType struct {
	Name       string
	fields     map[string]*FieldInfo
	fieldNames []string
}

type fieldCreator struct {
	Type Type
	Name string
}

func Field(typ Type, name string) *fieldCreator {
	return &fieldCreator{Type: typ, Name: name}
}

func Struct(name string, fields ...*fieldCreator) *StructType {
	fieldsMap := make(map[string]*FieldInfo)
	var fieldNames []string
	for i, field := range fields {
		fieldsMap[field.Name] = &FieldInfo{Type: field.Type, Idx: i}
		fieldNames = append(fieldNames, field.Name)
	}
	return &StructType{Name: name, fields: fieldsMap, fieldNames: fieldNames}
}

func (s *StructType) String() string {
	result := "struct " + s.Name + " { "
	var fields []string
	for _, name := range s.fieldNames {
		fields = append(fields, name+": "+s.fields[name].Type.String())
	}
	result += strings.Join(fields, "; ")
	result += " }"
	return result
}

func (s *StructType) isTastType() {}

func (s *StructType) FieldInfo(name string) (*FieldInfo, bool) {
	fieldInfo, ok := s.fields[name]
	return fieldInfo, ok
}

func (s *StructType) Fields() []string {
	return s.fieldNames
}

type TypedefType struct {
	Name    string
	Aliased Type
}

func Typedef(name string, aliased Type) *TypedefType {
	return &TypedefType{Name: name, Aliased: aliased}
}

func (t *TypedefType) String() string {
	return t.Name
}

func (t *TypedefType) isTastType() {}

// Op represents an operator in the TAST.
type Op int

const (
	OpInc Op = iota // ++ increment
	OpDec           // -- decrement
	OpMul           // * multiplication
	OpDiv           // / division
	OpMod           // % modulo
	OpAdd           // + addition
	OpSub           // - subtraction
	OpLt            // < less than
	OpGt            // > greater than
	OpLe            // <= less than or equal
	OpGe            // >= greater than or equal
	OpEq            // == equal
	OpNe            // != not equal
)

// String returns the symbol of the operator.
func (op Op) String() string {
	return [...]string{
		"++", // Inc
		"--", // Dec
		"*",  // Mul
		"/",  // Div
		"%",  // Mod
		"+",  // Add
		"-",  // Sub
		"<",  // Lt
		">",  // Gt
		"<=", // Le
		">=", // Ge
		"==", // Eq
		"!=", // Ne
	}[op]
}

// Name returns the type name of the operator.
func (op Op) Name() string {
	return [...]string{
		"OpInc",
		"OpDec",
		"OpMul",
		"OpDiv",
		"OpMod",
		"OpAdd",
		"OpSub",
		"OpLt",
		"OpGt",
		"OpLe",
		"OpGe",
		"OpEq",
		"OpNe",
	}[op]
}

// compile time interface implementation checks
var _ Type = BaseType(0)
var _ Type = (*ArrayType)(nil)
var _ FieldProvider = (*ArrayType)(nil)
var _ Type = (*StructType)(nil)
var _ FieldProvider = (*StructType)(nil)
var _ Type = (*TypedefType)(nil)
