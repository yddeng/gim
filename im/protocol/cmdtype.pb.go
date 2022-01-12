// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol/cmdtype.proto

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

type CmdType int32

const (
	CmdType_CmdNone                CmdType = 0
	CmdType_CmdUserLoginReq        CmdType = 1
	CmdType_CmdUserLoginResp       CmdType = 2
	CmdType_CmdHeartbeat           CmdType = 3
	CmdType_CmdGetUserInfoReq      CmdType = 4
	CmdType_CmdGetUserInfoResp     CmdType = 5
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
	CmdType_CmdDeleteFriendReq     CmdType = 409
	CmdType_CmdDeleteFriendResp    CmdType = 410
	CmdType_CmdNotifyDeleteFriend  CmdType = 411
	CmdType_CmdNotifyUserOnline    CmdType = 412
)

var CmdType_name = map[int32]string{
	0:   "CmdNone",
	1:   "CmdUserLoginReq",
	2:   "CmdUserLoginResp",
	3:   "CmdHeartbeat",
	4:   "CmdGetUserInfoReq",
	5:   "CmdGetUserInfoResp",
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
	409: "CmdDeleteFriendReq",
	410: "CmdDeleteFriendResp",
	411: "CmdNotifyDeleteFriend",
	412: "CmdNotifyUserOnline",
}

var CmdType_value = map[string]int32{
	"CmdNone":                0,
	"CmdUserLoginReq":        1,
	"CmdUserLoginResp":       2,
	"CmdHeartbeat":           3,
	"CmdGetUserInfoReq":      4,
	"CmdGetUserInfoResp":     5,
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
	"CmdDeleteFriendReq":     409,
	"CmdDeleteFriendResp":    410,
	"CmdNotifyDeleteFriend":  411,
	"CmdNotifyUserOnline":    412,
}

func (x CmdType) Enum() *CmdType {
	p := new(CmdType)
	*p = x
	return p
}

func (x CmdType) String() string {
	return proto.EnumName(CmdType_name, int32(x))
}

func (x *CmdType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CmdType_value, data, "CmdType")
	if err != nil {
		return err
	}
	*x = CmdType(value)
	return nil
}

func (CmdType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_469dd37e8cff621d, []int{0}
}

func init() {
	proto.RegisterEnum("message.IM.CmdType", CmdType_name, CmdType_value)
}

func init() { proto.RegisterFile("protocol/cmdtype.proto", fileDescriptor_469dd37e8cff621d) }

