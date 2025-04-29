package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aramyamal/javalette-to-llvm-compiler/gen/parser"
	"github.com/aramyamal/javalette-to-llvm-compiler/internal/typechk"
)

func main() {
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

	lexer := parser.NewJavaletteLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, 0)
	parser := parser.NewJavaletteParser(tokens)
	tree := parser.Prgm()

	typechk := typechk.NewTypeChecker()
	_, err = typechk.Typecheck(tree)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR")
		log.Fatalln(err)
	}
}
