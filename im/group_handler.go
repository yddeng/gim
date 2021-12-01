package im

import (
	"fmt"
	"github.com/yddeng/gim/im/pb"
	"github.com/yddeng/utils/log"
	"time"
)

func onCreateGroup(u *User, msg *Message) {
	req := msg.GetData().(*pb.CreateGroupReq)
	log.Debugf("user(%s) onCreateGroup %v", u.ID, req)

	nowUnix := time.Now().Unix()
	g := &Group{
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

	if err := WrapFunc(dbInsertGroup)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_Error})
			return
		}
		for _, m := range members {
			m.GroupID = g.ID
			m.ID = fmt.Sprintf("%d_%s", m.GroupID, m.UserID)
		}

		if err := WrapFunc(dbSetNxGroupMember)(func(err2 error) {
			if err != nil {
				log.Error(err)
				u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_Error})
				_ = taskPool.Submit(func() { _ = dbDeleteGroup(g.ID) })
				return
			}
			g.AddMember(members)
			addGroup(g)

			group := g.Pack()
			u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{
				Code:  pb.ErrCode_OK,
				Group: group,
			})

			notify := &pb.NotifyInvited{
				Group:  group,
				InitBy: u.ID,
			}
			g.Broadcast(notify, u.ID)
		}, members); err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_Busy})
		}
	}, g); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateGroupResp{Code: pb.ErrCode_Busy})
	}
}

func onGetGroupList(u *User, msg *Message) {
	//req := msg.GetData().(*pb.GetGroupListReq)
	//log.Debugf("user(%s) onGetGroupList %v", u.ID, req)
	if err := WrapFunc(dbGetUserGroups)(func(groups map[int64]*Member, err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.GetGroupListResp{Code: pb.ErrCode_Error})
			return
		}

		resp := &pb.GetGroupListResp{
			Groups: make([]*pb.Group, 0, len(groups)),
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
		u.SendToClient(msg.GetSeq(), &pb.GetGroupListResp{Code: pb.ErrCode_Busy})
	}
}

func onDissolveGroup(u *User, msg *Message) {
	req := msg.GetData().(*pb.DissolveGroupReq)
	log.Debugf("user(%s) onDissolveGroup %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.DissolveGroupResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if m, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.DissolveGroupResp{Code: pb.ErrCode_UserNotInGroup})
		return
	} else if m.Role != 1 {
		u.SendToClient(msg.GetSeq(), &pb.DissolveGroupResp{Code: pb.ErrCode_UserNotHasPermission})
		return
	}

	g.deleting = true
	if err := WrapFunc(dbDeleteGroup)(func(err error) {
		if err != nil {
			g.deleting = false
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.DissolveGroupResp{Code: pb.ErrCode_Error})
			return
		}
		notifyDissolve := &pb.NotifyDissolveGroup{
			GroupID: g.ID,
			InitBy:  u.ID,
		}
		members := make([]*Member, 0, len(g.Members))
		for uId, m := range g.Members {
			members = append(members, m)
			NotifyUser(uId, notifyDissolve)
		}

		removeGroup(g.ID)
		u.SendToClient(msg.GetSeq(), &pb.DissolveGroupResp{Code: pb.ErrCode_OK})
		_ = taskPool.Submit(func() { _ = dbDelGroupMember(members) })

	}, g.ID); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.DissolveGroupResp{Code: pb.ErrCode_Busy})
	}
}

func onAddMember(u *User, msg *Message) {
	req := msg.GetData().(*pb.AddMemberReq)
	log.Debugf("user(%s) onAddMember %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_UserNotInGroup})
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
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_OK, ExistIds: existIds})
		return
	}

	if err := WrapFunc(dbSetNxGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_Error})
			return
		}
		g.AddMember(members)
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_OK, ExistIds: existIds})

		group := g.Pack()
		notifyInvited := &pb.NotifyInvited{
			Group:  group,
			InitBy: u.ID,
		}

		for _, m := range members {
			NotifyUser(m.UserID, notifyInvited)
		}

		// 通知给群里其他人
		notifyJoined := &pb.NotifyMemberJoined{
			Group:   group,
			JoinIds: addIds,
			InitBy:  u.ID,
		}
		g.Broadcast(notifyJoined, addIds...)
	}, members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.AddMemberResp{Code: pb.ErrCode_Busy})
	}
}

func onRemoveMember(u *User, msg *Message) {
	req := msg.GetData().(*pb.RemoveMemberReq)
	log.Debugf("user(%s) onRemoveMember %v", u.ID, req)

	if len(req.GetRemoveIds()) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if m, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotInGroup})
		return
	} else if m.Role != 1 {
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_UserNotHasPermission})
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
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_OK})
		return
	}

	if err := WrapFunc(dbDelGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_Error})
			return
		}

		g.RemoveMember(members)

		group := g.Pack()
		notifyKicked := &pb.NotifyKicked{
			Group:    group,
			KickedBy: u.ID,
		}
		for _, m := range members {
			NotifyUser(m.UserID, notifyKicked)
		}

		// 通知给群里其他人
		notifyMemberLeft := &pb.NotifyMemberLeft{
			Group:    group,
			LeftIds:  rmIds,
			KickedBy: u.ID,
		}
		g.Broadcast(notifyMemberLeft)

	}, members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_Busy})
	}

}

