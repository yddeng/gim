package im

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
)

var convCache *lru.Cache = lru.New(5000)

type cacheConv struct {
	c *Conversation
}

func GetConversation(convID int64) *Conversation {
	v, ok := convCache.Get(convID)
	if ok {
		c := v.(*cacheConv)
		return c.c
	}

	if c, err := loadConversation(convID); err != nil {
		log.Error(err)
		return nil
	} else if c == nil {
		convCache.Add(convID, &cacheConv{c: nil})
		return nil
	} else if users, err := getConversationUsers(convID); err != nil {
		log.Error(err)
		return nil
	} else {
		c.Members = users
		convCache.Add(convID, &cacheConv{c: c})
		return c
	}
}

func addConversation(c *Conversation) {
	convCache.Add(c.ID, &cacheConv{c: c})
}

type Conversation struct {
	Type          pb.ConversationType // 对话类型
	ID            int64               // 全局唯一ID
	Creator       string              // 对话创建者
	CreateAt      int64               // 创建时间戳 秒
	Extra         map[string]string   // 附加属性
	LastMessageAt int64               // 最后一条消息的时间
	LastMessageID int64               // 最后一条消息的ID
	LastMessage   *pb.MessageInfo     // 最后一条消息
	Message       []*pb.MessageInfo
	Members       map[string]*Member
}

func (this *Conversation) Pack() *pb.Conversation {
	c := &pb.Conversation{
		Type:          this.Type,
		ID:            this.ID,
		LastMessageAt: this.LastMessageAt,
		LastMessageID: this.LastMessageID,
	}

	return c
}

func (this *Conversation) Broadcast(msg proto.Message, except ...string) {
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

func (this *Conversation) AddMember(members []*Member) {
	for _, m := range members {
		this.Members[m.UserID] = m
	}
}

func (this *Conversation) RemoveMember(members []*Member) {
	for _, m := range members {
		delete(this.Members, m.UserID)
	}
}

func loadConversation(id int64) (*Conversation, error) {
	sqlStr := `
SELECT * FROM "conversation_list" 
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, id)
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	var conv Conversation
	var extra []byte
	if err := rows.Scan(
		&conv.ID,
		&conv.Type,
		&conv.Creator,
		&conv.CreateAt,
		&extra,
		&conv.LastMessageID,
		&conv.LastMessageAt); err != nil {
		return nil, err
	}

	_ = json.Unmarshal(extra, &conv.Extra)
	return &conv, nil
}

func insertConversation(conv *Conversation) error {
	sqlStatement := `
INSERT INTO "conversation_list" (type,creator,create_at,extra)  
VALUES ($1,$2,$3,$4)
RETURNING id;`

	extra, _ := json.Marshal(conv.Extra)
	return db.SqlDB.QueryRow(sqlStatement,
		int32(conv.Type),
		conv.Creator,
		conv.CreateAt,
		extra).Scan(&conv.ID)
}

func updateConversation(conv *Conversation) error {
	sqlStr := `
UPDATE "conversation_list" 
SET extra = $1, last_message_id = $2, last_message_at = $3
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, conv.ID)
	extra, _ := json.Marshal(conv.Extra)
	_, err := db.SqlDB.Exec(sqlStatement, extra, conv.LastMessageID, conv.LastMessageAt)
	return err
}
