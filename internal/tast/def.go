package tast

// Def represents a definition node in the TAST.
type Def interface {
	TypedNode
	defNode()
}

// FuncDef represents a function definition in the TAST.
type FuncDef struct {
	Id   string // Function name
	Args []Arg  // Function arguments
	Stms []Stm  // Function body statements

	BaseTypedNode // Embeds type and source location information
}

func (*FuncDef) defNode() {}

// NewFuncDef creates a new FuncDef node with the given name, arguments, body
// statements, type, and source location information.
func NewFuncDef(
	id string,
	args []Arg,
	stms []Stm,
	typ Type,
	line int,
	col int,
	text string,
) *FuncDef {
	return &FuncDef{
		Id:   id,
		Args: args,
		Stms: stms,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that FuncDef implements Def
var _ Def = (*FuncDef)(nil)

// StructDef represents a struct definition in the TAST.
type StructDef struct {
	BaseTypedNode // Embed type and source location information
}

func (*StructDef) defNode() {}

// NewStructDef creates a new StructDef node with the given
// struct type and source location information.
func NewStructDef(
	structType *StructType,
	line,
	col int,
	text string,
) *StructDef {
	return &StructDef{
		BaseTypedNode: BaseTypedNode{
			typ:      structType,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that StructDef implements Def
var _ Def = (*StructDef)(nil)

// TypedefDef represents a typedef definition in the TAST.
type TypedefDef struct {
	Id            string // Typedef alias name.
	BaseTypedNode        // Embed type and source location information
}

func (*TypedefDef) defNode() {}

// NewTypedefDef creates a new typedef node with the given
// typedef type and source location information.
func NewTypedefDef(
	id string,
	typ Type,
	line,
	col int,
	text string,
) *TypedefDef {
	return &TypedefDef{
		Id: id,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that TypeDef implements Def
var _ Def = (*TypedefDef)(nil)
