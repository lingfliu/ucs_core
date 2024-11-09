package dao

import (
	"fmt"

	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/ulog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
 * @brief
 * DPoint 数据点位CURD接口
 * 基于TDengine实现
 */
type DNodeDao struct {
	Db *gorm.DB
}

func (dao *DNodeDao) Open() {
	db, err := gorm.Open(mysql.Open("user:password@tcp()"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dao.Db = db
}

func (dao *DNodeDao) Close() {
	//GORM 会自行管理连接池
}

func (dao *DNodeDao) Create() {
	db := dao.Db.Create(&model.DNode{})
	//create
	rows := db.RowsAffected
	ulog.Log().I("dnode_dao", fmt.Sprintf("create dnode success, rows affected: %d", rows))
}

func (dao *DNodeDao) Query(id int) *model.DNode {
	node := &model.DNode{}
	dao.Db.First(node, id)
	return node
}

func (dao *DNodeDao) QueryByClass(class int) []*model.DNode {
	nodes := make([]*model.DNode, 0)
	dao.Db.Where(nodes, "class = ?", class)
	return nodes
}

func (dao *DNodeDao) Insert(node *model.DNode) {
	dao.Db.Create(node)
}
