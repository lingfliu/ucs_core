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
 * DpDataDao 数据点位CRUD接口
 * 基于TDengine实现
 * 超表命名规则: dp_{node_class}_{dp_class}
 * 标签: dnode_id bigint | dpoint_id bigint | pos varchar
 */
type DpDataDao struct {
	TaosCli *rtdb.TaosCli
}

func NewDpDataDao(host string, dbName string, username string, password string) *DpDataDao {
	return &DpDataDao{TaosCli: rtdb.NewTaosCli(host, dbName, username, password)}
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
		stableName := CreateSTableName(node.Class, pt.Offset)
		colStr := "ts timestamp"
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
			// } else if dataClass == meta.DATA_CLASS_INT64 {
			// 	tdType = "bigint unsigned"
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
			if pt.DataMeta.ValAlias != nil && len(pt.DataMeta.ValAlias) > i && pt.DataMeta.ValAlias[i] != "" {
				//a valid column name
				colStr += fmt.Sprintf(", %s %s", pt.DataMeta.ValAlias[i], tdType)
			} else {
				colStr += fmt.Sprintf(", v%d %s", i, tdType)
			}
		}

		sql := fmt.Sprintf("create stable if not exists %s (%s) tags (%s)", stableName, colStr, CreateTags())

		ulog.Log().I("dpdao", fmt.Sprintf("create stable sql: %s", sql))

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

func CreateSTableName(dnodeClass string, offset int) string {
	return fmt.Sprintf("dp_%s_%d", dnodeClass, offset)
}

func CreateTableName(dnodeClass string, offset int, dnodeId int64, dpointOffset int) string {
	return fmt.Sprintf("dp_%s_%d_%d_%d", dnodeClass, offset, dnodeId, dpointOffset)
}

