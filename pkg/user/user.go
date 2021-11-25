package user

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"time"
)

var (
	users = map[string]*User{}
)

func GetUser(id string) *User {
	u, ok := users[id]
	if !ok {
		var err error
		u, err = loadUser(id)
		if err != nil {
			log.Error(err)
		}
		if u != nil {
			users[id] = u
		}
	}
	return u
}

type User struct {
	ID       string
	CreateAt int64
	UpdateAt int64
	Extra    map[string]string // 附加属性
	sess     dnet.Session
	exist    bool // 缓存击穿
}

func (this *User) SendToClient(seq uint32, msg proto.Message) {
	disposeHook(this, msg)

	if this.sess == nil {
		return
	}
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

	u, err := loadUser(id)
	if err != nil {
		log.Error(err)
		_ = sess.Send(codec.NewMessage(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_Error}))
		return
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

	if err := setNxUser(u); err != nil {
		log.Error(err)
		_ = sess.Send(codec.NewMessage(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_Error}))
		return
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
		_ = setNxUser(u)
	}
}
