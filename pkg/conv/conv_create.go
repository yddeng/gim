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
		Members:  make([]string, 0, len(req.GetMembers())),
		Name:     req.GetName(),
	}
	convID++

	c.Members = append(c.Members, u.ID)
	for _, id := range req.GetMembers() {
		// load 数据库
		if u2 := user.GetUser(id); u2 != nil {
			if u2 != u {
				c.Members = append(c.Members, id)
			}
		}

	}

	convMap[c.ID] = c

	conv := c.Pack()
	u.SendToClient(message.GetSeq(), &pb.CreateConversationResp{
		Code: pb.ErrCode_OK,
		Conv: conv,
	})

	notify := &pb.NotifyInvited{
		Conv:   conv,
		InitBy: u.ID,
	}
	c.Broadcast(notify, u.ID)
}

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdCreateConversationReq), onCreateConversation)
}
