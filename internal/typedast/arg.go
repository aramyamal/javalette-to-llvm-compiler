package typedast

type Arg interface {
	TypedNode
	argNode()
}

type ParamArg struct {
	typ Type
	Id  string

	line   int
	col int
	text   string
}

func (a ParamArg) Type() Type   { return a.typ }
func (a ParamArg) Line() int    { return a.line }
func (a ParamArg) Col() int  { return a.col }
func (a ParamArg) Text() string { return a.text }
func (*ParamArg) argNode()      {}

func NewParamArg(
	typ Type,
	id string,
	line int,
	col int,
	text string,
) ParamArg {
	return ParamArg{
		typ: typ,
		Id:  id,

		line:   line,
		col: col,
		text:   text,
	}
}

// check that ParamArg implements Arg
var _ Arg = (*ParamArg)(nil)
