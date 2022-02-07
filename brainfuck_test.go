package main

import "testing"

func TestOperator(t *testing.T) {
	var tape Tape

	t.Run("Brainfuck's tape is initialized with 30000 bytes of 0", func(t *testing.T) {
		wantLength := 30000
		if len(tape) != wantLength {
			t.Errorf("length: got %d want %d", len(tape), wantLength)
		}

		for i := 0; i < wantLength; i++ {
			if tape[i] != 0 {
				t.Errorf("tape is not initialized with all 0")
			}
		}
	})
}
