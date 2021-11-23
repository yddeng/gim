package conv

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/user"
)

var (
	convID  = uint64(1)
	convMap = map[uint64]*Conversation{}
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

func (this *Conversation) Broadcast(msg proto.Message, except map[string]struct{}) {
	for _, id := range this.Members {
		if _, ok := except[id]; !ok {
			u := user.GetUser(id)
			if u != nil {
				u.SendToClient(0, msg)
			}
		}
	}
}
