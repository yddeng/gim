package im

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
)

var userCache *lru.Cache = lru.New(5000)

type cacheUser struct {
	u *User
}

func GetUser(id string) *User {
	v, ok := userCache.Get(id)
	if ok {
		u := v.(*cacheUser)
		return u.u
	}

	u, err := loadUser(id)
	if err != nil {
		log.Error(err)
		return nil
	}

	userCache.Add(id, &cacheUser{u: u})
	return u
}

func addUser(u *User) {
	userCache.Add(u.ID, &cacheUser{u: u})
}

type User struct {
	ID       string
	CreateAt int64
	UpdateAt int64
	Extra    map[string]string // 附加属性
	Convs    map[int64]*Member // 会话列表
	sess     dnet.Session
}

func (this *User) SendToClient(seq uint32, msg proto.Message) {
	if this.sess == nil {
		return
	}
	if err := this.sess.Send(codec.NewMessage(seq, msg)); err != nil {
		this.sess.Close(err)
	}
}

func loadUser(key string) (*User, error) {
	sqlStr := `
SELECT * FROM "users" 
WHERE id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, key)
	log.Debug(sqlStatement)

	var u User
	var extra []byte
	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&u.ID, &u.CreateAt, &u.UpdateAt, &extra); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(extra, &u.Extra)
	return &u, nil
}

func setNxUser(u *User) error {
	sqlStatement := `
INSERT INTO "users" (id,create_at,update_at,extra)
VALUES($1, $2, $3, $4) 
ON conflict(id) DO 
UPDATE SET create_at = $2, update_at = $3, extra = $4;`
	smt, err := db.SqlDB.Prepare(sqlStatement)
	if err != nil {
		return err
	}

	extra, _ := json.Marshal(u.Extra)
	_, err = smt.Exec(u.ID, u.CreateAt, u.UpdateAt, extra)
	return err
}
