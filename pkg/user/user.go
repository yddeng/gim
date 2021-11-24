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
	ID       string
	CreateAt int64
	UpdateAt int64
	Attrs    map[string]string
	Convs    map[uint64]struct{}
	sess     dnet.Session
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
	ctx := session.Context()
	if ctx != nil {
		u := ctx.(*User)
		log.Infof("onClose user(%s) %s. ", u.ID, err)
		u.sess = nil
		delete(users, u.ID)
	}
}
