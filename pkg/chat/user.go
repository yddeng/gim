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

func (this *User) Reply(seq uint32, msg proto.Message) {
	this.sess.Send(codec.NewMessage(seq, msg))
}

func getUser(id string) *User {
	return users[id]
}
