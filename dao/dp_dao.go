package dao

import (
	"fmt"
	"time"

	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/model"
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
		dao.TaosCli.Exec(sql, msg.Mode, msg.DNodeId, idx, msg.Ts, v)
	}
}

func (dao *DpDao) Query(tic string, toc string) []*model.DPoint {
	//convert date string to int64
	tic_time, _ := time.Parse("2006-01-02 15:04:05.000", tic)
	tic_nano := tic_time.UnixNano()
	toc_time, _ := time.Parse("2006-01-02 15:04:05.000", toc)
	toc_nano := toc_time.UnixNano()

	sql := fmt.Sprintf("select * from dp where ts >= %d and ts <= %d", tic_nano, toc_nano)
	rows := dao.TaosCli.Query(sql)
	if rows == nil {
		ulog.Log().E("dpdao", "failed to query dp")
	} else {
		for rows.Next() {
			//read data
			var ts int64
			var v int
			var dnodeClass int
			var dnodeId int
			var dpOffsetIdx int
			err := rows.Scan(&ts, &v, &dnodeClass, &dnodeId, &dpOffsetIdx)
			if err != nil {
				ulog.Log().E("dpdao", "failed to scan dp")
			} else {
				ulog.Log().I("dpdao", fmt.Sprintf("ts: %d, v: %d, dnode_class: %d, dnode_id: %d, dp_offset_idx: %d", ts, v, dnodeClass, dnodeId, dpOffsetIdx))
			}
		}
	}

	return make([]*model.DPoint, 0)
}
