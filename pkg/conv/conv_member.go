package conv

import (
	"fmt"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
	"time"
)

func onAddMember(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.AddMemberReq)
	log.Debugf("user(%s) onAddMember %v", u.ID, req)

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if _, inConv := c.Members[u.ID]; !inConv {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	nowUnix := time.Now().Unix()
	members := make([]*CMember, 0, len(req.GetAddIds()))
	addIds := make([]string, 0, len(req.GetAddIds()))
	for _, id := range req.GetAddIds() {
		if _, exist := c.Members[id]; !exist {
			// load 数据库
			if u2 := user.GetUser(id); u2 != nil {
				members = append(members, &CMember{
					ID:       fmt.Sprintf("%d_%s", c.ID, id),
					ConvID:   c.ID,
					UserID:   id,
					CreateAt: nowUnix,
				})
				addIds = append(addIds, id)
			}
		}
	}

	if len(members) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_OK})
		return
	}

	if err := setNxConvUser(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(members)

	conv := c.Pack()
	for _, m := range members {
		if u2 := user.GetUser(m.UserID); u2 != nil {
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
	log.Debugf("user(%s) onRemoveMember %v", u.ID, req)

	if len(req.GetRemoveIds()) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if m, inConv := c.Members[u.ID]; !inConv {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotInConversation})
		return
	} else if m.Role != 1 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotHasPermission})
		return
	}

	rmIds := make([]string, 0, len(req.GetRemoveIds()))
	members := make([]*CMember, 0, len(req.GetRemoveIds()))
	for _, id := range req.GetRemoveIds() {
		if m, exist := c.Members[id]; exist {
			rmIds = append(rmIds, id)
			members = append(members, m)
		}
	}

	if len(members) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_OK})
		return
	}

	if err := delConvUser(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember(members)

	conv := c.Pack()
	for _, m := range members {
		if u2 := user.GetUser(m.UserID); u2 != nil {
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
	log.Debugf("user(%s) onJoin %v", u.ID, req)

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

	member := []*CMember{{
		ID:       fmt.Sprintf("%d_%s", c.ID, u.ID),
		ConvID:   c.ID,
		UserID:   u.ID,
		CreateAt: time.Now().Unix(),
	}}

	if err := setNxConvUser(member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(member)

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
	log.Debugf("user(%s) onQuit %v", u.ID, req)

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

	member := []*CMember{c.Members[u.ID]}
	if err := delConvUser(member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember(member)

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
