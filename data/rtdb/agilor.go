package rtdb

/*
#cgo CFLAGS: -I./lib_agilor/include
#cgo LDFLAGS: -L./lib_agilor -lagilor -Wl,-rpath=./lib_agilor
#include "./lib_agilor/src/agilor_wrap.h"
void* agilor_create(const char* name);
*/
import "C"

import (
	"strconv"

	"github.com/lingfliu/ucs_core/dd"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/ulog"
)

/****************************************************/
// C 结构体声明
/****************************************************/
// 定义agibool类型，假设这是一个布尔类型
type AgiBool bool

// 定义枚举类型
type AgiEnumValue struct {
	Class int16     // 0x0001：使用key, 0x0002：使用name，0x0003表示同时使用key,name
	Key   int16     // 枚举(值)
	Name  [128]byte // 枚举(字符串)
}

// 模拟union联合体类型
// TODO 用byte替代
type AgiValueUnion struct {
	RVal float32      // 浮点
	LVal int32        // 长整
	BVal AgiBool      // 开关
	SVal [128]byte    // sval  字符串
	EVal AgiEnumValue // 枚举
}

// 定义agilor_point_t结构体
type AgiPoint struct {
	Tag          [64]byte      // 测点标签 *
	Descriptor   [32]byte      // 测点描述 #
	Engunit      [16]byte      // 测点数据单位（安培、摄氏度等） #
	Id           int32         // 测点编号，由系统配置
	Class        uint8         // 菜单类型(R浮点数/S字符串/B开关/L整形/E枚举) * //TODO: 这里缺了常量声明
	Scan         uint8         // 测点扫描标识(0或>=0x80："禁止"， 1："输入", 2："输出" *
	TypicalValue float32       // 典型值 # @unused
	ValueUnion   AgiValueUnion // 点值 #
	EnumDesc     [128]byte     // 枚举描述 （"2:1,2,on:0,3,off"），暂时无用，[hp has not] @unused
	TimeDate     int64         // 时间戳 (ts)
	State        int32         // 点状态（点的质量、实时点值、缓冲的点值）
	// 由系统配置，覆盖添加时=old.state
	PointSource [32]byte  // 测点的数据源站(设备名) *
	SourceGroup [32]byte  // 测点的数据源结点组 #
	SourceTag   [128]byte // 测点的源标签 *

	UpperLimit float32 // 数据上限，用于压缩
	LowerLimit float32 // 数据下限，用于压缩

	PushRef1 uint16 // 实时推理规则标志 #
	RuleRef1 uint16 // 实时推理规则标志 #

	// 异常报告可确保Agilor接口只发送有意义的数据，而不是发送不必要的数据，从而加重系统的负担。
	// Exception reporting uses a simple deadband algorithm to determine whether
	// to send events to Agilor Data Archive. For each point, you can set
	// exception reporting specifications that create the deadband. The interface
	// ignores values that fall inside the deadband.
	// TODO: exc_xxx这3个参数，只对接口有效，还是内核中也使用这3个参数？
	ExcMin int64 // 实时数据处理最短间隔（接口参数）
	ExcMax int64 // 实时数据处理最大间隔（内核参数）
	// 不管是否压缩，点值是否变化，当timedate-last_timedate >= exc_max时强制存储数据
	ExcDev float32 // 实时数据处理偏差（接口参数）：
	// 当fabs(tagvalue.rval) < fabs(lptag->rval) * (1 - lptag->exc_dev)
	// 或(fabs(tagvalue.rval) > fabs(lptag->rval) * (1 - lptag->exc_dev)
	// 表示点值变化超过偏差。这时当点值变化超过偏差且与上次发送的点值时间戳之差>=exc_min
	// 时，即使是过滤发生，也会将点值发送到内核。

	/**************************************************/
	//N.B. 关于报警部分这里全部不使用
	/**************************************************/
	AlarmType  uint16  // 报警类型
	AlarmState uint16  // 状态报警
	AlarmHi    float32 // 上限报警
	AlarmLo    float32 // 下限报警
	AlarmHiHi  float32
	AlarmLoLo  float32

	PriorityHi   uint16 // 报警优先级，暂时不处理
	PriorityLo   uint16
	PriorityHiHi uint16
	PriorityLoLo uint16

	/**************************************************/
	//存储部分设置
	/**************************************************/
	Archive  AgiBool // 是否存储历史数据
	Compress AgiBool // 是否进行历史压缩 *，但type=O时，compress=agifalse
	Step     uint8   // 历史数据的插值形式（线形，台阶），compress=agitrue时有效
	HisIdx   int32   // 历史记录索引号，由系统配置

	// Compression testing
	CompMin int64   // 压缩最短间隔(压缩最小时间), compress minimum time，暂时无用
	CompMax int64   // 压缩最长间隔(压缩最大时间), compress maximum time，暂时无用
	CompDev float32 // 压缩灵敏度（压缩偏差）， compress deviation
	// 归档时压缩灵敏度=(upper_limit - lower_limit) * comp_dev

	LastVal      float32 // 上次数据存档的值 #  // TODO: 应该使用内存中的？
	LastTimeDate int64   // 上次数据存档的时间 # // TODO: 应该使用内存中的？
	CreateDate   int64   // 采集点创建日期，由系统配置
}

