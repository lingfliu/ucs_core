package dao

import (
	"fmt"

	"github.com/lingfliu/ucs_core/coder"
	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/ulog"
)

/**
 * @brief
 * DPoint 数据点位CURD接口
 * 基于TDengine实现
 */
type DpDao struct {
	taosCli *rtdb.TaosCli
}

func NewDpDao(host string, dbName string, username string, password string) *DpDao {
	return &DpDao{taosCli: rtdb.NewTaosCli(host, dbName, username, password)}
}

func (dao *DpDao) Open() {
	dao.taosCli.Open()
}

func (dao *DpDao) Init() {
	sql := fmt.Sprintf("create stable if not exists %s.dp (ts timestamp, v int) tags (dnode_class int, dnode_id int, dp_offset_idx int)", dao.taosCli.DbName)
	res := dao.taosCli.Exec(sql)
	if res < 0 {
		ulog.Log().E("dpdao", "failed to create stable dp")
	} else {
		ulog.Log().I("dpdao", "create stable dp success")
	}
}

func (dao *DpDao) Close() {
	dao.taosCli.Close()
}

func (dao *DpDao) Insert(msg *coder.DpMsg) {
	for idx, v := range msg.ValueList {
		tableName := fmt.Sprintf("dp_%d_%d", msg.DNodeId, idx)
		sql := fmt.Sprintf("insert into %s using dp tags(?,?,?) values (?, ?)", tableName)
		dao.taosCli.Exec(sql, msg.DNodeClass, msg.DNodeId, idx, msg.Ts, v)
	}
}
