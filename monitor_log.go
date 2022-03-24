package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type C_monitor_log struct {
	s_file__name string
	s_file__path string
}

func (t *C_monitor_log) Set_monitorlog(_s_file__path, _s_file__name string) os.File {

	t.Init__log(_s_file__path, _s_file__name)
	file, err := os.OpenFile(_s_file__name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	return *file
}

func (t *C_monitor_log) Init__log(_s_file__path, _s_file__name string) string {

	file := _s_file__path + _s_file__name
	return file
}
