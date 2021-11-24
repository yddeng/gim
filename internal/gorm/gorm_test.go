package gorm

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

type User struct {
	//ID   uint `gorm:"primaryKey"`
	Name string `gorm:"primaryKey"`
	Age  int
}

func TestOpen(t *testing.T) {
	db, err := Open("localhost", 5432, "yidongdeng", "dbuser", "123456")
	if err != nil {
		panic(err)
	}

	u := &User{
		Name: "ydd",
		Age:  22,
	}

	//result := db.Omit("created_at", "updated_at", "deleted_at").Create(&u)
	//result := db.Create(&u)

	sqlStr := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(&u)
	})
	fmt.Println(sqlStr)

	result := db.Create(&u)
	fmt.Println(result.Error)
}
