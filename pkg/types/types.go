package types

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

type Op int

const (
	OpInc Op = iota
	OpDec
	OpMul
	OpDiv
	OpMod
	OpAdd
	OpSub
	OpLt
	OpGt
	OpLe
	OpGe
	OpEq
	OpNe
)

func (op Op) String() string {
	return [...]string{
		"++", // Inc
		"--", // Dec
		"*",  // Mul
		"/",  // Div
		"%",  // Mod
		"+",  // Add
		"-",  // Sub
		"<",  // Lt
		">",  // Gt
		"<=", // Le
		">=", // Ge
		"==", // Eq
		"!=", // Ne
	}[op]
}

func (op Op) Name() string {
	return [...]string{
		"OpInc",
		"OpDec",
		"OpMul",
		"OpDiv",
		"OpMod",
		"OpAdd",
		"OpSub",
		"OpLt",
		"OpGt",
		"OpLe",
		"OpGe",
		"OpEq",
		"OpNe",
	}[op]
}
