package tast

type Type int

const (
	Int Type = iota
	Double
	Bool
	String
	Void
	Unknown
)

func (t Type) String() string {
	return [...]string{
		"Int",
		"Double",
		"Bool",
		"String",
		"Void",
		"Unknown",
	}[t]
}
