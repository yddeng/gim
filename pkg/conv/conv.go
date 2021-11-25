package conv

import (
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
	convID  = int64(1)
	convMap = map[int64]*Conversation{}
)

func GetConversation(convID int64) *Conversation {
	return convMap[convID]
}

type Conversation struct {
	Type          pb.ConversationType // 对话类型
	ID            int64               // 全局唯一ID
	Name          string              // 会话名
	Creator       string              // 对话创建者
	CreateAt      int64               // 创建时间戳 秒
	Members       map[string]struct{} // 成员列表
	LastMessageAt int64               // 最后一条消息的时间
	LastMessageID int64               // 最后一条消息的ID
	LastMessage   *pb.MessageInfo     // 最后一条消息
	Message       []*pb.MessageInfo
}

func (this *Conversation) Pack() *pb.Conversation {
	c := &pb.Conversation{
		Type:          this.Type,
		ID:            this.ID,
		Name:          this.Name,
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

func (this *Conversation) AddMember(ids []string) {
	for _, id := range ids {
		this.Members[id] = struct{}{}
	}
}

func (this *Conversation) RemoveMember(ids []string) {
	//f := func(s string, m *[]string) {
	//	for i, v := range *m {
	//		if v == s {
	//			*m = append((*m)[:i], (*m)[i+1:]...)
	//			break
	//		}
	//	}
	//}

	for _, id := range ids {
		delete(this.Members, id)
	}
}

func onCreateConversation(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.CreateConversationReq)
	log.Debugf("onCreateConversation %v", req)

	nowUnix := time.Now().Unix()
	c := &Conversation{
		ID:       convID,
		Creator:  u.ID,
		CreateAt: nowUnix,
		Members:  make(map[string]struct{}, len(req.GetMembers())),
		Name:     req.GetName(),
	}
	convID++

	c.Members[u.ID] = struct{}{}
	for _, id := range req.GetMembers() {
		// load 数据库
		if u2 := user.GetUser(id); u2 != nil {
			if u2 != u {
				c.Members[id] = struct{}{}
			}
		}

	}

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
