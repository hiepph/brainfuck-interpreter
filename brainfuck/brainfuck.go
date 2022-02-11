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

// Lex reads the source and collects only 8 acceppted tokens
func Lex(program io.Reader) (tokens []rune) {
	s := bufio.NewScanner(program)
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
// stores in an instruction array and
// traverses the array to evaluate the program.
func Interprete(program io.Reader, in io.Reader, out io.Writer) {
	tokens := Lex(program)

	tape := NewTape(in, out)

	instr := NewInstruction(tokens, tape)
	for ok := instr.Fetch(); ok; ok = instr.Fetch() {
	}
}
