package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func pgsqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return sql.Open("postgres", connStr)
}

func mysqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	return sql.Open("mysql", connStr)
}

func sqlOpen(sqlType string, host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	if sqlType == "mysql" {
		return mysqlOpen(host, port, dbname, user, password)
	} else {
		return pgsqlOpen(host, port, dbname, user, password)
	}
}

type Client struct {
	db      *sql.DB
	sqlType string
}

func NewClient(sqlType string, host string, port int, dbname string, user string, password string) (c *Client, err error) {
	c = new(Client)
	c.sqlType = sqlType
	c.db, err = sqlOpen(sqlType, host, port, dbname, user, password)
	if err != nil {
		return
	}

	err = c.db.Ping()
	if err != nil {
		return
	}
	return
}
