package tast

import "github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"

type Exp interface {
	TypedNode
	HasSideEffect() bool
	expNode()
}

type ParenExp struct {
	Exp Exp

	BaseTypedNode
}

func (*ParenExp) expNode()             {}
func (e ParenExp) HasSideEffect() bool { return e.Exp.HasSideEffect() }

func NewParenExp(
	exp Exp,
	typ types.Type,
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

func (*IntToDoubleExp) expNode()           {}
func (IntToDoubleExp) HasSideEffect() bool { return false }

func NewIntToDoubleExp(
	e Exp,
) *IntToDoubleExp {
	return &IntToDoubleExp{
		Exp: e,
		BaseTypedNode: BaseTypedNode{
			typ:      types.Double,
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

func (*BoolExp) expNode()           {}
func (BoolExp) Type() types.Type    { return types.Bool }
func (BoolExp) HasSideEffect() bool { return false }

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

func (*IntExp) expNode()           {}
func (IntExp) Type() types.Type    { return types.Int }
func (IntExp) HasSideEffect() bool { return false }

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

func (*DoubleExp) expNode()           {}
func (DoubleExp) Type() types.Type    { return types.Double }
func (DoubleExp) HasSideEffect() bool { return false }

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

func (*IdentExp) expNode()           {}
func (IdentExp) HasSideEffect() bool { return false }

func NewIdentExp(
	id string,
	typ types.Type,
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

func (*FuncExp) expNode()           {}
func (FuncExp) HasSideEffect() bool { return true }

func NewFuncExp(
	id string,
	exps []Exp,
	typ types.Type,
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

type StringExp struct {
	Value string

	BaseNode
}

func (*StringExp) expNode()           {}
func (StringExp) Type() types.Type    { return types.String }
func (StringExp) HasSideEffect() bool { return false }

func NewStringExp(
	value string,
	line int,
	col int,
	text string,
) *StringExp {
	return &StringExp{
		Value:    value,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that StringExp implements Exp
var _ Exp = (*StringExp)(nil)

type NegExp struct {
	Exp Exp

	BaseTypedNode
}

func (*NegExp) expNode()           {}
func (NegExp) HasSideEffect() bool { return false }

func NewNegExp(
	exp Exp,
	typ types.Type,
	line int,
	col int,
	text string,
) *NegExp {
	return &NegExp{
		Exp: exp,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that NegExp implements Exp
var _ Exp = (*NegExp)(nil)

type NotExp struct {
	Exp Exp

	BaseNode
}

func (*NotExp) expNode()           {}
func (*NotExp) Type() types.Type   { return types.Bool }
func (NotExp) HasSideEffect() bool { return false }

func NewNotExp(
	exp Exp,
	line int,
	col int,
	text string,
) *NotExp {
	return &NotExp{
		Exp:      exp,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that NotExp implements Exp
var _ Exp = (*NotExp)(nil)

type PostExp struct {
	Id string
	Op types.Op

	BaseTypedNode
}

func (*PostExp) expNode()           {}
func (PostExp) HasSideEffect() bool { return true }

func NewPostExp(
	id string,
	op types.Op,
	typ types.Type,
	line int,
	col int,
	text string,
) *PostExp {
	return &PostExp{
		Id: id,
		Op: op,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that PostExp implements Exp
var _ Exp = (*PostExp)(nil)

type PreExp struct {
	Id string
	Op types.Op

	BaseTypedNode
}

func (*PreExp) expNode()           {}
func (PreExp) HasSideEffect() bool { return true }

func NewPreExp(
	id string,
	op types.Op,
	typ types.Type,
	line int,
	col int,
	text string,
) *PreExp {
	return &PreExp{
		Id: id,
		Op: op,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that PreExp implements Exp
var _ Exp = (*PreExp)(nil)

type MulExp struct {
	LeftExp  Exp
	RightExp Exp
	Op       types.Op

	BaseTypedNode
}

func (*MulExp) expNode()           {}
func (MulExp) HasSideEffect() bool { return false }

func NewMulExp(
	leftExp Exp,
	rightExp Exp,
	op types.Op,
	typ types.Type,
	line int,
	col int,
	text string,
) *MulExp {
	return &MulExp{
		LeftExp:  leftExp,
		RightExp: rightExp,
		Op:       op,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that MulExp implements Exp
var _ Exp = (*MulExp)(nil)

type AddExp struct {
	LeftExp  Exp
	RightExp Exp
	Op       types.Op

	BaseTypedNode
}

func (*AddExp) expNode()           {}
func (AddExp) HasSideEffect() bool { return false }

func NewAddExp(
	leftExp Exp,
	rightExp Exp,
	op types.Op,
	typ types.Type,
	line int,
	col int,
	text string,
) *AddExp {
	return &AddExp{
		LeftExp:  leftExp,
		RightExp: rightExp,
		Op:       op,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that AddExp implements Exp
var _ Exp = (*AddExp)(nil)

type CmpExp struct {
	LeftExp  Exp
	RightExp Exp
	Op       types.Op

	BaseNode
}

func (*CmpExp) expNode()           {}
func (CmpExp) Type() types.Type    { return types.Bool }
func (CmpExp) HasSideEffect() bool { return false }

func NewCmpExp(
	leftExp Exp,
	rightExp Exp,
	op types.Op,
	line int,
	col int,
	text string,
) *CmpExp {
	return &CmpExp{
		LeftExp:  leftExp,
		RightExp: rightExp,
		Op:       op,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that CmpExp implements Exp
var _ Exp = (*CmpExp)(nil)

type AndExp struct {
	LeftExp  Exp
	RightExp Exp

	BaseNode
}

func (*AndExp) expNode()           {}
func (AndExp) Type() types.Type    { return types.Bool }
func (AndExp) HasSideEffect() bool { return false }

func NewAndExp(
	leftExp Exp,
	rightExp Exp,
	line int,
	col int,
	text string,
) *AndExp {
	return &AndExp{
		LeftExp:  leftExp,
		RightExp: rightExp,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that AndExp implements Exp
var _ Exp = (*AndExp)(nil)

type OrExp struct {
	LeftExp  Exp
	RightExp Exp

	BaseNode
}

func (*OrExp) expNode()           {}
func (OrExp) Type() types.Type    { return types.Bool }
func (OrExp) HasSideEffect() bool { return false }

func NewOrExp(
	leftExp Exp,
	rightExp Exp,
	line int,
	col int,
	text string,
) *OrExp {
	return &OrExp{
		LeftExp:  leftExp,
		RightExp: rightExp,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that OrExp implements Exp
var _ Exp = (*OrExp)(nil)

type AssignExp struct {
	Id  string
	Exp Exp

	BaseTypedNode
}

func (*AssignExp) expNode()           {}
func (AssignExp) HasSideEffect() bool { return true }

func NewAssignExp(
	id string,
	exp Exp,
	typ types.Type,
	line int,
	col int,
	text string,
) *AssignExp {
	return &AssignExp{
		Id:  id,
		Exp: exp,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that AssignExp implements Exp
var _ Exp = (*AssignExp)(nil)
