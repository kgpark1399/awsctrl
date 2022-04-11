package main

import (
	"monitor/monitor"
)

func main() {

	t := monitor.C_monitor{}

	// config 설정
	t.Set__db_conn("mysql", "root", "pwd", "127.0.0.1:3306", "monitor")
	t.Set__log("monitor.txt")

	defer t.DB_close()

	t.Run__Monitor(10)

}
