package main

import (
	"strings"

	"github.com/hiepph/brainfuck-interpreter/brainfuck"
)

func main() {
	in := "++>+++++"
	brainfuck.Interprete(strings.NewReader(in))
}
