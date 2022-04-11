package monitor

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor__db struct {
	C_database

	s_contact__name   string
	s_contact__mail   string
	s_contact__number string

	arrs_contact__number []string
	arrs_protocol        []string
	arrs_urls            []string
	arrs_data            []string
	arrs_use_compare     []string
	arrs_alert_date      []string
}

func (t *C_monitor__db) Init__monitor_db() error {
	var err error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DB URL 정보 호출 및 반환
func (t *C_monitor__db) Get__target_info() (protocol, url, data, use_compare, alert []string, err error) {

	err = t.Init__monitor_db()
	if err != nil {
		return
	}

	var website C_monitor
	var websites []C_monitor

	// DB URL STATUS 데이터 쿼리
	rows, err := t.db_conn.Query("SELECT protocol,url,data,use__compare,alert FROM target")
	if err != nil {
		log.Println("sql 데이터 쿼리 오류 : ", err)
		return

	}

	// URL, STATUS 데이터 각각 변수에 입력
	for rows.Next() {
		if err = rows.Scan(&website.s_protocol, &website.s_url, &website.s_data, &website.s_use_compare, &website.s_alert_date); err != nil {
			log.Println(err)
			return
		}

		websites = append(websites, website)
	}

	_t := C_monitor{}

	for _, target := range websites {
		_t.arrs_protocol = append(_t.arrs_protocol, target.s_protocol)
		_t.arrs_urls = append(_t.arrs_urls, target.s_url)
		_t.arrs_data = append(_t.arrs_data, target.s_data)
		_t.arrs_use_compare = append(_t.arrs_use_compare, target.s_use_compare)
		_t.arrs_alert_date = append(_t.arrs_alert_date, target.s_alert_date)

	}
	protocol = _t.arrs_protocol
	// 모니터링 대상 URL
	url = _t.arrs_urls
	// 문자열 비교 데이터
	data = _t.arrs_data
	// 문자열 비교 사용 여부
	use_compare = _t.arrs_use_compare
	// 알림 발송 여부 데이터
	alert = _t.arrs_alert_date
	return protocol, url, data, use_compare, alert, nil
}

// DB 관리자 정보 가져오기
func (t *C_monitor__db) Get__contact_info() (mail, number []string, err error) {

	err = t.Init__monitor_db()
	if err != nil {
		return
	}

	var contact C_monitor__db
	var contacts []C_monitor__db

	// DB URL, STATUS 쿼리
	rows, err := t.db_conn.Query("SELECT mail,mobile FROM contact")
	if err != nil {
		log.Println("sql 데이터 쿼리 오류 : ", err)
		return
	}

	// 쿼리 결과 변수 저장
	for rows.Next() {
		if err = rows.Scan(&contact.s_contact__mail, &contact.s_contact__number); err != nil {
			log.Println(err)
			return
		}
		contacts = append(contacts, contact)
	}

	var arrs_contact__number, arrs_contact__mail []string

	for _, target := range contacts {
		arrs_contact__mail = append(arrs_contact__mail, target.s_contact__mail)
		arrs_contact__number = append(arrs_contact__number, target.s_contact__number)
	}

	mail = arrs_contact__mail
	number = arrs_contact__number
	return mail, number, nil
}

// DB ULR 대상(Target) 추가
func (t *C_monitor__db) Create__url(_s_target__name, _s_target__url, _s_target__status string) error {

	err := t.Init__monitor_db()
	if err != nil {
		return err
	}

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO target (name,url,status) VALUES (?, ?, ?)", _s_target__name, _s_target__url, _s_target__status)
	if err != nil {
		log.Println("sql 데이터 쿼리 오류 : ", err)
		return err
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		log.Println("1 row inserted.")
	}
	return nil
}

// DB URL 담당자(Contact) 정보 추가
func (t *C_monitor__db) Create__contact(_s_contact__name, _s_contact__mail, _s_contact__number string) error {

	err := t.Init__monitor_db()
	if err != nil {
		return err
	}

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO contact (user,mail,mobile) VALUES (?, ?, ?)", _s_contact__name, _s_contact__mail, _s_contact__number)
	if err != nil {
		log.Println("sql 데이터 쿼리 오류", err)
		return err
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		log.Println("1 row inserted.")
	}
	return nil
}

// 알림 중복 방지를 위한 alert 상태 변경
func (t *C_monitor__db) Update__date(_s_url__alert_date, _s_hostname string) error {

	err := t.Init__monitor_db()
	if err != nil {
		return err
	}

	stmt, err := t.db_conn.Prepare("UPDATE target SET alert=? WHERE url=?")
	if err != nil {
		log.Println("sql 데이터 업데이트 오류 : ", err)
		return err
	}
	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec(_s_url__alert_date, _s_hostname) //Placeholder 파라미터 순서대로 전달
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
