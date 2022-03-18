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
	db, err := sql.Open("mysql", "root:devtools1!@tcp(3.34.1.156:3306)/monitor")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT url,status FROM target")

	if err != nil {
		panic(err.Error())
	}

	var websites []C_monitor

	for rows.Next() {
		if err := rows.Scan(&website.sUrl, &website.iStatus); err != nil {
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
func (t *C_monitor) ChagneStatus(_sUrl string) {
	db, err := sql.Open("mysql", "root:devtools1!@tcp(3.34.1.156:3306)/monitor")
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
	_, err = stmt.Exec(1, _sUrl) //Placeholder 파라미터 순서대로 전달
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
	rows, err := db.Query("SELECT url,status FROM target")

	if err != nil {
		panic(err.Error())
	}

	var groupstatus []C_monitor

	// 쿼리 조회 후 변수 저장
	for rows.Next() {
		if err := rows.Scan(&getstatus.sUrl, &getstatus.iStatus); err != nil {
			panic(err.Error())
		}

		groupstatus = append(groupstatus, getstatus)
	}

	// 쿼리 데이터 중 상태 데이터만 추출하여 배열 저장
	for _, target_s := range groupstatus {
		t.sStatusGrp = append(t.sStatusGrp, target_s.iStatus)
	}

	// 결과 데이터 반환
	result = t.sStatusGrp
	return result
}

func (t *C_monitor) CreateUrl(_sName, _sUrl string, _iStatus int) {
	db, err := sql.Open("mysql", "root:devtools1!@tcp(3.34.1.156:3306)/monitor")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	// INSERT 문 실행
	result, err := db.Exec("INSERT INTO target (name,url,status) VALUES (?, ?, ?)", _sName, _sUrl, _iStatus)
	if err != nil {
		fmt.Println(err)
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		fmt.Println("1 row inserted.")
	}
	fmt.Println(err)
}
