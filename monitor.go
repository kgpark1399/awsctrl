package monitor

type C_monitor struct {
	C_monitor__fuction
	n_rate int

	s_target__name     string
	s_target__protocol string
	s_target__url      string
}

// 모니터링 시스템 시작
func (t *C_monitor) Run_(n_rate int) error {

	err := t.Start__http_s_healthcheck(n_rate)
	if err != nil {
		return err
	}

	return nil
}

// 모니터링 대상 추가
func (t *C_monitor) Add__Monitoring_target(_s_target__name, s_target__protocol, _s_target__url string) error {

	err := t.Insert__Monitoring_target(_s_target__name, s_target__protocol, _s_target__url)
	if err != nil {
		return err
	}

	return nil
}

// 장애 알림 대상 추가
