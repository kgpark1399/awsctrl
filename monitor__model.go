package monitor

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor__db struct {
	C_database

	s_target__name string
	s_target__url  string

	s_contact__name   string
	s_contact__mail   string
	s_contact__number string

	arrs_contact__number []string
}

func (t *C_monitor__db) Init__monitor_db() (bool, error) {
	var err error
	if err != nil {
		return false, err
	}
	return true, nil
}

// DB URL 정보 호출 및 반환
func (t *C_monitor__db) Get__target_info() (url, data []string, use__ssl, use__compare, count []int, err error) {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	var website C_monitor
	var websites []C_monitor

	// DB URL STATUS 데이터 쿼리
	rows, err := t.db_conn.Query("SELECT url,data,use__ssl,use__compare,alert FROM target")
	if err != nil {
		fmt.Println(err)
	}

	// URL, STATUS 데이터 각각 변수에 입력
	for rows.Next() {
		if err := rows.Scan(&website.s_monitor__url, &website.s_monitor__data, &website.s_monitor__use__ssl, &website.s_monitor__use__compare, &website.n_monitor__alert_count); err != nil {
			fmt.Print(err)
		}

		websites = append(websites, website)
	}

	_t := C_monitor{}

	for _, target := range websites {
		_t.arrs_monitor__urls = append(_t.arrs_monitor__urls, target.s_monitor__url)
		_t.arrs_monitor__data = append(_t.arrs_monitor__data, target.s_monitor__data)
		_t.arrs_monitor__use__ssl = append(_t.arrs_monitor__use__ssl, target.s_monitor__use__ssl)
		_t.arrs_monitor__use__compare = append(_t.arrs_monitor__use__compare, target.s_monitor__use__compare)
		_t.arrn_monitor__alert = append(_t.arrn_monitor__alert, target.n_monitor__alert_count)

	}
	// 모니터링 대상 URL
	url = _t.arrs_monitor__urls
	// 문자열 비교 데이터
	data = _t.arrs_monitor__data
	// ssl 사용 여부
	use__ssl = _t.arrs_monitor__use__ssl
	// 문자열 비교 사용 여부
	use__compare = _t.arrs_monitor__use__compare
	// 알림 발송 여부 데이터
	count = _t.arrn_monitor__alert
	return url, data, use__ssl, use__compare, count, nil
}

// DB 관리자 정보 가져오기
func (t *C_monitor__db) Get__contact_info() (mail, number []string, err error) {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	var contact C_monitor__db
	var contacts []C_monitor__db

	// DB URL, STATUS 쿼리
	rows, err := t.db_conn.Query("SELECT mail,mobile FROM contact")

	if err != nil {
		fmt.Println(err)
	}

	// 쿼리 결과 변수 저장
	for rows.Next() {
		if err := rows.Scan(&contact.s_contact__mail, &contact.s_contact__number); err != nil {
			fmt.Println(err)
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

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO target (name,url,status) VALUES (?, ?, ?)", _s_target__name, _s_target__url, _s_target__status)
	if err != nil {
		fmt.Println(err)
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		fmt.Println("1 row inserted.")
	}
	return nil
}

// DB URL 담당자(Contact) 정보 추가
func (t *C_monitor__db) Create__contact(_s_contact__name, _s_contact__mail, _s_contact__number string) error {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	// INSERT 문 실행
	result, err := t.db_conn.Exec("INSERT INTO contact (user,mail,mobile) VALUES (?, ?, ?)", _s_contact__name, _s_contact__mail, _s_contact__number)
	if err != nil {
		fmt.Println(err)
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if n == 1 {
		fmt.Println("1 row inserted.")
	}
	return nil
}

// DB status 컬럼 데이터 변경 (check 성공 시)
func (t *C_monitor__db) Change_alert_count(_sUrl string) error {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	stmt, err := t.db_conn.Prepare("UPDATE target SET alert=? WHERE url=?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec(1, _sUrl) //Placeholder 파라미터 순서대로 전달
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
