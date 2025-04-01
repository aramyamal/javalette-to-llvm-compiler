package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
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
	var stream antlr.CharStream
	if len(os.Args) > 1 {
		filepath := os.Args[1]
		var err error
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
	pars := parser.NewJavaletteParser(tokens)
	pars.RemoveErrorListeners()
	pars.AddErrorListener(&errorListener{})

	tree := pars.Prgm()

	_, err := typechk.Typecheck(tree)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR")
		log.Fatalln(err)
	}

	// temporary, checking if parser works
	fmt.Fprintln(os.Stderr, "OK")
}
