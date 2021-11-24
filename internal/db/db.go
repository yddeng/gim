package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	sqlDB *sql.DB
)

func pgsqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return sql.Open("postgres", connStr)
}

func mysqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	return sql.Open("mysql", connStr)
}

func Open(sqlType string, host string, port int, dbname string, user string, password string) (err error) {
	if sqlType == "mysql" {
		sqlDB, err = mysqlOpen(host, port, dbname, user, password)
	} else {
		sqlDB, err = pgsqlOpen(host, port, dbname, user, password)
	}
	return
}

func DB() *sql.DB {
	return sqlDB
}
