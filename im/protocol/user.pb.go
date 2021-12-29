// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol/user.proto

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

type UserLoginReq struct {
	ID                   *string  `protobuf:"bytes,1,req,name=ID" json:"ID,omitempty"`
	Extras               []*Extra `protobuf:"bytes,2,rep,name=extras" json:"extras,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserLoginReq) Reset()         { *m = UserLoginReq{} }
func (m *UserLoginReq) String() string { return proto.CompactTextString(m) }
func (*UserLoginReq) ProtoMessage()    {}
func (*UserLoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe658c3c35e8c4bc, []int{0}
}

func (m *UserLoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLoginReq.Unmarshal(m, b)
}
func (m *UserLoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLoginReq.Marshal(b, m, deterministic)
}
func (m *UserLoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLoginReq.Merge(m, src)
}
func (m *UserLoginReq) XXX_Size() int {
	return xxx_messageInfo_UserLoginReq.Size(m)
}
func (m *UserLoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserLoginReq proto.InternalMessageInfo

func (m *UserLoginReq) GetID() string {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return ""
}

func (m *UserLoginReq) GetExtras() []*Extra {
	if m != nil {
		return m.Extras
	}
	return nil
}

type UserLoginResp struct {
	Code                 *ErrCode `protobuf:"varint,1,opt,name=code,enum=message.IM.ErrCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserLoginResp) Reset()         { *m = UserLoginResp{} }
func (m *UserLoginResp) String() string { return proto.CompactTextString(m) }
func (*UserLoginResp) ProtoMessage()    {}
func (*UserLoginResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe658c3c35e8c4bc, []int{1}
}

func (m *UserLoginResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLoginResp.Unmarshal(m, b)
}
func (m *UserLoginResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLoginResp.Marshal(b, m, deterministic)
}
func (m *UserLoginResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLoginResp.Merge(m, src)
}
func (m *UserLoginResp) XXX_Size() int {
	return xxx_messageInfo_UserLoginResp.Size(m)
}
func (m *UserLoginResp) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLoginResp.DiscardUnknown(m)
}

var xxx_messageInfo_UserLoginResp proto.InternalMessageInfo

func (m *UserLoginResp) GetCode() ErrCode {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ErrCode_OK
}

type Heartbeat struct {
	Timestamp            *int64   `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Heartbeat) Reset()         { *m = Heartbeat{} }
func (m *Heartbeat) String() string { return proto.CompactTextString(m) }
func (*Heartbeat) ProtoMessage()    {}
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe658c3c35e8c4bc, []int{2}
}

func (m *Heartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Heartbeat.Unmarshal(m, b)
}
func (m *Heartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Heartbeat.Marshal(b, m, deterministic)
}
func (m *Heartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Heartbeat.Merge(m, src)
}
func (m *Heartbeat) XXX_Size() int {
	return xxx_messageInfo_Heartbeat.Size(m)
}
func (m *Heartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_Heartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_Heartbeat proto.InternalMessageInfo

func (m *Heartbeat) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

type User struct {
	ID                   *string  `protobuf:"bytes,1,req,name=ID" json:"ID,omitempty"`
	Extras               []*Extra `protobuf:"bytes,2,rep,name=extras" json:"extras,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe658c3c35e8c4bc, []int{3}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetID() string {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return ""
}

func (m *User) GetExtras() []*Extra {
	if m != nil {
		return m.Extras
	}
	return nil
}

