package typechecker

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typedast"
)

func extractPosData(pr antlr.ParserRuleContext) (int, int, string) {
	return pr.GetStart().GetLine(), pr.GetStart().GetColumn(), pr.GetText()
}

func toAstType(fromType parser.ITypeContext) (typedast.Type, error) {
	parserChild := fromType.GetChild(0)
	switch parserChild.(type) {
	case *parser.IntTypeContext:
		return typedast.Int, nil
	case *parser.DoubleTypeContext:
		return typedast.Double, nil
	case *parser.BoolTypeContext:
		return typedast.Bool, nil
	case *parser.StringTypeContext:
		return typedast.String, nil
	case *parser.VoidTypeContext:
		return typedast.Void, nil
	default:
		return typedast.Unknown, fmt.Errorf(
			"type '%T' not yet implemented at %d:%d near '%s'",
			parserChild,
			fromType.GetStart().GetLine(),
			fromType.GetStart().GetColumn(),
		)
	}
}

