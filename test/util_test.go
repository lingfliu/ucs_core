package test

import (
	"fmt"
	"testing"

	"issc.io/mmitp/mmitp-core/utils"
)

func TestByte2Int(t *testing.T) {
	x := 600
	bs := make([]byte, 10)
	utils.Int2Byte(x, bs, 0, 5, true, true)
	y := utils.Byte2Int(bs, 0, 5, true, true)
	if x != y {
		t.Fatalf("expected %d, got %d", x, y)
	}

}

func TestAsciiStr2Hex(t *testing.T) {

	str := "hello Go"
	hex_str := utils.AsciiStr2Hex(str, " ")
	if hex_str != "68 65 6C 6C 6F 20 47 6F " {
		t.Fatalf("input %s, converted %s", str, hex_str)
	}
	// fmt.Println(hex_str)
}

func TestAsciiStr2Deci(t *testing.T) {

	str := "hello Go"
	deci_str := utils.AsciiStr2Deci(str, " ")
	if deci_str != "104 101 108 108 111 032 071 111 " {
		t.Fatalf("input %s, converted %s", str, deci_str)
	}
	fmt.Println(deci_str)
}