type GetUserInfoReq struct {
	UserIDs              []string `protobuf:"bytes,1,rep,name=userIDs" json:"userIDs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserInfoReq) Reset()         { *m = GetUserInfoReq{} }
func (m *GetUserInfoReq) String() string { return proto.CompactTextString(m) }
func (*GetUserInfoReq) ProtoMessage()    {}
func (*GetUserInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe658c3c35e8c4bc, []int{4}
}

func (m *GetUserInfoReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserInfoReq.Unmarshal(m, b)
}
func (m *GetUserInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserInfoReq.Marshal(b, m, deterministic)
}
func (m *GetUserInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserInfoReq.Merge(m, src)
}
func (m *GetUserInfoReq) XXX_Size() int {
	return xxx_messageInfo_GetUserInfoReq.Size(m)
}
func (m *GetUserInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserInfoReq proto.InternalMessageInfo

func (m *GetUserInfoReq) GetUserIDs() []string {
	if m != nil {
		return m.UserIDs
	}
	return nil
}

type GetUserInfoResp struct {
	Users                []*User  `protobuf:"bytes,1,rep,name=users" json:"users,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserInfoResp) Reset()         { *m = GetUserInfoResp{} }
func (m *GetUserInfoResp) String() string { return proto.CompactTextString(m) }
func (*GetUserInfoResp) ProtoMessage()    {}
func (*GetUserInfoResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe658c3c35e8c4bc, []int{5}
}

func (m *GetUserInfoResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserInfoResp.Unmarshal(m, b)
}
func (m *GetUserInfoResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserInfoResp.Marshal(b, m, deterministic)
}
func (m *GetUserInfoResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserInfoResp.Merge(m, src)
}
func (m *GetUserInfoResp) XXX_Size() int {
	return xxx_messageInfo_GetUserInfoResp.Size(m)
}
func (m *GetUserInfoResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserInfoResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserInfoResp proto.InternalMessageInfo

func (m *GetUserInfoResp) GetUsers() []*User {
	if m != nil {
		return m.Users
	}
	return nil
}

func init() {
	proto.RegisterType((*UserLoginReq)(nil), "message.IM.UserLoginReq")
	proto.RegisterType((*UserLoginResp)(nil), "message.IM.UserLoginResp")
	proto.RegisterType((*Heartbeat)(nil), "message.IM.Heartbeat")
	proto.RegisterType((*User)(nil), "message.IM.User")
	proto.RegisterType((*GetUserInfoReq)(nil), "message.IM.GetUserInfoReq")
	proto.RegisterType((*GetUserInfoResp)(nil), "message.IM.GetUserInfoResp")
}

func init() { proto.RegisterFile("protocol/user.proto", fileDescriptor_fe658c3c35e8c4bc) }

var fileDescriptor_fe658c3c35e8c4bc = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x95, 0xa4, 0xfc, 0xc9, 0x15, 0x02, 0xb8, 0x12, 0xb2, 0x10, 0x43, 0xe4, 0x01, 0x52,
	0x86, 0x20, 0x65, 0x82, 0x11, 0x08, 0x02, 0x4b, 0xb0, 0x58, 0x62, 0x61, 0x33, 0xed, 0x51, 0x55,
	0x22, 0x75, 0xb8, 0x33, 0x12, 0x1f, 0x1f, 0x39, 0xa1, 0xb4, 0x9d, 0x19, 0xcf, 0xef, 0x77, 0xef,
	0xf9, 0x1e, 0x8c, 0x5a, 0x72, 0xde, 0x4d, 0xdc, 0xc7, 0xe5, 0x17, 0x23, 0x95, 0xdd, 0x24, 0xa0,
	0x41, 0x66, 0x3b, 0xc3, 0x52, 0x3f, 0x9f, 0x1c, 0xff, 0x01, 0x48, 0x34, 0x71, 0x53, 0xec, 0x19,
	0xa5, 0x61, 0xef, 0x85, 0x91, 0x9e, 0xdc, 0x6c, 0xbe, 0x30, 0xf8, 0x29, 0x32, 0x88, 0x75, 0x2d,
	0xa3, 0x3c, 0x2e, 0x52, 0x13, 0xeb, 0x5a, 0x8c, 0x61, 0x1b, 0xbf, 0x3d, 0x59, 0x96, 0x71, 0x9e,
	0x14, 0xc3, 0xea, 0xa8, 0x5c, 0x99, 0x96, 0xf7, 0x41, 0x31, 0xbf, 0x80, 0xba, 0x82, 0xfd, 0x35,
	0x2b, 0x6e, 0xc5, 0x39, 0x0c, 0x42, 0x92, 0x8c, 0xf2, 0xa8, 0xc8, 0xaa, 0xd1, 0xc6, 0x26, 0xd1,
	0x9d, 0x9b, 0xa2, 0xe9, 0x00, 0x35, 0x86, 0xf4, 0x11, 0x2d, 0xf9, 0x37, 0xb4, 0x5e, 0x9c, 0x42,
	0xea, 0xe7, 0x0d, 0xb2, 0xb7, 0x4d, 0xdb, 0xad, 0x26, 0x66, 0xf5, 0xa0, 0x6e, 0x60, 0x10, 0x42,
	0xfe, 0xf3, 0xcf, 0x0b, 0xc8, 0x1e, 0xd0, 0x07, 0x17, 0xbd, 0x78, 0x77, 0xe1, 0x68, 0x09, 0x3b,
	0xa1, 0x36, 0x5d, 0xb3, 0x8c, 0xf2, 0xa4, 0x48, 0xcd, 0x72, 0x54, 0xd7, 0x70, 0xb0, 0xc1, 0x72,
	0x2b, 0xce, 0x60, 0x2b, 0xa8, 0x3d, 0x3a, 0xac, 0x0e, 0xd7, 0x83, 0x02, 0x68, 0x7a, 0xf9, 0x16,
	0x5e, 0x77, 0x97, 0x9d, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x47, 0x1a, 0x67, 0x02, 0x9f, 0x01,
	0x00, 0x00,
}
