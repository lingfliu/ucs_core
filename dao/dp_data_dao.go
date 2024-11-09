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
 * DpData 数据点位CRUD接口
 * 基于TDengine实现
 * 超表命名规则: dp_{node_class}_{dp_class}
 * 标签: dnode_id bigint | dpoint_id bigint | pos varchar
 */
type DpDataDao struct {
	ColNameList []string
	Template    *model.DPoint
	TaosCli     *rtdb.TaosCli
}

func NewDpDataDao(host string, dbName string, username string, password string) *DpDao {
	return &DpDao{TaosCli: rtdb.NewTaosCli(host, dbName, username, password)}
}

func (dao *DpDataDao) Open() {
	dao.TaosCli.Open()
}

func (dao *DpDataDao) TableExist(tableName string) bool {
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

func (dao *DpDataDao) CreateTableFromTemplate(template *model.DNodeTemplate) int {

	node := template.Template
	for _, pt := range node.DPointList {
		stableName := fmt.Sprintf("dp_%s_%s_%d", node.Class, pt.Class, pt.Offset)
		colStr := "(ts timestamp"
		//create type mapping
		dataClass := pt.DataMeta.DataClass
		tdType := "int"
		if dataClass == meta.DATA_CLASS_INT8 {
			tdType = "tinyint"
		} else if dataClass == meta.DATA_CLASS_UINT8 {
			tdType = "tinyint unsigned"
		} else if dataClass == meta.DATA_CLASS_INT16 {
			tdType = "smallint"
		} else if dataClass == meta.DATA_CLASS_UINT16 {
			tdType = "smallint unsigned"
		} else if dataClass == meta.DATA_CLASS_INT32 {
			tdType = "int"
		} else if dataClass == meta.DATA_CLASS_UINT32 {
			tdType = "int unsigned"
		} else if dataClass == meta.DATA_CLASS_INT64 {
			tdType = "bigint"
		} else if dataClass == meta.DATA_CLASS_INT64 {
			tdType = "bigint unsigned"
		} else if dataClass == meta.DATA_CLASS_FLOAT {
			tdType = "float"
		} else if dataClass == meta.DATA_CLASS_DOUBLE {
			tdType = "double"
		} else if dataClass == meta.DATA_CLASS_FLAG {
			tdType = "bool"
		} else {
			ulog.Log().E("dpdao", "unsupported data class")
			return -1
		}

		for i := 0; i < pt.DataMeta.Dimen; i++ {
			if template.ColNameList != nil && len(template.ColNameList) > i && template.ColNameList[i] != "" {
				//a valid column name
				colStr += fmt.Sprintf(", %s %s", template.ColNameList[i], tdType)
			} else {
				colStr += fmt.Sprintf(", v%d %s", i, tdType)
			}
		}

		sql := fmt.Sprintf("create stable if not exists %s (%s) tags (dnode_name nchar(32), dpoint_name nchar(32), dnode_id int, dpoint_offset int, dpoint_unit binary(16))", stableName, colStr)

		res := dao.TaosCli.Exec(sql)
		if res < 0 {
			ulog.Log().E("dpdao", fmt.Sprintf("failed to create stable %s", stableName))
			return -1
		} else {
			ulog.Log().I("dpdao", fmt.Sprintf("create stable %s success", stableName))
		}
	}

	return 0
}

func (dao *DpDataDao) Close() {
	dao.TaosCli.Close()
}

func (dao *DpDataDao) Insert(dmsg *msg.DMsg) int {
	tsList := make([]int64, dmsg.SampleLen)
	for i := 0; i < dmsg.SampleLen; i++ {
		tsList[i] = dmsg.Ts + int64(i)*dmsg.Sps
	}

	for idx, v := range dmsg.DataList {
		tableName := fmt.Sprintf("dp_%d_%d", dmsg.DNodeId, idx)

		dimen := v.Meta.Dimen
		colStr := "(?"
		i := 0
		for i < dimen {
			colStr += ", ?"
			i++
		}
		colStr += ")"

		valueList := make([][]any, dmsg.SampleLen)
		for i := 0; i < dmsg.SampleLen; i++ {
			for j := 0; j < dmsg.DataList[i].Meta.Dimen; j++ {
				if v.Meta.DataClass == meta.DATA_CLASS_INT16 {
					if dmsg.DataList[i].Meta.Msb {
						valueList[i][j] = int16(binary.BigEndian.Uint16(dmsg.DataList[i].Data[j*2 : (j+1)*2]))
					} else {
						valueList[i][j] = int16(binary.LittleEndian.Uint16(dmsg.DataList[i].Data[j*2 : (j+1)*2]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_INT32 {
					if dmsg.DataList[i].Meta.Msb {
						valueList[i][j] = int32(binary.BigEndian.Uint32(dmsg.DataList[i].Data[j*4 : (j+1)*4]))
					} else {
						valueList[i][j] = int32(binary.LittleEndian.Uint32(dmsg.DataList[i].Data[j*4 : (j+1)*4]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_INT64 {
					if dmsg.DataList[i].Meta.Msb {
						valueList[i][j] = int64(binary.BigEndian.Uint16(dmsg.DataList[i].Data[j*8 : (j+1)*8]))
					} else {
						valueList[i][j] = int64(binary.LittleEndian.Uint16(dmsg.DataList[i].Data[j*8 : (j+1)*8]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_FLOAT {
					if dmsg.DataList[i].Meta.Msb {
						valueList[i][j] = float32(binary.BigEndian.Uint32(dmsg.DataList[i].Data[j*4 : (j+1)*4]))
					} else {
						valueList[i][j] = float32(binary.LittleEndian.Uint32(dmsg.DataList[i].Data[j*4 : (j+1)*4]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_DOUBLE {
					if dmsg.DataList[i].Meta.Msb {
						valueList[i][j] = float64(binary.BigEndian.Uint64(dmsg.DataList[i].Data[j*8 : (j+1)*8]))
					} else {
						valueList[i][j] = float64(binary.LittleEndian.Uint64(dmsg.DataList[i].Data[j*8 : (j+1)*8]))
					}
				}
			}
		}
		sql := fmt.Sprintf("insert into %s using dp tags(?,?,?) values %s", tableName, colStr)
		for idx, ts := range tsList {
			values := valueList[idx]
			// anyValues[0] = dmsg.Mode
			// anyValues[1] = dmsg.DNodeId
			// anyValues[2] = dmsg.DataList[idx].Offset
			values[3] = ts

			res := dao.TaosCli.Exec(sql, values...)
			if res < 0 {
				ulog.Log().E("dpdao", fmt.Sprintf("failed to insert into %s", tableName))
				return -1
			} else {
				ulog.Log().I("dpdao", fmt.Sprintf("insert into %s success", tableName))
			}
		}
	}

	return 0
}

func (dao *DpDataDao) QueryDp(tic string, toc string, dnodeClass string, dnodeId int64, offset int, dataMeta *meta.DataMeta) []*model.DpData {

	dpList := make([]*model.DPoint, 0)
	//convert date string to int64
	tic_time, _ := time.Parse("2006-01-02 15:04:05.000", tic)
	tic_ms := tic_time.UnixNano() / 1000000
	toc_time, _ := time.Parse("2006-01-02 15:04:05.000", toc)
	toc_ms := toc_time.UnixNano() / 1000000

	tableName := fmt.Sprintf("dp_%s_%d_%d", dnodeClass, dnodeId, offset)

	sql := fmt.Sprintf("select * from %s where ts between %d and %d", "dp", tableName, tic_ms, toc_ms)
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
				ts := t.UnixNano()
				// ulog.Log().I("dpdao", fmt.Sprintf("ts: %d, v: %d, dnode_class: %d, dnode_id: %d, dp_offset_idx: %d", t.UnixNano()/1000000, values[0], dnodeClass, dnodeId, dpOffsetIdx))

				dp := &model.DnData{
					Ts:         ts,
					DpDataList: make([]*model.DPoint, 1),
					DataMeta:   dataMeta,
					Data:       make([]byte, dataMeta.Dimen*dataMeta.ByteLen),
				}

				i := 0
				for i < dataMeta.Dimen {
					if dataMeta.ByteLen == 1 {
					} else if dataMeta.ByteLen == 2 {
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
