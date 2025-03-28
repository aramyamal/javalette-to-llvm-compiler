package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	code, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error processing standard input:", err)
	}
	fmt.Println("hey gurl you just gave me this through stdin: ", code)
}
