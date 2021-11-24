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

	if len(req.GetAddIds()) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if inConv := c.HasUser(u.ID); !inConv {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	conv := c.Pack()
	var addIds = make([]string, 0, len(req.GetAddIds()))
	for _, id := range req.GetAddIds() {
		if exist := c.HasUser(id); !exist {
			// load 数据库
			if u2 := user.GetUser(id); u2 != nil {
				addIds = append(addIds, id)
				u2.SendToClient(0, &pb.NotifyInvited{
					Conv:   conv,
					InitBy: u.ID,
				})
			}
		}
	}

	u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_OK})

	if len(addIds) > 0 {
		c.AddMember(addIds)

		// 通知给群里其他人
		notifyJoined := &pb.NotifyMemberJoined{
			Conv:    conv,
			JoinIds: addIds,
			InitBy:  u.ID,
		}
		c.Broadcast(notifyJoined, addIds...)
	}

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

	if inConv := c.HasUser(u.ID); !inConv {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	// 是不是群主
	if c.Creator != u.ID {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotHasPermission})
		return
	}

	conv := c.Pack()
	var rmIds = make([]string, 0, len(req.GetRemoveIds()))
	for _, id := range req.GetRemoveIds() {
		if exist := c.HasUser(id); exist {
			rmIds = append(rmIds, id)

			// load 数据库
			if u2 := user.GetUser(id); u2 != nil {
				u2.SendToClient(0, &pb.NotifyKicked{
					Conv:     conv,
					KickedBy: u.ID,
				})
			}
		}
	}

	u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_OK})

	if len(rmIds) > 0 {
		c.RemoveMember(rmIds)
		// 通知给群里其他人
		notifyMemberLeft := &pb.NotifyMemberLeft{
			Conv:     conv,
			LeftIds:  rmIds,
			KickedBy: u.ID,
		}
		c.Broadcast(notifyMemberLeft)
	}

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
	if inConv := c.HasUser(u.ID); inConv {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_OK, Conv: conv})
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
	if inConv := c.HasUser(u.ID); !inConv {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_OK, Conv: conv})
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
