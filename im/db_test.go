package im

import (
	"github.com/yddeng/gim/im/protocol"
	"testing"
	"time"
)

func startService() {
	Service("../config.toml")
}

func TestUser(t *testing.T) {
	startService()
	u := &User{
		ID:       "ydd",
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
		Extra:    map[string]string{"name": "ydd", "age": "24"},
	}

	if err := dbSetNxUser(u); err != nil {
		t.Error(err)
	}

	u, err := dbLoadUser("ydd")
	if err != nil {
		t.Error(err)
	}
	t.Log(u)
}

func TestGroup(t *testing.T) {
	startService()
	conv := &Group{
		Type:     protocol.GroupType_Normal,
		ID:       0,
		Creator:  "ydd",
		CreateAt: time.Now().Unix(),
	}

	if err := dbInsertGroup(conv); err != nil {
		t.Error(err)
	}
	t.Log(conv.ID)

	conv2, err := dbLoadGroup(2)
	if err != nil {
		t.Error(err)
	}
	t.Log(conv2)

	conv.LastMessageID = 1
	conv.LastMessageAt = time.Now().Unix()
	if err := dbUpdateGroup(conv); err != nil {
		t.Error(err)
	}

}

func TestMember(t *testing.T) {
	startService()
	convID := int64(1)
	members := []*Member{{
		ID:       "1_ydd",
		GroupID:  convID,
		UserID:   "ydd",
		Nickname: "ydd",
		CreateAt: time.Now().Unix(),
		Mute:     1,
		Role:     0,
	}}

	if err := dbSetNxGroupMember(members); err != nil {
		t.Error(err)
	}

	convs, err := dbGetUserGroups("ydd")
	if err != nil {
		t.Error(err)
	}
	t.Log(convs)

	user, err := dbGetGroupMembers(convID)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	if err := dbDelGroupMember(members); err != nil {
		t.Error(err)
	}
}

func TestDate(t *testing.T) {
	t.Log(time.Now().AddDate(0, -2, -20).Format("20060102"))
}

func TestFriendDB(t *testing.T) {
	startService()

	f := &Friend{
		ID:       "1_2",
		UserID1:  "1",
		UserID2:  "2",
		CreateAt: time.Now().Unix(),
		Status:   FriendStatusAgree,
	}

	if err := dbSetNxFriend(f); err != nil {
		t.Error(err)
	}

	friends, err := dbLoadFriends("1")
	t.Log(friends, err)

	friends, err = dbLoadFriends("2")
	t.Log(friends, err)

	if err := dbDelFriend("1_2"); err != nil {
		t.Error(err)
	}
}
