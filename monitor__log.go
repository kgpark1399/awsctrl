package monitor

import (
	"fmt"
	"io"
	"log"
	"os"
)

type C_monitor__log struct {
	s_file__name string
}

func (t *C_monitor__log) Set_monitorlog(_s_file__name string) (os.File, error) {

	_bool, err := t.Init__log()
	if err != nil {
		fmt.Println(_bool, err)
	}

	file, err := os.OpenFile(_s_file__name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	return *file, nil
}

func (t *C_monitor__log) Init__log() (bool, error) {

	var err error
	if err != nil {
		return false, err
	}
	return true, nil
}
