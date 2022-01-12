package im

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/protocol"
	"github.com/yddeng/utils/log"
	"time"
)

func onCreateGroup(u *User, msg *Message) {
	req := msg.GetData().(*protocol.CreateGroupReq)
	log.Debugf("user(%s) onCreateGroup %v", u.ID, req)

	nowUnix := time.Now().Unix()
	g := &Group{
		Type:     protocol.GroupType_Normal,
		Creator:  u.ID,
		Extra:    req.GetExtras(),
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
		u.SendToClient(msg.GetSeq(), &protocol.CreateGroupResp{Code: protocol.ErrCode_RequestArgumentErr.Enum()})
		return
	}

	f := func(g *Group, members []*Member) error {
		if err := dbInsertGroup(g); err != nil {
			return err
		}
		for _, m := range members {
			m.GroupID = g.ID
			m.ID = fmt.Sprintf("%d_%s", m.GroupID, m.UserID)
		}
		if err := dbSetNxGroupMember(members); err != nil {
			_ = dbDeleteGroup(g.ID)
			return err
		}
		return nil
	}

	if err := WrapFunc(f)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.CreateGroupResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}
		g.AddMember(members)
		addGroup(g)

		group := g.Pack()
		u.SendToClient(msg.GetSeq(), &protocol.CreateGroupResp{
			Code:  protocol.ErrCode_OK.Enum(),
			Group: group,
		})

		notify := &protocol.NotifyInvited{
			Group:  group,
			InitBy: proto.String(u.ID),
		}
		g.Broadcast(notify, u.ID)
	}, g, members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.CreateGroupResp{Code: protocol.ErrCode_Busy.Enum()})
	}

}

func onGetGroupList(u *User, msg *Message) {
	//req := msg.GetData().(*protocol.GetGroupListReq)
	log.Debugf("user(%s) onGetGroupList", u.ID)
	if err := WrapFunc(dbGetUserGroups)(func(groups map[int64]*Member, err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.GetGroupListResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}

		resp := &protocol.GetGroupListResp{
			Groups: make([]*protocol.Group, 0, len(groups)),
		}

		for id := range groups {
			g := GetGroup(id)
			if g != nil {
				resp.Groups = append(resp.Groups, g.Pack())
			}
		}

		u.SendToClient(msg.GetSeq(), resp)
	}, u.ID); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.GetGroupListResp{Code: protocol.ErrCode_Busy.Enum()})
	}
}

func onDissolveGroup(u *User, msg *Message) {
	req := msg.GetData().(*protocol.DissolveGroupReq)
	log.Debugf("user(%s) onDissolveGroup %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.DissolveGroupResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	if m, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.DissolveGroupResp{Code: protocol.ErrCode_UserNotInGroup.Enum()})
		return
	} else if m.Role != 1 {
		u.SendToClient(msg.GetSeq(), &protocol.DissolveGroupResp{Code: protocol.ErrCode_UserNotHasPermission.Enum()})
		return
	}

	g.deleting = true
	if err := WrapFunc(dbDeleteGroup)(func(err error) {
		if err != nil {
			g.deleting = false
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.DissolveGroupResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}
		notifyDissolve := &protocol.NotifyDissolveGroup{
			GroupID: proto.Int64(g.ID),
			InitBy:  proto.String(u.ID),
		}
		members := make([]*Member, 0, len(g.Members))
		for uId, m := range g.Members {
			members = append(members, m)
			NotifyUser(uId, notifyDissolve)
		}

		removeGroup(g.ID)
		u.SendToClient(msg.GetSeq(), &protocol.DissolveGroupResp{Code: protocol.ErrCode_OK.Enum()})
		_ = taskPool.Submit(func() { _ = dbDelGroupMember(members) })

	}, g.ID); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.DissolveGroupResp{Code: protocol.ErrCode_Busy.Enum()})
	}
}

