package typechecker

import (
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parsing"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typedast"
)

func Typecheck(tree parsing.IPrgmContext) typedast.Prgm {
	tp := typedast.NewPrgm(make([]typedast.Def, 0))
	// env := NewEnvironment[typedast.Type]()
	// for _, def := range tree.AllDef() {
	// 	fmt.Printf("%#v\n", def)
	// }
	return tp
}
