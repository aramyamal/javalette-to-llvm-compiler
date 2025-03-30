package typedast

type Type int

const (
	Int Type = iota
	Double
	Bool
	String
	Void
)

func (t Type) String() string {
	return [...]string{
		"Int",
		"Double",
		"Bool",
		"String",
		"Void",
	}[t]
}
