package brainfuck

import (
	"bufio"
	"io"
)

var ACCEPTED_TOKENS = map[rune]bool{
	'+': true,
	'-': true,
	'[': true,
	']': true,
	'>': true,
	'<': true,
	'.': true,
	',': true,
}

func Lex(in io.Reader) (tokens []rune) {
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		ch := rune(s.Bytes()[0])
		if _, found := ACCEPTED_TOKENS[ch]; found {
			tokens = append(tokens, ch)
		}
	}

	return
}

// Interprete reads the source, breaks down into operation tokens,
// stores in an instruction array and traverses the array to evaluate the program
func Interprete(in io.Reader, out io.Writer) {
	tokens := Lex(in)

	tape := NewTape(out)

	instr := NewInstruction(tokens, tape)
	instr.Next()
}
