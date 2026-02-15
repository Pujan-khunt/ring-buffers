package main

import (
	"fmt"

	ringbuffer "github.com/Pujan-khunt/ring-buffers/ring-buffer"
)

func main() {
	rb := ringbuffer.NewRingBuffer(3)
	// arr := []byte{4, 8, 6, 7} // Will cause '7' to fail because of size mismatch
	// arr := []byte{} // Will cause read failure, since its empty
	arr := []byte{4, 8, 6}
	for _, val := range arr {
		if out := rb.Write(val); out == 1 {
			fmt.Printf("Write successful: %d\n", val)
		} else {
			fmt.Printf("Write Failure: %d\n", val)
		}
	}

	for range len(arr) {
		val := byte(0)
		if out := rb.Read(&val); out == 1 {
			fmt.Printf("Read successful: %d\n", val)
		} else {
			fmt.Printf("Read Failure: %d\n", val)
		}
	}
}
