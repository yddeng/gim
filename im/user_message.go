package im

import (
	"errors"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"time"
)

func onUserLogin(sess dnet.Session, msg *codec.Message) {
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

	addUser(u)
	sess.SetContext(u)

	u.SendToClient(msg.GetSeq(), &pb.UserLoginResp{Code: pb.ErrCode_OK})
}

func onSessionClose(session dnet.Session, err error) {
	ctx := session.Context()
	if ctx != nil {
		u := ctx.(*User)
		log.Infof("onClose user(%s) %s. ", u.ID, err)
		u.sess.SetContext(nil)
		u.sess = nil
	}
}
