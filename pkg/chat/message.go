package gim

import "github.com/yddeng/gim/internal/protocol"

type MessageEntity struct {
	MessageID      uint64 // 消息ID
	UserID         string // 用户ID
	ConversationID uint64 // 对话ID
	Message        *protocol.Message
	CreateAt       int64 // 时间
}