func onAddMember(u *User, msg *Message) {
	req := msg.GetData().(*protocol.AddMemberReq)
	log.Debugf("user(%s) onAddMember %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.AddMemberResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.AddMemberResp{Code: protocol.ErrCode_UserNotInGroup.Enum()})
		return
	}

	nowUnix := time.Now().Unix()
	members := make([]*Member, 0, len(req.GetAddIds()))
	addIds := make([]string, 0, len(req.GetAddIds()))
	existIds := make([]string, 0, len(req.GetAddIds()))
	for _, id := range req.GetAddIds() {
		if _, exist := g.Members[id]; !exist {
			// load 数据库
			if u2 := GetUser(id); u2 != nil {
				members = append(members, &Member{
					ID:       fmt.Sprintf("%d_%s", g.ID, id),
					GroupID:  g.ID,
					UserID:   id,
					CreateAt: nowUnix,
					UpdateAt: nowUnix,
				})
				addIds = append(addIds, id)
			}
		} else {
			existIds = append(existIds, id)
		}
	}

	if len(members) == 0 {
		u.SendToClient(msg.GetSeq(), &protocol.AddMemberResp{Code: protocol.ErrCode_OK.Enum(), ExistIds: existIds})
		return
	}

	if err := WrapFunc(dbSetNxGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.AddMemberResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}
		g.AddMember(members)
		u.SendToClient(msg.GetSeq(), &protocol.AddMemberResp{Code: protocol.ErrCode_OK.Enum(), ExistIds: existIds})

		group := g.Pack()
		notifyInvited := &protocol.NotifyInvited{
			Group:  group,
			InitBy: proto.String(u.ID),
		}

		for _, m := range members {
			NotifyUser(m.UserID, notifyInvited)
		}

		// 通知给群里其他人
		notifyJoined := &protocol.NotifyMemberJoined{
			Group:   group,
			JoinIds: addIds,
			InitBy:  proto.String(u.ID),
		}
		g.Broadcast(notifyJoined, addIds...)
	}, members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.AddMemberResp{Code: protocol.ErrCode_Busy.Enum()})
	}
}

func onRemoveMember(u *User, msg *Message) {
	req := msg.GetData().(*protocol.RemoveMemberReq)
	log.Debugf("user(%s) onRemoveMember %v", u.ID, req)

	if len(req.GetRemoveIds()) == 0 {
		u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_RequestArgumentErr.Enum()})
		return
	}

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	if m, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_UserNotInGroup.Enum()})
		return
	} else if m.Role != 1 {
		u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_UserNotHasPermission.Enum()})
		return
	}

	rmIds := make([]string, 0, len(req.GetRemoveIds()))
	members := make([]*Member, 0, len(req.GetRemoveIds()))
	for _, id := range req.GetRemoveIds() {
		if m, exist := g.Members[id]; exist {
			rmIds = append(rmIds, id)
			members = append(members, m)
		}
	}

	if len(members) == 0 {
		u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_OK.Enum()})
		return
	}

	if err := WrapFunc(dbDelGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}

		g.RemoveMember(members)

		group := g.Pack()
		notifyKicked := &protocol.NotifyKicked{
			Group:    group,
			KickedBy: proto.String(u.ID),
		}
		for _, m := range members {
			NotifyUser(m.UserID, notifyKicked)
		}

		// 通知给群里其他人
		notifyMemberLeft := &protocol.NotifyMemberLeft{
			Group:    group,
			LeftIds:  rmIds,
			KickedBy: proto.String(u.ID),
		}
		g.Broadcast(notifyMemberLeft)

	}, members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.RemoveMemberResp{Code: protocol.ErrCode_Busy.Enum()})
	}

}

func onJoin(u *User, msg *Message) {
	req := msg.GetData().(*protocol.JoinReq)
	log.Debugf("user(%s) onJoin %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.JoinResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	group := g.Pack()
	if _, isMember := g.Members[u.ID]; isMember {
		u.SendToClient(msg.GetSeq(), &protocol.JoinResp{Code: protocol.ErrCode_OK.Enum(), Group: group})
		return
	}

	nowUnix := time.Now().Unix()
	member := []*Member{{
		ID:       fmt.Sprintf("%d_%s", g.ID, u.ID),
		GroupID:  g.ID,
		UserID:   u.ID,
		CreateAt: nowUnix,
		UpdateAt: nowUnix,
	}}

	if err := WrapFunc(dbSetNxGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.JoinResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}
		g.AddMember(member)

		// 通知给群里其他人
		notifyJoined := &protocol.NotifyMemberJoined{
			Group:   group,
			JoinIds: []string{u.ID},
			InitBy:  proto.String(u.ID),
		}
		g.Broadcast(notifyJoined, u.ID)

	}, member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.JoinResp{Code: protocol.ErrCode_Busy.Enum()})
	}

}

func onQuit(u *User, msg *Message) {
	req := msg.GetData().(*protocol.QuitReq)
	log.Debugf("user(%s) onQuit %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.QuitResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	group := g.Pack()
	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.QuitResp{Code: protocol.ErrCode_OK.Enum(), Group: group})
		return
	}

	member := []*Member{g.Members[u.ID]}
	if err := WrapFunc(dbDelGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.QuitResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}
		g.RemoveMember(member)

		// 通知给群里其他人
		notifyMemberLeft := &protocol.NotifyMemberLeft{
			Group:    group,
			LeftIds:  []string{u.ID},
			KickedBy: proto.String(u.ID),
		}
		g.Broadcast(notifyMemberLeft)

	}, member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.QuitResp{Code: protocol.ErrCode_Busy.Enum()})
	}
}

