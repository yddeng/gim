package conv

import (
	"fmt"
	"github.com/yddeng/gim/config"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/internal/protocol/pb"
	"testing"
	"time"
)

func init() {
	config.LoadConfig("../../config/config.toml")
}

func TestDBConversation(t *testing.T) {
	conf := config.GetConfig().DBConfig
	db.Open(conf.SqlType, conf.Host, conf.Port, conf.Database, conf.User, conf.Password)

	conv := &Conversation{
		Type:     pb.ConversationType_Normal,
		ID:       0,
		Name:     "test",
		Creator:  "ydd",
		CreateAt: time.Now().Unix(),
		Members:  map[string]struct{}{"ydd": {}},
	}

	if err := setConversation(conv); err != nil {
		t.Error(err)
	}
	t.Log(conv.ID)

	conv2, err := loadConversation(2)
	if err != nil {
		t.Error(err)
	}
	t.Log(conv2)

	conv.Name = "test2"
	conv.LastMessageID = 1
	conv.LastMessageAt = time.Now().Unix()
	if err := updateConversation(conv); err != nil {
		t.Error(err)
	}

}

func TestDate(t *testing.T) {
	t.Log(time.Now().AddDate(0, -2, -20).Format("20060102"))
}

func TestMessage(t *testing.T) {
	conf := config.GetConfig().DBConfig
	db.Open(conf.SqlType, conf.Host, conf.Port, conf.Database, conf.User, conf.Password)

	msg := &pb.MessageInfo{
		Msg: &pb.Message{
			Text: "hello world",
		},
		UserID:   "ydd",
		CreateAt: time.Now().Unix(),
		MsgID:    1,
		Recalled: false,
	}

	for i := int64(1); i <= 20; i++ {
		msg.MsgID = i
		if err := insertMessage(1, msg); err != nil {
			t.Error(err)
		}
	}

	tableName := makeMessageTableName()
	limit := 10
	infos, err := loadMessageBatch(1, 15, limit, tableName)
	if err != nil {
		t.Error(err)
	}

	if len(infos) < limit {
		// 部分数据不在当前表中，应向前或向后查找
		t.Log(fmt.Sprintf("table %s not enough message", tableName))
	}
	for _, v := range infos {
		t.Log(v)
	}

}
