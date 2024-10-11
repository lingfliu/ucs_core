package rtdb

/*
#cgo CFLAGS: -I./lib_agilor/src
#cgo LDFLAGS: -L$./lib_agilor_wrap -lagilor_wrap -Wl, -rpath=./lib_agilor_wrap
void agilor_ucs_pt_create(ucs_pt_t* p);
*/

import "C"
