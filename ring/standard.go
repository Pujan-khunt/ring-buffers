package ring

import "errors"

var (
	ErrBufferFull  = errors.New("ring: buffer full")
	ErrInvalidSize = errors.New("ring: size must be greater than 0")
)

type StandardBuffer struct {
	buf    []byte
	size   uint
	length uint
	r      uint
	w      uint
}

func NewStandardBuffer(size uint) (*StandardBuffer, error) {
	if size == 0 {
		return nil, ErrInvalidSize
	}

	return &StandardBuffer{
		buf:  make([]byte, size),
		size: uint(size),
		r:    0,
		w:    0,
	}, nil
}

func (sb *StandardBuffer) Read(b []byte) (int, error) {
	n := len(b)
	available := sb.readableBytes()
	if n > available {
		n = available
	}

	for i := range n {
		b[i] = sb.buf[sb.r]
		sb.r = (sb.r + 1) % sb.size
		sb.length--
	}
	return n, nil
}

func (sb *StandardBuffer) Write(b []byte) (int, error) {
	n := len(b)
	available := sb.writableBytes()
	if n > available {
		n = available
	}

	for i := range n {
		sb.buf[sb.w] = b[i]
		sb.w = (sb.w + 1) % sb.size
		sb.length++
	}
	return n, nil
}

func (sb *StandardBuffer) Length() uint {
	return sb.length
}

func (sb *StandardBuffer) Capacity() uint {
	return sb.size
}

func (sb *StandardBuffer) readableBytes() int {
	return int(sb.length)
}

func (sb *StandardBuffer) writableBytes() int {
	return int(sb.size) - int(sb.length)
}
