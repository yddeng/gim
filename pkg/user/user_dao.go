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
	var attr, convs []byte
	rows, err := db.DB().Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&u.ID, &u.CreateAt, &u.UpdateAt, &attr, &convs); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(attr, &u.Attrs)
	_ = json.Unmarshal(attr, &u.Convs)

	return &u, nil
}

func setNxUser(u *User) error {
	sqlStatement := `
INSERT INTO "users" (id,create_at,update_at,attr,convs)
VALUES($1, $2, $3, $4, $5) 
ON conflict(id) DO 
UPDATE SET create_at = $2, update_at = $3, attr = $4, convs = $5;`
	smt, err := db.DB().Prepare(sqlStatement)
	if err != nil {
		return err
	}

	attr, _ := json.Marshal(u.Attrs)
	conv, _ := json.Marshal(u.Convs)
	_, err = smt.Exec(u.ID, u.CreateAt, u.UpdateAt, attr, conv)
	return err
}
