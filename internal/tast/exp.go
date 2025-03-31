package tast

type Exp interface {
	TypedNode
	expNode()
}

type BaseExp struct {
	typ  Type
	line int
	col  int
	text string
}

func (e BaseExp) Type() Type   { return e.typ }
func (e BaseExp) Line() int    { return e.line }
func (e BaseExp) Col() int     { return e.col }
func (e BaseExp) Text() string { return e.text }
func (*BaseExp) expNode()      {}

type ParenExp struct {
	Exp Exp
	BaseExp
}

func NewParenExp(
	exp Exp,
	typ Type,
	line int,
	col int,
	text string,
) *ParenExp {
	return &ParenExp{
		Exp: exp,
		BaseExp: BaseExp{
			typ:  typ,
			line: line,
			col:  col,
			text: text,
		},
	}
}

// check that ParenExp implements Exp
var _ Exp = (*ParenExp)(nil)
