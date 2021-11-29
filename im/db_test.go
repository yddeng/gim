package im

import (
	"fmt"
	"github.com/yddeng/gim/config"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/internal/protocol/pb"
	"testing"
	"time"
)

func init() {
	conf := config.LoadConfig("../../config/config.toml")
	dbConfig := conf.DBConfig
	db.Open(dbConfig.SqlType, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password)
}

func TestUser(t *testing.T) {
	u := &User{
		ID:       "ydd",
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
		Extra:    map[string]string{"name": "ydd", "age": "24"},
	}

	if err := setNxUser(u); err != nil {
		t.Error(err)
	}

	u, err := loadUser("ydd")
	if err != nil {
		t.Error(err)
	}
	t.Log(u)
}

func TestGroup(t *testing.T) {
	conv := &Group{
		Type:     pb.GroupType_Normal,
		ID:       0,
		Creator:  "ydd",
		CreateAt: time.Now().Unix(),
	}

	if err := insertGroup(conv); err != nil {
		t.Error(err)
	}
	t.Log(conv.ID)

	conv2, err := loadGroup(2)
	if err != nil {
		t.Error(err)
	}
	t.Log(conv2)

	conv.LastMessageID = 1
	conv.LastMessageAt = time.Now().Unix()
	if err := updateGroup(conv); err != nil {
		t.Error(err)
	}

}

func TestMember(t *testing.T) {
	convID := int64(1)
	members := []*Member{{
		ID:       "1_ydd",
		ConvID:   convID,
		UserID:   "ydd",
		Nickname: "ydd",
		CreateAt: time.Now().Unix(),
		Mute:     1,
		Role:     0,
	}}

	if err := setNxConvUser(members); err != nil {
		t.Error(err)
	}

	convs, err := getUserGroups("ydd")
	if err != nil {
		t.Error(err)
	}
	t.Log(convs)

	user, err := getGroupUsers(convID)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	if err := delConvUser(members); err != nil {
		t.Error(err)
	}
}

func TestDate(t *testing.T) {
	t.Log(time.Now().AddDate(0, -2, -20).Format("20060102"))
}

func TestMessage(t *testing.T) {
	msg := &pb.MessageInfo{
		Msg: &pb.Message{
			Text: "hello world",
		},
		UserID:   "ydd",
		CreateAt: time.Now().Unix(),
		MsgID:    1,
		Recalled: false,
	}
	tableName := makeMessageTableName()

	if exist := existMessageTable(tableName); !exist {
		t.Log("not exist")
		createMessageTable(tableName)
	}

	for i := int64(1); i <= 20; i++ {
		msg.MsgID = i
		if err := setNxMessage(1, msg, tableName); err != nil {
			t.Error(err)
		}
	}

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
