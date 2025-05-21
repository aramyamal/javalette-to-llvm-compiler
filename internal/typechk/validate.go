package typechk

import (
	"fmt"

	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/tast"
)

func (tc *TypeChecker) validateMainFunc(defs []parser.IDefContext) error {
	var mainFunc *parser.FuncDefContext
	for _, def := range defs {
		if funcDef, ok := def.(*parser.FuncDefContext); ok {
			if funcDef.Ident().GetText() == "main" {
				mainFunc = funcDef
				break
			}
		}
	}

	if mainFunc == nil {
		return fmt.Errorf("program has no entrypoint 'main'")
	}

	if len(mainFunc.AllArg()) != 0 {
		return fmt.Errorf("entrypoint 'main' may not have input variables")
	}

	if typ, err := tc.toTastType(mainFunc.Type_()); err != nil {
		return err
	} else if typ != tast.Int {
		return fmt.Errorf("'main' entrypoint function does not have type int")
	}

	return nil
}

func (tc *TypeChecker) validateDefs(
	defs []parser.IDefContext,
) error {

	// first pass to register struct names
	for _, def := range defs {
		switch d := def.(type) {
		case *parser.StructDefContext:
			name := d.Ident().GetText()
			// register struct name with placeholder type
			if ok := tc.env.ExtendStruct(name, tast.RegisterStruct(name)); !ok {
				return fmt.Errorf(
					"redefinition of struct '%s' at %d:%d",
					name, d.GetStart().GetLine(), d.GetStart().GetColumn(),
				)
			}
		// report unhandled types
		case *parser.TypedefDefContext:
			continue
		case *parser.FuncDefContext:
			continue // handled in last pass
		default:
			return fmt.Errorf(
				"validateDefs: unhandled def type %T at %d:%d near '%s'",
				d, d.GetStart().GetLine(), d.GetStart().GetColumn(), d.GetText(),
			)
		}
	}

	// second pass to register typedefs aliasing structs and other types
	for _, def := range defs {
		if d, ok := def.(*parser.TypedefDefContext); ok {
			alias := d.Type_(1).GetText()
			aliasType, err := tc.toTastType(d.Type_(0))
			if err != nil {
				return err
			}
			if ok := tc.env.ExtendTypedef(alias, tast.Pointer(aliasType)); !ok {
				return fmt.Errorf(
					"redefinition of typedef '%s' at %d:%d",
					alias, d.GetStart().GetLine(), d.GetStart().GetColumn(),
				)
			}
		}
	}

	// third pass to register correct struct fields
	for _, def := range defs {
		if d, ok := def.(*parser.StructDefContext); ok {
			name := d.Ident().GetText()
			fieldNames := make(map[string]struct{})
			var fields []*tast.FieldCreator
			for _, structField := range d.AllStructField() {

				// check for duplicate field names
				fieldName := structField.Ident().GetText()
				if _, exists := fieldNames[fieldName]; exists {
					return fmt.Errorf(
						"duplicate field name '%s' in struct '%s'",
						fieldName, name,
					)
				}
				fieldNames[fieldName] = struct{}{}

				fieldType, err := tc.toTastType(structField.Type_())
				if err != nil {
					return fmt.Errorf(
						"error resolving type of field '%s' in struct '%s': %v",
						structField.Ident().GetText(), name, err,
					)
				}
				fields = append(fields, tast.Field(fieldType, structField.Ident().GetText()))
			}
			tast.RegisterStruct(name, fields...)
		}
	}

	// last pass to handle functions
	for _, def := range defs {
		switch d := def.(type) {
		case *parser.FuncDefContext:
			name := d.Ident().GetText()
			returnType, err := tc.toTastType(d.Type_())
			if err != nil {
				return err
			}

			paramNames, params, err := tc.extractParams(d.AllArg())
			if err != nil {
				return err
			}

			if ok := tc.env.ExtendFunc(
				name, paramNames, params, returnType,
			); !ok {
				return fmt.Errorf(
					"redefinition of function '%s' at %d:%d",
					name, d.GetStart().GetLine(), d.GetStart().GetColumn(),
				)
			}
		}
	}
	return nil
}
