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
	register(&pb.Heartbeat{}, uint16(pb.CmdType_CmdHeartbeat))
	register(&pb.UserLoginReq{}, uint16(pb.CmdType_CmdUserLoginReq))
	register(&pb.UserLoginResp{}, uint16(pb.CmdType_CmdUserLoginResp))

	register(&pb.CreateGroupReq{}, uint16(pb.CmdType_CmdCreateGroupReq))
	register(&pb.CreateGroupResp{}, uint16(pb.CmdType_CmdCreateGroupResp))
	register(&pb.GetGroupListReq{}, uint16(pb.CmdType_CmdGetGroupListReq))
	register(&pb.GetGroupListResp{}, uint16(pb.CmdType_CmdGetGroupListResp))
	register(&pb.DissolveGroupReq{}, uint16(pb.CmdType_CmdDissolveGroupReq))
	register(&pb.DissolveGroupResp{}, uint16(pb.CmdType_CmdDissolveGroupResp))
	register(&pb.NotifyDissolveGroup{}, uint16(pb.CmdType_CmdNotifyDissolveGroup))

	register(&pb.AddMemberReq{}, uint16(pb.CmdType_CmdAddMemberReq))
	register(&pb.AddMemberResp{}, uint16(pb.CmdType_CmdAddMemberResp))
	register(&pb.RemoveMemberReq{}, uint16(pb.CmdType_CmdRemoveMemberReq))
	register(&pb.RemoveMemberResp{}, uint16(pb.CmdType_CmdRemoveMemberResp))
	register(&pb.JoinReq{}, uint16(pb.CmdType_CmdJoinReq))
	register(&pb.JoinResp{}, uint16(pb.CmdType_CmdJoinResp))
	register(&pb.QuitReq{}, uint16(pb.CmdType_CmdQuitReq))
	register(&pb.QuitResp{}, uint16(pb.CmdType_CmdQuitResp))
	register(&pb.NotifyMemberJoined{}, uint16(pb.CmdType_CmdNotifyMemberJoined))
	register(&pb.NotifyMemberLeft{}, uint16(pb.CmdType_CmdNotifyMemberLeft))
	register(&pb.NotifyInvited{}, uint16(pb.CmdType_CmdNotifyInvited))
	register(&pb.NotifyKicked{}, uint16(pb.CmdType_CmdNotifyKicked))
	register(&pb.GetGroupMembersReq{}, uint16(pb.CmdType_CmdGetGroupMembersReq))
	register(&pb.GetGroupMembersResp{}, uint16(pb.CmdType_CmdGetGroupMembersResp))

	register(&pb.SendMessageReq{}, uint16(pb.CmdType_CmdSendMessageReq))
	register(&pb.SendMessageResp{}, uint16(pb.CmdType_CmdSendMessageResp))
	register(&pb.NotifyMessage{}, uint16(pb.CmdType_CmdNotifyMessage))
	register(&pb.SyncMessageReq{}, uint16(pb.CmdType_CmdSyncMessageReq))
	register(&pb.SyncMessageResp{}, uint16(pb.CmdType_CmdSyncMessageResp))
	register(&pb.RecallMessageReq{}, uint16(pb.CmdType_CmdRecallMessageReq))
	register(&pb.RecallMessageResp{}, uint16(pb.CmdType_CmdRecallMessageResp))
	register(&pb.NotifyRecallMessage{}, uint16(pb.CmdType_CmdNotifyRecallMessage))

	register(&pb.AddFriendReq{}, uint16(pb.CmdType_CmdAddFriendReq))
	register(&pb.AddFriendResp{}, uint16(pb.CmdType_CmdAddFriendResp))
	register(&pb.AgreeFriendReq{}, uint16(pb.CmdType_CmdAgreeFriendReq))
	register(&pb.AgreeFriendResp{}, uint16(pb.CmdType_CmdAgreeFriendResp))
	register(&pb.GetFriendsReq{}, uint16(pb.CmdType_CmdGetFriendsReq))
	register(&pb.GetFriendsResp{}, uint16(pb.CmdType_CmdGetFriendsResp))
	register(&pb.NotifyAddFriend{}, uint16(pb.CmdType_CmdNotifyAddFriend))
	register(&pb.NotifyAgreeFriend{}, uint16(pb.CmdType_CmdNotifyAgreeFriend))

}
