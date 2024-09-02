package rtdb

import (
	"strconv"

	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/ulog"
)

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
	Data     []any
	Ts       int64
	Unit     string
}

func DPoint2AgilorDPoint(p *model.DPoint) *AgilorDPoint {
	return &AgilorDPoint{
		Id:       string(p.Id),
		Name:     p.Meta.Alias,
		DeviceId: string(p.ParentId),
		Type:     p.Meta.DataClass,
		Data:     p.Data,
		Ts:       p.Ts,
		Unit:     p.Meta.Unit,
	}
}

func AgilorDPoint2DPoint(ap *AgilorDPoint, meta *model.DataMeta) *model.DPoint {
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
		ParentId: parentId,
		Ts:       ap.Ts,
		Data:     ap.Data,
		Meta:     meta,
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

// type AlertEvent struct {
// 	Ts      int64
// 	Id      int64
// 	Level   int       //1,2,3...
// 	Snap    []*DPoint //数据快照
// 	Descrip string    //描述
// }

// type AlgTask struct {
// 	Ts     int64
// 	Id     int64
// 	Source []string
// 	Name   string
// }

// /**
//  * Algorithm result
//  */
// type AlgEvent struct {
// 	Ts     int64
// 	Id     int64
// 	TaskTs int64
// 	Name   string
// }

/**
 * AgilorCli for Agilor RTDB access
 */
type AgilorCli struct {
	Host     string
	Port     int
	Username string
	Passwd   string
	State    int
	c_cli    *C.agilor_cli //TODO: C的操作接口
}

func NewAgilorCli(host string, username string, password string) *AgilorCli {
	return &AgilorCli{
		Host:     host,
		Username: username,
		Passwd:   password,
	}
}

func (cli *AgilorCli) Open() {
}

func (cli *AgilorCli) Close() {
}

func (cli *AgilorCli) Insert(p *model.DPoint) {
	//转换为AgilorDPoint
	ap := DPoint2AgilorDPoint(p)
	cli.c_cli.insert(ap) //TODO: 这里调用C接口
}

/**
 * class, dnodeId, dpointId < 0 为无效参数
 */
func (cli *AgilorCli) Query(tic int64, toc int64, class int, dnodeId int64, dpointId int64, meta *model.DataMeta) []*model.DPoint {
	dpList := make([]*model.DPoint, 0)
	//TODO: 这里将查询到的AgilorDPoint转换为DPoint
	dp := &model.DPoint{
		Ts: 0,
	}
	dpList = append(dpList, dp)
	return dpList
}

func (cli *AgilorCli) QueryDNodeNow(id int64) *model.DPoint {
	return &model.DPoint{}
}

func (cli *AgilorCli) QueryDPointNow(id int64) *model.DPoint {
	return &model.DPoint{}
}

func (cli *AgilorCli) QueryDPointClassNow(class int) []*model.DPoint {
	return make([]*model.DPoint, 0)
}

/**
 * query aggregated data points
 * @param tic: start timestamp in us
 * @param toc: end timestamp in us
 * @param window: the window size model
 * @param step: the step size in us
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
 * @param class: 数据点位类型
 * class, dNodeId, dPointId 至少需要一个有效参数
 */
func (cli *AgilorCli) Delete(tic int64, toc int64, class int, dNodeId int64, dPointId int64) {
}

func (cli *AgilorCli) DropDNode(dNodeId int64) {
}

func (cli *AgilorCli) DropDPoint(dPointId int64) {
}

/**
 * 删除一类数据节点
 */
func (cli *AgilorCli) DropDNodeClass(dNodeClass int) {
}

/**
 * 删除指定一类数据点位
 */
func (cli *AgilorCli) DropDPointClass(dPointClass int) {

}

/**
 * 初始化调用，创建表
 */
func (cli *AgilorCli) CreateTable() {
}

/**
 * 运维用，删除表
 */
func (cli *AgilorCli) DropTable(name string) {

}

/****************************************************/
/* 以下为订阅相关接口，用于订阅数据变化事件 */
/****************************************************/
func (cli *AgilorCli) SubscribeDPoint(id int64) int {
	return 0
}

func (cli *AgilorCli) UnsubscribeDPoint(dPointId int64) int {
	return 0
}

func (cli *AgilorCli) SubscribeDpClass(class int) []int64 {
	//TODO: replace with dnode id
	return make([]int64, 0)
}

func (cli *AgilorCli) SubscribeDNode(dNodeId int64) int {
	return 0
}

func (cli *AgilorCli) UnsubscribeDNode(id int64) int {
	return 0
}