///////////////////////////////////////
///////////agilor_value_t////////////
//////////////////////////////////////

// 定义agilor_value_t结构体
type AgiValue struct {
	TimeDate int64 // 时间戳
	State    int32 // 状态 (Agpt_SetPointValue不需要设置state)
	Class    uint8 // 点值类型
	BlobSize int32
	Value    AgiValueUnion // 点值联合体
}

// 定义agilor_deviceinfo_t结构体
type AgiDeviceInfo struct {
	DeviceName [32]byte // 设备名称
	IsOnline   AgiBool  // 是否在线
	PointCount int32    // 测量点数量
}

// 定义agilor_devicepoint_t结构体
type AgiDevicePoint struct {
	LocalId   int32     // 本地重新分配的测点id
	Id        int32     // 测点id
	SourceTag [128]byte // 测点的源标签
	ExcDev    float32
	ExcMin    int64
	ExcMax    int64
	Class     uint16
	Scan      uint16
	TimeDate  int64
	State     int32
	Value     AgiValueUnion // 联合体
}

// C 函数接口

func AgiCreate(p *AgiDevicePoint) {

	// namePtr := C.CString(string(p.SourceTag[:]))

	// // fmt.Println(namePtr)

	// ptr := unsafe.Pointer(C.agilor_create(namePtr))

	// fmt.Println(ptr)

	// if ptr == nil {
	// 	ulog.Log().E("agilor", "create failed")
	// }
}

// func (AgiCreate)

/*************************************************/
/* 缩写声明：DNode*可缩写为Dn*， DPoint可缩写为Dp* */
/************************************************/

type AgilorDevice struct {
	Id         string
	Name       string
	DataPoints []*AgilorDPoint
}

/**
 * AgilorDPoint for Agilor RTDB data point
 * TODO: 按数据库定义补充字段
 */
type AgilorDPoint struct {
	Id       string
	Name     string
	DeviceId string
	Type     int
	Data     []byte
	Ts       int64
	Unit     string
}

func DPoint2AgilorDPoint(p *model.DPoint) *AgilorDPoint {
	return &AgilorDPoint{
		Id:       string(p.Id),
		Name:     p.DataMeta.Alias,
		DeviceId: string(p.NodeId),
		Type:     p.DataMeta.DataClass,
		Data:     p.Data,
		Ts:       p.Ts,
		Unit:     p.DataMeta.Unit,
	}
}

func AgilorDPoint2DPoint(ap *AgilorDPoint, meta *meta.DataMeta) *model.DPoint {
	var err error
	var id int64
	id, err = strconv.ParseInt(ap.Id, 10, 64)
	if err != nil {
		ulog.Log().E("agilor", "id parse error")
		return nil
	}

	var parentId int64
	parentId, err = strconv.ParseInt(ap.DeviceId, 10, 64)
	if err != nil {
		ulog.Log().E("agilor", "parentId parse error")
		return nil
	}

	return &model.DPoint{
		Id:       id,
		NodeId:   parentId,
		Ts:       ap.Ts,
		Data:     ap.Data,
		DataMeta: meta,
	}
}

const (
	AGGR_OP_MIN  = 0
	AGGR_OP_MAX  = 1
	AGGR_OP_MEAN = 2
	AGGR_OP_MED  = 3
	AGGR_OP_STD  = 4
	AGGR_OP_SUM  = 5
)

/**
 * AgilorCli for Agilor RTDB access
 */
type AgilorCli struct {
	Host     string //ip:port
	Ip       string
	Port     int
	Username string
	Passwd   string
	State    int
	// c_cli    *C.agilor_cli //TODO: C的操作接口
}

func NewAgilorCli(host string, username string, password string) *AgilorCli {
	return &AgilorCli{
		Host:     host,
		Username: username,
		Passwd:   password,
	}
}

/**
 * 连接数据库
 */
func (cli *AgilorCli) Open() {
}

