package typechk

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

func extractPosData(pr antlr.ParserRuleContext) (int, int, string) {
	return pr.GetStart().GetLine(), pr.GetStart().GetColumn(), pr.GetText()
}

func toIrType(fromType parser.ITypeContext) (types.Type, error) {
	parserChild := fromType.GetChild(0)
	switch parserChild.(type) {
	case *parser.IntTypeContext:
		return types.Int, nil
	case *parser.DoubleTypeContext:
		return types.Double, nil
	case *parser.BoolTypeContext:
		return types.Bool, nil
	case *parser.StringTypeContext:
		return types.String, nil
	case *parser.VoidTypeContext:
		return types.Void, nil
	default:
		return types.Unknown, fmt.Errorf(
			"type '%T' not yet implemented at %d:%d near '%s'",
			parserChild,
			fromType.GetStart().GetLine(),
			fromType.GetStart().GetColumn(),
			fromType.GetText(),
		)
	}
}

// Checks if actual type can be converted to expected type. Returns true if the
// conversion is valid.
func isConvertible(expected, actual types.Type) bool {
	switch expected {
	case types.Double:
		//return actual == ir.Int || actual == ir.Double
		return actual == types.Double
	case types.Int:
		return actual == types.Int
	case types.Bool:
		return actual == types.Bool
	case types.Void:
		return actual == types.Void
	case types.String:
		return actual == types.String
	default:
		return false
	}
}

// Determines the dominant type between two types for operations. For example,
// int + double = double
func dominantType(type1, type2 types.Type) (types.Type, error) {
	// if either type is Double, the result is Double

	if (type1 == types.Double && type2 == types.Int) ||
		(type1 == types.Double && type2 == types.Double) {
		return types.Double, nil
	}

	// same types return the same type
	if type1 == type2 {
		switch type1 {
		case types.Int, types.Bool, types.Void, types.Double:
			return type1, nil
		}
	}

	return types.Unknown, fmt.Errorf(
		"illegal implicit conversion between %v and %v",
		type1, type2,
	)
}
