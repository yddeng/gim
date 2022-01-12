package im

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/protocol"
	"github.com/yddeng/utils/log"
	"time"
)

func onAddFriend(u *User, msg *Message) {
	req := msg.GetData().(*protocol.AddFriendReq)
	log.Debugf("user(%s) onAddFriend %v", u.ID, req)

	friends := GetFriends(u.ID)
	addID := req.GetUserID()
	addUser := GetUser(addID)
	if addUser == nil {
		u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{Code: protocol.ErrCode_UserNotExist.Enum()})
		return
	}

	if f, ok := friends[addID]; ok {
		if f.Status == FriendStatusBoth ||
			(f.Status == FriendStatusU1U2 && u.ID == f.UserID1) ||
			(f.Status == FriendStatusU2U1 && u.ID == f.UserID2) {
			u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{})
		} else if f.Status == FriendStatusAgree {
			u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{Code: protocol.ErrCode_FriendAlreadyIsFriend.Enum()})
		} else {
			oldStatus := f.Status
			f.Status = FriendStatusBoth
			if err := WrapFunc(dbSetNxFriend)(func(err error) {
				if err != nil {
					f.Status = oldStatus
					log.Error(err)
					u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{Code: protocol.ErrCode_Error.Enum()})
					return
				}
				u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{})
				NotifyUser(addID, &protocol.NotifyAddFriend{UserID: proto.String(u.ID)})

			}, f); err != nil {
				log.Error(err)
				u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{Code: protocol.ErrCode_Busy.Enum()})
			}
		}
	} else {
		f = &Friend{CreateAt: time.Now().Unix()}
		if u.ID < addID {
			f.ID = fmt.Sprintf("%s_%s", u.ID, addID)
			f.UserID1 = u.ID
			f.UserID2 = addID
			f.Status = FriendStatusU1U2
		} else {
			f.ID = fmt.Sprintf("%s_%s", addID, u.ID)
			f.UserID1 = addID
			f.UserID2 = u.ID
			f.Status = FriendStatusU2U1
		}

		if err := WrapFunc(dbSetNxFriend)(func(err error) {
			if err != nil {
				log.Error(err)
				u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{Code: protocol.ErrCode_Error.Enum()})
				return
			}
			addFriend(u.ID, addID, f)
			addFriend(addID, u.ID, f)
			u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{})
			NotifyUser(addID, &protocol.NotifyAddFriend{UserID: proto.String(u.ID)})

		}, f); err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.AddFriendResp{Code: protocol.ErrCode_Busy.Enum()})
		}
	}
}

func onAgreeFriend(u *User, msg *Message) {
	req := msg.GetData().(*protocol.AgreeFriendReq)
	log.Debugf("user(%s) onAgreeFriend %v", u.ID, req)

	friends := GetFriends(u.ID)
	agreeID := req.GetUserID()
	isAgree := req.GetAgree()
	if f, ok := friends[agreeID]; ok {
		log.Debugf("user(%s) onAgreeFriend %v", u.ID, f)
		if f.Status == FriendStatusAgree {
			u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_FriendAlreadyIsFriend.Enum()})
		} else if (f.Status == FriendStatusU1U2 && u.ID == f.UserID1) ||
			(f.Status == FriendStatusU2U1 && u.ID == f.UserID2) {
			u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_FriendApplyClosed.Enum()})
		} else {
			if isAgree {
				oldStatus := f.Status
				f.Status = FriendStatusAgree
				if err := WrapFunc(dbSetNxFriend)(func(err error) {
					if err != nil {
						f.Status = oldStatus
						log.Error(err)
						u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_Error.Enum()})
						return
					}
					ff := GetFriends(agreeID)[u.ID]
					ff.Status = f.Status
					u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{})
					NotifyUser(agreeID, &protocol.NotifyAgreeFriend{UserID: proto.String(u.ID), Agree: proto.Bool(true)})

				}, f); err != nil {
					log.Error(err)
					u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_Busy.Enum()})
				}
			} else {
				if err := WrapFunc(dbDelFriend)(func(err error) {
					if err != nil {
						log.Error(err)
						u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_Error.Enum()})
						return
					}
					removeFriend(u.ID, agreeID)
					removeFriend(agreeID, u.ID)
					u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{})
					NotifyUser(agreeID, &protocol.NotifyAgreeFriend{UserID: proto.String(u.ID), Agree: proto.Bool(false)})

				}, f.ID); err != nil {
					log.Error(err)
					u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_Busy.Enum()})
				}
			}

		}
	} else {
		u.SendToClient(msg.GetSeq(), &protocol.AgreeFriendResp{Code: protocol.ErrCode_FriendApplyClosed.Enum()})
	}
}

func onDeleteFriend(u *User, msg *Message) {
	req := msg.GetData().(*protocol.DeleteFriendReq)
	log.Debugf("user(%s) onDeleteFriend %v", u.ID, req)

	friends := GetFriends(u.ID)
	deleteID := req.GetUserID()

	if f, ok := friends[deleteID]; ok {
		if err := WrapFunc(dbDelFriend)(func(err error) {
			if err != nil {
				log.Error(err)
				u.SendToClient(msg.GetSeq(), &protocol.DeleteFriendResp{Code: protocol.ErrCode_Error.Enum()})
				return
			}
			removeFriend(u.ID, deleteID)
			removeFriend(deleteID, u.ID)
			if f.Status == FriendStatusAgree {
				NotifyUser(deleteID, &protocol.NotifyDeleteFriend{UserID: proto.String(u.ID)})
			}

			u.SendToClient(msg.GetSeq(), &protocol.DeleteFriendResp{})

		}, f.ID); err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &protocol.DeleteFriendResp{Code: protocol.ErrCode_Busy.Enum()})
		}
	}
}

func onGetFriends(u *User, msg *Message) {
	//req := msg.GetData().(*protocol.GetFriendsReq)
	log.Debugf("user(%s) onGetFriends", u.ID)

	friends := GetFriends(u.ID)
	if len(friends) == 0 {
		u.SendToClient(msg.GetSeq(), &protocol.GetFriendsResp{})
		return
	}

	resp := &protocol.GetFriendsResp{
		Friends: make([]*protocol.Friend, 0, len(friends)),
	}
	for friendID, f := range friends {
		status := protocol.FriendStatus_Agree
		if f.Status == FriendStatusAgree {
		} else if f.Status == FriendStatusBoth ||
			(f.Status == FriendStatusU1U2 && friendID == f.UserID1) ||
			(f.Status == FriendStatusU2U1 && friendID == f.UserID2) {
			status = protocol.FriendStatus_Apply
		} else {
			continue
		}
		friend := &protocol.Friend{UserID: proto.String(friendID), Status: status.Enum()}
		resp.Friends = append(resp.Friends, friend)
	}

	u.SendToClient(msg.GetSeq(), resp)
}

func init() {
	registerHandler(uint16(protocol.CmdType_CmdAddFriendReq), onAddFriend)
	registerHandler(uint16(protocol.CmdType_CmdAgreeFriendReq), onAgreeFriend)
	registerHandler(uint16(protocol.CmdType_CmdGetFriendsReq), onGetFriends)
	registerHandler(uint16(protocol.CmdType_CmdDeleteFriendReq), onDeleteFriend)
}
