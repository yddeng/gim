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
	pb.Register("im", &proto.GetConversationUsersReq{}, uint16(proto.CmdType_CmdGetConversationUsersReq))
	pb.Register("im", &proto.GetConversationUsersResp{}, uint16(proto.CmdType_CmdGetConversationUsersResp))
	pb.Register("im", &proto.GetConversationListReq{}, uint16(proto.CmdType_CmdGetConversationListReq))
	pb.Register("im", &proto.GetConversationListResp{}, uint16(proto.CmdType_CmdGetConversationListResp))
	pb.Register("im", &proto.DissolveConversationReq{}, uint16(proto.CmdType_CmdDissolveConversationReq))
	pb.Register("im", &proto.DissolveConversationResp{}, uint16(proto.CmdType_CmdDissolveConversationResp))
	pb.Register("im", &proto.NotifyDissolveConversation{}, uint16(proto.CmdType_CmdNotifyDissolveConversation))

	pb.Register("im", &proto.AddMemberReq{}, uint16(proto.CmdType_CmdAddMemberReq))
	pb.Register("im", &proto.AddMemberResp{}, uint16(proto.CmdType_CmdAddMemberResp))
	pb.Register("im", &proto.RemoveMemberReq{}, uint16(proto.CmdType_CmdRemoveMemberReq))
	pb.Register("im", &proto.RemoveMemberResp{}, uint16(proto.CmdType_CmdRemoveMemberResp))
	pb.Register("im", &proto.JoinReq{}, uint16(proto.CmdType_CmdJoinReq))
	pb.Register("im", &proto.JoinResp{}, uint16(proto.CmdType_CmdJoinResp))
	pb.Register("im", &proto.QuitReq{}, uint16(proto.CmdType_CmdQuitReq))
	pb.Register("im", &proto.QuitResp{}, uint16(proto.CmdType_CmdQuitResp))
	pb.Register("im", &proto.NotifyMemberJoined{}, uint16(proto.CmdType_CmdNotifyMemberJoined))
	pb.Register("im", &proto.NotifyMemberLeft{}, uint16(proto.CmdType_CmdNotifyMemberLeft))
	pb.Register("im", &proto.NotifyInvited{}, uint16(proto.CmdType_CmdNotifyInvited))
	pb.Register("im", &proto.NotifyKicked{}, uint16(proto.CmdType_CmdNotifyKicked))

	pb.Register("im", &proto.SendMessageReq{}, uint16(proto.CmdType_CmdSendMessageReq))
	pb.Register("im", &proto.SendMessageResp{}, uint16(proto.CmdType_CmdSendMessageResp))
	pb.Register("im", &proto.NotifyMessage{}, uint16(proto.CmdType_CmdNotifyMessage))
	pb.Register("im", &proto.SyncMessageReq{}, uint16(proto.CmdType_CmdSyncMessageReq))
	pb.Register("im", &proto.SyncMessageResp{}, uint16(proto.CmdType_CmdSyncMessageResp))
	pb.Register("im", &proto.RecallMessageReq{}, uint16(proto.CmdType_CmdRecallMessageReq))
	pb.Register("im", &proto.RecallMessageResp{}, uint16(proto.CmdType_CmdRecallMessageResp))
	pb.Register("im", &proto.NotifyRecallMessage{}, uint16(proto.CmdType_CmdNotifyRecallMessage))
}
