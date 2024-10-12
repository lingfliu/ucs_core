package rtdb

/*
#cgo CFLAGS: -I./lib_agilor/src
#cgo LDFLAGS: -L./lib_agilor -lagilor_wrap -Wl,-rpath=./lib_agilor
#include<stdio.h>
#include<stdint.h>
#include<stdlib.h>
#include<string.h>
typedef struct {
	int64_t id;
	int64_t node_id;
    char tag[64];
    char descrip[32];
	void* pt_value;
	int64_t ts;
} ucs_pt_t;
void agilor_ucs_pt_create(ucs_pt_t* p);
void agilor_ucs_pt_insert(ucs_pt_t* p); //插入一个ucs点位值
int agilor_ucs_pt_query(char* tag, int64_t start_time, int64_t end_time, int64_t step, ucs_pt_t* pt_list); //查询范围内点位值
void agilor_ucs_pt_query_now(char* tag);
void agilor_ucs_pt_remove_before(char* tag, int64_t before_time);
*/
import "C"

//不能有空行

import (
	"fmt"
	"unsafe"

	"github.com/lingfliu/ucs_core/model"
)

/**
 * 调用创建表的C接口，
 * 这里仅将pt转换为C.ucs_pt_t, 再从c结构转换为agilor定义的结构
 * 存储在agilor数据库的内容仅包括：tag，name/descrip, value，其他全部放在mysql里面
 * agilor表的其他设置如单位，插值形式的设置，都在agilor_wrap.c中设置
 */
func CreatePtTable(pt *model.DPoint) {

	//拆表：
	dimen := pt.DataMeta.Dimen

	for i := 0; i < dimen; i++ {
		p := C.ucs_pt_t{}
		tag := fmt.Sprintf("%d-%d-%d", pt.Id, pt.NodeId, dimen) //标签
		p.id = C.int64_t(pt.Id)
		p.node_id = C.int64_t(pt.NodeId)

		cstr_tag := C.CString(tag)
		C.strncpy(&p.tag[0], cstr_tag, C.size_t(len(p.tag))-1)
		C.free(unsafe.Pointer(cstr_tag))

		cstr_name := C.CString(pt.Name)
		C.strncpy(&p.descrip[0], cstr_name, C.size_t(len(p.descrip)-1))
		C.free(unsafe.Pointer(cstr_name))

		C.agilor_ucs_pt_create(&p)

		// ulog.Log().I("agilor_c", "tag="+tag)
	}
}

func InserDData(pt *model.DPoint) {
	dimen := pt.DataMeta.Dimen
	sampleLen := pt.DataMeta.SampleLen

	//TODO: 这里根据一次的样本数量做多次插入, agilor本身不能存储除数值以外其他内容，所以如果出现idx不连续的，这里就不做插值补全了，直接计算ts后存储
	if pt.Mode == 0 {
		for i := 1; i < sampleLen; i++ {
			ts := pt.Ts + 1000000000/pt.Sps

			for i := 0; i < dimen; i++ {
				//拆成单点存储
				p := C.ucs_pt_t{}
				tag := fmt.Sprintf("%d-%d-%d", pt.Id, pt.NodeId, dimen) //标签
				p.id = C.int64_t(pt.Id)
				p.node_id = C.int64_t(pt.NodeId)

				cstr_tag := C.CString(tag)
				C.strncpy(&p.tag[0], cstr_tag, C.size_t(len(p.tag))-1)
				C.free(unsafe.Pointer(cstr_tag))

				cstr_name := C.CString(pt.Name)
				C.strncpy(&p.descrip[0], cstr_name, C.size_t(len(p.descrip)-1))
				C.free(unsafe.Pointer(cstr_name))

				p.ts = C.int64_t(ts)

				C.agilor_ucs_pt_insert(&p)

				// ulog.Log().I("agilor_c", "tag="+tag)
			}
		}
	}
}

func QueryData(pt model.DPoint, startTime int64, endTime int64, step int64) []*model.DPoint {
	dimen := pt.DataMeta.Dimen
	byteLen := pt.DataMeta.ByteLen

	listLen := 0
	dpList := make([]*model.DPoint, 0)

	for i := 0; i < dimen; i++ {
		//合并单点查询
		var p_list *C.ucs_pt_t
		p_ref := C.ucs_pt_t{}
		tag := fmt.Sprintf("%d-%d-%d", pt.Id, pt.NodeId, dimen) //标签
		cstr_tag := C.CString(tag)
		C.strncpy((*C.char)(unsafe.Pointer(&p_ref.tag[0])), cstr_tag, C.size_t(len(p_ref.tag))-1)
		C.free(unsafe.Pointer(cstr_tag))
		c_len := C.agilor_ucs_pt_query(&p_ref.tag[0], C.int64_t(startTime), C.int64_t(endTime), C.int64_t(step), p_list)

		listLen = int(c_len)

		if i == 0 {
			for k := 0; k < listLen; k++ {
				p := (*[1 << 30]C.ucs_pt_t)(unsafe.Pointer(p_list))[k]

				//首次构建DPointList
				dp := &model.DPoint{
					Id:       pt.Id,
					NodeId:   pt.NodeId,
					Name:     pt.Name,
					NodeAddr: pt.NodeAddr,
					Offset:   pt.Offset,
					Ts:       int64(p.ts),
					Idx:      -1,
					Session:  "",
					Mode:     pt.Mode,
					Sps:      pt.Sps,
					DataMeta: pt.DataMeta,
					Data:     make([]byte, dimen*byteLen),
				}

				//TODO: 将查询的数据映射回dp中
				// dp.Data[i*byteLen : (i+1)*byteLen] = p.pt_value

				dpList = append(dpList, dp)
			}
		} else {
			for k := 0; k < listLen; k++ {
				dp := dpList[k]
				//TODO: 将查询的数据映射回dp中
				dp.Data[i] = 0 //TODO: remove this line
				// dp.Data[i*byteLen : (i+1)*byteLen] = p.pt_value
			}
		}
	}

	return dpList
}
