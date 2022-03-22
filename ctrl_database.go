package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (t *C_db_config) Init_db(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) string {

	config := _s_db__id + ":" + _s_db__pwd + "@" + "tcp" + "(" + _s_db__hostname + ")" + "/" + _s_db__name + ")"
	return config

}

// DB Connection 설정
func (t *C_db_config) DB_conn(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) error {
	var err error
	config := t.Init_db(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)
	t.db_conn, err = sql.Open("mysql", config)
	if err != nil {
		return err
	}
	return nil
}

func (t *C_db_config) DB_close() error {
	return t.db_conn.Close()
}

// DB URL 정보 호출 및 반환
func (t *C_monitor) GetUrls() (result []string) {
	var website C_monitor
	var websites []C_monitor

	rows, err := t.db_conn.Query("SELECT url,status FROM target")

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		if err := rows.Scan(&website.s_url, &website.n_status); err != nil {
			panic(err.Error())
		}

		websites = append(websites, website)
	}

	for _, target := range websites {
		t.arrs_urls = append(t.arrs_urls, target.s_url)
	}

	result = t.arrs_urls
	return result
}

// DB status 컬럼 데이터 변경 (check 실패 시)
func (t *C_monitor) Chagne_status__false(_sUrl string) {

	stmt, err := t.db_conn.Prepare("UPDATE target SET status=? WHERE url=?")
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

func (t *C_monitor) Chagne_status__true(_sUrl string) {

	stmt, err := t.db_conn.Prepare("UPDATE target SET status=? WHERE url=?")
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
	var groupstatus []C_monitor

	// 쿼리 대상
	rows, err := t.db_conn.Query("SELECT url,status FROM target")

	if err != nil {
		panic(err.Error())
	}

	// 쿼리 조회 후 변수 저장
	for rows.Next() {
		if err := rows.Scan(&getstatus.s_url, &getstatus.n_status); err != nil {
			panic(err.Error())
		}

		groupstatus = append(groupstatus, getstatus)
	}

	// 쿼리 데이터 중 상태 데이터만 추출하여 배열 저장
	for _, target_s := range groupstatus {
		t.arrn_status_grp = append(t.arrn_status_grp, target_s.n_status)
	}

	// 결과 데이터 반환
	result = t.arrn_status_grp
	return result
}

func (t *C_monitor) CreateUrl(_sName, _sUrl string, _iStatus int) {

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO target (name,url,status) VALUES (?, ?, ?)", _sName, _sUrl, _iStatus)
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

// DB status 에러 상태 인 URL 데이터 호출
func (t *C_monitor) GetUrls_Err() (result []string) {
	var website C_monitor
	var websites []C_monitor

	// status = 1(에러) 값 쿼리
	rows, err := t.db_conn.Query("SELECT url FROM target WHERE status=1")

	if err != nil {
		panic(err.Error())
	}

	// 에러 상태인 URL 데이터 호출
	for rows.Next() {
		if err := rows.Scan(&website.s_url); err != nil {
			panic(err.Error())
		}

		websites = append(websites, website)
	}

	for _, target := range websites {
		t.arrs_urls = append(t.arrs_urls, target.s_url)
	}
	result = t.arrs_urls
	return result

}
