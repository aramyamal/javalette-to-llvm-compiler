package tast

type Node interface {
	Line() int
	Col() int
	Text() string
}

type TypedNode interface {
	Node
	Type() Type
}
