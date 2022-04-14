package monitor

import (
	"crypto/tls"
	"log"
	"time"
)

type C_ssl_cert struct {
	s_url      string
	s_hostName string
}

func (t *C_ssl_cert) Init() {
}

// SSL 인증서 유효성 만료일 체크
func (t *C_ssl_cert) Get__expiryPeriod(_s_url string) (result bool, err error) {

	url_port := _s_url + ":443"
	conn, err := tls.Dial("tcp", url_port, nil)
	if err != nil {
		log.Println(err)
		return false, err
	}

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter

	now := time.Now()
	before_month := expiry.AddDate(0, -1, 0)

	// 인증서 만료 한달 전 알림
	if before_month.Before(now) {
		log.Println(err)
		if err != nil {
			return false, err
		}
	} else {
		log.Println("SSL Certi OK, now : ", now, ", befor 1m expiry :", before_month)
	}

	return true, nil

}
