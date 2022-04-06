package monitor

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor struct {
	C_monitor__db
	C_monitor__log

	s_monitor__url  string
	s_monitor__data string

	s_monitor__use__ssl     int
	s_monitor__use__compare int
	n_monitor__alert_count  int
	n_monitor__rate         int

	arrs_monitor__urls         []string
	arrs_monitor__data         []string
	arrs_monitor__use__ssl     []int
	arrs_monitor__use__compare []int
	arrn_monitor__alert        []int
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

	// 반복 시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_monitor__rate))

	for range ticker.C {

		// DB에서 모니터링 대상 URL 호출
		target__url, target__data, target__use__ssl, target__use__compare, target__alert, err := t.Get__target_info()
		if err != nil {
			return err
		}

		// 모니터링 대상의 HTTPS 사용 유무 판단 후 모니터링 시작
		for i, url := range target__url {
			if target__use__ssl[i] == 0 {
				url__http := "http://" + url
				err = t.Run__url_check(url__http, url, target__data[i], target__use__compare[i], target__alert[i])
				if nil != err {
					return err
				}
			} else {
				url__https := "https://" + url
				url__https_port := url + ":443"
				err = t.Run__url_check(url__https, url, target__data[i], target__use__compare[i], target__alert[i])
				if err != nil {
					return err
				}
				err = t.Run__sslcheck(url__https_port, url, target__alert[i])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *C_monitor) Run__sslcheck(_s_url, _s_url_default string, _s_alert int) error {

	// SSL 인증서 유효성 체크
	conn, err := tls.Dial("tcp", _s_url, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// SSL 인증서와 호스트네임 비교
	err = conn.VerifyHostname(_s_url_default)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// SSL 인증서 만료 체크
	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter

	now := time.Now()
	before_month := expiry.AddDate(0, -1, 0)

	// 인증서 만료 한달 전 알림
	if before_month.Before(now) {
		message := "SSL Certificate Error , URL : " + _s_url
		log.Println("SSL Certi Error, now : ", now, ", befor 1m expiry :", before_month)

		if _s_alert == 0 {
			err = t.Run__alert(message, _s_url_default)
			if err != nil {
				return err
			} else {
				fmt.Print()
			}
		}

	} else {
		log.Println("SSL Certi OK, now : ", now, ", befor 1m expiry :", before_month)
	}

	return nil
}

func (t *C_monitor) Run__url_check(_s_url, _s_url_default, _s_data string, _s_compare, _s_alert int) error {

	// 모니터링 대상 URL http Get
	resp, err := http.Get(_s_url)
	if err != nil || resp.StatusCode >= 400 {
		message := "URL :" + _s_url + ", STATUS : ERR"

		// 장애 알림 중복 발송 체크 후 알림 발송
		if _s_alert == 0 {
			err = t.Run__alert(message, _s_url_default)
			if err != nil {
				return err
			} else {
				fmt.Print()
			}
		}
		//로그 찍기
		log.Println(message, _s_url)

	} else {
		// HTTP 접속 정상 로그 찍기
		log.Println("URL :", _s_url, ", STATUS :", resp.Status)

		if _s_compare == 1 {
			// http body 불러오기
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			// arr_string > string 변환
			s_body := string(body)

			//string 변환 된 HTTP Body의 불필요 줄바꿈 제거
			_body := strings.TrimRight(s_body, "\r\n")

			// DB의 URL Data와 HTTP Body 값 문자열 비교
			if strings.EqualFold(_body, _s_data) {
				fmt.Print()
			} else {
				message := "URL :" + _s_url + ", String Compare Err"
				log.Println(message)

				// 장애 알림 중복 발송 체크 후 알림 발송
				if _s_alert == 0 {
					err = t.Run__alert(message, _s_url)
					if err != nil {
						return err
					} else {
						fmt.Print()
					}
				}
			}
		}
	}
	return nil
}

// 메일 및 SMS 발송
func (t *C_monitor) Run__alert(_s_message, _s_url_default string) error {
	var err error

	// 중복 알림 발송 제한을 위한 카운트
	err = t.Change_alert_count(_s_url_default)
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
