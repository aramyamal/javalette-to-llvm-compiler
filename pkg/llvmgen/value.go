package llvmgen

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	String() string
}

type Global string

func (g Global) String() string {
	return "@" + string(g)
}

type Reg string

func (v Reg) String() string {
	return "%" + string(v)
}

type LitInt int

func (l LitInt) String() string {
	return strconv.Itoa(int(l))
}

type LitDouble float64

func (l LitDouble) String() string {
	llvmString := strconv.FormatFloat(float64(l), 'f', -1, 64)
	if !strings.Contains(llvmString, ".") {
		llvmString += ".0"
	}
	return llvmString
}

type LitBool bool

func (l LitBool) String() string {
	return strconv.FormatBool(bool(l))
}

type LitString string

func (l LitString) String() string {
	return `c"` + string(l) + `\00"`
}

type StructValue struct {
	typ    *StructType
	fields []Value
}

func Struct(typ *StructType, fields ...Value) *StructValue {
	if len(fields) != len(typ.Fields) {
		panic(fmt.Sprintf(
			"StructValue: field count mismatch, got %d but expected want %d)",
			len(fields), len(typ.Fields),
		))
	}
	return &StructValue{typ: typ, fields: fields}
}

func (s StructValue) String() string {
	var parts []string
	for i, v := range s.fields {
		parts = append(parts, s.typ.Fields[i].String()+" "+v.String())
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

type NullValue struct{}

func (n NullValue) String() string {
	return "null"
}

func Null() NullValue {
	return NullValue{}
}
