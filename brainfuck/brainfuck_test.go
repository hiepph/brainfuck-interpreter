package brainfuck

import (
	"bytes"
	"testing"
)

func TestInitialization(t *testing.T) {
	tape := NewTape(nil)

	wantLength := 30000
	gotLength := len(tape.data)
	if gotLength != wantLength {
		t.Errorf("length: got %d want %d", gotLength, wantLength)
	}

	for i := 0; i < wantLength; i++ {
		if tape.data[i] != 0 {
			t.Fatalf("tape is not initialized with all 0")
		}
	}
}

func TestOperators(t *testing.T) {
	t.Run("+ increments and - decrements the byte at data pointer", func(t *testing.T) {
		tape := NewTape(nil)

		tape.Command('+')

		assertPointer(t, tape, 0)
		assertValue(t, tape, 1)

		tape.Command('-')
		assertPointer(t, tape, 0)
		assertValue(t, tape, 0)
	})

	t.Run("> moves the pointer to the right and < moves it to the left", func(t *testing.T) {
		tape := NewTape(nil)

		tape.Command('>')
		assertPointer(t, tape, 1)

		tape.Command('>')
		assertPointer(t, tape, 2)

		tape.Command('<')
		assertPointer(t, tape, 1)
	})

	t.Run(". output the byte at the data pointer", func(t *testing.T) {
		out := &bytes.Buffer{}
		tape := NewTape(out)

		for i := 0; i < 72; i++ {
			tape.Command('+')
		}

		tape.Command('.')
		assertPointer(t, tape, 0)
		assertOutput(t, out.String(), "H")
	})
}

func assertPointer(t *testing.T, tape *Tape, want int) {
	t.Helper()

	got := tape.ptr
	if got != want {
		t.Errorf("pointer: got %d want %d", got, want)
	}
}

func assertValue(t *testing.T, tape *Tape, want int) {
	t.Helper()

	got := tape.data[0]
	if got != want {
		t.Errorf("value: got %d want %d", got, want)
	}
}

func assertOutput(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("output: got %s want %s", got, want)
	}
}
