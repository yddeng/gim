package conv

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"strings"
	"time"
)

func selectConversation(id uint64) (*Conversation, error) {
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

/*  -- ----------------------------
-- conv_member
-- ---------------------------- */

func setNxConvUser(cmember []*CMember) error {
	sqlStr := `
INSERT INTO "conv_member" (id, conv_id, user_id, nickname, create_at, mute, role)
VALUES %s
ON conflict(id) DO 
UPDATE SET nickname = excluded.nickname, mute = excluded.mute, role = excluded.role ;`

	values := make([]string, 0, len(cmember))
	for _, v := range cmember {
		values = append(values, fmt.Sprintf("('%s',%d,'%s','%s',%d,%d,%d)",
			v.ID, v.ConvID, v.UserID, v.Nickname, v.CreateAt, v.Mute, v.Role))
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(values, ","))
	log.Debug(sqlStatement)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func delConvUser(cmember []*CMember) error {
	sqlStr := `
DELETE FROM "conv_member" 
WHERE %s;`

	keys := make([]string, 0, len(cmember))
	for _, v := range cmember {
		keys = append(keys, fmt.Sprintf("id = '%s'", v.ID))
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(keys, " OR "))
	log.Debug(sqlStatement)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func getUserConversations(userID string) (map[int64]*CMember, error) {
	sqlStr := `
SELECT * FROM "conv_member" 
WHERE user_id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, userID)
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	convs := map[int64]*CMember{}
	defer rows.Close()
	for rows.Next() {
		var cm CMember
		if err = rows.Scan(
			&cm.ID,
			&cm.ConvID,
			&cm.UserID,
			&cm.Nickname,
			&cm.CreateAt,
			&cm.Mute,
			&cm.Role); err != nil {
			return nil, err
		}
		convs[cm.ConvID] = &cm
	}
	return convs, nil
}

func getConversationUsers(convID int64) (map[string]*CMember, error) {
	sqlStr := `
SELECT * FROM "conv_member" 
WHERE conv_id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, convID)
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	users := map[string]*CMember{}
	defer rows.Close()
	for rows.Next() {
		var cm CMember
		if err = rows.Scan(
			&cm.ID,
			&cm.ConvID,
			&cm.UserID,
			&cm.Nickname,
			&cm.CreateAt,
			&cm.Mute,
			&cm.Role); err != nil {
			return nil, err
		}
		users[cm.UserID] = &cm
	}
	return users, nil
}

/*  -- ----------------------------
-- message
-- ---------------------------- */

var (
	createMessageTableStr = `DROP TABLE IF EXISTS "%s";
CREATE TABLE "%s" (
    "id"           varchar(255) NOT NULL,
    "conv_id"      int8 NOT NULL ,
    "message_id"   int8 NOT NULL,
    "message"      bytea NOT NULL,
    PRIMARY KEY ("id")
);`
)

func makeMessageTableName() string {
	date := time.Now().Format("20060102")
	return "message_" + date
}

func existMessageTable(tableName string) bool {
	sqlStr := `
select count(*) from "%s";`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	smt, err := db.SqlDB.Prepare(sqlStatement)
	if err != nil {
		return false
	}
	row := smt.QueryRow()
	var count int
	err = row.Scan(&count)
	if err != nil {
		return false
	}
	return true
}

func createMessageTable(tableName string) error {
	sqlStatement := fmt.Sprintf(createMessageTableStr, tableName, tableName)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func setNxMessage(convID int64, msg *pb.MessageInfo, tableName string) error {
	sqlStr := `
INSERT INTO "%s" (id,conv_id,message_id,message)  
VALUES ($1,$2,$3,$4)
ON conflict(id) DO 
UPDATE SET message = $4;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	id := fmt.Sprintf("%d_%d", convID, msg.GetMsgID())
	data, _ := proto.Marshal(msg)
	_, err := db.SqlDB.Exec(sqlStatement, id, convID, msg.GetMsgID(), data)
	return err
}

func loadMessageBatch(convID int64, start, limit int, tableName string) ([]*pb.MessageInfo, error) {
	sqlStr := `
SELECT message FROM "%s" 
WHERE %s;`

	keys := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		seq := start + i
		keys = append(keys, fmt.Sprintf("id = '%d_%d'", convID, seq))
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, " OR "))
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	infos := make([]*pb.MessageInfo, 0, limit)
	defer rows.Close()
	for rows.Next() {
		var info pb.MessageInfo
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}

		_ = proto.Unmarshal(data, &info)
		infos = append(infos, &info)
	}

	return infos, nil
}
