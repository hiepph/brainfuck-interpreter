package brainfuck

import (
	"bufio"
	"io"
)

type Tape struct {
	data []int
	ptr  int
}

func NewTape() *Tape {
	return &Tape{make([]int, 30000, 30000), 0}
}

func Interprete(in io.Reader) {
	tape := NewTape()

	s := bufio.NewScanner(in)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		Command(tape, rune(s.Bytes()[0]))
	}
}

func Command(tape *Tape, op rune) {
	switch op {
	case '+':
		tape.data[tape.ptr]++
	case '-':
		tape.data[tape.ptr]--
	case '>':
		tape.ptr++
	case '<':
		tape.ptr--
	default:
		return
	}
}
