package user

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
)

var (
	users = map[string]*User{}
)

func GetUser(id string) *User {
	return users[id]
}

type User struct {
	ID         string
	Attrs      map[string]string
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

func (this *User) SendToClient(seq uint32, msg proto.Message) {
	if this.sess == nil {
		return
	}
	disposeHook(this, msg)
	if err := this.sess.Send(codec.NewMessage(seq, msg)); err != nil {
		this.sess.Close(err)
	}
}

func (this *User) OnNotifyInvited(notify *pb.NotifyInvited) {
	state := &ConversationState{
		ConversationID: notify.GetConv().GetID(),
		LastReadAt:     0,
		State:          stateWaitActive,
	}

	this.ConvStates[state.ConversationID] = state
	this.SendToClient(0, notify)
}

func (this *User) OnNotifyMessage(convID uint64) {

}

func OnUserLogin(sess dnet.Session, msg *codec.Message) {
	log.Infof("onUserLogin %v", msg)
	req := msg.GetData().(*pb.UserLoginReq)

	id := req.GetID()
	ctx := sess.Context()
	if ctx != nil {
		_ = sess.Send(codec.NewMessage(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_UserAlreadyLogin}))
		return
	}

	u := &User{
		ID:   id,
		sess: sess,
	}

	users[id] = u
	sess.SetContext(u)

	u.SendToClient(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_OK})
}

func OnClose(session dnet.Session, err error) {
	log.Infof("onClose %s. ", err)
	ctx := session.Context()
	if ctx != nil {
		u := ctx.(*User)
		u.sess.SetContext(nil)
		u.sess = nil
		delete(users, u.ID)
	}
}
