package tast

// Arg represents a function argument node in the TAST.
type Arg interface {
	TypedNode
	argNode()
}

// ParamArg represents a parameter argument in a function definition in the TAST.
type ParamArg struct {
	Id string // Parameter name

	BaseTypedNode // Embeds type and source location information
}

func (*ParamArg) argNode() {}

// NewParamArg creates a new ParamArg node with the given type, name, and source
// location information.
func NewParamArg(
	typ Type,
	id string,
	line int,
	col int,
	text string,
) *ParamArg {
	return &ParamArg{
		Id: id,
		BaseTypedNode: BaseTypedNode{
			typ:      typ,
			BaseNode: BaseNode{line: line, col: col, text: text},
		},
	}
}

// check that ParamArg implements Arg
var _ Arg = (*ParamArg)(nil)
