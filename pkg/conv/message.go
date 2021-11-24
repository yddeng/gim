package conv

import (
	"github.com/yddeng/gim/internal/protocol/pb"
)

type MessageEntity struct {
	ConversationID uint64 // 对话ID
	MessageInfo    *pb.MessageInfo
}
