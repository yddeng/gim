package conv

import (
	"encoding/json"
	"fmt"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/utils/log"
)

func loadConversation(id uint64) (*Conversation, error) {
	sqlStr := `
SELECT * FROM "conversation_list" 
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, id)
	log.Debug(sqlStatement)

	var conv Conversation
	var member []byte
	rows, err := db.DB().Query(sqlStatement)
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
		&conv.LastMessageAt,
		&member); err != nil {
		return nil, err
	}

	_ = json.Unmarshal(member, &conv.Members)

	return &conv, nil
}

func setConversation(conv *Conversation) error {
	sqlStatement := `
INSERT INTO "conversation_list" (type,name,creator,create_at,members)  
VALUES ($1,$2,$3,$4,$5)
RETURNING id;`

	member, _ := json.Marshal(conv.Members)
	return db.DB().QueryRow(sqlStatement,
		int32(conv.Type),
		conv.Name,
		conv.Creator,
		conv.CreateAt,
		member).Scan(&conv.ID)
}

func updateConversation(conv *Conversation) error {
	sqlStr := `
UPDATE "conversation_list" 
SET name = $1, last_message_id = $2, last_message_at = $3, members = $4
WHERE id = '%d';`

	sqlStatement := fmt.Sprintf(sqlStr, conv.ID)
	member, _ := json.Marshal(conv.Members)

	_, err := db.DB().Exec(sqlStatement, conv.Name, conv.LastMessageID, conv.LastMessageAt, member)
	return err
}

/*  -- ----------------------------
-- message
-- ---------------------------- */