func onJoin(u *User, msg *Message) {
	req := msg.GetData().(*pb.JoinReq)
	log.Debugf("user(%s) onJoin %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	group := g.Pack()
	if _, isMember := g.Members[u.ID]; isMember {
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_OK, Group: group})
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
			u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_Error})
			return
		}
		g.AddMember(member)

		// 通知给群里其他人
		notifyJoined := &pb.NotifyMemberJoined{
			Group:   group,
			JoinIds: []string{u.ID},
			InitBy:  u.ID,
		}
		g.Broadcast(notifyJoined, u.ID)

	}, member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.JoinResp{Code: pb.ErrCode_Busy})
	}

}

func onQuit(u *User, msg *Message) {
	req := msg.GetData().(*pb.QuitReq)
	log.Debugf("user(%s) onQuit %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	group := g.Pack()
	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_OK, Group: group})
		return
	}

	member := []*Member{g.Members[u.ID]}
	if err := WrapFunc(dbDelGroupMember)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_Error})
			return
		}
		g.RemoveMember(member)

		// 通知给群里其他人
		notifyMemberLeft := &pb.NotifyMemberLeft{
			Group:    group,
			LeftIds:  []string{u.ID},
			KickedBy: u.ID,
		}
		g.Broadcast(notifyMemberLeft)

	}, member); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.QuitResp{Code: pb.ErrCode_Busy})
	}
}

func onGetGroupMembers(u *User, msg *Message) {
	req := msg.GetData().(*pb.GetGroupMembersReq)
	log.Debugf("user(%s) onGetGroupMembers %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.GetGroupMembersResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.GetGroupMembersResp{Code: pb.ErrCode_UserNotInGroup})
		return
	}

	resp := &pb.GetGroupMembersResp{Members: make([]*pb.Member, 0, len(g.Members))}
	for _, m := range g.Members {
		resp.Members = append(resp.Members, m.Pack())
	}
	u.SendToClient(msg.GetSeq(), resp)
}

func onSendMessage(u *User, msg *Message) {
	req := msg.GetData().(*pb.SendMessageReq)
	log.Debugf("user(%s) onSendMessage %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	member, isMember := g.Members[u.ID]
	if !isMember {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_UserNotInGroup})
		return
	}

	g.LastMessageID += 1
	m := &pb.MessageInfo{
		Msg:      req.GetMsg(),
		UserID:   u.ID,
		CreateAt: time.Now().Unix(),
		MsgID:    g.LastMessageID,
	}

	if err := WrapFunc(messageDeliver.pushMessage)(func(err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_Error})
			return
		}

		g.LastMessageAt = m.GetCreateAt()
		member.UpdateAt = m.GetCreateAt()

		_ = taskPool.Submit(func() {
			_ = dbUpdateGroup(g)
			_ = dbSetNxGroupMember([]*Member{member})
		})

		group := g.Pack()
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_OK, Group: group})

		notifyMessage := &pb.NotifyMessage{
			Group:    group,
			MsgInfos: []*pb.MessageInfo{m},
		}
		g.Broadcast(notifyMessage)
	}, g.ID, m); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_Busy})
	}
}

func onSyncMessage(u *User, msg *Message) {
	req := msg.GetData().(*pb.SyncMessageReq)
	log.Debugf("user(%s) onSyncMessage %v", u.ID, req)

	g := GetGroup(req.GetGroupID())
	if g == nil {
		u.SendToClient(msg.GetSeq(), &pb.SyncMessageResp{Code: pb.ErrCode_GroupNotExist})
		return
	}

	if _, isMember := g.Members[u.ID]; !isMember {
		u.SendToClient(msg.GetSeq(), &pb.SyncMessageResp{Code: pb.ErrCode_UserNotInGroup})
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
		u.SendToClient(msg.GetSeq(), &pb.SyncMessageResp{Group: g.Pack()})
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

	if err := WrapFunc(messageDeliver.loadMessage)(func(infos []*pb.MessageInfo, err error) {
		if err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.SyncMessageResp{Code: pb.ErrCode_Error})
			return
		}

		resp := &pb.SyncMessageResp{Group: g.Pack(), MsgInfos: infos}
		u.SendToClient(msg.GetSeq(), resp)
	}, g.ID, ids); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.SyncMessageResp{Code: pb.ErrCode_Busy})
	}
}

func init() {
	registerHandler(uint16(pb.CmdType_CmdCreateGroupReq), onCreateGroup)
	registerHandler(uint16(pb.CmdType_CmdGetGroupListReq), onGetGroupList)
	registerHandler(uint16(pb.CmdType_CmdDissolveGroupReq), onDissolveGroup)

	registerHandler(uint16(pb.CmdType_CmdAddMemberReq), onAddMember)
	registerHandler(uint16(pb.CmdType_CmdRemoveMemberReq), onRemoveMember)
	registerHandler(uint16(pb.CmdType_CmdJoinReq), onJoin)
	registerHandler(uint16(pb.CmdType_CmdQuitReq), onQuit)
	registerHandler(uint16(pb.CmdType_CmdGetGroupMembersReq), onGetGroupMembers)

	registerHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
	registerHandler(uint16(pb.CmdType_CmdSyncMessageReq), onSyncMessage)
}
