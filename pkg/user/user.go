package user

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"time"
)

var (
	userMap = map[string]*User{}
)

func GetUser(id string) *User {
	u, ok := userMap[id]
	if !ok {
		var err error
		u, err = loadUser(id)
		if err != nil {
			log.Error(err)
			return nil
		}
		if u != nil {
			userMap[id] = u
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

	if err := setNxUser(u); err != nil {
		log.Error(err)
		_ = sess.Send(codec.NewMessage(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_Error}))
		return
	}

	userMap[id] = u
	sess.SetContext(u)

	u.SendToClient(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_OK})
}

func OnClose(session dnet.Session, err error) {
	ctx := session.Context()
	if ctx != nil {
		u := ctx.(*User)
		log.Infof("onClose user(%s) %s. ", u.ID, err)
		u.sess.SetContext(nil)
		u.sess = nil
		delete(userMap, u.ID)
		_ = setNxUser(u)
	}
}
