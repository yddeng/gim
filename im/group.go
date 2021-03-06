package im

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/protocol"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
)

var groupCache *lru.Cache

type cacheGroup struct {
	g *Group
}

func GetGroup(groupID int64) *Group {
	v, ok := groupCache.Get(groupID)
	if ok {
		g := v.(*cacheGroup).g
		if g != nil && g.deleting {
			return nil
		}
		return g
	}

	if g, err := dbLoadGroup(groupID); err != nil {
		log.Error(err)
		return nil
	} else if g == nil {
		groupCache.Add(groupID, &cacheGroup{g: nil})
		return nil
	} else if members, err := dbGetGroupMembers(groupID); err != nil {
		log.Error(err)
		return nil
	} else {
		g.Members = members
		groupCache.Add(groupID, &cacheGroup{g: g})
		return g
	}
}

func addGroup(g *Group) {
	groupCache.Add(g.ID, &cacheGroup{g: g})
}

func removeGroup(groupID int64) {
	groupCache.Remove(groupID)
}

type Group struct {
	Type          protocol.GroupType // 对话类型
	ID            int64              // 全局唯一ID
	Creator       string             // 对话创建者
	CreateAt      int64              // 创建时间戳 秒
	Extra         []*protocol.Extra  // 附加属性
	LastMessageAt int64              // 最后一条消息的时间
	LastMessageID int64              // 最后一条消息的ID
	Members       map[string]*Member
	deleting      bool // 正在移除
}

func (this *Group) Pack() *protocol.Group {
	g := &protocol.Group{
		Type:          this.Type.Enum(),
		ID:            proto.Int64(this.ID),
		LastMessageAt: proto.Int64(this.LastMessageAt),
		LastMessageID: proto.Int64(this.LastMessageID),
	}

	return g
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
			NotifyUser(id, msg)
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
	//log.Debug(sqlStatement)

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
