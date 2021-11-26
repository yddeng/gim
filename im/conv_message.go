package im

import (
	"fmt"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"time"
)

func onCreateConversation(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.CreateConversationReq)
	log.Debugf("user(%s) onCreateConversation %v", u.ID, req)

	nowUnix := time.Now().Unix()
	c := &Conversation{
		Type:     pb.ConversationType_Normal,
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
		u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{Code: pb.ErrCode_RequestArgumentErr})
		return
	}

	if err := insertConversation(c); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{Code: pb.ErrCode_Error})
		return
	}
	for _, m := range members {
		m.ConvID = c.ID
		m.ID = fmt.Sprintf("%d_%s", m.ConvID, m.UserID)
	}

	if err := setNxConvUser(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{Code: pb.ErrCode_Error})
		return
	}

	c.AddMember(members)
	addConversation(c)

	conv := c.Pack()
	u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{
		Code: pb.ErrCode_OK,
		Conv: conv,
	})

	notify := &pb.NotifyInvited{
		Conv:   conv,
		InitBy: u.ID,
	}
	c.Broadcast(notify, u.ID)
}

func onAddMember(u *User, msg *codec.Message) {
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
	members := make([]*Member, 0, len(req.GetAddIds()))
	addIds := make([]string, 0, len(req.GetAddIds()))
	for _, id := range req.GetAddIds() {
		if _, exist := c.Members[id]; !exist {
			// load 数据库
			if u2 := GetUser(id); u2 != nil {
				members = append(members, &Member{
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
		if u2 := GetUser(m.UserID); u2 != nil {
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

func onRemoveMember(u *User, msg *codec.Message) {
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

	if err := delConvUser(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.RemoveMemberResp{Code: pb.ErrCode_Error})
		return
	}

	c.RemoveMember(members)

	conv := c.Pack()
	for _, m := range members {
		if u2 := GetUser(m.UserID); u2 != nil {
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

func onJoin(u *User, msg *codec.Message) {
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

	member := []*Member{{
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

func onQuit(u *User, msg *codec.Message) {
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

	member := []*Member{c.Members[u.ID]}
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

func onSendMessage(u *User, msg *codec.Message) {
	req := msg.GetData().(*pb.SendMessageReq)
	log.Debugf("user(%s) onSendMessage %v", u.ID, req)

	c := GetConversation(req.GetConvID())
	if c == nil {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_ConversationNotExist})
		return
	}

	if _, inConv := c.Members[u.ID]; !inConv {
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_UserNotInConversation})
		return
	}

	msgID := c.LastMessageID + 1
	m := &pb.MessageInfo{
		Msg:      req.GetMsg(),
		UserID:   u.ID,
		CreateAt: time.Now().Unix(),
		MsgID:    msgID,
	}

	nowTableName := makeMessageTableName()
	if tableName != nowTableName {
		tableName = nowTableName
		if err := createMessageTable(tableName); err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_Error})
			return
		}
	}
	if err := setNxMessage(c.ID, m, tableName); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_Error})
		return
	}

	c.LastMessage = m
	c.LastMessageID = m.GetMsgID()
	c.LastMessageAt = m.GetCreateAt()
	c.Message = append(c.Message, m)

	updateConversation(c)

	conv := c.Pack()
	u.SendToClient(msg.GetSeq(), &pb.SendMessageResp{Code: pb.ErrCode_OK, Conv: conv})

	notifyMessage := &pb.NotifyMessage{
		Conv:     conv,
		MsgInfos: []*pb.MessageInfo{m},
	}
	c.Broadcast(notifyMessage)
}

func init() {
	registerHandler(uint16(pb.CmdType_CmdCreateConversationReq), onCreateConversation)
	registerHandler(uint16(pb.CmdType_CmdAddMemberReq), onAddMember)
	registerHandler(uint16(pb.CmdType_CmdRemoveMemberReq), onRemoveMember)
	registerHandler(uint16(pb.CmdType_CmdJoinReq), onJoin)
	registerHandler(uint16(pb.CmdType_CmdQuitReq), onQuit)
	registerHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
}
