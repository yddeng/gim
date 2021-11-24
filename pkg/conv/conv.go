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
	convID  = uint64(1)
	convMap = map[uint64]*Conversation{}
)

func GetConversation(convID uint64) *Conversation {
	return convMap[convID]
}

type Conversation struct {
	Type          pb.ConversationType // 对话类型
	ID            uint64              // 全局唯一ID
	Creator       string              // 对话创建者
	CreateAt      int64               // 创建时间戳 秒
	Members       []string            // 成员列表
	Name          string              // 会话名
	LastMessageAt *pb.MessageInfo     // 最后一条消息
	Message       []*pb.MessageInfo
}

func (this *Conversation) Pack() *pb.Conversation {
	c := &pb.Conversation{
		Type: this.Type,
		ID:   this.ID,
		Name: this.Name,
	}

	if this.LastMessageAt != nil {
		c.LastMessageTimestamp = this.LastMessageAt.GetCreateAt()
		c.LastMessageID = this.LastMessageAt.GetMsgID()
	}

	return c
}

func (this *Conversation) Broadcast(msg proto.Message, except ...string) {
	for _, id := range this.Members {
		if has := util.HasString(id, except); !has {
			u := user.GetUser(id)
			if u != nil {
				u.SendToClient(0, msg)
			}
		}
	}
}

func (this *Conversation) HasUser(id string) bool {
	return util.HasString(id, this.Members)
}

func (this *Conversation) AddMember(ids []string) {
	this.Members = append(this.Members, ids...)
}

func (this *Conversation) RemoveMember(ids []string) {
	f := func(s string, m *[]string) {
		for i, v := range *m {
			if v == s {
				*m = append((*m)[:i], (*m)[i+1:]...)
				break
			}
		}
	}

	for _, id := range ids {
		f(id, &this.Members)
	}
	log.Debug(this.Members)
}

func onCreateConversation(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.CreateConversationReq)
	log.Debugf("onCreateConversation %v", req)

	nowUnix := time.Now().Unix()
	c := &Conversation{
		ID:       convID,
		Creator:  u.ID,
		CreateAt: nowUnix,
		Members:  make([]string, 0, len(req.GetMembers())),
		Name:     req.GetName(),
	}
	convID++

	c.Members = append(c.Members, u.ID)
	for _, id := range req.GetMembers() {
		// load 数据库
		if u2 := user.GetUser(id); u2 != nil {
			if u2 != u {
				c.Members = append(c.Members, id)
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
