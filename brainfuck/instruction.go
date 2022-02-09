package brainfuck

import "fmt"

type Instruction struct {
	ptr    int
	tokens []rune
	tape   *Tape
}

func NewInstruction(tokens []rune, tape *Tape) *Instruction {
	return &Instruction{ptr: -1, tokens: tokens, tape: tape}
}

func (instr *Instruction) Fetch() {
	instr.ptr++
	instr.Execute()
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
func (instr *Instruction) Execute() {
	op := instr.tokens[instr.ptr]
	tape := instr.tape

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
		fmt.Fprintf(tape.out, "%c", instr.tape.data[tape.ptr])
	default:
		return
	}
}
