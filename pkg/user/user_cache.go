package user

import "github.com/yddeng/utils/log"

var (
	userMap = map[string]*User{}
)

type cacheUser struct {
	u *User
}

func GetUser(id string) *User {
	u, ok := userMap[id]
	if !ok {
		var err error
		u, err = loadUser(id)
		if err != nil {
			log.Error(err)
			return nil
		}
		if u != nil {
			userMap[id] = u
		}
	}
	return u
}
