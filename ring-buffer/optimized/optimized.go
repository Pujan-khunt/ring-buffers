package ringbuffer

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type OptimizedRB struct {
	data     []byte
	len      int
	readIdx  int
	writeIdx int
}

func NewRingBuffer(len int) (*OptimizedRB, error) {
	pageSize := os.Getpagesize()
	if len%pageSize != 0 {
		return nil, fmt.Errorf("length must be a multiple of page size: %d", pageSize)
	}

	file, err := os.CreateTemp("", "ring-buffer-shm")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if err := file.Truncate(int64(len)); err != nil {
		return nil, err
	}
	fd := int(file.Fd())

	// Reserve the block, not yet occupied.
	addr, _, errno := syscall.Syscall6(syscall.SYS_MMAP, 0, uintptr(2*len), syscall.PROT_NONE, syscall.MAP_ANON|syscall.MAP_PRIVATE, 0, 0)
	if errno != 0 {
		return nil, errno
	}

	// Occupy the first region [0...N)
	_, _, errno = syscall.Syscall6(
		syscall.SYS_MMAP,
		addr,
		uintptr(len),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_FIXED,
		uintptr(fd),
		0,
	)
	if errno != 0 {
		return nil, errno
	}

	// Occupy the second region [N, 2N)
	_, _, errno = syscall.Syscall6(
		syscall.SYS_MMAP,
		addr+uintptr(len),
		uintptr(len),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_FIXED,
		uintptr(fd),
		0,
	)
	if errno != 0 {
		return nil, errno
	}

	// Convert raw addr into byte slice
	var data []byte
	data = unsafe.Slice((*byte)(unsafe.Pointer(addr)), 2*len)

	return &OptimizedRB{
		data: data,
		len:  len,
	}, nil
}

// Write writes data to the buffer.
func (rb *OptimizedRB) Write(p []byte) (int, error) {
	n := len(p)

	// Simply copy into the slice at the current write index.
	// Even if the write index is at the end, the MMU maps it to the start automatically.
	copy(rb.data[rb.writeIdx:], p)

	rb.writeIdx += n

	// If we have advanced into the mirror region, snap back.
	if rb.writeIdx >= rb.len {
		rb.writeIdx -= rb.len
	}

	return n, nil
}

// Read reads data from the buffer.
func (rb *OptimizedRB) Read(p []byte) (int, error) {
	n := len(p)

	copy(p, rb.data[rb.readIdx:])

	rb.readIdx += n

	// Snap back if we drift into the mirror region
	if rb.readIdx >= rb.len {
		rb.readIdx -= rb.len
	}

	return n, nil
}
