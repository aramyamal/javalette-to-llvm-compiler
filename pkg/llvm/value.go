package llvm

import "strconv"

type Value interface {
	String() string
}

type Literal interface {
	Value
	literal()
}

type Global string

func (g Global) String() string {
	return "@" + string(g)
}

type Var string

func (v Var) String() string {
	return "%" + string(v)
}

type LitInt int

func (l LitInt) String() string {
	return strconv.Itoa(int(l))
}

func (LitInt) literal() {}

type LitDouble float64

func (l LitDouble) String() string {
	return strconv.FormatFloat(float64(l), 'f', -1, 64)
}

func (LitDouble) literal() {}

type LitBool bool

func (l LitBool) String() string {
	return strconv.FormatBool(bool(l))
}

func (LitBool) literal() {}

type LitString string

func (l LitString) String() string {
	panic("string literal not yet implemented")
}

// TODO global cstring handling
