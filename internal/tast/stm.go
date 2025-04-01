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
	items []Item,
	line int,
	col int,
	text string,
) *DeclsStm {
	return &DeclsStm{
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

// VoidReturnStm is a return statement node in the typed abstract syntax tree
// that returns void
type VoidReturnStm struct {
	BaseNode
}

func (*VoidReturnStm) stmNode() {}

func NewVoidReturnStm(
	line int,
	col int,
	text string,
) *VoidReturnStm {
	return &VoidReturnStm{
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that ReturnStm implements Stm
var _ Stm = (*VoidReturnStm)(nil)

// WhileStm is while statement node in the typed abstract syntax tree
type WhileStm struct {
	Exp Exp
	Stm Stm

	BaseNode
}

func (*WhileStm) stmNode() {}

func NewWhileStm(
	exp Exp,
	stm Stm,
	line int,
	col int,
	text string,
) *WhileStm {
	return &WhileStm{
		Exp:      exp,
		Stm:      stm,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that WhileStm implements Stm
var _ Stm = (*WhileStm)(nil)

// BlockStm is block statement node in the typed abstract syntax tree containing
// a list of other statement nodes
type BlockStm struct {
	Stms []Stm

	BaseNode
}

func (*BlockStm) stmNode() {}

func NewBlockStm(
	stms []Stm,
	line int,
	col int,
	text string,
) *BlockStm {
	return &BlockStm{
		Stms:     stms,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that BlockStm implements Stm
var _ Stm = (*BlockStm)(nil)

// IfStm is an if statement node in the typed abstract syntax tree containing
// a then statement and an else statement that is nil if there is no else branch
type IfStm struct {
	Exp     Exp
	ThenStm Stm
	ElseStm Stm

	BaseNode
}

func (*IfStm) stmNode() {}

func NewIfStm(
	exp Exp,
	thenStm Stm,
	elseStm Stm,
	line int,
	col int,
	text string,
) *IfStm {
	return &IfStm{
		Exp:      exp,
		ThenStm:  thenStm,
		ElseStm:  elseStm,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that IfStm implements Stm
var _ Stm = (*IfStm)(nil)

// BlankStm is an empty statement node in the typed abstract syntax tree
type BlankStm struct {
	BaseNode
}

func (*BlankStm) stmNode() {}

func NewBlankStm(
	line int,
	col int,
	text string,
) *BlankStm {
	return &BlankStm{BaseNode: BaseNode{line: line, col: col, text: text}}
}

// ensure that BlankStm implements Stm
var _ Stm = (*BlankStm)(nil)
