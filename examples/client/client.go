package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/im"
	"github.com/yddeng/gim/im/protocol"
	"github.com/yddeng/utils/log"
	"os"
	"strings"
	"time"
)

var handler = map[uint16]func(dnet.Session, *im.Message){}

func registerHandler(cmdType protocol.CmdType, h func(dnet.Session, *im.Message)) {
	cmd := uint16(cmdType)
	if _, ok := handler[cmd]; ok {
		panic(fmt.Sprintf("cmd %d is alreadly register. ", cmd))
	}
	handler[cmd] = h
}

func dispatchMessage(sess dnet.Session, msg *im.Message) {
	if h, ok := handler[msg.GetCmd()]; ok {
		h(sess, msg)
	}
}

func main() {
	//address := "10.128.2.123:43210"
	address := "81.69.172.73:41210"
	c, err := dnet.DialTCP(address, 0)
	if err != nil {
		panic(err)
	}

	sess := dnet.NewTCPSession(c,
		dnet.WithCodec(im.Codec{}),
		dnet.WithErrorCallback(func(session dnet.Session, err error) {
			fmt.Println("onError", err)
		}),
		dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
			msg := data.(*im.Message)
			dispatchMessage(session, msg)
		}),
		dnet.WithCloseCallback(func(session dnet.Session, reason error) {
			fmt.Println("onClose", reason)
		}),
	)

	go func() {
		for {
			time.Sleep(time.Second * 5)
			if err := sess.Send(im.NewMessage(0, &protocol.Heartbeat{})); err != nil {
				panic(err)
			}
		}
	}()

	userID := os.Args[1]
	if err := sess.Send(im.NewMessage(1, &protocol.UserLoginReq{
		ID:     proto.String(userID),
		Extras: []*protocol.Extra{{Key: proto.String(userID), Value: proto.String(userID)}},
	})); err != nil {
		panic(err)
	}

	registerHandler(protocol.CmdType_CmdUserLoginResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("UserLoginResp %v", msg.GetData().(*protocol.UserLoginResp))
	})
	registerHandler(protocol.CmdType_CmdCreateGroupResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("CreateGroupResp %v", msg.GetData().(*protocol.CreateGroupResp))
	})
	registerHandler(protocol.CmdType_CmdAddMemberResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("AddMemberResp %v", msg.GetData().(*protocol.AddMemberResp))
	})
	registerHandler(protocol.CmdType_CmdRemoveMemberResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("RemoveMemberResp %v", msg.GetData().(*protocol.RemoveMemberResp))
	})
	registerHandler(protocol.CmdType_CmdJoinResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("JoinResp %v", msg.GetData().(*protocol.JoinResp))
	})
	registerHandler(protocol.CmdType_CmdQuitResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("QuitResp %v", msg.GetData().(*protocol.QuitResp))
	})
	registerHandler(protocol.CmdType_CmdSendMessageResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("SendMessageResp %v", msg.GetData().(*protocol.SendMessageResp))
	})
	registerHandler(protocol.CmdType_CmdNotifyInvited, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyInvited %v", msg.GetData().(*protocol.NotifyInvited))
	})
	registerHandler(protocol.CmdType_CmdNotifyKicked, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyKicked %v", msg.GetData().(*protocol.NotifyKicked))
	})
	registerHandler(protocol.CmdType_CmdNotifyMemberJoined, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyMemberJoined %v", msg.GetData().(*protocol.NotifyMemberJoined))
	})
	registerHandler(protocol.CmdType_CmdNotifyMemberLeft, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyMemberLeft %v", msg.GetData().(*protocol.NotifyMemberLeft))
	})
	registerHandler(protocol.CmdType_CmdNotifyMessage, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyMessage %v", msg.GetData().(*protocol.NotifyMessage))
	})
	registerHandler(protocol.CmdType_CmdGetGroupMembersResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("GetGroupMembersResp %v", msg.GetData().(*protocol.GetGroupMembersResp))
	})
	registerHandler(protocol.CmdType_CmdSyncMessageResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("SyncMessageResp %v", msg.GetData().(*protocol.SyncMessageResp))
	})
	registerHandler(protocol.CmdType_CmdGetGroupListResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("GetGroupListResp %v", msg.GetData().(*protocol.GetGroupListResp))
	})
	registerHandler(protocol.CmdType_CmdAddFriendResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("AddFriendResp %v", msg.GetData().(*protocol.AddFriendResp))
	})
	registerHandler(protocol.CmdType_CmdNotifyAgreeFriend, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyAgreeFriend %v", msg.GetData().(*protocol.NotifyAgreeFriend))
	})
	registerHandler(protocol.CmdType_CmdAgreeFriendResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("AgreeFriendResp %v", msg.GetData().(*protocol.AgreeFriendResp))
	})
	registerHandler(protocol.CmdType_CmdNotifyAddFriend, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyAddFriend %v", msg.GetData().(*protocol.NotifyAddFriend))
	})
	registerHandler(protocol.CmdType_CmdDeleteFriendResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("DeleteFriendResp %v", msg.GetData().(*protocol.DeleteFriendResp))
	})
	registerHandler(protocol.CmdType_CmdNotifyDeleteFriend, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyDeleteFriend %v", msg.GetData().(*protocol.NotifyDeleteFriend))
	})
	registerHandler(protocol.CmdType_CmdGetFriendsResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("GetFriendsResp %v", msg.GetData().(*protocol.GetFriendsResp))
	})
	registerHandler(protocol.CmdType_CmdGetUserInfoResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("GetUserInfoResp %v", msg.GetData().(*protocol.GetUserInfoResp))
	})
	registerHandler(protocol.CmdType_CmdNotifyUserOnline, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyUserOnline %v", msg.GetData().(*protocol.NotifyUserOnline))
	})

	go func() {
		for {
			cmd(sess)
		}
	}()
	select {}
}

