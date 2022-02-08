package brainfuck

import (
	"bufio"
	"fmt"
	"io"
)

type Tape struct {
	data []int
	ptr  int
	out  io.Writer
}

func NewTape(out io.Writer) *Tape {
	return &Tape{make([]int, 30000, 30000), 0, out}
}

// Interprete reads the source and through 4 stages runs the program.
// 1. Lexer: splits the source into tokens
// 2. Parser: outlines rules of the tokens
// 3. AST: hierarchical representation of the tokens
// 4. Interpreter: traverse the tree and run the operations
func Interprete(in io.Reader, out io.Writer) {
	tape := NewTape(out)

	s := bufio.NewScanner(in)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		tape.Command(rune(s.Bytes()[0]))
	}
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
