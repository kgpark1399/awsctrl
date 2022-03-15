package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type C_database struct {
	sDBtype string
	sDBid   string
	sDBpwd  string
	sAddr   string
	sDBname string
}

func ConnectDB() {
	db, err := sql.Open("mysql", "root:devtools1!@tcp(3.34.1.156:3306)/monitor")
	if err != nil {
		log.Fatal("database open error:", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("SELECT url FROM target WHERE id = 1").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
