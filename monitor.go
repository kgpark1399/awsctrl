package main

import (
	"fmt"
	"net/http"
	"time"
)

type C_monitor struct {
	sUrl    string
	sPorts  string
	sData   string
	sStatus string

	iRate int
}

func Run_checkUrl(_sUrl string, _iRate int) {
	cMonitor := New_C_checkUrl()
	cMonitor.checkUrl(_sUrl, _iRate)

}

func New_C_checkUrl() *C_monitor {
	c := &C_monitor{}
	return c
}

func (t *C_monitor) checkUrl(_sUrl string, _iRate int) string {

	// http 상태 조회 , 오류 시 red 변경
	resp, err := http.Get(_sUrl)
	if err != nil {
		fmt.Println(err)
		t.sStatus = "red"
	}

	// http 상태 체크 주기 설정
	ticker := time.NewTicker(time.Second * time.Duration(_iRate))

	// 상태 체크 동작
	for time := range ticker.C {
		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("HTTP Status is in the 2xx range", time)
			t.sStatus = "green"
		} else {
			fmt.Println("Web Error", time)
			t.sStatus = "red"
		}
	}

	// 종료 테스트
	time.Sleep(time.Second * 5)
	ticker.Stop()
	fmt.Println("Ticker stopped")

	return t.sStatus
}
