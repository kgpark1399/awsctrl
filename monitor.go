package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type C_monitor struct {
	sUrl    string
	sPorts  string
	sData   string
	sStatus string

	iRate int
}

func Run_checkUrl(_sUrl string, _iRate int) {
	cMonitor := New_C_monitor()
	cMonitor.checkUrl(_sUrl, _iRate)

}

func New_C_monitor() *C_monitor {
	c := &C_monitor{}
	return c
}

func (t *C_monitor) checkUrl(_sUrl string, _iRate int) {
	resp, err := http.Get(_sUrl)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * time.Duration(_iRate))

	go func() {
		for time := range ticker.C {
			fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				fmt.Println("HTTP Status is in the 2xx range", time)
			} else {
				fmt.Println("Web Error", time)
			}
		}
	}()
	time.Sleep(time.Second * 5)
	ticker.Stop()
	fmt.Println("Ticker stopped")

}
