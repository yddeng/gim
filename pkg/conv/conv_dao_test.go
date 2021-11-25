package conv

import (
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
