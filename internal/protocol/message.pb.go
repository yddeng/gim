// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Conversation struct {
	ID                   uint64   `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Conversation) Reset()         { *m = Conversation{} }
func (m *Conversation) String() string { return proto.CompactTextString(m) }
func (*Conversation) ProtoMessage()    {}
func (*Conversation) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Conversation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Conversation.Unmarshal(m, b)
}
func (m *Conversation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Conversation.Marshal(b, m, deterministic)
}
func (m *Conversation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Conversation.Merge(m, src)
}
func (m *Conversation) XXX_Size() int {
	return xxx_messageInfo_Conversation.Size(m)
}
func (m *Conversation) XXX_DiscardUnknown() {
	xxx_messageInfo_Conversation.DiscardUnknown(m)
}

var xxx_messageInfo_Conversation proto.InternalMessageInfo

func (m *Conversation) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Conversation) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// 创建对话
type CreateConversationReq struct {
	Members              []string `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Unique               bool     `protobuf:"varint,3,opt,name=unique,proto3" json:"unique,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateConversationReq) Reset()         { *m = CreateConversationReq{} }
func (m *CreateConversationReq) String() string { return proto.CompactTextString(m) }
func (*CreateConversationReq) ProtoMessage()    {}
func (*CreateConversationReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *CreateConversationReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateConversationReq.Unmarshal(m, b)
}
func (m *CreateConversationReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateConversationReq.Marshal(b, m, deterministic)
}
func (m *CreateConversationReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateConversationReq.Merge(m, src)
}
func (m *CreateConversationReq) XXX_Size() int {
	return xxx_messageInfo_CreateConversationReq.Size(m)
}
func (m *CreateConversationReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateConversationReq.DiscardUnknown(m)
}

var xxx_messageInfo_CreateConversationReq proto.InternalMessageInfo

func (m *CreateConversationReq) GetMembers() []string {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *CreateConversationReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateConversationReq) GetUnique() bool {
	if m != nil {
		return m.Unique
	}
	return false
}

type CreateConversationResp struct {
	Ok                   bool          `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Conv                 *Conversation `protobuf:"bytes,2,opt,name=conv,proto3" json:"conv,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CreateConversationResp) Reset()         { *m = CreateConversationResp{} }
func (m *CreateConversationResp) String() string { return proto.CompactTextString(m) }
func (*CreateConversationResp) ProtoMessage()    {}
func (*CreateConversationResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *CreateConversationResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateConversationResp.Unmarshal(m, b)
}
func (m *CreateConversationResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateConversationResp.Marshal(b, m, deterministic)
}
func (m *CreateConversationResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateConversationResp.Merge(m, src)
}
func (m *CreateConversationResp) XXX_Size() int {
	return xxx_messageInfo_CreateConversationResp.Size(m)
}
func (m *CreateConversationResp) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateConversationResp.DiscardUnknown(m)
}

var xxx_messageInfo_CreateConversationResp proto.InternalMessageInfo

func (m *CreateConversationResp) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *CreateConversationResp) GetConv() *Conversation {
	if m != nil {
		return m.Conv
	}
	return nil
}

type NotifyInvited struct {
	Conv                 *Conversation `protobuf:"bytes,1,opt,name=conv,proto3" json:"conv,omitempty"`
	InitBy               string        `protobuf:"bytes,2,opt,name=initBy,proto3" json:"initBy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *NotifyInvited) Reset()         { *m = NotifyInvited{} }
func (m *NotifyInvited) String() string { return proto.CompactTextString(m) }
func (*NotifyInvited) ProtoMessage()    {}
func (*NotifyInvited) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}

func (m *NotifyInvited) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyInvited.Unmarshal(m, b)
}
func (m *NotifyInvited) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyInvited.Marshal(b, m, deterministic)
}
func (m *NotifyInvited) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyInvited.Merge(m, src)
}
func (m *NotifyInvited) XXX_Size() int {
	return xxx_messageInfo_NotifyInvited.Size(m)
}
func (m *NotifyInvited) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyInvited.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyInvited proto.InternalMessageInfo

func (m *NotifyInvited) GetConv() *Conversation {
	if m != nil {
		return m.Conv
	}
	return nil
}

func (m *NotifyInvited) GetInitBy() string {
	if m != nil {
		return m.InitBy
	}
	return ""
}

// 发送消息
type Message struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type SendMessageReq struct {
	Msg                  *Message `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMessageReq) Reset()         { *m = SendMessageReq{} }
func (m *SendMessageReq) String() string { return proto.CompactTextString(m) }
func (*SendMessageReq) ProtoMessage()    {}
func (*SendMessageReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{5}
}

func (m *SendMessageReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMessageReq.Unmarshal(m, b)
}
func (m *SendMessageReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMessageReq.Marshal(b, m, deterministic)
}
func (m *SendMessageReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageReq.Merge(m, src)
}
func (m *SendMessageReq) XXX_Size() int {
	return xxx_messageInfo_SendMessageReq.Size(m)
}
func (m *SendMessageReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageReq.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageReq proto.InternalMessageInfo

func (m *SendMessageReq) GetMsg() *Message {
	if m != nil {
		return m.Msg
	}
	return nil
}

type SendMessageResp struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMessageResp) Reset()         { *m = SendMessageResp{} }
func (m *SendMessageResp) String() string { return proto.CompactTextString(m) }
func (*SendMessageResp) ProtoMessage()    {}
func (*SendMessageResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{6}
}

func (m *SendMessageResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMessageResp.Unmarshal(m, b)
}
func (m *SendMessageResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMessageResp.Marshal(b, m, deterministic)
}
func (m *SendMessageResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageResp.Merge(m, src)
}
func (m *SendMessageResp) XXX_Size() int {
	return xxx_messageInfo_SendMessageResp.Size(m)
}
func (m *SendMessageResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageResp.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageResp proto.InternalMessageInfo

func (m *SendMessageResp) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

type NotifyMessage struct {
	Conv                 *Conversation `protobuf:"bytes,1,opt,name=conv,proto3" json:"conv,omitempty"`
	Msg                  *Message      `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *NotifyMessage) Reset()         { *m = NotifyMessage{} }
