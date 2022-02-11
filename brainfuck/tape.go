package brainfuck

import "io"

type Tape struct {
	data []int8
	ptr  int
	out  io.Writer
	in   io.Reader
}

func NewTape(in io.Reader, out io.Writer) *Tape {
	return &Tape{
		data: make([]int8, 30000, 30000),
		ptr:  0,
		in:   in,
		out:  out,
	}
}
