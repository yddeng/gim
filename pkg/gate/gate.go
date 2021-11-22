package gate

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/codec/pb"
	"github.com/yddeng/gim/pkg/user"
	"net"
	"time"
)

var (
	msgHandler = map[uint16]func(*user.User, *codec.Message){}
)

func RegisterHandler(msg proto.Message, h func(*user.User, *codec.Message)) {
	cmd := pb.GetIdByName("im", proto.MessageName(msg))

	if _, ok := msgHandler[cmd]; ok {
		panic(fmt.Sprintf("cmd %d is alreadly register. ", cmd))
	}
	msgHandler[cmd] = h
}

func dispatchMessage(u *user.User, msg *codec.Message) {
	cmd := msg.GetCmd()
	if h, ok := msgHandler[cmd]; ok {
		h(u, msg)
	}
}

func StartTCPGateway(address string) error {
	return dnet.ServeTCPFunc(address, func(conn net.Conn) {
		fmt.Println("new client", conn.RemoteAddr().String())

		_ = dnet.NewTCPSession(conn,
			dnet.WithTimeout(time.Second*5, 0), // 超时
			dnet.WithCodec(codec.NewCodec("im")),
			dnet.WithErrorCallback(func(session dnet.Session, err error) {
				fmt.Println("onError", err)
			}),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				msg := data.(*codec.Message)
				switch msg.GetCmd() {
				case 1:
					user.OnUserLogin(session, msg)
				default:
					ctx := session.Context()
					if ctx == nil {
						session.Close(errors.New("no user"))
						return
					}
					dispatchMessage(ctx.(*user.User), msg)
				}

			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				fmt.Println("onClose", reason)
			}))
	})
}
