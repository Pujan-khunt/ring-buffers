package ring

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

type MagicBuffer struct {
	data     []byte
	len      int
	readIdx  int
	writeIdx int
}

func NewMagicBuffer(requestedSize int) (*MagicBuffer, error) {
	pageSize := os.Getpagesize()

	// Ensure buffer size remains a multiple of processor page size.
	if requestedSize%pageSize != 0 {
		requestedSize = ((requestedSize / pageSize) + 1) * pageSize
	}

	// Creates an anonymous file which behaves as a regular file and can be modified, truncated and memory mapped.
	// It lives in RAM and doesn't have a filesystem backing but rather a volatile backing storage (RAM).
	// It is released once all references to the file are dropped i.e. once all the memory maps and the fd are closed.
	fd, err := unix.MemfdCreate("ring-buffer-mem", 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)

	// Resize anonymous file to twice the size of buffer
	if err = unix.Ftruncate(fd, 2*int64(requestedSize)); err != nil {
		return nil, err
	}

	// Reserve the placeholder virtual address space that is required for the buffer.
	// It is required so two virtual pages don't end up at different locations.
	// This reservation will ensure that both virtual pages end up adjacent to each other.
	reserved, err := unix.Mmap(-1, 0, 2*requestedSize, unix.PROT_NONE, unix.MAP_ANONYMOUS|unix.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}

	// Retrieve the address to calculate starting address of the second half of the reserved memory region
	// Converts a pointer into an integer so that pointer arithmetic can be performed.
	addr := uintptr(unsafe.Pointer(&reserved[0]))

	// Map the first half of the virtual address space
	_, _, errno := unix.Syscall6(unix.SYS_MMAP, addr, uintptr(requestedSize), uintptr(unix.PROT_READ|unix.PROT_WRITE), uintptr(unix.MAP_SHARED|unix.MAP_FIXED), uintptr(fd), 0)
	if errno != 0 {
		unix.Munmap(reserved)
		return nil, errno
	}

	// Map the second half of the virtual address space
	_, _, errno = unix.Syscall6(
		unix.SYS_MMAP,
		addr+uintptr(requestedSize), // The address: Start + Size
		uintptr(requestedSize),      // Length
		uintptr(unix.PROT_READ|unix.PROT_WRITE),
		uintptr(unix.MAP_SHARED|unix.MAP_FIXED),
		uintptr(fd),
		0,
	)
	if errno != 0 {
		unix.Munmap(reserved)
		return nil, errno
	}

	data := unsafe.Slice((*byte)(unsafe.Pointer(addr)), 2*requestedSize)

	return &MagicBuffer{
		data: data,
		len:  requestedSize,
	}, nil
}
