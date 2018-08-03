package user

import (
	"fmt"
	"time"
	"net/http"
	"net/url"
	"encoding/gob"
	"public/log"
	"public/token"
	"module/global"
)

const (
	AuthToken		=	"authtoken"
	AuthRedirect	=	"redirect"
	refererHeader	=	"Referer"
	DefaultRedirect	=	"/"
	oauth2ClientId	=	"client_id"
	oauth2Code		=	"code"
	oauth2State		=	"state"
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
func (a *Authorize) Info() string {
	return ""
}


func Auth(w http.ResponseWriter, r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	// save redirect
	redirect := r.Header.Get("Referer")
	if redirect == ""  {
		redirect = DefaultRedirect
	}
	// is login
	user := sess.Get("user")	
	if user != nil {
		http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
		return
	}
	// oauth2 request
	r.ParseForm()
	var oauthinfo string
	if r.Form["client_id"] != nil {
		a := Authorize{
			Response_type:	r.Form["response_type"][0],
			Client_id:		r.Form["client_id"][0],
			Redirect_uri:	r.Form["redirect_uri"][0],
			Scope:			r.Form["scope"],
			State:			r.Form["state"][0],
		}
		if a.Check() {
			// load oauth2
			oauthinfo = r.Form["client_id"][0] + r.Form["state"][0] 
			log.Info(oauthinfo)
		} else {
			// invalid clientid error
			http.Error(w,http.StatusText(403),http.StatusUnauthorized)
		}
	}
	// login
	global.Template(w,"auth/user/auth.html",map[string]interface{}{"Redirect": redirect,"Oauth2": oauthinfo})
}

func Authpass(w http.ResponseWriter, r *http.Request) {
	// check pass
	r.ParseForm()
	var uid int
	login := r.Form.Get("login")
	pass := r.Form.Get("pass")
	log.Info(login, pass)
	err := global.Sql.QueryRow("SELECT UID FROM tb_auth_oauth2_pass WHERE Login=? AND Pass=?",login,pass).Scan(&uid)
	if err != nil || uid == 0 {
		log.Info(err)
		return
	}
	// oauth2 requert
	oauthinfo := r.Form.Get("oauth2")
	if len(oauthinfo) != 0 {
		// get code
		code := "code???"
		// callback
		log.Info(oauthinfo)
		http.Redirect(w, r, oauthinfo + code, http.StatusFound)
		return
	}
	// simple login
	hmacSampleSecret := []byte("secret")
	token.AkSet("secret", hmacSampleSecret)
	to := token.NewWithClaims( token.SigningMethodHS256, &token.MapClaims{
		"uid":		uid,
		"ak":		hmacSampleSecret,
		"expires":	time.Now().Add(100 * time.Second).Unix(),
	})
	t , _ := to.SignedString(hmacSampleSecret)
	r.Form.Set(AuthToken, t)
	Login(w,r)
}
