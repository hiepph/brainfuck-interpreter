package brainfuck

import (
	"bytes"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	program := `++
[ >- ]`

	got := Lex(strings.NewReader(program))
	want := []rune{'+', '+', '[', '>', '-', ']'}

	assertTokens(t, got, want)
}

func NewTokens(program string) []rune {
	return Lex(strings.NewReader(program))
}

func TestInstruction(t *testing.T) {
	t.Run("Do nothing if there are no instructions left to fetch", func(t *testing.T) {
		program := "+"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // +
		assertCannotFetch(t, instr.Fetch())
	})

	t.Run("moves forward to the next token", func(t *testing.T) {
		tokens := NewTokens("+-")
		instr := NewInstruction(tokens, NewTape(nil, nil))

		instr.Fetch()
		assertInstructionPointer(t, instr, 0)

		instr.Fetch()
		assertInstructionPointer(t, instr, 1)
	})

	t.Run("+ increments and - decrements the byte at data pointer", func(t *testing.T) {
		tokens := NewTokens("+-")
		instr := NewInstruction(tokens, NewTape(nil, nil))

		instr.Fetch()
		assertTapePointer(t, instr.tape, 0)
		assertTapeValue(t, instr.tape, 1)

		instr.Fetch()
		assertTapePointer(t, instr.tape, 0)
		assertTapeValue(t, instr.tape, 0)
	})

	t.Run("> moves the pointer to the right and < moves it to the left", func(t *testing.T) {
		tokens := NewTokens("><")
		instr := NewInstruction(tokens, NewTape(nil, nil))

		instr.Fetch()
		assertTapePointer(t, instr.tape, 1)

		instr.Fetch()
		assertTapePointer(t, instr.tape, 0)
	})

	t.Run(". output the byte at the data pointer", func(t *testing.T) {
		var program strings.Builder
		for i := 0; i < 72; i++ {
			program.WriteRune('+')
		}
		program.WriteRune('.')
		tokens := NewTokens(program.String())

		out := &bytes.Buffer{}
		instr := NewInstruction(tokens, NewTape(nil, out))

		for ok := instr.Fetch(); ok; ok = instr.Fetch() {
		}
		assertTapePointer(t, instr.tape, 0)
		assertOutput(t, out.String(), "H")
	})

	t.Run("[ moves the pointer forward if the current byte is non-zero", func(t *testing.T) {
		program := "+[>+]"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: [
		assertInstructionPointer(t, instr, 2)
	})

	t.Run("[ jumps the pointer after the matching ] if current byte is zero", func(t *testing.T) {
		program := "[>+]+"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: [ -> 4: +
		assertInstructionPointer(t, instr, 4)
	})

	t.Run("] moves the pointer foward if the current byte is zero", func(t *testing.T) {
		program := "]+"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch()
		assertInstructionPointer(t, instr, 1)
		assertTapeValue(t, instr.tape, 1)
	})

	t.Run("] moves the pointer after the matching [ if the current byte is non-zero", func(t *testing.T) {
		program := "++[-]"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: +
		instr.Fetch() // 2: [ -> 3: -
		assertInstructionPointer(t, instr, 3)
		assertTapeValue(t, instr.tape, 1)
		instr.Fetch() // 4: ] -> 3: -
		assertInstructionPointer(t, instr, 3)
		assertTapeValue(t, instr.tape, 0)
	})

	t.Run(`[ jumps the ponter after the matching ]
if the current byte is zero (double stacking brackets)`, func(t *testing.T) {
		program := "[+[>+++]+]+"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: [ -> 10: +
		assertInstructionPointer(t, instr, 10)
	})

	t.Run(`[ jumps the pointer after the matching ]
if the current byte is zero (triple stacking brackets)`, func(t *testing.T) {
		program := "+[-[>[+]+]+]"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: [ -> 2: -
		instr.Fetch() // 3: [ -> 10: ]
		assertInstructionPointer(t, instr, 10)
	})

	t.Run(`] moves the pointer after the matching [
if the current byte is non-zero (double stacking brackets)`, func(t *testing.T) {
		program := "+[-[]+]"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: [ -> 2: -
		instr.Fetch() // 3: [ -> 5: +
		assertInstructionPointer(t, instr, 5)
		assertInstructionValue(t, instr, '+')
		instr.Fetch() // 6: ] -> 2: -
		assertInstructionPointer(t, instr, 2)
		assertInstructionValue(t, instr, '-')
	})

	t.Run(`] moves the pointer after the matching [
if the current byte is non-zero (triple stacking brackets)`, func(t *testing.T) {
		program := "+[[-[]+]+]"
		instr := NewInstruction(NewTokens(program), NewTape(nil, nil))

		instr.Fetch() // 0: +
		instr.Fetch() // 1: [ -> 2: [ -> 3: -
		instr.Fetch() // 4: [ -> 6: +
		instr.Fetch() // 7: [ -> 3: -
		assertInstructionPointer(t, instr, 3)
		assertInstructionValue(t, instr, '-')
	})

	t.Run(", receives one char from input and stores it in the current byte", func(t *testing.T) {
		in := &bytes.Buffer{}
		instr := NewInstruction(NewTokens(","), NewTape(in, nil))

		in.WriteRune('a')
		instr.Fetch()
		assertInstructionPointer(t, instr, 0)
		assertTapeValue(t, instr.tape, 'a')
	})

	t.Run(", reads multiple chars from input", func(t *testing.T) {
		in := strings.NewReader("ab")
		instr := NewInstruction(NewTokens(",>,"), NewTape(in, nil))

		for i := 0; i < 3; i++ {
			instr.Fetch()
		}
		assertTapePointer(t, instr.tape, 1)
		assertTapeValue(t, instr.tape, 'b')
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

func assertInstructionValue(t *testing.T, instr *Instruction, want rune) {
	t.Helper()

	got := instr.tokens[instr.ptr]
	if got != want {
		t.Errorf("pointer: got %c want %c", got, want)
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

	got := tape.data[tape.ptr]
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
		name    string
		program string
		input   string
		want    string
	}{
		{
			name:    "Single character",
			program: "++ > +++++ [ <+ >- ] ++++ ++++ [ <+++ +++ >- ] < . ",
			want:    "7",
		},
		{
			name:    "string",
			program: ">+++++++++[<++++++>-]<...",
			want:    "666",
		},
		{
			name: "Hello World",
			program: `++++++++[>++++[>++>+++>+++>+<<<<-]>+>->+>>+[<]<-]>>.>
>---.+++++++..+++.>.<<-.>.+++.------.--------.>+.>++.`,
			want: "Hello World!\n",
		},
		{
			name:    "echo input string",
			program: ",[.,]",
			input:   "abc",
			want:    "abc",
		},
		{
			name:    "reverse input string",
			program: ">,[>,]<[.<]",
			input:   "abc",
			want:    "cba",
		},
		{
			name: "wc",
			program: `
>>>+>>>>>+>>+>>+[<<],[
    -[-[-[-[-[-[-[-[<+>-[>+<-[>-<-[-[-[<++[<++++++>-]<
        [>>[-<]<[>]<-]>>[<+>-[<->[-]]]]]]]]]]]]]]]]
    <[-<<[-]+>]<<[>>>>>>+<<<<<<-]>[>]>>>>>>>+>[
        <+[
            >+++++++++<-[>-<-]++>[<+++++++>-[<->-]+[+>>>>>>]]
            <[>+<-]>[>>>>>++>[-]]+<
        ]>[-<<<<<<]>>>>
    ],
]+<++>>>[[+++++>>>>>>]<+>+[[<++++++++>-]<.<<<<<]>>>>>>>>]
`,
			input: `This is a line.
This is another line.`,
			want: "	1	8	37\n",
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			in := strings.NewReader(test.input)
			out := &bytes.Buffer{}
			instr := NewInstruction(NewTokens(test.program), NewTape(in, out))

			for ok := instr.Fetch(); ok; ok = instr.Fetch() {
			}
			assertOutput(t, out.String(), test.want)
		})
	}
}
