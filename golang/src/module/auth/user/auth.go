package user

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/gob"
	"public/log"
	"module/global"
	"module/tools"
)



func init() {
	gob.Register(Authorize{})
}

type Authorize struct {
	Name			string
	Response_type	string
	Client_id		string
	Redirect_uri	string
	Scope			[]string
	State			string
}

func (a *Authorize) Redirect(code string) string {
	arg := url.Values{
		"code":     {code},
		"state": 	{a.State},
	}
	return fmt.Sprintf("%s?%s",a.Redirect_uri,arg.Encode())
}

func (a *Authorize) Check() bool {	
	return global.Sql.QueryRow("SELECT Name FROM tb_auth_oauth2_app WHERE ClientID=? AND Callback=?;",a.Client_id,a.Redirect_uri).Scan(&a.Name) == nil
}

func Auth(w http.ResponseWriter, r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	defer sess.SessionRelease(w)
	// save redirect
	redirect := r.Header.Get("Referer")
	if redirect == ""  {
		redirect = "/"
	}
	sess.Set("redirect",redirect)
	// is login
	user := sess.Get("user")	
	if user != nil {
//		Login(w,r)
		http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
		return
	}
	// oauth2 request
	r.ParseForm()
	if r.Form["client_id"] != nil {
		a := Authorize{
			Response_type:	r.Form["response_type"][0],
			Client_id:		r.Form["client_id"][0],
			Redirect_uri:	r.Form["redirect_uri"][0],
			Scope:			r.Form["scope"],
			State:			r.Form["state"][0],
		}
		if a.Check() {
			// 
			sess.Set("authorize",a)
			log.Json(sess.Get("authorize"))
		} else {
			// invalid clientid error
		}
		// invalid clientid
	}
	// login
	tools.Template(w,"auth/user/auth.html",map[string]interface{}{})
}

func Authpass(w http.ResponseWriter, r *http.Request) {
	sess,err := global.Session.SessionStart(w, r)
	defer sess.SessionRelease(w)
	r.ParseForm()
	// check pass
	var uid int
	login := r.Form["login"][0]
	pass := r.Form["pass"][0]
	err = global.Sql.QueryRow("SELECT UID FROM tb_auth_oauth2_pass WHERE Login=? AND Pass=?",login,pass).Scan(&uid)
	if err != nil || uid == 0 {
		log.Info(err)
		return
	}
	// oauth2 requert
	a := sess.Get("authorize")
	if a != nil {
		// get code
		code := "code???"
		// callback
		au := a.(Authorize)
		log.Info(au.Redirect(code))
		http.Redirect(w, r, au.Redirect(code), http.StatusFound)
		return
	}
	// simple login
	sess.Set("user", *(NewUser(uid)))
	Login(w,r)
	// return code state
}
