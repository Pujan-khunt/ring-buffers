package ringbuffer

type RingBuffer interface {
	Read(val *byte) int
	Write(val byte) int
	isEmpty() bool
	IsFull() bool
}
