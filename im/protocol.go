package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/protocol"
	"io"
	"reflect"
)

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

func (this *Message) GetSeq() uint32 {
	return this.seq
}

const (
	SizeLen   = 2
	SizeSeq   = 4
	SizeFlag  = 2
	SizeCmd   = 2
	SizeTeach = 2 // 教学引导ID
	SizeErr   = 2
	SizeHead  = SizeLen + SizeSeq + SizeFlag + SizeCmd + SizeTeach + SizeErr
)

type Codec struct{}

func readHeader(buff []byte) (uint16, uint32, uint16) {
	var seq uint32
	var cmd, length uint16

	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &length)
	binary.Read(buffer, binary.BigEndian, &seq)
	binary.Read(buffer, binary.BigEndian, &cmd)
	binary.Read(buffer, binary.BigEndian, &cmd)

	return length - (SizeHead - SizeLen), seq, cmd
}

//解码
func (_ Codec) Decode(reader io.Reader) (interface{}, error) {
	hdr := make([]byte, SizeHead)
	if _, err := io.ReadFull(reader, hdr); err != nil {
		return nil, err
	}

	length, seq, cmd := readHeader(hdr)
	if length < 0 {
		return nil, fmt.Errorf("Message too large. ")
	}

	buff := make([]byte, length)
	if _, err := io.ReadFull(reader, buff); err != nil {
		return nil, err
	}

	msg, err := unmarshal(cmd, buff)
	if err != nil {
		return nil, err
	}

	return &Message{
		cmd:  cmd,
		seq:  seq,
		data: msg,
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
		seqNo = msg.seq
		cmd, data, err = marshal(msg.GetData())
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid type:%s. ", reflect.TypeOf(o).String())
	}

	length := uint16(len(data))
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, uint16(length+(SizeHead-SizeLen)))
	binary.Write(buffer, binary.BigEndian, seqNo)
	binary.Write(buffer, binary.BigEndian, uint16(0))
	binary.Write(buffer, binary.BigEndian, cmd)
	binary.Write(buffer, binary.BigEndian, uint16(0))
	binary.Write(buffer, binary.BigEndian, uint16(0))
	buffer.Write(data)

	return buffer.Bytes(), nil
}

var (
	cmd2Type = map[uint16]reflect.Type{}
	type2Cmd = map[reflect.Type]uint16{}
)

//根据名字注册实例(注意函数非线程安全，需要在初始化阶段完成所有消息的Register)
func register(msg proto.Message, id uint16) {
	if _, ok := cmd2Type[id]; ok {
		panic(fmt.Sprintf("id %d id areadly register. ", id))
	}

	tt := reflect.TypeOf(msg)
	cmd2Type[id] = tt
	type2Cmd[tt] = id
}

func marshal(o interface{}) (uint16, []byte, error) {
	tt := reflect.TypeOf(o)
	id, ok := type2Cmd[tt]
	if !ok {
		return 0, nil, fmt.Errorf("marshal type: %s undefined. ", reflect.TypeOf(o))
	}
	data, err := proto.Marshal(o.(proto.Message))
	return id, data, err
}

func unmarshal(cmd uint16, buff []byte) (interface{}, error) {
	tt, ok := cmd2Type[cmd]
	if !ok {
		return nil, fmt.Errorf("unmarshal cmd: %d undefined. ", cmd)
	}

	//反序列化的结构
	msg := reflect.New(tt.Elem()).Interface()
	err := proto.Unmarshal(buff, msg.(proto.Message))
	return msg, err
}

