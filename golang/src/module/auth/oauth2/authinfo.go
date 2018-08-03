package oauth2

import (
	"time"
	"encoding/gob"
	"module/global"
	"public/token"
)

// Save callback get user info
type AuthInfo struct {
	Source 	int
	Id		string
	Name	string
	Email	string
	Uid 	int
}

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

// User info to jwt
func (au *AuthInfo) GetJwt() (string, error) {
	hmacSampleSecret := []byte("secret")
	token.AkSet("secret", hmacSampleSecret)
	to := token.NewWithClaims( token.SigningMethodHS256, &token.MapClaims{
		"source":	au.Source,
		"sourcename":	au.GetSource(),
		"id":		au.Id,
		"name":		au.Name,
		"uid":		au.Uid,
		"ak":		hmacSampleSecret,
		"expires":	time.Now().Add(1000 * time.Second).Unix(),
	})
	return to.SignedString(hmacSampleSecret)
}

// Get user auth type
func (au *AuthInfo) GetSource() string {
	return oauth2source[au.Source]
}
