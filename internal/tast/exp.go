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

type BoolExp struct {
	Value bool

	BaseNode
}

func (*BoolExp) expNode() {}

func (BoolExp) Type() Type { return Bool }

func NewBoolExp(
	value bool,
	line int,
	col int,
	text string,
) *BoolExp {
	return &BoolExp{
		Value:    value,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that BoolExp implements Exp
var _ Exp = (*BoolExp)(nil)

type IntExp struct {
	Value int

	BaseNode
}

func (*IntExp) expNode() {}

func (IntExp) Type() Type { return Int }

func NewIntExp(
	value int,
	line int,
	col int,
	text string,
) *IntExp {
	return &IntExp{
		Value:    value,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that IntExp implements Exp
var _ Exp = (*IntExp)(nil)

type DoubleExp struct {
	Value float64

	BaseNode
}

func (*DoubleExp) expNode() {}

func (DoubleExp) Type() Type { return Double }

func NewDoubleExp(
	value float64,
	line int,
	col int,
	text string,
) *DoubleExp {
	return &DoubleExp{
		Value:    value,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that DoubleExp implements Exp
var _ Exp = (*DoubleExp)(nil)

type IdentExp struct {
	Id string

	BaseTypedNode
}

func (*IdentExp) expNode() {}

func NewIdentExp(
	id string,
	typ Type,
	line int,
	col int,
	text string,
) *IdentExp {
	return &IdentExp{
		Id: id,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that IdentExp implements Exp
var _ Exp = (*IdentExp)(nil)

type FuncExp struct {
	Id   string
	Exps []Exp

	BaseTypedNode
}

func (*FuncExp) expNode() {}

func NewFuncExp(
	id string,
	exps []Exp,
	typ Type,
	line int,
	col int,
	text string,
) *FuncExp {
	return &FuncExp{
		Id:   id,
		Exps: exps,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that FuncExp implements Exp
var _ Exp = (*FuncExp)(nil)
