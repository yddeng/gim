package gate

import (
	"errors"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
	"net"
)

func StartTCPGateway(address string) error {
	log.Infof("start tcp gateway on address:%s. ", address)
	return dnet.ServeTCPFunc(address, func(conn net.Conn) {
		log.Debug("new client", conn.RemoteAddr().String())

		_ = dnet.NewTCPSession(conn,
			//dnet.WithTimeout(time.Second*5, 0), // 超时
			dnet.WithCodec(codec.NewCodec("im")),
			//dnet.WithErrorCallback(func(session dnet.Session, err error) {
			//	fmt.Println("onError", err)
			//}),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				msg := data.(*codec.Message)
				switch msg.GetCmd() {
				case uint16(pb.CmdType_CmdUserLoginReq):
					user.OnUserLogin(session, msg)
				default:
					ctx := session.Context()
					if ctx == nil {
						session.Close(errors.New("user is not login. "))
						return
					}
					dispatchMessage(ctx.(*user.User), msg)
				}

			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				user.OnClose(session, reason)
			}))
	})
}
