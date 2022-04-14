package monitor

import (
	"log"
	"net/http"
)

type C_http_status_code struct {
	c_list_target C_list_target

	s_url string
}

func (t *C_http_status_code) Init() {

}

// HTTP 상태 및 코드 확인
func (t *C_http_status_code) Get(_s_url string) (result bool, err error) {

	http_url := "http://" + _s_url
	resp, err := http.Get(http_url)
	if err != nil || resp.StatusCode >= 400 {
		log.Println(err)
		return false, err
	} else {
		log.Println("URL :", _s_url, ", Status code : ", resp.StatusCode, "Connecton OK")
	}

	return true, nil
}
