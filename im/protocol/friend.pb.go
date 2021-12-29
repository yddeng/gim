// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol/friend.proto

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

type FriendStatus int32

const (
	FriendStatus_Agree FriendStatus = 0
	FriendStatus_Apply FriendStatus = 1
)

var FriendStatus_name = map[int32]string{
	0: "Agree",
	1: "Apply",
}

var FriendStatus_value = map[string]int32{
	"Agree": 0,
	"Apply": 1,
}

func (x FriendStatus) Enum() *FriendStatus {
	p := new(FriendStatus)
	*p = x
	return p
}

func (x FriendStatus) String() string {
	return proto.EnumName(FriendStatus_name, int32(x))
}

func (x *FriendStatus) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FriendStatus_value, data, "FriendStatus")
	if err != nil {
		return err
	}
	*x = FriendStatus(value)
	return nil
}

func (FriendStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{0}
}

type Friend struct {
	Status               *FriendStatus `protobuf:"varint,1,opt,name=status,enum=message.IM.FriendStatus" json:"status,omitempty"`
	UserID               *string       `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	Extras               []*Extra      `protobuf:"bytes,3,rep,name=extras" json:"extras,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Friend) Reset()         { *m = Friend{} }
func (m *Friend) String() string { return proto.CompactTextString(m) }
func (*Friend) ProtoMessage()    {}
func (*Friend) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{0}
}

func (m *Friend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Friend.Unmarshal(m, b)
}
func (m *Friend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Friend.Marshal(b, m, deterministic)
}
func (m *Friend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Friend.Merge(m, src)
}
func (m *Friend) XXX_Size() int {
	return xxx_messageInfo_Friend.Size(m)
}
func (m *Friend) XXX_DiscardUnknown() {
	xxx_messageInfo_Friend.DiscardUnknown(m)
}

var xxx_messageInfo_Friend proto.InternalMessageInfo

func (m *Friend) GetStatus() FriendStatus {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return FriendStatus_Agree
}

func (m *Friend) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *Friend) GetExtras() []*Extra {
	if m != nil {
		return m.Extras
	}
	return nil
}

type AddFriendReq struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddFriendReq) Reset()         { *m = AddFriendReq{} }
func (m *AddFriendReq) String() string { return proto.CompactTextString(m) }
func (*AddFriendReq) ProtoMessage()    {}
func (*AddFriendReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{1}
}

func (m *AddFriendReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddFriendReq.Unmarshal(m, b)
}
func (m *AddFriendReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddFriendReq.Marshal(b, m, deterministic)
}
func (m *AddFriendReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddFriendReq.Merge(m, src)
}
func (m *AddFriendReq) XXX_Size() int {
	return xxx_messageInfo_AddFriendReq.Size(m)
}
func (m *AddFriendReq) XXX_DiscardUnknown() {
	xxx_messageInfo_AddFriendReq.DiscardUnknown(m)
}

var xxx_messageInfo_AddFriendReq proto.InternalMessageInfo

func (m *AddFriendReq) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

type AddFriendResp struct {
	Code                 *ErrCode `protobuf:"varint,1,opt,name=code,enum=message.IM.ErrCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddFriendResp) Reset()         { *m = AddFriendResp{} }
func (m *AddFriendResp) String() string { return proto.CompactTextString(m) }
func (*AddFriendResp) ProtoMessage()    {}
func (*AddFriendResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{2}
}

func (m *AddFriendResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddFriendResp.Unmarshal(m, b)
}
func (m *AddFriendResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddFriendResp.Marshal(b, m, deterministic)
}
func (m *AddFriendResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddFriendResp.Merge(m, src)
}
func (m *AddFriendResp) XXX_Size() int {
	return xxx_messageInfo_AddFriendResp.Size(m)
}
func (m *AddFriendResp) XXX_DiscardUnknown() {
	xxx_messageInfo_AddFriendResp.DiscardUnknown(m)
}

var xxx_messageInfo_AddFriendResp proto.InternalMessageInfo

func (m *AddFriendResp) GetCode() ErrCode {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ErrCode_OK
}

type NotifyAddFriend struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyAddFriend) Reset()         { *m = NotifyAddFriend{} }
func (m *NotifyAddFriend) String() string { return proto.CompactTextString(m) }
func (*NotifyAddFriend) ProtoMessage()    {}
func (*NotifyAddFriend) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{3}
}

func (m *NotifyAddFriend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyAddFriend.Unmarshal(m, b)
}
func (m *NotifyAddFriend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyAddFriend.Marshal(b, m, deterministic)
}
func (m *NotifyAddFriend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyAddFriend.Merge(m, src)
}
func (m *NotifyAddFriend) XXX_Size() int {
	return xxx_messageInfo_NotifyAddFriend.Size(m)
}
func (m *NotifyAddFriend) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyAddFriend.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyAddFriend proto.InternalMessageInfo

func (m *NotifyAddFriend) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

type AgreeFriendReq struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	Agree                *bool    `protobuf:"varint,2,opt,name=agree" json:"agree,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AgreeFriendReq) Reset()         { *m = AgreeFriendReq{} }
func (m *AgreeFriendReq) String() string { return proto.CompactTextString(m) }
func (*AgreeFriendReq) ProtoMessage()    {}
func (*AgreeFriendReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{4}
}

func (m *AgreeFriendReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgreeFriendReq.Unmarshal(m, b)
}
func (m *AgreeFriendReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgreeFriendReq.Marshal(b, m, deterministic)
}
func (m *AgreeFriendReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgreeFriendReq.Merge(m, src)
}
func (m *AgreeFriendReq) XXX_Size() int {
	return xxx_messageInfo_AgreeFriendReq.Size(m)
}
func (m *AgreeFriendReq) XXX_DiscardUnknown() {
	xxx_messageInfo_AgreeFriendReq.DiscardUnknown(m)
}

var xxx_messageInfo_AgreeFriendReq proto.InternalMessageInfo

func (m *AgreeFriendReq) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *AgreeFriendReq) GetAgree() bool {
	if m != nil && m.Agree != nil {
		return *m.Agree
	}
	return false
}

