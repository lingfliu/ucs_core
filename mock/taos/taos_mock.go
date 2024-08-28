package main

import (
	"database/sql"
	"strconv"

	"github.com/lingfliu/ucs_core/data/rtdb"
	"github.com/lingfliu/ucs_core/ulog"
)

func main() {
	ulog.Config(ulog.LOG_LEVEL_DEBUG, "", false)
	var taoDsn = "root:taosdata@tcp(62.234.16.239:6030)/ucs_srv"
	taos, err := sql.Open("taosSql", taoDsn)
	if err != nil {
		ulog.Log().E("tas", "failed to connect to taos")
		return
	}
	defer taos.Close()
	ulog.Log().I("tas", "connected to taos")

	// //test: create stable
	// _, err = taos.Exec("create stable if not exists ucs_srv.eval_demo (ts timestamp, val float) tags (node_id binary(32), offset_idx int)")
	// if err != nil {
	// 	ulog.Log().E("tas", "failed to create stable, err: "+err.Error())
	// 	return
	// }
	// ulog.Log().I("tas", "create stable success")

	// //test: show stable
	// rows, err := taos.Query("show stables")
	// if err != nil {
	// 	ulog.Log().E("tas", "failed to query taos")
	// 	return
	// }
	// //read stable names
	// var v string
	// for rows.Next() {
	// 	err := rows.Scan(&v)
	// 	if err != nil {
	// 		ulog.Log().E("tas", "failed to get stable list, err:"+err.Error())
	// 	} else {
	// 		ulog.Log().I("tas", "stable list: "+v)
	// 	}
	// 	if err != nil {
	// 		ulog.Log().E("tas", "failed to get columns")
	// 		return
	// 	}
	// 	// for _, col := range cols {
	// 	// 	row_str += col + " "
	// 	// }
	// 	// ulog.Log().I("tas", "query taos success, table list: "+row_str)
	// }
	// defer rows.Close()

	// //test: create table
	// for id := 0; id < 10; id++ {
	// 	op := fmt.Sprintf("create table if not exists node_%d using ucs_srv.eval_demo tags (\"001\", %d)", id, id)
	// 	_, err = taos.Exec(op)
	// 	if err != nil {
	// 		ulog.Log().E("tas", "failed to create table, err: "+err.Error())
	// 	} else {
	// 		ulog.Log().I("tas", "create table success: "+strconv.Itoa(id))
	// 	}

	// }

	// //show tables
	// rows, err = taos.Query("show tables")
	// for rows.Next() {
	// 	err := rows.Scan(&v)
	// 	if err != nil {
	// 		ulog.Log().E("tas", "failed to get table list, err:"+err.Error())
	// 	} else {
	// 		ulog.Log().I("tas", "table list: "+v)
	// 	}
	// }

	//test: remove tables
	rows, err := taos.Query("show tables")
	tables := make([]string, 0)
	for rows.Next() {
		var v string
		err := rows.Scan(&v)
		if err != nil {
			ulog.Log().E("tas", "failed to get table list, err:"+err.Error())
		} else {
			ulog.Log().I("tas", "table list: "+v)
			tables = append(tables, v)
		}
	}
	for _, table := range tables {
		res, err := taos.Exec("drop table " + table)
		if err != nil {
			ulog.Log().E("tas", "failed to drop table, err: "+err.Error())
		} else {
			ulog.Log().I("tas", "drop table success: "+table)
			affected, err := res.RowsAffected()
			if err != nil {
				ulog.Log().E("tas", "failed to get affected rows, err: "+err.Error())
			} else {
				//convert int64 to string
				ulog.Log().I("tas", "affected rows: "+strconv.FormatInt(affected, 10))
			}
		}
	}

	// //test: insert
	// for i := 0; i < 10; i++ {
	// 	for j := 0; j < 10; j++ {
	// 		insertSql := fmt.Sprintf("insert into ucs_srv.node_%d values (now, %d)", i, j)
	// 		res, err := taos.Exec(insertSql)
	// 		if err != nil {
	// 			ulog.Log().E("tas", "failed to insert data, err: "+err.Error())
	// 		} else {
	// 			affected, err := res.RowsAffected()
	// 			if err != nil {
	// 				ulog.Log().E("tas", "failed to get affected rows, err: "+err.Error())
	// 			} else {
	// 				ulog.Log().I("tas", "insert data success: "+strconv.Itoa(i)+", "+strconv.Itoa(j)+", affected rows: "+strconv.Itoa(int(affected)))
	// 			}
	// 		}
	// 	}
	// }

	// cli := &rtdb.TaosCli{
	// 	Host:     "62.234.16.239:6030/ucs_srv",
	// 	Username: "root",
	// 	Password: "taosdata",
	// }

	// go _task_connect_test(cli)
	// for {
	// 	time.Sleep(1 * time.Second)
	// }
}

func _task_connect_test(cli *rtdb.TaosCli) {
	cli.Open()
	cli.ShowSTables()
	cli.CreateSTable("ucs", "eval_demo", "ts timestamp, val float", "node_id string, offset int")
	for i := 0; i < 100; i++ {
		cli.Insert("eval_demo", []string{"ts", "val", "node_id", "offset"}, []string{"now", "3.14", "node1"})
	}
	defer cli.Close()
}
