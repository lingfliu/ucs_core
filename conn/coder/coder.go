package coder

// some reserved attr
const (
	CODE_CLASS_INT    = 0
	CODE_CLASS_FLOAT  = 1
	CODE_CLASS_BOOL   = 2
	CODE_CLASS_STRING = 3

	//other special class
	CODE_CLASS_DECI  = 4 //deimal format as ["1234"(int),"56"(frac)] in string format
	CODE_CLASS_ENUM  = 5 //enums are cross-byte uint values, not used
	CODE_CLASS_UNION = 6 //union of different types
)

type CodeAttrSpec struct {
	Name     string
	Class    int
	ByteLen  int // for bool values, is the bit position (0-7) in the byte
	Offset   int // in byte
	Size     int // >/= 0 for fixed length array, -1 for variable length
	Encoding string
	Unsigned bool   //for int only
	Msb      bool   //for numbers only
	LenSpec  string //identify the length field for variable length array, by default ""
	Reserved bool   //reserved field
}

type CodeMsgSpec struct {
	Name        string
	Class       int
	MetaList    []*CodeAttrSpec //msg spec meta
	PayloadList []*CodeAttrSpec
	Varlen      bool //if the msg has variable payload length
}

type Coder interface {
	Reset()
	Encode(msg *UMsg, bs []byte) int   //encode msg to byte, return length of the msg
	PushDecode(bs []byte, n int) *UMsg //push bs into the coder and try to decode from the ringbuffer return nil if decaode fails
	FastDecode(bs []byte) *UMsg        //fast decode without passing through ring buffer
}
