package conv

import (
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
)

func onAddMember(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.AddMemberReq)
	log.Debugf("onAddMember %v", req)

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if _, inConv := c.Members[u.ID]; !inConv {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	var idsMap = make(map[string]int, len(req.GetAddIds()))
	var addIds = make([]string, 0, len(req.GetAddIds()))
	for _, id := range req.GetAddIds() {
		if _, exist := c.Members[id]; !exist {
			// load 数据库
			if u2 := user.GetUser(id); u2 != nil {
				idsMap[id] = 0
				addIds = append(addIds, id)
			}
		}
	}

	if len(addIds) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_OK})
		return
	}

	if err := setNxConvUser(c.ID, idsMap); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(addIds)

	conv := c.Pack()
	for id := range idsMap {
		if u2 := user.GetUser(id); u2 != nil {
			u2.SendToClient(0, &pb.NotifyInvited{
				Conv:   conv,
				InitBy: u.ID,
			})
		}
	}

	// 通知给群里其他人
	notifyJoined := &pb.NotifyMemberJoined{
		Conv:    conv,
		JoinIds: addIds,
		InitBy:  u.ID,
	}
	c.Broadcast(notifyJoined, addIds...)

}

func onRemoveMember(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.RemoveMemberReq)
	log.Debugf("onRemoveMember %v", req)

	if len(req.GetRemoveIds()) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if role, inConv := c.Members[u.ID]; !inConv {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotInConversation})
		return
	} else if role != 1 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotHasPermission})
		return
	}

	var rmIds = make([]string, 0, len(req.GetRemoveIds()))
	for _, id := range req.GetRemoveIds() {
		if _, exist := c.Members[id]; exist {
			rmIds = append(rmIds, id)
		}
	}

	if len(rmIds) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_OK})
		return
	}

	if err := delConvUser(c.ID, rmIds); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember(rmIds)

	conv := c.Pack()
	for _, id := range rmIds {
		if u2 := user.GetUser(id); u2 != nil {
			u2.SendToClient(0, &pb.NotifyKicked{
				Conv:     conv,
				KickedBy: u.ID,
			})
		}
	}

	// 通知给群里其他人
	notifyMemberLeft := &pb.NotifyMemberLeft{
		Conv:     conv,
		LeftIds:  rmIds,
		KickedBy: u.ID,
	}
	c.Broadcast(notifyMemberLeft)

}

func onJoin(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.JoinReq)
	log.Debugf("onJoin %v", req)

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	conv := c.Pack()
	if _, inConv := c.Members[u.ID]; inConv {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_OK, Conv: conv})
		return
	}

	if err := setNxConvUser(c.ID, map[string]int{u.ID: 0}); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember([]string{u.ID})

	// 通知给群里其他人
	notifyJoined := &pb.NotifyMemberJoined{
		Conv:    conv,
		JoinIds: []string{u.ID},
		InitBy:  u.ID,
	}
	c.Broadcast(notifyJoined, u.ID)

}

func onQuit(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.QuitReq)
	log.Debugf("onQuit %v", req)

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	conv := c.Pack()
	if _, inConv := c.Members[u.ID]; !inConv {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_OK, Conv: conv})
		return
	}

	if err := delConvUser(c.ID, []string{u.ID}); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember([]string{u.ID})

	// 通知给群里其他人
	notifyMemberLeft := &pb.NotifyMemberLeft{
		Conv:     conv,
		LeftIds:  []string{u.ID},
		KickedBy: u.ID,
	}
	c.Broadcast(notifyMemberLeft)

}

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdAddMemberReq), onAddMember)
	gate.RegisterHandler(uint16(pb.CmdType_CmdRemoveMemberReq), onRemoveMember)
	gate.RegisterHandler(uint16(pb.CmdType_CmdJoinReq), onJoin)
	gate.RegisterHandler(uint16(pb.CmdType_CmdQuitReq), onQuit)
}
