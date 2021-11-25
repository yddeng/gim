package conv

import (
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
	"time"
)

var tableName string

func InitMessageTable() {
	tableName = makeMessageTableName()
	if exist := existMessageTable(tableName); !exist {
		createMessageTable(tableName)
	}
}

func onSendMessage(u *user.User, msg *codec.Message) {
	req := msg.GetData().(*pb.SendMessageReq)
	log.Debugf("onSendMessage %v", req)

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
	if err := insertMessage(c.ID, m, tableName); err != nil {
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
	gate.RegisterHandler(uint16(pb.CmdType_CmdSendMessageReq), onSendMessage)
}
