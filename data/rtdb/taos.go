package rtdb

import (
	"database/sql"
	"fmt"

	"github.com/lingfliu/ucs_core/ulog"
	_ "github.com/taosdata/driver-go/v3/taosSql"
)

const (
	TAOS_STATE_DISCONNECTED = 0
	TAOS_STATE_CONNECTED    = 1
)

type TaosCli struct {
	Host     string
	Username string
	Password string
	taos     *sql.DB

	Io chan int
}

func NewTaosCli(Host string, Username string, Password string) *TaosCli {
	return &TaosCli{
		Host:     Host,
		Username: Username,
		Password: Password,
		Io:       make(chan int),
	}
}

func (cli *TaosCli) Open() {
	taos, err := sql.Open("taosSql", fmt.Sprintf("%s:%s@tcp(%s)", cli.Username, cli.Password, cli.Host))
	if err != nil {
		ulog.Log().E("tas", "failed to connect to taos")
	}
	cli.taos = taos
	// cli.Io <- 1
}

func (cli *TaosCli) Close() {
	cli.taos.Close()
}

// func (cli *TaosCli) ShowDatabases() {
// 	rows, err := cli.taos.Query("show databases")
// 	if err != nil {
// 		ulog.Log().E("tas", "failed to query taos")
// 		return
// 	}
// 	defer rows.Close()
// }

func (cli *TaosCli) ShowSTables() {
	taos, _ := sql.Open("taosSql", fmt.Sprintf("%s:%s@tcp(%s)", cli.Username, cli.Password, cli.Host))

	rows, err := taos.Query("show stables")
	if err != nil {
		ulog.Log().E("tas", "failed to query taos")
		return
	}
	defer rows.Close()
}
func (cli *TaosCli) ShowTables() {
	rows, err := cli.taos.Query("show tables")
	if err != nil {
		ulog.Log().E("tas", "failed to query taos")
		return
	}
	defer rows.Close()
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

func (cli *TaosCli) CreateSTable(dbName string, tableName string, columns string, tag_columns string) {
	cli.taos.Exec(fmt.Sprintf("create database if not exist %s", dbName))
	cli.taos.Exec(fmt.Sprintf("use %s", dbName))
	cli.taos.Exec(fmt.Sprintf("create stable if not exist %s(%s) tags(%s)", tableName, columns, tag_columns))
}

func (cli *TaosCli) CreateTable(dbName string, stableName string, tableName string, tags []string) {
	cli.taos.Exec(fmt.Sprintf("create database if not exist %s", dbName))
	cli.taos.Exec(fmt.Sprintf("use %s", dbName))
	tagStr := ""
	for _, tag := range tags {
		tagStr += tag + " "
	}
	cli.taos.Exec(fmt.Sprintf("create table if not exist using %s, %s(%s)", stableName, tableName, tagStr))
}

func (cli *TaosCli) Insert(tableName string, columns []string, tags []string) {
	columnStr := ""
	for _, column := range columns {
		columnStr += column + " "
	}
	tagStr := ""
	for _, tag := range tags {
		tagStr += tag + " "
	}
	cli.taos.Exec(fmt.Sprintf("insert into %s values(%s) tags(%s)", tableName, columnStr, tagStr))
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
