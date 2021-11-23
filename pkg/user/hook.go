package user

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/protocol/pb"
)

var (
	hooks = map[string]func(user *User, msg proto.Message){}
)

func registerHook(name string, h func(user *User, msg proto.Message)) {
	hooks[name] = h
}

func disposeHook(user *User, msg proto.Message) {
	name := proto.MessageName(msg)
	if h, ok := hooks[name]; ok {
		h(user, msg)
	}
}

func onNotifyKicked(user *User, msg proto.Message) {

}

func init() {
	registerHook(proto.MessageName(&pb.NotifyKicked{}), onNotifyKicked)
}
