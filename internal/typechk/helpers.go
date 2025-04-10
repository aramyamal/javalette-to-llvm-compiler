package typechk

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/ir"
)

func extractPosData(pr antlr.ParserRuleContext) (int, int, string) {
	return pr.GetStart().GetLine(), pr.GetStart().GetColumn(), pr.GetText()
}

func toIrType(fromType parser.ITypeContext) (ir.Type, error) {
	parserChild := fromType.GetChild(0)
	switch parserChild.(type) {
	case *parser.IntTypeContext:
		return ir.Int, nil
	case *parser.DoubleTypeContext:
		return ir.Double, nil
	case *parser.BoolTypeContext:
		return ir.Bool, nil
	case *parser.StringTypeContext:
		return ir.String, nil
	case *parser.VoidTypeContext:
		return ir.Void, nil
	default:
		return ir.Unknown, fmt.Errorf(
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
func isConvertible(expected, actual ir.Type) bool {
	switch expected {
	case ir.Double:
		//return actual == ir.Int || actual == ir.Double
		return actual == ir.Double
	case ir.Int:
		return actual == ir.Int
	case ir.Bool:
		return actual == ir.Bool
	case ir.Void:
		return actual == ir.Void
	case ir.String:
		return actual == ir.String
	default:
		return false
	}
}

// Determines the dominant type between two types for operations. For example,
// int + double = double
func dominantType(type1, type2 ir.Type) (ir.Type, error) {
	// if either type is Double, the result is Double

	if (type1 == ir.Double && type2 == ir.Int) ||
		(type1 == ir.Double && type2 == ir.Double) {
		return ir.Double, nil
	}

	// same types return the same type
	if type1 == type2 {
		switch type1 {
		case ir.Int, ir.Bool, ir.Void, ir.Double:
			return type1, nil
		}
	}

	return ir.Unknown, fmt.Errorf(
		"illegal implicit conversion between %v and %v",
		type1, type2,
	)
}
