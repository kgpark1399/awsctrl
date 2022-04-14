package monitor

import (
	"fmt"
	"log"
	"net"
	"time"

	"gopkg.in/ini.v1"
)

type C_list_target struct {

	// 모니터링 대상
	s_server__url     []string
	n_runCycle_second int

	// 장애 알림 대상
	s_notice__mail   []string
	s_notice__mobile []string

	// 메일 연동 정보 (ini)
	s_mail__id   string
	s_mail__pwd  string
	s_mail__host string
	s_mail__port string
}

// ini 파일 존재 여부 체크
func (t *C_list_target) Init() (file *ini.File, err error) {

	file, err = ini.Load("config.ini")
	if err != nil {
		log.Println("[ERROR] Fail to read config.ini file : ", err)
		return file, err
	}

	return file, nil
}

// 모니터링 대상 리스트 로딩
func (t *C_list_target) Get__monitor() (url []string, cycle int, err error) {

	cfg, err := t.Init()
	if err != nil {
		return
	}

	title := "monitor"
	t.s_server__url = cfg.Section(title).Key("S_server__url").Strings(",")
	t.n_runCycle_second = cfg.Section(title).Key("N_runCycle_second").MustInt()

	return t.s_server__url, t.n_runCycle_second, nil
}

// 장애 알림 대상 리스트 로딩
func (t *C_list_target) Get__alert_notice_contact() (mail, mobile []string, err error) {

	cfg, err := t.Init()
	if err != nil {
		return
	}

	title := "notification"
	t.s_notice__mail = cfg.Section(title).Key("S_notice__mail").Strings(",")
	t.s_notice__mobile = cfg.Section(title).Key("S_notice__mobile").Strings(",")

	fmt.Println(t.s_notice__mail, t.s_notice__mobile)
	return t.s_notice__mail, t.s_notice__mobile, nil

}

// SMTP 접속 정보 리스트 로딩
func (t *C_list_target) Get__smtp_acess_info() (id, pwd, host, port string, err error) {

	// config ini 파일 읽기
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Println("[ERROR] Fail to read config.ini file : ", err)
		return "", "", "", "", err
	}

	title := "smtp"
	t.s_mail__id = cfg.Section(title).Key("S_mail__id").String()
	t.s_mail__pwd = cfg.Section(title).Key("S_mail__pwd").String()
	t.s_mail__host = cfg.Section(title).Key("S_mail__host").String()
	t.s_mail__port = cfg.Section(title).Key("S_mail__port").String()

	smtpserver := t.s_mail__host + ":" + t.s_mail__port

	// SMTP 서버 통신 체크
	conn, err := net.DialTimeout("tcp", smtpserver, 3*time.Second)
	if err != nil {
		log.Println("[ERROR] Fail to connect smtp server : ", err, conn)
		return "", "", "", "", err
	}

	return t.s_mail__id, t.s_mail__pwd, t.s_mail__host, t.s_mail__port, nil
}