type AgreeFriendResp struct {
	Code                 *ErrCode `protobuf:"varint,1,opt,name=code,enum=message.IM.ErrCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AgreeFriendResp) Reset()         { *m = AgreeFriendResp{} }
func (m *AgreeFriendResp) String() string { return proto.CompactTextString(m) }
func (*AgreeFriendResp) ProtoMessage()    {}
func (*AgreeFriendResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{5}
}

func (m *AgreeFriendResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgreeFriendResp.Unmarshal(m, b)
}
func (m *AgreeFriendResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgreeFriendResp.Marshal(b, m, deterministic)
}
func (m *AgreeFriendResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgreeFriendResp.Merge(m, src)
}
func (m *AgreeFriendResp) XXX_Size() int {
	return xxx_messageInfo_AgreeFriendResp.Size(m)
}
func (m *AgreeFriendResp) XXX_DiscardUnknown() {
	xxx_messageInfo_AgreeFriendResp.DiscardUnknown(m)
}

var xxx_messageInfo_AgreeFriendResp proto.InternalMessageInfo

func (m *AgreeFriendResp) GetCode() ErrCode {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ErrCode_OK
}

type NotifyAgreeFriend struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	Agree                *bool    `protobuf:"varint,2,opt,name=agree" json:"agree,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyAgreeFriend) Reset()         { *m = NotifyAgreeFriend{} }
func (m *NotifyAgreeFriend) String() string { return proto.CompactTextString(m) }
func (*NotifyAgreeFriend) ProtoMessage()    {}
func (*NotifyAgreeFriend) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{6}
}

func (m *NotifyAgreeFriend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyAgreeFriend.Unmarshal(m, b)
}
func (m *NotifyAgreeFriend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyAgreeFriend.Marshal(b, m, deterministic)
}
func (m *NotifyAgreeFriend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyAgreeFriend.Merge(m, src)
}
func (m *NotifyAgreeFriend) XXX_Size() int {
	return xxx_messageInfo_NotifyAgreeFriend.Size(m)
}
func (m *NotifyAgreeFriend) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyAgreeFriend.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyAgreeFriend proto.InternalMessageInfo

func (m *NotifyAgreeFriend) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *NotifyAgreeFriend) GetAgree() bool {
	if m != nil && m.Agree != nil {
		return *m.Agree
	}
	return false
}

type DeleteFriendReq struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteFriendReq) Reset()         { *m = DeleteFriendReq{} }
func (m *DeleteFriendReq) String() string { return proto.CompactTextString(m) }
func (*DeleteFriendReq) ProtoMessage()    {}
func (*DeleteFriendReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{7}
}

func (m *DeleteFriendReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteFriendReq.Unmarshal(m, b)
}
func (m *DeleteFriendReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteFriendReq.Marshal(b, m, deterministic)
}
func (m *DeleteFriendReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteFriendReq.Merge(m, src)
}
func (m *DeleteFriendReq) XXX_Size() int {
	return xxx_messageInfo_DeleteFriendReq.Size(m)
}
func (m *DeleteFriendReq) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteFriendReq.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteFriendReq proto.InternalMessageInfo

func (m *DeleteFriendReq) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

