package user;

import (
	"encoding/gob"
	"module/global"
)

type User struct {
	Uid		int
	Name	string
	Status	int
	Level	int
}

func init() {
	gob.Register(User{})
}

func NewUser(uid int) *User {
	var u User = User{Uid: uid}
	global.Sql.QueryRow("SELECT Name,Status,Level FROM tb_auth_user_info WHERE UID=?",uid).Scan(&u.Name,&u.Status,&u.Level)
	return &u
}
