package utils

import (
	"sync"
	// "golang.org/x/sync/semaphore"
)

type ByteRingBuffer struct {
	buff     []byte
	l        int //buffer length
	h        int //head position
	t        int //tail position
	Capacity int

	// rlck semaphore.Weighted
	wlck sync.Mutex
}

func NewByteRingBuffer(buff_len int) *ByteRingBuffer {
	rb := &ByteRingBuffer{
		buff:     make([]byte, buff_len),
		l:        buff_len,
		h:        0,
		t:        0,
		Capacity: 0,
	}
	return rb
}

/**
 * Push data to the ring buffer
 * return true if the buffer is overflowed
 */
func (rb *ByteRingBuffer) Push(bs []byte, n int) int {
	//parameter check
	if n > rb.l {
		//over push not supported
		return -1
	}

	if n > len(bs) {
		n = len(bs)
	}

	if rb.l-rb.t >= n {
		//whole copy
		copy(rb.buff[rb.t:rb.t+n], bs)

		//update the head position
		if rb.t+n >= rb.h && rb.t < rb.h {
			rb.h = rb.t
			if rb.h == rb.l {
				rb.h = 0
			}
		}

		//update the tail position
		rb.t = rb.t + n
		if rb.t == rb.l {
			rb.t = 0
		}
		if rb.t >= rb.h {
			rb.h = rb.t + 1
			if rb.h == rb.l {
				rb.h = 0
			}
		}

	} else {
		//split copy
		copy(rb.buff[rb.t:rb.l], bs[0:rb.l-rb.t])
		copy(rb.buff[0:n-rb.l+rb.t], bs[rb.l-rb.t:n])

		//updat the tail position
		rb.t = n - rb.l + rb.t - 1
		//update the head position
		if rb.t >= rb.h {
			rb.h = rb.t + 1
			if rb.h == rb.l {
				rb.h = 0
			}
		}
	}

	//update capacity
	if rb.Capacity+n >= rb.l {
		rb.Capacity = rb.l
	} else {
		rb.Capacity += n
	}
	return n
}

func (rb *ByteRingBuffer) Pop(bs []byte, n int) int {
	//parameter check
	if n > rb.l {
		return -1 //over pop not supported
	}
	if n > len(bs) {
		n = len(bs)
	}
	if rb.Capacity < n {
		n = rb.Capacity
	}

	if rb.l-rb.h > n {
		//whole copy
		copy(bs, rb.buff[rb.h:rb.h+n])

		rb.h = rb.h + n
	} else {
		//split copy
		copy(bs, rb.buff[rb.h:rb.l])
		copy(bs[rb.l-rb.h:], rb.buff[0:n-rb.l+rb.h])

		//update tail position
		rb.h = n - rb.l + rb.h
		if rb.t <= rb.h {
			rb.t = rb.h
		}
	}

	//update capacity
	rb.Capacity -= n
	if rb.Capacity < 0 {
		rb.Capacity = 0
	}

	return n
}

func (rb *ByteRingBuffer) Availability() int {
	return rb.l - rb.Capacity
}

/**
 * Read data from the ring buffer without moving the head pointer
 */
func (rb *ByteRingBuffer) Peek(bs []byte, n int) int {
	//parameter check
	if n > rb.l {
		return -1 //over pop not supported
	}

	if n > len(bs) {
		n = len(bs)
	} else if rb.Capacity < n {
		n = rb.Capacity
	}

	if rb.l-rb.h > n {
		//whole copy
		copy(bs, rb.buff[rb.h:rb.h+n])

	} else {
		//split copy
		copy(bs, rb.buff[rb.h:rb.l])
		copy(bs[rb.l-rb.h:], rb.buff[0:n-rb.l+rb.h])

	}
	return n
}

/**
 * Drop n bytes from the head
 */
func (rb *ByteRingBuffer) Drop(n int) int {
	//parameter check
	if n > rb.l {
		return -1 //over drop not supported
	}
	if rb.Capacity < n {
		n = rb.Capacity
	}

	if rb.l-rb.h > n {
		rb.h = rb.h + n
	} else {
		rb.h = n - rb.l + rb.h
	}

	//update capacity
	rb.Capacity -= n
	if rb.Capacity < 0 {
		rb.Capacity = 0
	}
	return n
}

func (rb *ByteRingBuffer) Flush() {
	rb.Capacity = 0
	rb.h = 0
	rb.t = 0
}
