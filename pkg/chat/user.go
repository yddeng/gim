package gim

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
)

var (
	sess2User = map[dnet.Session]*User{}
	users     = map[string]*User{}
)

type User struct {
	ID   string
	sess dnet.Session
}

func (this *User) Reply(msg proto.Message) {
	codec.NewMessage()
}
