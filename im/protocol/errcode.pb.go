// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol/errcode.proto

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

type ErrCode int32

const (
	ErrCode_OK                    ErrCode = 0
	ErrCode_Error                 ErrCode = 1
	ErrCode_Busy                  ErrCode = 2
	ErrCode_UserAlreadyLogin      ErrCode = 10
	ErrCode_GroupNotExist         ErrCode = 11
	ErrCode_UserNotInGroup        ErrCode = 12
	ErrCode_UserNotExist          ErrCode = 13
	ErrCode_UserNotHasPermission  ErrCode = 14
	ErrCode_RequestArgumentErr    ErrCode = 1001
	ErrCode_FriendAlreadyIsFriend ErrCode = 1101
	ErrCode_FriendApplyClosed     ErrCode = 1102
)

var ErrCode_name = map[int32]string{
	0:    "OK",
	1:    "Error",
	2:    "Busy",
	10:   "UserAlreadyLogin",
	11:   "GroupNotExist",
	12:   "UserNotInGroup",
	13:   "UserNotExist",
	14:   "UserNotHasPermission",
	1001: "RequestArgumentErr",
	1101: "FriendAlreadyIsFriend",
	1102: "FriendApplyClosed",
}

var ErrCode_value = map[string]int32{
	"OK":                    0,
	"Error":                 1,
	"Busy":                  2,
	"UserAlreadyLogin":      10,
	"GroupNotExist":         11,
	"UserNotInGroup":        12,
	"UserNotExist":          13,
	"UserNotHasPermission":  14,
	"RequestArgumentErr":    1001,
	"FriendAlreadyIsFriend": 1101,
	"FriendApplyClosed":     1102,
}

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}

func (x *ErrCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ErrCode_value, data, "ErrCode")
	if err != nil {
		return err
	}
	*x = ErrCode(value)
	return nil
}

func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7a5460248a277878, []int{0}
}

type Extra struct {
	Key                  *string  `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value                *string  `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Extra) Reset()         { *m = Extra{} }
func (m *Extra) String() string { return proto.CompactTextString(m) }
func (*Extra) ProtoMessage()    {}
func (*Extra) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a5460248a277878, []int{0}
}

func (m *Extra) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Extra.Unmarshal(m, b)
}
func (m *Extra) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Extra.Marshal(b, m, deterministic)
}
func (m *Extra) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Extra.Merge(m, src)
}
func (m *Extra) XXX_Size() int {
	return xxx_messageInfo_Extra.Size(m)
}
func (m *Extra) XXX_DiscardUnknown() {
	xxx_messageInfo_Extra.DiscardUnknown(m)
}

var xxx_messageInfo_Extra proto.InternalMessageInfo

func (m *Extra) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Extra) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func init() {
	proto.RegisterEnum("ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterType((*Extra)(nil), "Extra")
}

func init() { proto.RegisterFile("protocol/errcode.proto", fileDescriptor_7a5460248a277878) }

var fileDescriptor_7a5460248a277878 = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x90, 0xc1, 0x4e, 0xf2, 0x50,
	0x10, 0x85, 0xff, 0xf6, 0xb7, 0x52, 0x46, 0x20, 0xc3, 0xa4, 0x62, 0xe3, 0xca, 0xb8, 0x32, 0x2e,
	0xe4, 0x19, 0x80, 0x54, 0x25, 0x1a, 0x34, 0x24, 0x6e, 0xdc, 0x35, 0x74, 0x42, 0x1a, 0x4b, 0xa7,
	0xce, 0xbd, 0xd7, 0xd0, 0x87, 0xd4, 0x27, 0xf0, 0x05, 0x7c, 0x0b, 0x43, 0x5b, 0x77, 0x73, 0xbe,
	0xf3, 0x2d, 0x26, 0x07, 0x26, 0x95, 0x8a, 0x95, 0x8d, 0x14, 0x53, 0x56, 0xdd, 0x48, 0xc6, 0x37,
	0x0d, 0xb8, 0x9c, 0x42, 0x90, 0xec, 0xad, 0xa6, 0x84, 0xf0, 0xff, 0x8d, 0xeb, 0xd8, 0xbb, 0xf0,
	0xaf, 0xfa, 0xeb, 0xc3, 0x49, 0x11, 0x04, 0x1f, 0x69, 0xe1, 0x38, 0xf6, 0x1b, 0xd6, 0x86, 0xeb,
	0x6f, 0x0f, 0x7a, 0x89, 0xea, 0x42, 0x32, 0xa6, 0x63, 0xf0, 0x9f, 0x1e, 0xf0, 0x1f, 0xf5, 0x21,
	0x48, 0x54, 0x45, 0xd1, 0xa3, 0x10, 0x8e, 0xe6, 0xce, 0xd4, 0xe8, 0x53, 0x04, 0xf8, 0x62, 0x58,
	0x67, 0x85, 0x72, 0x9a, 0xd5, 0x8f, 0xb2, 0xcd, 0x4b, 0x04, 0x1a, 0xc3, 0xf0, 0x4e, 0xc5, 0x55,
	0x2b, 0xb1, 0xc9, 0x3e, 0x37, 0x16, 0x4f, 0x88, 0x60, 0x74, 0x10, 0x57, 0x62, 0x97, 0x65, 0xd3,
	0xe1, 0x80, 0x10, 0x06, 0x1d, 0x6b, 0xad, 0x21, 0xc5, 0x10, 0x75, 0xe4, 0x3e, 0x35, 0xcf, 0xac,
	0xbb, 0xdc, 0x98, 0x5c, 0x4a, 0x1c, 0xd1, 0x19, 0xd0, 0x9a, 0xdf, 0x1d, 0x1b, 0x3b, 0xd3, 0xad,
	0xdb, 0x71, 0x69, 0x13, 0x55, 0xfc, 0xe9, 0xd1, 0x39, 0x9c, 0xde, 0x6a, 0xce, 0x65, 0xd6, 0xfd,
	0xb0, 0x34, 0x6d, 0xc4, 0xcf, 0x90, 0x26, 0x30, 0xee, 0xba, 0xaa, 0x2a, 0xea, 0x45, 0x21, 0x86,
	0x33, 0xfc, 0x0a, 0xe7, 0xf0, 0x1a, 0xfe, 0x2d, 0xf5, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x57, 0xd8,
	0x10, 0x40, 0x34, 0x01, 0x00, 0x00,
}
