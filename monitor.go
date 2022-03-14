package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type C_monitor struct {
	s_url   string
	s_rate  string
	s_ports string
	s_data  string
}

func (t *C_monitor) status_check(_s_url string, _s_rate int) {

	resp, err := http.Get(_s_url)
	if err != nil {
		// t.Send("Node01 Server Connection Fail", "Server Connection Fail")
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * time.Duration(_s_rate))

	go func() {
		for time := range ticker.C {
			fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				fmt.Println("HTTP Status is in the 2xx range", time)
			} else {
				// t.Send("Node01 HTTP Status Err", "HTTP Status Err")
				fmt.Println("Web Error", time)
			}
		}
	}()
	time.Sleep(time.Second * 5)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func (t *C_monitor) port_check(_s_url, _s_ports string) {
	address := net.JoinHostPort(_s_url, _s_ports)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
	} else if conn != nil {
		defer conn.Close()
		fmt.Printf("%s:%s is opened \n", _s_url, _s_ports)
	}
}

func (t *C_monitor) string_check(_s_url, _s_data string) {
	resp, err := http.Get(_s_url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// http body 값 변수 저장 및 string 변환
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)

}