func cmd(sess dnet.Session) {
	fmt.Println("1:createGroup 2:addMember 3:removeMember 4:join 5:quit 6:send 7:getMembers 8:syncMessage 9:getGroupList")
	fmt.Println("11:addFriend 12:agreeFriend 13:getFriends 14:deleteFriend 15:getUserInfo")
	fmt.Printf("==>")
	var k int
	fmt.Scan(&k)
	switch k {
	case 1:
		fmt.Printf("CreateGroup users:")
		var users string
		fmt.Scan(&users)
		us := strings.Split(users, "&")
		_ = sess.Send(im.NewMessage(0, &protocol.CreateGroupReq{
			Members: us,
		}))
	case 2:
		fmt.Printf("AddMember id,users :")
		var id int
		var users string
		fmt.Scan(&id, &users)
		us := strings.Split(users, "&")
		_ = sess.Send(im.NewMessage(0, &protocol.AddMemberReq{
			GroupID: proto.Int64(int64(id)),
			AddIds:  us,
		}))
	case 3:
		fmt.Printf("RemoveMember id,users :")
		var id int
		var users string
		fmt.Scan(&id, &users)
		us := strings.Split(users, "&")
		_ = sess.Send(im.NewMessage(0, &protocol.RemoveMemberReq{
			GroupID:   proto.Int64(int64(id)),
			RemoveIds: us,
		}))
	case 4:
		fmt.Printf("Join id :")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(im.NewMessage(0, &protocol.JoinReq{
			GroupID: proto.Int64(int64(id)),
		}))
	case 5:
		fmt.Printf("Quit id :")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(im.NewMessage(0, &protocol.QuitReq{
			GroupID: proto.Int64(int64(id)),
		}))
	case 6:
		fmt.Printf("Send id,msg:")
		var id int
		var msg string
		fmt.Scan(&id, &msg)
		_ = sess.Send(im.NewMessage(0, &protocol.SendMessageReq{
			GroupID: proto.Int64(int64(id)),
			Msg:     &protocol.Message{Text: proto.String(msg)},
		}))
	case 7:
		fmt.Printf("GetMember id:")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(im.NewMessage(0, &protocol.GetGroupMembersReq{
			GroupID: proto.Int64(int64(id)),
		}))
	case 8:
		fmt.Printf("Send id,startID,limit:")
		var id, startID, limit int
		fmt.Scan(&id, &startID, &limit)
		_ = sess.Send(im.NewMessage(0, &protocol.SyncMessageReq{
			GroupID:  proto.Int64(int64(id)),
			StartID:  proto.Int64(int64(startID)),
			Limit:    proto.Int32(int32(limit)),
			OldToNew: proto.Bool(false),
		}))
	case 9:
		fmt.Printf("Send :")
		_ = sess.Send(im.NewMessage(0, &protocol.GetGroupListReq{}))
	case 11:
		fmt.Printf("Send user:")
		var userID string
		fmt.Scan(&userID)
		_ = sess.Send(im.NewMessage(0, &protocol.AddFriendReq{
			UserID: proto.String(userID),
		}))
	case 12:
		fmt.Printf("Send userID,agree:")
		var userID string
		fmt.Scan(&userID)
		var ok int
		fmt.Scan(&ok)
		_ = sess.Send(im.NewMessage(0, &protocol.AgreeFriendReq{
			UserID: proto.String(userID),
			Agree:  proto.Bool(ok != 0),
		}))
	case 13:
		fmt.Printf("Send :")
		_ = sess.Send(im.NewMessage(0, &protocol.GetFriendsReq{}))
	case 14:
		fmt.Printf("Send userID:")
		var userID string
		fmt.Scan(&userID)
		_ = sess.Send(im.NewMessage(0, &protocol.DeleteFriendReq{
			UserID: proto.String(userID),
		}))
	case 15:
		fmt.Printf("Send userID:")
		var userID string
		fmt.Scan(&userID)
		_ = sess.Send(im.NewMessage(0, &protocol.GetUserInfoReq{
			UserIDs: []string{userID},
		}))
	}
}
