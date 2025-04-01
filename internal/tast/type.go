package tast

type Type int

const (
	Unknown Type = iota
	Int
	Double
	Bool
	String
	Void
)

func (t Type) String() string {
	return [...]string{
		"Unknown",
		"Int",
		"Double",
		"Bool",
		"String",
		"Void",
	}[t]
}
