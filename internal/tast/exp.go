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

type IntToDoubleExp struct {
	Exp Exp

	BaseTypedNode
}

func (*IntToDoubleExp) expNode() {}

func NewIntToDoubleExp(
	e Exp,
) *IntToDoubleExp {
	return &IntToDoubleExp{
		Exp: e,
		BaseTypedNode: BaseTypedNode{
			typ:      Double,
			BaseNode: BaseNode{line: e.Line(), col: e.Col(), text: e.Text()},
		},
	}
}

// check that IntToDoubleExp implements Exp
var _ Exp = (*IntToDoubleExp)(nil)
