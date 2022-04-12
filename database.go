package monitor

import (
	"database/sql"
	"log"

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
func (t *C_database) Init(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) (string, error) {
	var err error
	config := _s_db__id + ":" + _s_db__pwd + "@" + "tcp" + "(" + _s_db__hostname + ")" + "/" + _s_db__name
	if err != nil {
		log.Println(err)
		return "", err
	}
	return config, nil
}

// DB Connection 시작
func (t *C_database) Db__conn(_s_db__type, _s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name string) error {

	config, err := t.Init(_s_db__id, _s_db__pwd, _s_db__hostname, _s_db__name)
	if err != nil {
		return err
	}

	t.db_conn, err = sql.Open(_s_db__type, config)
	if err != nil {
		log.Println("[ERROR] Failed to database conn : ", err)
		return err
	}

	return nil
}

// DB Connection 종료
func (t *C_database) Db__close() error {
	return t.db_conn.Close()
}
