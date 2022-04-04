package monitor

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 기능이 없어도 init 을 넣고, false 와 err 를 호출하도록

type C_monitor struct {
	C_monitor__db
	C_monitor__log

	s_monitor__url  string
	s_monitor__name string
	s_monitor__data string

	n_monitor__rate         int
	n_monitor__expired_days int

	arrs_monitor__urls []string
	arrs_monitor__data []string
}

func (t *C_monitor) Init_check() (bool, error) {
	var err error
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	return true, nil
}

// URL HTTP 상태 체크 실행
func (t *C_monitor) Run__Monitor(_n_monitor__rate int) error {

	_, err := t.Init_check()
	if err != nil {
		return err
	}

	// DB에서 모니터링 대상 URL 호출
	target_url, target_data, err := t.Get__urls()
	if err != nil {
		return err
	}

	// 반복 시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_monitor__rate))

	for range ticker.C {

		// 모니터링 대상 URL arr_string > string 변환
		for i, url := range target_url {
			resp, err := http.Get(url)

			// HTTP 접속 오류 및 Status code 400 이상이면 오류
			if err != nil || resp.StatusCode >= 400 {
				message := "URL :" + url + "STATUS : ERR"

				// 장애 알림 발송(Mail,SMS)
				t.Run__alert(message)

				//로그 찍기
				log.Println(message)

			} else {
				// HTTP 접속 정상의 경우, 로그 찍기
				log.Println("URL :", url, ", STATUS :", resp.Status)

				// URL 대상 HTTP Body 값 호출
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				// arr_string > string 변환
				s_body := string(body)

				//string 변환 된 HTTP Body의 불필요 줄바꿈 제거
				_body := strings.TrimRight(s_body, "\r\n")

				// DB의 URL Data와 HTTP Body 값 문자열 비교
				if strings.EqualFold(_body, target_data[i]) {
					fmt.Print()
				} else {
					message := "URL :" + url + ", String Compare Err"
					t.Run__alert(message)
				}
			}
		}
	}
	return nil
}

// 메일 및 SMS 발송
func (t *C_monitor) Run__alert(_s_message string) error {
	var err error

	_, err = t.Init_check()
	if err != nil {
		return err
	}

	c_sendmail := C_Sendmail{}

	// DB 연락처, 메일 데이터 쿼리하여 변수 저장
	mail, number, err := t.Get__contact_info()
	if err != nil {
		return err
	}

	// 메일 발송 함수 실행
	err = c_sendmail.Send_mail(_s_message, mail)
	if err != nil {
		return err
	}

	// 연락처 string 변환 후 SMS 발송
	for _, _number := range number {
		err = Send_sns(_s_message, _number)
		if err != nil {
			return err
		}
	}
	return nil
}
