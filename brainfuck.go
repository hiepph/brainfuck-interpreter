package main

import (
	"bufio"
	"fmt"
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
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanRunes)

	for s.Scan() {
		fmt.Println(s.Text())
	}
}

func command(tape *Tape, op rune) {
	tape.data[tape.ptr]++
}
