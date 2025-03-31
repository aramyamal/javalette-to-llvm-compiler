package tast

type Op int

const (
	OpInc Op = iota
	OpDec
	OpMul
	OpDiv
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
