package tast

type Def interface {
	TypedNode
	defNode()
}

type FuncDef struct {
	Id   string
	Args []Arg
	Stms []Stm

	BaseTypedNode
}

func (*FuncDef) defNode() {}

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
