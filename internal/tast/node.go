package tast

import types "github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"

type Node interface {
	Line() int
	Col() int
	Text() string
}

type BaseNode struct {
	line int
	col  int
	text string
}

func (n BaseNode) Line() int    { return n.line }
func (n BaseNode) Col() int     { return n.col }
func (n BaseNode) Text() string { return n.text }

// typed node of typed ast
type TypedNode interface {
	Node
	Type() types.Type
}

type BaseTypedNode struct {
	BaseNode
	typ types.Type
}

func (n BaseTypedNode) Type() types.Type { return n.typ }