func onGetGroupMembers(u *User, msg *Message) {
	req := msg.GetData().(*protocol.GetGroupMembersReq)
	log.Debugf("user(%s) onGetGroupMembers %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.GetGroupMembersResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.GetGroupMembersResp{Code: protocol.ErrCode_UserNotInGroup.Enum()})
		return
	}

	resp := &protocol.GetGroupMembersResp{Members: make([]*protocol.Member, 0, len(g.Members))}
	for _, m := range g.Members {
		resp.Members = append(resp.Members, m.Pack())
	}
	u.SendToClient(msg.GetSeq(), resp)
}

func onSendMessage(u *User, msg *Message) {
	req := msg.GetData().(*protocol.SendMessageReq)
	log.Debugf("user(%s) onSendMessage %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.SendMessageResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	member, isMember := g.Members[u.ID]
	if !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.SendMessageResp{Code: protocol.ErrCode_UserNotInGroup.Enum()})
		return
	}

	g.LastMessageID += 1
	m := &protocol.MessageInfo{
		Msg:      req.GetMsg(),
		UserID:   proto.String(u.ID),
		CreateAt: proto.Int64(time.Now().Unix()),
		MsgID:    proto.Int64(g.LastMessageID),
	}

	if err := WrapFunc(messageDeliver.pushMessage)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.SendMessageResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}

		g.LastMessageAt = m.GetCreateAt()
		member.UpdateAt = m.GetCreateAt()

		_ = taskPool.Submit(func() {
			_ = dbUpdateGroup(g)
			_ = dbSetNxGroupMember([]*Member{member})
		})

		group := g.Pack()
		u.SendToClient(msg.GetSeq(), &protocol.SendMessageResp{Code: protocol.ErrCode_OK.Enum(), Group: group})

		notifyMessage := &protocol.NotifyMessage{
			Group:    group,
			MsgInfos: []*protocol.MessageInfo{m},
		}
		g.Broadcast(notifyMessage)
	}, g.ID, m); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.SendMessageResp{Code: protocol.ErrCode_Busy.Enum()})
	}
}

func onSyncMessage(u *User, msg *Message) {
	req := msg.GetData().(*protocol.SyncMessageReq)
	log.Debugf("user(%s) onSyncMessage %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &protocol.SyncMessageResp{Code: protocol.ErrCode_GroupNotExist.Enum()})
		return
	}

	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &protocol.SyncMessageResp{Code: protocol.ErrCode_UserNotInGroup.Enum()})
		return
	}
	limit := req.GetLimit()
	if limit == 0 || limit > 50 {
		limit = 50
	}

	ids := make([]int64, 0, limit)
	for i := int64(0); i < int64(limit); i++ {
		id := req.GetStartID()
		if req.GetOldToNew() {
			id = id + i
		} else {
			id = id - i
		}
		if id >= 1 && id <= g.LastMessageID {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		u.SendToClient(msg.GetSeq(), &protocol.SyncMessageResp{Group: g.Pack()})
		return
	}

	// 反转，使之从小到大排序
	if !req.GetOldToNew() {
		newIds := make([]int64, 0, len(ids))
		for i := len(ids) - 1; i >= 0; i-- {
			newIds = append(newIds, ids[i])
		}
		ids = newIds
	}

	if err := WrapFunc(messageDeliver.loadMessage)(func(infos []*protocol.MessageInfo, err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.SyncMessageResp{Code: protocol.ErrCode_Error.Enum()})
			return
		}

		resp := &protocol.SyncMessageResp{Group: g.Pack(), MsgInfos: infos}
		u.SendToClient(msg.GetSeq(), resp)
	}, g.ID, ids); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &protocol.SyncMessageResp{Code: protocol.ErrCode_Busy.Enum()})
	}
}

func init() {
	registerHandler(uint16(protocol.CmdType_CmdCreateGroupReq), onCreateGroup)
	registerHandler(uint16(protocol.CmdType_CmdGetGroupListReq), onGetGroupList)
	registerHandler(uint16(protocol.CmdType_CmdDissolveGroupReq), onDissolveGroup)

	registerHandler(uint16(protocol.CmdType_CmdAddMemberReq), onAddMember)
	registerHandler(uint16(protocol.CmdType_CmdRemoveMemberReq), onRemoveMember)
	registerHandler(uint16(protocol.CmdType_CmdJoinReq), onJoin)
	registerHandler(uint16(protocol.CmdType_CmdQuitReq), onQuit)
	registerHandler(uint16(protocol.CmdType_CmdGetGroupMembersReq), onGetGroupMembers)

	registerHandler(uint16(protocol.CmdType_CmdSendMessageReq), onSendMessage)
	registerHandler(uint16(protocol.CmdType_CmdSyncMessageReq), onSyncMessage)

	registerHandler(uint16(protocol.CmdType_CmdGetUserInfoReq), onGetUserInfo)
}
