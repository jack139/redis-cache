package helper

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/sijms/go-ora/v2"
)

var (
	conn_ora *sql.DB
)

func ora_init() error {
	var err error

	conn_ora, err = sql.Open("oracle", Settings.Server.ORA_CONNECTION)
	if err!=nil {
		return err
	}

	log.Println("Oracle connected.")

	return nil
}


func Ora_test() (string,  error) {
	stmt, err := conn_ora.Prepare("select table_name, num_rows FROM user_tables WHERE table_name LIKE :1 ")
	if err!=nil {
		return "", fmt.Errorf("Failed Prepare: %s", err)
	}
	defer stmt.Close()

	// 替换变量 :1
	rows, err := stmt.Query("DEF%")
	if err!=nil {
		return "", fmt.Errorf("Failed Query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item sql.NullString
		var num sql.NullString // 可能为 NULL
		err = rows.Scan(&item, &num)
		if err!=nil {
			return "", fmt.Errorf("Failed Scan: %s", err)
		}
		
		log.Printf("fetch item: %v %v", item.String, num.String)
	}

	return "", nil

}
