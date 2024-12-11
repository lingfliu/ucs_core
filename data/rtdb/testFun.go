package main

/*
#cgo CFLAGS: -I./lib_agilor/include
#cgo LDFLAGS: -L./lib_agilor/lib -lagilor_wrap -Wl,-rpath=./lib_agilor/lib
#include "agilor.h"
#include "agilor_defs.h"
#include "agilor_wrap.h"  
#include<stdio.h>
#include<stdint.h>
#include<stdlib.h>
#include<stdbool.h>
#include<string.h>

*/
import "C"
import (
    "fmt"
    "time"
    "unsafe"
)
func main() {
    var pt C.ucs_pt_t
    C.strncpy(&pt.tag[0], C.CString("GoPoint_2"), C.size_t(len("GoPoint_2")))
    pt.tag[len("GoPoint_2")] = 0
    pt.v_type = 9
    now := time.Now()
    pt.ts = C.int64_t(now.Unix() * 1000)
println("Timestamp in milliseconds:", pt.ts)

    threadID := C.uint64_t(12345)
    throughFirewall:=false
    tf := C.agibool(throughFirewall)

    server := C.CString("Agilor")
    host_addr := C.CString("192.168.66.27")
    username := C.CString("test")
    password := C.CString("123")
    port := C.uint32_t(3955)

    // 启动 Agilor 客户端 连接数据库前必须调用Startup
    startupResult := C.c_Agcn_Startup(threadID, tf)
    if startupResult != 0 {
        fmt.Println("c_Agcn_Startup failed with error code:", startupResult)
        return
    }

    // 连接到服务器
    connectResult := C.c_Agcn_Connect(server, host_addr, username, password, port)
    if connectResult != 0 {
        fmt.Println("c_Agcn_Connect failed with error code:", connectResult)
        C.c_Agcn_Cleanup()
        return
    }

    //************ 创建点位函数测试************
   C.agilor_ucs_pt_create(&pt)


    // 断开连接
    disconnectResult := C.c_Agcn_Disconnect(server)
    if disconnectResult != 0 {
        fmt.Println("c_Agcn_Disconnect failed with error code:", disconnectResult)
    }

    // 清理
    cleanupResult := C.c_Agcn_Cleanup()
    if cleanupResult != 0 {
        fmt.Println("c_Agcn_Cleanup failed with error code:", cleanupResult)
    }

    // 释放 C 字符串
    C.free(unsafe.Pointer(server))
    C.free(unsafe.Pointer(host_addr))
    C.free(unsafe.Pointer(username))
    C.free(unsafe.Pointer(password))
}