package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor struct {
	sId     string
	sUrl    string
	sUrls   []string
	sStatus string
}

func (t *C_monitor) CheckUrl(_sUrl string, c chan<- C_monitor) {
	// url get 요청
	resp, err := http.Get(_sUrl)
	status := resp.Status
	// 에러 발생 또는 상태코드가 400과 같거나 큰 경우 에러처리
	if err != nil || resp.StatusCode >= 400 {
		status = "ERR"

	}
	c <- C_monitor{sUrl: _sUrl, sStatus: status}
}

func (t *C_monitor) CheckUrl_multi(_sUrl []string) {

	// URL 접속 결과를 담을 map 선언
	results := make(map[string]string)

	// chanel 생성
	c := make(chan C_monitor)

	// 인자로 받을 url 접속 시도
	for _, url := range _sUrl {
		go t.CheckUrl(url, c)
	}

	// 결과물 map 저장
	for i := 0; i < len(_sUrl); i++ {
		result := <-c
		results[result.sUrl] = result.sStatus
	}

	// 결과물 출력
	for url, status := range results {
		fmt.Println(url, status)
	}

}

func Run_CheckUrl(_sUrl []string) {
	run := C_monitor{}
	run.CheckUrl_multi(_sUrl)
}
