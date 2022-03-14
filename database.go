package main

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type C_database struct {
	s_dbid   string
	s_dbpwd  string
	s_addr   string
	s_dbname string
}

func (t *C_database) GetConnector(_s_dbid, _s_dbpwd, _s_addr, _s_dbname string) *sql.DB {
	cfg := mysql.Config{
		User:                 _s_dbid,
		Passwd:               _s_dbpwd,
		Net:                  "tcp",
		Addr:                 _s_addr,
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               _s_dbname,
	}
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		panic(err)
	}
	db := sql.OpenDB(connector)
	return db
}
