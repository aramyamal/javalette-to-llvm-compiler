package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parsing"
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
	// stop program, print ERROR to stderr and return status code 1
	log.Fatalln("ERROR")
}

func main() {
	isStdinput := true
	var stream antlr.CharStream
	if len(os.Args) > 1 {
		filepath := os.Args[1]
		var err error
		stream, err = antlr.NewFileStream(filepath)
		if err != nil {
			log.Fatalf("Error reading file %v: %v", filepath, err)
		}
		isStdinput = false
	} else {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Error processing standard input:", err)
		}
		stream = antlr.NewInputStream(string(input))
	}

	lexer := parsing.NewJavaletteLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewJavaletteParser(tokens)

	// add custom error listener
	p.AddErrorListener(&errorListener{})

	tree := p.Program()

	// temporary, checking if parsing works
	fmt.Fprintln(os.Stderr, "OK")
	fmt.Println(isStdinput, tree.ToStringTree(nil, p))
}
