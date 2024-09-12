package impl

import (
	"encoding/binary"
	"fmt"

	"github.com/lingfliu/ucs_core/dao"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/ulog"
)

type TehuNodeDao struct {
	dao.DpDao
}

const (
	stableName = "th_node"
)

func (dao *TehuNodeDao) Create() int {

	sql := fmt.Sprintf("create stable if not exists %s.dp (ts timestamp, temp int, humi int) tags (dnode_id int, dp_offset int, )", stableName)
	res := dao.TaosCli.Exec(sql)
	if res < 0 {
		ulog.Log().E("dpdao", "failed to create stable dp")
	} else {
		ulog.Log().I("dpdao", "create stable dp success")
	}
	return res
}

func (dao *TehuNodeDao) Insert(p *model.DPoint) {
	//子表命名方式 ${stablename}_${nodeid}_${dp_offset}
	temp := binary.BigEndian.Uint32(p.Data[:3])
	humi := binary.BigEndian.Uint32(p.Data[4:])
	tableName := fmt.Sprintf("%s_%d_%d", stableName, p.NodeId, p.Offset)
	sql := fmt.Sprintf("insert into %s using %s values(?, ?, ?) tags(?, ?)", tableName, stableName)
	dao.TaosCli.Exec(sql, p.Ts, int(temp), int(humi), p.NodeId, p.Offset)
}

func (dao *TehuNodeDao) Query(nodeId int, tic int64, toc int64) []*model.DPoint {
	sql := fmt.Sprintf("select * from %s.dp where dnode_id=%d and ts>=%d and ts<=%d", stableName, nodeId, tic, toc)
	res := dao.TaosCli.Query(sql)
	if res == nil {
		return nil
	}
	defer res.Close()
	var ret []*model.DPoint
	for res.Next() {
		p := &model.DPoint{}
		res.Scan(&p.Ts, &p.Data, &p.NodeId, &p.Offset)
		ret = append(ret, p)
	}
	return ret
}

/**
 * 聚合查询
 * @param nodeId
 * @param tic
 * @param toc
 * @param op: 操作符： avg, std, sum, max, min, count
 */
func (dao *TehuNodeDao) AggrQuery(nodeId int, tic int64, toc int64, op int, window int, step int) []*model.DPoint {
	sql := fmt.Sprintf("select %s(temp), %s(humi) from %s.dp where dnode_id=? and ts>=? and ts<=? group by dnode_id, dp_offset/10", op, op, stableName, nodeId, tic, toc)
	res := dao.TaosCli.Query(sql)
	if res == nil {
		return nil
	}
	defer res.Close()
	var ret []*model.DPoint
	for res.Next() {
		p := &model.DPoint{}
		res.Scan(&p.Ts, &p.Data, &p.NodeId, &p.Offset)
		ret = append(ret, p)
	}
	return ret
}

// TODO: 删除t时间以前所有数据
func (dao *TehuNodeDao) DeleteBefore(nodeId int64, t int64) {

}

// TODO: 删除子表数据
func (dao *TehuNodeDao) Drop(nodeId int64, offset int) {
	// sql := fmt.Sprintf("drop ")
}
