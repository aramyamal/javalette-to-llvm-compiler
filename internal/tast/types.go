package tast

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
	return typeSummary(a.Elem) + "[]"
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

type StructType struct {
	Name string
}

// global mapping of struct fields to structs to prevent recursion problems
var structFields = map[string][]*FieldCreator{}

func RegisterStruct(name string, fields ...*FieldCreator) *StructType {
	structFields[name] = fields
	return &StructType{Name: name}
}

type FieldCreator struct {
	Type Type
	Name string
}

func Field(typ Type, name string) *FieldCreator {
	return &FieldCreator{Type: typ, Name: name}
}

func (s *StructType) String() string {
	str := s.Name
	if fields, ok := structFields[s.Name]; ok {
		str += " { "
		for i, field := range fields {
			if i > 0 {
				str += "; "
			}
			str += typeSummary(field.Type) + " " + field.Name
		}
		str += " }"
	}
	return str
}

func (s *StructType) isTastType() {}

func (s *StructType) FieldInfo(name string) (*FieldInfo, bool) {
	fields, ok := structFields[s.Name]
	if !ok {
		return nil, false
	}
	for i, f := range fields {
		if f.Name == name {
			return &FieldInfo{Type: f.Type, Idx: i}, true
		}
	}
	return nil, false
}

func (s *StructType) Fields() []string {
	fields, ok := structFields[s.Name]
	if !ok {
		return nil
	}
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names
}

type TypedefType struct {
	Name    string
	Aliased Type
}

func Typedef(name string, aliased Type) *TypedefType {
	return &TypedefType{Name: name, Aliased: aliased}
}

func (t *TypedefType) String() string {
	return "typedef:" + t.Name + "(" + typeSummary(t.Aliased) + ")"
}

func (t *TypedefType) isTastType() {}

type PointerType struct {
	Elem Type
}

func (p *PointerType) String() string {
	return p.Elem.String() + "*"
}
func (p *PointerType) isTastType() {}

func Pointer(elem Type) *PointerType {
	return &PointerType{Elem: elem}
}

func typeSummary(typ Type) string {
	switch t := typ.(type) {
	case *StructType:
		return "struct " + t.Name
	case *PointerType:
		// only print one level that is pointer to struct or base type
		switch elem := t.Elem.(type) {
		case *StructType:
			return "struct " + elem.Name + "*"
		default:
			return typeSummary(elem) + "*"
		}
	case *ArrayType:
		return typeSummary(t.Elem) + "[]"
	case BaseType:
		return t.String()
	case *TypedefType:
		return t.Name
	default:
		return "unknown"
	}
}

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
var _ Type = (*PointerType)(nil)
