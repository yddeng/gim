package conv

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/gim/pkg/util"
)

var (
	convID  = uint64(1)
	convMap = map[uint64]*Conversation{}
)

func GetConversation(convID uint64) *Conversation {
	return convMap[convID]
}

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

func (this *Conversation) Pack() *pb.Conversation {
	c := &pb.Conversation{
		Type: this.Type,
		ID:   this.ID,
		Name: this.Name,
	}

	if this.LastMessageAt != nil {
		c.LastMessageTimestamp = this.LastMessageAt.MessageInfo.GetCreateAt()
		c.LastMessageID = this.LastMessageAt.MessageInfo.GetMsgID()
	}

	return c
}

func (this *Conversation) Broadcast(msg proto.Message, except ...string) {
	for _, id := range this.Members {
		if has := util.HasString(id, except); !has {
			u := user.GetUser(id)
			if u != nil {
				u.SendToClient(0, msg)
			}
		}
	}
}

func (this *Conversation) HasUser(id string) bool {
	return util.HasString(id, this.Members)
}

func (this *Conversation) AddMember(ids []string) {
	this.Members = append(this.Members, ids...)
}

func (this *Conversation) RemoveMember(ids []string) {
	f := func(s string, m *[]string) {
		for i, v := range *m {
			if v == s {
				*m = append((*m)[:i], (*m)[i+1:]...)
				break
			}
		}
	}

	for _, id := range ids {
		f(id, &this.Members)
	}
}
