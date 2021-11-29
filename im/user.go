package im

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/im/pb"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
	"time"
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

	u, err := dbLoadUser(id)
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
	Groups   map[int64]*Member // 会话列表
	sess     dnet.Session
}

func (this *User) online() bool {
	return this.sess != nil
}

func (this *User) SendToClient(seq uint32, msg proto.Message) {
	if this.sess == nil {
		return
	}
	if err := this.sess.Send(NewMessage(seq, msg)); err != nil {
		this.sess.Close(err)
	}
}

func onUserLogin(sess dnet.Session, msg *Message) {
	log.Infof("onUserLogin %v", msg)
	req := msg.GetData().(*pb.UserLoginReq)

	id := req.GetID()
	ctx := sess.Context()
	if ctx != nil {
		_ = sess.Send(NewMessage(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_UserAlreadyLogin}))
		return
	}

	u := GetUser(id)
	if u != nil && u.sess != nil {
		u.sess.SetContext(nil)
		u.sess.Close(errors.New("session kicked. "))
		u.sess = nil
	}

	nowUnix := time.Now().Unix()
	if u == nil {
		u = &User{
			ID:       id,
			CreateAt: nowUnix,
		}
	}

	u.UpdateAt = nowUnix
	u.Extra = req.GetExtra()
	u.sess = sess

	if err := dbSetNxUser(u); err != nil {
		log.Error(err)
		_ = sess.Send(NewMessage(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_Error}))
		return
	}

	addUser(u)
	sess.SetContext(u)

	u.SendToClient(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_OK})
}

func dbLoadUser(key string) (*User, error) {
	sqlStr := `
SELECT * FROM "users" 
WHERE id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, key)
	log.Debug(sqlStatement)

	var u User
	var extra []byte
	rows, err := sqlDB.Query(sqlStatement)
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

func dbSetNxUser(u *User) error {
	sqlStatement := `
INSERT INTO "users" (id,create_at,update_at,extra)
VALUES($1, $2, $3, $4) 
ON conflict(id) DO 
UPDATE SET create_at = $2, update_at = $3, extra = $4;`
	smt, err := sqlDB.Prepare(sqlStatement)
	if err != nil {
		return err
	}

	extra, _ := json.Marshal(u.Extra)
	_, err = smt.Exec(u.ID, u.CreateAt, u.UpdateAt, extra)
	return err
}

func init() {
	registerUserHandler(uint16(pb.CmdType_CmdUserLoginReq), onUserLogin)
}
