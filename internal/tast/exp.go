package tast

type Exp interface {
	TypedNode
	expNode()
}

type ParenExp struct {
	Exp Exp

	BaseTypedNode
}

func (*ParenExp) expNode() {}

func NewParenExp(
	exp Exp,
	typ Type,
	line int,
	col int,
	text string,
) *ParenExp {
	return &ParenExp{
		Exp: exp,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that ParenExp implements Exp
var _ Exp = (*ParenExp)(nil)
