package im

import (
	"fmt"
	"github.com/yddeng/gim/im/pb"
	"strings"
)

type Member struct {
	ID       string
	GroupID  int64
	UserID   string
	Nickname string
	CreateAt int64
	UpdateAt int64
	Mute     int // 禁言
	Role     int // 会话角色
}

func (this *Member) Pack() *pb.Member {
	msg := &pb.Member{
		UserID:   this.UserID,
		Nickname: this.Nickname,
		Mute:     this.Mute == 1,
		Role:     int32(this.Role),
		CreateAt: this.CreateAt,
		UpdateAt: this.UpdateAt,
	}

	u := GetUser(this.UserID)
	if u != nil && u.online() {
		msg.Online = true
	}

	return msg
}

func dbSetNxGroupMember(cmember []*Member) error {
	sqlStr := `
INSERT INTO "group_member" (id, group_id, user_id, nickname, create_at, update_at, mute, role)
VALUES %s
ON conflict(id) DO 
UPDATE SET nickname = excluded.nickname, update_at = excluded.update_at, mute = excluded.mute, role = excluded.role ;`

	values := make([]string, 0, len(cmember))
	for _, v := range cmember {
		values = append(values, fmt.Sprintf("('%s',%d,'%s','%s',%d,%d,%d,%d)",
			v.ID, v.GroupID, v.UserID, v.Nickname, v.CreateAt, v.UpdateAt, v.Mute, v.Role))
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(values, ","))
	//log.Debug(sqlStatement)
	_, err := sqlDB.Exec(sqlStatement)
	return err
}

func dbDelGroupMember(cmember []*Member) error {
	sqlStr := `
DELETE FROM "group_member" 
WHERE %s;`

	keys := make([]string, 0, len(cmember))
	for _, v := range cmember {
		keys = append(keys, fmt.Sprintf("id = '%s'", v.ID))
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(keys, " OR "))
	//log.Debug(sqlStatement)
	_, err := sqlDB.Exec(sqlStatement)
	return err
}

func dbGetUserGroups(userID string) (map[int64]*Member, error) {
	sqlStr := `
SELECT * FROM "group_member" 
WHERE user_id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, userID)
	//log.Debug(sqlStatement)

	rows, err := sqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	groups := map[int64]*Member{}
	defer rows.Close()
	for rows.Next() {
		var cm Member
		if err = rows.Scan(
			&cm.ID,
			&cm.GroupID,
			&cm.UserID,
			&cm.Nickname,
			&cm.CreateAt,
			&cm.UpdateAt,
			&cm.Mute,
			&cm.Role); err != nil {
			return nil, err
		}
		groups[cm.GroupID] = &cm
	}
	return groups, nil
}

func dbGetGroupMembers(groupID int64) (map[string]*Member, error) {
	sqlStr := `
SELECT * FROM "group_member" 
WHERE group_id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, groupID)
	//log.Debug(sqlStatement)

	rows, err := sqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	users := map[string]*Member{}
	defer rows.Close()
	for rows.Next() {
		var cm Member
		if err = rows.Scan(
			&cm.ID,
			&cm.GroupID,
			&cm.UserID,
			&cm.Nickname,
			&cm.CreateAt,
			&cm.UpdateAt,
			&cm.Mute,
			&cm.Role); err != nil {
			return nil, err
		}
		users[cm.UserID] = &cm
	}
	return users, nil
}
