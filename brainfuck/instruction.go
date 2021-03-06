package brainfuck

import "fmt"

type Instruction struct {
	ptr    int
	tokens []rune
	tape   *Tape
	level  int
}

func NewInstruction(tokens []rune, tape *Tape) *Instruction {
	return &Instruction{ptr: -1, tokens: tokens, tape: tape, level: 0}
}

// Fetch increase moves the pointer to the next tokens.
// Returns whether it successfully fetch and execute one.
func (instr *Instruction) Fetch() bool {
	instr.ptr++
	if instr.ptr >= len(instr.tokens) {
		return false
	}

	instr.Execute()
	return true
}

// Execute modifies the memory tape or reads from input or writes to output
// depending on the operator.
// >: increments the data pointer (points to the next cell on the right).
// <: decrements the data pointer (points to the previous cell on the left).
// +: increases by one the byte at the data pointer.
// -: decreases by one the byte at the data pointer.
// .: output the byte at the data pointer, using the ASCII character encoding.
// ,: accept one byte of input, storing its value in the byte at the pointer.
// [: if the byte at the data pointer is zero, jump forward to the command
// after ']'; move forward to the next command otherwise.
// ]: if the byte at the data pointer is nonzero, jump back to the command
// after matching '['; move forward to the next command otherwise.
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
	case ',':
		buf := make([]byte, 1)
		tape.in.Read(buf)
		tape.data[tape.ptr] = int8(buf[0])
	case '[':
		if tape.data[tape.ptr] != 0 {
			instr.level++
			instr.Fetch()
		} else {
			currentLevel := instr.level

			for {
				instr.ptr++
				if instr.tokens[instr.ptr] == '[' {
					instr.level++
				}

				if instr.tokens[instr.ptr] == ']' {
					if currentLevel == instr.level {
						break
					}
					instr.level--
				}
			}
			instr.Fetch()
		}
	case ']':
		if tape.data[tape.ptr] == 0 {
			instr.level--
			instr.Fetch()
		} else {
			currentLevel := instr.level

			for {
				instr.ptr--
				if instr.tokens[instr.ptr] == ']' {
					instr.level--
				}

				if instr.tokens[instr.ptr] == '[' {
					if currentLevel == instr.level {
						break
					}
					instr.level++
				}
			}
			instr.Fetch()
		}
	default:
		return
	}
}
