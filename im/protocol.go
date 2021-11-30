package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/pb"
	"io"
	"reflect"
)

type MessageType int

const (
	MESSAGE_UESR   MessageType = 1
	MESSAGE_GROUP  MessageType = 2
	MESSAGE_FRIEND MessageType = 3
)

type Message struct {
	data    interface{}
	cmd     uint16
	seq     uint32
	msgType MessageType
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

func (this *Message) GetSeq() uint32 {
	return this.seq
}

func (this *Message) GetType() MessageType {
	return this.msgType
}

const (
	seqSize       = 4
	cmdSize       = 2
	bodySize      = 4
	HeadSize      = seqSize + cmdSize + bodySize
	MaxBufferSize = 64 * 1024
)

type Codec struct{}

func readHeader(buff []byte) (uint32, uint16, uint32) {
	var seq, length uint32
	var cmd uint16

	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &seq)
	binary.Read(buffer, binary.BigEndian, &cmd)
	binary.Read(buffer, binary.BigEndian, &length)

	return seq, cmd, length
}

//解码
func (_ Codec) Decode(reader io.Reader) (interface{}, error) {
	hdr := make([]byte, HeadSize)
	if _, err := io.ReadFull(reader, hdr); err != nil {
		return nil, err
	}

	seq, cmd, length := readHeader(hdr)
	if length < 0 || length >= MaxBufferSize {
		return nil, fmt.Errorf("Message too large. ")
	}

	buff := make([]byte, length)
	if _, err := io.ReadFull(reader, buff); err != nil {
		return nil, err
	}

	msgType, msg, err := unmarshal(cmd, buff)
	if err != nil {
		return nil, err
	}

	return &Message{
		cmd:     cmd,
		seq:     seq,
		data:    msg,
		msgType: msgType,
	}, nil
}

//编码
func (_ Codec) Encode(o interface{}) ([]byte, error) {
	var seqNo uint32
	var cmd uint16
	var data []byte
	var err error

	switch o.(type) {
	case *Message:
		msg := o.(*Message)
		cmd, data, err = marshal(msg.GetData())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid type:%s. ", reflect.TypeOf(o).String())
	}

	length := uint32(len(data))
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, seqNo)
	binary.Write(buffer, binary.BigEndian, cmd)
	binary.Write(buffer, binary.BigEndian, length)
	buffer.Write(data)

	return buffer.Bytes(), nil
}

var (
	cmd2Type = map[uint16]*_protocol{}
	type2Cmd = map[reflect.Type]*_protocol{}
)

type _protocol struct {
	cmd     uint16
	tt      reflect.Type
	msgType MessageType
}

//根据名字注册实例(注意函数非线程安全，需要在初始化阶段完成所有消息的Register)
func register(msg proto.Message, id uint16, msgType MessageType) {
	if _, ok := cmd2Type[id]; ok {
		panic(fmt.Sprintf("id %d id areadly register. ", id))
	}

	tt := reflect.TypeOf(msg)
	_p := &_protocol{
		cmd:     id,
		tt:      tt,
		msgType: msgType,
	}

	cmd2Type[id] = _p
	type2Cmd[tt] = _p
}

func marshal(o interface{}) (uint16, []byte, error) {
	tt := reflect.TypeOf(o)
	_p, ok := type2Cmd[tt]
	if !ok {
		return 0, nil, fmt.Errorf("marshal type: %s undefined. ", reflect.TypeOf(o))
	}
	data, err := proto.Marshal(o.(proto.Message))
	return _p.cmd, data, err
}

func unmarshal(cmd uint16, buff []byte) (MessageType, interface{}, error) {
	_p, ok := cmd2Type[cmd]
	if !ok {
		return 0, nil, fmt.Errorf("unmarshal cmd: %d undefined. ", cmd)
	}

	//反序列化的结构
	msg := reflect.New(_p.tt.Elem()).Interface()
	err := proto.Unmarshal(buff, msg.(proto.Message))
	return _p.msgType, msg, err
}

