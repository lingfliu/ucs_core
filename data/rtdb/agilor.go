package rtdb

/*
#cgo CFLAGS: -I${SRCDIR}/data/rtdb/lib_agilor/src
//#cgo LDFLAGS: -L${SRCDIR}/data/rtdb/lib_agilor -lagilor_wrap -Wl,-rpath=${SRCDIR}/data/rtdb/lib_agilor
void agilor_ucs_pt_create(ucs_pt_t* p);
*/

import "C"

func CreatePt() {
	var p C.ucs_pt_t
	C.agilor_ucs_pt_create(&p)
}
