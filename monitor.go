package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// URL HTTP 상태 체크 실행

func Run_check__url(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) {
	t := C_monitor{}
	t.check__url(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)
}

func Run_check__status(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) {
	t := C_monitor{}
	t.check__status(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)
}

// URL HTTP 상태 체크 실행
func (t *C_monitor) check__url(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) {

	// DB 접속
	t.DB_conn(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)

	// DB에서 모니터링 대상 URL 호출
	target := t.GetUrls()

	// 반복 시간 설정
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for time := range ticker.C {

			// 모니터링 대상 URL string 으로 변경 후 http 상태 조회
			for _, url := range target {
				resp, err := http.Get(url)
				if err != nil || resp.StatusCode >= 400 {
					// http status 오류의 경우 DB status 값을 0(false)로 변경
					fmt.Println("URL :", url, ", STATUS : ERR ")
					t.Chagne_status__false(url)
				} else {
					// http status 정상의 경우 DB status 값을 0(false)로 변경
					fmt.Println("URL :", url, ", STATUS :", resp.Status)
					t.Chagne_status__true(url)

				}
				// 로그 시간 출력
				fmt.Println(time)
			}
		}
	}()
	time.Sleep(time.Second * 20)
	ticker.Stop()
	fmt.Println("monitor check stopped")
	t.DB_close()
}

// DB URL status 값 체크 실행
func (t *C_monitor) check__status(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) {
	t.DB_conn(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for time := range ticker.C {

			// DB status Err 상태 URL 데이터 호출
			var arrs_target []string
			arrs_target = t.GetUrls_Err()

			target := t.GetStatus()
			for _, status := range target {
				if status == 0 {
					fmt.Print()
				} else {
					for _, s_target := range arrs_target {
						fmt.Println("URL : ", s_target, ",  HTTP STATUS : ERROR", ", Time :", time)
						s_target := "URL :" + s_target + "Error"
						Send_sns(s_target)
					}

				}
			}
		}
	}()
	time.Sleep(time.Second * 20)
	ticker.Stop()
	fmt.Println("Ticker stopped")
	t.DB_close()

}