func init() {
	register(&protocol.Heartbeat{}, uint16(protocol.CmdType_CmdHeartbeat))
	register(&protocol.UserLoginReq{}, uint16(protocol.CmdType_CmdUserLoginReq))
	register(&protocol.UserLoginResp{}, uint16(protocol.CmdType_CmdUserLoginResp))
	register(&protocol.GetUserInfoReq{}, uint16(protocol.CmdType_CmdGetUserInfoReq))
	register(&protocol.GetUserInfoResp{}, uint16(protocol.CmdType_CmdGetUserInfoResp))

	register(&protocol.CreateGroupReq{}, uint16(protocol.CmdType_CmdCreateGroupReq))
	register(&protocol.CreateGroupResp{}, uint16(protocol.CmdType_CmdCreateGroupResp))
	register(&protocol.GetGroupListReq{}, uint16(protocol.CmdType_CmdGetGroupListReq))
	register(&protocol.GetGroupListResp{}, uint16(protocol.CmdType_CmdGetGroupListResp))
	register(&protocol.DissolveGroupReq{}, uint16(protocol.CmdType_CmdDissolveGroupReq))
	register(&protocol.DissolveGroupResp{}, uint16(protocol.CmdType_CmdDissolveGroupResp))
	register(&protocol.NotifyDissolveGroup{}, uint16(protocol.CmdType_CmdNotifyDissolveGroup))

	register(&protocol.AddMemberReq{}, uint16(protocol.CmdType_CmdAddMemberReq))
	register(&protocol.AddMemberResp{}, uint16(protocol.CmdType_CmdAddMemberResp))
	register(&protocol.RemoveMemberReq{}, uint16(protocol.CmdType_CmdRemoveMemberReq))
	register(&protocol.RemoveMemberResp{}, uint16(protocol.CmdType_CmdRemoveMemberResp))
	register(&protocol.JoinReq{}, uint16(protocol.CmdType_CmdJoinReq))
	register(&protocol.JoinResp{}, uint16(protocol.CmdType_CmdJoinResp))
	register(&protocol.QuitReq{}, uint16(protocol.CmdType_CmdQuitReq))
	register(&protocol.QuitResp{}, uint16(protocol.CmdType_CmdQuitResp))
	register(&protocol.NotifyMemberJoined{}, uint16(protocol.CmdType_CmdNotifyMemberJoined))
	register(&protocol.NotifyMemberLeft{}, uint16(protocol.CmdType_CmdNotifyMemberLeft))
	register(&protocol.NotifyInvited{}, uint16(protocol.CmdType_CmdNotifyInvited))
	register(&protocol.NotifyKicked{}, uint16(protocol.CmdType_CmdNotifyKicked))
	register(&protocol.GetGroupMembersReq{}, uint16(protocol.CmdType_CmdGetGroupMembersReq))
	register(&protocol.GetGroupMembersResp{}, uint16(protocol.CmdType_CmdGetGroupMembersResp))

	register(&protocol.SendMessageReq{}, uint16(protocol.CmdType_CmdSendMessageReq))
	register(&protocol.SendMessageResp{}, uint16(protocol.CmdType_CmdSendMessageResp))
	register(&protocol.NotifyMessage{}, uint16(protocol.CmdType_CmdNotifyMessage))
	register(&protocol.SyncMessageReq{}, uint16(protocol.CmdType_CmdSyncMessageReq))
	register(&protocol.SyncMessageResp{}, uint16(protocol.CmdType_CmdSyncMessageResp))
	register(&protocol.RecallMessageReq{}, uint16(protocol.CmdType_CmdRecallMessageReq))
	register(&protocol.RecallMessageResp{}, uint16(protocol.CmdType_CmdRecallMessageResp))
	register(&protocol.NotifyRecallMessage{}, uint16(protocol.CmdType_CmdNotifyRecallMessage))

	register(&protocol.AddFriendReq{}, uint16(protocol.CmdType_CmdAddFriendReq))
	register(&protocol.AddFriendResp{}, uint16(protocol.CmdType_CmdAddFriendResp))
	register(&protocol.AgreeFriendReq{}, uint16(protocol.CmdType_CmdAgreeFriendReq))
	register(&protocol.AgreeFriendResp{}, uint16(protocol.CmdType_CmdAgreeFriendResp))
	register(&protocol.GetFriendsReq{}, uint16(protocol.CmdType_CmdGetFriendsReq))
	register(&protocol.GetFriendsResp{}, uint16(protocol.CmdType_CmdGetFriendsResp))
	register(&protocol.NotifyAddFriend{}, uint16(protocol.CmdType_CmdNotifyAddFriend))
	register(&protocol.NotifyAgreeFriend{}, uint16(protocol.CmdType_CmdNotifyAgreeFriend))
	register(&protocol.DeleteFriendReq{}, uint16(protocol.CmdType_CmdDeleteFriendReq))
	register(&protocol.DeleteFriendResp{}, uint16(protocol.CmdType_CmdDeleteFriendResp))
	register(&protocol.NotifyDeleteFriend{}, uint16(protocol.CmdType_CmdNotifyDeleteFriend))
	register(&protocol.NotifyUserOnline{}, uint16(protocol.CmdType_CmdNotifyUserOnline))

}
