package im

import (
	"errors"
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"net"
)

func StartTCPGateway(address string) error {
	log.Infof("start tcp gateway on address:%s. ", address)
	return dnet.ServeTCPFunc(address, func(conn net.Conn) {
		log.Debug("new tcp client", conn.RemoteAddr().String())
		_ = createSession(conn)
	})
}

func StartWSGateway(address string) error {
	log.Infof("start ws gateway on address:%s. ", address)
	return dnet.ServeWSFunc(address, func(conn net.Conn) {
		log.Debug("new ws client", conn.RemoteAddr().String())
		_ = createSession(conn)
	})
}

func createSession(conn net.Conn) dnet.Session {
	return dnet.NewTCPSession(conn,
		//dnet.WithTimeout(time.Second*5, 0), // 超时
		dnet.WithCodec(codec.NewCodec("im")),
		//dnet.WithErrorCallback(func(session dnet.Session, err error) {
		//	fmt.Println("onError", err)
		//}),
		dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
			dispatchMessage(session, data.(*codec.Message))
		}),
		dnet.WithCloseCallback(func(session dnet.Session, reason error) {
			onSessionClose(session, reason)
		}))
}

var msgHandler = map[uint16]func(*User, *codec.Message){}

func registerHandler(cmd uint16, h func(*User, *codec.Message)) {
	if _, ok := msgHandler[cmd]; ok {
		panic(fmt.Sprintf("cmd %d is alreadly register. ", cmd))
	}
	msgHandler[cmd] = h
}

func dispatchMessage(session dnet.Session, msg *codec.Message) {
	switch msg.GetCmd() {
	case uint16(pb.CmdType_CmdUserLoginReq):
		onUserLogin(session, msg)
	default:
		ctx := session.Context()
		if ctx == nil {
			session.Close(errors.New("user is not login. "))
			return
		}

		cmd := msg.GetCmd()
		u := ctx.(*User)
		if h, ok := msgHandler[cmd]; ok {
			h(u, msg)
		}
	}
}
