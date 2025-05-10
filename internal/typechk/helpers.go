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

func toTastBaseType(fromType parser.IBaseTypeContext) (tast.Type, error) {
	switch fromType.GetChild(0).(type) {
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

func toTastType(fromType parser.ITypeContext) (tast.Type, error) {
	baseType := fromType.BaseType()
	t, err := toTastBaseType(baseType)
	if err != nil {
		return nil, err
	}
	arraySuffixes := fromType.AllArraySuffix()
	for range arraySuffixes {
		t = tast.Array(t)
	}
	return t, nil
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
	// if either type is Double, the result is Double

	if (type1 == tast.Double && type2 == tast.Int) ||
		(type1 == tast.Double && type2 == tast.Double) {
		return tast.Double, nil
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
