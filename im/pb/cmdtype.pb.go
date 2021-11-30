// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmdtype.proto

package pb

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

type CmdType int32

const (
	CmdType_CmdNone                CmdType = 0
	CmdType_CmdUserLoginReq        CmdType = 1
	CmdType_CmdUserLoginResp       CmdType = 2
	CmdType_CmdCreateGroupReq      CmdType = 101
	CmdType_CmdCreateGroupResp     CmdType = 102
	CmdType_CmdGetGroupMembersReq  CmdType = 111
	CmdType_CmdGetGroupMembersResp CmdType = 112
	CmdType_CmdGetGroupListReq     CmdType = 121
	CmdType_CmdGetGroupListResp    CmdType = 122
	CmdType_CmdDissolveGroupReq    CmdType = 123
	CmdType_CmdDissolveGroupResp   CmdType = 124
	CmdType_CmdNotifyDissolveGroup CmdType = 125
	CmdType_CmdAddMemberReq        CmdType = 201
	CmdType_CmdAddMemberResp       CmdType = 202
	CmdType_CmdRemoveMemberReq     CmdType = 203
	CmdType_CmdRemoveMemberResp    CmdType = 204
	CmdType_CmdJoinReq             CmdType = 211
	CmdType_CmdJoinResp            CmdType = 212
	CmdType_CmdQuitReq             CmdType = 213
	CmdType_CmdQuitResp            CmdType = 224
	CmdType_CmdNotifyInvited       CmdType = 231
	CmdType_CmdNotifyMemberJoined  CmdType = 232
	CmdType_CmdNotifyKicked        CmdType = 233
	CmdType_CmdNotifyMemberLeft    CmdType = 234
	CmdType_CmdSendMessageReq      CmdType = 301
	CmdType_CmdSendMessageResp     CmdType = 302
	CmdType_CmdNotifyMessage       CmdType = 303
	CmdType_CmdSyncMessageReq      CmdType = 311
	CmdType_CmdSyncMessageResp     CmdType = 312
	CmdType_CmdRecallMessageReq    CmdType = 313
	CmdType_CmdRecallMessageResp   CmdType = 314
	CmdType_CmdNotifyRecallMessage CmdType = 315
	CmdType_CmdAddFriendReq        CmdType = 401
	CmdType_CmdAddFriendResp       CmdType = 402
	CmdType_CmdAgreeFriendReq      CmdType = 403
	CmdType_CmdAgreeFriendResp     CmdType = 404
	CmdType_CmdGetFriendsReq       CmdType = 405
	CmdType_CmdGetFriendsResp      CmdType = 406
	CmdType_CmdNotifyAddFriend     CmdType = 407
	CmdType_CmdNotifyAgreeFriend   CmdType = 408
)

var CmdType_name = map[int32]string{
	0:   "CmdNone",
	1:   "CmdUserLoginReq",
	2:   "CmdUserLoginResp",
	101: "CmdCreateGroupReq",
	102: "CmdCreateGroupResp",
	111: "CmdGetGroupMembersReq",
	112: "CmdGetGroupMembersResp",
	121: "CmdGetGroupListReq",
	122: "CmdGetGroupListResp",
	123: "CmdDissolveGroupReq",
	124: "CmdDissolveGroupResp",
	125: "CmdNotifyDissolveGroup",
	201: "CmdAddMemberReq",
	202: "CmdAddMemberResp",
	203: "CmdRemoveMemberReq",
	204: "CmdRemoveMemberResp",
	211: "CmdJoinReq",
	212: "CmdJoinResp",
	213: "CmdQuitReq",
	224: "CmdQuitResp",
	231: "CmdNotifyInvited",
	232: "CmdNotifyMemberJoined",
	233: "CmdNotifyKicked",
	234: "CmdNotifyMemberLeft",
	301: "CmdSendMessageReq",
	302: "CmdSendMessageResp",
	303: "CmdNotifyMessage",
	311: "CmdSyncMessageReq",
	312: "CmdSyncMessageResp",
	313: "CmdRecallMessageReq",
	314: "CmdRecallMessageResp",
	315: "CmdNotifyRecallMessage",
	401: "CmdAddFriendReq",
	402: "CmdAddFriendResp",
	403: "CmdAgreeFriendReq",
	404: "CmdAgreeFriendResp",
	405: "CmdGetFriendsReq",
	406: "CmdGetFriendsResp",
	407: "CmdNotifyAddFriend",
	408: "CmdNotifyAgreeFriend",
}

