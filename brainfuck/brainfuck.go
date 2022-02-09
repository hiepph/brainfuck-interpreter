package brainfuck

import (
	"bufio"
	"fmt"
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

type Tape struct {
	data []int
	ptr  int
	out  io.Writer
}

type Instruction struct {
	ptr    int
	tokens []rune
	tape   *Tape
}

func NewTape(out io.Writer) *Tape {
	return &Tape{data: make([]int, 30000, 30000), ptr: 0, out: out}
}

func NewInstruction(tokens []rune, tape *Tape) *Instruction {
	return &Instruction{ptr: 0, tokens: tokens, tape: tape}
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
// stores in an instruction array and traverses the array to runs the program
func Interprete(in io.Reader, out io.Writer) {
	tokens := Lex(in)

	tape := NewTape(out)

	instr := NewInstruction(tokens, tape)
	instr.Next()
}

// Command modifies the memory tape or reads from input or writes to output
// depending on the operator.
//
// >: increments the data pointer (points to the next cell on the right)
// <: decrements the data pointer (points to the previous cell on the left)
// +: increases by one the byte at the data pointer
// -: decreases by one the byte at the data pointer
// .: output the byte at the data pointer, using the ASCII character encoding
// [: if the byte at the data pointer is zero, jump forward to the command after
//    ']'; move forward to the next command otherwise.
// ]: if the byte at the data pointer is nonzero, jump back to the command
//    after matching '['; move forward to the next command otherwise.
func (tape *Tape) Command(op rune) {
	switch op {
	case '+':
		tape.data[tape.ptr]++
	case '-':
		tape.data[tape.ptr]--
	case '>':
		tape.ptr++
	case '<':
		tape.ptr--
	case '.':
		fmt.Fprintf(tape.out, "%c", tape.data[tape.ptr])
	default:
		return
	}
}

func (instr *Instruction) Next() {
}
