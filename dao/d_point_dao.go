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
type DPointDao struct {
	Db *gorm.DB
}

func (dao *DPointDao) Open() {
	db, err := gorm.Open(mysql.Open("user:password@tcp()"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dao.Db = db
}

func (dao *DPointDao) Close() {
	//GORM 会自行管理连接池
}

func (dao *DPointDao) Create() {
	db := dao.Db.Create(&model.DNode{})
	//create
	rows := db.RowsAffected
	ulog.Log().I("dnode_dao", fmt.Sprintf("create dnode success, rows affected: %d", rows))
}

func (dao *DPointDao) Query(id int) *model.DPoint {
	node := &model.DPoint{}
	dao.Db.First(node, id)
	return node
}

func (dao *DPointDao) QueryByClass(class int) []*model.DPoint {
	nodes := make([]*model.DPoint, 0)
	dao.Db.Where(nodes, "class = ?", class)
	return nodes
}

func (dao *DPointDao) Insert(node *model.DPoint) {
	dao.Db.Create(node)
}
