package protocol

import (
	"github.com/yddeng/gim/internal/codec/pb"
	proto "github.com/yddeng/gim/internal/protocol/pb"
)

func init() {
	pb.Register("im", &proto.UserLoginReq{}, uint16(proto.CmdType_CmdUserLoginReq))
	pb.Register("im", &proto.UserLoginResp{}, uint16(proto.CmdType_CmdUserLoginResp))

	pb.Register("im", &proto.CreateGroupReq{}, uint16(proto.CmdType_CmdCreateGroupReq))
	pb.Register("im", &proto.CreateGroupResp{}, uint16(proto.CmdType_CmdCreateGroupResp))
	pb.Register("im", &proto.GetGroupMembersReq{}, uint16(proto.CmdType_CmdGetGroupMembersReq))
	pb.Register("im", &proto.GetGroupMembersResp{}, uint16(proto.CmdType_CmdGetGroupMembersResp))
	pb.Register("im", &proto.GetGroupListReq{}, uint16(proto.CmdType_CmdGetGroupListReq))
	pb.Register("im", &proto.GetGroupListResp{}, uint16(proto.CmdType_CmdGetGroupListResp))
	pb.Register("im", &proto.DissolveGroupReq{}, uint16(proto.CmdType_CmdDissolveGroupReq))
	pb.Register("im", &proto.DissolveGroupResp{}, uint16(proto.CmdType_CmdDissolveGroupResp))
	pb.Register("im", &proto.NotifyDissolveGroup{}, uint16(proto.CmdType_CmdNotifyDissolveGroup))

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
