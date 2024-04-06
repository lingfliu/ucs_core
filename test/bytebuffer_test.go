package test

import (
	"testing"

	"github.com/lingfliu/ucs_core/utils"
)

func TestRingBuffer(t *testing.T) {
	rb := utils.NewByteRingBuffer(10)

	bs1 := []byte("123456")
	bs2 := []byte("7890ab")
	rb.Push(bs1, 6)
	rb.Push(bs2, 5)
	bs := make([]byte, 10)

	rb.Peek(bs, 5)
	if string(bs[:5]) != "23456" {
		t.Fatalf("expected 12345, got %s", string(bs))
	}

	rb.Pop(bs, 5)
	if string(bs[:5]) != "23456" {
		t.Fatalf("expected 12345, got %s", string(bs))
	}

	rb.Pop(bs, 5)
	if string(bs[:5]) != "7890a" {
		t.Fatalf("expected 7890a, got %s", string(bs))
	}

	bs3 := []byte("cdefg")
	rb.Push(bs3, 5)
	rb.Pop(bs, 5)
	if string(bs[:5]) != "cdefg" {
		t.Fatalf("expected bcdef, got %s", string(bs))
	}

	bs4 := []byte("hijklmnp")
	for i := 0; i < 100; i++ {
		rb.Push(bs4, 8)
		rb.Pop(bs, 8)
		if string(bs[:8]) != "hijklmnp" {
			t.Fatalf("expected hijklmnp, got %s", string(bs))
		}
	}

}
