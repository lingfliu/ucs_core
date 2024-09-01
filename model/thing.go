package model

import (
	"github.com/lingfliu/ucs_core/conn"
	"github.com/lingfliu/ucs_core/model/gis"
)

/**
 * 物模型
 *  1. 基本属性：id，归属，描述，
 *  2. conn属性：地址，协议，网关
 *  3. 时空属性
 *  4. 基础属性：在线状态，错误码
 *  5. 设备属性：int，float, bool
 *  6. 数据-事件
 */

type BaseThing struct {
	Id   string
	Mac  string
	Name string //名称
	Desc string //描述

	//地址
	Url     string //ip or domain
	AddrMaj string //used for non-ip network
	AddrMin string //used for non-ip network
	Proto   string
	Conn    *conn.Conn //nil if direct connect to the server

	Gw ThingGw //nil if direct connect to the server

	//全局坐标
	GPos gis.GPos
	//速度
	Velo [3]float64 //速度：x,y,z unit: m/s

	//基础属性
	Online  bool //在线状态
	ErrCode int  //是否正常工作

	//states
	States map[string]*PropMeta //状态列表
	//events
	Events map[string]*PropMeta //事件名列表
	//data
	Datas map[string]*PropMeta //数据列表
}

type Thing interface {
}
