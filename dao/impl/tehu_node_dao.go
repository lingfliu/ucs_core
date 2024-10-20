package impl

import (
	"encoding/binary"
	"fmt"

	"github.com/lingfliu/ucs_core/dao"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/ulog"
)

type TehuNodeDao struct {
	dao.DpDao
}

const (
	stableName = "tehu_node"
)

func (dao *TehuNodeDao) GenerateTemplate() *model.DNode {
	//generate a template for tehunode
	tehuNode := &model.DNode{
		Id:        0,
		Name:      stableName,
		Addr:      "",
		Mode:      model.DNODE_MODE_AUTO,
		DPointSet: make(map[int64]*model.DPoint),
	}

	//temperate dpoint
	tehuNode.DPointSet[0] = &model.DPoint{
		Id:     0,
		Name:   "温度",
		NodeId: 0,
		Offset: 0,

		Ts:      0,
		Idx:     0,
		Session: "",
		DataMeta: &meta.DataMeta{
			DataClass: meta.DATA_CLASS_INT,
			ByteLen:   4,
			Dimen:     1,
			SampleLen: 1,
			Alias:     "",
		},
		Data: make([]byte, 4),
	}

	//humidity dpoint
	tehuNode.DPointSet[1] = &model.DPoint{
		Id:     0,
		Name:   "湿度",
		NodeId: 0,
		Offset: 1,

		Ts:      0,
		Idx:     0,
		Session: "",
		DataMeta: &meta.DataMeta{
			DataClass: meta.DATA_CLASS_INT,
			ByteLen:   4,
			Dimen:     1,
			SampleLen: 1,
			Alias:     "",
		},
		Data: make([]byte, 4),
	}
	return tehuNode
}

func (d *TehuNodeDao) TableExist() bool {
	for _, v := range d.GenerateTemplate().DPointSet {
		if !d.DpDao.TableExist(fmt.Sprintf("%s_%d_%d", stableName, v.NodeId, v.Offset)) {
			return false
		}
	}
	return true
}

// 初始化表
func (d *TehuNodeDao) InitTable() int {

	dd := dao.DpDao{}

	for _, v := range d.GenerateTemplate().DPointSet {
		//TODO: return error
		res := dd.InitTable(v)
		if res < 0 {
			ulog.Log().E("dpdao", "failed to init table")
			return -1
		}
	}
	ulog.Log().I("dpdao", "init table success")
	return 0
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
