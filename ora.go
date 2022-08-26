package main

import (
	//"time"
	"log"

	"database/sql"
	_ "github.com/sijms/go-ora/v2"
)


func main() {

	conn, err := sql.Open("oracle", "oracle://system:oracle@localhost:49161/xe")
	if err!=nil {
		log.Fatalf("Failed sql.Open: %s", err)
	}
	defer conn.Close()

	//stmt, err := conn.Prepare("SELECT col_1, col_2, col_3 FROM table WHERE col_1 = :1 or col_2 = :2")
	//if err!=nil {
	//	log.Fatalf("Failed Prepare: %s", err)
	//}
	//defer stmt.Close()

	// suppose we have 2 params one time.Time and other is double
	//rows, err := stmt.Query(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC), 9.2)
	rows, err := conn.Query("select to_char(sysdate,'yyyy-mm-dd hh24:mi:ss') AS name from dual")
	if err!=nil {
		log.Fatalf("Failed Query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err!=nil {
			log.Printf("Failed Scan: %s", err)
		} else {
			log.Printf("fetch item: %s", name)
		}
	}

}
