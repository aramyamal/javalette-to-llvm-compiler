package tast

type Prgm struct {
	Node
	Defs []Def
}

func (*Prgm) Line() int    { return 0 }
func (*Prgm) Col() int     { return 0 }
func (*Prgm) Text() string { return "" }

func NewPrgm(
	defs []Def,
) *Prgm {
	return &Prgm{Defs: defs}
}
