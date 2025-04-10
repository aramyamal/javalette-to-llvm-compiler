package tast

import "github.com/aramyamal/javalette-to-llvm-compiler/pkg/ir"

type Arg interface {
	TypedNode
	argNode()
}

type ParamArg struct {
	Id string

	BaseTypedNode
}

func (*ParamArg) argNode() {}

func NewParamArg(
	typ ir.Type,
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
