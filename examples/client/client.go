package main

import (
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol"
)

func initHandler() {

}

func dispatchMessage(sess dnet.Session, msg *codec.Message) {
	fmt.Println(msg.GetData())
	switch msg.GetData().(type) {
	case *protocol.UserLoginResp:
		createConversation(sess)
	case *protocol.CreateConversationResp:
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

	login(sess)

	select {}
}

func login(session dnet.Session) {
	fmt.Printf("账号：")
	var id string
	fmt.Scan(&id)
	session.Send(codec.NewMessage(1, &protocol.UserLoginReq{
		ID: id,
	}))
}

func createConversation(session dnet.Session) {
	session.Send(codec.NewMessage(1, &protocol.CreateConversationReq{
		Members: []string{"111"},
	}))
}
