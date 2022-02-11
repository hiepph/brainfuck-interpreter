package main

import (
	"os"

	"github.com/hiepph/brainfuck-interpreter/brainfuck"
)

func main() {
	file, _ := os.Open("examples/echo.bf")
	brainfuck.Interprete(file, os.Stdin, os.Stdout)
}