var CmdType_value = map[string]int32{
	"CmdNone":                0,
	"CmdUserLoginReq":        1,
	"CmdUserLoginResp":       2,
	"CmdCreateGroupReq":      101,
	"CmdCreateGroupResp":     102,
	"CmdGetGroupMembersReq":  111,
	"CmdGetGroupMembersResp": 112,
	"CmdGetGroupListReq":     121,
	"CmdGetGroupListResp":    122,
	"CmdDissolveGroupReq":    123,
	"CmdDissolveGroupResp":   124,
	"CmdNotifyDissolveGroup": 125,
	"CmdAddMemberReq":        201,
	"CmdAddMemberResp":       202,
	"CmdRemoveMemberReq":     203,
	"CmdRemoveMemberResp":    204,
	"CmdJoinReq":             211,
	"CmdJoinResp":            212,
	"CmdQuitReq":             213,
	"CmdQuitResp":            224,
	"CmdNotifyInvited":       231,
	"CmdNotifyMemberJoined":  232,
	"CmdNotifyKicked":        233,
	"CmdNotifyMemberLeft":    234,
	"CmdSendMessageReq":      301,
	"CmdSendMessageResp":     302,
	"CmdNotifyMessage":       303,
	"CmdSyncMessageReq":      311,
	"CmdSyncMessageResp":     312,
	"CmdRecallMessageReq":    313,
	"CmdRecallMessageResp":   314,
	"CmdNotifyRecallMessage": 315,
	"CmdAddFriendReq":        401,
	"CmdAddFriendResp":       402,
	"CmdAgreeFriendReq":      403,
	"CmdAgreeFriendResp":     404,
	"CmdGetFriendsReq":       405,
	"CmdGetFriendsResp":      406,
	"CmdNotifyAddFriend":     407,
	"CmdNotifyAgreeFriend":   408,
}

func (x CmdType) String() string {
	return proto.EnumName(CmdType_name, int32(x))
}

func (CmdType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_df7da9134e4cbc1b, []int{0}
}

func init() {
	proto.RegisterEnum("CmdType", CmdType_name, CmdType_value)
}

func init() { proto.RegisterFile("cmdtype.proto", fileDescriptor_df7da9134e4cbc1b) }

