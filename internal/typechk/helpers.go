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

func toAstType(fromType parser.ITypeContext) (tast.Type, error) {
	parserChild := fromType.GetChild(0)
	switch parserChild.(type) {
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
			parserChild,
			fromType.GetStart().GetLine(),
			fromType.GetStart().GetColumn(),
		)
	}
}
