package gate

import (
	"errors"
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/user"
)

var msgHandler = map[uint16]func(*user.User, *codec.Message){}

func RegisterHandler(cmd uint16, h func(*user.User, *codec.Message)) {
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
		u := ctx.(*user.User)
		if h, ok := msgHandler[cmd]; ok {
			h(u, msg)
		}
	}
}

func onUserLogin(session dnet.Session, msg *codec.Message) {
	//task.PostTask(func() {
	user.OnUserLogin(session, msg)
	//})
}

func onSessionClose(session dnet.Session, reason error) {
	user.OnClose(session, reason)
}
