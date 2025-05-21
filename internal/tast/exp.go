package tast

// Exp represents an expression node in the typed abstract syntax tree (TAST).
type Exp interface {
	TypedNode
	HasSideEffect() bool
	IsLValue() bool
	expNode()
}

// ParenExp represents a parenthesis expression node in the TAST.
type ParenExp struct {
	Exp Exp // The inner expression.

	BaseTypedNode // Embeds type and source location information
}

func (*ParenExp) expNode()             {}
func (e ParenExp) HasSideEffect() bool { return e.Exp.HasSideEffect() }
func (e ParenExp) IsLValue() bool      { return e.Exp.IsLValue() }

// NewParenExp creates a new ParenExp node with the given inner expression,
// type, and source location.
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

// NullPtrExp represents a null pointer expression node in the TAST.
type NullPtrExp struct {
	BaseTypedNode // Embeds type and source location information
}

func (*NullPtrExp) expNode()             {}
func (e NullPtrExp) HasSideEffect() bool { return false }
func (e NullPtrExp) IsLValue() bool      { return false }

// NewNullPtrExp creates a new NullPtrExp node with the given type and source
// location.
func NewNullPtrExp(
	typ Type,
	line int,
	col int,
	text string,
) *NullPtrExp {
	return &NullPtrExp{
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that NullPtrExp implements Exp
var _ Exp = (*NullPtrExp)(nil)

// IntToDoubleExp represents an integer-to-double conversion expression node in
// the TAST.
type IntToDoubleExp struct {
	Exp Exp // The integer expression to convert

	BaseTypedNode // Embeds type and source location information
}

func (*IntToDoubleExp) expNode()           {}
func (IntToDoubleExp) HasSideEffect() bool { return false }
func (IntToDoubleExp) IsLValue() bool      { return false }

// NewIntToDoubleExp creates a new IntToDoubleExp node with the given expression.
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

// BoolExp represents a boolean literal expression node in the TAST.
type BoolExp struct {
	Value bool // The boolean value

	BaseNode // Embeds source location information
}

func (*BoolExp) expNode() {}

// Type returns the type of the expression (Bool).
func (BoolExp) Type() Type          { return Bool }
func (BoolExp) HasSideEffect() bool { return false }
func (BoolExp) IsLValue() bool      { return false }

// NewBoolExp creates a new BoolExp node with the given boolean and source
// location.
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

// IntExp represents an integer literal expression node in the TAST.
type IntExp struct {
	Value int // The integer value

	BaseNode // Embeds source location information
}

func (*IntExp) expNode() {}

// Type returns the type of the expression (Int).
func (IntExp) Type() Type          { return Int }
func (IntExp) HasSideEffect() bool { return false }
func (IntExp) IsLValue() bool      { return false }

// NewIntExp creates a new IntExp node with the given integer and source
// location.
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

// DoubleExp represents a double literal expression node in the TAST.
type DoubleExp struct {
	Value float64 // The double value

	BaseNode // Embeds source location information
}

func (*DoubleExp) expNode() {}

// Type returns the type of the expression (Double).
func (DoubleExp) Type() Type          { return Double }
func (DoubleExp) HasSideEffect() bool { return false }
func (DoubleExp) IsLValue() bool      { return false }

// NewDoubleExp creates a new DoubleExp node with the given value and source location.
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

// NewArrExp represents array allocation expression in the TAST.
type NewArrExp struct {
	Exps []Exp // Array index expressions

	BaseTypedNode // Embeds type and source location information
}

func (*NewArrExp) expNode()           {}
func (NewArrExp) HasSideEffect() bool { return true }
func (NewArrExp) IsLValue() bool      { return false }

// NewNewArrExp creates a new NewArrExp node in the TAST with dimension
// expressions, type, and source location.
func NewNewArrExp(
	exps []Exp,
	typ Type,
	line int,
	col int,
	text string,
) *NewArrExp {
	return &NewArrExp{
		Exps: exps,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that NewArrExp implements Exp
var _ Exp = (*NewArrExp)(nil)

// NewStructExp represents struct allocation expression in the TAST.
type NewStructExp struct {
	BaseTypedNode // Embeds type and source location information
}

func (*NewStructExp) expNode()           {}
func (NewStructExp) HasSideEffect() bool { return true }
func (NewStructExp) IsLValue() bool      { return false }

// NewNewStructExp creates a new NewStructExp node in the TAST with given
// struct type name, type, and source location.
func NewNewStructExp(
	typ Type,
	line int,
	col int,
	text string,
) *NewStructExp {
	return &NewStructExp{
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that NewStructExp implements Exp
var _ Exp = (*NewStructExp)(nil)

// IdentExp represents an identifier expression node in the TAST.
type IdentExp struct {
	Id string // Identifier name

	BaseTypedNode // Embeds type and source location information
}

func (*IdentExp) expNode()           {}
func (IdentExp) HasSideEffect() bool { return false }
func (IdentExp) IsLValue() bool      { return true }

// NewIdentExp creates a new IdentExp node with the given identifier, type, and source location.
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

// FuncExp represents a function call expression node in the TAST.
type FuncExp struct {
	Id   string // Function name
	Exps []Exp  // Function arguments

	BaseTypedNode // Embeds type and source location information
}

func (*FuncExp) expNode()           {}
func (FuncExp) HasSideEffect() bool { return true }
func (FuncExp) IsLValue() bool      { return false }

// NewFuncExp creates a new FuncExp node with the given function name,
// arguments, type, and source location.
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

// ArrIndexExp represents an array element access expression in the TAST.
type ArrIndexExp struct {
	Exp     Exp   // Array expression
	IdxExps []Exp // Index expressions for each dimension

	BaseTypedNode // Embeds type and source location information
}

func (*ArrIndexExp) expNode()           {}
func (ArrIndexExp) HasSideEffect() bool { return false }
func (ArrIndexExp) IsLValue() bool      { return true }

// NewArrIndexExp creates a new ArrIndexExp node with the given array expression,
// index expressions, result type, and source location.
func NewArrIndexExp(
	exp Exp,
	idxExps []Exp,
	typ Type,
	line int,
	col int,
	text string,
) *ArrIndexExp {
	return &ArrIndexExp{
		Exp:     exp,
		IdxExps: idxExps,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that ArrIndexExp implements Exp
var _ Exp = (*ArrIndexExp)(nil)

// FieldExp represents a field access expression in the TAST.
type FieldExp struct {
	Exp  Exp    // Expression whose field to access
	Name string // Name of the field to access

	BaseTypedNode // Embeds type and source location information
}

func (*FieldExp) expNode()           {}
func (FieldExp) HasSideEffect() bool { return false }
func (FieldExp) IsLValue() bool      { return true }

// NewFieldExp creates a new FieldExp node with the given expression
// whose field to access, the name of the field, type, and source location.
func NewFieldExp(
	exp Exp,
	name string,
	typ Type,
	line int,
	col int,
	text string,
) *FieldExp {
	return &FieldExp{
		Exp:  exp,
		Name: name,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that FieldExp implements Exp
var _ Exp = (*FieldExp)(nil)

// DerefExp represents a field pointer dereference expression in the TAST.
type DerefExp struct {
	Exp  Exp    // Expression whose field to dereference
	Name string // Name of the field to dereference

	BaseTypedNode // Embeds type and source location information
}

func (*DerefExp) expNode()           {}
func (DerefExp) HasSideEffect() bool { return false }
func (DerefExp) IsLValue() bool      { return true }

// NewDerefExp creates a new DerefExp node with the given expression
// whose field to dereference, the name of the field, type, and source location.
func NewDerefExp(
	exp Exp,
	name string,
	typ Type,
	line int,
	col int,
	text string,
) *DerefExp {
	return &DerefExp{
		Exp:  exp,
		Name: name,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that DerefExp implements Exp
var _ Exp = (*DerefExp)(nil)

// StringExp represents a string literal expression node in the TAST.
type StringExp struct {
	Value string // String value

	BaseNode // Embeds source location information
}

func (*StringExp) expNode()           {}
func (StringExp) Type() Type          { return String }
func (StringExp) HasSideEffect() bool { return false }
func (StringExp) IsLValue() bool      { return false }

// NewStringExp creates a new StringExp node with the given value and source
// location.
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

// NegExp represents a negation expression node in the TAST.
type NegExp struct {
	Exp Exp // Expression to negate

	BaseTypedNode // Embeds type and source location information
}

func (*NegExp) expNode()           {}
func (NegExp) HasSideEffect() bool { return false }
func (NegExp) IsLValue() bool      { return false }

// NewNegExp creates a new NegExp node with the given expression, type, and
// source location.
func NewNegExp(
	exp Exp,
	typ Type,
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

// NotExp represents a logical NOT expression node in the TAST.
type NotExp struct {
	Exp Exp // Expression to NOT

	BaseNode // Embeds source location information
}

func (*NotExp) expNode()           {}
func (*NotExp) Type() Type         { return Bool }
func (NotExp) HasSideEffect() bool { return false }
func (NotExp) IsLValue() bool      { return false }

// NewNotExp creates a new NotExp node with the given expression and source
// location.
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

// PostExp represents a post-increment or post-decrement expression node in the
// TAST.
type PostExp struct {
	Exp Exp // Expression being incremented or decremented
	Op  Op  // Operation (OpInc or OpDec)

	BaseTypedNode // Embeds type and source location information
}

func (*PostExp) expNode()           {}
func (PostExp) HasSideEffect() bool { return true }
func (PostExp) IsLValue() bool      { return false }

// NewPostExp creates a new PostExp node with the given identifier, operation,
// type, and source location.
func NewPostExp(
	exp Exp,
	op Op,
	typ Type,
	line int,
	col int,
	text string,
) *PostExp {
	return &PostExp{
		Exp: exp,
		Op:  op,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that PostExp implements Exp
var _ Exp = (*PostExp)(nil)

// PreExp represents a pre-increment or pre-decrement expression node in the
// TAST.
type PreExp struct {
	Exp Exp // Expression being incremented or decremented
	Op  Op  // Operation (OpInc or OpDec)

	BaseTypedNode // Embeds type and source location information
}

func (*PreExp) expNode()           {}
func (PreExp) HasSideEffect() bool { return true }
func (PreExp) IsLValue() bool      { return false }

// NewPreExp creates a new PreExp node with the given identifier, operation,
// type, and source location.
func NewPreExp(
	exp Exp,
	op Op,
	typ Type,
	line int,
	col int,
	text string,
) *PreExp {
	return &PreExp{
		Exp: exp,
		Op:  op,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that PreExp implements Exp
var _ Exp = (*PreExp)(nil)

// MulExp represents a multiplication, division, or modulo expression node in
// the TAST.
type MulExp struct {
	LeftExp  Exp // Left operand expression
	RightExp Exp // Right operand expression
	Op       Op  // Operation (OpMul, OpDiv, or OpMod)

	BaseTypedNode // Embeds type and source location information
}

func (*MulExp) expNode()           {}
func (MulExp) HasSideEffect() bool { return false }
func (MulExp) IsLValue() bool      { return false }

// NewMulExp creates a new MulExp node with the given operands, operation, type,
// and source location.
func NewMulExp(
	leftExp Exp,
	rightExp Exp,
	op Op,
	typ Type,
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

// AddExp represents an addition or subtraction expression node in the TAST.
type AddExp struct {
	LeftExp  Exp // Left operand expression
	RightExp Exp // Right operand expression
	Op       Op  // Operation (OpAdd or OpSub)

	BaseTypedNode // Embeds type and source location information
}

func (*AddExp) expNode()           {}
func (AddExp) HasSideEffect() bool { return false }
func (AddExp) IsLValue() bool      { return false }

// NewAddExp creates a new AddExp node with the given operands, operation, type,
// and source location.
func NewAddExp(
	leftExp Exp,
	rightExp Exp,
	op Op,
	typ Type,
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

// CmpExp represents a comparison expression node in the TAST.
type CmpExp struct {
	LeftExp  Exp // Left operand expression
	RightExp Exp // Right operand expression
	Op       Op  // Comparison operation (e.g., OpLt, OpGt, OpEq, etc.)

	BaseNode // Embeds source location information
}

func (*CmpExp) expNode()           {}
func (CmpExp) Type() Type          { return Bool }
func (CmpExp) HasSideEffect() bool { return false }
func (CmpExp) IsLValue() bool      { return false }

// NewCmpExp creates a new CmpExp node with the given operands, operation, and
// source location.
func NewCmpExp(
	leftExp Exp,
	rightExp Exp,
	op Op,
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

// AndExp represents a logical AND expression node in the TAST.
type AndExp struct {
	LeftExp  Exp // Left operand expression
	RightExp Exp // Right operand expression

	BaseNode // Embeds source location information
}

func (*AndExp) expNode()           {}
func (AndExp) Type() Type          { return Bool }
func (AndExp) HasSideEffect() bool { return false }
func (AndExp) IsLValue() bool      { return false }

// NewAndExp creates a new AndExp node with the given operands and source
// location.
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

// OrExp represents a logical OR expression node in the TAST.
type OrExp struct {
	LeftExp  Exp // Left operand expression
	RightExp Exp // Right operand expression

	BaseNode // Embeds source location information
}

func (*OrExp) expNode()           {}
func (OrExp) Type() Type          { return Bool }
func (OrExp) HasSideEffect() bool { return false }
func (OrExp) IsLValue() bool      { return false }

// NewOrExp creates a new OrExp node with the given operands and source location.
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

// AssignExp represents an assignment expression node in the TAST.
type AssignExp struct {
	ExpLhs Exp // Expression being assigned to
	Exp    Exp // Expression being assigned

	BaseTypedNode // Embeds type and source location information
}

func (*AssignExp) expNode()           {}
func (AssignExp) HasSideEffect() bool { return true }
func (AssignExp) IsLValue() bool      { return false }

// NewAssignExp creates a new AssignExp node with the given identifier,
// expression, type, and source location.
func NewAssignExp(
	expLhs Exp,
	exp Exp,
	typ Type,
	line int,
	col int,
	text string,
) *AssignExp {
	return &AssignExp{
		ExpLhs: expLhs,
		Exp:    exp,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that AssignExp implements Exp
var _ Exp = (*AssignExp)(nil)