func init() {
	register(&pb.UserLoginReq{}, uint16(pb.CmdType_CmdUserLoginReq), MESSAGE_UESR)
	register(&pb.UserLoginResp{}, uint16(pb.CmdType_CmdUserLoginResp), MESSAGE_UESR)

	register(&pb.CreateGroupReq{}, uint16(pb.CmdType_CmdCreateGroupReq), MESSAGE_GROUP)
	register(&pb.CreateGroupResp{}, uint16(pb.CmdType_CmdCreateGroupResp), MESSAGE_GROUP)
	register(&pb.GetGroupListReq{}, uint16(pb.CmdType_CmdGetGroupListReq), MESSAGE_GROUP)
	register(&pb.GetGroupListResp{}, uint16(pb.CmdType_CmdGetGroupListResp), MESSAGE_GROUP)
	register(&pb.DissolveGroupReq{}, uint16(pb.CmdType_CmdDissolveGroupReq), MESSAGE_GROUP)
	register(&pb.DissolveGroupResp{}, uint16(pb.CmdType_CmdDissolveGroupResp), MESSAGE_GROUP)
	register(&pb.NotifyDissolveGroup{}, uint16(pb.CmdType_CmdNotifyDissolveGroup), MESSAGE_GROUP)

	register(&pb.AddMemberReq{}, uint16(pb.CmdType_CmdAddMemberReq), MESSAGE_GROUP)
	register(&pb.AddMemberResp{}, uint16(pb.CmdType_CmdAddMemberResp), MESSAGE_GROUP)
	register(&pb.RemoveMemberReq{}, uint16(pb.CmdType_CmdRemoveMemberReq), MESSAGE_GROUP)
	register(&pb.RemoveMemberResp{}, uint16(pb.CmdType_CmdRemoveMemberResp), MESSAGE_GROUP)
	register(&pb.JoinReq{}, uint16(pb.CmdType_CmdJoinReq), MESSAGE_GROUP)
	register(&pb.JoinResp{}, uint16(pb.CmdType_CmdJoinResp), MESSAGE_GROUP)
	register(&pb.QuitReq{}, uint16(pb.CmdType_CmdQuitReq), MESSAGE_GROUP)
	register(&pb.QuitResp{}, uint16(pb.CmdType_CmdQuitResp), MESSAGE_GROUP)
	register(&pb.NotifyMemberJoined{}, uint16(pb.CmdType_CmdNotifyMemberJoined), MESSAGE_GROUP)
	register(&pb.NotifyMemberLeft{}, uint16(pb.CmdType_CmdNotifyMemberLeft), MESSAGE_GROUP)
	register(&pb.NotifyInvited{}, uint16(pb.CmdType_CmdNotifyInvited), MESSAGE_GROUP)
	register(&pb.NotifyKicked{}, uint16(pb.CmdType_CmdNotifyKicked), MESSAGE_GROUP)
	register(&pb.GetGroupMembersReq{}, uint16(pb.CmdType_CmdGetGroupMembersReq), MESSAGE_GROUP)
	register(&pb.GetGroupMembersResp{}, uint16(pb.CmdType_CmdGetGroupMembersResp), MESSAGE_GROUP)

	register(&pb.SendMessageReq{}, uint16(pb.CmdType_CmdSendMessageReq), MESSAGE_GROUP)
	register(&pb.SendMessageResp{}, uint16(pb.CmdType_CmdSendMessageResp), MESSAGE_GROUP)
	register(&pb.NotifyMessage{}, uint16(pb.CmdType_CmdNotifyMessage), MESSAGE_GROUP)
	register(&pb.SyncMessageReq{}, uint16(pb.CmdType_CmdSyncMessageReq), MESSAGE_GROUP)
	register(&pb.SyncMessageResp{}, uint16(pb.CmdType_CmdSyncMessageResp), MESSAGE_GROUP)
	register(&pb.RecallMessageReq{}, uint16(pb.CmdType_CmdRecallMessageReq), MESSAGE_GROUP)
	register(&pb.RecallMessageResp{}, uint16(pb.CmdType_CmdRecallMessageResp), MESSAGE_GROUP)
	register(&pb.NotifyRecallMessage{}, uint16(pb.CmdType_CmdNotifyRecallMessage), MESSAGE_GROUP)

	register(&pb.NotifyRecallMessage{}, uint16(pb.CmdType_CmdNotifyRecallMessage), MESSAGE_GROUP)

}
