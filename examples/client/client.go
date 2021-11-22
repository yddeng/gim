package main

import (
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/internal/codec"
	"time"
)

func main() {
	address := "127.0.0.1:43210"
	c, err := dnet.DialTCP(address, 0)
	if err != nil {
		panic(err)
	}

	sess := dnet.NewTCPSession(c,
		dnet.WithTimeout(time.Second*5, 0), // 超时
		dnet.WithCodec(codec.NewCodec("im")),
		dnet.WithErrorCallback(func(session dnet.Session, err error) {
			fmt.Println("onError", err)
		}),
		dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
			msg := data.(*codec.Message)
			dispatchMessage(session, msg)
		}),
		dnet.WithCloseCallback(func(session dnet.Session, reason error) {
			fmt.Println("onClose", reason)
		}),
	)

}
