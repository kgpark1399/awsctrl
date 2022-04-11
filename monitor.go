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
	// DB, log 구조체 상속
	C_monitor__db
	C_monitor__log

	s_protocol    string
	s_url         string
	s_hostname    string
	s_data        string
	s_message     string
	s_use_compare string
	s_alert_date  string

	n_rate int
}

func (t *C_monitor) Init_check() error {
	var err error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// URL HTTP 상태 체크 실행
func (t *C_monitor) Run__Monitor(_n_rate int) error {

	err := t.Init_check()
	if err != nil {
		return err
	}

	// 반복 시간 설정
	ticker := time.NewTicker(time.Second * time.Duration(_n_rate))

	for range ticker.C {

		// DB에서 모니터링 대상 URL 호출
		target__protocol, target__url, target__data, target__use__compare, target__alert, err := t.Get__target_info()
		if err != nil {
			return err
		}

		// HTTP Stutus 체크
		for i, url := range target__url {
			url__protocol := target__protocol[i] + url
			url__port := url + ":443"
			err = t.Run__url_check(url__protocol, url, target__data[i], target__use__compare[i], target__alert[i])
			if nil != err {
				return err
			}

			// HTTPS 사용 시 인증서 유효성 및 만료기간 체크
			if target__protocol[i] == "https://" {
				err = t.Run__sslcheck(url__port, url, target__alert[i])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// SSL 인증서 유효성 및 만료일 체크
func (t *C_monitor) Run__sslcheck(_s_url, _s_hostname, _s_alert_date string) error {
	err := t.Init_check()
	if err != nil {
		return err
	}

	nowtime := time.Now().Format("2006-01-02")

	// SSL 인증서 유효성 체크
	conn, err := tls.Dial("tcp", _s_url, nil)
	if err != nil {
		log.Println("SSL 인증서 체크 오류 : ", err)
		return err
	}

	// SSL 인증서와 호스트네임 비교
	err = conn.VerifyHostname(_s_hostname)
	if err != nil {
		log.Println("SSL 인증서와 호스트네임 매칭 오류 : ", err)
		return err
	}

	// SSL 인증서 만료 체크
	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter

	now := time.Now()
	before_month := expiry.AddDate(0, -1, 0)

	// 인증서 만료 한달 전 알림
	if before_month.Before(now) {
		message := "SSL 인증서 만료 한달 전 입니다. , URL : " + _s_url
		log.Println("SSL Certi Error, now : ", now, ", befor 1m expiry :", before_month)

		if strings.EqualFold(_s_alert_date, nowtime) {
			fmt.Print()
		} else {
			err = t.Run__alert(message, _s_hostname)
			if err != nil {
				return err
			}
		}

	} else {
		log.Println("SSL Certi OK, now : ", now, ", befor 1m expiry :", before_month)
	}

	return nil
}

// HTTT/S 상태 및 문자열 체크
func (t *C_monitor) Run__url_check(_s_url, _s_hostname, _s_data, _s_use_compare, _s_alert_date string) error {

	err := t.Init_check()
	if err != nil {
		return err
	}

	nowtime := time.Now().Format("2006-01-02")

	// 모니터링 대상 URL http Get
	resp, err := http.Get(_s_url)
	if err != nil || resp.StatusCode >= 400 {
		message := "URL :" + _s_url + ", STATUS : ERR"

		// 장애 알림 중복 발송 체크 후 알림 발송
		if strings.EqualFold(_s_alert_date, nowtime) {
			fmt.Print()
		} else {
			err = t.Run__alert(message, _s_hostname)
			if err != nil {
				return err
			}
		}
		//로그 찍기
		log.Println(message, _s_url)

	} else {
		// HTTP 접속 정상 로그 찍기
		log.Println("URL :", _s_url, ", STATUS :", resp.Status)

		if _s_use_compare == "Y" {
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
				if strings.EqualFold(_s_alert_date, nowtime) {
					fmt.Print()
				} else {
					err = t.Run__alert(message, _s_hostname)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
