package monitor

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type C_telegram struct {
	n_chat__id      string
	s_bot__id       string
	s_chat__message string
}

func Send_telegram(_n_chat__id, s_bot__id, _s_chat__message string) (string, error) {

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
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
