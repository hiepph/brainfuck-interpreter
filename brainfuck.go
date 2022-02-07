package main

import (
	"bufio"
	"io"
)

type Tape struct {
	data []int
	ptr  int
}

func newTape() *Tape {
	return &Tape{make([]int, 30000, 30000), 0}
}

func interprete(in io.Reader) {
	tape := newTape()

	s := bufio.NewScanner(in)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		command(tape, rune(s.Bytes()[0]))
	}
}

func command(tape *Tape, op rune) {
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
