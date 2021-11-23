package user

import "github.com/golang/protobuf/proto"

var (
	hooks = map[string]func(user *User, msg proto.Message){}
)

func registerHook(name string, h func(user *User, msg proto.Message)) {

}

func disposeHook(user *User, msg proto.Message) {

}

func init() {

}
