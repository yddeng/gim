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
	ID         string
	CreateAt   int64
	ConvStates map[uint64]*ConversationState
	sess       dnet.Session
}

const (
	stateWaitActive = iota
	stateActive
	stateNotify
	stateRemove
)

type ConversationState struct {
	ConversationID uint64
	LastReadAt     uint64 // 最后阅读的消息ID
	State          int    // 状态
}

func (this *User) Reply(seq uint32, msg proto.Message) {
	this.sess.Send(codec.NewMessage(seq, msg))
}

func (this *User) Tick() {

}

func (this *User) OnNotifyInvited(notify *protocol.NotifyInvited) {
	state := &ConversationState{
		ConversationID: notify.GetConv().GetID(),
		LastReadAt:     0,
		State:          stateWaitActive,
	}

	this.ConvStates[state.ConversationID] = state
	this.Reply(0, notify)
}

func (this *User) OnNotifyMessage(convID uint64) {

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
