package gate

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/codec/pb"
	"net"
	"time"
)

var msgHandler = map[uint16]func(dnet.Session, *codec.Message){}

func RegisterHandler(msg proto.Message, h func(dnet.Session, *codec.Message)) {
	cmd := pb.GetIdByName("im", proto.MessageName(msg))

	if _, ok := msgHandler[cmd]; ok {
		panic(fmt.Sprintf("cmd %d is alreadly register. ", cmd))
	}
	msgHandler[cmd] = h
}

func dispatchMessage(session dnet.Session, msg *codec.Message) {
	cmd := msg.GetCmd()
	if h, ok := msgHandler[cmd]; ok {
		h(session, msg)
	}
}

func StartTCPGateway(address string) error {
	return dnet.ServeTCPFunc(address, func(conn net.Conn) {
		fmt.Println("new client", conn.RemoteAddr().String())
		_ = dnet.NewTCPSession(conn,
			dnet.WithTimeout(time.Second*5, 0), // 超时
			dnet.WithCodec(codec.NewCodec("im", "", "")),
			dnet.WithErrorCallback(func(session dnet.Session, err error) {
				fmt.Println("onError", err)
			}),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				dispatchMessage(session, data.(*codec.Message))
			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				fmt.Println("onClose", reason)
			}))
	})
}
