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

	s_protocol    string
	s_url         string
	s_data        string
	s_use_compare string
	s_alert_date  string

	arrs_contact__number []string
	arrs_protocol        []string
	arrs_urls            []string
	arrs_data            []string
	arrs_use_compare     []string
	arrs_alert_date      []string
}

func (t *C_monitor__db) Init() error {
	var err error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 모니터링 대상 정보 호출
func (t *C_monitor__db) Get__monitoring_target() (protocol, url, data, use_compare, alert []string, err error) {

	err = t.Init()
	if err != nil {
		return
	}

	var website C_monitor__db
	var websites []C_monitor__db

	// DB URL STATUS 데이터 쿼리
	rows, err := t.db_conn.Query("SELECT protocol,url,data,use__compare,alert FROM target")
	if err != nil {
		log.Println("[ERROR] Failed to database query : ", err)
		return

	}

	// URL, STATUS 데이터 각각 변수에 입력
	for rows.Next() {
		if err = rows.Scan(&website.s_protocol, &website.s_url, &website.s_data, &website.s_use_compare, &website.s_alert_date); err != nil {
			log.Println("[ERROR] Failed to database rows scane : ", err)
			return
		}

		websites = append(websites, website)
	}

	var arrs_protocol, arrs_urls, arrs_data, arrs_use_compare, arrs_alert_date []string

	for _, target := range websites {
		arrs_protocol = append(arrs_protocol, target.s_protocol)
		arrs_urls = append(arrs_urls, target.s_url)
		arrs_data = append(arrs_data, target.s_data)
		arrs_use_compare = append(arrs_use_compare, target.s_use_compare)
		arrs_alert_date = append(arrs_alert_date, target.s_alert_date)

	}

	return arrs_protocol, arrs_urls, arrs_data, arrs_use_compare, arrs_alert_date, nil
}

// 장애 알림 대상 호출
func (t *C_monitor__db) Get__Alert_Notification_target() (mail, number []string, err error) {

	err = t.Init()
	if err != nil {
		return
	}

	var contact C_monitor__db
	var contacts []C_monitor__db

	// DB URL, STATUS 쿼리
	rows, err := t.db_conn.Query("SELECT mail,mobile FROM contact")
	if err != nil {
		log.Println("[ERROR] Failed to database query : ", err)
		return
	}

	// 쿼리 결과 변수 저장
	for rows.Next() {
		if err = rows.Scan(&contact.s_contact__mail, &contact.s_contact__number); err != nil {
			log.Println("[ERROR] Failed to database rows scane : ", err)
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

// 모니터링 대상 추가
func (t *C_monitor__db) Insert__Monitoring_target(_s_target__name, _s_target__url, _s_target__status string) error {

	err := t.Init()
	if err != nil {
		return err
	}

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO target (name,url,status) VALUES (?, ?, ?)", _s_target__name, _s_target__url, _s_target__status)
	if err != nil {
		log.Println("[ERROR] Failed to database insert : ", err)
		return err
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		log.Println("1 row inserted.")
	}
	return nil
}

// 장애 알림 대상 추가
func (t *C_monitor__db) Insert__alert_notification_target(_s_contact__name, _s_contact__mail, _s_contact__number string) error {

	err := t.Init()
	if err != nil {
		return err
	}

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO contact (user,mail,mobile) VALUES (?, ?, ?)", _s_contact__name, _s_contact__mail, _s_contact__number)
	if err != nil {
		log.Println("[ERROR] Failed to database insert : ", err)
		return err
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		log.Println("1 row inserted.")
	}
	return nil
}

// 장애 알림 중복 발송 제한을 위한 알림 발송 기록
func (t *C_monitor__db) Update__alert_date(_s_url__alert_date, _s_hostname string) error {

	err := t.Init()
	if err != nil {
		return err
	}

	stmt, err := t.db_conn.Prepare("UPDATE target SET alert=? WHERE url=?")
	if err != nil {
		log.Println("[ERROR] Failed to database update : ", err)
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
