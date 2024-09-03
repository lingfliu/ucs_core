package model

import (
	"github.com/lingfliu/ucs_core/conn"
	"github.com/lingfliu/ucs_core/model/gis"
)

const (
	ERR_CODE_OK       = 0 //正常
	ERR_CODE_OVERFLOW = 1 //溢出
	ERR_CODE_LOWPOWER = 2 //低电量
	ERR_CODE_DEAD     = 3 //死机
	ERR_CODE_DRIFT    = 4 //漂移
)

/**
 * 基础物模型
 *  1. 基本属性：id，归属，描述，
 *  2. conn属性：地址，协议，网关
 *  3. 时空属性
 *  4. 基础属性：在线状态，错误码
 *  5. 设备属性：int，float, bool
 *  6. 数据-事件
 */

type Thing struct {
	Id       int64
	ParentId int64
	Mac      int64
	Name     string //名称
	Descrip  string //描述(型号，设备商)

	//连接属性 (直连、网关)
	ConnCfg   *conn.ConnCfg
	Conn      *conn.Conn
	Addr      string  //ip or url
	Gw        *GwNode //nil if direct connect to the server
	OffsetIdx int     //偏置地址索引

	//位置
	GPos *gis.GPos
	LPos *gis.LPos
	Velo *gis.Velo //速度：x,y,z unit: m/s

	//基础状态
	Online  bool //在线状态
	ErrCode int  //是否正常工作

	PropSet map[string]*NodeProp //静态属性
	DpSet   map[string]*DPoint   //数据/状态列表
	CpSet   map[string]*CtlPoint //控制列表
}
