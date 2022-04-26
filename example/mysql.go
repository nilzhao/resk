package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsName := "root:123456@tcp(127.0.0.1:3306)/envelope?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println(err)
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(3)
	db.SetConnMaxLifetime(7 * time.Hour)
	defer db.Close()

	fmt.Println(db.Query("select now()"))
}