var fileDescriptor_469dd37e8cff621d = []byte{
	// 558 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x54, 0xc9, 0x72, 0xd3, 0x40,
	0x10, 0x45, 0x12, 0x14, 0xa9, 0x09, 0x55, 0x69, 0x26, 0xb6, 0xb3, 0xf0, 0x07, 0x1c, 0xc2, 0x37,
	0x04, 0xa5, 0x30, 0x06, 0x3b, 0x14, 0x01, 0x2e, 0xdc, 0x1c, 0x4f, 0xdb, 0xa5, 0x42, 0xd2, 0x34,
	0xea, 0x89, 0xab, 0xcc, 0xf2, 0x11, 0xec, 0xeb, 0x2f, 0x00, 0x47, 0xb6, 0x2f, 0x60, 0xb9, 0x01,
	0x77, 0x4e, 0x6c, 0x5f, 0x41, 0x69, 0x34, 0xb2, 0xa4, 0x24, 0x37, 0xd7, 0x7b, 0xaf, 0x5f, 0xbf,
	0xee, 0x1e, 0x4b, 0x74, 0x28, 0xd3, 0x46, 0x8f, 0x74, 0x7c, 0x66, 0x94, 0x28, 0x33, 0x23, 0xdc,
	0xb0, 0x80, 0x14, 0x09, 0x32, 0x0f, 0x27, 0xb8, 0xd1, 0x1b, 0x9c, 0xfe, 0xb5, 0x20, 0x8e, 0x87,
	0x89, 0xba, 0x3a, 0x23, 0x94, 0x8b, 0xf6, 0xe7, 0xb6, 0x4e, 0x11, 0x8e, 0xc8, 0x65, 0xb1, 0x14,
	0x26, 0xea, 0x1a, 0x63, 0xd6, 0xd7, 0x93, 0x28, 0xdd, 0xc1, 0x9b, 0xe0, 0xc9, 0x96, 0x80, 0x26,
	0xc8, 0x04, 0xbe, 0x04, 0x71, 0x22, 0x4c, 0xd4, 0x79, 0x1c, 0x66, 0x66, 0x17, 0x87, 0x06, 0x02,
	0xd9, 0x16, 0x27, 0xc3, 0x44, 0x75, 0xd1, 0xe4, 0xd2, 0x5e, 0x3a, 0xd6, 0x79, 0xf9, 0x51, 0xd9,
	0x11, 0x72, 0x3f, 0xcc, 0x04, 0xc7, 0x9c, 0x3c, 0xcc, 0x70, 0x68, 0xb0, 0x9b, 0xe9, 0x3d, 0xca,
	0xe5, 0xe8, 0xe4, 0x0d, 0x98, 0x09, 0xc6, 0x72, 0x4d, 0xb4, 0x0b, 0x1b, 0x0b, 0x0e, 0x30, 0xd9,
	0xc5, 0x8c, 0xf3, 0x12, 0x2d, 0xd7, 0x45, 0xe7, 0x30, 0x8a, 0x09, 0xa8, 0xea, 0x6e, 0xb9, 0x7e,
	0xc4, 0x26, 0xaf, 0x99, 0xc9, 0x15, 0xb1, 0x7c, 0x00, 0x67, 0x82, 0x5b, 0x8e, 0xd8, 0x8a, 0x98,
	0x75, 0x3c, 0xad, 0x82, 0xdd, 0x96, 0xab, 0xa2, 0x75, 0x90, 0x60, 0x82, 0x3b, 0xae, 0xff, 0xb6,
	0x36, 0xd1, 0x78, 0xd6, 0xe0, 0xe1, 0xae, 0x6c, 0xd9, 0x8d, 0x6e, 0x2a, 0x55, 0xc4, 0xca, 0xad,
	0x3e, 0x79, 0xb2, 0x6d, 0x57, 0x5a, 0x43, 0x99, 0xe0, 0xb3, 0x27, 0x57, 0x6c, 0xd8, 0x1d, 0x4c,
	0xf4, 0x14, 0x2b, 0xfd, 0x17, 0x4f, 0xae, 0xda, 0x50, 0x4d, 0x82, 0x09, 0xbe, 0x7a, 0x72, 0x49,
	0x88, 0x30, 0x51, 0x17, 0x74, 0x71, 0xac, 0x6f, 0x9e, 0x04, 0xb1, 0x38, 0x07, 0x98, 0xe0, 0x7b,
	0x29, 0xb9, 0xbc, 0x17, 0xd9, 0xd1, 0x7f, 0x94, 0x92, 0x02, 0x60, 0x82, 0x9f, 0x65, 0x9e, 0x62,
	0x82, 0x5e, 0x3a, 0x8d, 0x0c, 0x2a, 0xf8, 0xed, 0xc9, 0x75, 0xbb, 0xf3, 0x02, 0x2e, 0xda, 0xe6,
	0xbe, 0xa8, 0xe0, 0x8f, 0xe7, 0x06, 0x2b, 0xb8, 0x8b, 0xd1, 0xe8, 0x06, 0x2a, 0xf8, 0x5b, 0x06,
	0xad, 0x57, 0xf4, 0x71, 0x6c, 0xe0, 0x9f, 0x27, 0x3b, 0xf6, 0xdc, 0x57, 0x30, 0x55, 0x83, 0xe2,
	0x21, 0xe6, 0x61, 0x5e, 0xf9, 0x6e, 0xe6, 0x06, 0xce, 0x04, 0xaf, 0xfd, 0x46, 0x26, 0x47, 0xc1,
	0x1b, 0xbf, 0xf4, 0x99, 0xa5, 0xa3, 0x9a, 0xcf, 0xdb, 0xb9, 0x4f, 0x1d, 0x67, 0x82, 0x77, 0xfe,
	0x7c, 0x77, 0xa3, 0x61, 0x1c, 0xd7, 0x4a, 0xde, 0xfb, 0x72, 0xcd, 0x5e, 0x74, 0x1f, 0xc3, 0x04,
	0x1f, 0x7c, 0x79, 0xaa, 0x76, 0xd2, 0x86, 0x00, 0x3e, 0xfa, 0xd5, 0x4d, 0xcf, 0x65, 0x11, 0xa6,
	0x2a, 0x77, 0xbb, 0x17, 0x54, 0x37, 0x2d, 0x51, 0x26, 0xb8, 0x1f, 0xb8, 0xbc, 0x9b, 0x93, 0x0c,
	0xb1, 0x92, 0x3f, 0x08, 0x5c, 0xde, 0x06, 0xce, 0x04, 0x0f, 0x4b, 0x9f, 0x2e, 0x9a, 0x02, 0xb6,
	0x6f, 0xfc, 0x51, 0xe9, 0x53, 0x87, 0x99, 0xe0, 0x71, 0xe9, 0x53, 0x24, 0x9d, 0x37, 0x87, 0x27,
	0x81, 0x9b, 0xce, 0x11, 0x55, 0x1b, 0x78, 0x5a, 0xd6, 0x6c, 0x61, 0x8c, 0xa6, 0x16, 0xea, 0x59,
	0xe0, 0x76, 0xd5, 0x24, 0x98, 0xe0, 0x79, 0xd0, 0x78, 0x0a, 0x75, 0x1e, 0x5e, 0x04, 0x8d, 0xa3,
	0xe7, 0x7f, 0xf2, 0x4b, 0x69, 0x1c, 0xa5, 0x08, 0x2f, 0x83, 0xb3, 0xe2, 0xfa, 0x42, 0xf9, 0x39,
	0xfa, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x27, 0x33, 0x3b, 0x62, 0x99, 0x04, 0x00, 0x00,
}
