package tast

type Item interface {
	Node
	itemNode()
}

type NoInitItem struct {
	Id string
	BaseNode
}

func (*NoInitItem) itemNode() {}

func NewNoInitItem(
	id string,
	line int,
	col int,
	text string,
) *NoInitItem {
	return &NoInitItem{
		Id:       id,
		BaseNode: BaseNode{line: line, col: col, text: text},
	}
}

// check that NoInitItem implements Item
var _ Item = (*NoInitItem)(nil)
