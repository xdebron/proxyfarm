package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var querier_db, _ = sql.Open("mysql", DB_CONN_STR)
var mysql_queue = make(chan string, 20000)

func fire_and_forget_querier() {
	for {
		query := <-mysql_queue
		_, err := querier_db.Exec(query)
		if err != nil {
			fmt.Println("querier.go mysql query queue: ", err)
		}
	}
}
