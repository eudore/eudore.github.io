package user

import (
	"time"
	"net/http"
	"public/log"
	"module/global"
	"module/tools"
	"module/auth/oauth2"
)

var (
	stmtQueryUserId		=	tools.Stmt(global.Sql,"SELECT UID FROM tb_auth_user_info WHERE Name=?;")
	stmtInsertUser		=	tools.Stmt(global.Sql,"INSERT INTO tb_auth_user_info(Name,LoginIP,SiginTime) VALUES(?,?,?);")
	stmtUpdateSignUp	=	tools.Stmt(global.Sql,"UPDATE tb_auth_oauth2_login SET UID=?,Stats=? WHERE Source=? AND OID=?;")
)

func Signup(w http.ResponseWriter, r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	a := sess.Get("oauth2")
	if a != nil {
		// oauth2 sign up
		a2 := a.(oauth2.AuthInfo)
		tools.Template(w,"auth/user/signup.html",map[string]interface{}{"AuthInfo": a2,"Source": a2.GetSource()})
	}else {
		// user sign up
	}
}

func SignupSubmit(w http.ResponseWriter, r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	name := r.URL.RawQuery
	// create user
	_,err := stmtInsertUser.Exec(name,tools.Ipint(tools.GetRealClientIP(r)),time.Now())	
	if err != nil {
		// have user
		log.Info(err)	
		log.Info(name,"存在")
		return
	}
	a := sess.Get("oauth2")
	if a != nil {
		// user oauth2 login
		a2 := a.(oauth2.AuthInfo)
		stmtQueryUserId.QueryRow(name).Scan(&a2.Uid)
		_,err := stmtUpdateSignUp.Exec(a2.Uid,0,a2.Source,a2.Id)
		log.Info(err)
	}else {
		// create user pass
		// user login
	}
	Login(w,r)
}
