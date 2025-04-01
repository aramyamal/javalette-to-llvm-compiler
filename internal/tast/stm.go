package tast

// Stm represents an untyped statement node in the typed abstract syntax tree
type Stm interface {
	Node
	stmNode()
}

// ExpStm is an expression statement node in the typed abstract syntax tree
type ExpStm struct {
	Exp Exp

	BaseNode
}

func (*ExpStm) stmNode() {}

func NewExpStm(
	exp Exp,
	line int,
	col int,
	text string,
) *ExpStm {
	return &ExpStm{
		Exp:      exp,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that ExpStm implements Stm
var _ Stm = (*ExpStm)(nil)

// DeclsStm is a declaration statement node in the typed abstract syntax tree
type DeclsStm struct {
	Type  Type
	Items []Item

	BaseNode
}

func (*DeclsStm) stmNode() {}

func NewDeclsStm(
	typ Type,
	items []Item,
	line int,
	col int,
	text string,
) *DeclsStm {
	return &DeclsStm{
		Type:     typ,
		Items:    items,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that DeclsStm implements Stm
var _ Stm = (*DeclsStm)(nil)

// ReturnStm is a return statement node in the typed abstract syntax tree
type ReturnStm struct {
	Exp Exp

	BaseNode
}

func (*ReturnStm) stmNode() {}

func NewReturnStm(
	exp Exp,
	line int,
	col int,
	text string,
) *ReturnStm {
	return &ReturnStm{
		Exp:      exp,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that ReturnStm implements Stm
var _ Stm = (*ReturnStm)(nil)
