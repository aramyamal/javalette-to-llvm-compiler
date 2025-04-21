package tast

// Stm represents an untyped statement node in the TAST.
type Stm interface {
	Node
	stmNode()
}

// ExpStm represents an expression statement node in the TAST.
type ExpStm struct {
	Exp Exp // The expression

	BaseNode // Embeds source location information
}

func (*ExpStm) stmNode() {}

// NewExpStm creates a new ExpStm node with the given expression and source
// location.
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

// DeclsStm represents a variable declaration statement node in the TAST.
// Each variable declared in this statement is represented as an Item type.
type DeclsStm struct {
	Type  Type   // Declared type for all variables in this statement
	Items []Item // List of variable declarations

	BaseNode // Embeds source location information
}

func (*DeclsStm) stmNode() {}

// NewDeclsStm creates a new DeclsStm node with the given variable declarations
// and source location.
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

// ReturnStm represents a return statement node with a value in the TAST.
type ReturnStm struct {
	Type Type // Return type
	Exp  Exp  // Return expression

	BaseNode // Embeds source location information
}

func (*ReturnStm) stmNode() {}

// NewReturnStm creates a new ReturnStm node with the given type, expression,
// and source location.
func NewReturnStm(
	typ Type,
	exp Exp,
	line int,
	col int,
	text string,
) *ReturnStm {
	return &ReturnStm{
		Type:     typ,
		Exp:      exp,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// ensure that ReturnStm implements Stm
var _ Stm = (*ReturnStm)(nil)

// VoidReturnStm represents a return statement node with no value (void) in the
// TAST.
type VoidReturnStm struct {
	BaseNode // Embeds source location information
}

func (*VoidReturnStm) stmNode() {}

// NewVoidReturnStm creates a new VoidReturnStm node with the given source
// location.
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

// WhileStm represents a while statement node in the TAST.
type WhileStm struct {
	Exp Exp // Condition expression
	Stm Stm // Body statement

	BaseNode // Embeds source location information
}

func (*WhileStm) stmNode() {}

// NewWhileStm creates a new WhileStm node with the given condition, body, and
// source location.
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

// BlockStm is block statement node in the TAST containing a list of other
// statement nodes.
type BlockStm struct {
	Stms []Stm // List of statements

	BaseNode // Embeds source location information
}

func (*BlockStm) stmNode() {}

// NewBlockStm creates a new BlockStm node with the given statements and source
// location.
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

// IfStm represents an if statement node in the TAST, with an optional else
// branch.
type IfStm struct {
	Exp     Exp // Condition expression
	ThenStm Stm // Then branch statement
	ElseStm Stm // Else branch statement (nil if absent)

	BaseNode // Embeds source location information
}

func (*IfStm) stmNode() {}

// NewIfStm creates a new IfStm node with the given condition, branches, and
// source location.
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

// BlankStm represents an empty statement node in the TAST.
type BlankStm struct {
	BaseNode // Embeds source location information
}

func (*BlankStm) stmNode() {}

// NewBlankStm creates a new BlankStm node with the given source location.
func NewBlankStm(
	line int,
	col int,
	text string,
) *BlankStm {
	return &BlankStm{BaseNode: BaseNode{line: line, col: col, text: text}}
}

// ensure that BlankStm implements Stm
var _ Stm = (*BlankStm)(nil)
