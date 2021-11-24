package user

import (
	"fmt"
	"github.com/yddeng/gim/internal/db"
	"testing"
	"time"
)

func TestLoadUser(t *testing.T) {
	db.Open("pgsql", "localhost", 5432, "yidongdeng", "dbuser", "123456")
	u, err := LoadUser("ydd")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(u)
}

func TestSetNxUser(t *testing.T) {
	db.Open("pgsql", "localhost", 5432, "yidongdeng", "dbuser", "123456")

	u := &User{
		ID:       "ydd",
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
		Attrs:    map[string]string{"name": "ydd", "age": "24"},
		Convs:    map[uint64]struct{}{123123: {}},
	}

	if err := SetNxUser(u); err != nil {
		t.Error(err)
	}

	u, err := LoadUser("ydd")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(u)
}