func (m *NotifyMessage) String() string { return proto.CompactTextString(m) }
func (*NotifyMessage) ProtoMessage()    {}
func (*NotifyMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{7}
}

func (m *NotifyMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyMessage.Unmarshal(m, b)
}
func (m *NotifyMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyMessage.Marshal(b, m, deterministic)
}
func (m *NotifyMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyMessage.Merge(m, src)
}
func (m *NotifyMessage) XXX_Size() int {
	return xxx_messageInfo_NotifyMessage.Size(m)
}
func (m *NotifyMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyMessage proto.InternalMessageInfo

func (m *NotifyMessage) GetConv() *Conversation {
	if m != nil {
		return m.Conv
	}
	return nil
}

func (m *NotifyMessage) GetMsg() *Message {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*Conversation)(nil), "Conversation")
	proto.RegisterType((*CreateConversationReq)(nil), "CreateConversationReq")
	proto.RegisterType((*CreateConversationResp)(nil), "CreateConversationResp")
	proto.RegisterType((*NotifyInvited)(nil), "NotifyInvited")
	proto.RegisterType((*Message)(nil), "Message")
	proto.RegisterType((*SendMessageReq)(nil), "SendMessageReq")
	proto.RegisterType((*SendMessageResp)(nil), "SendMessageResp")
	proto.RegisterType((*NotifyMessage)(nil), "NotifyMessage")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0xcf, 0x6b, 0xfa, 0x30,
	0x14, 0x27, 0x55, 0xb4, 0x7d, 0xdf, 0x6f, 0x1d, 0x04, 0x26, 0x65, 0x30, 0xa8, 0x39, 0xf5, 0x30,
	0x7a, 0x70, 0xff, 0x81, 0x7a, 0xe9, 0xc6, 0x7a, 0xc8, 0x6e, 0x83, 0x1d, 0xaa, 0xbe, 0x49, 0x70,
	0x4d, 0xb4, 0x89, 0x65, 0xfe, 0xf7, 0x23, 0x59, 0x3a, 0x94, 0x39, 0xd8, 0xed, 0xf3, 0xf2, 0xde,
	0xe7, 0xc7, 0x7b, 0x81, 0xb8, 0x46, 0xad, 0xab, 0x0d, 0xe6, 0xbb, 0x46, 0x19, 0xc5, 0xa6, 0xf0,
	0x7f, 0xae, 0x64, 0x8b, 0x8d, 0xae, 0x8c, 0x50, 0x92, 0x8e, 0x20, 0x28, 0x16, 0x09, 0x49, 0x49,
	0xd6, 0xe7, 0x41, 0xb1, 0xa0, 0x14, 0xfa, 0x65, 0x55, 0x63, 0x12, 0xa4, 0x24, 0x8b, 0xb8, 0xc3,
	0xec, 0x15, 0xae, 0xe7, 0x0d, 0x56, 0x06, 0x4f, 0x99, 0x1c, 0xf7, 0x34, 0x81, 0x61, 0x8d, 0xf5,
	0x12, 0x1b, 0x9d, 0x90, 0xb4, 0x97, 0x45, 0xbc, 0x2b, 0xad, 0x8c, 0x3c, 0x91, 0xb1, 0x98, 0x8e,
	0x61, 0x70, 0x90, 0x62, 0x7f, 0xc0, 0xa4, 0x97, 0x92, 0x2c, 0xe4, 0xbe, 0x62, 0x8f, 0x30, 0xbe,
	0x24, 0xaf, 0x77, 0x36, 0x9c, 0xda, 0xba, 0x70, 0x21, 0x0f, 0xd4, 0x96, 0x4e, 0xa0, 0xbf, 0x52,
	0xb2, 0x75, 0xaa, 0xff, 0xa6, 0x71, 0x7e, 0x46, 0x70, 0x2d, 0xf6, 0x00, 0x71, 0xa9, 0x8c, 0x78,
	0x3b, 0x16, 0xb2, 0x15, 0x06, 0xd7, 0xdf, 0x1c, 0xf2, 0x2b, 0xc7, 0x06, 0x13, 0x52, 0x98, 0xd9,
	0xd1, 0xc7, 0xf5, 0x15, 0xbb, 0x85, 0xe1, 0xd3, 0xd7, 0xf1, 0xec, 0x3e, 0x06, 0x3f, 0x8c, 0x53,
	0x89, 0xb8, 0xc3, 0xec, 0x0e, 0x46, 0xcf, 0x28, 0xd7, 0x7e, 0xc4, 0xde, 0xe3, 0x06, 0x7a, 0xb5,
	0xde, 0x78, 0xab, 0x30, 0xef, 0x3a, 0xf6, 0x91, 0x4d, 0xe0, 0xea, 0x6c, 0xfa, 0xe7, 0x7a, 0xac,
	0xec, 0xb2, 0x77, 0xae, 0x7f, 0xc8, 0xee, 0x2d, 0x83, 0x0b, 0x96, 0x33, 0x78, 0x09, 0xdd, 0xa7,
	0xaf, 0xd4, 0xfb, 0x72, 0xe0, 0xd0, 0xfd, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4d, 0xf1, 0xf6,
	0xc6, 0x0f, 0x02, 0x00, 0x00,
}
