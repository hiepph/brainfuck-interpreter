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

func Interprete(in io.Reader, out io.Writer) {
	tape := NewTape(out)

	s := bufio.NewScanner(in)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		tape.Command(rune(s.Bytes()[0]))
	}
}

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
