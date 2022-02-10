package brainfuck

import (
	"bytes"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `++
[ >- ]`

	got := Lex(strings.NewReader(input))
	want := []rune{'+', '+', '[', '>', '-', ']'}

	assertTokens(t, got, want)
}

func NewTokens(input string) []rune {
	return Lex(strings.NewReader(input))
}

func TestInstruction(t *testing.T) {
	t.Run("Do nothing if there are no instructions left to fetch", func(t *testing.T) {
		input := "+"
		instr := NewInstruction(NewTokens(input), NewTape(nil))

		instr.Fetch() // +
		assertCannotFetch(t, instr.Fetch())
	})

	t.Run("moves forward to the next token", func(t *testing.T) {
		tokens := NewTokens("+-")
		instr := NewInstruction(tokens, NewTape(nil))

		instr.Fetch()
		assertInstructionPointer(t, instr, 0)

		instr.Fetch()
		assertInstructionPointer(t, instr, 1)
	})

	t.Run("+ increments and - decrements the byte at data pointer", func(t *testing.T) {
		tokens := NewTokens("+-")
		instr := NewInstruction(tokens, NewTape(nil))

		instr.Fetch()
		assertTapePointer(t, instr.tape, 0)
		assertTapeValue(t, instr.tape, 1)

		instr.Fetch()
		assertTapePointer(t, instr.tape, 0)
		assertTapeValue(t, instr.tape, 0)
	})

	t.Run("> moves the pointer to the right and < moves it to the left", func(t *testing.T) {
		tokens := NewTokens("><")
		instr := NewInstruction(tokens, NewTape(nil))

		instr.Fetch()
		assertTapePointer(t, instr.tape, 1)

		instr.Fetch()
		assertTapePointer(t, instr.tape, 0)
	})

	t.Run(". output the byte at the data pointer", func(t *testing.T) {
		var input strings.Builder
		for i := 0; i < 72; i++ {
			input.WriteRune('+')
		}
		input.WriteRune('.')
		tokens := NewTokens(input.String())

		out := &bytes.Buffer{}
		instr := NewInstruction(tokens, NewTape(out))

		for ok := instr.Fetch(); ok; ok = instr.Fetch() {
		}
		assertTapePointer(t, instr.tape, 0)
		assertOutput(t, out.String(), "H")
	})

	t.Run("[ moves the pointer forward if the current byte is non-zero", func(t *testing.T) {
		input := "+[>+]"
		instr := NewInstruction(NewTokens(input), NewTape(nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: [
		assertInstructionPointer(t, instr, 2)
	})

	t.Run("[ jumps the pointer after the matching ] if current byte is zero", func(t *testing.T) {
		input := "[>+]+"
		instr := NewInstruction(NewTokens(input), NewTape(nil))

		instr.Fetch() // 0: [ -> 4: +
		assertInstructionPointer(t, instr, 4)
	})

	t.Run("] moves the pointer foward if the current byte is zero", func(t *testing.T) {
		input := "]+"
		instr := NewInstruction(NewTokens(input), NewTape(nil))

		instr.Fetch()
		assertInstructionPointer(t, instr, 1)
		assertTapeValue(t, instr.tape, 1)
	})

	t.Run("] moves the pointer after the matching [ if the current byte is non-zero", func(t *testing.T) {
		input := "++[-]"
		instr := NewInstruction(NewTokens(input), NewTape(nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: +
		instr.Fetch() // 2: [ -> 3: -
		assertInstructionPointer(t, instr, 3)
		assertTapeValue(t, instr.tape, 1)
		instr.Fetch() // 4: ] -> 3: -
		assertInstructionPointer(t, instr, 3)
		assertTapeValue(t, instr.tape, 0)
	})
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

func assertTapeValue(t *testing.T, tape *Tape, want int8) {
	t.Helper()

	got := tape.data[0]
	if got != want {
		t.Errorf("value: got %d want %d", got, want)
	}
}

func assertOutput(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("output: got %q want %q", got, want)
	}
}

func assertCannotFetch(t *testing.T, isAvailableToFetch bool) {
	t.Helper()
	if isAvailableToFetch {
		t.Errorf("Shouldn't fetch anymore instruction.")
	}
}

func TestIntegration(t *testing.T) {
	table := []struct {
		name  string
		input string
		want  string
	}{
		{
			// refer: https://www.wikiwand.com/en/Brainfuck
			name:  "Single character",
			input: "++ > +++++ [ <+ >- ] ++++ ++++ [ <+++ +++ >- ] < . ",
			want:  "7",
		},
		// 		{
		// 			name: "Hello World",
		// 			input: `++++++++[>++++[>++>+++>+++>+<<<<-]>+>->+>>+[<]<-]>>.>
		// >---.+++++++..+++.>.<<-.>.+++.------.--------.>+.>++.`,
		// 			want: "Hello World!",
		// 		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			instr := NewInstruction(NewTokens(test.input), NewTape(out))

			for ok := instr.Fetch(); ok; ok = instr.Fetch() {
			}
			assertOutput(t, out.String(), test.want)
		})
	}
}
