package main

import (
	"fmt"
	"os"

	"github.com/hiepph/brainfuck-interpreter/brainfuck"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: bf [file]")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR: open file %q", filename)
		os.Exit(1)
	}
	brainfuck.Interprete(file, os.Stdin, os.Stdout)
}
