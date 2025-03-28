package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parsing"
)

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
		stdinput, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Error processing standard input:", err)
		}
		stream = antlr.NewInputStream(string(stdinput))
	}

	lexer := parsing.NewJavaletteLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewJavaletteParser(tokens)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Program()

	// temporary, checking if it works
	fmt.Println(isStdinput, tree.ToStringTree(nil, p))
}
