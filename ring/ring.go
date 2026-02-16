package ring

import "io"

// Buffer defines the behavior of a ring buffer.
// It inherits [io.Reader] and [io.Writer] for compatibility.
type Buffer interface {
	io.Reader
	io.Writer

	// Returns all available bytes in the buffer without advancing the read pointer.
	// Useful for peeking.
	Bytes() []byte

	// Returns the number of readable bytes.
	Length() uint

	// Returns the total size of the buffer.
	Capacity() uint

	// Clears out entire buffer
	Reset()
}
