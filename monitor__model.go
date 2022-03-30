package monitor

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type C_monitor__db struct {
	C_database

	s_target__name   string
	s_target__url    string
	s_target__status string

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

// DB status 컬럼 데이터 조회
func (t *C_monitor__db) Get__status() (status, url []string, err error) {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	var get C_monitor
	var gets []C_monitor

	// DB URL, STATUS 쿼리
	rows, err := t.db_conn.Query("SELECT url,status FROM target")
	if err != nil {
		fmt.Println(err)
	}

	// 쿼리 결과 변수 저장
	for rows.Next() {
		if err := rows.Scan(&get.s_monitor__url, &get.s_monitor__status); err != nil {
			fmt.Print(err)
		}
		gets = append(gets, get)
	}

	// 쿼리 결과를 받아줄 변수 생성
	var arrs_monitor__urls, arrs_monitor__status_grp []string

	// URL, STATUS 데이터 각각 변수에 입력
	for _, target := range gets {
		arrs_monitor__urls = append(arrs_monitor__urls, target.s_monitor__url)
		arrs_monitor__status_grp = append(arrs_monitor__status_grp, target.s_monitor__status)
	}

	// 결과 데이터 반환
	url = arrs_monitor__urls
	status = arrs_monitor__status_grp
	return url, status, nil
}

// DB URL 정보 호출 및 반환
func (t *C_monitor__db) Get__urls() (result []string, err error) {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	var website C_monitor
	var websites []C_monitor

	// DB URL STATUS 데이터 쿼리
	rows, err := t.db_conn.Query("SELECT url FROM target")
	if err != nil {
		fmt.Println(err)
	}

	// URL, STATUS 데이터 각각 변수에 입력
	for rows.Next() {
		if err := rows.Scan(&website.s_monitor__url); err != nil {
			fmt.Print(err)
		}

		websites = append(websites, website)
	}

	_t := C_monitor{}
	for _, target := range websites {
		_t.arrs_monitor__urls = append(_t.arrs_monitor__urls, target.s_monitor__url)
	}

	// 결과 반환
	result = _t.arrs_monitor__urls
	return result, nil

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

// DB status 컬럼 데이터 변경 (check 실패 시)
func (t *C_monitor__db) Change_status__false(_sUrl string) error {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	stmt, err := t.db_conn.Prepare("UPDATE target SET status=? WHERE url=?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec("false", _sUrl) //Placeholder 파라미터 순서대로 전달
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// DB status 컬럼 데이터 변경 (check 성공 시)
func (t *C_monitor__db) Change_status__true(_sUrl string) error {

	_bool, err := t.Init__monitor_db()
	if err != nil {
		fmt.Println(_bool, err)
	}

	stmt, err := t.db_conn.Prepare("UPDATE target SET status=? WHERE url=?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec("true", _sUrl) //Placeholder 파라미터 순서대로 전달
	if err != nil {
		fmt.Println(err)
	}
	return nil
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
