package tast

// Prgm represents the root node of a Javalette program in the TAST.
type Prgm struct {
	Defs []Def // Definitions
}

func (*Prgm) Line() int    { return 0 }
func (*Prgm) Col() int     { return 0 }
func (*Prgm) Text() string { return "" }

// NewPrgm creates a new Prgm node with the given definitions.
func NewPrgm(
	defs []Def,
) *Prgm {
	return &Prgm{Defs: defs}
}
