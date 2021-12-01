package main

import (
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/im"
	"github.com/yddeng/gim/im/pb"
	"github.com/yddeng/utils/log"
	"os"
	"strings"
	"time"
)

var handler = map[uint16]func(dnet.Session, *im.Message){}

func registerHandler(cmdType pb.CmdType, h func(dnet.Session, *im.Message)) {
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
	address := "127.0.0.1:43210"
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
			_ = sess.Send(im.NewMessage(0, &pb.Heartbeat{}))
		}
	}()

	userID := os.Args[1]
	sess.Send(im.NewMessage(1, &pb.UserLoginReq{
		ID: userID,
	}))

	registerHandler(pb.CmdType_CmdUserLoginResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("UserLoginResp %v", msg.GetData().(*pb.UserLoginResp))
	})
	registerHandler(pb.CmdType_CmdCreateGroupResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("CreateGroupResp %v", msg.GetData().(*pb.CreateGroupResp))
	})
	registerHandler(pb.CmdType_CmdAddMemberResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("AddMemberResp %v", msg.GetData().(*pb.AddMemberResp))
	})
	registerHandler(pb.CmdType_CmdRemoveMemberResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("RemoveMemberResp %v", msg.GetData().(*pb.RemoveMemberResp))
	})
	registerHandler(pb.CmdType_CmdJoinResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("JoinResp %v", msg.GetData().(*pb.JoinResp))
	})
	registerHandler(pb.CmdType_CmdQuitResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("QuitResp %v", msg.GetData().(*pb.QuitResp))
	})
	registerHandler(pb.CmdType_CmdSendMessageResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("SendMessageResp %v", msg.GetData().(*pb.SendMessageResp))
	})
	registerHandler(pb.CmdType_CmdNotifyInvited, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyInvited %v", msg.GetData().(*pb.NotifyInvited))
	})
	registerHandler(pb.CmdType_CmdNotifyKicked, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyKicked %v", msg.GetData().(*pb.NotifyKicked))
	})
	registerHandler(pb.CmdType_CmdNotifyMemberJoined, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyMemberJoined %v", msg.GetData().(*pb.NotifyMemberJoined))
	})
	registerHandler(pb.CmdType_CmdNotifyMemberLeft, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyMemberLeft %v", msg.GetData().(*pb.NotifyMemberLeft))
	})
	registerHandler(pb.CmdType_CmdNotifyMessage, func(session dnet.Session, msg *im.Message) {
		log.Debugf("NotifyMessage %v", msg.GetData().(*pb.NotifyMessage))
	})
	registerHandler(pb.CmdType_CmdGetGroupMembersResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("GetGroupMembersResp %v", msg.GetData().(*pb.GetGroupMembersResp))
	})
	registerHandler(pb.CmdType_CmdSyncMessageResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("SyncMessageResp %v", msg.GetData().(*pb.SyncMessageResp))
	})
	registerHandler(pb.CmdType_CmdGetGroupListResp, func(session dnet.Session, msg *im.Message) {
		log.Debugf("GetGroupListResp %v", msg.GetData().(*pb.GetGroupListResp))
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
	fmt.Printf("==>")
	var k int
	fmt.Scan(&k)
	switch k {
	case 1:
		fmt.Printf("CreateGroup users:")
		var users string
		fmt.Scan(&users)
		us := strings.Split(users, "&")
		_ = sess.Send(im.NewMessage(0, &pb.CreateGroupReq{
			Members: us,
		}))
	case 2:
		fmt.Printf("AddMember id,users :")
		var id int
		var users string
		fmt.Scan(&id, &users)
		us := strings.Split(users, "&")
		_ = sess.Send(im.NewMessage(0, &pb.AddMemberReq{
			GroupID: int64(id),
			AddIds:  us,
		}))
	case 3:
		fmt.Printf("RemoveMember id,users :")
		var id int
		var users string
		fmt.Scan(&id, &users)
		us := strings.Split(users, "&")
		_ = sess.Send(im.NewMessage(0, &pb.RemoveMemberReq{
			GroupID:   int64(id),
			RemoveIds: us,
		}))
	case 4:
		fmt.Printf("Join id :")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(im.NewMessage(0, &pb.JoinReq{
			GroupID: int64(id),
		}))
	case 5:
		fmt.Printf("Quit id :")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(im.NewMessage(0, &pb.QuitReq{
			GroupID: int64(id),
		}))
	case 6:
		fmt.Printf("Send id,msg:")
		var id int
		var msg string
		fmt.Scan(&id, &msg)
		_ = sess.Send(im.NewMessage(0, &pb.SendMessageReq{
			GroupID: int64(id),
			Msg:     &pb.Message{Text: msg},
		}))
	case 7:
		fmt.Printf("GetMember id:")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(im.NewMessage(0, &pb.GetGroupMembersReq{
			GroupID: int64(id),
		}))
	case 8:
		fmt.Printf("Send id,startID,limit:")
		var id, startID, limit int
		fmt.Scan(&id, &startID, &limit)
		_ = sess.Send(im.NewMessage(0, &pb.SyncMessageReq{
			GroupID:  int64(id),
			StartID:  int64(startID),
			Limit:    int32(limit),
			OldToNew: false,
		}))
	case 9:
		fmt.Printf("Send :")
		_ = sess.Send(im.NewMessage(0, &pb.GetGroupListReq{}))
	}
}
