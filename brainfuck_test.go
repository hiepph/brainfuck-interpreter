package main

import "testing"

func TestInitialization(t *testing.T) {
	tape := newTape()

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
		tape := newTape()

		command(tape, '+')

		assertPointer(t, tape, 0)
		assertValue(t, tape, 1)

		command(tape, '-')
		assertPointer(t, tape, 0)
		assertValue(t, tape, 0)
	})

	t.Run("> moves the pointer to the right and < moves it to the left", func(t *testing.T) {
		tape := newTape()

		command(tape, '>')
		assertPointer(t, tape, 1)

		command(tape, '>')
		assertPointer(t, tape, 2)

		command(tape, '<')
		assertPointer(t, tape, 1)
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
