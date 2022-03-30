package monitor

import (
	"fmt"
	"net/smtp"
)

type C_Sendmail struct {
	// 메일 서버 연동 셋팅
	// s_mail__id   string
	// s_mail__pwd  string
	// s_mail__host string

	// 메일 발송 옵션
	s_mail__from string
	s_mail__body string

	arrs_mail__to []string
}

// 메일 로그인 정보 입력
func (t *C_Sendmail) Init__mail() (smtp.Auth, error) {
	var err error

	// 메일 로그인 정보
	auth := smtp.PlainAuth("", "kgpark@devtools.kr", "bmtvoyrsqimessqi", "smtp.gmail.com")

	if err != nil {
		fmt.Print(err)
		return auth, err
	}

	return auth, nil
}

// 메일 발송 함수
func (t *C_Sendmail) Send_mail(_s_mail__body string, arrs_mail__to []string) error {

	auth, err := t.Init__mail()
	if err != nil {
		return err
	}

	from := "kgpark@devtools.kr"
	to := arrs_mail__to

	// 메시지 작성
	headerSubject := "Subject: [Alert] Server Error\r\n"
	headerBlank := "\r\n"
	body := _s_mail__body
	msg := []byte(headerSubject + headerBlank + body)

	// 메일 발송
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return err
}
