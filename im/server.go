package im

import (
	"errors"
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
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
			ctx := session.Context()
			if ctx != nil {
				u := ctx.(*User)
				log.Infof("onClose user(%s) %s. ", u.ID, reason)
				u.sess.SetContext(nil)
				u.sess = nil
			}
		}))
}

var (
	userHandler  = map[uint16]func(session dnet.Session, msg *codec.Message){}
	groupHandler = map[uint16]func(*User, *codec.Message){}
)

func registerGroupHandler(cmd uint16, h func(*User, *codec.Message)) {
	if _, ok := groupHandler[cmd]; ok {
		panic(fmt.Sprintf("group handler cmd %d is alreadly register. ", cmd))
	}
	groupHandler[cmd] = h
}

func registerUserHandler(cmd uint16, h func(session dnet.Session, msg *codec.Message)) {
	if _, ok := userHandler[cmd]; ok {
		panic(fmt.Sprintf("user handler cmd %d is alreadly register. ", cmd))
	}
	userHandler[cmd] = h
}

func dispatchMessage(session dnet.Session, msg *codec.Message) {
	cmd := msg.GetCmd()

	if h, ok := userHandler[cmd]; ok {
		h(session, msg)
	} else if h2, ok := groupHandler[cmd]; ok {
		ctx := session.Context()
		if ctx == nil {
			session.Close(errors.New("user is not login. "))
			return
		}
		h2(ctx.(*User), msg)
	}
}
