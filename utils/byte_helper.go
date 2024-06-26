package utils

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

/**
 * Convert a byte array to an integer
 * @param b byte array
 * @param offset start position
 * @l length
 * @param signed true if the integer is signed
 * @param msb true if the most significant byte is first
 */
const sysSize = strconv.IntSize / 8

func Byte2Int(bs []byte, offset int, l int, signed bool, msb bool) int {
	var x int = 0
	if msb {
		for i := 0; i < l; i++ {
			v := int(bs[offset+i] & 0xff)
			v <<= 8 * (l - i - 1)
			if signed {
				v <<= 8 * (sysSize - l)
				v >>= 8 * (sysSize - l)
			}
			x |= int(v)
		}
	} else {
		for i := l - 1; i >= 0; i-- {
			v := uint(bs[offset+i] & 0xff)
			v <<= 8 * (l - i - 1)
			if signed {
				v <<= 8 * (sysSize - l)
				v >>= 8 * (sysSize - l)
			}
			x |= int(v)
		}
	}
	return x
}

/**
 * Convert an integer to a byte array
 */
func Int2Byte(x int, bs []byte, offset int, l int, signed bool, msb bool) {
	if msb {
		for i := 0; i < l; i++ {
			bs[offset+i] = byte((x >> uint(8*(l-i-1))) & 0xFF)
		}
	} else {
		for i := l - 1; i >= 0; i-- {
			bs[offset+i] = byte((x >> uint(8*(l-i-1))) & 0xFF)
		}
	}
}

/*
*
* convert ascii string to hex string, zero padding by default
@param str ascii string
@param seperator default is space " "
*/
func AsciiStr2Hex(str string, seperator string) string {
	var hex_code = []byte("0123456789ABCDEF")

	if seperator == "" {
		seperator = " "
	}

	var hex []byte = make([]byte, len(str)*3)
	for i := 0; i < len(str); i += 1 {
		b := str[i]
		//TODO convert b to hex
		hex[i*3] = hex_code[(b&0xF0)>>4]

		hex[i*3+1] = hex_code[(b&0x0F)&0x0F]

		hex[i*3+2] = seperator[0]
	}
	return string(hex)
}

/*
*
* convert ascii string to decimal string, zero padding by default
@param str ascii string
@param seperator default is space " "
*/
func AsciiStr2Deci(str string, seperator string) string {
	var deci_code = []byte("0123456789")

	if seperator == "" {
		seperator = " "
	}

	var deci []byte = make([]byte, len(str)*4)
	for i := 0; i < len(str); i += 1 {
		b := str[i] & 0xff

		deci[i*4] = deci_code[b/100]
		deci[i*4+1] = deci_code[int((b/10)%10)]
		deci[i*4+2] = deci_code[int(b%10)]
		deci[i*4+3] = seperator[0]
	}
	return string(deci)

}

func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Float2Byte(x float32, bs []byte, offset int, msb bool) {
	bits := math.Float32bits(x)
	if msb {
		binary.BigEndian.PutUint32(bs[offset:], bits)
	} else {
		binary.LittleEndian.PutUint32(bs[offset:], bits)
	}
}

func Byte2Float(bs []byte, offset int, msb bool) float32 {
	bits := uint32(0)
	if msb {
		bits = binary.BigEndian.Uint32(bs[offset:])
	} else {
		bits = binary.LittleEndian.Uint32(bs[offset:])
	}
	return math.Float32frombits(bits)
}

func Double2Byte(x float64, bs []byte, offset int, msb bool) {
	bits := math.Float64bits(x)
	if msb {
		binary.BigEndian.PutUint64(bs[offset:], bits)
	} else {
		binary.LittleEndian.PutUint64(bs[offset:], bits)
	}
}

func Byte2Double(bs []byte, offset int, msb bool) float64 {
	bits := uint64(0)
	if msb {
		bits = binary.BigEndian.Uint64(bs[offset:])
	} else {
		bits = binary.LittleEndian.Uint64(bs[offset:])
	}
	return math.Float64frombits(bits)
}

/*
*
* insert bool onto a byte
@param blen bit position (0-7)
*/
func Bool2Byte(x bool, bs []byte, offset int, blen int) {
	if x {
		bs[offset] |= 1 << blen
	} else {
		bs[offset] &= 0 << blen
	}
}

func Byte2Bool(bs []byte, offset int, blen int) bool {
	x := bs[offset] & (1 << blen)
	if x > 0 {
		return true
	} else {
		return false
	}
}

func String2Byte(s string, bs []byte, offset int, strlen int) {
	l := strings.Count(s, "") - 1
	if l > strlen {
		//truncate string
		ss := s[:strlen]
		copy(bs[offset:], []byte(ss))
	} else {
		copy(bs[offset:], []byte(s))
	}
}

func Byte2String(bs []byte, offset int, strlen int) string {
	return string(bs[offset : offset+strlen])
}
