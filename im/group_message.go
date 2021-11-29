package im

import (
	"fmt"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"time"
)

func onCreateGroup(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.CreateGroupReq)
	log.Debugf("user(%s) onCreateGroup %v", u.ID, req)

	nowUnix := time.Now().Unix()
	c := &Group{
		Type:     pb.GroupType_Normal,
		Creator:  u.ID,
		Extra:    req.GetExtra(),
		CreateAt: nowUnix,
		Members:  make(map[string]*Member, 16),
	}

	members := make([]*Member, 0, len(req.GetMembers())+1)
	members = append(members, &Member{UserID: u.ID, CreateAt: nowUnix, Role: 1})
	for _, id := range req.GetMembers() {
		if u2 := GetUser(id); u2 != nil {
			if u.ID != id {
				members = append(members, &Member{UserID: id, CreateAt: nowUnix, Role: 0})
			}
		}
	}

	if len(members) < 2 {
		u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	if err := insertGroup(c); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_Error})
		return
	}
	for _, m := range members {
		m.GroupID = c.ID
		m.ID = fmt.Sprintf("%d_%s", m.GroupID, m.UserID)
	}

	if err := setNxGroupMember(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(members)
	addGroup(c)

	group := c.Pack()
	u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{
		Code:  pb.ErrCode_OK,
		Group: group,
	})

	notify := &pb.NotifyInvited{
		Group:  group,
		InitBy: u.ID,
	}
	c.Broadcast(notify, u.ID)
}

func onAddMember(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.AddMemberReq)
	log.Debugf("user(%s) onAddMember %v", u.ID, req)

	c := GetGroup(req.GetGroupID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if _, isMember := c.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_UserNotInGroup})
		return
	}

	nowUnix := time.Now().Unix()
	members := make([]*Member, 0, len(req.GetAddIds()))
	addIds := make([]string, 0, len(req.GetAddIds()))
	for _, id := range req.GetAddIds() {
		if _, exist := c.Members[id]; !exist {
			// load 数据库
			if u2 := GetUser(id); u2 != nil {
				members = append(members, &Member{
					ID:       fmt.Sprintf("%d_%s", c.ID, id),
					GroupID:  c.ID,
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

	if err := setNxGroupMember(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(members)

	group := c.Pack()
	for _, m := range members {
		if u2 := GetUser(m.UserID); u2 != nil {
			u2.SendToClient(0, &pb.NotifyInvited{
				Group:  group,
				InitBy: u.ID,
			})
		}
	}

	// 通知给群里其他人
	notifyJoined := &pb.NotifyMemberJoined{
		Group:   group,
		JoinIds: addIds,
		InitBy:  u.ID,
	}
	c.Broadcast(notifyJoined, addIds...)

}

func onRemoveMember(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.RemoveMemberReq)
	log.Debugf("user(%s) onRemoveMember %v", u.ID, req)

	if len(req.GetRemoveIds()) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	c := GetGroup(req.GetGroupID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if m, isMember := c.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotInGroup})
		return
	} else if m.Role != 1 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotHasPermission})
		return
	}

	rmIds := make([]string, 0, len(req.GetRemoveIds()))
	members := make([]*Member, 0, len(req.GetRemoveIds()))
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

	if err := delGroupMember(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember(members)

	group := c.Pack()
	for _, m := range members {
		if u2 := GetUser(m.UserID); u2 != nil {
			u2.SendToClient(0, &pb.NotifyKicked{
				Group:    group,
				KickedBy: u.ID,
			})
		}
	}

	// 通知给群里其他人
	notifyMemberLeft := &pb.NotifyMemberLeft{
		Group:    group,
		LeftIds:  rmIds,
		KickedBy: u.ID,
	}
	c.Broadcast(notifyMemberLeft)

}

func onJoin(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.JoinReq)
	log.Debugf("user(%s) onJoin %v", u.ID, req)

	c := GetGroup(req.GetGroupID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	group := c.Pack()
	if _, isMember := c.Members[u.ID]; isMember {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_OK, Group: group})
		return
	}

	member := []*Member{{
		ID:       fmt.Sprintf("%d_%s", c.ID, u.ID),
		GroupID:  c.ID,
		UserID:   u.ID,
		CreateAt: time.Now().Unix(),
	}}

	if err := setNxGroupMember(member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(member)

	// 通知给群里其他人
	notifyJoined := &pb.NotifyMemberJoined{
		Group:   group,
		JoinIds: []string{u.ID},
		InitBy:  u.ID,
	}
	c.Broadcast(notifyJoined, u.ID)

}

func onQuit(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.QuitReq)
	log.Debugf("user(%s) onQuit %v", u.ID, req)

	c := GetGroup(req.GetGroupID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	group := c.Pack()
	if _, isMember := c.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_OK, Group: group})
		return
	}

	member := []*Member{c.Members[u.ID]}
	if err := delGroupMember(member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember(member)

	// 通知给群里其他人
	notifyMemberLeft := &pb.NotifyMemberLeft{
		Group:    group,
		LeftIds:  []string{u.ID},
		KickedBy: u.ID,
	}
	c.Broadcast(notifyMemberLeft)

}

func onSendMessage(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.SendMessageReq)
	log.Debugf("user(%s) onSendMessage %v", u.ID, req)

	c := GetGroup(req.GetGroupID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	member, isMember := c.Members[u.ID]
	if !isMember {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_UserNotInGroup})
		return
	}

	msgID := c.LastMessageID + 1
	m := &pb.MessageInfo{
		Msg:      req.GetMsg(),
		UserID:   u.ID,
		CreateAt: time.Now().Unix(),
		MsgID:    msgID,
	}

	if err := messageDeliver.pushMessage(c.ID, m); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_Error})
		return
	}

	c.LastMessage = m
	c.LastMessageID = m.GetMsgID()
	c.LastMessageAt = m.GetCreateAt()
	c.Message = append(c.Message, m)
	member.UpdateAt = m.GetCreateAt()

	updateGroup(c)
	setNxGroupMember([]*Member{member})

	group := c.Pack()
	u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_OK, Group: group})

	notifyMessage := &pb.NotifyMessage{
		Group:    group,
		MsgInfos: []*pb.MessageInfo{m},
	}
	c.Broadcast(notifyMessage)
}

func onGetGroupMembers(u *User, msg *codec.Message) {

}

func init() {
	registerGroupHandler(uint16(pb.CmdType_CmdCreateGroupReq), onCreateGroup)
	registerGroupHandler(uint16(pb.CmdType_CmdAddMemberReq), onAddMember)
	registerGroupHandler(uint16(pb.CmdType_CmdRemoveMemberReq), onRemoveMember)
	registerGroupHandler(uint16(pb.CmdType_CmdJoinReq), onJoin)
	registerGroupHandler(uint16(pb.CmdType_CmdQuitReq), onQuit)
	registerGroupHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
	registerGroupHandler(uint16(pb.CmdType_CmdGetGroupMembersReq), onGetGroupMembers)
}
