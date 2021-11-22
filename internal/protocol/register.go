package protocol

import "github.com/yddeng/gim/internal/codec/pb"

func init() {
	pb.Register("im", &CreateConversationReq{}, 1001)
	pb.Register("im", &CreateConversationResp{}, 1002)
	pb.Register("im", &NotifyInvited{}, 1003)

	pb.Register("im", &SendMessageReq{}, 1011)
	pb.Register("im", &SendMessageResp{}, 1012)
	pb.Register("im", &NotifyMessage{}, 1013)
}
