package dao

import (
	"fmt"

	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/model/msg"
	"github.com/lingfliu/ucs_core/ulog"
)

/**
 * @brief
 * DPoint 数据点位CURD接口
 * 基于TDengine实现
 */
type DpDao struct {
	TaosCli *rtdb.TaosCli
}

func NewDpDao(host string, dbName string, username string, password string) *DpDao {
	return &DpDao{TaosCli: rtdb.NewTaosCli(host, dbName, username, password)}
}

func (dao *DpDao) Open() {
	dao.TaosCli.Open()
}

func (dao *DpDao) Init() {
	sql := fmt.Sprintf("create stable if not exists %s.dp (ts timestamp, v int) tags (dnode_class int, dnode_id int, dp_offset_idx int)", dao.TaosCli.DbName)
	res := dao.TaosCli.Exec(sql)
	if res < 0 {
		ulog.Log().E("dpdao", "failed to create stable dp")
	} else {
		ulog.Log().I("dpdao", "create stable dp success")
	}
}

func (dao *DpDao) Close() {
	dao.TaosCli.Close()
}

// TODO: 需要实现泛化，否则需要硬编码逐个数据结构进行实现
func (dao *DpDao) Insert(msg *msg.DMsg) {
	for idx, v := range msg.DataSet {
		tableName := fmt.Sprintf("dp_%d_%d", msg.DNodeId, idx)
		sql := fmt.Sprintf("insert into %s using dp tags(?,?,?) values (?, ?)", tableName)
		dao.TaosCli.Exec(sql, msg.Class, msg.DNodeId, idx, msg.Ts, v)
	}
}
