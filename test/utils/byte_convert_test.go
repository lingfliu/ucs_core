package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lingfliu/ucs_core/utils"
)

func TestByte2Int(t *testing.T) {
	var x int = 0
	var y int = 0

	x = 15668
	y = 0
	bs1 := make([]byte, 4)

	utils.Int2Byte(x, bs1, 0, 4, false, true)
	y = utils.Byte2Int(bs1, 0, 4, false, true)

	if x != y {
		t.Fatalf(fmt.Sprintf("expected %d, got %d", x, y))
	}

	x = -15668
	y = 0

	utils.Int2Byte(x, bs1, 0, 3, true, true)
	y = utils.Byte2Int(bs1, 0, 3, true, true)

	if x != y {
		t.Fatalf(fmt.Sprintf("expected %d, got %d", x, y))
	}

}

func TestByte2Float(t *testing.T) {
	var x float32 = 0
	var y float32 = 0

	x = 15.668
	y = 0
	bs1 := make([]byte, 16)

	utils.Float2Byte(x, bs1, 0, true)
	y = utils.Byte2Float(bs1, 0, true)

	if x != y {
		t.Fatalf(fmt.Sprintf("expected %v, got %v", x, y))
	}

	var x2 float64 = -15.668
	var y2 float64 = 0

	utils.Double2Byte(x2, bs1, 0, true)
	y2 = utils.Byte2Double(bs1, 0, true)

	if x != y {
		t.Fatalf(fmt.Sprintf("expected %v, got %v", x2, y2))
	}

}

func TestByte2Bool(t *testing.T) {
	var x bool = true
	var y bool = false

	bs1 := make([]byte, 10)

	utils.Bool2Byte(x, bs1, 3, 5)
	y = utils.Byte2Bool(bs1, 3, 5)

	if x != y {
		t.Fatalf(fmt.Sprintf("expected %t, got %t", x, y))
	}
}

func TestByte2String(t *testing.T) {
	var x string = "hello"
	var y string = ""

	bs1 := make([]byte, 100)
	utils.String2Byte(x, bs1, 23, strings.Count(x, "")-1)
	y = utils.Byte2String(bs1, 23, 5)

	if x != y {
		t.Fatalf(fmt.Sprintf("expected %s, got %s", x, y))
	}

}
