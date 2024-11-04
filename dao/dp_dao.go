package dao

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
	"github.com/lingfliu/ucs_core/model/msg"
	"github.com/lingfliu/ucs_core/ulog"
)

/**
 * @brief
 * DPoint 数据点位CRUD接口
 * 基于TDengine实现
 * 超表命名规则: dp_{node_class}
 * 标签: dnode_id bigint | dpoint_id bigint | pos varchar
 */
type DpDao struct {
	Template *model.DPoint
	TaosCli  *rtdb.TaosCli
}

func NewDpDao(host string, dbName string, username string, password string, template model.DPoint) *DpDao {
	return &DpDao{TaosCli: rtdb.NewTaosCli(host, dbName, username, password)}
}

func (dao *DpDao) Open() {
	dao.TaosCli.Open()
}

func (dao *DpDao) TableExist(tableName string) bool {
	sql := fmt.Sprintf("show tables like '%s'", tableName)
	rows := dao.TaosCli.Query(sql)
	if rows == nil {
		ulog.Log().E("dpdao", "no table found")
		return false
	} else {
		defer rows.Close()
		if rows.Next() {
			return true
		}
	}
	return false
}

func (dao *DpDao) InitTable(template *model.DPoint) int {

	var valueClass string
	if template.DataMeta.DataClass == meta.DATA_CLASS_INT32 {
		valueClass = "int"
	} else if template.DataMeta.DataClass == meta.DATA_CLASS_FLOAT {
		valueClass = "float"
	} else {
		ulog.Log().E("dpdao", "unsupported data class")
		return -1
	}

	dimen := template.DataMeta.Dimen

	colStr := "(ts timestamp"

	for i := 0; i < dimen; i++ {
		colStr += fmt.Sprintf(", v%d %s", i, valueClass)
	}
	colStr += ")"

	stableName := fmt.Sprintf("dp")

	sql := fmt.Sprintf("create stable if not exists %s %s tags (dnode_name nchar(32), dpoint_name nchar(32), dnode_id int, dp_offset int, dpoint_unit binary(16))", stableName, colStr)

	res := dao.TaosCli.Exec(sql)
	if res < 0 {
		ulog.Log().E("dpdao", fmt.Sprintf("failed to create stable %s", stableName))
	} else {
		ulog.Log().I("dpdao", fmt.Sprintf("create stable %s success", stableName))
	}
	return res
}

func (dao *DpDao) Close() {
	dao.TaosCli.Close()
}

func (dao *DpDao) Insert(dmsg *msg.DMsg) {
	for idx, v := range dmsg.DataSet {
		tableName := fmt.Sprintf("dp_%d_%d", dmsg.DNodeId, idx)

		dimen := v.Meta.Dimen
		colStr := "(?"
		i := 0
		for i < dimen {
			colStr += ", ?"
			i++
		}
		colStr += ")"

		valueList, tsList := v.AsInt32(dmsg.Ts, dmsg.Sps, true)

		for idx, ts := range tsList {
			values := valueList[idx]
			sql := fmt.Sprintf("insert into %s using dp tags(?,?,?) values %s", tableName, colStr)
			anyValues := make([]any, len(values)+4)
			anyValues[0] = dmsg.Mode
			anyValues[1] = dmsg.DNodeId
			anyValues[2] = dmsg.Offset
			anyValues[3] = ts
			for i, v := range values {
				anyValues[i+4] = v
			}
			dao.TaosCli.Exec(sql, anyValues...)
		}
	}
}

func (dao *DpDao) Query(tic string, toc string, dnodeId int64, offset int, dataMeta *meta.DataMeta) []*model.DPoint {

	dpList := make([]*model.DPoint, 0)
	//convert date string to int64
	tic_time, _ := time.Parse("2006-01-02 15:04:05.000", tic)
	tic_ms := tic_time.UnixNano() / 1000000
	toc_time, _ := time.Parse("2006-01-02 15:04:05.000", toc)
	toc_ms := toc_time.UnixNano() / 1000000

	// tableName := fmt.Sprintf("dp_%d_%d", dnodeId, offset)
	sql := fmt.Sprintf("select * from %s where ts between %d and %d", "dp", tic_ms, toc_ms)
	rows := dao.TaosCli.Query(sql)
	if rows == nil {
		ulog.Log().E("dpdao", "failed to query dp")
	} else {
		defer rows.Close()
		for rows.Next() {
			//read data
			var ts string
			var dnodeClass int
			var dnodeId int64
			var dpOffset int32
			scanned := make([]any, dataMeta.Dimen+4)
			values := make([]int, 4)

			scanned[0] = &ts
			i := 0
			for i < dataMeta.Dimen {
				scanned[i+1] = &values[i]
				i++
			}
			scanned[5] = &dnodeClass
			scanned[6] = &dnodeId
			scanned[7] = &dpOffset

			err := rows.Scan(scanned...)

			if err != nil {
				ulog.Log().E("dpdao", "failed to scan dp")
			} else {
				t, _ := time.Parse("2006-01-02T15:04:05.000+08:00", ts)
				t_ms := t.UnixNano()
				// ulog.Log().I("dpdao", fmt.Sprintf("ts: %d, v: %d, dnode_class: %d, dnode_id: %d, dp_offset_idx: %d", t.UnixNano()/1000000, values[0], dnodeClass, dnodeId, dpOffsetIdx))

				dp := &model.DPoint{
					NodeId:   dnodeId,
					Offset:   offset,
					Ts:       t_ms,
					DataMeta: dataMeta,
					Data:     make([]byte, dataMeta.Dimen*dataMeta.ByteLen),
				}

				i := 0
				for i < dataMeta.Dimen {
					if dataMeta.ByteLen == 2 {
						binary.BigEndian.PutUint16(dp.Data[i*dataMeta.ByteLen:(i+1)*dataMeta.ByteLen], uint16(values[i]))
					} else if dataMeta.ByteLen == 4 {
						binary.BigEndian.PutUint32(dp.Data[i*dataMeta.ByteLen:(i+1)*dataMeta.ByteLen], uint32(values[i]))
					} else if dataMeta.ByteLen == 8 {
						binary.BigEndian.PutUint64(dp.Data[i*dataMeta.ByteLen:(i+1)*dataMeta.ByteLen], uint64(values[i]))
					}
					i++
				}

				dpList = append(dpList, dp)

			}
		}
	}

	return dpList
}
