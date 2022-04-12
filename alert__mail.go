package monitor

import (
	"log"
	"net"
	"net/smtp"
	"time"

	"gopkg.in/ini.v1"
)

type C_sendmail struct {

	// 메일 연동 정보 (ini)
	s_mail__id   string
	s_mail__pwd  string
	s_mail__host string
	s_mail__port string

	// 메일 발송 옵션
	s_mail__from string
	s_mail__body string

	arrs_mail__to []string
}

func (t *C_sendmail) Init__sendmail() (id, pwd, host, port string, err error) {

	// config ini 파일 읽기
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Println("Fail to read config.ini file : ", err)
		return "", "", "", "", err
	}

	// 메일 관련 정보 저장
	t.s_mail__id = cfg.Section("SMTP").Key("S_mail__id").String()
	t.s_mail__pwd = cfg.Section("SMTP").Key("S_mail__pwd").String()
	t.s_mail__host = cfg.Section("SMTP").Key("S_mail__host").String()
	t.s_mail__port = cfg.Section("SMTP").Key("S_mail__port").String()

	smtpserver := t.s_mail__host + ":" + t.s_mail__port

	// SMTP 서버 통신 체크
	conn, err := net.DialTimeout("tcp", smtpserver, 3*time.Second)
	if err != nil {
		log.Println("Fail to connect smtp server : ", err, conn)
		return "", "", "", "", err
	}

	return t.s_mail__id, t.s_mail__pwd, t.s_mail__host, t.s_mail__port, nil
}

func (t *C_sendmail) Set__mail(_s_mail__body string, arrs_mail__to []string) error {

	// ini 파일 메일 발송 정보 호출
	id, pwd, host, port, err := t.Init__sendmail()
	if err != nil {
		return err
	}

	// 메일 송수신 정보 입력
	auth := smtp.PlainAuth("", id, pwd, host)
	from := id
	to := arrs_mail__to

	// 메시지 작성
	headerSubject := "Subject: [Alert] Server Error\r\n"
	headerBlank := "\r\n"
	body := _s_mail__body
	msg := []byte(headerSubject + headerBlank + body)

	// 메일 발송
	smtp_server := host + ":" + port
	err = smtp.SendMail(smtp_server, auth, from, to, msg)
	if err != nil {
		log.Println("Fail to connect smtp server  : ", err)
		return err
	}
	return err
}

func Send__mail(_s_mail__body string, arrs_mail__to []string) error {

	t := C_sendmail{}
	err := t.Set__mail(_s_mail__body, arrs_mail__to)
	if err != err {
		return err
	}
	return nil
}