type DeleteFriendResp struct {
	Code                 *ErrCode `protobuf:"varint,1,opt,name=code,enum=message.IM.ErrCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteFriendResp) Reset()         { *m = DeleteFriendResp{} }
func (m *DeleteFriendResp) String() string { return proto.CompactTextString(m) }
func (*DeleteFriendResp) ProtoMessage()    {}
func (*DeleteFriendResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{8}
}

func (m *DeleteFriendResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteFriendResp.Unmarshal(m, b)
}
func (m *DeleteFriendResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteFriendResp.Marshal(b, m, deterministic)
}
func (m *DeleteFriendResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteFriendResp.Merge(m, src)
}
func (m *DeleteFriendResp) XXX_Size() int {
	return xxx_messageInfo_DeleteFriendResp.Size(m)
}
func (m *DeleteFriendResp) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteFriendResp.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteFriendResp proto.InternalMessageInfo

func (m *DeleteFriendResp) GetCode() ErrCode {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ErrCode_OK
}

type NotifyDeleteFriend struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyDeleteFriend) Reset()         { *m = NotifyDeleteFriend{} }
func (m *NotifyDeleteFriend) String() string { return proto.CompactTextString(m) }
func (*NotifyDeleteFriend) ProtoMessage()    {}
func (*NotifyDeleteFriend) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{9}
}

func (m *NotifyDeleteFriend) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyDeleteFriend.Unmarshal(m, b)
}
func (m *NotifyDeleteFriend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyDeleteFriend.Marshal(b, m, deterministic)
}
func (m *NotifyDeleteFriend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyDeleteFriend.Merge(m, src)
}
func (m *NotifyDeleteFriend) XXX_Size() int {
	return xxx_messageInfo_NotifyDeleteFriend.Size(m)
}
func (m *NotifyDeleteFriend) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyDeleteFriend.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyDeleteFriend proto.InternalMessageInfo

func (m *NotifyDeleteFriend) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

type GetFriendsReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFriendsReq) Reset()         { *m = GetFriendsReq{} }
func (m *GetFriendsReq) String() string { return proto.CompactTextString(m) }
func (*GetFriendsReq) ProtoMessage()    {}
func (*GetFriendsReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{10}
}

func (m *GetFriendsReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFriendsReq.Unmarshal(m, b)
}
func (m *GetFriendsReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFriendsReq.Marshal(b, m, deterministic)
}
func (m *GetFriendsReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFriendsReq.Merge(m, src)
}
func (m *GetFriendsReq) XXX_Size() int {
	return xxx_messageInfo_GetFriendsReq.Size(m)
}
func (m *GetFriendsReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFriendsReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetFriendsReq proto.InternalMessageInfo

type GetFriendsResp struct {
	Friends              []*Friend `protobuf:"bytes,1,rep,name=friends" json:"friends,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetFriendsResp) Reset()         { *m = GetFriendsResp{} }
func (m *GetFriendsResp) String() string { return proto.CompactTextString(m) }
func (*GetFriendsResp) ProtoMessage()    {}
func (*GetFriendsResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_92c06ba2c1dbe2eb, []int{11}
}

func (m *GetFriendsResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFriendsResp.Unmarshal(m, b)
}
func (m *GetFriendsResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFriendsResp.Marshal(b, m, deterministic)
}
func (m *GetFriendsResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFriendsResp.Merge(m, src)
}
func (m *GetFriendsResp) XXX_Size() int {
	return xxx_messageInfo_GetFriendsResp.Size(m)
}
func (m *GetFriendsResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFriendsResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetFriendsResp proto.InternalMessageInfo

func (m *GetFriendsResp) GetFriends() []*Friend {
	if m != nil {
		return m.Friends
	}
	return nil
}

func init() {
	proto.RegisterEnum("message.IM.FriendStatus", FriendStatus_name, FriendStatus_value)
	proto.RegisterType((*Friend)(nil), "message.IM.Friend")
	proto.RegisterType((*AddFriendReq)(nil), "message.IM.AddFriendReq")
	proto.RegisterType((*AddFriendResp)(nil), "message.IM.AddFriendResp")
	proto.RegisterType((*NotifyAddFriend)(nil), "message.IM.NotifyAddFriend")
	proto.RegisterType((*AgreeFriendReq)(nil), "message.IM.AgreeFriendReq")
	proto.RegisterType((*AgreeFriendResp)(nil), "message.IM.AgreeFriendResp")
	proto.RegisterType((*NotifyAgreeFriend)(nil), "message.IM.NotifyAgreeFriend")
	proto.RegisterType((*DeleteFriendReq)(nil), "message.IM.DeleteFriendReq")
	proto.RegisterType((*DeleteFriendResp)(nil), "message.IM.DeleteFriendResp")
	proto.RegisterType((*NotifyDeleteFriend)(nil), "message.IM.NotifyDeleteFriend")
	proto.RegisterType((*GetFriendsReq)(nil), "message.IM.GetFriendsReq")
	proto.RegisterType((*GetFriendsResp)(nil), "message.IM.GetFriendsResp")
}

func init() { proto.RegisterFile("protocol/friend.proto", fileDescriptor_92c06ba2c1dbe2eb) }

var fileDescriptor_92c06ba2c1dbe2eb = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xc1, 0x4b, 0x02, 0x41,
	0x14, 0xc6, 0xdb, 0xcc, 0x4d, 0x5f, 0xea, 0xea, 0x54, 0xb2, 0x74, 0x92, 0x21, 0x4a, 0x43, 0xb6,
	0xf0, 0x14, 0x05, 0x81, 0x65, 0x85, 0x87, 0x3a, 0x4c, 0xb7, 0x6e, 0x8b, 0xfb, 0x14, 0xc1, 0x9a,
	0x75, 0x66, 0x84, 0x3c, 0xf4, 0xbf, 0xc7, 0xcc, 0xac, 0x3a, 0x12, 0x22, 0xde, 0xf6, 0xed, 0xfc,
	0xde, 0xf7, 0x7d, 0x33, 0x7c, 0x70, 0x9a, 0x0a, 0xae, 0xf8, 0x80, 0x4f, 0xae, 0x87, 0x62, 0x8c,
	0xdf, 0x49, 0x64, 0x66, 0x02, 0x5f, 0x28, 0x65, 0x3c, 0xc2, 0xa8, 0xff, 0x76, 0x56, 0x5f, 0x22,
	0x28, 0xc4, 0x80, 0x27, 0x68, 0x19, 0xfa, 0x0b, 0xfe, 0x8b, 0xd9, 0x21, 0x37, 0xe0, 0x4b, 0x15,
	0xab, 0x99, 0x0c, 0xbd, 0x86, 0xd7, 0xac, 0x74, 0xc2, 0x68, 0xb5, 0x1e, 0x59, 0xe6, 0xc3, 0x9c,
	0xb3, 0x8c, 0x23, 0x75, 0xf0, 0x67, 0x12, 0x45, 0xbf, 0x17, 0xee, 0x37, 0xbc, 0x66, 0x91, 0x65,
	0x13, 0x69, 0x81, 0x8f, 0x3f, 0x4a, 0xc4, 0x32, 0xcc, 0x35, 0x72, 0xcd, 0xa3, 0x4e, 0xcd, 0x55,
	0x7a, 0xd6, 0x27, 0x2c, 0x03, 0xe8, 0x05, 0x94, 0xba, 0x49, 0x62, 0xd5, 0x19, 0x4e, 0x1d, 0x49,
	0xcf, 0x95, 0xa4, 0xb7, 0x50, 0x76, 0x38, 0x99, 0x92, 0x4b, 0x38, 0xd0, 0xb7, 0xc8, 0xb2, 0x1e,
	0xaf, 0x39, 0x08, 0xf1, 0xc4, 0x13, 0x64, 0x06, 0xa0, 0x2d, 0x08, 0xde, 0xb9, 0x1a, 0x0f, 0xe7,
	0xcb, 0xfd, 0x8d, 0x26, 0x0f, 0x50, 0xe9, 0x8e, 0x04, 0xe2, 0xd6, 0x38, 0xe4, 0x04, 0xf2, 0xb1,
	0x26, 0xcd, 0xc5, 0x0b, 0xcc, 0x0e, 0xf4, 0x0e, 0x82, 0xb5, 0xfd, 0x5d, 0x62, 0x76, 0xa1, 0x96,
	0xc5, 0x5c, 0x29, 0xec, 0x68, 0xdf, 0x82, 0xa0, 0x87, 0x13, 0x54, 0xdb, 0xf3, 0xd3, 0x7b, 0xa8,
	0xae, 0xa3, 0xbb, 0x44, 0x6d, 0x03, 0xb1, 0x51, 0x5d, 0x89, 0x8d, 0x56, 0x01, 0x94, 0x5f, 0x51,
	0x59, 0x48, 0x32, 0x9c, 0xea, 0x57, 0x76, 0x7f, 0xc8, 0x94, 0xb4, 0xe1, 0xd0, 0xf6, 0x56, 0x57,
	0x4f, 0x17, 0x86, 0xfc, 0xaf, 0x1e, 0x5b, 0x20, 0x57, 0xe7, 0x50, 0x72, 0xdb, 0x48, 0x8a, 0x90,
	0x37, 0x6f, 0x56, 0xdd, 0x33, 0x9f, 0x69, 0x3a, 0x99, 0x57, 0xbd, 0x47, 0xf8, 0x2c, 0x2c, 0x1a,
	0xff, 0x17, 0x00, 0x00, 0xff, 0xff, 0x10, 0xaa, 0x3c, 0x19, 0x1f, 0x03, 0x00, 0x00,
}
