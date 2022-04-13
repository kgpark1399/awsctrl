package main

import (
	"log"
	"monitor/monitor"
	"os"
	"os/signal"
)

func Init() {
	// 모든 Init 함수들을 호출하여 사전 테스트
	// 실패 시 프로그램 정지
}

func main() {

	t := monitor.C_monitor{}
	t.Conn("mysql", "root", "pwd", "127.0.0.1:3306", "monitor")
	defer t.Close()

	t.Enable_log("monitor.txt")

	t.Run_(10)

	channel := make(chan os.Signal, 1)

	signal.Notify(channel)

	for sig := range channel {
		log.Println("Got signal:", sig)
		os.Exit(1)
	}
}
