package tast

// Type represents a Javalette type in the typed abstract syntax tree (TAST).
type Type int

const (
	Unknown Type = iota
	Int
	Double
	Bool
	String
	Void
)

// String returns the string representation of the Type.
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

// Op represents an operator in the TAST.
type Op int

const (
	OpInc Op = iota // ++ increment
	OpDec           // -- decrement
	OpMul           // * multiplication
	OpDiv           // / division
	OpMod           // % modulo
	OpAdd           // + addition
	OpSub           // - subtraction
	OpLt            // < less than
	OpGt            // > greater than
	OpLe            // <= less than or equal
	OpGe            // >= greater than or equal
	OpEq            // == equal
	OpNe            // != not equal
)

// String returns the symbol of the operator.
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

// Name returns the type name of the operator.
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
