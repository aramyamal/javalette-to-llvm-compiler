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