func CreateTags() string {
	return "dnode_id bigint, dpoint_offset int, dnode_name nchar(32), dpoint_alias nchar(32), dpoint_unit binary(16)"
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
		stableName := CreateSTableName(dmsg.DNodeClass, v.Offset)
		tableName := CreateTableName(dmsg.DNodeClass, v.Offset, dmsg.DNodeId, v.Offset)

		dimen := v.Meta.Dimen
		colStr := "(?"
		i := 0
		for i < dimen {
			colStr += ", ?"
			i++
		}
		colStr += ")"
		sql := fmt.Sprintf("insert into %s using %s tags(?,?,?,?,?) values%s", tableName, stableName, colStr) //TODO: stable & table differentiation

		valueIdx := 6
		for tidx, ts := range tsList {
			values := make([]any, dmsg.DataList[idx].Meta.Dimen+6) //every sample
			baseIdx := tidx * dmsg.DataList[idx].Meta.ByteLen * dmsg.DataList[idx].Meta.Dimen
			for j := 0; j < dmsg.DataList[idx].Meta.Dimen; j++ {
				if v.Meta.DataClass == meta.DATA_CLASS_INT16 {
					if dmsg.DataList[idx].Meta.Msb {
						values[j+valueIdx] = int16(binary.BigEndian.Uint16(dmsg.DataList[idx].Data[baseIdx+j*2 : baseIdx+(j+1)*2]))
					} else {
						values[j+valueIdx] = int16(binary.LittleEndian.Uint16(dmsg.DataList[idx].Data[baseIdx+j*2 : baseIdx+(j+1)*2]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_INT32 {
					if dmsg.DataList[idx].Meta.Msb {
						values[j+valueIdx] = int32(binary.BigEndian.Uint32(dmsg.DataList[idx].Data[baseIdx+j*4 : baseIdx+(j+1)*4]))
					} else {
						values[j+valueIdx] = int32(binary.LittleEndian.Uint32(dmsg.DataList[idx].Data[baseIdx+j*4 : baseIdx+(j+1)*4]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_INT64 {
					if dmsg.DataList[idx].Meta.Msb {
						values[j+valueIdx] = int64(binary.BigEndian.Uint16(dmsg.DataList[idx].Data[baseIdx+j*8 : baseIdx+(j+1)*8]))
					} else {
						values[j+valueIdx] = int64(binary.LittleEndian.Uint16(dmsg.DataList[idx].Data[baseIdx+j*8 : baseIdx+(j+1)*8]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_FLOAT {
					if dmsg.DataList[idx].Meta.Msb {
						values[j+valueIdx] = float32(binary.BigEndian.Uint32(dmsg.DataList[idx].Data[baseIdx+j*4 : baseIdx+(j+1)*4]))
					} else {
						values[j+valueIdx] = float32(binary.LittleEndian.Uint32(dmsg.DataList[idx].Data[baseIdx+j*4 : baseIdx+(j+1)*4]))
					}
				} else if v.Meta.DataClass == meta.DATA_CLASS_DOUBLE {
					if dmsg.DataList[idx].Meta.Msb {
						values[j+valueIdx] = float64(binary.BigEndian.Uint64(dmsg.DataList[idx].Data[baseIdx+j*8 : baseIdx+(j+1)*8]))
					} else {
						values[j+valueIdx] = float64(binary.LittleEndian.Uint64(dmsg.DataList[idx].Data[baseIdx+j*8 : baseIdx+(j+1)*8]))
					}
				}
			}

			// sql := fmt.Sprintf("insert into %s using %s tags(%d,%d,\"%s\",\"%s\",\"%s\") values(%d, %f)", tableName, stableName, dmsg.DNodeId, dmsg.DataList[idx].Offset, dmsg.DNodeName, dmsg.DataList[idx].PtAlias, dmsg.DataList[idx].Meta.Unit, ts, values[6].(float32)) //TODO: stable & table differentiation
			values[0] = dmsg.DNodeId
			values[1] = dmsg.DataList[idx].Offset
			values[2] = "\"" + dmsg.DNodeName + "\""
			values[3] = "\"" + dmsg.DataList[idx].PtAlias + "\""
			values[4] = "\"" + dmsg.DataList[idx].Meta.Unit + "\""
			values[5] = ts

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

func (dao *DpDataDao) Query(tic string, toc string, dnodeClass string, dnodeId int64, dpointOffset int, limit int, offset int, dataMeta *meta.DataMeta) ([]int64, *model.DpData) {

	tsList := make([]int64, 0)
	data := make([]any, 0)

	//convert date string to int64
	tic_time, _ := time.Parse("2006-01-02 15:04:05.000", tic)
	tic_ms := tic_time.UnixMilli()
	toc_time, _ := time.Parse("2006-01-02 15:04:05.000", toc)
	toc_ms := toc_time.UnixMilli()

	stableName := CreateSTableName(dnodeClass, dpointOffset)
	sql := fmt.Sprintf("select * from %s where ts between %d and %d and dnode_id = %d and dpoint_offset = %d limit %d offset %d", stableName, tic_ms, toc_ms, dnodeId, dpointOffset, limit, offset)
	rows := dao.TaosCli.Query(sql)
	if rows == nil {
		ulog.Log().E("dpdao", "failed to query dp")
	} else {
		defer rows.Close()
		for rows.Next() {
			//read data
			var ts string
			var dnodeId int64
			var dpOffset int32
			var dnodeName string
			var dpAlias string
			var dpUnit string //not used
			scanned := make([]any, dataMeta.Dimen+6)

			values := make([]any, dataMeta.Dimen)

			scanned[0] = &ts
			i := 0
			for i < dataMeta.Dimen {
				scanned[i+1] = &values[i]
				i++
			}

			sidx := 1 + dataMeta.Dimen
			scanned[sidx] = &dnodeId
			scanned[sidx+1] = &dpOffset
			scanned[sidx+2] = &dnodeName
			scanned[sidx+3] = &dpAlias
			scanned[sidx+4] = &dpUnit

			err := rows.Scan(scanned...)

			if err != nil {
				ulog.Log().E("dpdao", "failed to scan dp")
			} else {
				// Handle both with and without milliseconds
				t, err := time.Parse("2006-01-02T15:04:05.000+08:00", ts)
				if err != nil {
					t, _ = time.Parse("2006-01-02T15:04:05+08:00", ts)
				}
				ts_ms := t.UnixMilli()

				tsList = append(tsList, ts_ms)
				data = append(data, values)
			}
		}
	}

	return tsList, &model.DpData{
		Offset:   offset,
		DataMeta: dataMeta,
		Data:     data,
	}
}

func (dao *DpDataDao) AggrQuery(tic string, toc string, dnodeClass string, dnodeId int64, offset int, idx int, dataMeta *meta.DataMeta, window int64, step int64, ops []int) ([]int64, *model.DpData) {
	//convert date string to int64
	tic_time, _ := time.Parse("2006-01-02 15:04:05.000", tic)
	tic_ms := tic_time.UnixMilli()
	toc_time, _ := time.Parse("2006-01-02 15:04:05.000", toc)
	toc_ms := toc_time.UnixMilli()
	stableName := CreateSTableName(dnodeClass, offset)

	var opCodes []string
	for _, op := range ops {
		opCodes = append(opCodes, dao.TaosCli.GenAggrOpCode(op))
	}

	colStr := ""
	var colName string
	if dataMeta.ValAlias != nil && len(dataMeta.ValAlias) > idx {
		colName = dataMeta.ValAlias[idx]
	} else {
		colName = fmt.Sprintf("v%d", idx)
	}
	for _, op := range opCodes {
		colStr += fmt.Sprintf("%s(%s),", op, colName)
	}
	// for i := 0; i < dataMeta.Dimen; i++ {
	// 	colStr += fmt.Sprintf("%s(v%d)", opCodes[i], i)
	// }
	colStr = colStr[:len(colStr)-1]
	ulog.Log().I("main", fmt.Sprintf("Generated column name: %s", colName))

	sql := fmt.Sprintf("select _wstart, _wend, %s from %s where ts between %d and %d interval(%d) sliding(%d)", colStr, stableName, tic_ms, toc_ms, window, step)
	rows := dao.TaosCli.Query(sql)
	tsList := make([]int64, 0)
	aggrList := make([]any, 0)
	if rows == nil {
		ulog.Log().E("dpdao", "failed to query dp")
	} else {
		defer rows.Close()
		for rows.Next() {
			//read data
			var tStartTime, tEndTime time.Time
			data := make([]float32, len(opCodes))
			scanned := make([]any, len(opCodes)+2)
			scanned[0] = &tStartTime
			scanned[1] = &tEndTime
			for i := range data {
				scanned[2+i] = &data[i] // 为每列分配指针
			}
			if err := rows.Scan(scanned...); err != nil {
				ulog.Log().E("main", fmt.Sprintf("Error scanning row: %v", err))
				continue
			}
			tStart := tStartTime.UnixMilli()
			tEnd := tEndTime.UnixMilli()
			tsList = append(tsList, tStart)
			aggrList = append(aggrList, data)
			ulog.Log().I("main", fmt.Sprintf("Parsed Row: tStart=%d, tEnd=%d, data=%+v", tStart, tEnd, data))
		}
	}

	return tsList, &model.DpData{
		Offset:   offset,
		DataMeta: dataMeta,
		Data:     aggrList,
	}
}

func (dao *DpDataDao) DeleteBefore(dnodeClass string, dnodeId int64, offset int, ts string) int {
	//convert date string to int64
	ts_time, _ := time.Parse("2006-01-02 15:04:05.000", ts)
	ts_ms := ts_time.UnixNano() / 1000000

	tableName := CreateSTableName(dnodeClass, offset)
	sql := fmt.Sprintf("delete from %s where ts < %d and dnode_id = %d and dp_offset = %d", tableName, ts_ms, dnodeId, offset)
	return dao.TaosCli.Exec(sql)
}

func (dao *DpDataDao) DeleteBetween(dnodeClass string, dnodeId int64, offset int, tic string, toc string) int {
	//convert date string to int64
	tic_time, _ := time.Parse("2006-01-02 15:04:05.000", tic)
	tic_ms := tic_time.UnixNano() / 1000000
	toc_time, _ := time.Parse("2006-01-02 15:04:05.000", toc)
	toc_ms := toc_time.UnixNano() / 1000000

	tableName := CreateSTableName(dnodeClass, offset)
	sql := fmt.Sprintf("delete from %s where ts between %d and %d and dnode_id = %d and dp_offset = %d", tableName, tic_ms, toc_ms, dnodeId, offset)
	return dao.TaosCli.Exec(sql)
}
