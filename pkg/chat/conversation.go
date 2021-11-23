package gim

import (
	"fmt"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"time"
)

var (
	convID      = uint64(1)
	convsations = map[uint64]*Conversation{}
)

type Conversation struct {
	Type          pb.ConversationType // 对话类型
	ID            uint64              // 全局唯一ID
	Creator       string              // 对话创建者
	CreateAt      int64               // 创建时间戳 秒
	Members       []string            // 成员列表
	Name          string              // 会话名
	LastMessageAt *MessageEntity      // 最后一条消息
	Message       []*MessageEntity
}

func onCreateConversation(u *user.User, message *codec.Message) {
	req := message.GetData().(*pb.CreateConversationReq)
	fmt.Printf("onCreateConversation %v\n", req)

	c := &Conversation{
		ID:       convID,
		Creator:  u.ID,
		CreateAt: time.Now().Unix(),
		Members:  req.GetMembers(),
		Name:     req.GetName(),
	}
	convID++

	exist := false
	for _, uid := range c.Members {
		if uid == u.ID {
			exist = true
		} else {
			u2 := user.GetUserByID(uid)
			if u2 == nil {
				u.Reply(message.GetSeq(), &pb.CreateConversationResp{Code: pb.ErrCode_UserNotExist})
				return
			}
		}
	}

	if !exist {
		c.Members = append(c.Members, u.ID)
	}

	convsations[c.ID] = c

	conv := &pb.Conversation{
		ID:   c.ID,
		Name: c.Name,
	}

	u.Reply(message.GetSeq(), &pb.CreateConversationResp{
		Code: pb.ErrCode_OK,
		Conv: conv,
	})

	notify := &pb.NotifyInvited{
		Conv:   conv,
		InitBy: u.ID,
	}
	for _, uid := range c.Members {
		if uid != u.ID {
			u := user.GetUserByID(uid)
			if u != nil {
				u.OnNotifyInvited(notify)
			}
		}
	}

}

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
	gate.RegisterHandler(uint16(pb.CmdType_CmdCreateConversationReq), onCreateConversation)
	gate.RegisterHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
}
