package monitor

import "time"

// 메일 및 SMS 발송
func (t *C_monitor) Run__alert(_s_monotor__message, _s_monitor__hostname string) error {
	err := t.Init_check()
	if err != nil {
		return err
	}

	nowtime := time.Now().Format("2006-01-02")
	// 중복 알림 발송 제한을 발송 날짜 기록
	err = t.Update__date(nowtime, _s_monitor__hostname)
	if err != nil {
		return err
	}

	// DB 연락처, 메일 데이터 쿼리하여 변수 저장
	mail, number, err := t.Get__contact_info()
	if err != nil {
		return err
	}

	// 메일 발송 함수 실행
	err = Send__mail(_s_monotor__message, mail)
	if err != nil {
		return err
	}

	// 연락처 string 변환 후 SMS 발송
	for _, _number := range number {
		err = Send__sms(_s_monotor__message, _number)
		if err != nil {
			return err
		}
	}
	return nil
}
