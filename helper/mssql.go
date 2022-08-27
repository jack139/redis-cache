package helper

import (
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	conn_mssql *sql.DB
)

func mssql_init() error {
	var err error

	conn_mssql, err = sql.Open("mssql", Settings.Server.MSSQL_CONNECTION)
	if err != nil {
		return err
	}

	log.Println("MS-SQL connected.")

	return nil
}


func Mssql_test() (string,  error) {
	stmt, err := conn_mssql.Prepare("select name from master.dbo.sysdatabases WHERE name LIKE :1")
	if err!=nil {
		return "", fmt.Errorf("Failed Prepare: %s", err)
	}
	defer stmt.Close()

	// 替换变量 :1
	rows, err := stmt.Query("m%")
	if err!=nil {
		return "", fmt.Errorf("Failed Query: %s", err)
	}
	defer rows.Close()

	var sqlData []string
	for rows.Next() {
		var item sql.NullString
		err = rows.Scan(&item)
		if err!=nil {
			return "", fmt.Errorf("Failed Scan: %s", err)
		}

		sqlData = append(sqlData, item.String)
	}

	// json 返回
	msgBody, err := json.Marshal(map[string]interface{}{
		"key": "123",
		"value": sqlData,
	})
	if err != nil {
		return "", fmt.Errorf("Failed Json: %s", err)
	}

	return string(msgBody), nil

}