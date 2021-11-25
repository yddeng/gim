package user

import (
	"encoding/json"
	"fmt"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/utils/log"
)

func loadUser(key string) (*User, error) {
	sqlStr := `
SELECT * FROM "users" 
WHERE id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, key)
	log.Debug(sqlStatement)

	var u User
	var extra []byte
	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&u.ID, &u.CreateAt, &u.UpdateAt, &extra); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(extra, &u.Extra)
	return &u, nil
}

func setNxUser(u *User) error {
	sqlStatement := `
INSERT INTO "users" (id,create_at,update_at,extra)
VALUES($1, $2, $3, $4) 
ON conflict(id) DO 
UPDATE SET create_at = $2, update_at = $3, extra = $4;`
	smt, err := db.SqlDB.Prepare(sqlStatement)
	if err != nil {
		return err
	}

	extra, _ := json.Marshal(u.Extra)
	_, err = smt.Exec(u.ID, u.CreateAt, u.UpdateAt, extra)
	return err
}
