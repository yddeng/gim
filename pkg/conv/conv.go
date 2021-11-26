package conv

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/gim/pkg/util"
	"github.com/yddeng/utils/log"
	"time"
)

var (
	convMap = map[int64]*Conversation{}
)

func GetConversation(convID int64) *Conversation {
	return convMap[convID]
}

type Conversation struct {
	Type          pb.ConversationType // 对话类型
	ID            int64               // 全局唯一ID
	Creator       string              // 对话创建者
	CreateAt      int64               // 创建时间戳 秒
	Extra         map[string]string   // 附加属性
	LastMessageAt int64               // 最后一条消息的时间
	LastMessageID int64               // 最后一条消息的ID
	LastMessage   *pb.MessageInfo     // 最后一条消息
	Message       []*pb.MessageInfo
	Members       map[string]*CMember
}

func (this *Conversation) Pack() *pb.Conversation {
	c := &pb.Conversation{
		Type:          this.Type,
		ID:            this.ID,
		LastMessageAt: this.LastMessageAt,
		LastMessageID: this.LastMessageID,
	}

	return c
}

func (this *Conversation) Broadcast(msg proto.Message, except ...string) {
	for id := range this.Members {
		if has := util.HasString(id, except); !has {
			u := user.GetUser(id)
			if u != nil {
				u.SendToClient(0, msg)
			}
		}
	}
}

func (this *Conversation) AddMember(members []*CMember) {
	for _, m := range members {
		this.Members[m.UserID] = m
	}
}

func (this *Conversation) RemoveMember(members []*CMember) {
	for _, m := range members {
		delete(this.Members, m.UserID)
	}
}

func onCreateConversation(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.CreateConversationReq)
	log.Debugf("user(%s) onCreateConversation %v", u.ID, req)

	nowUnix := time.Now().Unix()
	c := &Conversation{
		Type:     pb.ConversationType_Normal,
		Creator:  u.ID,
		Extra:    req.GetExtra(),
		CreateAt: nowUnix,
		Members:  make(map[string]*CMember, 16),
	}

	members := make([]*CMember, 0, len(req.GetMembers())+1)
	members = append(members, &CMember{UserID: u.ID, CreateAt: nowUnix, Role: 1})
	for _, id := range req.GetMembers() {
		if u2 := user.GetUser(id); u2 != nil {
			if u.ID != id {
				members = append(members, &CMember{UserID: id, CreateAt: nowUnix, Role: 0})
			}
		}
	}

	if len(members) < 2 {
		u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{
			Code: pb.ErrCode_RequestArgumentErr,
		})
		return
	}

	if err := insertConversation(c); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{
			Code: pb.ErrCode_Error,
		})
		return
	}
	for _, m := range members {
		m.ConvID = c.ID
		m.ID = fmt.Sprintf("%d_%s", m.ConvID, m.UserID)
	}

	if err := setNxConvUser(members); err != nil {
		log.Error(err)
		u.SendToClient(msg.GetSeq(), &pb.CreateConversationResp{
			Code: pb.ErrCode_Error,
		})
		return
	}

	c.AddMember(members)
	convMap[c.ID] = c

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

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdCreateConversationReq), onCreateConversation)
}
