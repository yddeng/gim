package gate

import (
	"fmt"
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/pkg/user"
)

var msgHandler = map[uint16]func(*user.User, *codec.Message){}

func RegisterHandler(cmd uint16, h func(*user.User, *codec.Message)) {
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
