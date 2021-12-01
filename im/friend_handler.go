package im

import (
	"fmt"
	"github.com/yddeng/gim/im/pb"
	"github.com/yddeng/utils/log"
	"time"
)

func onAddFriend(u *User, msg *Message) {
	req := msg.GetData().(*pb.AddFriendReq)
	log.Debugf("user(%s) onAddFriend %v", u.ID, req)

	friends := GetFriends(u.ID)
	addID := req.GetUserID()
	addUser := GetUser(addID)
	if addUser == nil {
		u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{Code: pb.ErrCode_UserNotExist})
		return
	}

	if f, ok := friends[addID]; ok {
		if f.Status == FriendStatusBoth ||
			(f.Status == FriendStatusU1U2 && u.ID == f.UserID1) ||
			(f.Status == FriendStatusU2U1 && u.ID == f.UserID2) {
			u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{})
		} else if f.Status == FriendStatusAgree {
			u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{Code: pb.ErrCode_FriendAlreadyIsFriend})
		} else {
			oldStatus := f.Status
			f.Status = FriendStatusBoth
			if err := WrapFunc(dbSetNxFriend)(func(err error) {
				if err != nil {
					f.Status = oldStatus
					log.Error(err)
					u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{Code: pb.ErrCode_Error})
					return
				}
				u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{})
				addUser.SendToClient(0, &pb.NotifyAddFriend{UserID: u.ID})

			}, f); err != nil {
				log.Error(err)
				u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{Code: pb.ErrCode_Error})
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
				u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{Code: pb.ErrCode_Error})
				return
			}
			addFriend(u.ID, f)
			u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{})
			addUser.SendToClient(0, &pb.NotifyAddFriend{UserID: u.ID})

		}, f); err != nil {
			log.Error(err)
			u.SendToClient(msg.GetSeq(), &pb.AddFriendResp{Code: pb.ErrCode_Error})
		}
	}
}

func onAgreeFriend(u *User, msg *Message) {
	req := msg.GetData().(*pb.AgreeFriendReq)
	log.Debugf("user(%s) onAgreeFriend %v", u.ID, req)

	friends := GetFriends(u.ID)
	agreeID := req.GetUserID()
	isAgree := req.GetAgree()
	if f, ok := friends[agreeID]; ok {
		if f.Status == FriendStatusAgree {
			u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_FriendAlreadyIsFriend})
		} else if (f.Status == FriendStatusU1U2 && u.ID == f.UserID1) ||
			(f.Status == FriendStatusU2U1 && u.ID == f.UserID2) {
			u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_FriendApplyClosed})
		} else {
			if isAgree {
				oldStatus := f.Status
				f.Status = FriendStatusAgree
				if err := WrapFunc(dbSetNxFriend)(func(err error) {
					if err != nil {
						f.Status = oldStatus
						log.Error(err)
						u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_Error})
						return
					}
					u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{})
					GetUser(agreeID).SendToClient(0, &pb.NotifyAgreeFriend{UserID: u.ID, Agree: true})

				}, f); err != nil {
					log.Error(err)
					u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_Error})
				}
			} else {
				if err := WrapFunc(dbDelFriend)(func(err error) {
					if err != nil {
						log.Error(err)
						u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_Error})
						return
					}
					removeFriend(u.ID, agreeID)
					u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{})
					GetUser(agreeID).SendToClient(0, &pb.NotifyAgreeFriend{UserID: u.ID, Agree: false})

				}, f.ID); err != nil {
					log.Error(err)
					u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_Error})
				}
			}

		}
	} else {
		u.SendToClient(msg.GetSeq(), &pb.AgreeFriendResp{Code: pb.ErrCode_FriendApplyClosed})
	}
}

func onGetFriends(u *User, msg *Message) {
	//req := msg.GetData().(*pb.GetFriendsReq)
	//log.Debugf("user(%s) onGetFriends %v", u.ID, req)

	friends := GetFriends(u.ID)
	if len(friends) == 0 {
		u.SendToClient(msg.GetSeq(), &pb.GetFriendsResp{})
		return
	}

	resp := &pb.GetFriendsResp{
		Friends: make([]*pb.Friend, 0, len(friends)),
	}
	for friendID, f := range friends {
		status := pb.FriendStatus_Agree
		if f.Status == FriendStatusAgree {
		} else if f.Status == FriendStatusBoth ||
			(f.Status == FriendStatusU1U2 && friendID == f.UserID1) ||
			(f.Status == FriendStatusU2U1 && friendID == f.UserID2) {
			status = pb.FriendStatus_Apply
		} else {
			continue
		}
		friend := &pb.Friend{UserID: friendID, Status: status}
		resp.Friends = append(resp.Friends, friend)
	}

	u.SendToClient(msg.GetSeq(), resp)
}

func init() {
	registerHandler(uint16(pb.CmdType_CmdAddFriendReq), onAddFriend)
	registerHandler(uint16(pb.CmdType_CmdAgreeFriendReq), onAgreeFriend)
	registerHandler(uint16(pb.CmdType_CmdGetFriendsReq), onGetFriends)
}
