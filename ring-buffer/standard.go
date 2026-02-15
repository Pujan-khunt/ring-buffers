package ringbuffer

type StandardRB struct {
	data     []byte // Slice holding the data of the ringbuffer
	len      uint   // Length of the ringbuffer
	readIdx  uint   // Index from where the ringbuffer data would be read.
	writeIdx uint   // Index from where the ringbuffer data would be written.
}

func NewRingBuffer(len int) *StandardRB {
	data := make([]byte, len+1)
	return &StandardRB{
		data:     data,
		len:      uint(len + 1),
		readIdx:  0,
		writeIdx: 0,
	}
}

func (rb *StandardRB) Read(val *byte) int {
	if rb.readIdx == rb.writeIdx {
		return 0
	}
	*val = rb.data[rb.readIdx]
	rb.readIdx = (rb.readIdx + 1) % rb.len
	return 1
}

func (rb *StandardRB) Write(val byte) int {
	if (rb.writeIdx+1)%rb.len == rb.readIdx {
		return 0
	}
	rb.data[rb.writeIdx] = val
	rb.writeIdx = (rb.writeIdx + 1) % rb.len
	return 1
}

func (rb *StandardRB) isEmpty() bool {
	return rb.writeIdx == rb.readIdx
}

func (rb *StandardRB) isFull() bool {
	return (rb.writeIdx+1)%rb.len == rb.readIdx
}
