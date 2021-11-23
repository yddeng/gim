package conv

import (
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
	"time"
)

func onCreateConversation(u *user.User, message *codec.Message) {
	req := message.GetData().(*pb.CreateConversationReq)
	log.Debugf("onCreateConversation %v", req)

	nowUnix := time.Now().Unix()
	c := &Conversation{
		ID:       convID,
		Creator:  u.ID,
		CreateAt: nowUnix,
		Members:  req.GetMembers(),
		Name:     req.GetName(),
	}
	convID++

	exist := false
	for _, uid := range c.Members {
		if uid == u.ID {
			exist = true
		} else {
			u2 := user.GetUser(uid)
			if u2 == nil {
				u.SendToClient(message.GetSeq(), &pb.CreateConversationResp{Code: pb.ErrCode_UserNotExist})
				return
			}
		}
	}

	if !exist {
		c.Members = append(c.Members, u.ID)
	}

	convMap[c.ID] = c

	conv := &pb.Conversation{
		ID:   c.ID,
		Name: c.Name,
	}

	u.SendToClient(message.GetSeq(), &pb.CreateConversationResp{
		Code: pb.ErrCode_OK,
		Conv: conv,
	})

	notify := &pb.NotifyInvited{
		Conv:   conv,
		InitBy: u.ID,
	}
	c.Broadcast(notify, map[string]struct{}{u.ID: {}})
}

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdCreateConversationReq), onCreateConversation)
}
