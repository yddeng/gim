package conv

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"strings"
	"time"
)

func loadConversation(id uint64) (*Conversation, error) {
	sqlStr := `
SELECT * FROM "conversation_list" 
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, id)
	log.Debug(sqlStatement)

	var conv Conversation
	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(
		&conv.ID,
		&conv.Type,
		&conv.Name,
		&conv.Creator,
		&conv.CreateAt,
		&conv.LastMessageID,
		&conv.LastMessageAt); err != nil {
		return nil, err
	}

	return &conv, nil
}

func insertConversation(conv *Conversation) error {
	sqlStatement := `
INSERT INTO "conversation_list" (type,name,creator,create_at)  
VALUES ($1,$2,$3,$4)
RETURNING id;`

	return db.SqlDB.QueryRow(sqlStatement,
		int32(conv.Type),
		conv.Name,
		conv.Creator,
		conv.CreateAt).Scan(&conv.ID)
}

func updateConversation(conv *Conversation) error {
	sqlStr := `
UPDATE "conversation_list" 
SET name = $1, last_message_id = $2, last_message_at = $3
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, conv.ID)
	_, err := db.SqlDB.Exec(sqlStatement, conv.Name, conv.LastMessageID, conv.LastMessageAt)
	return err
}

/*  -- ----------------------------
-- conv_user
-- ---------------------------- */

func setNxConvUser(convID int64, userID string, role int) error {
	sqlStatement := `
INSERT INTO "conv_user" (id,conv_id, user_id, role)
VALUES($1, $2, $3, $4) 
ON conflict(id) DO 
UPDATE SET conv_id = $2, user_id = $3, role = $4;`
	smt, err := db.SqlDB.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%d_%s", convID, userID)
	_, err = smt.Exec(id, convID, userID, role)
	return err
}

func getUserConversations(userID string) (map[int64]int, error) {
	sqlStr := `
SELECT conv_id, role FROM "conv_user" 
WHERE user_id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, userID)
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	convs := map[int64]int{}
	defer rows.Close()
	for rows.Next() {
		var convID int64
		var role int
		err = rows.Scan(&convID, &role)
		if err != nil {
			return nil, err
		}
		convs[convID] = role
	}
	return convs, nil
}

func getConversationUsers(convID int64) (map[string]int, error) {
	sqlStr := `
SELECT user_id, role FROM "conv_user" 
WHERE conv_id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, convID)
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	users := map[string]int{}
	defer rows.Close()
	for rows.Next() {
		var userID string
		var role int
		err = rows.Scan(&userID, &role)
		if err != nil {
			return nil, err
		}
		users[userID] = role
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

func createMessageTable(tableName string) error {
	sqlStatement := fmt.Sprintf(createMessageTableStr, tableName, tableName)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func insertMessage(convID int64, msg *pb.MessageInfo, tableName string) error {
	//tableName := makeMessageTableName()
	//if tableName != nowMessageTable {
	//	nowMessageTable = tableName
	//	if err := createMessageTable(tableName); err != nil {
	//		return err
	//	}
	//}

	sqlStr := `
INSERT INTO "%s" (id,conv_id,message_id,message)  
VALUES ($1,$2,$3,$4);`

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
