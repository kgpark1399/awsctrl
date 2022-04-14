package monitor

import (
	"log"
	"net"
	"net/smtp"
	"time"
)

type C_notice__mail struct {
	c_list_target C_list_target

	s_mail__from  string
	s_mail__title string
	s_mail__body  string

	arrs_mail__to []string
}

// 메일 발송
func (t *C_notice__mail) Send(s_mail__title, _s_mail__body string, arrs_mail__to []string) error {

	// ini 파일 메일 발송 정보 호출
	id, pwd, host, port, err := t.Init()
	if err != nil {
		return err
	}

	// 메일 송수신 정보 입력
	auth := smtp.PlainAuth("", id, pwd, host)
	from := id
	to := arrs_mail__to

	// 메시지 작성
	headerSubject := s_mail__title
	headerBlank := "\r\n"
	body := _s_mail__body
	msg := []byte(headerSubject + headerBlank + body)

	// 메일 발송
	smtp_server := host + ":" + port
	err = smtp.SendMail(smtp_server, auth, from, to, msg)
	if err != nil {
		log.Println("[ERROR] Fail to send mail, check your account&SMTP information : ", err)
		return err
	}
	return err
}

// 메일 발송 서버 정보 ini 호출 및 SMTP 서버 통신 테스트
func (t *C_notice__mail) Init() (id, pwd, host, port string, err error) {

	id, pwd, host, port, err = t.c_list_target.Get__smtp_acess_info()
	if err != nil {
		return "", "", "", "", err
	}

	smtpserver := host + ":" + port

	// SMTP 서버 통신 체크
	conn, err := net.DialTimeout("tcp", smtpserver, 3*time.Second)
	if err != nil {
		log.Println("[ERROR] Fail to connect smtp server : ", err, conn)
		return "", "", "", "", err
	}

	return id, pwd, host, port, nil
}
