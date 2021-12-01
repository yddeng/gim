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
	CmdType_CmdHeartbeat           CmdType = 3
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
	3:   "CmdHeartbeat",
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
	"CmdHeartbeat":           3,
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
	// 504 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0xc9, 0x8e, 0xd3, 0x4e,
	0x10, 0xc6, 0xff, 0xb6, 0xff, 0x62, 0xa4, 0x1e, 0xd0, 0x14, 0x35, 0x49, 0x66, 0xe1, 0x0d, 0x38,
	0x70, 0xe1, 0x09, 0x06, 0x23, 0xc2, 0x92, 0x19, 0x89, 0x01, 0x2e, 0xdc, 0x92, 0x74, 0x25, 0xb2,
	0x88, 0xed, 0xc2, 0xd5, 0x13, 0x29, 0x2c, 0x0f, 0xc1, 0xce, 0x4b, 0x00, 0x47, 0xb6, 0x27, 0x60,
	0x11, 0x17, 0xe0, 0xce, 0x8d, 0xe5, 0x29, 0x90, 0xdb, 0xdd, 0xb1, 0x3d, 0xc3, 0x2d, 0xfa, 0x7d,
	0x55, 0x5f, 0x7f, 0x55, 0x15, 0xab, 0x13, 0xe3, 0x54, 0x9b, 0x05, 0xd3, 0x19, 0x2e, 0x72, 0x93,
	0x9f, 0xfe, 0xb2, 0xa2, 0x56, 0xe2, 0x54, 0x5f, 0x5f, 0x30, 0xe1, 0xaa, 0xfd, 0xb9, 0x97, 0x67,
	0x04, 0xff, 0xe1, 0xba, 0x5a, 0x8b, 0x53, 0x7d, 0x43, 0xa8, 0x18, 0xe4, 0xd3, 0x24, 0xdb, 0xa7,
	0xdb, 0x10, 0x60, 0x47, 0x41, 0x1b, 0x0a, 0x43, 0x88, 0xa0, 0x8e, 0xc7, 0xa9, 0xbe, 0x48, 0xc3,
	0xc2, 0x8c, 0x68, 0x68, 0x20, 0xc2, 0xae, 0x3a, 0x19, 0xa7, 0x3a, 0x2e, 0x68, 0x68, 0xa8, 0x5f,
	0xe4, 0x07, 0x5c, 0xb6, 0x13, 0xf6, 0x14, 0x1e, 0xc6, 0xc2, 0x30, 0xc1, 0x2d, 0xd5, 0x8d, 0x53,
	0xdd, 0x27, 0x63, 0xe1, 0x2e, 0xa5, 0x23, 0x2a, 0xa4, 0x6c, 0xc9, 0x71, 0x5b, 0xf5, 0xfe, 0x25,
	0x09, 0x03, 0x3b, 0x3b, 0xaf, 0x0d, 0x12, 0x31, 0x65, 0xcf, 0x02, 0x37, 0xd4, 0xfa, 0x11, 0x2e,
	0x0c, 0x77, 0x9c, 0x70, 0x3e, 0x11, 0xc9, 0x67, 0xf3, 0x3a, 0xd8, 0x5d, 0xdc, 0x54, 0x9d, 0xa3,
	0x82, 0x30, 0xdc, 0x73, 0xef, 0xef, 0xe5, 0x26, 0x99, 0x2c, 0x5a, 0x3a, 0xdc, 0xc7, 0x8e, 0x5d,
	0xd1, 0x8e, 0xd6, 0x55, 0xac, 0xd2, 0xea, 0x43, 0x80, 0x5d, 0xbb, 0xa3, 0x06, 0x15, 0x86, 0x8f,
	0x01, 0x6e, 0xd8, 0xb0, 0xfb, 0x94, 0xe6, 0x73, 0xaa, 0xeb, 0x3f, 0x05, 0xb8, 0x69, 0x43, 0xb5,
	0x05, 0x61, 0xf8, 0x1c, 0xe0, 0x9a, 0x52, 0x71, 0xaa, 0x2f, 0xe7, 0xd5, 0xf6, 0xbf, 0x06, 0x08,
	0x6a, 0x75, 0x09, 0x84, 0xe1, 0x9b, 0x2f, 0xb9, 0x7a, 0x90, 0xd8, 0xd1, 0xbf, 0xfb, 0x92, 0x0a,
	0x08, 0xc3, 0x0f, 0x9f, 0xa7, 0x9a, 0xe0, 0x52, 0x36, 0x4f, 0x0c, 0x69, 0xf8, 0x19, 0xe0, 0xb6,
	0xdd, 0x79, 0x85, 0xab, 0x67, 0x4b, 0x5f, 0xd2, 0xf0, 0x2b, 0x70, 0x83, 0x55, 0xda, 0x95, 0x64,
	0x7c, 0x8b, 0x34, 0xfc, 0xf6, 0x41, 0x9b, 0x1d, 0x03, 0x9a, 0x18, 0xf8, 0x13, 0x60, 0xcf, 0x9e,
	0xfb, 0x1a, 0x65, 0x7a, 0x97, 0x44, 0x86, 0x53, 0x2a, 0xc3, 0xbc, 0x08, 0xdd, 0xcc, 0x2d, 0x2e,
	0x0c, 0x2f, 0xc3, 0x56, 0x26, 0x27, 0xc1, 0xab, 0xd0, 0xfb, 0x2c, 0xb2, 0x71, 0xc3, 0xe7, 0xf5,
	0xd2, 0xa7, 0xc9, 0x85, 0xe1, 0x4d, 0xb8, 0xdc, 0xdd, 0x78, 0x38, 0x9b, 0x35, 0x5a, 0xde, 0x86,
	0xb8, 0x65, 0x2f, 0x7a, 0x48, 0x11, 0x86, 0x77, 0x21, 0x9e, 0x6a, 0x9c, 0xb4, 0x55, 0x00, 0xef,
	0xc3, 0xfa, 0xa6, 0x17, 0x8a, 0x84, 0x32, 0x5d, 0xba, 0x3d, 0x88, 0xea, 0x9b, 0x7a, 0x2a, 0x0c,
	0x0f, 0x23, 0x97, 0x77, 0x67, 0x5a, 0x10, 0xd5, 0xe5, 0x8f, 0x22, 0x97, 0xb7, 0xc5, 0x85, 0xe1,
	0xb1, 0xf7, 0xe9, 0x93, 0xa9, 0xb0, 0xfd, 0x8f, 0x3f, 0xf1, 0x3e, 0x4d, 0x2c, 0x0c, 0x4f, 0xbd,
	0x4f, 0x95, 0x74, 0xf9, 0x38, 0x3c, 0x8b, 0xdc, 0x74, 0x4e, 0xa8, 0x9f, 0x81, 0xe7, 0xd1, 0xb9,
	0xff, 0x6f, 0x86, 0x3c, 0x1a, 0x1d, 0xb3, 0x5f, 0xf7, 0xd9, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xae, 0x1c, 0x3e, 0x18, 0xee, 0x03, 0x00, 0x00,
}
