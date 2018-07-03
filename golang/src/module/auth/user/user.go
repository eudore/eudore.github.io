package user;

import (
	"encoding/gob"
	"module/global"
	"module/tools"
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

func Reload() error{
	stmtQueryUserId		=	tools.Stmt(global.Sql,"SELECT UID FROM tb_auth_user_info WHERE Name=?;")
	stmtInsertUser		=	tools.Stmt(global.Sql,"INSERT INTO tb_auth_user_info(Name,LoginIP,SiginTime) VALUES(?,?,?);")
	stmtUpdateSignUp	=	tools.Stmt(global.Sql,"UPDATE tb_auth_oauth2_login SET UID=?,Stats=? WHERE Source=? AND OID=?;")
	return nil
}