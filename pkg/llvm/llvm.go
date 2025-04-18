package llvm

import (
	"fmt"
	"io"
	"strings"
)

type LLVMWriter struct {
	writer io.Writer
}

type FuncParam struct {
	Type Type
	Name Reg
}

func Param(typ Type, name string) FuncParam {
	return FuncParam{Type: typ, Name: Reg(name)}
}

type FuncArg struct {
	Type  Type
	Value Value
}

func Arg(typ Type, value Value) FuncArg {
	return FuncArg{Type: typ, Value: value}
}

type PhiPair struct {
	Val   Value
	Label string
}

func Phi(val Value, lab string) PhiPair {
	return PhiPair{Val: val, Label: lab}
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
	inputs ...FuncParam,
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
	llvmInstr := fmt.Sprintf("\n%s:\n", name)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Label(name string) error {
	return w.Block(name)
}

func (w *LLVMWriter) Br(label string) error {
	llvmInstr := fmt.Sprintf("\tbr label %%%s\n", label)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) BrIf(
	typ Type,
	cond Value,
	iftrue string,
	iffalse string,
) error {
	if typ != I1 {
		return fmt.Errorf("Br: cannot branch on non-boolean values")
	}
	llvmInstr := fmt.Sprintf(
		"\tbr i1 %s, label %%%s, label %%%s\n", cond.String(), iftrue, iffalse,
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Phi(
	des Reg,
	typ Type,
	phiPairs ...PhiPair,
) error {
	if len(phiPairs) == 0 {
		return fmt.Errorf("Phi: must have at least one incoming value")
	}
	var froms []string
	for _, phiPair := range phiPairs {
		froms = append(froms,
			fmt.Sprintf("[ %s, %%%s ]", phiPair.Val.String(), phiPair.Label),
		)
	}
	llvmInstr := fmt.Sprintf(
		"\t%s = phi %s %s\n",
		des.String(), typ.String(), strings.Join(froms, ", "),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Ret(typ Type, val ...Value) error {
	var llvmInstr string
	if typ == Void {
		llvmInstr = "\tret void\n"
	} else {
		if len(val) == 0 {
			return fmt.Errorf("Ret: non-void return type requires a value")
		}
		llvmInstr = fmt.Sprintf("\tret %s %s\n", typ.String(), val[0].String())
	}
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

// Call emits a function call. If typ is llvm.Void, des is ignored
func (w *LLVMWriter) Call(
	des Reg,
	typ Type,
	funcName Global,
	args ...FuncArg,
) error {
	var argsStrs []string
	for _, arg := range args {
		argsStrs = append(argsStrs, arg.Type.String()+" "+arg.Value.String())
	}
	fmtArgs := strings.Join(argsStrs, ", ")
	var llvmInstr string
	if typ == Void {
		llvmInstr = fmt.Sprintf(
			"\tcall void %s(%s)\n",
			funcName.String(), fmtArgs,
		)
	} else {
		llvmInstr = fmt.Sprintf(
			"\t%s = call %s %s(%s)\n",
			des.String(), typ.String(), funcName.String(), fmtArgs,
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Sub(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = sub i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fsub double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'sub'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Add(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = add i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fadd double %s, %s\n",
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

func (w *LLVMWriter) Mul(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = mul i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fmul double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'mul'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Div(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = sdiv i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fdiv double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'div'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Rem(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string

	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = srem i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = frem double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'div'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) Xor(des Reg, typ Type, lhs, rhs Value) error {
	llvmInstr := fmt.Sprintf(
		"\t%s = xor %s %s, %s\n",
		des.String(), typ.String(), lhs.String(), rhs.String(),
	)
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) CmpLt(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string
	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp slt i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fcmp olt double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'cmp lt'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) CmpLe(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string
	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp sle i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fcmp ole double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'cmp le'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) CmpGt(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string
	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp sgt i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fcmp ogt double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'cmp gt'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) CmpGe(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string
	switch typ {
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp sge i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fcmp oge double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'cmp ge'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) CmpEq(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string
	switch typ {
	case I1:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp eq i1 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp eq i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fcmp oeq double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'cmp eq'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}

func (w *LLVMWriter) CmpNe(des Reg, typ Type, lhs, rhs Value) error {
	var llvmInstr string
	switch typ {
	case I1:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp ne i1 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case I32:
		llvmInstr = fmt.Sprintf(
			"\t%s = icmp ne i32 %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	case Double:
		llvmInstr = fmt.Sprintf(
			"\t%s = fcmp one double %s, %s\n",
			des.String(), lhs.String(), rhs.String(),
		)
	default:
		return fmt.Errorf(
			"unsupported type '%s' for LLVM instruction 'cmp ne'",
			typ.String(),
		)
	}
	_, err := w.writer.Write([]byte(llvmInstr))
	return err
}
