package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor struct {
	C_monitor_db
	C_monitor_log

	s_monitor__url    string
	s_monitor__name   string
	s_monitor__status string

	n_monitor__rate int

	arrs_monitor__urls       []string
	arrs_monitor__status_grp []string
}

// URL HTTP 상태 체크 실행
func (t *C_monitor) Run_check__url(_n_monitor__rate int) {

	// DB에서 모니터링 대상 URL 호출
	target := t.Get__urls()

	// 반복 시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_monitor__rate))

	for range ticker.C {

		// 모니터링 대상 URL string 으로 변경 후 http 상태 조회
		for _, url := range target {
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode >= 400 {
				// http status 오류의 경우 DB status 값을 0(false)로 변경
				log.Println("URL :", url, ", STATUS : ERR ")
				t.Change_status__false(url)
			} else {
				// http status 정상의 경우 DB status 값을 0(false)로 변경
				log.Println("URL :", url, ", STATUS :", resp.Status)
				t.Change_status__true(url)
			}
		}
	}
}

// DB URL status 값 체크 및 알림 발송
func (t *C_monitor) Run_check__status(_n_monitor__rate_min int) {

	// 반복시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_monitor__rate_min))

	for range ticker.C {

		// DB url, status 값 호출 및 string 변환
		url, status := t.Get__status()
		for i, _status := range status {
			// status 값 체크하여 false의 경우 알림 발송
			if _status == "true" {
				fmt.Print()
			} else {
				log.Println("====== URL :", url[i], ", SERVER STATUS :", _status, "======")
			}
		}
	}
}
