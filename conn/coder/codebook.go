package coder

import (
	"encoding/json"
)

const DATA_CLASS_UINT = 1
const DATA_CLASS_INT = 2
const DATA_CLASS_FLOAT = 4 //IEEE 754 float32 or float64
const DATA_CLASS_BOOL = 5  //flags
const DATA_CLASS_STRING = 6

type CodeMeta struct {
	Name      string
	Offset    int
	DataClass int
	Dimen     int    //byte per element
	Size      int    //0 for variable length
	Vals      []byte //for fixed content
	AttrSize  bool   //if true, Size = 1, and should declare the related payload attr idx
	AttrIdx   int    //idx of the related variable payload, if =-1, it is the idx of the whole payload
}

type Code struct {
	Class     int
	ClassName string
	Metas     []*CodeMeta
	Payloads  []*CodeMeta
	MinSize   int
}

const CODE_PROC_RULE_EQUAL = 0 //if int / byte value matches
const CODE_PROC_RULE_MORE = 1  //if a int value is larger
const CODE_PROC_RULE_LESS = 2  //if a int value is less
const CODE_PROC_RULE_SET = 3   //if a flag value is set / unset

type CodeProcRule struct {
	Rule       int
	AttrName   string
	AttrOffset int
	AttrVal    int
	AttrSpan   int //slice span for given array, only effective for equal rule
}

type CodeProc struct {
	CodeName  string
	CodeClass int
	Side      int //0 for sender, 1 for receiver
	Timeout   int64
	Retry     int
	Rule      []*CodeProcRule
	Prev      *CodeProc
	Next      []*CodeProc
}

type Codebook struct {
	//commons for all codes
	Header  []byte
	MetaLen int
	Metas   []*CodeMeta
	HasErc  bool
	ErcLen  int
	Tail    []byte

	//code collection
	Codes map[string]*Code

	//TODO procedure and routine declaration
	//procedures
	Procs map[string]*CodeProc

	//routines (for cyclic procedures)
	Routines []string
}

func NewCodebookFromJson(jsonStr string) *Codebook {
	var cb Codebook
	err := json.Unmarshal([]byte(jsonStr), &cb)
	if err != nil {
		panic(err)
	}

	//TODO handle the err
	return &cb
}

func (cb *Codebook) Prepare() {
	for _, code := range cb.Codes {
		minSize := cb.CalcMinCodeSize(code.ClassName)
		code.MinSize = minSize
	}
}

// calculate the minimum code size for a given code
func (cb *Codebook) CalcMinCodeSize(name string) int {
	size := 0
	code := cb.Codes[name]
	for _, meta := range code.Metas {
		size += meta.Size
	}

	for _, meta := range code.Payloads {
		size += meta.Size
	}
	return size
}

func (cb *Codebook) GetCodeClassOffset() (int, *CodeMeta) {
	for _, meta := range cb.Metas {
		if meta.DataClass > 0 {
			return meta.Offset, meta
		}
	}
	return -1, nil
}

func (cb *Codebook) HasClass(class int) bool {
	for _, meta := range cb.Metas {
		if meta.DataClass == class {
			return true
		}
	}
	return false
}

func (cb *Codebook) Varlen(class int) bool {
	for _, meta := range cb.Metas {
		if meta.DataClass == class && meta.Size == 0 {
			return true
		}
	}
	return false
}

func (cb *Codebook) GetMetas(class int) []*CodeMeta {
	var metas []*CodeMeta
	for _, meta := range cb.Metas {
		if meta.DataClass == class {
			metas = append(metas, meta)
		}
	}
	return metas
}

func (cb *Codebook) GetPayloads(class int) []*CodeMeta {
	var metas []*CodeMeta
	for _, meta := range cb.Metas {
		if meta.DataClass == class {
			metas = append(metas, meta)
		}
	}
	return metas
}
