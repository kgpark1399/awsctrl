package main

import (
	"sync"

	"monitor/monitor"
)

func main() {
	var wait sync.WaitGroup
	wait.Add(2)

	t := monitor.C_monitor{}

	// DB 접속
	t.DB_conn("mysql", "root", "devtools1!", "127.0.0.1:3306", "monitor")
	defer t.DB_close()

	// Monitor 로그 저장
	// t.Set_logfile("monitor.txt")

	// 메인 기능 동작
	go t.Monitor__checkUrl(5)
	go t.Monitor__checkStatus(10)

	wait.Wait()
}
