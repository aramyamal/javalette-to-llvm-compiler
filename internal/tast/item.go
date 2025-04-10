package tast

import "github.com/aramyamal/javalette-to-llvm-compiler/pkg/ir"

type Item interface {
	TypedNode
	itemNode()
}

type NoInitItem struct {
	Id string

	BaseTypedNode
}

func (*NoInitItem) itemNode() {}

func NewNoInitItem(
	id string,
	typ ir.Type,
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

type InitItem struct {
	Id  string
	Exp Exp

	BaseTypedNode
}

func (*InitItem) itemNode() {}

func NewInitItem(
	id string,
	exp Exp,
	typ ir.Type,
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