/**
 * 关闭数据库
 */
func (cli *AgilorCli) Close() {
}

/**
 * 拆分创建数据点位
 */
func (cli *AgilorCli) CreateDPoint(p *model.DPoint) {
	for i := 0; i < p.DataMeta.Dimen; i++ {
		//TODO: 创建nodeId-dpointId-i标签的
		//createAgiPoint()
	}
}

/**
 * 删除数据点位
 */
func (cli *AgilorCli) DeleteDPoint(pointId int64) {
	//TODO: 调用maskTag接口找到对应的表并删除
}

/**
 * 创建设备节点对应的点位
 */
func (cli *AgilorCli) CreateDNode(node *model.DNode) {
	for _, p := range node.DPointSet {
		cli.CreateDPoint(p)
	}
}

/**
 * 删除设备节点
 */
func (cli *AgilorCli) DeleteDNode(nodeId int64) {
	//TODO: 调用maskTag接口找到对应的表并删除
}

/**
 * 查询数据节点, id=0查询所有
 */
func (cli *AgilorCli) QueryDNode(id int64) []*model.DNode {
	//TODO:
	return make([]*model.DNode, 0)
}

/**
 * 插入数据
 */
func (cli *AgilorCli) Insert(p *model.DPoint) {
	//转换为AgilorDPoint
	// ap := DPoint2AgilorDPoint(p)
	// cli.c_cli.insert(ap) //TODO: 这里调用C接口

	for i := 0; i < p.DataMeta.Dimen; i++ {
		//TODO: 调用C接口插入一个数据
		// tag := fmt.Sprintf("%d-%d-%d", p.NodeId, p.Id, i)
		// AgiInsert(tag, p.Data[i], p.Ts)
	}
}

/**
 * 查询单个历史点位数据
 * class, dpointId < 0 为无效参数
 */
func (cli *AgilorCli) QueryDPoint(tic int64, toc int64, dpointId int64, meta *meta.DataMeta) []*model.DPoint {
	dpList := make([]*model.DPoint, 0)
	//TODO: 这里将查询到的AgilorDPoint转换为DPoint
	return dpList
}

/**
 * 查询指定节点所有的数据快照
 */
func (cli *AgilorCli) QueryDNodeNow(id int64, class int) []*model.DPoint {
	//TODO: 调用query与转换接口
	return make([]*model.DPoint, 0)
}

/**
 * 查询指定数据点位（类型）数据快照
 */
func (cli *AgilorCli) QueryDPointNow(id int64, class int) []*model.DPoint {
	//TODO: 调用query与转化接口
	points := make([]*model.DPoint, 0)
	return points
}

/**
 * query aggregated data points
 * @param tic: start timestamp in ns
 * @param toc: end timestamp in us
 * @param window: the window size in ns
 * @param step: the step size in ns
 * @param class: the class of data points
 */
func (cli *AgilorCli) AggregateQuery(tic int64, toc int64, window int64, step int64, class int, op int) []model.DPoint {
	//TODO: 调用接口，如接口不支持，则需要实现手动计算
	switch op {
	case AGGR_OP_MIN:
	case AGGR_OP_MAX:
	case AGGR_OP_MEAN:
	case AGGR_OP_MED:
	case AGGR_OP_STD:
	case AGGR_OP_SUM:
	}

	return make([]model.DPoint, 0)
}

/**
 * 删除数据
 * @param tic: 开始时间 //TODO: 一般设置为0，优先实现删除toc时间之前所有历史数据
 * @param toc: 结束时间
 * dNodeId, dPointId 至少需要一个有效参数
 */
func (cli *AgilorCli) Delete(tic int64, toc int64, dNodeId int64, dPointId int64) {
}

/*
*
* 删除节点(类型）以及其下所有的数据点位数据
* 不单独删除数据点位因为点位和节点是强绑定关系
@param dNodeId: 节点id
@param dNodeClass: 节点类型
以上两个参数至少需要一个有效
*/
func (cli *AgilorCli) DropDNode(dNodeId int64, dNodeClass int) {
	//TODO: 调用接口删除
}

/****************************************************/
/* 以下为订阅相关接口，用于订阅数据变化事件 */
/****************************************************/
func (cli *AgilorCli) SubscribeDPoint(id int64, callback func(*dd.DdZeroMsg)) int {
	return 0
}

func (cli *AgilorCli) UnsubscribeDPoint(dPointId int64) int {
	return 0
}

func (cli *AgilorCli) SubscribeDNode(dNodeId int64) int {
	return 0
}

func (cli *AgilorCli) UnsubscribeDNode(id int64) int {
	return 0
}
