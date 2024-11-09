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
	DbName   string
	Username string
	Password string
	taos     *sql.DB

	Io chan int
}

func NewTaosCli(host string, dbName string, username string, password string) *TaosCli {
	return &TaosCli{
		Host:     host,
		DbName:   dbName,
		Username: username,
		Password: password,
	}
}

func (cli *TaosCli) Open() {
	taos, err := sql.Open("taosSql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cli.Username, cli.Password, cli.Host, cli.DbName))
	if err != nil {
		ulog.Log().E("taos", "failed to connect to taos, err: "+err.Error())
	}
	cli.taos = taos
}

func (cli *TaosCli) Close() {
	cli.taos.Close()
}

func (cli *TaosCli) ShowDatabases() []string {
	databaseList := make([]string, 0)
	rows, err := cli.taos.Query("show databases")
	if err != nil {
		ulog.Log().E("taos", "failed to query taos, err: "+err.Error())
		return databaseList
	}
	defer rows.Close()

	for rows.Next() {
		var v string
		rows.Scan(&v)
		databaseList = append(databaseList, v)
	}
	return databaseList
}

func (cli *TaosCli) ShowSTables() []string {
	stableList := make([]string, 0)

	rows, err := cli.taos.Query("show stables")
	if err != nil {
		ulog.Log().E("taos", "failed to query taos, err: "+err.Error())
		return stableList
	}
	defer rows.Close()

	for rows.Next() {
		var v string
		rows.Scan(&v)
		stableList = append(stableList, v)
	}
	return stableList
}

func (cli *TaosCli) ShowTables() []string {
	tableList := make([]string, 0)
	rows, err := cli.taos.Query("show tables")
	if err != nil {
		ulog.Log().E("taos", "failed to query taos, err: "+err.Error())
		return tableList
	}
	defer rows.Close()
	for rows.Next() {
		var v string
		rows.Scan(&v)
		tableList = append(tableList, v)
	}
	return tableList
}

func (cli *TaosCli) Exec(sql string, args ...any) int {
	res, err := cli.taos.Exec(sql, args...)
	if err != nil {
		ulog.Log().E("taos", "failed to exec taos, err: "+err.Error())
		return -1
	} else {
		affected, err := res.RowsAffected()
		if err != nil {
			ulog.Log().E("taos", "failed to get affected rows, err: "+err.Error())
			return -1
		} else {
			ulog.Log().I("taos", "done, affected rows: "+fmt.Sprintf("%d", affected))
			return 0
		}
	}
}

func (cli *TaosCli) Query(sql string, args ...any) *sql.Rows {
	rows, err := cli.taos.Query(sql, args...)
	if err != nil {
		ulog.Log().E("taos", fmt.Sprintf("failed to query %s, err: "+err.Error(), sql))
		return nil
	}
	return rows
}

func (cli *TaosCli) CreateSTable(stableName string, columns string, tag_columns string) {
	cli.taos.Exec(fmt.Sprintf("create stable if not exist %s.%s(%s) tags(%s)", cli.DbName, stableName, columns, tag_columns))
}

func (cli *TaosCli) CreateTable(tableName string, stableName string, tags []string) {
	tagStr := ""
	tagStr += tags[0]
	for i := 1; i < len(tags); i++ {
		tagStr += tags[i] + ","
	}
	cli.taos.Exec(fmt.Sprintf("create table if not exist %s using %s.%s tags(%s)", tableName, cli.DbName, stableName, tagStr))
}
