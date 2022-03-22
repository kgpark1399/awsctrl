package main

import (
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type C_monitor struct {
	C_db_config

	s_url  string
	s_name string

	n_status   int
	n_rate_min int

	arrs_urls       []string
	arrn_status_grp []int
}

type C_db_config struct {
	db_conn *sql.DB

	s_db__id       string
	s_db__pwd      string
	s_db__hostname string
	s_db__name     string
}

type C_sns struct {
	cfg aws.Config

	sTitle   string
	sRegion  string
	sAcid    string
	sAckey   string
	sSession string
	sTopic   string
}
