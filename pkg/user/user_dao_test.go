package user

import (
	"fmt"
	"github.com/yddeng/gim/config"
	"github.com/yddeng/gim/internal/db"
	"testing"
	"time"
)

func init() {
	config.LoadConfig("../../config/config.toml")
}

func TestLoadUser(t *testing.T) {
	conf := config.GetConfig().DBConfig
	db.Open(conf.SqlType, conf.Host, conf.Port, conf.Database, conf.User, conf.Password)
	u, err := LoadUser("ydd")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(u)
}

func TestSetNxUser(t *testing.T) {
	conf := config.GetConfig().DBConfig
	db.Open(conf.SqlType, conf.Host, conf.Port, conf.Database, conf.User, conf.Password)

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
