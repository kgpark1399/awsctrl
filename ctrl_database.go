package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB URL 정보 호출 및 반환
func (t *C_monitor) GetUrls() (result []string) {
	var website C_monitor
	db, err := sql.Open("mysql", "root:devtools1!@tcp(127.0.0.1:3306)/monitor")
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

	result = t.sUrls

	return result
}

// DB status 컬럼 데이터 변경 (check 실패 시)
func (t *C_monitor) Status(_sUrl string) {
	db, err := sql.Open("mysql", "root:devtools1!@tcp(127.0.0.1:3306)/monitor")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE target SET status=? WHERE url=?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec(0, _sUrl) //Placeholder 파라미터 순서대로 전달
	if err != nil {
		fmt.Println(err)
	}
}

// DB status 컬럼 데이터 조회
func (t *C_monitor) GetStatus() (result []int) {
	var getstatus C_monitor
	db, err := sql.Open("mysql", "root:devtools1!@tcp(3.34.1.156:3306)/monitor")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 쿼리 대상
	rows, err := db.Query("SELECT id,url,status FROM target")

	if err != nil {
		panic(err.Error())
	}

	var groupstatus []C_monitor

	// 쿼리 조회 후 변수 저장
	for rows.Next() {
		if err := rows.Scan(&getstatus.sId, &getstatus.sUrl, &getstatus.sStatus); err != nil {
			panic(err.Error())
		}

		groupstatus = append(groupstatus, getstatus)
	}

	// 쿼리 데이터 중 상태 데이터만 추출하여 배열 저장
	for _, target := range groupstatus {
		t.sStatusGrp = append(t.sStatusGrp, target.sStatus)
	}

	// 결과 데이터 반환
	result = t.sStatusGrp
	return result
}
