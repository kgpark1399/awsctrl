package monitor

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/ini.v1"
)

type C_notice__telegram struct {
	n_chat__id      string
	s_bot__id       string
	s_chat__message string
}

// 텔레그램 접속 정보 호출
func (t *C_notice__telegram) Init() error {

	//config.ini 파일 읽기
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return nil
	}

	title := "telegram"
	t.n_chat__id = cfg.Section(title).Key("N_chat__id").String()
	t.s_bot__id = cfg.Section(title).Key("S_bot__id").String()

	fmt.Println(t.n_chat__id, t.s_bot__id)
	return nil

}

func (t *C_notice__telegram) Send(_n_chat__id, s_bot__id, _s_chat__message string) error {

	log.Printf("Sending %s to chat_id: %s", _s_chat__message, _n_chat__id)
	telegram_api := "https://api.telegram.org/bot" + s_bot__id + "/sendMessage"
	response, err := http.PostForm(
		telegram_api,
		url.Values{
			"chat_id": {_n_chat__id},
			"text":    {_s_chat__message},
		})

	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return nil
}
