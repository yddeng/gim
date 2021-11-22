package user

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol"
	"time"
)

var (
	sess2User = map[dnet.Session]*User{}
	users     = map[string]*User{}
)

func GetUserBySession(sess dnet.Session) *User {
	return sess2User[sess]
}

func GetUserByID(id string) *User {
	return users[id]
}

type User struct {
	ID       string
	CreateAt int64
	sess     dnet.Session
}

func (this *User) Reply(seq uint32, msg proto.Message) {
	this.sess.Send(codec.NewMessage(seq, msg))
}

func OnUserLogin(sess dnet.Session, msg *codec.Message) {
	fmt.Printf("onUserLogin %v\n", msg)
	req := msg.GetData().(*protocol.UserLoginReq)

	id := req.GetID()
	u := GetUserByID(id)
	if u != nil {
		sess.Send(codec.NewMessage(msg.GetSeq(), &protocol.UserLoginResp{Ok: false}))
		return
	}

	u = &User{
		ID:       id,
		CreateAt: time.Now().Unix(),
		sess:     sess,
	}

	users[id] = u
	sess.SetContext(u)

	u.Reply(msg.GetSeq(), &protocol.UserLoginResp{Ok: true})
}

func OnClose(session dnet.Session, err error) {
	fmt.Println("onClose", err)
	ctx := session.Context()
	if ctx != nil {
		u := ctx.(*User)
		delete(users, u.ID)
	}
}
