package monitor

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

type C_list_target struct {

	// 모니터링 대상
	s_monitor__url       []string
	n_monitor__cycle_sec int

	// 장애 알림 대상
	s_notice__mail   []string
	s_notice__mobile []string
}

// ini 파일 존재 여부 체크
func (t *C_list_target) Init() (cfg *ini.File, err error) {

	cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Println("[ERROR] Not found config.ini file : ", err)
		return cfg, err
	}

	return cfg, nil
}

// 모니터링 대상 리스트 로딩
func (t *C_list_target) Get__monitor() (url []string, cycle int, err error) {

	cfg, err := t.Init()
	if err != nil {
		return
	}

	title := "monitor_target"
	t.s_monitor__url = cfg.Section(title).Key("S_monitor__url").Strings(",")
	t.n_monitor__cycle_sec = cfg.Section(title).Key("N_monitor__cycle_sec").MustInt()

	return t.s_monitor__url, t.n_monitor__cycle_sec, nil
}

// 장애 알림 대상 리스트 로딩
func (t *C_list_target) Get__alert_contact() (mail, mobile []string, err error) {

	cfg, err := t.Init()
	if err != nil {
		return
	}

	title := "alert_target"
	t.s_notice__mail = cfg.Section(title).Key("S_notice__mail").Strings(",")
	t.s_notice__mobile = cfg.Section(title).Key("S_notice__mobile").Strings(",")

	fmt.Println(t.s_notice__mail, t.s_notice__mobile)
	return t.s_notice__mail, t.s_notice__mobile, nil

}
