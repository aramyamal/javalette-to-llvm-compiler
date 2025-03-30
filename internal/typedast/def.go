package typedast

type Def interface {
	TypedNode
	defNode()
}

type FuncDef struct {
	Id   string
	Args []Arg
	Stms []Stm
	typ  Type

	line int
	col  int
	text string
}

func (f FuncDef) Type() Type   { return f.typ }
func (f FuncDef) Line() int    { return f.line }
func (f FuncDef) Col() int     { return f.col }
func (f FuncDef) Text() string { return f.text }
func (*FuncDef) defNode()      {}

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
		typ:  typ,
		Id:   id,
		Args: args,
		Stms: stms,

		line: line,
		col:  col,
		text: text,
	}
}

// check that FuncDef implements Def
var _ Def = (*FuncDef)(nil)
