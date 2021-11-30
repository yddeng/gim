package im

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/pb"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
)

var groupCache *lru.Cache

type cacheGroup struct {
	c *Group
}

func GetGroup(groupID int64) *Group {
	v, ok := groupCache.Get(groupID)
	if ok {
		c := v.(*cacheGroup)
		return c.c
	}

	if c, err := dbLoadGroup(groupID); err != nil {
		log.Error(err)
		return nil
	} else if c == nil {
		groupCache.Add(groupID, &cacheGroup{c: nil})
		return nil
	} else if members, err := dbGetGroupMembers(groupID); err != nil {
		log.Error(err)
		return nil
	} else {
		c.Members = members
		groupCache.Add(groupID, &cacheGroup{c: c})
		return c
	}
}

func addGroup(c *Group) {
	groupCache.Add(c.ID, &cacheGroup{c: c})
}

func removeGroup(groupID int64) {
	groupCache.Remove(groupID)
}

type Group struct {
	Type          pb.GroupType      // 对话类型
	ID            int64             // 全局唯一ID
	Creator       string            // 对话创建者
	CreateAt      int64             // 创建时间戳 秒
	Extra         map[string]string // 附加属性
	LastMessageAt int64             // 最后一条消息的时间
	LastMessageID int64             // 最后一条消息的ID
	Members       map[string]*Member
}

func (this *Group) Pack() *pb.Group {
	c := &pb.Group{
		Type:          this.Type,
		ID:            this.ID,
		LastMessageAt: this.LastMessageAt,
		LastMessageID: this.LastMessageID,
	}

	return c
}

func (this *Group) Broadcast(msg proto.Message, except ...string) {
	for id := range this.Members {
		has := false
		for _, id2 := range except {
			if id == id2 {
				has = true
				break
			}
		}
		if !has {
			u := GetUser(id)
			if u != nil {
				u.SendToClient(0, msg)
			}
		}
	}
}

func (this *Group) AddMember(members []*Member) {
	for _, m := range members {
		this.Members[m.UserID] = m
	}
}

func (this *Group) RemoveMember(members []*Member) {
	for _, m := range members {
		delete(this.Members, m.UserID)
	}
}

func dbLoadGroup(id int64) (*Group, error) {
	sqlStr := `
SELECT * FROM "groups" 
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, id)
	log.Debug(sqlStatement)

	rows, err := sqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	var group Group
	var extra []byte
	if err := rows.Scan(
		&group.ID,
		&group.Type,
		&group.Creator,
		&group.CreateAt,
		&extra,
		&group.LastMessageID,
		&group.LastMessageAt); err != nil {
		return nil, err
	}

	_ = json.Unmarshal(extra, &group.Extra)
	return &group, nil
}

func dbInsertGroup(group *Group) error {
	sqlStatement := `
INSERT INTO "groups" (type,creator,create_at,extra)  
VALUES ($1,$2,$3,$4)
RETURNING id;`

	extra, _ := json.Marshal(group.Extra)
	return sqlDB.QueryRow(sqlStatement,
		int32(group.Type),
		group.Creator,
		group.CreateAt,
		extra).Scan(&group.ID)
}

func dbUpdateGroup(group *Group) error {
	sqlStr := `
UPDATE "groups" 
SET extra = $1, last_message_id = $2, last_message_at = $3
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, group.ID)
	extra, _ := json.Marshal(group.Extra)
	_, err := sqlDB.Exec(sqlStatement, extra, group.LastMessageID, group.LastMessageAt)
	return err
}

func dbDeleteGroup(groupID int64) error {
	sqlStr := `
DELETE from "groups"
WHERE id = %d;`

	sqlStatement := fmt.Sprintf(sqlStr, groupID)
	_, err := sqlDB.Exec(sqlStatement)
	return err
}
