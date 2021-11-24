package gorm

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// gorm.Model 的定义
type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func pgsqlOpen(host string, port int, dbname string, user string, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return sql.Open("postgres", connStr)
}

func Open(host string, port int, dbname string, user string, password string) (*gorm.DB, error) {
	sqlDB, err := pgsqlOpen(host, port, dbname, user, password)
	if err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	return gormDB, err
}
