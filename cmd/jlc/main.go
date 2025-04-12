package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/codegen"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typechk"
)

// custom error listener
type errorListener struct {
	*antlr.DefaultErrorListener
}

func (e *errorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol any,
	line int,
	column int,
	msg string,
	err antlr.RecognitionException) {
	// print ERROR to stderr, and return status code 1
	fmt.Fprintln(os.Stderr, "ERROR")
	os.Exit(1)
}

func main() {
	outputFile := flag.String("o", "", "Output file (default: stdout)")
	flag.Parse()
	args := flag.Args()

	var err error
	var stream antlr.CharStream
	if len(args) > 0 {
		filepath := args[0]
		stream, err = antlr.NewFileStream(filepath)
		if err != nil {
			log.Fatalf("Error reading file %v: %v", filepath, err)
		}
	} else {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Error processing standard input:", err)
		}
		stream = antlr.NewInputStream(string(input))
	}

	errorList := &errorListener{}

	lexer := parser.NewJavaletteLexer(stream)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorList)

	tokens := antlr.NewCommonTokenStream(lexer, 0)
	parser := parser.NewJavaletteParser(tokens)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errorListener{})

	tree := parser.Prgm()

	typechk := typechk.NewTypeChecker()
	tast, err := typechk.Typecheck(tree)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR")
		log.Fatalln(err)
	}

	var writer io.Writer
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Error creating output file %v: %v", *outputFile, err)
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}

	codegen := codegen.NewCodeGenerator(writer)
	if err := codegen.GenerateCode(tast); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR")
		log.Fatalln(err)
	}

	fmt.Fprintln(os.Stderr, "OK")
}
