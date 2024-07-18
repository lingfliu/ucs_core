package utils

import (
	"sync"
)

/**
 * A simple byte buffer wrapper
 */
type ByteArrayRingBuffer struct {
	buff [][]byte
	l    int //buffer length
	c    int //head position

	// rlck semaphore.Weighted
	wlck sync.Mutex
}

func NewByteArrayRingBuffer(ring_size int, buff_size int) *ByteArrayRingBuffer {
	rb := &ByteArrayRingBuffer{
		buff: make([][]byte, ring_size),
		l:    ring_size,
		c:    0,
	}

	for i := 0; i < ring_size; i++ {
		rb.buff[i] = make([]byte, buff_size)
	}
	return rb
}

func (rb *ByteArrayRingBuffer) Push(bs []byte) {
}

func (rb *ByteArrayRingBuffer) Pop() []byte {
	return nil
}

func (rb *ByteArrayRingBuffer) Len() int {
	return 0
}

// func (rb *ByteArrayRingBuffer) Curr() []byte {
// 	return rb.buff[rb.c]
// }

// func (rb *ByteArrayRingBuffer) Next() []byte {
// 	if rb.c == rb.l-1 {
// 		rb.c = 0
// 	} else {
// 		rb.c++
// 	}

// 	return rb.buff[rb.c]
// }
