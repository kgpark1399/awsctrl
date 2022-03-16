package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func Run_GetUrls() {
	run := C_monitor{}
	run.GetUrls()
}

func (t *C_monitor) GetUrls() {
	var website C_monitor
	db, err := sql.Open("mysql", "root:devtools1!@tcp(3.34.1.156:3306)/monitor")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,url,status FROM target")

	if err != nil {
		panic(err.Error())
	}

	var websites []C_monitor

	for rows.Next() {
		if err := rows.Scan(&website.sId, &website.sUrl, &website.sStatus); err != nil {
			panic(err.Error())
		}

		websites = append(websites, website)
	}

	for _, target := range websites {
		t.sUrls = append(t.sUrls, target.sUrl)
	}
	fmt.Println(t.sUrls)
}
