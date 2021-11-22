package codec

import "github.com/golang/protobuf/proto"

type Message struct {
	data interface{}
	cmd  uint16
	seq  uint32
}

func NewMessage(seq uint32, data proto.Message) *Message {
	return &Message{seq: seq, data: data}
}

func (this *Message) GetData() interface{} {
	return this.data
}

func (this *Message) GetCmd() uint16 {
	return this.cmd
}
