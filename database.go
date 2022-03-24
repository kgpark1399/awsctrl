package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type C_database struct {
	db_conn *sql.DB

	s_db__type     string
	s_db__id       string
	s_db__pwd      string
	s_db__hostname string
	s_db__name     string
}

// DB Config init
func (t *C_database) Init_db(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) string {

	config := _s_db__id + ":" + _s_db__pwd + "@" + "tcp" + "(" + _s_db__hostname + ")" + "/" + _s_db__name
	return config

}

// DB Connection 시작
func (t *C_database) DB_conn(_s_db__type, _s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) error {
	var err error
	config := t.Init_db(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)
	t.db_conn, err = sql.Open(_s_db__type, config)
	if err != nil {
		return err
	}
	return nil
}

// DB Connection 종료
func (t *C_database) DB_close() error {
	return t.db_conn.Close()
}
