package rtdb

import (
	"database/sql"
	"fmt"

	"github.com/lingfliu/ucs_core/ulog"
	_ "github.com/taosdata/driver-go/v3/taosSql"
)

type TaosCli struct {
	Host     string
	Username string
	Password string
	taos     *sql.DB
}

func (cli *TaosCli) Open() {
	taos, err := sql.Open("taosSql", fmt.Sprintf("%s:%s@tcp(%s)", cli.Username, cli.Password, cli.Host))
	if err != nil {
		ulog.Log().E("tas", "failed to connect to taos")
	}
	defer taos.Close()
	cli.taos = taos
}

func (cli *TaosCli) Exec(sql string) {
	_, err := cli.taos.Exec(sql)
	if err != nil {
		ulog.Log().E("tas", "failed to exec taos")
	}
}

func (cli *TaosCli) Query(sql string) *sql.Rows {
	rows, err := cli.taos.Query(sql)
	if err != nil {
		ulog.Log().E("tas", "failed to query taos")
		return nil
	}
	return rows
}

func (cli *TaosCli) CreateTable(dbName string, tableName string, columns string) {
	cli.taos.Exec(fmt.Sprintf("create database if not exist %s", dbName))
	cli.taos.Exec(fmt.Sprintf("use %s", dbName))
	cli.taos.Exec(fmt.Sprintf("create table if not exist %s(%s)", tableName, columns))
}

func (cli *TaosCli) Insert(dbName string, tableName string, columns string, values string) {
	cli.taos.Exec(fmt.Sprintf("insert into %s.%s(%s) values(%s)", dbName, tableName, columns, values))
}

func (cli *TaosCli) QueryAll(dbName string, tableName string) {
	rows, err := cli.taos.Query(fmt.Sprintf("select * from %s.%s", dbName, tableName))
	if err != nil {
		ulog.Log().E("tas", "failed to query taos")
	}
	defer rows.Close()
}

func (cli *TaosCli) DeleteById(dbName string, tableName string, id string) {
	cli.taos.Exec(fmt.Sprintf("delete from %s.%s where id=%s", dbName, tableName, id))
}

func (cli *TaosCli) Close() {
	cli.taos.Close()
}
