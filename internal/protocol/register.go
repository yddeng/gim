package protocol

import (
	"github.com/yddeng/gim/internal/codec/pb"
	proto "github.com/yddeng/gim/internal/protocol/pb"
)

func init() {
	pb.Register("im", &proto.UserLoginReq{}, uint16(proto.CmdType_CmdUserLoginReq))
	pb.Register("im", &proto.UserLoginResp{}, uint16(proto.CmdType_CmdUserLoginResp))

	pb.Register("im", &proto.CreateConversationReq{}, uint16(proto.CmdType_CmdCreateConversationReq))
	pb.Register("im", &proto.CreateConversationResp{}, uint16(proto.CmdType_CmdCreateConversationResp))
	pb.Register("im", &proto.NotifyInvited{}, uint16(proto.CmdType_CmdNotifyInvited))

	pb.Register("im", &proto.SendMessageReq{}, uint16(proto.CmdType_CmdSendMessageReq))
	pb.Register("im", &proto.SendMessageResp{}, uint16(proto.CmdType_CmdSendMessageResp))
	pb.Register("im", &proto.NotifyMessage{}, uint16(proto.CmdType_CmdNotifyMessage))
}
