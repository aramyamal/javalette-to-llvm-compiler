package typechk

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func extractPosData(pr antlr.ParserRuleContext) (int, int, string) {
	return pr.GetStart().GetLine(), pr.GetStart().GetColumn(), pr.GetText()
}

func (tc *TypeChecker) toTastBaseType(fromType parser.IBaseTypeContext) (tast.Type, error) {
	switch t := fromType.GetChild(0).(type) {
	case *parser.IntTypeContext:
		return tast.Int, nil
	case *parser.DoubleTypeContext:
		return tast.Double, nil
	case *parser.BoolTypeContext:
		return tast.Bool, nil
	case *parser.StringTypeContext:
		return tast.String, nil
	case *parser.VoidTypeContext:
		return tast.Void, nil
	case *parser.CustomTypeContext:
		name := t.Ident().GetText()
		var baseType tast.Type
		var found bool
		if baseType, found = tc.env.LookupTypedef(name); !found {
			if baseType, found = tc.env.LookupStruct(name); !found {
				return nil, fmt.Errorf(
					"type '%s' not defined at %d:%d", name,
					t.GetStart().GetLine(), t.GetStart().GetColumn(),
				)
			}
		}
		return baseType, nil

	default:
		return tast.Unknown, fmt.Errorf(
			"type '%T' not yet implemented at %d:%d near '%s'",
			fromType,
			fromType.GetStart().GetLine(),
			fromType.GetStart().GetColumn(),
			fromType.GetText(),
		)
	}
}

func (tc *TypeChecker) toTastType(fromType parser.ITypeContext) (tast.Type, error) {
	switch t := fromType.(type) {
	case *parser.PrimitiveTypeContext:
		baseType := t.BaseType()
		typ, err := tc.toTastBaseType(baseType)
		if err != nil {
			return nil, err
		}
		arraySuffixes := t.AllArraySuffix()
		for range arraySuffixes {
			typ = tast.Array(typ)
		}
		return typ, nil
	default:
		return nil, fmt.Errorf("unhandled type '%T'", fromType)
	}
}

// Checks if actual type can be converted to expected type. Returns true if the
// conversion is valid.
func isConvertible(expected, actual tast.Type) bool {

	// handle array type recursively
	expectedArr, expectedIsArr := expected.(*tast.ArrayType)
	actualArr, actualIsArr := actual.(*tast.ArrayType)
	if expectedIsArr && actualIsArr {
		return isConvertible(expectedArr.Elem, actualArr.Elem)
	}
	if expectedIsArr || actualIsArr {
		return false
	}

	// handle pointer types
	expectedPtr, expectedIsPtr := expected.(*tast.PointerType)
	actualPtr, actualIsPtr := actual.(*tast.PointerType)
	if expectedIsPtr && actualIsPtr {
		return isConvertible(expectedPtr.Elem, actualPtr.Elem)
	}
	if expectedIsPtr || actualIsPtr {
		return false
	}

	// handle typedefs, allow if both are typedefs with the same name
	expectedTypedef, expectedIsTypedef := expected.(*tast.TypedefType)
	actualTypedef, actualIsTypedef := actual.(*tast.TypedefType)
	if expectedIsTypedef {
		return isConvertible(expectedTypedef.Aliased, actual)
	}
	if actualIsTypedef {
		return isConvertible(expected, actualTypedef.Aliased)
	}

	// handle struct types only by name
	expectedStruct, expectedIsStruct := expected.(*tast.StructType)
	actualStruct, actualIsStruct := actual.(*tast.StructType)
	if expectedIsStruct && actualIsStruct {
		return expectedStruct.Name == actualStruct.Name
	}
	if expectedIsStruct || actualIsStruct {
		return false
	}

	// then handle base types
	switch expected {
	case tast.Double:
		//return actual == ir.Int || actual == ir.Double
		return actual == tast.Double
	case tast.Int:
		return actual == tast.Int
	case tast.Bool:
		return actual == tast.Bool
	case tast.Void:
		return actual == tast.Void
	case tast.String:
		return actual == tast.String
	default:
		return false
	}
}

// Determines the dominant type between two tast.for operations. For example,
// int + double = double
func dominantType(type1, type2 tast.Type) (tast.Type, error) {

	// handle typedefs by unwrapping and recursing this function
	if t1, ok := type1.(*tast.TypedefType); ok {
		return dominantType(t1.Aliased, type2)
	}
	if t2, ok := type2.(*tast.TypedefType); ok {
		return dominantType(type1, t2.Aliased)
	}

	// handle pointers by making sure both are pointers then recurse on element
	if p1, ok1 := type1.(*tast.PointerType); ok1 {
		if p2, ok2 := type2.(*tast.PointerType); ok2 {
			elemType, err := dominantType(p1.Elem, p2.Elem)
			if err != nil {
				return tast.Unknown, err
			}
			return tast.Pointer(elemType), nil
		}
		return tast.Unknown, fmt.Errorf(
			"cannot mix pointer and non-pointer types: %v, %v", type1, type2,
		)
	}
	if _, ok := type2.(*tast.PointerType); ok {
		// only type2 is pointer, type1 is not
		return tast.Unknown, fmt.Errorf(
			"cannot mix pointer and non-pointer types: %v, %v", type1, type2,
		)
	}

	// if both are struct types, check that names are the same
	if t1, ok1 := type1.(*tast.StructType); ok1 {
		if t2, ok2 := type2.(*tast.StructType); ok2 {
			if t1.Name == t2.Name {
				return type1, nil
			}
			return tast.Unknown, fmt.Errorf(
				"illegal implicit conversion between struct %v and %v",
				t1.Name, t2.Name,
			)
		}
	}

	// same tast.return the same type
	if type1 == type2 {
		switch type1 {
		case tast.Int, tast.Bool, tast.Void, tast.Double:
			return type1, nil
		}
	}

	return tast.Unknown, fmt.Errorf(
		"illegal implicit conversion between %v and %v",
		type1, type2,
	)
}

func extractIncDecOp(opCtx parser.IIncDecOpContext) (tast.Op, error) {
	switch opCtx.(type) {
	case *parser.IncContext:
		return tast.OpInc, nil
	case *parser.DecContext:
		return tast.OpDec, nil
	default:
		return 0, fmt.Errorf("unhandled inc/dec operator type %T", opCtx)
	}
}

func promoteExp(exp tast.Exp, typ tast.Type) tast.Exp {
	if exp.Type() == tast.Int && typ == tast.Double {
		return tast.NewIntToDoubleExp(exp)
	}
	return exp
}

