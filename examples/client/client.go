package main

import (
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"os"
	"strings"
)

var handler = map[uint16]func(dnet.Session, *codec.Message){}

func registerHandler(cmdType pb.CmdType, h func(dnet.Session, *codec.Message)) {
	cmd := uint16(cmdType)
	if _, ok := handler[cmd]; ok {
		panic(fmt.Sprintf("cmd %d is alreadly register. ", cmd))
	}
	handler[cmd] = h
}

func dispatchMessage(sess dnet.Session, msg *codec.Message) {
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
		dnet.WithCodec(codec.NewCodec("im")),
		dnet.WithErrorCallback(func(session dnet.Session, err error) {
			fmt.Println("onError", err)
		}),
		dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
			msg := data.(*codec.Message)
			dispatchMessage(session, msg)
		}),
		dnet.WithCloseCallback(func(session dnet.Session, reason error) {
			fmt.Println("onClose", reason)
		}),
	)

	userID := os.Args[1]
	sess.Send(codec.NewMessage(1, &pb.UserLoginReq{
		ID: userID,
	}))

	registerHandler(pb.CmdType_CmdUserLoginResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("UserLoginResp %v", msg.GetData().(*pb.UserLoginResp))
	})
	registerHandler(pb.CmdType_CmdCreateConversationResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("CreateConversationResp %v", msg.GetData().(*pb.CreateConversationResp))
	})
	registerHandler(pb.CmdType_CmdAddMemberResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("AddMemberResp %v", msg.GetData().(*pb.AddMemberResp))
	})
	registerHandler(pb.CmdType_CmdRemoveMemberResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("RemoveMemberResp %v", msg.GetData().(*pb.RemoveMemberResp))
	})
	registerHandler(pb.CmdType_CmdJoinResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("JoinResp %v", msg.GetData().(*pb.JoinResp))
	})
	registerHandler(pb.CmdType_CmdQuitResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("QuitResp %v", msg.GetData().(*pb.QuitResp))
	})
	registerHandler(pb.CmdType_CmdSendMessageResp, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("SendMessageResp %v", msg.GetData().(*pb.SendMessageResp))
	})
	registerHandler(pb.CmdType_CmdNotifyInvited, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("NotifyInvited %v", msg.GetData().(*pb.NotifyInvited))
	})
	registerHandler(pb.CmdType_CmdNotifyKicked, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("NotifyKicked %v", msg.GetData().(*pb.NotifyKicked))
	})
	registerHandler(pb.CmdType_CmdNotifyMemberJoined, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("NotifyMemberJoined %v", msg.GetData().(*pb.NotifyMemberJoined))
	})
	registerHandler(pb.CmdType_CmdNotifyMemberLeft, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("NotifyMemberLeft %v", msg.GetData().(*pb.NotifyMemberLeft))
	})
	registerHandler(pb.CmdType_CmdNotifyMessage, func(session dnet.Session, msg *codec.Message) {
		log.Debugf("NotifyMessage %v", msg.GetData().(*pb.NotifyMessage))
	})

	go func() {
		for {
			cmd(sess)
		}
	}()
	select {}
}

func cmd(sess dnet.Session) {
	fmt.Println("1:createConversation 2:addMember 3:removeMember 4:join 5:quit 6:send")
	fmt.Printf("==>")
	var k int
	fmt.Scan(&k)
	switch k {
	case 1:
		fmt.Printf("CreateConversation users:")
		var users string
		fmt.Scan(&users)
		us := strings.Split(users, "&")
		_ = sess.Send(codec.NewMessage(0, &pb.CreateConversationReq{
			Members: us,
		}))
	case 2:
		fmt.Printf("AddMember id,users :")
		var id int
		var users string
		fmt.Scan(&id, &users)
		us := strings.Split(users, "&")
		_ = sess.Send(codec.NewMessage(0, &pb.AddMemberReq{
			ConvID: int64(id),
			AddIds: us,
		}))
	case 3:
		fmt.Printf("RemoveMember id,users :")
		var id int
		var users string
		fmt.Scan(&id, &users)
		us := strings.Split(users, "&")
		_ = sess.Send(codec.NewMessage(0, &pb.RemoveMemberReq{
			ConvID:    int64(id),
			RemoveIds: us,
		}))
	case 4:
		fmt.Printf("Join id :")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(codec.NewMessage(0, &pb.JoinReq{
			ConvID: int64(id),
		}))
	case 5:
		fmt.Printf("Quit id :")
		var id int
		fmt.Scan(&id)
		_ = sess.Send(codec.NewMessage(0, &pb.QuitReq{
			ConvID: int64(id),
		}))
	case 6:
		fmt.Printf("Send id,msg:")
		var id int
		var msg string
		fmt.Scan(&id, &msg)
		_ = sess.Send(codec.NewMessage(0, &pb.SendMessageReq{
			ConvID: int64(id),
			Msg:    &pb.Message{Text: msg},
		}))
	}
}
