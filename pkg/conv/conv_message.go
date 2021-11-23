package conv

import (
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"time"
)

func onSendMessage(u *user.User, message *codec.Message) {
	req := message.GetData().(*pb.SendMessageReq)

	c := convsations[req.GetConvID()]
	if c == nil {
		u.Reply(message.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	exist := false
	for _, id := range c.Members {
		if id == u.ID {
			exist = true
		}
	}

	if !exist {
		u.Reply(message.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	messageID := uint64(1)
	if c.LastMessageAt != nil {
		messageID = c.LastMessageAt.MessageID + 1
	}

	m := &MessageEntity{
		MessageID:      messageID,
		UserID:         u.ID,
		ConversationID: c.ID,
		Message:        req.GetMsg(),
		CreateAt:       time.Now().Unix(),
	}

	c.LastMessageAt = m
	c.Message = append(c.Message, m)

	conv := &pb.Conversation{
		ID:   c.ID,
		Name: c.Name,
	}
	info := &pb.MessageInfo{
		Msg:      req.GetMsg(),
		UserID:   u.ID,
		CreateAt: m.CreateAt,
	}
	notifyMessage := &pb.NotifyMessage{
		Conv:     conv,
		MsgInfos: []*pb.MessageInfo{info},
	}

	for _, id := range c.Members {
		if u2 := user.GetUserByID(id); u2 != nil {
			u2.Reply(0, notifyMessage)
		}
	}

}

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
}
