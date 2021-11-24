package conv

import (
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
	"time"
)

func onSendMessage(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.SendMessageReq)
	log.Debugf("onSendMessage %v", req)

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if inConv := c.HasUser(u.ID); !inConv {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	msgID := uint64(1)
	if c.LastMessageAt != nil {
		msgID = c.LastMessageAt.MessageInfo.GetMsgID() + 1
	}

	m := &MessageEntity{
		MessageInfo: &pb.MessageInfo{
			Msg:      req.GetMsg(),
			UserID:   u.ID,
			CreateAt: time.Now().Unix(),
			MsgID:    msgID,
		},
		ConversationID: c.ID,
	}

	c.LastMessageAt = m
	c.Message = append(c.Message, m)

	conv := c.Pack()
	u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_OK, Conv: conv})

	notifyMessage := &pb.NotifyMessage{
		Conv:     conv,
		MsgInfos: []*pb.MessageInfo{m.MessageInfo},
	}
	c.Broadcast(notifyMessage)
}

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
}
