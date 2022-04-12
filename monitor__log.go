package monitor

import (
	"io"
	"log"
	"os"
)

type C_monitor__log struct {
	s_file__name string
}

// 모니터링 로그 활성화
func (t *C_monitor__log) Enable_log(_s_file__name string) error {

	err := t.Init()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(_s_file__name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("[ERROR] Failed to open log file : ", err)
		return err
	}

	// 커맨드 출력과 로그 동시 기록
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	return nil
}

// 모니터링 로그 초기화
func (t *C_monitor__log) Init() error {

	var err error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
