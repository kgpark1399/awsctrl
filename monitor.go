package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor struct {
	sId   string
	sUrl  string
	SName string

	iStatus int

	sUrls      []string
	sStatusGrp []int
}

// URL HTTP 상태 체크 실행
func Run_ChkeckUrl() {

	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for t := range ticker.C {
			_t := C_monitor{}
			target := _t.GetUrls()
			for _, url := range target {
				_t.CheckUrl(url)
				fmt.Println(t)
			}
		}
	}()
	time.Sleep(time.Minute * 10)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

// DB URL status 값 체크 실행
func Run_CheckStatus() {

	ticker := time.NewTicker(time.Minute * 3)
	go func() {
		for t := range ticker.C {
			_t := C_monitor{}
			target := _t.GetStatus()
			for _, status := range target {
				_t.CheckStatus(status)
				fmt.Println(t)
			}
		}
	}()
	time.Sleep(time.Minute * 10)
	ticker.Stop()
	fmt.Println("Ticker stopped")

}

// ---------------------------------------------------------------------- //

// URL HTTP 상태 체크 기능
func (t *C_monitor) CheckUrl(_sUrl string) {
	// url get 요청
	resp, err := http.Get(_sUrl)

	// 에러 발생 또는 상태코드가 400과 같거나 큰 경우 에러처리
	if err != nil || resp.StatusCode >= 400 {
		fmt.Println("STATUS : ERR", "URL :", _sUrl)
		t.ChagneStatus(_sUrl)

	} else {
		fmt.Println("STATUS :", resp.Status, "URL :", _sUrl)

	}
}

// DB URL status 값 체크 기능
func (t *C_monitor) CheckStatus(_sStatus int) {

	if _sStatus == 0 {
		fmt.Println("OK")

	} else {
		fmt.Println("err")

	}
}
