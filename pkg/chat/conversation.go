package gim

import (
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol"
	"github.com/yddeng/gim/pkg/gate"
	"time"
)

var (
	convID      = uint64(1)
	convsations = map[uint64]*Conversation{}
)

type ConversationType int

const (
	ConversationType_Default   = iota // 普通对话
	ConversationType_Transient        // 临时对话
	ConversationType_System           // 系统对话
)

type Conversation struct {
	Type     ConversationType // 对话类型
	ID       uint64           // 全局唯一ID
	Creator  string           // 对话创建者
	CreateAt int64            // 创建时间戳 秒
	Members  []string         // 成员列表
	Name     string           // 会话名
}

func onCreateConversation(session dnet.Session, message *codec.Message) {
	user := sess2User[session]
	req := message.GetData().(*protocol.CreateConversationReq)

	c := &Conversation{
		ID:       convID,
		Creator:  user.ID,
		CreateAt: time.Now().Unix(),
		Members:  req.GetMembers(),
		Name:     req.GetName(),
	}
	convID++

	exist := false
	for _, uid := range c.Members {
		if uid == user.ID {
			exist = true
			break
		}
	}

	if !exist {
		c.Members = append(c.Members, user.ID)
	}

	convsations[c.ID] = c

	conv := &protocol.Conversation{
		ID:   c.ID,
		Name: c.Name,
	}

	user.Reply(message.GetSeq(), &protocol.CreateConversationResp{
		Ok:   true,
		Conv: conv,
	})

	notify := &protocol.NotifyInvited{
		Conv:   conv,
		InitBy: user.ID,
	}
	for _, uid := range c.Members {
		if uid != user.ID {
			u := getUser(uid)
			if u != nil {
				u.Reply(0, notify)
			}
		}
	}

}

func init() {
	gate.RegisterHandler(&protocol.CreateConversationReq{}, onCreateConversation)
}
