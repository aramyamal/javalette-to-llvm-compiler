package llvm

import (
	"fmt"
	"io"
	"strings"

	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/types"
)

type LLVMWriter struct {
	writer io.Writer
}

type Param struct {
	Type types.Type
	Name Var
}

func NewParam(typ types.Type, name string) Param {
	return Param{Type: typ, Name: Var(name)}
}

func NewLLVMWriter(w io.Writer) *LLVMWriter {
	return &LLVMWriter{writer: w}
}

func (w *LLVMWriter) Newline() error {
	_, err := w.writer.Write([]byte("\n"))
	return err
}

func (w *LLVMWriter) Declare(
	returns types.Type,
	funcName Global,
	inputs ...types.Type,
) error {
	var inputTypes []string
	for _, input := range inputs {
		inputTypes = append(inputTypes, toLlvm(input))
	}
	llvmInputs := strings.Join(inputTypes, ", ")

	llvmInstr := fmt.Sprintf(
		"declare %s %s(%s)\n",
		toLlvm(returns),
		funcName.String(),
		llvmInputs,
	)
	if _, err := w.writer.Write([]byte(llvmInstr)); err != nil {
		return err
	}
	return nil
}

func (w *LLVMWriter) StartDefine(
	returns types.Type,
	funcName Global,
	inputs ...Param,
) error {
	var llvmParams []string
	for _, param := range inputs {
		llvmParams = append(llvmParams, toLlvm(param.Type)+param.Name.String())
	}
	llvmInstr := fmt.Sprintf(
		"define %s %s(%s){\n",
		toLlvm(returns), Global(funcName), strings.Join(llvmParams, ", "),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) EndDefine() error {
	_, err := w.writer.Write([]byte("}\n\n"))
	return err
}

func (w *LLVMWriter) Block(name string) error {
	llvmInstr := fmt.Sprintf("%s:\n", name)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Label(name string) error {
	return w.Block(name)
}

func (w *LLVMWriter) Ret(typ types.Type, val Value) error {
	llvmInstr := fmt.Sprintf("\tret %s %s\n", toLlvm(typ), val.String())
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Constant(des Var, typ types.Type, lit Literal) error {
	llvmInstr := fmt.Sprintf(
		"\t%s = %s %s\n",
		des.String(), toLlvm(typ), lit.String(),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err

}

func (w *LLVMWriter) Add(typ types.Type, des Var, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case types.Int:
		llvmInstr = fmt.Sprintf(
			"%s = add i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case types.Double:
		llvmInstr = fmt.Sprintf(
			"%s = fadd double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupperted type '%s' for LLVM instruction 'add'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func toLlvm(typ types.Type) string {
	switch typ {
	case types.Int:
		return "i32"
	case types.Bool:
		return "i1"
	case types.Double:
		return "double"
	case types.String:
		return "i8*"
	case types.Void:
		return "void"
	default:
		panic(fmt.Sprintf(
			"Conversion of type %s to LLVM not supported",
			typ.String(),
		))
	}
}
