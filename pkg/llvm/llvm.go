package llvm

import (
	"fmt"
	"io"
	"strings"
)

type LLVMWriter struct {
	writer io.Writer
}

type Param struct {
	Type Type
	Name Reg
}

func NewParam(typ Type, name string) Param {
	return Param{Type: typ, Name: Reg(name)}
}

func NewLLVMWriter(w io.Writer) *LLVMWriter {
	return &LLVMWriter{writer: w}
}

func (w *LLVMWriter) Newline() error {
	_, err := w.writer.Write([]byte("\n"))
	return err
}

func (w *LLVMWriter) Declare(
	returns Type,
	funcName Global,
	inputs ...Type,
) error {
	var inputTypes []string
	for _, input := range inputs {
		inputTypes = append(inputTypes, input.String())
	}
	llvmInputs := strings.Join(inputTypes, ", ")

	llvmInstr := fmt.Sprintf(
		"declare %s %s(%s)\n",
		returns.String(),
		funcName.String(),
		llvmInputs,
	)
	if _, err := w.writer.Write([]byte(llvmInstr)); err != nil {
		return err
	}
	return nil
}

func (w *LLVMWriter) StartDefine(
	returns Type,
	funcName Global,
	inputs ...Param,
) error {
	var llvmParams []string
	for _, param := range inputs {
		llvmParams = append(
			llvmParams,
			param.Type.String()+" "+param.Name.String(),
		)
	}
	llvmInstr := fmt.Sprintf(
		"define %s %s(%s){\n",
		returns.String(), funcName.String(), strings.Join(llvmParams, ", "),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) EndDefine() error {
	_, err := w.writer.Write([]byte("}\n"))
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

func (w *LLVMWriter) Ret(typ Type, val Value) error {
	llvmInstr := fmt.Sprintf("\tret %s %s\n", typ.String(), val.String())
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Constant(des Reg, typ Type, lit Value) error {
	llvmInstr := fmt.Sprintf(
		"\t%s = %s %s\n",
		des.String(), typ.String(), lit.String(),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) GetElementPtr(
	des Reg,
	typ Type,
	from Global,
	idx ...int,
) error {
	var indices []string
	for _, i := range idx {
		indices = append(indices, fmt.Sprintf("i32 %d", i))
	}
	t := typ.String()
	llvmInstr := fmt.Sprintf(
		"\t%s = getelementptr %s, %s* %s, %s\n",
		des.String(), t, t, from.String(), strings.Join(indices, ", "),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) InternalConstant(name Global, typ Type, val Value) error {
	llvmInstr := fmt.Sprintf(
		"%s = internal constant %s %s\n",
		name.String(), typ.String(), val.String(),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Alloca(des Reg, typ Type) error {
	llvmInstr := fmt.Sprintf("\t%s = alloca %s\n", des.String(), typ.String())
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Store(typ Type, value Value, ptr Reg) error {
	llvmType := typ.String()
	llvmInstr := fmt.Sprintf(
		"\tstore %s %s, %s* %s\n",
		llvmType, value.String(), llvmType, ptr.String(),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Load(des Reg, typ Type, ptr Reg) error {
	llvmType := typ.String()
	llvmInstr := fmt.Sprintf(
		"\t%s = load %s, %s* %s, align %d\n",
		des.String(), llvmType, llvmType, ptr.String(), typ.alignment(),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Add(typ Type, des Reg, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"%s = add i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"%s = fadd double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'add'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}
