package user;

import (
	"encoding/gob"
	"database/sql"
	"module/global"
)

var (
	stmtQueryUserId			*sql.Stmt
	stmtQueryUserInfo		*sql.Stmt
	stmtQueryUserPolicy		*sql.Stmt
	stmtInsertUser			*sql.Stmt
	stmtUpdateSignUp		*sql.Stmt
)

type User struct {
	Uid		int
	Name	string
	Status	int
	Level	int
	Policy	string
}

func NewUser(uid int) *User {
	var u User = User{Uid: uid}
	global.Sql.QueryRow("SELECT Name,Status,Level FROM tb_auth_user_info WHERE UID=?",uid).Scan(&u.Name,&u.Status,&u.Level)
	u.LoadPolicy()
	return &u
}

func (u *User) LoadPolicy() error {
	return stmtQueryUserPolicy.QueryRow(u.Uid).Scan(&u.Policy)
}

func Reload() error{
	stmtQueryUserId		=	global.Stmt("SELECT UID FROM tb_auth_user_info WHERE Name=?;")
	stmtQueryUserInfo	=	global.Stmt("SELECT Name,Status,Level FROM tb_auth_user_info WHERE UID=?;")
	stmtQueryUserPolicy	=	global.Stmt("SELECT GROUP_CONCAT(PID ORDER BY `Index`) FROM tb_auth_ram_bind WHERE UID=?;")
	stmtInsertUser		=	global.Stmt("INSERT INTO tb_auth_user_info(Name,LoginIP,SiginTime) VALUES(?,?,?);")
	stmtUpdateSignUp	=	global.Stmt("UPDATE tb_auth_oauth2_login SET UID=?,Stats=? WHERE Source=? AND OID=?;")
	return nil
}

func init() {
	gob.Register(User{})
}