package model

import (
	"github.com/lingfliu/ucs_core/conn"
	"github.com/lingfliu/ucs_core/model/gis"
)

const (
	//错误码
	ERR_CODE_OK       = 0 //正常
	ERR_CODE_OVERFLOW = 1 //溢出
	ERR_CODE_LOWPOWER = 2 //低电量
	ERR_CODE_DEAD     = 3 //死机
	ERR_CODE_DRIFT    = 4 //漂移
)

//deprecated
/**
 * 基础物模型：基于IO读写控制的节点模型
 * 1. 基本属性：id，归属，描述，
 * 2. 连接属性：地址，协议，网关
 * 3. 位置属性
 * 4. 基础属性：在线状态，错误码
 * 5. 子节点挂载：
 * 6. 数据点位
 */

// TODO: remove
type UNode struct {
	Id       int64 //数据库检索用
	ParentId int64
	Mac      string            //可选
	Name     string            //名称
	Class    string            //类型
	Descrip  string            //概要描述(型号，设备商等)
	PropSet  map[string]string //静态属性

	//连接属性
	Addr    string        //url地址
	ConnCfg *conn.ConnCfg //连接配置

	//位置，速度信息，当节点为装备、人员时为实时位置信息，否则为装配位置配置信息（静态）
	GPos *gis.GPos
	LPos *gis.LPos
	Velo *gis.Velo

	//状态属性
	Online  bool //在线状态
	ErrCode int  //错误码

	Children map[string]*UNode  //子节点, id as the key
	PointSet map[string]*UPoint //数据点位
}
