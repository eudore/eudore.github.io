package oauth2

import (
	"errors"
	"encoding/gob"
	"module/global"
)

type AuthInfo struct {
	Source 	int
	Id		string
	Name	string
	Email	string
	Uid 	int
}

var (
	ErrStats 		=	 errors.New("stats invalid")
)

func init() {
	gob.Register(AuthInfo{})
}

func (au *AuthInfo) getuid() (err error) {
	var stats int = -1
	err = global.Sql.QueryRow("SELECT UID,Stats FROM tb_auth_oauth2_login WHERE Source=? AND OID=?;",au.Source,au.Id).Scan(&au.Uid,&stats)
	if err != nil {
		if stmt, err := global.Sql.Prepare("INSERT tb_auth_oauth2_login(Source,Name,Email,OID,Stats) VALUES(?,?,?,?,?);");err==nil{
			_, err = stmt.Exec(au.Source,au.Name,au.Email,au.Id,1)	
		}
		return
	}
	return
}

func (au *AuthInfo) GetSource() string {
	return oauth2source[au.Source]
}