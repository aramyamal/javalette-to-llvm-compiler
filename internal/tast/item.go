package tast

// Item represents a single variable declaration in a declaration statement of
// the typed abstract syntax tree (TAST).
type Item interface {
	TypedNode
	itemNode()
}

// NoInitItem represents a variable declaration without an initializer in a
// declaration statement in the TAST.
type NoInitItem struct {
	Id string // Variable name

	BaseTypedNode // Embeds type and source location information
}

func (*NoInitItem) itemNode() {}

// NewNoInitItem creates a new NoInitItem with the given name, type, and source
// location.
func NewNoInitItem(
	id string,
	typ Type,
	line int,
	col int,
	text string,
) *NoInitItem {
	return &NoInitItem{
		Id: id,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that NoInitItem implements Item
var _ Item = (*NoInitItem)(nil)

// InitItem represents a variable declaration with an initializer expression in
// a declaration statement in the TAST.
type InitItem struct {
	Id  string // Variable name
	Exp Exp    // Initialization expression

	BaseTypedNode // Embeds type and source location information
}

func (*InitItem) itemNode() {}

// NewInitItem creates a new InitItem with the given name, initializer
// expression, type, and source location.
func NewInitItem(
	id string,
	exp Exp,
	typ Type,
	line int,
	col int,
	text string,
) *InitItem {
	return &InitItem{
		Id:  id,
		Exp: exp,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that InitItem implements Item
var _ Item = (*InitItem)(nil)