var fileDescriptor_df7da9134e4cbc1b = []byte{
	// 489 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0xcb, 0x6e, 0xd3, 0x50,
	0x10, 0x86, 0xb1, 0x83, 0xa8, 0x74, 0x2a, 0xd4, 0x61, 0x9a, 0xa4, 0x17, 0xde, 0x80, 0x05, 0x1b,
	0x9e, 0xa0, 0x18, 0x11, 0x01, 0x69, 0x25, 0x0a, 0x6c, 0xd8, 0x25, 0x39, 0x93, 0xc8, 0x22, 0xf6,
	0x19, 0x3c, 0x6e, 0x24, 0x73, 0x79, 0x08, 0xee, 0xbc, 0x04, 0xb0, 0xe4, 0xfa, 0x02, 0x5c, 0x76,
	0xc0, 0x9e, 0x1d, 0x97, 0xa7, 0x40, 0x3e, 0x3e, 0x27, 0xb6, 0xdb, 0xee, 0xa2, 0xef, 0x9f, 0xf3,
	0xcf, 0x3f, 0x33, 0xb1, 0x3a, 0x3d, 0x49, 0x74, 0x5e, 0x30, 0x9d, 0xe7, 0xcc, 0xe4, 0xe6, 0xdc,
	0xa7, 0x15, 0xb5, 0x12, 0x25, 0xfa, 0x66, 0xc1, 0x84, 0xab, 0xf6, 0xe7, 0x9e, 0x49, 0x09, 0x4e,
	0xe0, 0xba, 0x5a, 0x8b, 0x12, 0x7d, 0x4b, 0x28, 0x1b, 0x9a, 0x59, 0x9c, 0xee, 0xd3, 0x5d, 0x08,
	0xb0, 0xab, 0xa0, 0x0d, 0x85, 0x21, 0xc4, 0x9e, 0x3a, 0x13, 0x25, 0x3a, 0xca, 0x68, 0x94, 0xd3,
	0x20, 0x33, 0x07, 0x5c, 0x16, 0x13, 0xf6, 0x15, 0x1e, 0xc6, 0xc2, 0x30, 0xc5, 0x2d, 0xd5, 0x8b,
	0x12, 0x3d, 0xa0, 0xdc, 0xc2, 0x5d, 0x4a, 0xc6, 0x94, 0x49, 0xf9, 0xc4, 0xe0, 0xb6, 0xea, 0x1f,
	0x27, 0x09, 0x03, 0x3b, 0x3b, 0xaf, 0x0d, 0x63, 0xc9, 0xcb, 0x37, 0x05, 0x6e, 0xa8, 0xf5, 0x23,
	0x5c, 0x18, 0xee, 0x39, 0xe1, 0x52, 0x2c, 0x62, 0xe6, 0x8b, 0x3a, 0xd8, 0x7d, 0xdc, 0x54, 0xdd,
	0xa3, 0x82, 0x30, 0x3c, 0x70, 0xfd, 0xf7, 0x4c, 0x1e, 0x4f, 0x8b, 0x96, 0x0e, 0x0f, 0xb1, 0x6b,
	0x17, 0xb2, 0xa3, 0x75, 0x15, 0xab, 0xb4, 0xfa, 0x1c, 0x60, 0xcf, 0x6e, 0xa4, 0x41, 0x85, 0xe1,
	0x4b, 0x80, 0x1b, 0x36, 0xec, 0x3e, 0x25, 0x66, 0x41, 0x75, 0xfd, 0xd7, 0x00, 0x37, 0x6d, 0xa8,
	0xb6, 0x20, 0x0c, 0xdf, 0x02, 0x5c, 0x53, 0x2a, 0x4a, 0xf4, 0x55, 0x53, 0xed, 0xfa, 0x7b, 0x80,
	0xa0, 0x56, 0x97, 0x40, 0x18, 0x7e, 0xf8, 0x92, 0xeb, 0x07, 0xb1, 0x1d, 0xfd, 0xa7, 0x2f, 0xa9,
	0x80, 0x30, 0xfc, 0xf2, 0x79, 0xaa, 0x09, 0xae, 0xa4, 0x8b, 0x38, 0x27, 0x0d, 0xbf, 0x03, 0xdc,
	0xb6, 0x3b, 0xaf, 0x70, 0xd5, 0xb6, 0xf4, 0x25, 0x0d, 0x7f, 0x02, 0x37, 0x58, 0xa5, 0x5d, 0x8b,
	0x27, 0x77, 0x48, 0xc3, 0x5f, 0x1f, 0xb4, 0xf9, 0x62, 0x48, 0xd3, 0x1c, 0xfe, 0x05, 0xd8, 0xb7,
	0xe7, 0xbe, 0x41, 0xa9, 0xde, 0x25, 0x91, 0xd1, 0x8c, 0xca, 0x30, 0xaf, 0x42, 0x37, 0x73, 0x8b,
	0x0b, 0xc3, 0xeb, 0xb0, 0x95, 0xc9, 0x49, 0xf0, 0x26, 0xf4, 0x3e, 0x45, 0x3a, 0x69, 0xf8, 0xbc,
	0x5d, 0xfa, 0x34, 0xb9, 0x30, 0xbc, 0x0b, 0x97, 0xbb, 0x9b, 0x8c, 0xe6, 0xf3, 0xc6, 0x93, 0xf7,
	0x21, 0x6e, 0xd9, 0x8b, 0x1e, 0x52, 0x84, 0xe1, 0x43, 0x88, 0x67, 0x1b, 0x27, 0x6d, 0x15, 0xc0,
	0xc7, 0xb0, 0xbe, 0xe9, 0xe5, 0x2c, 0xa6, 0x54, 0x97, 0x6e, 0x8f, 0x3a, 0xf5, 0x4d, 0x3d, 0x15,
	0x86, 0xc7, 0x1d, 0x97, 0x77, 0x67, 0x96, 0x11, 0xd5, 0xe5, 0x4f, 0x3a, 0x2e, 0x6f, 0x8b, 0x0b,
	0xc3, 0x53, 0xef, 0x33, 0xa0, 0xbc, 0xc2, 0xf6, 0x3f, 0xfe, 0xcc, 0xfb, 0x34, 0xb1, 0x30, 0x3c,
	0xf7, 0x3e, 0x55, 0xd2, 0x65, 0x73, 0x78, 0xd1, 0x71, 0xd3, 0x39, 0xa1, 0x6e, 0x03, 0x2f, 0x3b,
	0x17, 0x4f, 0xde, 0x0e, 0x79, 0x3c, 0x3e, 0x65, 0xbf, 0xe5, 0x0b, 0xff, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xdd, 0xb3, 0x01, 0x3c, 0xdc, 0x03, 0x00, 0x00,
}
