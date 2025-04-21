// Package tast defines the typed abstract syntax tree (TAST) for the Javalette
// language. The TAST is produced by the type checker and is later used for
// efficient code generation.
//
// The package provides interfaces and types for the nodes of the TAST that
// are specialized into different language constructs, all annotated with
// type (if applicable) and source location information.
package tast

// Node represents a node in the TAST with source location information.
//
// Line returns the line number in the source code where the node appears.
// Col returns the column number in the source code where the node appears.
// Text returns the original source text corresponding to the node.
type Node interface {
	Line() int
	Col() int
	Text() string
}

// BaseNode provides a basic implementation of the Node interface, storing
// source location information.
type BaseNode struct {
	line int    // Line number of source code
	col  int    // Column number of source code
	text string // Text of source code
}

func (n *BaseNode) Line() int    { return n.line }
func (n *BaseNode) Col() int     { return n.col }
func (n *BaseNode) Text() string { return n.text }

// TypedNode represents a node in the TAST that is type annotated.
// It embeds Node and adds a method to retrieve the node's type.
//
// Type returns the type of the node.
type TypedNode interface {
	Node
	Type() Type
}

// BaseTypedNode provides a basic implementation of the TypedNode interface.
type BaseTypedNode struct {
	BaseNode // Embeds source location information

	typ Type // Annotated type of the node
}

func (n *BaseTypedNode) Type() Type { return n.typ }
