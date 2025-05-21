package codegen

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
	"github.com/aramyamal/javalette-to-llvm-compiler/pkg/llvmgen"
)

func (cg *CodeGenerator) compileDef(def tast.Def) error {
	cg.env.EnterContext()
	switch d := def.(type) {
	case *tast.FuncDef:
		return cg.compileFuncDef(d)
	case *tast.StructDef:
		return cg.compileStructDef(d)
	case *tast.TypedefDef:
		return nil
	default:
		return fmt.Errorf(
			"compileDef: unhandled def type %T at %d:%d near '%s'",
			d, d.Line(), d.Col(), d.Text(),
		)
	}
}

func (cg *CodeGenerator) compileFuncDef(d *tast.FuncDef) error {
	params, err := cg.extractParams(d.Args)
	if err != nil {
		return err
	}
	cg.write.StartDefine(cg.toLlvmRetType(d.Type()), llvmgen.Global(d.Id), params...)
	cg.write.Label("entry")
	for _, param := range params {
		cg.emitVarAlloc(string(param.Name), param.Type, param.Name)
	}
	for _, stm := range d.Stms {
		if err := cg.compileStm(stm); err != nil {
			return err
		}
	}
	return cg.write.EndDefine()
}

func (cg *CodeGenerator) compileStructDef(d *tast.StructDef) error {
	structType, ok := d.Type().(*tast.StructType)
	if !ok {
		return fmt.Errorf(
			"internal compiler error in compileStructDef: "+
				"badly defined StructDef node, expected node type"+
				"tast.StructType but got %T at %d:%d near %s",
			d.Type(), d.Line(), d.Col(), d.Text(),
		)
	}

	fields := structType.Fields()
	fieldLlvmTypes := make([]llvmgen.Type, len(fields))
	for i, fieldName := range fields {
		fieldInfo, ok := structType.FieldInfo(fieldName)
		if !ok {
			return fmt.Errorf(
				"internal compiler error in compileStructDef: "+
					"unable to access struct field %s from %T",
				fieldName, structType,
			)
		}
		fieldLlvmTypes[i] = cg.toLlvmType(fieldInfo.Type)
	}
	return cg.emitTypeDecl(llvmgen.StructDef(
		structType.Name, fieldLlvmTypes...,
	))
}

func (cg *CodeGenerator) extractParams(args []tast.Arg) ([]llvmgen.FuncParam, error) {
	var params []llvmgen.FuncParam
	for _, arg := range args {
		switch a := arg.(type) {
		case *tast.ParamArg:
			params = append(params, llvmgen.Param(cg.toLlvmRetType(a.Type()), a.Id))
		default:
			return nil, fmt.Errorf(
				"extractParams: unhandled Arg type %T at %d:%d near '%s'",
				arg, arg.Line(), arg.Col(), arg.Text(),
			)
		}
	}
	return params, nil
}
