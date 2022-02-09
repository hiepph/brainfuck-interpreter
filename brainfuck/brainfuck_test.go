package brainfuck

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	stream := `++
[ >- ]`

	got := Lex(strings.NewReader(stream))
	want := []rune{'+', '+', '[', '>', '-', ']'}

	assertTokens(t, got, want)
}

func NewTokens(stream string) []rune {
	return Lex(strings.NewReader(stream))
}

func TestInstruction(t *testing.T) {
	t.Run("moves forward to the next token", func(t *testing.T) {
		tokens := NewTokens("+-")
		instr := NewInstruction(tokens, NewTape(nil))

		instr.Next()
		assertInstructionPointer(t, instr, 0)

		instr.Next()
		assertInstructionPointer(t, instr, 1)
	})

	t.Run("+ increments and - decrements the byte at data pointer", func(t *testing.T) {
		tokens := NewTokens("+-")
		instr := NewInstruction(tokens, NewTape(nil))

		instr.Next()
		assertTapePointer(t, instr.tape, 0)
		assertTapeValue(t, instr.tape, 1)

		instr.Next()
		assertTapePointer(t, instr.tape, 0)
		assertTapeValue(t, instr.tape, 0)
	})

	t.Run("> moves the pointer to the right and < moves it to the left", func(t *testing.T) {
		tokens := NewTokens("><")
		instr := NewInstruction(tokens, NewTape(nil))

		instr.Next()
		assertTapePointer(t, instr.tape, 1)

		instr.Next()
		assertTapePointer(t, instr.tape, 0)
	})

	// t.Run(". output the byte at the data pointer", func(t *testing.T) {
	// 	out := &bytes.Buffer{}
	// 	tape := NewTape(out)

	// 	for i := 0; i < 72; i++ {
	// 		tape.Command('+')
	// 	}

	// 	tape.Command('.')
	// 	assertPointer(t, tape, 0)
	// 	assertOutput(t, out.String(), "H")
	// })

	// t.Run("[ skip")
}

func assertTokens(t testing.TB, got, want []rune) {
	if len(got) != len(want) {
		t.Fatalf("Lexing error: different length, got %d want %d", len(got), len(want))
	}

	for i, _ := range got {
		if got[i] != want[i] {
			t.Fatalf("Token[%d] is wrong: got %c want %c", i, got[i], want[i])
		}
	}
}

func assertInstructionPointer(t *testing.T, instr *Instruction, want int) {
	t.Helper()

	got := instr.ptr
	if got != want {
		t.Errorf("pointer: got %d want %d", got, want)
	}
}

func assertTapePointer(t *testing.T, tape *Tape, want int) {
	t.Helper()

	got := tape.ptr
	if got != want {
		t.Errorf("pointer: got %d want %d", got, want)
	}
}

func assertTapeValue(t *testing.T, tape *Tape, want int) {
	t.Helper()

	got := tape.data[0]
	if got != want {
		t.Errorf("value: got %d want %d", got, want)
	}
}
