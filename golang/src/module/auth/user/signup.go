package user

import (
	"time"
	"net/http"
	"public/log"
	"module/global"
	"module/tools"
	"module/auth/oauth2"
)


func Signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	authtoken := r.Form.Get(AuthToken)
	redirect := r.Form.Get(AuthRedirect)
	log.Info(authtoken)
	//if authtoken != "" {
		// oauth2 sign up
		err := global.Template(w,"auth/user/signup.html",map[string]interface{}{"Redirect": redirect, "Source": oauth2.GetSourceName(2)})
		log.Info(err)
	//}else {
		// user sign up
	//}
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
