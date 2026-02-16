package ring

import (
	"slices"
	"testing"
)

func TestBuffer(t *testing.T) {
	type step struct {
		action      string // "READ" or "WRITE"
		input       []byte // Data to write, or buffer to read
		expectCount int    // Count of bytes to expect on read/write
		expectData  []byte // Data to expect in case of reads.
	}
	tests := []struct {
		name  string
		steps []step
		size  uint
	}{
		{
			name: "write full and read full",
			size: uint(5),
			steps: []step{
				{action: "WRITE", input: []byte{1, 2, 3, 4, 5}, expectCount: 5, expectData: nil},            // write full
				{action: "READ", input: make([]byte, 5), expectCount: 5, expectData: []byte{1, 2, 3, 4, 5}}, // read full
			},
		},
		{
			name: "write partial and read full",
			size: uint(5),
			steps: []step{
				{action: "WRITE", input: []byte{1, 2, 3}, expectCount: 3, expectData: nil},                  // write partial
				{action: "READ", input: make([]byte, 5), expectCount: 3, expectData: []byte{1, 2, 3, 0, 0}}, // read full
			},
		},
		{
			name: "write full and read partial",
			size: uint(5),
			steps: []step{
				{action: "WRITE", input: []byte{1, 2, 3, 4, 5}, expectCount: 5, expectData: nil},      // write full
				{action: "READ", input: make([]byte, 3), expectCount: 3, expectData: []byte{1, 2, 3}}, // read partial
			},
		},
		{
			name:  "warp around logic",
			size:  5,
			steps: []step{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rb, _ := NewStandardBuffer(tc.size)

			for i, step := range tc.steps {
				var err error
				var n int

				switch step.action {
				case "READ":
					n, err = rb.Read(step.input)
				case "WRITE":
					n, err = rb.Write(step.input)
				}

				if err != nil {
					t.Errorf("Error received from READ/WRITE method: %v", err)
				}
				if n != step.expectCount {
					t.Errorf("Step %d (%s): got count %d, want %d", i, step.action, n, step.expectCount)
				}

				if step.action == "READ" && step.expectData != nil {
					if !slices.Equal(step.input, step.expectData) {
						t.Errorf("Step %d (READ): got %v, want: %v", i, step.input, step.expectData)
					}
				}
			}
		})
	}
}
