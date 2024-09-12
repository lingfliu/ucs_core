package spec

import (
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
)

/**
 * 温湿度传感器
 * 固定规格: 1个温度、1个湿度，均为4字节， 归属于一个数据点位
 */
func NewTehuNode(id int64, name string, addr string) *model.DNode {
	dmeta := &meta.DataMeta{
		ByteLen: 4,
		Dimen:   2,
		Alias:   "温湿度",
		Unit:    "℃/%",
	}
	node := &model.DNode{
		Id:        id,
		Name:      name,
		Addr:      addr,
		Desc:      "温湿度传感器",
		PropSet:   make(map[string]string),
		DPointSet: make(map[int64]*model.DPoint),
		State:     model.DNODE_STATE_OK,
	}

	point := model.NewDPoint(0, id, 0, 0, 0, dmeta, []byte{0, 0, 0, 0, 0, 0, 0, 0})

	node.DPointSet[0] = point

	return node
}
