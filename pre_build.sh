#!/bin/sh 
echo "Pre-build Agilor wrapper library"
gcc ./data/rtdb/lib_agilor/src/agilor_wrap.c -fpic -shared -o ./data/rtdb/lib_agilor/libagilor_wrap.so
go build -x