package im

import (
	"fmt"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
)

var friendCache *lru.Cache

type cacheFriend struct {
	fs map[string]*Friend
}

func GetFriends(userID string) map[string]*Friend {
	v, ok := friendCache.Get(userID)
	if ok {
		c := v.(*cacheFriend)
		return c.fs
	}

	if fs, err := dbLoadFriends(userID); err != nil {
		log.Error(err)
		return nil
	} else {
		friendCache.Add(userID, &cacheFriend{fs: fs})
		return fs
	}
}

func addFriend(userID, friendID string, f *Friend) {
	c, ok := friendCache.Get(userID)
	if ok {
		cf := c.(*cacheFriend)
		cf.fs[friendID] = f
	}
}

func removeFriend(userID, friendID string) {
	c, ok := friendCache.Get(userID)
	if ok {
		cf := c.(*cacheFriend)
		delete(cf.fs, friendID)
	}
}

const (
	FriendStatusU1U2  = 1
	FriendStatusU2U1  = 2
	FriendStatusBoth  = 3
	FriendStatusAgree = 4
)

type Friend struct {
	ID       string
	UserID1  string
	UserID2  string
	CreateAt int64
	Status   int // 1:u1->u2 2:u2->u1 3:both 4:friend
}

func dbLoadFriends(userID string) (map[string]*Friend, error) {
	sqlStr := `
SELECT * FROM "friend" 
WHERE user1_id = '%s' OR user2_id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, userID, userID)
	//log.Debug(sqlStatement)

	rows, err := sqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fs := map[string]*Friend{}
	for rows.Next() {
		var f Friend
		if err := rows.Scan(
			&f.ID,
			&f.UserID1,
			&f.UserID2,
			&f.CreateAt,
			&f.Status); err != nil {
			return nil, err
		}

		if f.UserID1 == userID {
			fs[f.UserID2] = &f
		} else {
			fs[f.UserID1] = &f
		}
	}

	return fs, nil
}

func dbSetNxFriend(f *Friend) error {
	sqlStatement := `
INSERT INTO "friend" (id,user1_id,user2_id,create_at,status)
VALUES($1, $2, $3, $4,$5) 
ON conflict(id) DO 
UPDATE SET  create_at = $4, status = $5;`
	smt, err := sqlDB.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(f.ID, f.UserID1, f.UserID2, f.CreateAt, f.Status)
	return err
}

func dbDelFriend(id string) error {
	sqlStr := `
DELETE from "friend"
WHERE id = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, id)
	_, err := sqlDB.Exec(sqlStatement)
	return err
}
